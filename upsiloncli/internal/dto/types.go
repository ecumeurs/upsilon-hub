package dto

import "time"

// Position represents 2D coordinates.
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Entity represents a tactical unit on the board.
type Entity struct {
	ID       string   `json:"id"`
	PlayerID string   `json:"player_id"`
	Name     string   `json:"name"`
	HP       int      `json:"hp"`
	MaxHP    int      `json:"max_hp"`
	Attack   int      `json:"attack"`
	Defense  int      `json:"defense"`
	Move     int      `json:"move"`
	MaxMove  int      `json:"max_move"`
	Position Position `json:"position"`
}

// Cell represents a single tile on the grid.
type Cell struct {
	EntityID string `json:"entity_id"`
	Obstacle bool   `json:"obstacle"`
}

// Grid is the tactical map layout.
type Grid struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Cells  [][]Cell `json:"cells"`
}

// Turn represents an entry in the initiative timeline.
type Turn struct {
	PlayerID string `json:"player_id"`
	EntityID string `json:"entity_id"`
	Delay    int    `json:"delay"`
}

// BoardState is the full DTO for the tactical situation.
type BoardState struct {
	Entities        []Entity  `json:"entities"`
	Grid            Grid      `json:"grid"`
	Turn            []Turn    `json:"turn"`
	CurrentPlayerID string    `json:"current_player_id"`
	CurrentEntityID string    `json:"current_entity_id"`
	Timeout         time.Time `json:"timeout"`
	StartTime       time.Time `json:"start_time"`
	WinnerID        string    `json:"winner_id"`
}

// Participant links a player UUID to a team and nickname.
type Participant struct {
	PlayerID string `json:"player_id"`
	Nickname string `json:"nickname"`
	Team     int    `json:"team"`
}

// GameResponse is the expanded response from GET /api/v1/game/{id}
type GameResponse struct {
	MatchID      string        `json:"match_id"`
	GameMode     string        `json:"game_mode"`
	GameState    BoardState    `json:"game_state"`
	Participants []Participant `json:"participants"`
}
