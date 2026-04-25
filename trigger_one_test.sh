#!/bin/bash
# trigger_one_test.sh - Run a single E2E or Edge Case test locally
# Usage: ./trigger_one_test.sh <script_name_or_path>

set -e

# 1. Pre-flight Check
if [ ! -f "./check_services.sh" ]; then
    echo "ERROR: check_services.sh not found."
    exit 1
fi

echo "Checking local services..."
if ! ./check_services.sh; then
    echo "ERROR: Some services are down. Please start the stack before running tests."
    exit 1
fi

# 2. Build CLI (ensure it exists)
CLI="./upsiloncli/bin/upsiloncli"
if [ ! -f "$CLI" ]; then
    echo "CLI binary not found. Building..."
    cd upsiloncli && go build -o bin/upsiloncli cmd/upsiloncli/main.go && cd ..
fi

# 3. Parameter handling
SCRIPT=$1
if [ -z "$SCRIPT" ]; then
    echo "Usage: ./trigger_one_test.sh <script_name_or_path>"
    echo "Example: ./trigger_one_test.sh e2e_combat_turn_management"
    exit 1
fi

# Find the script if only name was provided
if [ ! -f "$SCRIPT" ]; then
    FOUND=$(find upsiloncli/tests/scenarios -name "${SCRIPT}.js" | head -n 1)
    if [ -n "$FOUND" ]; then
        SCRIPT=$FOUND
    else
        echo "ERROR: Script not found: $SCRIPT"
        exit 1
    fi
fi

NAME=$(basename "$SCRIPT" .js)
LOG_DIR="upsiloncli/tests/logs"
LOG_FILE="$LOG_DIR/${NAME}.log"
mkdir -p "$LOG_DIR"

# 4. Determine agent count (following CI logic)
AGENTS=1
if [[ "$NAME" == *"pvp"* ]] || [[ "$NAME" == *"coordination"* ]] || [[ "$NAME" == *"combat"* ]] || [[ "$NAME" == *"friendly_fire"* ]] || [[ "$NAME" == *"resolution_standard"* ]] || [[ "$NAME" == *"progression_constraints"* ]]; then
    AGENTS=2
fi
if [[ "$NAME" == *"2v2"* ]]; then
    AGENTS=4
fi

echo "=================================================="
echo "Running $NAME (Agents: $AGENTS)... "
echo "Log: $LOG_FILE"
echo "=================================================="

# Construct paths array for the farm
PATHS=""
for i in $(seq 1 "$AGENTS"); do
    PATHS="$PATHS $SCRIPT"
done

# Run the farm with --local flag
# No timeout here to allow debugging, but follow same structure as CI
if "$CLI" --local --farm -L "$LOG_DIR" $PATHS ; then
    echo -e "\033[32m[PASSED]\033[0m"
    echo "[SCENARIO_RESULT: PASSED]" >> "$LOG_FILE"
    exit 0
else
    echo -e "\033[31m[FAILED]\033[0m"
    echo "[SCENARIO_RESULT: FAILED]" >> "$LOG_FILE"
    echo "Check logs at: $LOG_FILE"
    exit 1
fi
