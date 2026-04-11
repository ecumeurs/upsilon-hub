#!/bin/bash

# Configuration
LOG_DIR="logs"
PID_FILE=".services.pids"

# Ensure log directory exists
mkdir -p "$LOG_DIR"

# Clear existing PIDs file
> "$PID_FILE"

echo "---------------------------------------"
echo "Starting Upsilon Stack in Background..."
echo "---------------------------------------"

# Function to start a service and record PID
start_service() {
    local name=$1
    local dir=$2
    local command=$3
    local log_file=$4

    echo "[+] Starting $name..."
    cd "$dir" || exit 1
    nohup $command > "../$LOG_DIR/$log_file" 2>&1 &
    local pid=$!
    echo "$pid" >> "../$PID_FILE"
    cd ..
    echo "    $name started (PID: $pid)"
}

# 1. Laravel API
start_service "Laravel API" "battleui" "php artisan serve" "laravel.log"

# 2. Reverb Server
start_service "Reverb Server" "battleui" "php artisan reverb:start" "reverb.log"

# 3. Vue Frontend
start_service "Vue Frontend" "battleui" "npm run dev" "vite.log"

# 4. Upsilon Engine (Go)
start_service "Upsilon Engine" "upsilonapi" "go run main.go" "engine.log"

echo "---------------------------------------"
echo "All services are running."
echo "Logs: $LOG_DIR/"
echo "Stop: ./stop_services.sh"
echo "---------------------------------------"
