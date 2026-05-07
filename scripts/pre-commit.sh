#!/bin/bash

# Upsilon Pre-commit Verification Script
# This script runs all checks required for CI to pass.

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=======================================${NC}"
echo -e "${BLUE}   Upsilon CI Pre-check Validator      ${NC}"
echo -e "${BLUE}=======================================${NC}"

# TRACKING
GO_SYNC="SKIP"
GO_VET="FAIL"
GO_TEST="FAIL"
PHP_TEST="FAIL"
HEALTH_CHECK="FAIL"

# 0. CODE HEALTH CHECK
echo -e "\n${YELLOW}[0/5] Running Code Health Check (DISABLED)...${NC}"
# if ./scripts/code_health_check.py; then
#     HEALTH_CHECK="PASS"
#     echo -e "${GREEN}✓ Code health standards met${NC}"
# else
#     HEALTH_CHECK="FAIL"
#     echo -e "${RED}✗ Code health check failed${NC}"
# fi
HEALTH_CHECK="SKIP"

# 1. GO WORKSPACE SYNC
echo -e "\n${YELLOW}[1/5] Syncing Go Workspace...${NC}"
if go work sync; then
    GO_SYNC="PASS"
    echo -e "${GREEN}✓ Workspace synchronized${NC}"
else
    GO_SYNC="FAIL"
    echo -e "${RED}✗ Workspace sync failed${NC}"
fi

# 2. GO VET
echo -e "\n${YELLOW}[2/5] Running Go Vet...${NC}"
MODULES="./upsilonapi/... ./upsiloncli/... ./upsilonbattle/... ./upsilonmapdata/... ./upsilonmapmaker/... ./upsilonserializer/... ./upsilontools/..."
if go vet $MODULES; then
    GO_VET="PASS"
    echo -e "${GREEN}✓ Go Vet passed${NC}"
else
    echo -e "${RED}✗ Go Vet found issues${NC}"
fi

# 3. GO TEST
echo -e "\n${YELLOW}[3/5] Running Go Unit Tests...${NC}"
if go test -timeout 30s $MODULES; then
    GO_TEST="PASS"
    echo -e "${GREEN}✓ Go Tests passed${NC}"
else
    echo -e "${RED}✗ Go Tests failed${NC}"
fi

# 4. PHP TESTS (Conditional)
echo -e "\n${YELLOW}[4/5] Running PHP Tests (BattleUI)...${NC}"
if [ -d "battleui" ]; then
    if command -v php >/dev/null 2>&1 && command -v composer >/dev/null 2>&1; then
        cd battleui
        if [ -f "artisan" ]; then
            if php artisan test --parallel; then
                PHP_TEST="PASS"
                echo -e "${GREEN}✓ PHP Tests passed${NC}"
            else
                echo -e "${RED}✗ PHP Tests failed${NC}"
            fi
        else
            echo -e "${YELLOW}! No artisan file found in battleui${NC}"
            PHP_TEST="MISSING"
        fi
        cd ..
    else
        echo -e "${YELLOW}! PHP or Composer not found (Skipping tests)${NC}"
        PHP_TEST="NO_ENV"
    fi
else
    echo -e "${YELLOW}! battleui directory not found${NC}"
    PHP_TEST="MISSING"
fi

# SUMMARY TABLE
echo -e "\n${BLUE}=======================================${NC}"
echo -e "${BLUE}           SUMMARY REPORT              ${NC}"
echo -e "${BLUE}=======================================${NC}"

format_status() {
    if [ "$1" == "PASS" ]; then
        echo -e "${GREEN}PASS${NC}"
    elif [ "$1" == "NO_ENV" ] || [ "$1" == "MISSING" ]; then
        echo -e "${YELLOW}$1${NC}"
    else
        echo -e "${RED}FAIL${NC}"
    fi
}

echo -e "Code Health Check : $(format_status $HEALTH_CHECK)"
echo -e "Go Workspace Sync : $(format_status $GO_SYNC)"
echo -e "Go Linting (Vet)  : $(format_status $GO_VET)"
echo -e "Go Unit Tests     : $(format_status $GO_TEST)"
echo -e "PHP Unit Tests    : $(format_status $PHP_TEST)"
echo -e "${BLUE}=======================================${NC}"

if [ "$HEALTH_CHECK" == "PASS" ] && [ "$GO_SYNC" == "PASS" ] && [ "$GO_VET" == "PASS" ] && [ "$GO_TEST" == "PASS" ]; then
    echo -e "\n${GREEN}READY TO COMMIT (Go checks passed)${NC}"
    exit 0
else
    echo -e "\n${RED}PLEASE FIX ERRORS BEFORE COMMITTING${NC}"
    exit 1
fi
