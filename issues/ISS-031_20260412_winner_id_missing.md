# Issue: Combat Concludes but WinnerID is Not Communicated to DTOs

**ID:** `20260412_winner_id_missing`
**Ref:** `ISS-031`
**Date:** 2026-04-12
**Severity:** High
**Status:** Open
**Component:** `upsilonbattle/battlearena/ruler/rules`
**Affects:** `upsilonapi/api`, `battleui`

---

## Summary

When a battle concludes naturally (all enemies eliminated or player forfeits), the engine transitions to the `Finished` state, but it fails to store the winning player's UUID in the `GameState`. Consequently, the API DTOs sent to clients always have a `null` `winner_id`, preventing bots and the WebUI from recognizing that the game has ended.

---

## Technical Description

### Background
The engine tracks battle state in `Ruler.CurrentState`. When a victory condition is met, the ruler should identify the winner and make this record available for state queries.

### The Problem Scenario
1. A battle ends (e.g., via `EndOfTurn` logic when only one controller remains).
2. `ruler.go` logs `##### END OF BATTLE! #####` and sets `CurrentState = Finished`.
3. `api.NewBoardState` is called to produce the JSON representation.
4. `BoardState.WinnerID` is left empty because `GameState` does not have a field to persist it, and the constructor doesn't populate it.
5. Clients receive `winner_id: null` and continue their battle loops indefinitely.

### Where This Pattern Exists Today
- [/workspace/upsilonbattle/battlearena/ruler/rules/gamestate.go](/workspace/upsilonbattle/battlearena/ruler/rules/gamestate.go:12-19) - Missing `WinnerID` field.
- [/workspace/upsilonapi/api/output.go](/workspace/upsilonapi/api/output.go:115-164) - `NewBoardState` does not populate `WinnerID`.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High — Battles Never End for clients |
| Detectability | High — Bots time out or loop forever |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Add `WinnerID` to `GameState` and populate it in `NewBoardState`.
**Medium term:** Ensure all state-changing rules update the `WinnerID` on transition to `Finished`.

---

## References

- [ruler.go](/workspace/upsilonbattle/battlearena/ruler/ruler.go)
- [output.go](/workspace/upsilonapi/api/output.go)
