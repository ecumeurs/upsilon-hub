---
id: rule_code_health_monitoring
human_name: "Code Health Monitoring"
type: RULE
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 2
tags: [governance, linting, quality]
parents:
  - [[contract_upsilon_contract]]
dependents: []
---

# Code Health Monitoring

## INTENT
Ensure the codebase remains maintainable by enforcing strict limits on file size, function complexity, and documentation coverage while maintaining ATD traceability.

## THE RULE / LOGIC
Every source file in the project (Go, Python, PHP, JS, Vue) must adhere to the following health metrics:

1. **File Length (Lines of Code):**
   - Warning: > 300 LOC
   - Error: > 500 LOC
2. **Function Complexity:**
   - Nesting depth must not exceed 3 levels within a single function.
3. **Documentation Coverage (Function Intent):**
   - Every function must have at least one descriptive comment immediately preceding its definition.
   - **Exposed/Public Functions:** Error if missing.
   - **Private Functions:** Warning if missing.
   - **ATD Tags Exclusion:** Traceability tags (@spec-link, @test-link) do NOT satisfy this intent requirement.
   - **Body Exclusion:** Documentation inside the function body is ignored for this metric and generally discouraged in favor of clean logic.
4. **ATD Presence:**
   - Error: < 2 ATD links (@spec-link or @test-link).
   - Warning: > 5 ATD links (potential responsibility bloating).
   - Error: > 10 ATD links (violates Single Responsibility Principle).
5. **ATD Validity:**
   - All referenced ATD IDs must exist and be known to the ATD system.

**Exemptions:**
- Use `@lint-ignore-file-bloating` to skip file length checks.
- Use `@lint-ignore-complexity` to skip nesting checks.
- Use `@lint-ignore-documentation` to skip comment intent checks.
- Use `@lint-ignore-atd` to skip ATD presence/validity checks.

## TECHNICAL INTERFACE (The Bridge)
- **Script:** `scripts/code_health_check.py`
- **Hook:** Integrated into `scripts/pre-commit.sh`
- **Code Tag:** `@spec-link [[rule_code_health_monitoring]]`

## EXPECTATION (For Testing)
- Running `scripts/code_health_check.py` returns non-zero exit code if any "Error" threshold is breached.
- The script correctly identifies "phantom" ATD links that don't exist in `docs/`.
- The script respects `@lint-ignore-*` comments.
