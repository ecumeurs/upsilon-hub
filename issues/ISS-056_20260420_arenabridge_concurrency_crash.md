# Issue: ArenaBridge Concurrency Crash in GetBoardState

**ID:** `20260420_arenabridge_concurrency_crash`
**Ref:** `ISS-056`
**Date:** 2026-04-20
**Severity:** Critical
**Status:** In Progress
**Component:** `upsilonapi/bridge`
**Affects:** `upsilonapi`, `battleui`

---

## Summary

`ArenaBridge.GetBoardState` crashes with a `concurrent map iteration and map write` fatal error because it accesses the Ruler's GameState entities directly while the Ruler actor is modifying them.

---

## Technical Description

### Background
The `ArenaBridge` holds a reference to `arena` objects, which contain the `Ruler` instance. The `Ruler` owns the `GameState`.

### The Problem Scenario
1. The `Ruler` actor is processing a turn and modifying the `Entities` map.
2. A user makes a request that triggers `GetBoardState` (e.g. via `forwardToWebhook` or a direct GET request).
3. `GetBoardState` iterates over `arena.Ruler.GameState.Entities`.
4. Go runtime detects concurrent access and panics.

### Where This Pattern Exists Today
`upsilonapi/bridge/bridge.go:150`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Critical (Crashes the API) |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Implement a mutex-protected deep copy of the entity state or use a message-based approach to request the state from the Ruler (ensuring serial access via the actor's queue).  
**Medium term:** Strictly enforce ownership rules where the Ruler's internal state is Never accessible directly from outside its own actor loop.

---

## References

- [ui_investigation.md](file:///home/bastien/work/upsilon/projbackend/ui_investigation.md)
- [upsilonapi_crash/engine.log](file:///home/bastien/work/upsilon/projbackend/upsilonapi_crash/engine.log)
