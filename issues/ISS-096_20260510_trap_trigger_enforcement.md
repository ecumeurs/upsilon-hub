# Issue: Trap Trigger Type Enforcement

**ID:** `20260510_trap_trigger_enforcement`
**Ref:** `ISS-096`
**Date:** 2026-05-10
**Severity:** Medium
**Status:** Open
**Component:** `upsilonbattle/battlearena/ruler/rules`
**Affects:** `positionaleffect.go`, `ProcessPositionalEffects`

---

## Summary

Traps (positional effects) currently fail silently if they are missing a `TriggerType` property. This makes misconfigured traps difficult to debug as they simply never fire without any log or error feedback.

---

## Technical Description

### Background
Positional effects are attached to grid cells and intended to fire when specific triggers occur (e.g., `OnEnter`, `OnExit`, `OnTurn`). The matching logic depends on the presence of a `property.TriggerType` SkillProperty on the effect.

### The Problem Scenario
In `positionaleffect.go`, the `processSinglePositionalEffect` function performs the following check:

```go
triggerProp := eff.GetProperty(property.TriggerType)
if triggerProp == nil {
    return
}
```

If an effect is added to the grid but is missing the `TriggerType` property (due to a configuration error in the map data or a skill definition), it will be ignored during processing. There is no error log indicating that an effect was skipped due to missing metadata.

### Where This Pattern Exists Today
- [positionaleffect.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/positionaleffect.go#L47-L50)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium |
| Detectability | Low — traps just "don't work" without explanation |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Add a debug or warning log in `processSinglePositionalEffect` when an effect is found at a position but lacks a `TriggerType`.  
**Medium term:** Implement a validation layer when adding effects to the grid (`PositionalEffects` map) to ensure they have the mandatory trigger metadata.  
**Long term:** Use a structured `Trap` type or specialized factory that enforces these properties at the type level.

---

## Extra Data

Discovered during code health remediation and investigation of `rules_iss066_test.go`.

---

## References

- [positionaleffect.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/positionaleffect.go)
- [rules_iss066_test.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/rules_iss066_test.go)
