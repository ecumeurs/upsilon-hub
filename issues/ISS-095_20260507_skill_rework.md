# Issue: Skill Rework & Generation Standardization

**ID:** `20260507_skill_rework`
**Ref:** `ISS-095`
**Date:** 2026-05-07
**Severity:** Medium
**Status:** Open
**Component:** `upsilontypes/entity/skill/skillgenerator`
**Affects:** `upsilonapi`, `battleui`, `upsilonbattle`

---

## Summary

The current skill generation and execution framework suffers from several design inconsistencies and technical debt:
1. **Zero-Value Clutter**: Generated skills include redundant properties (e.g., `Damage: 0`, `Delay: 0`) in the API output due to explicit initialization and lack of pruning.
2. **Incomplete Mechanics**: Status effects like Stun and Poison rely on separate Power and Chance properties, but producers often set only one, resulting in mechanically dead skills.
3. **Ambiguous Defaults**: `Damage` defaulting to 100% makes non-damaging skills visually noisy in the engine, but changing it to 0 breaks existing balance bands.
4. **Targeting Interpretation**: `Range: 0` is interpreted as "Self/Current Cell" but often appears unintentionally in DOT/Stun skills that should have a baseline range.

---

## Technical Description

### 1. The Stun/Poison Pairing Risk
In `effectapplicator.ApplyDirectEffect`, status effects only trigger if `Power - Resistance > 0` AND a `RandomInt(0, 100) < Chance` roll succeeds.
- **Problem**: Many skills are generated with `StunPower: 5` but `StunChance: 0` (default), or vice versa. This renders the effect impossible to apply.
- **Requirement**: Any producer adding a status effect must ensure both properties are non-zero.

### 2. Zero-Value Property Clutter
The `skillgenerator.blueprint` explicitly initializes critical properties to 0 to avoid inheriting engine defaults.
- **Problem**: When serialized by the API, these 0-values appear in the JSON payload, creating noise for the UI and complicating rehydration logic.
- **Proposed Solution**: A `Prune()` layer in the generator should remove properties that match their engine defaults before finalization.

### 3. Damage Default & Balancing
- **Status Quo**: `Damage` defaults to 100.
- **Risk**: Changing the default to 0 (to allow "clean" utility skills) shifts the Skill Weight (SW) of all existing skills by 100 points, breaking the current Grade Bands (e.g., Grade I floor of 60 PSW).
- **Decision**: Revert to 100 for now, but investigate a structural rework of the weight calculator to support a 0-based baseline.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium — Visual clutter and "dead" mechanics like stuns that never proc. |
| Detectability | High — Observed in API payloads and battle logs. |
| Current mitigant | Explicit property setting in some producers. |

---

## Recommended Fix

**Short term:**
- Implement a `Prune()` method in the blueprint to strip default-value properties.
- Update `producer_stun` and `producer_dot` to ensure functional pairing (Power + Chance).

**Medium term:**
- Implement the "Reaction and Counter" triggering logic (originally identified in this issue).
- Re-balance the Skill Weight bands to support a 0-default for Damage and other primary levers.

**Long term:**
- Unified `Condition` system for skill activation and effect application.

---

## Extra Data

**Targeting Findings**:
- `Range: 0` is valid but should be restricted to "Self-Targeting" skills.
- `TargetingMechanics` defaults to "Anywhere", but "Line of Sight" (LoS) is preferred for standard balance.

---

## References

- [upsilontypes/entity/skill/skillgenerator/blueprint.go](file:///workspace/upsilontypes/entity/skill/skillgenerator/blueprint.go)
- [upsilonbattle/battlearena/property/effect/effectapplicator/effectapplicator.go](file:///workspace/upsilonbattle/battlearena/property/effect/effectapplicator/effectapplicator.go)
- [upsilontypes/property/def/skill.go](file:///workspace/upsilontypes/property/def/skill.go)
