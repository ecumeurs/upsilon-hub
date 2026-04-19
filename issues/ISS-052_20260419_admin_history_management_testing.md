# Issue: Implement Full E2E Admin Test Suite

**ID:** `20260419_admin_history_management_testing`
**Ref:** `ISS-052`
**Date:** 2026-04-19
**Severity:** Medium
**Status:** Open
**Component:** `upsiloncli/tests/scenarios/`
**Affects:** CI Coverage, Feature Stability, Admin Security

---

## Summary

The Administrative suite currently lacks comprehensive E2E validation. While some logic and UI assets exist, most administrative use cases (`uc_admin_login`, `ui_admin_dashboard`, `uc_admin_history_management`) are not verified in the CI pipeline.

This issue tracks the creation of a unified or modular suite of scenarios that verify the entire administrative lifecycle: from initial seeding and login to dashboard navigation and history maintenance.

---

## Technical Description

### Background
Currently, ONLY **Admin User Management (CR-15)** is tracked and reported in the E2E CI compliance matrix. The following atoms are "Testing Orphans":
1. `uc_admin_login`: Authentication logic exists but no E2E flow validates it.
2. `ui_admin_dashboard`: Visibility and role-based redirection are unverified.
3. `uc_admin_history_management`: No tests (dependency on [[ISS-051]]).
4. `infra_seed_admin`: No E2E proof that the admin user seeded in CI is functional.

### The Problem Scenario
A developer might break the Admin Login redirect or the Dashboard security middleware, and CI would remain Green because no scenario exercises these routes with an elevated privilege token.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | High (Security bypass, broken admin tools) |
| Detectability | Low (Until manual inspection or production failure) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** 
1. Create `e2e_admin_full_lifecycle.js` or separate specific scenarios for:
   - **Login Verification:** Assert successful login with seeded `admin` credentials.
   - **Dashboard Access:** Verify authorized access to `/admin/dashboard`.
   - **Access Denial:** Verify non-admin users/guests are blocked from admin routes (using `edge_auth_non_admin_access.js`).
2. Implement `e2e_admin_history_management.js` once history endpoints are ready.

**Medium term:** Update `tests/ci_report.sh` to include CR-mappings for these new tests.

---

## References

- [uc_admin_history_management.atom.md](../../docs/uc_admin_history_management.atom.md)
- [e2e_admin_user_management.js](../../upsiloncli/tests/scenarios/e2e_admin_user_management.js)
