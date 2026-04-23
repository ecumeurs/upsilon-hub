# Issue: Skill Weight & Grading System Implementation

**ID:** `20260422_skill_weight_grading_system`
**Ref:** `ISS-065`
**Date:** 2026-04-22
**Severity:** High
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/entity/skill`
**Affects:** `battleui`, `upsilonapi`, `upsiloncli`

---

## Summary

Implement a mathematical Skill Weight (SW) system for balanced skill design, automatic skill grading (I-V), and credit cost calculation. Every skill must have Net SW = 0 (benefits = costs), with grades determined by Total Positive SW.

---

## Technical Description

### Background

Current skill system has properties but no mathematical framework for balance. Skills are randomly generated without weight considerations, making skill selection and shop pricing impossible.

### The Problem Scenario

1. **Design Phase**: Creator wants to design a "Fireball" skill
2. **Balance Challenge**: How much should it cost? What grade is it?
3. **No Framework**: Current system provides no guidance for these decisions
4. **Player Impact**: Random skill generation creates imbalance

### Skill Weight Framework

**Benefits Table (Adds SW):**
- Damage Multiplier: +10 SW per 10% (100% damage = +100 SW)
- Critical Chance: +2 SW per 1% (+25% crit = +50 SW)
- Range Extension: +10 SW per cell > 1 (Range 3 = +20 SW)
- Zone/AoE: +50 SW per extra cell
- Target Anywhere: +40 SW
- Stun Chance: +2 SW per 1%
- Poison Power: +15 SW per 1 dmg/turn
- Heal: +15 SW per 10% base
- Shield: +10 SW per point

**Payments Table (Reduces SW):**
- Delay (+100): -100 SW (baseline)
- Extra Delay: -10 SW per +10 delay
- Channeling: -15 SW per 10 delay
- MP Leech: -15 SW per 1 MP
- SP Leech: -10 SW per 1 SP
- HP Leech: -20 SW per 1 HP
- Cooldown: -25 SW per turn

**Grading System:**
- Grade I: 0-150 Positive SW
- Grade II: 151-300 Positive SW
- Grade III: 301-500 Positive SW
- Grade IV: 501-750 Positive SW
- Grade V: 750+ Positive SW

**Credit Cost:** Total Positive SW × 2

### Where This Pattern Exists Today

- `upsilonbattle/battlearena/entity/skill/skill.go` - Current skill structure
- `upsilonbattle/battlearena/entity/skill/skillgenerator/skillgenerator.go` - Random generation
- `upsilonbattle/battlearena/property/propertyenum.go` - Property definitions

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Implement SkillWeight calculator function and grading algorithm. Update skill generator to produce balanced skills.

**Medium term:** Create skill template system with pre-balanced skill definitions. Integrate with shop pricing.

**Long term:** Build skill balancing tools for designers and automatic imbalance detection.

---

## References

- `third_party_reply.md` - Skill Weight mathematical framework
- `upsilonbattle/battlearena/property/propertyenum.go` - Property definitions
- `docs/domain_skill_system.atom.md` - Skill system domain
