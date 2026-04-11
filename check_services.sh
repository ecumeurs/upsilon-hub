#!/bin/bash

# Configuration
PID_FILE=".services.pids"
STATUS_CODE=0

if [ ! -f "$PID_FILE" ]; then
    echo "ERROR: No services are currently tracked (missing $PID_FILE)."
    exit 1
fi

echo "--- Upsilon Service Status ---"

# Read info and check status
while IFS='|' read -r name pid log port; do
    if [ ! -z "$pid" ]; then
        # Check if process exists
        if ps -p "$pid" > /dev/null; then
            # Validate port using ss (socket statistics)
            if [ ! -z "$port" ] && ss -Hlntp "sport = :$port" | grep -q "pid=$pid,"; then
                echo "[RUNNING] $name (PID: $pid, Port: $port)"
            elif [ ! -z "$port" ]; then
                echo "[PENDING] $name (PID: $pid, Waiting for Port $port...)"
                # STATUS_CODE=1 # Consider it running but not yet listening
            else
                echo "[RUNNING] $name (PID: $pid)"
            fi
        else
            echo "[ DOWN  ] $name (PID: $pid)"
            STATUS_CODE=1
        fi
    fi
done < "$PID_FILE"

if [ $STATUS_CODE -eq 0 ]; then
    echo "------------------------------"
    echo "All services are operational."
else
    echo "------------------------------"
    echo "WARNING: One or more services are DOWN."
fi

exit $STATUS_CODE
