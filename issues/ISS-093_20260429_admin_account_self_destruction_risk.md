# Issue: Administrative Account Self-Destruction Risk

**ID:** `20260429_admin_account_self_destruction_risk`
**Ref:** `ISS-093`
**Date:** 2026-04-29
**Severity:** Critical
**Status:** Open
**Component:** `battleui/app/Http/Controllers/API/AdminController.php` (or equivalent Laravel admin handler)
**Affects:** All administrative functionality and CI stability.

---

## Summary

The administrative API allows users with admin privileges to target any account for anonymization or deletion, including their own. When an admin self-anonymizes, their session is invalidated and their credentials (account name/password) are overwritten, effectively locking the system out of administrative access. This is currently causing cascading failures in the CI suite when tests accidentally target the `admin` user.

---

## Technical Description

### Background

Administrative users can manage the user registry via the following routes:
- `POST /api/v1/admin/users/{account_name}/anonymize`
- `DELETE /api/v1/admin/users/{account_name}`

These routes are intended for GDPR compliance and system maintenance.

### The Problem Scenario

1.  A test script (e.g., `e2e_admin_user_management.js`) logs in as `admin`.
2.  The script fetches the user list.
3.  The script identifies the "last user" in the list as a target for anonymization.
4.  If the registry is small or sorted such that `admin` is last, the script calls `admin_user_anonymize("admin")`.
5.  The backend processes the request:
    - Overwrites `admin` PII and credentials.
    - Invalidates active sessions for `admin`.
6.  The test script fails at logout (Unauthenticated).
7.  All subsequent tests fail at login (Invalid Credentials).

### Where This Pattern Exists Today

- **Scripting Layer:** `upsiloncli/tests/scenarios/e2e_admin_user_management.js` (targets the last user without checking if it is the admin).
- **Backend Layer:** The controller handling `admin_user_anonymize` does not check if `target_user == current_user` or if the `target_user` has protected roles.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High — currently triggered in CI. |
| Impact if triggered | Critical — total loss of administrative access. |
| Detectability | High — manifests as `Invalid credentials` for admin login. |
| Current mitigant | None. |

---

## Recommended Fix

**Short term:**
Add a guard clause in the backend controller to prevent an admin from targeting themselves.
```php
if ($targetUser->id === auth()->id()) {
    return response()->error("Self-destruction is not permitted.", 403);
}
```

**Medium term:**
Introduce a "Protected Account" flag or role-based restriction where users with the `SuperAdmin` role cannot be anonymized or deleted via the standard admin API (requires a manual DB intervention or a multi-admin approval flow).

**Long term:**
Implement a multi-signature or approval-based system for destructive administrative actions on privileged accounts.

---

## References

- [e2e_admin_user_management.js](file:///workspace/upsiloncli/tests/scenarios/e2e_admin_user_management.js)
- [admin.go](file:///workspace/upsiloncli/internal/endpoint/admin.go)
- [ISS-081](file:///workspace/issues/ISS-081_20260425_cross_stack_error_handling.md) (related to error propagation)
