package api

import (
	"time"

	"github.com/ecumeurs/upsilonapi/stdmessage"
	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/ecumeurs/upsilonmapdata/grid/cell"
	"github.com/ecumeurs/upsilonmapdata/grid/position"
	"github.com/google/uuid"
)

// @spec-link [[api_go_battle_engine]]

type ArenaActionResponse struct {
	Status string `json:"status"`
}

type ArenaStartResponse struct {
	ArenaID      string     `json:"arena_id"`
	InitialState BoardState `json:"initial_state"`
}

type ActiveMatchStatsResponse struct {
	ActiveCount int `json:"active_count"`
}

// @spec-link [[entity_grid]]

type Cell struct {
	EntityID string `json:"entity_id"` // if any
	Obstacle bool   `json:"obstacle"`  // if any
}

// Grid: A 2D array of cells; for our purpose as in this implementation, the height will be fixed at 1 for every cell giving us a flat map.
type Grid struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Cells  [][]Cell `json:"cells"` // Cells are stored in width-major order.
}

type Turn struct {
	PlayerID string `json:"player_id"`
	Delay    int    `json:"delay"`
	EntityID string `json:"entity_id"`
}

type BoardState struct {
	Entities        []Entity  `json:"entities"`
	Grid            Grid      `json:"grid"`
	Turn            []Turn    `json:"turn"`
	CurrentPlayerID string    `json:"current_player_id"`
	CurrentEntityID string    `json:"current_entity_id"`
	Timeout         time.Time `json:"timeout"` // End of turn date.
	StartTime       time.Time `json:"start_time"`
	WinnerID        string    `json:"winner_id"` // if any, the game is done; based on player id.
	Players         []Player  `json:"players"`   // Full roster
}

// ArenaEvent is the payload for the webhook
type ArenaEvent struct {
	MatchID   string     `json:"match_id"`   // targeted match
	EventType string     `json:"event_type"` // Board State Change, Turn Started, Battle Start, Battle End
	PlayerID  string     `json:"player_id"`  // if set, targeted player
	EntityID  string     `json:"entity_id"`  // if set, targeted entity
	Data      BoardState `json:"data"`       // event specific data (board change)
	Timeout   time.Time  `json:"timeout"`    // End of turn date.
}

type ArenaActionResponseMessage = stdmessage.StandardMessage[ArenaActionResponse, stdmessage.MetaNil]
type ArenaStartResponseMessage = stdmessage.StandardMessage[ArenaStartResponse, stdmessage.MetaNil]
type ArenaEventMessage = stdmessage.StandardMessage[ArenaEvent, stdmessage.MetaNil]

// NewError creates a new StandardMessage with the given error.
func NewError(requestId string, err string) stdmessage.StandardMessage[stdmessage.DataNil, stdmessage.MetaNil] {
	return stdmessage.StandardMessage[stdmessage.DataNil, stdmessage.MetaNil]{
		RequestID: requestId,
		Message:   err,
		Meta:      stdmessage.MetaNil{},
		Success:   false,
		Data:      stdmessage.DataNil{},
	}
}

// NewSuccess creates a new StandardMessage with the given data.
func NewSuccess[T any](requestId string, msg string, data T) stdmessage.StandardMessage[T, stdmessage.MetaNil] {
	return stdmessage.StandardMessage[T, stdmessage.MetaNil]{
		RequestID: requestId,
		Message:   msg,
		Meta:      stdmessage.MetaNil{},
		Success:   true,
		Data:      data,
	}
}

// NewEntity creates a new Entity from the given entity (upsilonbattle's)
func NewEntity(entity entity.Entity) Entity {
	return Entity{
		ID:       entity.ID.String(),
		PlayerID: entity.ControllerID.String(),
		Name:     entity.Name,
		HP:       entity.GetPropertyC("HP").GetValue(),
		MaxHP:    entity.GetPropertyC("HP").GetMaxValue(),
		Attack:   entity.GetPropertyI("Attack").I(),
		Defense:  entity.GetPropertyI("Defense").I(),
		Move:     entity.GetPropertyC("Movement").GetValue(),
		MaxMove:  entity.GetPropertyC("Movement").GetMaxValue(),
		Position: Position{X: entity.Position.X, Y: entity.Position.Y},
	}
}

// NewBoardState creates a new BoardState DTO from internal state.
func NewBoardState(matchID uuid.UUID, g *grid.Grid, entities []entity.Entity, players []Player, ts turner.TurnState, startTime time.Time, timeout time.Time) BoardState {
	bs := BoardState{
		StartTime:       startTime,
		Timeout:         timeout,
		CurrentEntityID: ts.CurrentEntityTurn.String(),
		Players:         players,
	}

	// Map Grid
	bs.Grid = Grid{
		Width:  g.Width,
		Height: g.Length,
		Cells:  make([][]Cell, g.Width),
	}
	for x := 0; x < g.Width; x++ {
		bs.Grid.Cells[x] = make([]Cell, g.Length)
		for y := 0; y < g.Length; y++ {
			z := g.TopMostCellAt(x, y)
			cl, ok := g.CellAt(position.New(x, y, z))
			if ok {
				bs.Grid.Cells[x][y] = Cell{
					EntityID: cl.EntityID.String(),
					Obstacle: cl.Type == cell.Obstacle,
				}
				if cl.EntityID == uuid.Nil {
					bs.Grid.Cells[x][y].EntityID = ""
				}
			}
		}
	}

	entityToPlayer := make(map[uuid.UUID]string)
	for _, e := range entities {
		entityToPlayer[e.ID] = e.ControllerID.String()
		bs.Entities = append(bs.Entities, NewEntity(e))

		if e.ID == ts.CurrentEntityTurn {
			bs.CurrentPlayerID = e.ControllerID.String()
		}
	}

	for _, t := range ts.RemainingTurns {
		bs.Turn = append(bs.Turn, Turn{
			EntityID: t.EntityId.String(),
			PlayerID: entityToPlayer[t.EntityId],
			Delay:    t.Delay,
		})
	}

	return bs
}
