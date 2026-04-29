# CI Failure Report - 2026-04-29

**Summary:** The local suite completed with **76 Passed** and **15 Failed** results. The failures are clustered around authentication infrastructure regressions, state machine desynchronization, and error-key drift.

## Failure Severity Overview

| Severity | Count | Primary Categories |
| :--- | :--- | :--- |
| **Critical** | 5 | Admin Authentication / Infrastructure |
| **High** | 6 | Match State Synchronization, Progression Logic |
| **Medium** | 4 | Error Key Drift, Validation Assertions |

---

## 1. Critical Failures: Authentication & Infrastructure

These failures block administrative operations and setup, rendering multiple E2E scenarios untestable.

### **[ISSUE] Admin Login: Invalid Credentials**
- **Affected Tests:**
  - `e2e_exotic_weapon_dual_path`
  - `e2e_friendly_fire_skill_test`
  - `e2e_item_grants_skill`
  - `edge_admin_skill_template_not_found`
- **Symptom:** `[CALL_ERROR] Route admin_login failed: Invalid credentials.`
- **Root Cause Analysis:** The `jsAdminSection` helper in the CLI defaults to `AdminPassword123!`. The "Creating CI Admin User" step in the CI setup appears to be out of sync with this default or the `UPSILON_ADMIN_PASSWORD` environment variable is misconfigured.
- **Impact:** Complete blockage of all administrative setup/teardown in tests.

### **[ISSUE] Session Desynchronization (Unauthenticated)**
- **Affected Test:** `e2e_admin_user_management`
- **Symptom:** `Route auth_logout failed: -- DEBUG MODE -- Unauthenticated.`
- **Root Cause Analysis:** The test likely lost its session token or the token expired/was invalidated before the logout call.

---

## 2. High Severity: State Machine & Engine Logic

These represent regressions in core game behavior or critical synchronization between the API bridge and the Go engine.

### **[ISSUE] Match State Desynchronization**
- **Affected Tests:**
  - `edge_attack_wrong_controller_with_2` ("Game is not in progress")
  - `edge_match_action_after_end_with_2` ("Match ended unexpectedly", "arena not found")
  - `edge_match_queue_while_in_match_with_2` ("Match state should still be accessible")
- **Root Cause Analysis:** Likely race conditions in the `GameController` or `Matchmaking` service where the arena is disposed of before the final actions are processed, or state transitions are not being propagated correctly to all agents.

### **[ISSUE] Progression Validation Bypass**
- **Affected Test:** `edge_prog_allocation_no_wins`
- **Symptom:** `Assertion Failed: HP changed after failed upgrade (Expected: 30, Actual: 31)`
- **Root Cause Analysis:** The engine's progression system is allowing attribute allocation even when the requirement (e.g., having wins) is not met. This is a logic regression in the `progression` package.

---

## 3. Medium Severity: Error Keys & Validation

These are "soft" failures where the system is behaving safely (rejecting illegal actions) but signaling the rejection with the wrong identifier or status code.

### **[ISSUE] Error Key Drift (ISS-080)**
- **Affected Test:** `edge_movement_path_not_adjacent`
- **Symptom:** `Expected entity.path.notadjacent, Actual: entity.path.notvalid`
- **Root Cause Analysis:** The engine has consolidated or renamed internal error keys. Per **ISS-080**, there is currently no canonical list or automated check for these keys, leading to drift between the engine and test assertions.

### **[ISSUE] Validation Assertion Mismatch**
- **Affected Tests:**
  - `edge_equip_unowned_character`
  - `edge_equip_unowned_item`
- **Symptom:** `Assertion Failed: Error must be 403 Forbidden or 404 Not Found`
- **Root Cause Analysis:** The API is likely returning a different status code (possibly 422 or 401) or the error envelope structure has changed in a way that the test's `assertResponse` helper no longer recognizes the rejection.

---

## Recommendations

1.  **Immediate Fix:** Reconcile CI Admin User password seeding with `upsiloncli` defaults.
2.  **Engine Audit:** Investigate the progression point allocation logic in `upsilonbattle/battlearena/progression`.
3.  **Traceability:** Align error keys for movement in the engine to restore test stability for `edge_movement_path_not_adjacent`.
4.  **Issue Updates:** Link these results to [ISS-080](file:///workspace/issues/ISS-080_20260425_error_key_atd_and_envelope.md) and [ISS-081](file:///workspace/issues/ISS-081_20260425_cross_stack_error_handling.md).
