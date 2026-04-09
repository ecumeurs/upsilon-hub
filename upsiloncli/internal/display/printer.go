// Package display handles all terminal output formatting:
// curl commands, JSON responses, system messages, and board rendering.
package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/ecumeurs/upsiloncli/internal/dto"
)

// ANSI color codes
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"

	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BgRed   = "\033[41m"
	BgGreen = "\033[42m"
)

// Printer handles formatted terminal output.
type Printer struct{}

// NewPrinter creates a new terminal printer.
func NewPrinter() *Printer {
	return &Printer{}
}

// Curl prints the equivalent curl command for an API request.
func (p *Printer) Curl(method, url string, headers http.Header, body []byte) {
	fmt.Println()
	fmt.Printf("%s[CURL]%s ", Cyan+Bold, Reset)

	var parts []string
	parts = append(parts, "curl", "-X", method)

	for key, vals := range headers {
		if key == "Content-Type" || key == "Accept" || key == "Authorization" {
			parts = append(parts, "-H", fmt.Sprintf("'%s: %s'", key, vals[0]))
		}
	}

	if len(body) > 0 {
		parts = append(parts, "-d", fmt.Sprintf("'%s'", string(body)))
	}

	parts = append(parts, fmt.Sprintf("'%s'", url))
	fmt.Println(Dim + strings.Join(parts, " ") + Reset)
}

// Response prints the HTTP status and pretty-printed JSON body.
func (p *Printer) Response(statusCode int, body []byte) {
	color := Green
	if statusCode >= 400 {
		color = Red
	} else if statusCode >= 300 {
		color = Yellow
	}

	fmt.Printf("%s[REPLY %d]%s ", color+Bold, statusCode, Reset)

	// Pretty-print JSON
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "  ", "  "); err == nil {
		fmt.Println()
		fmt.Println("  " + Dim + pretty.String() + Reset)
	} else {
		fmt.Println(Dim + string(body) + Reset)
	}
}

// System prints a system-level informational message.
func (p *Printer) System(msg string) {
	fmt.Printf("%s[SYSTEM]%s %s\n", Yellow+Bold, Reset, msg)
}

// Warn prints a warning message.
func (p *Printer) Warn(msg string) {
	fmt.Printf("%s[WARN]%s %s\n", Red+Bold, Reset, msg)
}

// WebSocket prints a received WebSocket event.
func (p *Printer) WebSocket(eventType string, payload []byte) {
	fmt.Printf("\n%s[WS]%s %s event received.\n", Magenta+Bold, Reset, eventType)
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, payload, "  ", "  "); err == nil {
		fmt.Println("  " + Dim + pretty.String() + Reset)
	}
}

// Prompt displays a prompt for user input with an optional default value.
func (p *Printer) Prompt(name, hint, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("  %s%s%s [default: %s%s%s]: ", Bold, name, Reset, Green, defaultVal, Reset)
	} else if hint != "" {
		fmt.Printf("  %s%s%s (%s): ", Bold, name, Reset, hint)
	} else {
		fmt.Printf("  %s%s%s: ", Bold, name, Reset)
	}
	return "" // actual reading is done by the caller
}

// RouteTable prints the endpoint registry as a table.
func (p *Printer) RouteTable(routes []RouteInfo) {
	fmt.Println()
	fmt.Printf("  %s%-25s %-8s %-40s %s%s\n", Bold, "ROUTE NAME", "VERB", "PATH", "DESCRIPTION", Reset)
	fmt.Printf("  %s%s%s\n", Dim, strings.Repeat("─", 100), Reset)
	for _, r := range routes {
		authMark := " "
		if r.Auth {
			authMark = "🔒"
		}
		fmt.Printf("  %-25s %-8s %-40s %s %s\n", Green+r.Name+Reset, r.Method, Dim+r.Path+Reset, r.Description, authMark)
	}
	fmt.Println()
}

// RouteInfo is used by RouteTable to describe a registered endpoint.
type RouteInfo struct {
	Name        string
	Method      string
	Path        string
	Description string
	Auth        bool
}

// SessionInfo prints the current session state.
func (p *Printer) SessionInfo(data map[string]string) {
	fmt.Println()
	fmt.Printf("  %sSession Context%s\n", Bold, Reset)
	fmt.Printf("  %s%s%s\n", Dim, strings.Repeat("─", 40), Reset)
	for k, v := range data {
		fmt.Printf("  %-20s %s\n", Cyan+k+Reset, v)
	}
	fmt.Println()
}
// Board renders the tactical map and entity status table.
func (p *Printer) Board(bs *dto.BoardState, currentUserID string, participants []dto.Participant) {
	if bs == nil {
		p.Warn("No board state available.")
		return
	}

	// 1. Setup Roles and Symbols
	var myTeam int
	var actualCurrentUserID string
	for _, part := range participants {
		// Matching current user by ID or by matching nickname if we have it in session
		if part.PlayerID == currentUserID || (currentUserID != "" && part.Nickname == currentUserID) {
			myTeam = part.Team
			actualCurrentUserID = part.PlayerID // Normalize to the ID used in the match
			break
		}
	}

	// Fallback if no exact match (e.g. currentUserID is account_name but participants uses UUIDs)
	if actualCurrentUserID == "" && currentUserID != "" {
		// Try to find by nickname in session if available
		for _, part := range participants {
			if part.PlayerID == currentUserID {
				actualCurrentUserID = part.PlayerID
				myTeam = part.Team
				break
			}
		}
	}

	var allyID string
	var enemies []string
	nicknames := make(map[string]string)
	for _, part := range participants {
		nicknames[part.PlayerID] = part.Nickname
		if part.PlayerID == actualCurrentUserID {
			continue
		}
		if part.Team == myTeam {
			allyID = part.PlayerID
		} else {
			enemies = append(enemies, part.PlayerID)
		}
	}

	// Group entities by player
	playerEntities := make(map[string][]dto.Entity)
	for _, ent := range bs.Entities {
		if ent.HP > 0 {
			playerEntities[ent.PlayerID] = append(playerEntities[ent.PlayerID], ent)
		}
	}
	// Sort for deterministic symbols
	for pid := range playerEntities {
		sort.Slice(playerEntities[pid], func(i, j int) bool {
			return playerEntities[pid][i].ID < playerEntities[pid][j].ID
		})
	}

	entitySymbols := make(map[string]string)
	entityColors := make(map[string]string)

	assign := func(pID string, syms []string, color string) {
		for i, ent := range playerEntities[pID] {
			if i < len(syms) {
				entitySymbols[ent.ID] = syms[i]
				entityColors[ent.ID] = color
			} else {
				entitySymbols[ent.ID] = "?"
				entityColors[ent.ID] = color
			}
		}
	}

	assign(actualCurrentUserID, []string{"A", "B", "C"}, Green)
	if allyID != "" {
		assign(allyID, []string{"a", "b", "c"}, Green)
	}
	if len(enemies) > 0 {
		assign(enemies[0], []string{"X", "Y", "Z"}, Red)
	}
	if len(enemies) > 1 {
		assign(enemies[1], []string{"x", "y", "z"}, Red)
	}

	// 2. Render Grid
	fmt.Println()
	fmt.Printf("  %sTACTICAL FEED — MATCH DATA%s\n", Cyan+Bold, Reset)
	fmt.Printf("  %s%s%s\n", Dim, strings.Repeat("─", 40), Reset)

	// Top border
	fmt.Print("    ")
	for x := 0; x < bs.Grid.Width; x++ {
		fmt.Printf("%2d", x)
	}
	fmt.Println()

	for y := 0; y < bs.Grid.Height; y++ {
		fmt.Printf("%2d │", y)
		for x := 0; x < bs.Grid.Width; x++ {
			cell := bs.Grid.Cells[x][y]
			if cell.EntityID != "" {
				sym := entitySymbols[cell.EntityID]
				color := entityColors[cell.EntityID]
				if cell.EntityID == bs.CurrentEntityID {
					fmt.Printf("%s%s%s ", color+Bold+BgGreen, sym, Reset) // Highlight current turn? No, let's keep it simple
				} else {
					fmt.Printf("%s%s%s ", color+Bold, sym, Reset)
				}
			} else if cell.Obstacle {
				fmt.Printf("%s#%s ", Dim, Reset)
			} else {
				fmt.Printf("%s.%s ", Dim, Reset)
			}
		}
		fmt.Println("│")
	}

	// Map entity ID to delay
	delays := make(map[string]int)
	for _, t := range bs.Turn {
		delays[t.EntityID] = t.Delay
	}

	// 3. Entity List
	fmt.Println()
	fmt.Printf("  %s%-3s %-15s %-12s %-10s %-7s %-5s %s\n", Bold, "ID", "UNIT NAME", "OWNER", "HP/MAX", "MVT", "DELAY", Reset)
	fmt.Printf("  %s%s%s\n", Dim, strings.Repeat("─", 70), Reset)

	// Sort entities by symbol for the list
	displayEnts := bs.Entities
	sort.Slice(displayEnts, func(i, j int) bool {
		symI := entitySymbols[displayEnts[i].ID]
		symJ := entitySymbols[displayEnts[j].ID]
		return symI < symJ
	})

	for _, ent := range displayEnts {
		if ent.HP <= 0 {
			continue // Hide dead units
		}
		sym := entitySymbols[ent.ID]
		color := entityColors[ent.ID]
		owner := nicknames[ent.PlayerID]
		if owner == "" {
			owner = "System/AI"
		}
		if ent.ID == bs.CurrentEntityID {
			fmt.Print(Cyan + "> " + Reset)
		} else {
			fmt.Print("  ")
		}
		
		delayStr := fmt.Sprintf("%d", delays[ent.ID])

		fmt.Printf("%s%s%s %-15s %-12s %-10s %-7d %-5s\n",
			color+Bold, sym, Reset,
			ent.Name,
			owner,
			fmt.Sprintf("%d/%d", ent.HP, ent.MaxHP),
			ent.Move,
			delayStr,
		)
	}
	fmt.Println()
}
