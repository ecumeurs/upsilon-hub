package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ecumeurs/upsiloncli/internal/cli"
	"github.com/ecumeurs/upsiloncli/internal/endpoint"
	"github.com/ecumeurs/upsiloncli/internal/script"
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

	defaultBaseURL := os.Getenv("UPSILON_BASE_URL")
	if defaultBaseURL == "" {
		defaultBaseURL = "http://localhost:8000"
	}
	baseURL := flag.String("base-url", defaultBaseURL, "Laravel API base URL")
	auto := flag.Bool("auto", false, "Run full journey in autopilot mode")
	persist := flag.Bool("persist", false, "Load/save session to .upsilon_session.json")
	flag.BoolVar(persist, "P", false, "Load/save session to .upsilon_session.json (shorthand)")
	farm := flag.Bool("farm", false, "Execute multiple bot scripts in parallel")
	logDir := flag.String("logs", "", "Directory to store individual agent log files")
	flag.StringVar(logDir, "L", "", "Directory to store individual agent log files (shorthand)")
	flag.Parse()

	if *auto {
		fmt.Println("Autopilot mode — not yet implemented.")
		os.Exit(0)
	}

	if *farm {
		if flag.NArg() == 0 {
			fmt.Println("Error: --farm requires at least one script path.")
			os.Exit(1)
		}
		// Register endpoints
		reg := endpoint.NewRegistry()
		endpoint.RegisterAll(reg)
		script.RunFarm(*baseURL, reg, flag.Args(), *logDir)
		return
	}

	app := cli.New(*baseURL, *persist)

	// If there are remaining arguments, treat them as a single command line
	if flag.NArg() > 0 {
		app.ExecuteDirect(flag.Args())
		return
	}

	app.Run()
}
