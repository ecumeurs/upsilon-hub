# Issue: Starting Stats & Character Progression Redesign

**ID:** `20260422_starting_stats_progression`
**Ref:** `ISS-071`
**Date:** 2026-04-22
**Severity:** High
**Status:** Open
**Component:** `upsilonapi/api`, `battleui`
**Affects**: Character creation, progression system, skill system balance

---

## Summary

Redesign character starting stats and progression system from V1 (4 random points, base stats 3/1/1/1) to V2 (100 CP point-buy system, x10 baseline stats: HP 30-50, Attack 10, Defense 5, Movement 3). This enables meaningful skill percentages and character variety.

---

## Technical Description

### Background

V1 character system uses extremely low base integers (Attack 1, HP 3) that make percentage-based skill modifiers mathematically invisible. 120% damage on Attack 1 = 1.2 damage → rounds to 1 (useless skill). Current progression uses simple +1 point per win with limited character variety.

### The Problem Scenario

1. **Skill System Imbalance**: Percentage modifiers don't work with low base stats
2. **Limited Character Variety**: 4 random points don't create meaningful differences
3. **Progression Too Linear**: +1 point per win doesn't scale well
4. **Movement Overpowered**: Hard-lock "once every 5 wins" feels arbitrary

### V2 Stat Redesign: x10 Baseline

| Attribute | V1 Base | **V2 Base** | Rationale |
| :--- | :--- | :--- | :--- |
| **HP** | 3 | **30 - 50** | Characters survive 3-4 hits instead of 1-2 |
| **Attack** | 1 | **10** | Clean baseline for percentages (10% = 1 damage) |
| **Defense** | 1 | **5** | Lower than Attack ensures "1 damage rule" doesn't trigger constantly |
| **Movement** | 1 | **3** | Standard TRPG baseline, positioning matters without crossing whole board |

### V2 Progression: 100 CP Point-Buy System

**Instead of 4 random points, players get 100 Character Points (CP) to spend:**

**Standard Attributes:**
- **HP (+1):** Cost 1 CP (Linear, cheap)
- **Attack (+1):** Cost 5 CP (Direct damage scaling is powerful)
- **Defense (+1):** Cost 5 CP (Direct damage mitigation is equally powerful)

**Exotic Attributes:**
- **Critical Chance (+1%):** Cost 10 CP (High value, caps easily)
- **Critical Multiplier (+5%):** Cost 5 CP
- **Jump Height (+1):** Cost 15 CP (Drastically alters terrain navigation)

**Movement Premium:**
- **Movement (+1 cell):** Cost 30 CP (Most potent stat, expensive = natural restriction)
- **Note:** Eliminates need for "once every 5 wins" movement hard-lock

### Progression Rules Update

**Character Creation:**
- **Starting Pool:** 100 CP to spend on top of V2 base stats
- **Random Distribution:** Players allocate points (no more random 4-point system)
- **Character Variety:** Meaningful stat differences possible

**Win Rewards:**
- **CP per Win:** +10 CP (instead of +1 point)
- **Total Cap:** 100 + (total_wins × 10)
- **Character Power:** Scales dramatically with V2 stat system

**Skill Integration:**
- **Skill Viability:** Percentage modifiers now meaningful (120% of 10 attack = 12 damage)
- **Stat Scaling:** Higher base stats make skill effects noticeable
- **Balance Framework:** Skill Weight system works with x10 baseline

### Where This Pattern Exists Today

- `docs/rule_character_create_character.atom.md` - Current V1 character creation rules
- `docs/rule_progression.atom.md` - Current V1 progression rules
- `upsilonapi/api/input.go` - Player/Character data structures

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | None (complete system redesign needed) |

---

## Recommended Fix

**Short term:** Update character creation schema for x10 baseline stats. Implement 100 CP point-buy system. Update progression rewards to +10 CP per win.

**Medium term:** Update all existing characters to V2 stat system (migration script). Rebalance skill effects for new baseline. Adjust AI stat allocation for V2.

**Long term:** Implement exotic attribute progression (Crit Chance, Multiplier, Jump Height). Create character stat respec system. Balance point costs based on gameplay data.

---

## References

- `third_party_reply.md` - Complete stat redesign analysis and rationale
- `docs/rule_character_create_character.atom.md` - Current V1 rules
- `docs/rule_progression.atom.md` - Current V1 progression
- `V2_ARCHITECTURAL_DECISIONS.md` - Integrated V2 decisions
