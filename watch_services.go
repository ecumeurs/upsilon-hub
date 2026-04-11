// @spec-link [[watch_services]]
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Service struct {
	Name    string
	PID     string
	LogFile string
	Port    string
	CPU     string
	Mem     string
	Errors  []string
	Alive   bool
}

const (
	PidFile = ".services.pids"
	LogDir  = "logs"
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
	ClearScreen = "\033[H\033[2J"
)

func main() {
	for {
		services, err := loadServices()
		if err != nil {
			fmt.Printf("Error loading services: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		updateMetrics(services)
		checkLogs(services)

		render(services)

		time.Sleep(2 * time.Second)
	}
}

func loadServices() ([]*Service, error) {
	file, err := os.Open(PidFile)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", PidFile, err)
	}
	defer file.Close()

	var services []*Service
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if len(parts) >= 2 {
			svc := &Service{
				Name:  parts[0],
				PID:   parts[1],
				Alive: true,
			}
			if len(parts) >= 3 {
				svc.LogFile = parts[2]
			}
			if len(parts) >= 4 {
				svc.Port = parts[3]
			}
			services = append(services, svc)
		}
	}
	return services, scanner.Err()
}

func updateMetrics(services []*Service) {
	for _, svc := range services {
		// Get CPU and Mem using ps
		out, err := exec.Command("ps", "-p", svc.PID, "-o", "%cpu,%mem", "--no-headers").Output()
		if err != nil {
			svc.Alive = false
			svc.CPU = "0.0"
			svc.Mem = "0.0"
			continue
		}

		fields := strings.Fields(string(out))
		if len(fields) >= 2 {
			svc.CPU = fields[0]
			svc.Mem = fields[1]
		}
	}
}

func checkLogs(services []*Service) {
	for _, svc := range services {
		if svc.LogFile == "" {
			continue
		}

		path := filepath.Join(LogDir, svc.LogFile)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		// Only check if modified in the last 2 minutes (buffer)
		if time.Since(info.ModTime()) > 2*time.Minute {
			continue
		}

		// Read the last 100 lines
		cmd := exec.Command("tail", "-n", "100", path)
		out, err := cmd.Output()
		if err != nil {
			continue
		}

		lines := strings.Split(string(out), "\n")
		var recentErrors []string
		errorKeywords := []string{"error", "panic", "fatal", "exception"}

		for _, line := range lines {
			lowerLine := strings.ToLower(line)
			isError := false
			for _, kw := range errorKeywords {
				if strings.Contains(lowerLine, kw) {
					isError = true
					break
				}
			}

			if isError {
				// To keep it simple, we just take the last 3 unique errors
				recentErrors = append(recentErrors, line)
			}
		}

		if len(recentErrors) > 3 {
			svc.Errors = recentErrors[len(recentErrors)-3:]
		} else {
			svc.Errors = recentErrors
		}
	}
}

func render(services []*Service) {
	fmt.Print(ClearScreen)
	fmt.Printf("%s%s=== UPSILON SERVICE MONITOR ===%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("Time: %s\n\n", time.Now().Format("15:04:05"))

	fmt.Printf("%-20s | %-8s | %-6s | %-6s | %-6s | %-10s\n", "SERVICE", "PID", "PORT", "CPU%", "MEM%", "STATUS")
	fmt.Println(strings.Repeat("-", 70))

	for _, svc := range services {
		status := fmt.Sprintf("%sOK%s", ColorGreen, ColorReset)
		if !svc.Alive {
			status = fmt.Sprintf("%sDOWN%s", ColorRed, ColorReset)
		} else if len(svc.Errors) > 0 {
			status = fmt.Sprintf("%sERRORS%s", ColorYellow, ColorReset)
		}

		nameColor := ColorReset
		if !svc.Alive {
			nameColor = ColorRed
		} else if len(svc.Errors) > 0 {
			nameColor = ColorYellow
		}

		portDisplay := svc.Port
		if portDisplay == "" {
			portDisplay = "N/A"
		}

		fmt.Printf("%s%-20s%s | %-8s | %-6s | %-6s | %-6s | %-10s\n", 
			nameColor, svc.Name, ColorReset, svc.PID, portDisplay, svc.CPU, svc.Mem, status)
	}

	hasErrors := false
	for _, svc := range services {
		if len(svc.Errors) > 0 {
			if !hasErrors {
				fmt.Printf("\n%s%s--- RECENT ERRORS (< 1m) ---%s\n", ColorBold, ColorYellow, ColorReset)
				hasErrors = true
			}
			for _, errLine := range svc.Errors {
				// Clean up common long lines or prefixes if needed, but let's keep it raw for now
				if len(errLine) > 100 {
					errLine = errLine[:97] + "..."
				}
				fmt.Printf("[%s%s%s] %s\n", ColorCyan, svc.Name, ColorReset, strings.TrimSpace(errLine))
			}
		}
	}

	fmt.Printf("\n%sPress Ctrl+C to exit%s\n", ColorYellow, ColorReset)
}
