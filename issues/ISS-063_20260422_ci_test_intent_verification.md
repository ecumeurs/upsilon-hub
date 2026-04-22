# Issue: CI Test Intent Verification Mismatch

**ID:** `20260422_ci_test_intent_verification`
**Ref:** `ISS-063`
**Date:** 2026-04-22
**Severity:** Medium
**Status:** Resolved
**Component:** `tests/`, `CI.md`, `upsiloncli/tests/scenarios/`
**Affects:** CI pipeline reporting, test coverage tracking, and documentation accuracy

---

## Summary

CI.md claims 100% implementation of all 48 edge case tests and 17 E2E customer scenarios, but the actual CI reporting scripts (edge_case_report.sh and ci_report.sh) contain outdated hardcoded coverage statistics and fail to reflect the true test implementation status. This creates a false sense of test coverage and prevents accurate tracking of CI test compliance.

---

## Technical Description

### Background

The CI testing framework is documented in CI.md with two main test categories:
- **Edge Case Tests**: 48 tests validating API boundaries, validation rules, and error handling
- **E2E Customer Scenarios**: 17 scenarios mapping to business requirements from the Conformity Matrix

Test scripts are automatically discovered in `upsiloncli/tests/scenarios/` using prefixes:
- `edge_*.js` for edge case tests
- `e2e_*.js` for end-to-end customer scenarios

### The Problem Scenario

**Mismatch 1: Hardcoded Coverage Statistics**

In `tests/edge_case_report.sh` (lines 186-196), coverage statistics are hardcoded:
```bash
count_category 9 2 "Movement Validation"  # Claims only 2 of 9 implemented
count_category 10 2 "Attack Validation"  # Claims only 2 of 10 implemented
count_category 6 0 "Character & Progression"  # Claims 0 of 6 implemented
```

However, CI.md (line 106) states: "**Implementation Progress**: 48/48 tests (100%) ✅ All tests implemented"

**Mismatch 2: Undocumented E2E Tests**

CI.md documents 17 E2E customer scenarios (CR-01 to CR-17), but `upsiloncli/tests/scenarios/` contains 19 E2E test files:
- 17 documented in CI.md ✅
- 2 undocumented: `e2e_admin_full_lifecycle.js`, `e2e_admin_history_management.js`

**Mismatch 3: Tests Are Actually Implemented**

Investigation shows that:
- All 48 edge case test files exist with proper `@spec-link` tags
- All 17 documented E2E test files exist with proper `@spec-link` tags
- The tests appear to be well-implemented with appropriate assertions

The issue is NOT missing tests, but **inaccurate reporting** of test status.

### Where This Pattern Exists Today

| File | Lines | Issue |
|---|---|---|
| `CI.md` | 53-106 | Claims 100% implementation but doesn't verify against actual files |
| `tests/edge_case_report.sh` | 186-196 | Hardcoded coverage numbers instead of dynamic calculation |
| `tests/ci_report.sh` | 64-81 | Hardcoded mapping of 17 scenarios (missing 2 undocumented tests) |

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High — CI reports are generated on every run |
| Impact if triggered | Medium — False sense of security, unclear test status, misleading coverage metrics |
| Detectability | Low — Reports appear valid until manually audited |
| Current mitigant | None; developers must manually verify test existence |

**Specific Risks:**
1. **False Confidence**: CI reports showing 0% implementation in categories where tests exist
2. **Undocumented Tests**: 2 E2E tests not tracked in CI compliance matrix
3. **Maintenance Burden**: Adding new tests requires manual updates to hardcoded numbers
4. **Audit Trail Break**: No reliable way to verify which tests are actually implemented

---

## Tasks

### 1. Proof Read All CI Tests for Intent Appropriateness

**Priority:** High  
**Owner:** Unassigned  
**Status:** Completed  
**Completion Date:** 2026-04-22  
**Findings:** See [test_quality_review_report.md](test_quality_review_report.md)

**Description:**
Systematically review all 48 edge case tests and 19 E2E scenarios to verify that:
- Test implementations actually validate the documented intent
- Assertions are comprehensive and appropriate for the test purpose
- Edge cases truly cover boundary conditions, not just happy paths
- Customer scenarios accurately represent real-world usage patterns
- `@spec-link` tags reference the correct ATD atoms
- Test descriptions match their actual behavior

**Approach:**
1. Read each test file in `upsiloncli/tests/scenarios/`
2. Compare test implementation against:
   - EC-ID or CR-ID documentation in CI.md
   - Referenced ATD atom specifications
   - Business requirement intent
3. Document any mismatches, insufficient assertions, or missing coverage
4. Create follow-up issues for any test quality problems found

**Expected Deliverables:**
- Test review matrix showing each test's alignment score
- List of tests that need improvement or rewriting
- Updated test files with better assertions and documentation
- Corrected `@spec-link` references where needed

---

## Recommended Fix

**Short term (Documentation & Scripts):**
1. Update `tests/edge_case_report.sh` to dynamically calculate coverage from actual test files:
   ```bash
   # Instead of: count_category 9 2 "Movement Validation"
   # Use: count_category 9 $(find upsiloncli/tests/scenarios -name "edge_movement_*.js" | wc -l) "Movement Validation"
   ```
2. Update CI.md to reflect actual test count (19 E2E scenarios, not 17)
3. Add the 2 undocumented E2E tests to CI.md's Customer Requirement Mapping table

**Medium term (Automation):**
1. Modify report generators to auto-discover test files and count implementations
2. Add validation step in CI that checks CI.md matches actual test files
3. Generate CI.md tables automatically from test file metadata

**Long term (Architecture):**
1. Implement a test registry system where tests self-report their metadata
2. Use test file frontmatter or JSDoc comments for ID/Category mapping
3. Create a single source of truth for test definitions that CI.md, report generators, and CI workflows all reference

---

## References

- [CI.md](../CI.md) — Main CI documentation with outdated test counts
- [tests/edge_case_report.sh](../tests/edge_case_report.sh) — Edge case report generator with hardcoded statistics
- [tests/ci_report.sh](../tests/ci_report.sh) — E2E report generator with incomplete mapping
- [upsiloncli/tests/scenarios/](../upsiloncli/tests/scenarios/) — Actual test implementation directory
- [.agent/rules/issues.md](../.agent/rules/issues.md) — Issue filing procedure reference
