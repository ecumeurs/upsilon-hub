# Issue: Extract Properties and Skill Weight to Shared Library

**ID:** `20260425_extract_properties_shared_library`
**Ref:** `ISS-085`
**Date:** 2026-04-25
**Severity:** Medium
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/property`
**Affects:** `upsilonapi`, `upsilonbattle`, `upsiloncli`, `skillgenerator`

---

## Summary

Extract the core property system and skill weight mathematical framework from the `upsilonbattle` engine into a standalone shared Go library. This architectural move decouple "Combat Execution" from "World/Meta Generation," enabling specialized services (like a Token-based Skill Crafting API) to use engine-compliant math without importing the entire battle simulation logic.

---

## Technical Description

### Background
Currently, all combat arithmetic rules (Properties) and balancing algorithms (Skill Weight) live inside the `upsilonbattle` repository. This was efficient during early development, but creates a "DRY" violation as we expand the meta-game.

### The Problem Scenario
1. **World Meta-Game**: We want a "Regional Generation API" that crafts skills based on location and player "Tokens."
2. **Dependency Bloat**: For this API to ensure skills are balanced, it must currently import `upsilonbattle`.
3. **Engine Fragility**: Changes to the world-building logic should not require recompiling or risking the stability of the core battle engine.

### Where This Pattern Exists Today
- `/workspace/upsilonbattle/battlearena/property/` - Core property definitions.
- `/workspace/upsilonbattle/battlearena/entity/skill/skillweight/` - Balancing math.
- `/workspace/upsilonbattle/battlearena/entity/skill/skillgenerator/` - Procedural generation (to be moved).

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (as meta-game grows) |
| Impact if triggered | High architectural debt |
| Detectability | Medium - seen as increasing complexity in `upsilonapi` |
| Current mitigant | Shared Go workspace, but logically tightly coupled |

---

## Recommended Fix

**Short term:**
- Extract `battlearena/property` into a new package/module (e.g., `github.com/ecumeurs/upsilon-mechanics`).
- Relocate `skillweight` to this same package.

**Medium term:**
- Update `upsilonbattle` and `upsilonapi` to import the new shared library.
- Refactor the `skillgenerator` to live in a dedicated "Generation" domain.

**Long term:**
- Implement the "Token Combination" and "Regional Context" features within the new standalone Generation service.

---

## References

- [ISS-065 (Skill Weight)](ISS-065_20260422_skill_weight_grading_system.md)
- [ISS-073 (Roguelike Skill System)](ISS-073_20260423_roguelike_skill_system.md)
- `upsilonbattle/battlearena/property/property.go`
- `upsilonbattle/battlearena/entity/skill/skillweight/skillweight.go`
