package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ecumeurs/upsiloncli/internal/cli"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if it exists
	_ = godotenv.Load()

	appKey := os.Getenv("REVERB_APP_KEY")
	if appKey == "" {
		fmt.Printf("\033[31m\033[1m[ERROR]\033[0m Mandatory environment variable REVERB_APP_KEY is missing.\n")
		fmt.Println("Please set it in your system environment or a .env file.")
		os.Exit(1)
	}

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
