package script

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"github.com/dop251/goja"
	"github.com/ecumeurs/upsiloncli/internal/dto"
	"github.com/ecumeurs/upsiloncli/internal/endpoint"
)

func (a *Agent) bindJSAPI() {
	upsilonObj := map[string]interface{}{
		"call":         a.jsCall,
		"waitForEvent": a.jsWaitForEvent,
		"getContext":   a.jsGetContext,
		"setContext":   a.jsSetContext,
		"log":          a.jsLog,

		// New lifecycle and assertion methods
		"onTeardown":   a.jsOnTeardown,
		"assert":       a.jsAssert,

		// New Shared State Methods
		"setShared": a.jsSetShared,
		"getShared": a.jsGetShared,

		// New Flow Control Methods
		"sleep": a.jsSleep,

		// Pathfinding
		"findPath":         a.jsFindPath,
		"planTravelToward": a.jsPlanTravelToward,

		// Environment
		"getEnv": a.jsGetEnv,
		
		// Tactical Helpers
		"myPlayer":            a.jsMyPlayer,
		"currentPlayer":       a.jsCurrentPlayer,
		"currentCharacter":    a.jsCurrentCharacter,
		"myCharacters":        a.jsMyCharacters,
		"myAllies":            a.jsMyAllies,
		"myAlliesCharacters":  a.jsMyAlliesCharacters,
		"myFoes":              a.jsMyFoes,
		"myFoesCharacters":    a.jsMyFoesCharacters,
		"cellContentAt":       a.jsCellContentAt,
	}
	a.VM.Set("upsilon", upsilonObj)
}

func (a *Agent) jsLog(msg interface{}) {
	a.Display.Print(fmt.Sprintf("%v", msg))
}

func (a *Agent) jsCall(routeName string, params map[string]interface{}) (interface{}, error) {
	ep := a.Registry.Get(routeName)
	if ep == nil {
		return nil, fmt.Errorf("unknown route: %s", routeName)
	}

	// Convert JS params to string map expected by endpoint.Execute
	inputs := make(map[string]string)
	for k, v := range params {
		inputs[k] = fmt.Sprintf("%v", v)
	}

	resp, err := ep.ExecuteRaw(a.Client, a.Session, inputs)
	if err != nil {
		return nil, err
	}

	// Capture session state (tokens, IDs) from response
	endpoint.SyncSession(resp, a.Session)

	// Ensure WebSockets are synced if auth happened (token might have been set)
	a.Listener.Sync()

	return resp.Data, nil
}

func (a *Agent) jsGetContext(key string) string {
	val, _ := a.Session.Get(key)
	return val
}

func (a *Agent) jsSetContext(key, value string) {
	a.Session.Set(key, value)
}

func (a *Agent) jsWaitForEvent(eventName string, timeoutMs int) (interface{}, error) {
	return a.Listener.WaitForData(a.Ctx, eventName, timeoutMs)
}

// jsOnTeardown stores a JS callback to be executed later
func (a *Agent) jsOnTeardown(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) > 0 {
		if fn, ok := goja.AssertFunction(call.Arguments[0]); ok {
			a.TeardownHook = fn
		}
	}
	return goja.Undefined()
}

// jsAssert throws a JS exception if the condition is false
func (a *Agent) jsAssert(condition bool, msg string) {
	if !condition {
		// Panic inside a Goja bridged function causes a catchable JS exception
		panic(a.VM.ToValue(fmt.Sprintf("Assertion Failed: %s", msg)))
	}
}

func (a *Agent) jsSetShared(key string, value interface{}) {
	a.Shared.Set(key, value)
}

func (a *Agent) jsGetShared(key string) interface{} {
	val, ok := a.Shared.Get(key)
	if !ok {
		return nil
	}
	return val
}

// jsSleep pauses the current agent's goroutine without affecting others.
func (a *Agent) jsSleep(ms int) {
	select {
	case <-a.Ctx.Done():
		return
	case <-time.After(time.Duration(ms) * time.Millisecond):
	}
}

func (a *Agent) flattenEntities(board *dto.BoardState) []dto.Entity {
	var all []dto.Entity
	if board == nil { return all }
	for _, p := range board.Players {
		all = append(all, p.Entities...)
	}
	return all
}

func (a *Agent) jsFindPath(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 3 {
		return a.VM.ToValue(nil)
	}

	var start, end dto.Position
	var board dto.BoardState

	// Marshal/Unmarshal is the most reliable way to convert deep JS objects to Go DTOs
	startBytes, _ := json.Marshal(call.Arguments[0].Export())
	json.Unmarshal(startBytes, &start)

	endBytes, _ := json.Marshal(call.Arguments[1].Export())
	json.Unmarshal(endBytes, &end)

	boardBytes, _ := json.Marshal(call.Arguments[2].Export())
	json.Unmarshal(boardBytes, &board)
	
	// Inject flattened entities for pathfinding algorithms that expect them
	board.Entities = a.flattenEntities(&board)

	path := FindPath(&board, start, end)
	
	// Ensure proper JSON mapping for the return value
	var result interface{}
	pathBytes, _ := json.Marshal(path)
	json.Unmarshal(pathBytes, &result)

	return a.VM.ToValue(result)
}

// @spec-link [[api_plan_travel_toward]]
func (a *Agent) jsPlanTravelToward(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 3 {
		return a.VM.ToValue(nil)
	}

	entityID := call.Arguments[0].String()

	var target dto.Position
	targetBytes, _ := json.Marshal(call.Arguments[1].Export())
	json.Unmarshal(targetBytes, &target)

	var board dto.BoardState
	boardBytes, _ := json.Marshal(call.Arguments[2].Export())
	json.Unmarshal(boardBytes, &board)

	// Inject flattened entities
	board.Entities = a.flattenEntities(&board)

	path := PlanTravelToward(&board, entityID, target)
	
	// Ensure proper JSON mapping for the return value
	var result interface{}
	pathBytes, _ := json.Marshal(path)
	json.Unmarshal(pathBytes, &result)

	return a.VM.ToValue(result)
}

func (a *Agent) jsGetEnv(key string) string {
	return os.Getenv(key)
}

// --- Tactical Utility Implementations ---

func (a *Agent) jsMyPlayer() interface{} {
	parts := a.Session.Participants()
	for _, p := range parts {
		if p.IsSelf {
			return p
		}
	}
	return nil
}

func (a *Agent) jsCurrentPlayer() interface{} {
	board := a.Session.LastBoard()
	if board == nil { return nil }
	
	if board.CurrentPlayerIsSelf {
		return a.jsMyPlayer()
	}

	// Find owner of the current entity
	for _, p := range board.Players {
		for _, e := range p.Entities {
			if e.ID == board.CurrentEntityID {
				return p
			}
		}
	}

	return nil
}

func (a *Agent) jsCurrentCharacter() interface{} {
	board := a.Session.LastBoard()
	if board == nil || board.Players == nil { return nil }
	for _, p := range board.Players {
		for _, e := range p.Entities {
			if e.ID == board.CurrentEntityID {
				return e
			}
		}
	}
	return nil
}

func (a *Agent) jsMyCharacters() []dto.Entity {
	board := a.Session.LastBoard()
	if board == nil { return nil }
	var mine []dto.Entity
	for _, p := range board.Players {
		if p.IsSelf {
			mine = append(mine, p.Entities...)
		}
	}
	return mine
}

func (a *Agent) jsMyAllies() []dto.Player {
	board := a.Session.LastBoard()
	if board == nil { return nil }

	var allies []dto.Player
	var myTeam int
	found := false
	for _, p := range board.Players {
		if p.IsSelf {
			myTeam = p.Team
			found = true
			break
		}
	}

	if !found { return nil }
	
	for _, p := range board.Players {
		if p.Team == myTeam && !p.IsSelf {
			allies = append(allies, p)
		}
	}
	return allies
}

func (a *Agent) jsMyAlliesCharacters() []dto.Entity {
	allies := a.jsMyAllies()
	var all []dto.Entity
	for _, p := range allies {
		all = append(all, p.Entities...)
	}
	return all
}

func (a *Agent) jsMyFoes() []dto.Player {
	board := a.Session.LastBoard()
	if board == nil { return nil }

	var foes []dto.Player
	var myTeam int
	found := false
	for _, p := range board.Players {
		if p.IsSelf {
			myTeam = p.Team
			found = true
			break
		}
	}

	if !found { return nil }
	
	for _, p := range board.Players {
		if p.Team != myTeam {
			foes = append(foes, p)
		}
	}
	return foes
}

func (a *Agent) jsMyFoesCharacters() []dto.Entity {
	foes := a.jsMyFoes()
	var all []dto.Entity
	for _, p := range foes {
		all = append(all, p.Entities...)
	}
	return all
}

func (a *Agent) jsCellContentAt(x, y int) interface{} {
	board := a.Session.LastBoard()
	if board == nil { return nil }
	
	if y < 0 || y >= len(board.Grid.Cells) || x < 0 || x >= len(board.Grid.Cells[0]) {
		return nil
	}
	
	cell := board.Grid.Cells[y][x]
	var foundEntity *dto.Entity
	if cell.EntityID != "" {
		for _, p := range board.Players {
			for _, e := range p.Entities {
				if e.ID == cell.EntityID {
					foundEntity = &e
					break
				}
			}
			if foundEntity != nil { break }
		}
	}
	
	return map[string]interface{}{
		"obstacle": cell.Obstacle,
		"entity":   foundEntity,
	}
}
