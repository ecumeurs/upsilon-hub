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
	flag.Parse()

	if *auto {
		fmt.Println("Autopilot mode — not yet implemented.")
		os.Exit(0)
	}

	app := cli.New(*baseURL)
	app.Run()
}
