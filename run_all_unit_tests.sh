#!/bin/bash

# run_all_unit_tests.sh - Comprehensive Unit Test Runner for UpsilonBattle

# Setup colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "--- Running All Unit Tests ---\n"

# 1. Run Go Tests
echo ">>> Executing Go Tests..."
# We target specific modules to ensure coverage across the workspace
go test ./upsilonapi/... ./upsilonbattle/... ./upsiloncli/... ./upsilonmapdata/... ./upsilonmapmaker/... ./upsilonserializer/... ./upsilontools/... ./upsilontypes/...
GO_EXIT=$?

# 2. Run PHP Tests
echo -e "\n>>> Executing PHP Tests (Excluding External API Roundtrips)..."
cd battleui || exit 1
# Exclude Roundtrip tests which require the Go engine to be running as a service
php artisan test --exclude-filter UpsilonApiRoundtripTest
PHP_EXIT=$?
cd ..

# 3. Final Sanctions
echo -e "\n=== SANCTIONS ==="

if [ $GO_EXIT -eq 0 ]; then
    echo -e "GO: ${GREEN}PASSED${NC}"
else
    echo -e "GO: ${RED}FAILED${NC}"
fi

if [ $PHP_EXIT -eq 0 ]; then
    echo -e "PHP: ${GREEN}PASSED${NC}"
else
    echo -e "PHP: ${RED}FAILED${NC}"
fi

# Exit with success only if both passed
if [ $GO_EXIT -eq 0 ] && [ $PHP_EXIT -eq 0 ]; then
    exit 0
else
    exit 1
fi
