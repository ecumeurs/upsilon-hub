# Issue: Internal User ID Exposure in Battle State DTOs

**ID:** `Ref_20260513_battle_engine_internal_user_id_leak`
**Ref:** `ISS-098`
**Date:** 2026-05-13
**Severity:** High
**Status:** Open
**Component:** `upsilonapi/api`
**Affects:** `battleui` frontend, E2E tests, and any public API consumer observing the board state.

---

## Summary

The Go battle engine's API bridge currently leaks internal User UUIDs via the `player_id` field in tactical `Entity` DTOs. This violates the `[[requirement_customer_user_id_privacy]]` mandate which prohibits exposing database primary keys to the client to prevent primary key enumeration and protect user privacy.

---

## Technical Description

### Background
According to `[[requirement_customer_user_id_privacy]]`, all internal User identifiers must be masked. The frontend should only receive masked identifiers (like nicknames or session-specific keys) rather than raw database UUIDs.

### The Problem Scenario
When the `upsilonapi` bridge prepares the `BoardState` DTO to send to Laravel (and subsequently to the frontend via WebSockets), it maps the engine's internal `ControllerID` directly to the JSON `player_id` field.

1. Engine executes a turn or event.
2. `upsilonapi/api/output.go` constructs a new `BoardState`.
3. `NewEntity` function is called for each tactical unit.
4. The internal `ControllerID` (a raw User UUID) is assigned to `PlayerID`.

```go
// upsilonapi/api/output.go:208
res := &Entity{
    ID:             ent.ID.String(),
    PlayerID:       ent.ControllerID.String(), // LEAK: Internal User UUID
    Team:           int(ent.Team),
    // ...
}
```

### Where This Pattern Exists Today
- [upsilonapi/api/output.go:208](file:///workspace/upsilonapi/api/output.go#L208)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High — revealed in every board state WebSocket message or API poll |
| Current mitigant | None. The DTO field is populated unconditionally. |

---

## Recommended Fix

**Short term:** Implement an `IdMasker` utility in `upsilonapi` that maps internal `ControllerID`s to a transient, match-specific `TacticalPlayerID` or simply use the Player's Team/Index if unique.

**Medium term:** Update `NewEntity` to use the masked ID instead of `ent.ControllerID.String()`. Ensure the `BoardStateResource` in Laravel also respects this masking.

**Long term:** Align with `[[arch_api_id_masking_gateway]]` to ensure all cross-service boundaries perform automatic ID transformation.

---

## References

- [requirement_customer_user_id_privacy.atom.md](file:///workspace/docs/requirement_customer_user_id_privacy.atom.md)
- [upsilonapi/api/output.go](file:///workspace/upsilonapi/api/output.go)
- [e2e_battle_starts_privacy_check.js](file:///workspace/upsiloncli/tests/scenarios/e2e_battle_starts_privacy_check.js)
