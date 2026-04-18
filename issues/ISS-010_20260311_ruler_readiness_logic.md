# Issue: Ruler readiness trigger enhancements

**ID:** `20260311_ruler_readiness_logic`
**Ref:** `ISS-010`
**Date:** 2026-03-11
**Severity:** Low
**Status:** Resolved
**Component:** `github.com/ecumeurs/upsilonbattle/battlearena/ruler`
**Affects:** `github.com/ecumeurs/upsilonapi/bridge`

---

## Summary

The current readiness trigger for the `Ruler` (the `BattleStart` notification) only checks if the required number of controllers are connected. It does not verify that other essential game components, such as entities being registered/assigned and a grid being present, are correctly initialized.

---

## Technical Description

### Background

The `Ruler` waits for a specific state before starting the battle. Currently, this state is solely defined by the number of registered controllers reaching `NbControllers`.

### The Problem Scenario

In `upsilonbattle/battlearena/ruler/ruler.go`:

```go
// ruler.go:201
if len(r.GameState.Controllers) == r.NbControllers {
    r.NotifyActor(message.Create(nil, rulermethods.BattleStart{}, nil))
}
```

If controllers connect but the grid isn't set or entities haven't been assigned (e.g., due to an error in `bridge.go` or a race condition), the battle might start in an invalid state.

### Where This Pattern Exists Today

- `upsilonbattle/battlearena/ruler/ruler.go:201-203`
- `upsilonbattle/battlearena/ruler/ruler.go:210` (Readiness check for entity turn)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium (Panic or invalid game state at start) |
| Detectability | High — game will fail to start correctly |
| Current mitigant | Bridge logic usually sets these up before controller registration |

---

## Recommended Fix

**Short term:** Enhance the conditional check in `addController` and `controllerBattleReady` to also verify `r.GameState.Grid != nil` and that the number of registered entities matches expectations.

**Medium term:** Implement a formal "Readiness" state machine or check function that encapsulates all requirements for starting a battle.

**Long term:** N/A

---

## References

- [ruler.go](file:///workspace/upsilonbattle/battlearena/ruler/ruler.go)
