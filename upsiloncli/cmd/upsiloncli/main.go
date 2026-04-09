package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ecumeurs/upsiloncli/internal/cli"
)

func main() {
	baseURL := flag.String("base-url", "http://localhost:8000", "Laravel API base URL")
	auto := flag.Bool("auto", false, "Run full journey in autopilot mode")
	persist := flag.Bool("persist", false, "Load/save session to .upsilon_session.json")
	flag.BoolVar(persist, "P", false, "Load/save session to .upsilon_session.json (shorthand)")
	flag.Parse()

	if *auto {
		fmt.Println("Autopilot mode — not yet implemented.")
		os.Exit(0)
	}

	app := cli.New(*baseURL, *persist)

	// If there are remaining arguments, treat them as a single command line
	if flag.NArg() > 0 {
		app.ExecuteDirect(flag.Args())
		return
	}

	app.Run()
}
