# Issue: Ruler ownership bypass in bridge.go and public GameState

**ID:** `20260311_ruler_ownership_bypass`
**Ref:** `ISS-009`
**Date:** 2026-03-11
**Severity:** Low
**Status:** Open
**Component:** `github.com/ecumeurs/upsilonbattle/battlearena/ruler`
**Affects:** `github.com/ecumeurs/upsilonapi/bridge`

---

## Summary

In `bridge.go`'s `StartArena` function, the `Ruler`'s ownership of game resources is bypassed by directly manipulating its `GameState`. This violates the actor model's encapsulation rules, where state changes should occur through messages. Furthermore, `Ruler` exposes its `GameState` publicly, facilitating this bypass.

---

## Technical Description

### Background

The `Ruler` is implemented as an `Actor`. In a proper actor-based architecture, any modification to the actor's state (grid, entities, etc.) should be performed by sending messages to the actor's queue. Similarly, retrieving state should be done via requests to ensure consistency and prevent data races.

### The Problem Scenario

In `upsilonapi/bridge/bridge.go`, the `StartArena` function directly modifies the `Ruler`'s state:

```go
// bridge.go
battleArena.Ruler.SetGrid(gridgenerator.GeneratePlainSquare(10, 10)) // Direct call
battleArena.Ruler.SetNbControllers(len(start.Players))               // Direct call
...
battleArena.Ruler.AddEntity(e) // Direct call modifying GameState
...
return battleArena.Ruler.ID,
    battleArena.Ruler.GameState.Grid, // Direct access
    res,
    battleArena.Ruler.GameState.Turner.GetTurnState(), // Direct access
    nil
```

This bypasses the actor's message loop. While "okay" during initial setup when nothing is moving, it creates a precedent for breaking encapsulation.

### Where This Pattern Exists Today

- `upsilonapi/bridge/bridge.go:47-48`, `73`, `109-111`, `115-117`
- `upsilonbattle/battlearena/ruler/ruler.go:46` (`GameState` field is public)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (Data races, inconsistent state) |
| Detectability | Medium — manifests as non-deterministic bugs or crashes |
| Current mitigant | Limited to initial setup where concurrency is low |

---

## Recommended Fix

**Short term:** Add appropriate methods or notifications to `Ruler` to handle setting the grid and adding entities via messages. Use existing `GetState`, `GetGridState`, and `GetEntitiesState` calls for reply creation in `bridge.go`.

**Medium term:** Perform an impact analysis and make `Ruler.GameState` private. Ensure all external access is mediated via the actor's interface.

**Long term:** Enforce strict actor encapsulation across the entire system.

---

## References

- [bridge.go](file:///workspace/upsilonapi/bridge/bridge.go)
- [ruler.go](file:///workspace/upsilonbattle/battlearena/ruler/ruler.go)
