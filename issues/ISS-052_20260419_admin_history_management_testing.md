# Issue: Implement E2E Tests for Admin Match History Management

**ID:** `20260419_admin_history_management_testing`
**Ref:** `ISS-052`
**Date:** 2026-04-19
**Severity:** Medium
**Status:** Open
**Component:** `upsiloncli/tests/scenarios/`
**Affects:** CI Coverage, Feature Stability

---

## Summary

The Administrative Match History Management feature lacks any E2E test coverage. Even after the backend is implemented, we need automated validation to ensure administrators can successfully list and purge history without regressions.

---

## Technical Description

### Background
The ATD specification `uc_admin_history_management` defines two test scenarios:
1. `TestAdminListMatches`: Verify an admin can see the match history.
2. `TestAdminPurgeHistory`: Verify an admin can delete old records.

### The Problem Scenario
There are currently no `.js` scenario files in `upsiloncli/tests/scenarios/` that target these functionalities. Consequently, the CI pipeline does not verify these administrative rules.

### Where This Pattern Exists Today
- **Existing Admin Tests:** `upsiloncli/tests/scenarios/e2e_admin_user_management.js`
- **Missing Tests:** `e2e_admin_history_management.js` and potentially `edge_admin_history_purge_empty.js`.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium (Regressions in admin tools, unverified security boundaries) |
| Detectability | High (Missing from CI reports) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Create a new E2E scenario `e2e_admin_history_management.js` that:
1. Logs in as a seeded admin.
2. Checks for presence of match history.
3. Triggers a purge and verifies records are removed (using a mock date if necessary).

**Medium term:** Add edge cases for unauthorized access to history (non-admin users).

---

## References

- [uc_admin_history_management.atom.md](../../docs/uc_admin_history_management.atom.md)
- [e2e_admin_user_management.js](../../upsiloncli/tests/scenarios/e2e_admin_user_management.js)
