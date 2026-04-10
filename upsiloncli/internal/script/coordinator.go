package script

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dop251/goja"
	"github.com/ecumeurs/upsiloncli/internal/endpoint"
)

func RunFarm(baseURL string, reg *endpoint.Registry, scriptPaths []string, logDir string) {
	var wg sync.WaitGroup

	for i, path := range scriptPaths {
		wg.Add(1)
		go func(agentIdx int, scriptPath string) {
			defer wg.Done()

			agentID := fmt.Sprintf("Bot-%02d", agentIdx+1)
			
			var logger *os.File
			if logDir != "" {
				fileName := fmt.Sprintf("%s.log", agentID)
				f, err := os.Create(filepath.Join(logDir, fileName))
				if err != nil {
					fmt.Printf("[Farm] Error creating log file for %s: %v\n", agentID, err)
					logger = os.Stdout
				} else {
					logger = f
					defer f.Close()
				}
			} else {
				logger = os.Stdout
			}

			agent := NewAgent(agentID, baseURL, reg, logger)
			agent.Listener.Start()
			
			scriptData, err := os.ReadFile(scriptPath)
			if err != nil {
				fmt.Fprintf(logger, "[%s] Error reading script: %v\n", agentID, err)
				return
			}

			_, err = agent.VM.RunString(string(scriptData))
			if err != nil {
				if jsErr, ok := err.(*goja.Exception); ok {
					fmt.Fprintf(logger, "[%s] JS Exception: %v\n", agentID, jsErr.String())
				} else {
					fmt.Fprintf(logger, "[%s] Execution failed: %v\n", agentID, err)
				}
			}
			
			fmt.Fprintf(logger, "[%s] Script execution finished.\n", agentID)
		}(i, path)
	}

	wg.Wait()
	fmt.Println("All agents have finished execution.")
}
