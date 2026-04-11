#!/bin/bash
# run_all_battles.sh - Automated Tactical Engine Test Suite

set -e

# Configuration
CLI="./bin/upsiloncli"
SCRIPT="samples/pvp_bot_battle.js"
LOG_DIR="tests/logs"

# Ensure log directory exists
mkdir -p "$LOG_DIR"

# Ensure bin exists
if [ ! -f "$CLI" ]; then
    echo "Building CLI..."
    go build -o bin/upsiloncli cmd/upsiloncli/main.go
fi

run_test() {
    local mode=$1
    local agents=$2
    local log_file="$LOG_DIR/${mode}.log"
    
    echo "--------------------------------------------------"
    echo "TESTING MODE: $mode (Agents: $agents)"
    echo "--------------------------------------------------"
    
    # Construct paths array for the farm
    local paths=""
    for i in $(seq 1 "$agents"); do
        paths="$paths $SCRIPT"
    done
    
    # Set environment variable for the bot script
    export UPSILON_GAME_MODE="$mode"
    
    # Run the farm and capture output
    # We use a timeout to prevent hanging tests if the engine stalls
    timeout 300 $CLI --farm $paths > "$log_file" 2>&1 || true
    
    # Check for success indicators
    if grep -q "Game Over! Winner:" "$log_file"; then
        local winner=$(grep "Game Over! Winner:" "$log_file" | head -n 1 | awk -F': ' '{print $2}')
        echo -e "\033[32m[SUCCESS]\033[0m Match concluded! Winner: $winner"
    elif grep -q "STALEMATE" "$log_file"; then
        echo -e "\033[33m[SUCCESS]\033[0m Match concluded in a DRAW."
    else
        echo -e "\033[31m[FAILURE]\033[0m $mode did not conclude naturally or timed out."
        echo "Check $log_file for details."
        tail -n 20 "$log_file"
        exit 1
    fi
    echo ""
}

# Run the battery of tests
run_test "1v1_PVE" 1
run_test "2v2_PVE" 2
run_test "1v1_PVP" 2
run_test "2v2_PVP" 4

echo "=================================================="
echo -e "\033[32m\033[1mALL TESTS PASSED SUCCESSFULLY!\033[0m"
echo "=================================================="
