# Issue: Skill Inspection Command

**ID:** `20260423_skill_inspection`
**Ref:** `ISS-077`
**Date:** 2026-04-23
**Severity:** Medium
**Status:** Open
**Component:** `upsiloncli`, `battleui`, `upsilonapi`
**Affects:** CLI commands, UI modals, skill information display

---

## Summary

Implement skill inspection functionality allowing players to view detailed skill properties, effects, and statistics via CLI command and UI modal. This provides transparency into skill mechanics while respecting privacy boundaries for unowned skills.

---

## Technical Description

### Background

Players need to understand their skills to make strategic decisions during equipment and battle. Current system lacks a way to inspect skill details outside of the selection modal. CLI also needs a command to view skill information for debugging and analysis.

### The Problem Scenario

1. **What does this skill do?**: No way to see detailed skill properties
2. **How much damage?**: No damage formula breakdown shown
3. **Status effects unclear**: How does poison work? What's the duration?
4. **CLI debugging**: No way to inspect skill JSON structure
5. **Privacy concern**: Should not expose details for skills player doesn't own

### Skill Inspection Architecture

**Privacy Rules:**
- Players can only inspect skills they own
- Enemy character skill details are hidden
- Inspect command requires authentication check
- Skill templates visible to all for browsing
- Owned skill details include acquisition date, reforge count

**Inspection Data:**
```go
type SkillDetail struct {
    ID          uuid.UUID
    Name        string
    Grade       string         // I-V
    TotalSW     int            // Skill weight
    Cost        int            // Credit cost
    Properties  []PropertyInfo  // All properties with values
    Effects     []EffectInfo     // All effects with descriptions
    Stats       SkillStats      // Usage, damage dealt, credits earned
}

type Property struct {
    Name        string
    Value       interface{}
    Description string
}

type EffectInfo struct {
    Type      string       // Damage, Heal, Stun, etc.
    Value      interface{}   // Magnitude
    Duration   int         // Turns (if applicable)
    Target    string       // Self, Enemy, Anywhere
}
```

### CLI Integration

**Commands:**
```bash
# Inspect owned skill
skill inspect <skillId>

# List all owned skills with details
skill list --details

# Get skill usage statistics
skill stats <skillId>
```

**Output Format:**
```
=== FIREBALL ===
Grade: II (Total SW: 250)
Cost: 500 credits

Properties:
  - Damage: +15 (150 SW)
  - Range: 3 cells (+20 SW)
  - Cooldown: 2 turns (-50 SW)

Effects:
  - Fire Damage: 15 HP
  - Area: 3x3 cone (+100 SW)
  - Burn: 5 HP/turn for 2 turns (+50 SW)

Stats:
  - Total Uses: 142
  - Damage Dealt: 2,150 HP
  - Credits Earned: 1,420
```

### UI Integration

**Skill Detail Modal:**
- Skill name with grade badge (I-V colors)
- Properties section with each property and value
- Effects section with all effects and tooltips
- Statistics section with usage metrics
- "Equip This Skill" button (if slots available)
- Close button

**Access Control:**
- Modal only opens for owned skills
- Attempting to inspect unowned skill returns error
- Skills in shop shown without stats (public catalog)

### API Endpoints

**Skill Inspection:**
- `GET /api/v1/skill/{id}/inspect` - Full details for owned skill
- `GET /api/v1/skills/templates` - Public catalog without stats

**Validation:**
- Skill ID must belong to player
- Cannot inspect skills from other players' characters
- Skill templates don't require ownership

### Where This Pattern Exists Today

- `docs/rule_skill_grading_system.atom.md` - Grade I-V definitions
- `docs/mech_skill_selection_progression.atom.md` - Skill acquisition
- `docs/mech_skill_reforging_mechanic.atom.md` - Skill modifications
- `ISS-073` (Roguelike Skill System) - Skill inventory to inspect
- `docs/rule_character_skill_slots.atom.md` - Slot limits (related to equipment)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Low |
| Detectability | High |
| Current mitigant | Skills exist, just need inspection endpoint |

---

## Recommended Fix

**Short term:** Implement CLI `skill inspect` command. Add `GET /api/v1/skill/{id}/inspect` endpoint. Create basic skill detail modal.

**Medium term:** Add usage statistics tracking per skill. Create `skill list --details` command with formatted output. Add filtering and sorting to skill list.

**Long term:** Advanced statistics (average damage per use, hit rate, win rate with skill). Comparison tool between skills. Export skill data to JSON.

---

## References

- `ISS-065` (Skill Weight & Grading) - Skills have properties and weights
- `ISS-067` (Credit Economy) - Skills cost credits, track usage
- `docs/mechanic_mech_skill_weight_calculator.atom.md` - Property value calculations
