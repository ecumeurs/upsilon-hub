package script

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/dop251/goja"
	"github.com/ecumeurs/upsiloncli/internal/endpoint"
)

func RunFarm(baseURL string, reg *endpoint.Registry, scriptPaths []string, logDir string) {
	var wg sync.WaitGroup
	sharedStore := NewSharedStore()

	var agents []*Agent
	var agentsMu sync.Mutex

	// Catch SIGINT/SIGTERM to allow graceful teardown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		sig := <-sigChan
		fmt.Printf("\n[Farm] Received %v. Interrupting agents for cleanup...\n", sig)
		
		agentsMu.Lock()
		for _, a := range agents {
			// Interrupt the VM execution. This causes RunString and any blocking bridge 
			// calls (like sleep or waitForEvent) to return with an error.
			a.VM.Interrupt(fmt.Errorf("interrupted by %v", sig))
		}
		agentsMu.Unlock()
	}()

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

			agent := NewAgent(agentID, baseURL, reg, logger, sharedStore)
			
			agentsMu.Lock()
			agents = append(agents, agent)
			agentsMu.Unlock()

			agent.Listener.Start()
			
			// GUARANTEED TEARDOWN BLOCK
			defer func() {
				if agent.TeardownHook != nil {
					// Execute the JS teardown function safely
					_, err := agent.TeardownHook(goja.Undefined())
					if err != nil {
						fmt.Fprintf(logger, "[%s] Teardown hook failed: %v\n", agentID, err)
					}
				}
				// Ensure WebSocket is closed cleanly
				agent.Listener.Stop() 
			}()
			
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
	fmt.Println("All agents have finished execution and cleanup.")
}
