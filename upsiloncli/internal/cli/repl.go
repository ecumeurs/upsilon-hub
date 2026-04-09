// Package cli provides the interactive REPL loop and command dispatcher.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ecumeurs/upsiloncli/internal/api"
	"github.com/ecumeurs/upsiloncli/internal/display"
	"github.com/ecumeurs/upsiloncli/internal/endpoint"
	"github.com/ecumeurs/upsiloncli/internal/session"
)

// CLI is the interactive command-line application.
type CLI struct {
	Session  *session.Session
	Client   *api.Client
	Printer  *display.Printer
	Registry *endpoint.Registry
	Scanner  *bufio.Scanner
}

// New creates a new CLI instance.
func New(baseURL string) *CLI {
	sess := session.New()
	printer := display.NewPrinter()
	client := api.NewClient(baseURL, sess, printer)
	reg := endpoint.NewRegistry()
	endpoint.RegisterAll(reg)

	return &CLI{
		Session:  sess,
		Client:   client,
		Printer:  printer,
		Registry: reg,
		Scanner:  bufio.NewScanner(os.Stdin),
	}
}

// Run starts the interactive REPL loop.
func (c *CLI) Run() {
	c.printBanner()

	for {
		fmt.Printf("\n%s[%s]%s > ", display.Cyan, c.Session.String(), display.Reset)

		if !c.Scanner.Scan() {
			break
		}

		input := strings.TrimSpace(c.Scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmd := strings.ToLower(parts[0])
		args := parts[1:]

		switch cmd {
		case "exit", "quit", "q":
			fmt.Println("Goodbye.")
			return

		case "help", "h", "?":
			c.printHelp()

		case "routes":
			c.Printer.RouteTable(c.Registry.List())

		case "call":
			if len(args) == 0 {
				c.Printer.Warn("Usage: call <route_name>")
				continue
			}
			c.executeEndpoint(args[0])

		case "jwt":
			if len(args) == 0 {
				// Display current JWT
				token := c.Session.Token()
				if token == "" {
					c.Printer.System("No JWT token set.")
				} else {
					c.Printer.System(fmt.Sprintf("Current JWT: %s", token))
				}
			} else {
				// Override JWT
				c.Session.SetToken(args[0])
				c.Printer.Warn("JWT manually overridden. All further requests will use this token.")
			}

		case "session":
			c.Printer.SessionInfo(c.Session.Dump())

		case "redraw":
			c.Printer.System("Board redraw — not yet implemented (pending WebSocket integration).")

		default:
			// Check if it's a valid route_name shortcut
			if ep := c.Registry.Get(cmd); ep != nil {
				c.executeEndpoint(cmd)
			} else {
				c.Printer.Warn(fmt.Sprintf("Unknown command: %q. Type 'help' for available commands.", cmd))
			}
		}
	}
}

// executeEndpoint runs an endpoint by name, prompting for parameters.
func (c *CLI) executeEndpoint(name string) {
	ep := c.Registry.Get(name)
	if ep == nil {
		c.Printer.Warn(fmt.Sprintf("Unknown route: %q. Use 'routes' to list available endpoints.", name))
		return
	}

	// Check auth requirement
	if ep.Auth() && !c.Session.HasToken() {
		c.Printer.Warn("This endpoint requires authentication. Use 'call auth_login' or 'call auth_register' first.")
		return
	}

	// Collect parameters
	params := ep.Params()
	inputs := make(map[string]string)

	for _, p := range params {
		// Resolve default from session context
		defaultVal := ""
		if p.ContextKey != "" {
			if v, ok := c.Session.Get(p.ContextKey); ok {
				defaultVal = v
			}
		}

		// Prompt user
		value := c.prompt(p.Name, p.Hint, defaultVal, p.Required)
		inputs[p.Name] = value
	}

	// Execute
	if err := ep.Execute(c.Client, c.Session, inputs); err != nil {
		c.Printer.Warn(fmt.Sprintf("Error: %v", err))
		return
	}

	// Suggest next routes
	next := ep.Next()
	if len(next) > 0 {
		var formatted []string
		for _, n := range next {
			if nxtEp := c.Registry.Get(n); nxtEp != nil {
				formatted = append(formatted, display.Green+n+display.Reset)
			}
		}
		if len(formatted) > 0 {
			fmt.Printf("\n  %sSuggested next steps:%s %s\n", display.Dim, display.Reset, strings.Join(formatted, ", "))
		}
	}
}

// prompt asks the user for a value, showing the default if available.
func (c *CLI) prompt(name, hint, defaultVal string, required bool) string {
	for {
		if defaultVal != "" {
			fmt.Printf("  %s%s%s [default: %s%s%s]: ", display.Bold, name, display.Reset, display.Green, defaultVal, display.Reset)
		} else if hint != "" {
			fmt.Printf("  %s%s%s (%s): ", display.Bold, name, display.Reset, hint)
		} else {
			fmt.Printf("  %s%s%s: ", display.Bold, name, display.Reset)
		}

		c.Scanner.Scan()
		value := strings.TrimSpace(c.Scanner.Text())

		if value == "" && defaultVal != "" {
			return defaultVal
		}
		if value == "" && required {
			c.Printer.Warn(fmt.Sprintf("%s is required.", name))
			continue
		}
		return value
	}
}

func (c *CLI) printBanner() {
	fmt.Println()
	fmt.Println(display.Cyan + display.Bold + "  ╔══════════════════════════════════════════╗" + display.Reset)
	fmt.Println(display.Cyan + display.Bold + "  ║       ⌬ UpsilonCLI — API Explorer        ║" + display.Reset)
	fmt.Println(display.Cyan + display.Bold + "  ╚══════════════════════════════════════════╝" + display.Reset)
	fmt.Printf("  Target: %s%s%s\n", display.Dim, c.Client.BaseURL, display.Reset)
	fmt.Println()
	fmt.Println("  Type 'help' for commands, 'routes' to see all endpoints.")
}

func (c *CLI) printHelp() {
	fmt.Println()
	fmt.Printf("  %sAvailable Commands%s\n", display.Bold, display.Reset)
	fmt.Printf("  %s%s%s\n", display.Dim, strings.Repeat("─", 50), display.Reset)
	fmt.Printf("  %-22s %s\n", display.Green+"routes"+display.Reset, "List all API endpoints with route_name identifiers")
	fmt.Printf("  %-22s %s\n", display.Green+"call <route_name>"+display.Reset, "Execute an endpoint interactively")
	fmt.Printf("  %-22s %s\n", display.Green+"jwt"+display.Reset, "Display current JWT token")
	fmt.Printf("  %-22s %s\n", display.Green+"jwt <token>"+display.Reset, "Manually override the JWT (for testing)")
	fmt.Printf("  %-22s %s\n", display.Green+"session"+display.Reset, "Display current session context")
	fmt.Printf("  %-22s %s\n", display.Green+"redraw"+display.Reset, "Re-render last known tactical board")
	fmt.Printf("  %-22s %s\n", display.Green+"help"+display.Reset, "Show this help message")
	fmt.Printf("  %-22s %s\n", display.Green+"exit"+display.Reset, "Quit the CLI")
	fmt.Println()
	fmt.Printf("  %sTip:%s You can also type a route_name directly (e.g., 'auth_login').\n", display.Dim, display.Reset)
}
