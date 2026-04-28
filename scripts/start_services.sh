#!/bin/bash
cd "$(dirname "$0")/.."

# Configuration
LOG_DIR="logs"
PID_FILE=".services.pids"
PORTS_VERIFY_TIMEOUT=10 # Max seconds to wait for a port

# Ensure log directory exists
mkdir -p "$LOG_DIR"

echo "---------------------------------------"
echo "Starting Upsilon Stack (Authoritative Mode)"
echo "---------------------------------------"

# 1. Authoritative Cleanup
./scripts/stop_services.sh

# Clear existing PIDs file
> "$PID_FILE"

# Function to start a service and verify it is listening
start_service() {
    local name=$1
    local dir=$2
    local command=$3
    local log_file=$4
    local port=$5

    echo "[+] Starting $name on port $port..."
    cd "$dir" || exit 1
    
    # Run in background
    nohup $command > "../$LOG_DIR/$log_file" 2>&1 &
    local shell_pid=$!
    
    # Verification loop
    echo -n "    Verifying..."
    local count=0
    local confirmed_pid=""
    while [ $count -lt $PORTS_VERIFY_TIMEOUT ]; do
        # Try to find the PID actually listening on the port
        confirmed_pid=$(ss -tulpn | grep ":$port " | grep -oP 'users:\(\("\S+",pid=\K\d+' | head -n 1)
        
        if [ ! -z "$confirmed_pid" ]; then
            echo " OK (PID: $confirmed_pid)"
            # Format: NAME|PID|LOG_FILE|PORT
            echo "$name|$confirmed_pid|$log_file|$port" >> "../$PID_FILE"
            cd ..
            return 0
        fi
        
        sleep 1
        echo -n "."
        count=$((count + 1))
    done

    echo " FAILED"
    echo "[ERROR] $name failed to start on port $port after $PORTS_VERIFY_TIMEOUT seconds."
    echo "Check logs/ $log_file for details."
    cd ..
    exit 1
}


# 1. Laravel API
start_service "Laravel API" "battleui" "php artisan serve --host=0.0.0.0 --port=8000" "laravel.log" 8000

# 2. Reverb Server
start_service "Reverb Server" "battleui" "php artisan reverb:start" "reverb.log" 8080

# 3. Vue Frontend
start_service "Vue Frontend" "battleui" "npm run dev" "vite.log" 5173

# 4. Upsilon Engine (Go)
start_service "Upsilon Engine" "upsilonapi" "./bin/upsilonapi" "engine.log" 8081

echo "---------------------------------------"
echo "All services are running and verified."
echo "Logs: $LOG_DIR/"
echo "Stop: ./scripts/stop_services.sh"
echo "---------------------------------------"

