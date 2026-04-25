#!/bin/bash
# trigger_all_tests.sh - Run all E2E and Edge Case tests locally

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

# 3. Execution Setup
SCENARIO_DIR="upsiloncli/tests/scenarios"
LOG_DIR="upsiloncli/tests/logs"
mkdir -p "$LOG_DIR"

echo "=================================================="
echo "      UPSILON LOCAL TEST TRIGGER (E2E + EDGE)"
echo "=================================================="

# 0. Purge stale data
echo "Purging stale match data..."
"$CLI" --local --quiet upsiloncli/tests/scenarios/util_purge_all.js > /dev/null 2>&1 || true

FAILED_TESTS=""
PASSED_COUNT=0
FAILED_COUNT=0

run_test() {
    local script=$1
    local name=$(basename "$script" .js)
    local log_file="$LOG_DIR/${name}.log"
    
    # Determine agent count (following CI logic)
    local agents=1
    if [[ "$name" == *"pvp"* ]] || [[ "$name" == *"coordination"* ]] || [[ "$name" == *"combat"* ]] || [[ "$name" == *"resolution_standard"* ]] || [[ "$name" == *"progression_constraints"* ]] || [[ "$name" == *"progression_post_win"* ]] || [[ "$name" == *"out_of_turn"* ]]; then
        agents=2
    fi
    
    if [[ "$name" == *"2v2"* ]] || [[ "$name" == *"targeting_rules"* ]] || [[ "$name" == *"friendly_fire"* ]]; then
        agents=4
    fi

    echo -n "Running $name (Agents: $agents)... "

    # Construct paths array for the farm
    local paths=""
    for i in $(seq 1 "$agents"); do
        paths="$paths $script"
    done

    # Run the farm with --local flag
    if timeout 180 "$CLI" --local --farm -L "$LOG_DIR" $paths > /dev/null 2>&1; then
        echo -e "\033[32m[PASSED]\033[0m"
        echo "[SCENARIO_RESULT: PASSED]" >> "$log_file"
        PASSED_COUNT=$((PASSED_COUNT + 1))
    else
        echo -e "\033[31m[FAILED]\033[0m"
        echo "[SCENARIO_RESULT: FAILED]" >> "$log_file"
        FAILED_COUNT=$((FAILED_COUNT + 1))
        FAILED_TESTS="$FAILED_TESTS $name"
        # Don't exit on failure, continue to other tests
    fi
}

# 4. Run Suite
# Start with E2E
echo "--- Running E2E Scenarios ---"
for script in $(ls $SCENARIO_DIR/e2e_*.js | sort); do
    run_test "$script"
done

# Then Edge Cases
echo ""
echo "--- Running Edge Case Tests ---"
for script in $(ls $SCENARIO_DIR/edge_*.js | sort); do
    run_test "$script"
done

# 5. Summary
echo ""
echo "=================================================="
echo "Local Suite Results:"
echo "  Passed: $PASSED_COUNT"
echo "  Failed: $FAILED_COUNT"
echo "=================================================="

# 6. Report Generation
echo "Generating reports..."
if [ -f "./tests/ci_report.sh" ]; then
    ./tests/ci_report.sh > ci_report.md 2>/dev/null
    echo "  -> ci_report.md"
fi

if [ -f "./tests/edge_case_report.sh" ]; then
    ./tests/edge_case_report.sh > edge_case_report.md 2>/dev/null
    echo "  -> edge_case_report.md"
fi

if [ $FAILED_COUNT -gt 0 ]; then
    echo "FAILED: $FAILED_TESTS"
    exit 1
fi

exit 0
