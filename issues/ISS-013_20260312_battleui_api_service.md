# Issue: BattleUI API Communication Service Ownership

**ID:** `20260312_battleui_api_service`
**Ref:** `ISS-013`
**Date:** 2026-03-12
**Severity:** High
**Status:** Resolved
**Component:** `battleui/app/Services/UpsilonApiService.php`
**Affects:** `GameController`, `MatchMakingController`, `BattleController`

---

## Summary

The `battleui` component currently lacks a centralized service to handle communication with the `upsilonapi` (Go Battle Engine). Communication is currently ad-hoc or missing in placeholder controllers. A dedicated `UpsilonApiService` must own this responsibility to ensure consistent DTO mapping, robust error handling, and correct unpacking of API responses into `battleui` entities.

---

## Technical Description

### Background
The `battleui` gateway acts as a proxy for the `upsilonapi`. It needs to send commands (start arena, action) and receive state updates. These interactions involve JSON payloads that must match the Go API's expectations.

### The Problem Scenario
Currently, controllers like `GameController.php` are empty. When implementing them, there is a risk of:
1. Fragmented API logic duplicated across controllers.
2. Inconsistent DTO (API Resource) structures compared to the Go `api` package.
3. Brittle response handling that doesn't validate if Go's `BoardState` matches Laravel's expected entities.

```
[BattleUI Controller] --(Untyped JSON)--> [UpsilonAPI]
[BattleUI Controller] <--(Raw Response)-- [UpsilonAPI]
(Risk: Type mismatch or lack of validation)
```

### Where This Pattern Exists Today
- `battleui/app/Http/Controllers/API/GameController.php` (Empty)
- `battleui/app/Http/Controllers/API/MatchMakingController.php` (Empty)
- `battleui/routes/api.php` (Points to empty methods)

---

## Risk Assessment

| Factor              | Value                                                            |
| ------------------- | ---------------------------------------------------------------- |
| Likelihood          | High                                                             |
| Impact if triggered | High                                                             |
| Detectability       | Medium — manifests as runtime errors or inconsistent game state. |
| Current mitigant    | None (Placeholders only).                                        |

---

## Recommended Fix

**Short term:** Define the `UpsilonApiService` interface and create DTO classes (Laravel Resources/Data Objects) that mirror the Go API structures defined in `api_go_battle_start.atom.md` and `api_go_battle_action.atom.md`.

**Medium term:** Implement the `UpsilonApiService` using Laravel's `Http` client, ensuring all communication is wrapped in the `api_standard_envelope`.

**Long term:** Automate DTO synchronization or use a shared specification (like TypeSpec) to generate clients.

---

## References

- [api_go_battle_engine.atom.md](file:///workspace/docs/api_go_battle_engine.atom.md)
- [api_go_battle_start.atom.md](file:///workspace/docs/api_go_battle_start.atom.md)
- [api_go_battle_action.atom.md](file:///workspace/docs/api_go_battle_action.atom.md)
- [api_go_webhook_callback.atom.md](file:///workspace/docs/api_go_webhook_callback.atom.md)
- [GameController.php](file:///workspace/battleui/app/Http/Controllers/API/GameController.php)
- [ISS-014](file:///workspace/issues/ISS-014_20260312_battleui_db_models.md) (Missing Database Models in BattleUI)
