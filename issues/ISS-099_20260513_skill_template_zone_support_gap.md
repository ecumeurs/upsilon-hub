# Issue: Limited Zone Targeting Support in Go Engine

**ID:** `20260513_skill_template_zone_support_gap`
**Ref:** `ISS-099`
**Date:** 2026-05-13
**Severity:** Medium
**Status:** Open
**Component:** `upsilontypes/property/def/skill.go`
**Affects:** `upsilonapi/bridge/bridge_utils.go`, `upsilontypes/property/def/skill.go`

---

## Summary

The current engine infrastructure lacks support for defining dynamic Area of Effect (AoE) zones via the skill targeting system. While the Laravel management interface correctly stores targeting data as a JSON object, the Go engine's `ZoneProperty` only supports hardcoded "Single" and "Neighbours" patterns. It fails to parse more complex definitions like "Circle:3" or "Square:2" provided by the bridge.

---

## Technical Description

### Background

AoE skills (Behavior: "Direct" or "Reaction" with a "Zone" targeting property) require a `ZoneProperty` to define which tiles are affected relative to the target. These patterns (Circle, Square, Line, etc.) are passed from Laravel as part of the `targeting` JSON field.

### The Problem Scenario

1.  A Skill Template is created in Laravel with `targeting` JSON containing `{"Zone": "Circle:3"}`.
2.  The `upsilonapi` bridge receives this payload and maps it to a `def.ZoneProperty`.
3.  The `ZoneProperty.Set(p interface{})` method in Go is too primitive:

```go
// upsilontypes/property/def/skill.go
func (bh *ZoneProperty) Set(p interface{}) {
	if s, ok := p.(string); ok {
		bh.PatternType = s
		if s == "Single" {
			bh.ZonePattern = pattern.Single()
		} else if s == "Neighbours" {
			bh.ZonePattern = pattern.Neighbours()
		} else {
			// FALLBACK: Silently fails to support custom radius/patterns
			bh.ZonePattern = pattern.Single()
		}
	}
}
```

4.  Any pattern string other than "Single" or "Neighbours" results in a single-target fallback, rendering the AoE definition useless.

### Where This Pattern Exists Today

-   `upsilontypes/property/def/skill.go:177` (Limited setter logic)
-   `upsilonapi/bridge/bridge_utils.go:94` (Targeting property map construction)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium — Prevents creation of standard AoE abilities (Fireball, Blast). |
| Detectability | High — Skills simply don't hit multiple targets as expected. |
| Current mitigant | None. Skills are limited to single-tile effects regardless of configuration. |

---

## Recommended Fix

**Short term:** Implement a basic string parser in `ZoneProperty.Set` to support "Circle:R" and "Square:R" using the existing `pattern` package.

**Medium term:** Refactor `ZoneProperty` to handle structured targeting payloads if we move beyond simple strings.

**Long term:** Sync the spatial pattern registry between Laravel and Go to ensure design-time validation matches simulation-time execution.

---

## Extra Data

The problem was identified during friendly fire testing, where an AoE skill with radius 6 was expected to skip allies, but only hit the primary target due to the engine fallback.

---

## References

- [skill.go](file:///workspace/upsilontypes/property/def/skill.go)
- [bridge_utils.go](file:///workspace/upsilonapi/bridge/bridge_utils.go)
- [pattern.go](file:///workspace/upsilonmapdata/grid/position/pattern/pattern.go)
