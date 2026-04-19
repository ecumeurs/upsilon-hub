# Issue: Implement Administrative Match History Management (Backend & Frontend)

**ID:** `20260419_admin_history_management_impl`
**Ref:** `ISS-051`
**Date:** 2026-04-19
**Severity:** Medium
**Status:** Open
**Component:** `battleui/app/Http/Controllers/API/AdminController.php`, `battleui/resources/js/Pages/Admin/`
**Affects:** Admin Dashboard, Database Performance, Admin UI

---

## Summary

The Administrative Match History Management use case (`uc_admin_history_management`) is documented in ATD but has no implementation in the backend or frontend. Administrators currently lack the tools (API and UI) to review match outcomes or purge old history.

This issue also covers the verification and potential refinement of the **Admin Dashboard** and **Admin User Management** frontend pages to ensure they are fully operational and correctly linked to the new history management views.

---

## Technical Description

### Background
The system is required to provide administrators with a unified management interface:
1. **Match History UI:** A dedicated view to list all completed matches.
2. **Purge Action:** A frontend button to trigger history purging (older than 90 days).
3. **Operational Dashboards:** Ensure `Admin/Dashboard.vue` and `Admin/UserManagement.vue` are integrated with the system logic.

### The Problem Scenario
- An administrator logs into the dashboard and attempts to access the "Match History" section.
- The request fails or shows an empty/unimplemented state because the API endpoints defined in the ATD (`GET /admin/history`, `POST /admin/history/purge`) do not exist.

### Where This Pattern Exists Today
- **ATD Atom:** `uc_admin_history_management.atom.md`
- **Related Controller:** `battleui/app/Http/Controllers/API/AdminController.php` (already handles user management but lacks history methods).

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (Database bloat, lack of auditability) |
| Detectability | High (Feature simply missing from UI/API) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** 
- Implement `GET /api/admin/history` and `POST /api/admin/history/purge` in `AdminController.php`.
- Create `battleui/resources/js/Pages/Admin/History.vue` to display the match list and provide the purge button.

**Medium term:** 
- Update `Admin/Dashboard.vue` with a navigation link to the History page.
- Audit `Admin/UserManagement.vue` for functional parity with current backend logic.

**Long term:** 
- Automate history purging via a scheduled Laravel task.

---

## References

- [uc_admin_history_management.atom.md](../../docs/uc_admin_history_management.atom.md)
- [AdminController.php](../../battleui/app/Http/Controllers/API/AdminController.php)
