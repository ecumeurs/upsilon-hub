#!/bin/bash
cd "$(dirname "$0")/.."
# zombie_killer.sh - Forcefully clean up hanging Upsilon processes
echo "Scanning for zombie/hanging processes..."
TARGETS=("upsiloncli" "upsilonapi" "upsilonbattle")
for target in "${TARGETS[@]}"; do
    PIDS=$(pgrep -f "$target")
    if [ ! -z "$PIDS" ]; then
        echo "Found $target PIDs: $PIDS. Killing..."
        pkill -9 -f "$target"
    else
        echo "No $target processes found."
    fi
done
echo "Cleanup complete."
