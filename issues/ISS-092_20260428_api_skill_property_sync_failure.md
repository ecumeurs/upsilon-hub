# Issue: API Bridge Skill Property Sync Failure

**ID:** `20260428_api_skill_property_sync_failure`
**Ref:** `ISS-092`
**Date:** 2026-04-28
**Severity:** High
**Status:** Open
**Component:** `upsilonapi/bridge`, `upsilontypes/property/def`
**Affects:** Admin skill creation, skill-item application, engine-side targeting validation.

---

## Summary

The API bridge (`upsilonapi`) and the Go engine (`battlearena`) are currently out of sync regarding complex skill properties (`Range`, `TargetType`, `TargetingMechanics`). The `PropertyDTO` structure used for cross-service communication lacks the fields to represent complex structs, and the corresponding engine properties have empty `Set()` implementations, preventing external configuration from reaching the logic layer.

---

## Technical Description

### Background
Skills in the Upsilon engine use specialized property types to handle targeting logic:
- `RangeProperty` (handles Min/Max range)
- `TargetTypeProperty` (handles enums like `EnemyOnly`, `FriendOnly`)
- `TargetingMechanicsProperty` (handles `Line of Sight`, etc.)

### The Problem Scenario
1.  **Strict Typing Mismatch:** `api.PropertyDTO` only supports `Value`, `FValue`, `Max`, `BValue`, and `SValue`. It cannot naturally represent a `RangeProperty` which requires two distinct integer values (Min/Max).
2.  **Deaf Setters:** The `Set(interface{})` methods for these properties in `upsilontypes/property/def/skill.go` are currently empty.
3.  **Bridge Omission:** The `setSkillPropValue` helper in `upsilonapi/bridge/bridge.go` does not contain mapping logic for these non-primitive property types.

As a result, configuring a skill with `Range: 3` via the Admin API results in the engine defaulting to `Range(1, 1)`, causing unexpected "Target is not in range" errors in E2E tests even when valid ranges are provided.

### Where This Pattern Exists Today
- `upsilontypes/property/def/skill.go:249` (`RangeProperty.Set`)
- `upsilontypes/property/def/skill.go:388` (`TargetTypeProperty.Set`)
- `upsilonapi/bridge/bridge.go:513` (`setSkillPropValue`)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High — Admin-created skills won't behave as intended |
| Detectability | High — E2E tests fail on range/targeting validation |
| Current mitigant | Workaround: Move bots to Distance 1 in tests to bypass range checks |

---

## Recommended Fix

**Short term:** 
- Refactor `RangeProperty` to implement `IntCounterProperty` where `Value` maps to `MinRange` and `MaxValue` maps to `MaxRange`.
- Implement `Set(interface{})` in `TargetTypeProperty` and `TargetingMechanicsProperty` to accept and parse strings.

**Medium term:** 
- Move `TargetType` and `Mechanics` to use `DefaultStringProperty` with a factory-based enum validation layer to ensure content safety.
- Update `upsilonapi/bridge` to explicitly handle these specialized property mappings.

**Long term:** 
- Standardize the `Property` interface and `PropertyDTO` to support composite types or more flexible serialization patterns.

---

## References

- [upsilontypes/property/def/skill.go](file:///workspace/upsilontypes/property/def/skill.go)
- [upsilonapi/bridge/bridge.go](file:///workspace/upsilonapi/bridge/bridge.go)
- [e2e_friendly_fire_skill_test.js](file:///workspace/upsiloncli/tests/scenarios/e2e_friendly_fire_skill_test.js)
