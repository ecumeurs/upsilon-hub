// Package display handles all terminal output formatting:
// curl commands, JSON responses, system messages, and board rendering.
package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
