package endpoint

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ecumeurs/upsiloncli/internal/api"
	"github.com/ecumeurs/upsiloncli/internal/session"
	"github.com/ecumeurs/upsiloncli/internal/dto"
)

// --- Profile ---

// ProfileGet implements Endpoint for GET /api/v1/profile.
type ProfileGet struct{}

func (e *ProfileGet) Name() string        { return "profile_get" }
func (e *ProfileGet) Description() string { return "Get player profile with stats" }
func (e *ProfileGet) Method() string      { return "GET" }
func (e *ProfileGet) Path() string        { return "/api/v1/profile" }
func (e *ProfileGet) Auth() bool          { return true }
func (e *ProfileGet) Params() []Param     { return nil }

func (e *ProfileGet) Next() []string {
	return []string{"profile_characters"}
}

func (e *ProfileGet) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Get(e.Path())
	if err != nil {
		return err
	}
	// Extract character IDs from response data for context
	extractCharacterIDs(resp, sess)
	return nil
}

// ProfileCharacters implements Endpoint for GET /api/v1/profile/characters.
type ProfileCharacters struct{}

func (e *ProfileCharacters) Name() string        { return "profile_characters" }
func (e *ProfileCharacters) Description() string { return "List character roster" }
func (e *ProfileCharacters) Method() string      { return "GET" }
func (e *ProfileCharacters) Path() string        { return "/api/v1/profile/characters" }
func (e *ProfileCharacters) Auth() bool          { return true }
func (e *ProfileCharacters) Params() []Param     { return nil }

func (e *ProfileCharacters) Next() []string {
	return []string{"profile_character", "character_reroll", "character_upgrade", "character_rename", "matchmaking_join"}
}

func (e *ProfileCharacters) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Get(e.Path())
	if err != nil {
		return err
	}
	// Store first character_id in context for convenience
	if data, ok := resp.Data.([]interface{}); ok && len(data) > 0 {
		if char, ok := data[0].(map[string]interface{}); ok {
			if id, ok := char["id"].(string); ok {
				sess.Set("character_id", id)
			}
		}
		// Also store all character IDs
		var ids []string
		for _, item := range data {
			if char, ok := item.(map[string]interface{}); ok {
				if id, ok := char["id"].(string); ok {
					ids = append(ids, id)
				}
			}
		}
		if len(ids) > 0 {
			sess.Set("character_ids", strings.Join(ids, ","))
		}
	}
	return nil
}

// ProfileCharacter implements Endpoint for GET /api/v1/profile/character/{characterId}.
type ProfileCharacter struct{}

func (e *ProfileCharacter) Name() string        { return "profile_character" }
func (e *ProfileCharacter) Description() string { return "Get specific character details" }
func (e *ProfileCharacter) Method() string      { return "GET" }
func (e *ProfileCharacter) Path() string        { return "/api/v1/profile/character/{characterId}" }
func (e *ProfileCharacter) Auth() bool          { return true }
func (e *ProfileCharacter) Params() []Param {
	return []Param{
		{Name: "characterId", Hint: "character UUID", Required: true, ContextKey: "character_id"},
	}
}

func (e *ProfileCharacter) Next() []string {
	return []string{"profile_characters", "character_upgrade"}
}

func (e *ProfileCharacter) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	path := strings.ReplaceAll(e.Path(), "{characterId}", inputs["characterId"])
	_, err := client.Get(path)
	return err
}

// CharacterReroll implements Endpoint for POST /api/v1/profile/character/{characterId}/reroll.
type CharacterReroll struct{}

func (e *CharacterReroll) Name() string        { return "character_reroll" }
func (e *CharacterReroll) Description() string { return "Reroll character stats (new accounts)" }
func (e *CharacterReroll) Method() string      { return "POST" }
func (e *CharacterReroll) Path() string        { return "/api/v1/profile/character/{characterId}/reroll" }
func (e *CharacterReroll) Auth() bool          { return true }
func (e *CharacterReroll) Params() []Param {
	return []Param{
		{Name: "characterId", Hint: "character UUID", Required: true, ContextKey: "character_id"},
	}
}

func (e *CharacterReroll) Next() []string {
	return []string{"profile_character"}
}

func (e *CharacterReroll) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	path := strings.ReplaceAll(e.Path(), "{characterId}", inputs["characterId"])
	_, err := client.Post(path, nil)
	return err
}

// CharacterUpgrade implements Endpoint for POST /api/v1/profile/character/{characterId}/upgrade.
type CharacterUpgrade struct{}

func (e *CharacterUpgrade) Name() string    { return "character_upgrade" }
func (e *CharacterUpgrade) Description() string { return "Allocate stat points from wins" }
func (e *CharacterUpgrade) Method() string  { return "POST" }
func (e *CharacterUpgrade) Path() string    { return "/api/v1/profile/character/{characterId}/upgrade" }
func (e *CharacterUpgrade) Auth() bool      { return true }
func (e *CharacterUpgrade) Params() []Param {
	return []Param{
		{Name: "characterId", Hint: "character UUID", Required: true, ContextKey: "character_id"},
		{Name: "hp", Hint: "HP increment (int)"},
		{Name: "attack", Hint: "Attack increment (int)"},
		{Name: "defense", Hint: "Defense increment (int)"},
		{Name: "movement", Hint: "Movement increment (int)"},
	}
}

func (e *CharacterUpgrade) Next() []string {
	return []string{"profile_character"}
}

func (e *CharacterUpgrade) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	path := strings.ReplaceAll(e.Path(), "{characterId}", inputs["characterId"])
	stats := make(map[string]string)
	for _, k := range []string{"hp", "attack", "defense", "movement"} {
		if v := inputs[k]; v != "" {
			stats[k] = v
		}
	}
	body := map[string]interface{}{"stats": stats}
	_, err := client.Post(path, body)
	return err
}

// CharacterRename implements Endpoint for POST /api/v1/profile/character/{characterId}/rename.
type CharacterRename struct{}

func (e *CharacterRename) Name() string        { return "character_rename" }
func (e *CharacterRename) Description() string { return "Rename a character" }
func (e *CharacterRename) Method() string      { return "POST" }
func (e *CharacterRename) Path() string        { return "/api/v1/profile/character/{characterId}/rename" }
func (e *CharacterRename) Auth() bool          { return true }
func (e *CharacterRename) Params() []Param {
	return []Param{
		{Name: "characterId", Hint: "character UUID", Required: true, ContextKey: "character_id"},
		{Name: "name", Hint: "new character name", Required: true},
	}
}

func (e *CharacterRename) Next() []string {
	return []string{"profile_character"}
}

func (e *CharacterRename) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	path := strings.ReplaceAll(e.Path(), "{characterId}", inputs["characterId"])
	_, err := client.Post(path, map[string]string{"name": inputs["name"]})
	return err
}

// CharacterDelete implements Endpoint for DELETE /api/v1/profile/character/{characterId}.
type CharacterDelete struct{}

func (e *CharacterDelete) Name() string        { return "character_delete" }
func (e *CharacterDelete) Description() string { return "Delete a character" }
func (e *CharacterDelete) Method() string      { return "DELETE" }
func (e *CharacterDelete) Path() string        { return "/api/v1/profile/character/{characterId}" }
func (e *CharacterDelete) Auth() bool          { return true }
func (e *CharacterDelete) Params() []Param {
	return []Param{
		{Name: "characterId", Hint: "character UUID", Required: true, ContextKey: "character_id"},
	}
}

func (e *CharacterDelete) Next() []string {
	return []string{"profile_characters"}
}

func (e *CharacterDelete) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	path := strings.ReplaceAll(e.Path(), "{characterId}", inputs["characterId"])
	_, err := client.Delete(path)
	return err
}

// --- Matchmaking ---

// MatchmakingJoin implements Endpoint for POST /api/v1/matchmaking/join.
type MatchmakingJoin struct{}

func (e *MatchmakingJoin) Name() string        { return "matchmaking_join" }
func (e *MatchmakingJoin) Description() string { return "Enter the matchmaking queue" }
func (e *MatchmakingJoin) Method() string      { return "POST" }
func (e *MatchmakingJoin) Path() string        { return "/api/v1/matchmaking/join" }
func (e *MatchmakingJoin) Auth() bool          { return true }
func (e *MatchmakingJoin) Params() []Param {
	return []Param{
		{Name: "game_mode", Hint: "1v1_PVP|1v1_PVE|2v2_PVP|2v2_PVE", Required: true},
	}
}

func (e *MatchmakingJoin) Next() []string {
	return []string{"matchmaking_status", "game_state"}
}

func (e *MatchmakingJoin) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Post(e.Path(), map[string]string{
		"game_mode": inputs["game_mode"],
	})
	if err != nil {
		return err
	}
	// Capture match_id if matched immediately
	if data, ok := resp.Data.(map[string]interface{}); ok {
		if matchID, ok := data["match_id"].(string); ok && matchID != "" {
			sess.Set("match_id", matchID)
			client.Printer.System(fmt.Sprintf("Match found! match_id: %s", matchID))
		}
		if status, ok := data["status"].(string); ok {
			sess.Set("matchmaking_status", status)
		}
	}
	return nil
}

// MatchmakingStatus implements Endpoint for GET /api/v1/matchmaking/status.
type MatchmakingStatus struct{}

func (e *MatchmakingStatus) Name() string        { return "matchmaking_status" }
func (e *MatchmakingStatus) Description() string { return "Poll matchmaking queue status" }
func (e *MatchmakingStatus) Method() string      { return "GET" }
func (e *MatchmakingStatus) Path() string        { return "/api/v1/matchmaking/status" }
func (e *MatchmakingStatus) Auth() bool          { return true }
func (e *MatchmakingStatus) Params() []Param     { return nil }

func (e *MatchmakingStatus) Next() []string {
	return []string{"matchmaking_join", "game_state", "matchmaking_leave"}
}

func (e *MatchmakingStatus) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Get(e.Path())
	if err != nil {
		return err
	}
	if data, ok := resp.Data.(map[string]interface{}); ok {
		if matchID, ok := data["match_id"].(string); ok && matchID != "" {
			sess.Set("match_id", matchID)
		}
		if status, ok := data["status"].(string); ok {
			sess.Set("matchmaking_status", status)
		}
	}
	return nil
}

// MatchmakingLeave implements Endpoint for DELETE /api/v1/matchmaking/leave.
type MatchmakingLeave struct{}

func (e *MatchmakingLeave) Name() string        { return "matchmaking_leave" }
func (e *MatchmakingLeave) Description() string { return "Leave the matchmaking queue" }
func (e *MatchmakingLeave) Method() string      { return "DELETE" }
func (e *MatchmakingLeave) Path() string        { return "/api/v1/matchmaking/leave" }
func (e *MatchmakingLeave) Auth() bool          { return true }
func (e *MatchmakingLeave) Params() []Param     { return nil }

func (e *MatchmakingLeave) Next() []string {
	return []string{"matchmaking_join", "profile_get"}
}

func (e *MatchmakingLeave) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	resp, err := client.Delete(e.Path())
	if err != nil {
		return err
	}
	sess.Delete("matchmaking_status")
	client.Printer.System(fmt.Sprintf("Leave: %s", resp.Message))
	return nil
}

// --- Stats ---

// StatsWaiting implements Endpoint for GET /api/v1/match/stats/waiting.
type StatsWaiting struct{}

func (e *StatsWaiting) Name() string        { return "stats_waiting" }
func (e *StatsWaiting) Description() string { return "Get waiting players count" }
func (e *StatsWaiting) Method() string      { return "GET" }
func (e *StatsWaiting) Path() string        { return "/api/v1/match/stats/waiting" }
func (e *StatsWaiting) Auth() bool          { return true }
func (e *StatsWaiting) Params() []Param     { return nil }

func (e *StatsWaiting) Next() []string {
	return []string{"matchmaking_join"}
}

func (e *StatsWaiting) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	_, err := client.Get(e.Path())
	return err
}

// StatsActive implements Endpoint for GET /api/v1/match/stats/active.
type StatsActive struct{}

func (e *StatsActive) Name() string        { return "stats_active" }
func (e *StatsActive) Description() string { return "Get active match count" }
func (e *StatsActive) Method() string      { return "GET" }
func (e *StatsActive) Path() string        { return "/api/v1/match/stats/active" }
func (e *StatsActive) Auth() bool          { return true }
func (e *StatsActive) Params() []Param     { return nil }

func (e *StatsActive) Next() []string {
	return []string{"matchmaking_join"}
}

func (e *StatsActive) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	_, err := client.Get(e.Path())
	return err
}

// --- Game Proxy ---

// GameState implements Endpoint for GET /api/v1/game/{id}.
type GameState struct{}

func (e *GameState) Name() string        { return "game_state" }
func (e *GameState) Description() string { return "Get cached board state for a match" }
func (e *GameState) Method() string      { return "GET" }
func (e *GameState) Path() string        { return "/api/v1/game/{id}" }
func (e *GameState) Auth() bool          { return true }
func (e *GameState) Params() []Param {
	return []Param{
		{Name: "id", Hint: "match UUID", Required: true, ContextKey: "match_id"},
	}
}

func (e *GameState) Next() []string {
	return []string{"game_action", "game_state"}
}

func (e *GameState) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	path := strings.ReplaceAll(e.Path(), "{id}", inputs["id"])
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	// Parse custom response structure
	var game dto.GameResponse
	// Re-marshal/unmarshal because resp.Data is interface{}
	dataBytes, _ := json.Marshal(resp.Data)
	if err := json.Unmarshal(dataBytes, &game); err == nil {
		sess.SetParticipants(game.Participants)
		sess.SetLastBoard(&game.GameState)
		client.Printer.System("Tactical state synchronized.")
		client.Printer.Board(&game.GameState, sess.UserIdentifier(), game.Participants)
	}

	return nil
}

// json tag helper
func (e *GameState) marshal(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}

// GameAction implements Endpoint for POST /api/v1/game/{id}/action.
type GameAction struct{}

func (e *GameAction) Name() string        { return "game_action" }
func (e *GameAction) Description() string { return "Send tactical action (MOVE/ATTACK/PASS/FORFEIT)" }
func (e *GameAction) Method() string      { return "POST" }
func (e *GameAction) Path() string        { return "/api/v1/game/{id}/action" }
func (e *GameAction) Auth() bool          { return true }
func (e *GameAction) Params() []Param {
	return []Param{
		{Name: "id", Hint: "match UUID", Required: true, ContextKey: "match_id"},
		{Name: "entity_id", Hint: "acting entity UUID", Required: true, ContextKey: "current_entity_id"},
		{Name: "type", Hint: "MOVE|ATTACK|PASS|FORFEIT", Required: true},
		{Name: "target_coords", Hint: "x,y coordinates (e.g. 3,2)"},
	}
}

func (e *GameAction) Next() []string {
	return []string{"game_state", "game_action"}
}

func (e *GameAction) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	path := strings.ReplaceAll(e.Path(), "{id}", inputs["id"])

	body := map[string]interface{}{
		"player_id": inputs["player_id"],
		"entity_id": inputs["entity_id"],
		"type":      inputs["type"],
	}

	// Parse target_coords "x,y" into [{x, y}]
	if tc := inputs["target_coords"]; tc != "" {
		var coords []map[string]int
		parts := strings.Split(tc, ";") // support multiple coords separated by ;
		for _, part := range parts {
			xy := strings.Split(strings.TrimSpace(part), ",")
			if len(xy) == 2 {
				x, y := 0, 0
				fmt.Sscanf(xy[0], "%d", &x)
				fmt.Sscanf(xy[1], "%d", &y)
				coords = append(coords, map[string]int{"x": x, "y": y})
			}
		}
		body["target_coords"] = coords
	}

	_, err := client.Post(path, body)
	return err
}

// --- Help ---

// HelpEndpoint implements Endpoint for GET /api/v1/help.
type HelpEndpoint struct{}

func (e *HelpEndpoint) Name() string        { return "api_help" }
func (e *HelpEndpoint) Description() string { return "Get machine-readable API documentation" }
func (e *HelpEndpoint) Method() string      { return "GET" }
func (e *HelpEndpoint) Path() string        { return "/api/v1/help" }
func (e *HelpEndpoint) Auth() bool          { return false }
func (e *HelpEndpoint) Params() []Param     { return nil }

func (e *HelpEndpoint) Next() []string {
	return nil
}

func (e *HelpEndpoint) Execute(client *api.Client, sess *session.Session, inputs map[string]string) error {
	_, err := client.Get(e.Path())
	return err
}

// --- Helpers ---

func extractCharacterIDs(resp *api.Response, sess *session.Session) {
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return
	}

	// Capture user info if present (flat or nested)
	if id, ok := data["user_id"].(string); ok {
		sess.Set("user_id", id)
	}
	if name, ok := data["account_name"].(string); ok {
		sess.Set("account_name", name)
	}

	// Deep capture if nested in "user" object
	if user, ok := data["user"].(map[string]interface{}); ok {
		if id, ok := user["id"].(string); ok {
			sess.Set("user_id", id)
		}
		if name, ok := user["account_name"].(string); ok {
			sess.Set("account_name", name)
		}
	}

	if chars, ok := data["characters"].([]interface{}); ok && len(chars) > 0 {
		var ids []string
		for _, item := range chars {
			if char, ok := item.(map[string]interface{}); ok {
				if id, ok := char["id"].(string); ok {
					ids = append(ids, id)
				}
			}
		}
		if len(ids) > 0 {
			sess.Set("character_id", ids[0])
			sess.Set("character_ids", strings.Join(ids, ","))
		}
	}
}

// RegisterAll populates the registry with all known endpoints.
func RegisterAll(reg *Registry) {
	// Auth
	reg.Register(&AuthLogin{})
	reg.Register(&AuthRegister{})
	reg.Register(&AuthLogout{})
	reg.Register(&AuthUpdate{})
	reg.Register(&AuthPassword{})
	reg.Register(&AuthExport{})
	reg.Register(&AuthDelete{})

	// Profile & Characters
	reg.Register(&ProfileGet{})
	reg.Register(&ProfileCharacters{})
	reg.Register(&ProfileCharacter{})
	reg.Register(&CharacterReroll{})
	reg.Register(&CharacterUpgrade{})
	reg.Register(&CharacterRename{})
	reg.Register(&CharacterDelete{})

	// Matchmaking
	reg.Register(&MatchmakingJoin{})
	reg.Register(&MatchmakingStatus{})
	reg.Register(&MatchmakingLeave{})

	// Stats
	reg.Register(&StatsWaiting{})
	reg.Register(&StatsActive{})

	// Game Proxy
	reg.Register(&GameState{})
	reg.Register(&GameAction{})

	// Help
	reg.Register(&HelpEndpoint{})
}
