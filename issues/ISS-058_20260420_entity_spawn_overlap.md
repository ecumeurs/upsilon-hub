# Issue: Entity Spawning Overlap

**ID:** `20260420_entity_spawn_overlap`
**Ref:** `ISS-058`
**Date:** 2026-04-20
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi/ruler`
**Affects:** `upsilonapi`, `battleui`

---

## Summary

In some cases, multiple entities are spawned on the same tile at the start of a match, which violates game rules and causes visual glitching.

---

## Technical Description

### Background
The `Ruler` or `EntityGenerator` is responsible for placing entities on the grid during match initialization.

### The Problem Scenario
Observed in `ui_investigation.md`: "two entity spawnned at the same tile."

### Where This Pattern Exists Today
Logic responsible for initial entity placement in `upsilonapi`.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Low |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Add a validation check during spawning to ensuring the target tile is empty.  
**Medium term:** Refactor spawning logic to use a shuffle-based placement on available tiles.

---

## References

- [ui_investigation.md](file:///home/bastien/work/upsilon/projbackend/ui_investigation.md)
