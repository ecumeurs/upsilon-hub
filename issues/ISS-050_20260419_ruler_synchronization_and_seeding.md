# Issue: Ruler Synchronization Race and Seeding Fragility

**ID:** `20260419_ruler_synchronization_and_seeding`
**Ref:** `ISS-050`
**Date:** 2026-04-19
**Severity:** High
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/ruler`
**Affects:** `upsilonbattle/battlearena/controller`, `ruler_test.go`

---

## Summary

The `Ruler` component currently suffers from two interlinked stability issues: 
1. A race condition where controllers signal "BattleReady" before fully reconciling their local view of entities and turn order.
2. An over-reliance on fixed random seeding (`tools.SeedWith(42)`) in the test suite, which makes logic assertions brittle to any modification of the entity generator or turn-management architecture.

---

## Technical Description

### Background
The battle lifecycle requires a handshake:
1. `AddController` (Handshake Start)
2. `SetQueue` (Assignment of entities)
3. `ControllerBattleReady` (Handshake Complete)
4. `BattleStart` (First Turn Trigger)

### The Problem Scenario
In the current implementation of `FakeController` (used for testing), the `ControllerBattleReady` signal is sent immediately after receiving the `GetGridStateReply`. However, the `GetEntitiesStateReply` (which contains critical information about the `Turner` and assigned entities) may still be in flight. 

If the `Ruler` receives `ControllerBattleReady` from all controllers, it triggers `BattleStart` and the first turn. If the controller hasn't processed its entities yet, it may receive a `ControllerNextTurn` for an entity it doesn't "know" yet, leading to assertion failures or hangs in the test's event loop.

Furthermore, `TestRulerNextTurnSkipsDeadEntity` and others hardcode expected Entity IDs based on a fixed random seed. Any change to how `Turner` initializes (e.g., adding a property to entities) shifts the random sequence, causing these tests to fail even if the underlying logic is correct.

### Where This Pattern Exists Today
- `upsilonbattle/battlearena/ruler/ruler_test.go:184-192` (`GetGridStateReply` handler)
- `upsilonbattle/battlearena/ruler/ruler_test.go:47` (`init` seeding)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (intermittent CI failures) |
| Impact if triggered | High (battle hangs, incorrect turn assignment) |
| Detectability | Medium (manifests as timeouts or ID mismatch panics) |
| Current mitigant | Fixed seed (42) and manual turnaround checks in Ruler |

---

## Recommended Fix

**Short term:** Update `FakeController` to only signal `ControllerBattleReady` once both `GetGridState` and `GetEntitiesState` have successfully returned.  
**Medium term:** Identify which tests rely on turner for their logic and refactor them to be independent of the random seed; after all in these test we mostly don't care which entity plays, just that it can handle actions. For turn specific logic like forfeit or shotclock, we can access the turner prior deciding which entity plays next for our test and ensure it's appropriate for the test.

---

## References

- [ruler.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/ruler/ruler.go)
- [ruler_test.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/ruler/ruler_test.go)
- [ISS-046](file:///home/bastien/work/upsilon/projbackend/issues/ISS-046_dead_entity_hang.md) (Related regression)
