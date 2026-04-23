# Issue: AI Archetype Enhancement & Progression

**ID:** `20260422_ai_archetype_enhancement`
**Ref:** `ISS-069`
**Date:** 2026-04-22
**Severity:** Medium
**Status:** Open
**Component:** `upsilonbattle/battlearena/controller/controllers`
**Affects**: `upsilonbattle/battlearena/ruler`, matchmaking system

---

## Summary

Enhance AI system with four distinct archetypes (Fighter, Ranger, Support, Sneak) that follow player progression rules and use archetype-specific stat allocation and skill selection. Implement team composition limits (max 1 support, max 1 sneak per AI team).

---

## Technical Description

### Background

Current AI system has single AggressiveController archetype with basic target selection and pathfinding. No skill usage, difficulty scaling, or team composition logic.

### The Problem Scenario

1. **All AI acts identically**: Single aggressive behavior, no variety
2. **No progression**: AI doesn't scale with player levels
3. **Simple tactics**: No skill usage, limited strategic depth
4. **Unbalanced teams**: No composition constraints create unbalanced matches

### AI Archetype System

**Four Archetypes:**

**Fighter Controller:**
- **Skills:** High damage melee skills, defensive skills, charge abilities
- **Stat Priority:** Attack > Defense > Movement > HP
- **Tactics:** Direct approach, aggressive positioning, focus on damage

**Ranger Controller:**
- **Skills:** Ranged damage skills, trap placement, movement/positioning skills
- **Stat Priority:** Attack > Movement > Accuracy > Defense
- **Tactics:** Kiting behavior, maintain range, use terrain advantages

**Support Controller:**
- **Skills:** Healing, shielding, buff/debuff application, ally positioning
- **Stat Priority:** MP/SP > Defense > HP > Attack
- **Tactics:** Stay near allies, protect weak team members, prioritize healing

**Sneak Controller:**
- **Skills:** Backstabbing bonuses, movement skills, poison/stun application, stealth
- **Stat Priority:** Movement > Attack > Dodge > CritChance
- **Tactics:** Flanking, target weak enemies, avoid direct confrontation

**Progression System:**
- **Same Rules:** +1 point per win, max 10 + total_wins
- **Movement Restriction:** +1 movement once every 5 levels
- **Skill Selection:** Choose from archetype-appropriate skill pools
- **Level Matching:** AI level matches average player level

**Team Composition Rules:**
```go
func (team *Team) ValidateComposition() bool {
    supportCount := 0
    sneakCount := 0
    
    for _, char := range team.Characters {
        if char.ControllerType == Support {
            supportCount++
        } else if char.ControllerType == Sneak {
            sneakCount++
        }
    }
    
    return supportCount <= 1 && sneakCount <= 1
}
```

### Where This Pattern Exists Today

- `upsilonbattle/battlearena/controller/controllers/aggressive.go` - Current AI controller
- `upsilonbattle/battlearena/controller/controller.go` - Base controller structure
- `docs/rule_progression.atom.md` - Progression rules

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | AggressiveController provides good base to extend |

---

## Recommended Fix

**Short term:** Create four archetype controllers extending base controller. Implement archetype-specific skill pools. Add team composition validation.

**Medium term:** Implement archetype stat allocation algorithms. Connect AI to player progression rules. Add AI skill selection logic.

**Long term:** Implement advanced AI behaviors (tactical positioning, skill combinations). Create difficulty scaling beyond level matching. Add AI personality and communication.

---

## References

- `V2_ARCHITECTURAL_DECISIONS.md` - AI progression decisions
- `upsilonbattle/battlearena/controller/controllers/aggressive.go` - Current AI implementation
- `docs/rule_progression.atom.md` - Shared progression rules
