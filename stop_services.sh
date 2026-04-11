#!/bin/bash

# Configuration
PID_FILE=".services.pids"

if [ ! -f "$PID_FILE" ]; then
    echo "Summary: No tracked services appear to be running (missing $PID_FILE)."
    exit 0
fi

echo "---------------------------------------"
echo "Stopping Upsilon Stack..."
echo "---------------------------------------"

# Read PIDs and kill them
while IFS= read -r pid; do
    if [ ! -z "$pid" ]; then
        # Check if process exists before killing
        if ps -p "$pid" > /dev/null; then
            echo "[-] Stopping process $pid and its children..."
            # Try to kill the process group
            pkill -P "$pid" 2>/dev/null
            kill "$pid" 2>/dev/null
            sleep 0.5
            if ps -p "$pid" > /dev/null; then
                kill -9 "$pid" 2>/dev/null
            fi
        else
            echo "[!] Process $pid already stopped or defunct."
        fi
    fi
done < "$PID_FILE"

# Clean up
rm "$PID_FILE"

echo "---------------------------------------"
echo "All tracked services stopped."
echo "---------------------------------------"
