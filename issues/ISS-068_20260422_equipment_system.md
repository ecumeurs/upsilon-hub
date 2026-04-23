# Issue: Equipment System & Weapon-as-Skill

**ID:** `20260422_equipment_system`
**Ref:** `ISS-068`
**Date:** 2026-04-22
**Severity:** Medium
**Status:** Open
**Component:** `upsilonbattle/battlearena/property`, `upsilonapi/api`
**Affects**: `battleui`, character progression

---

## Summary

Implement 3-slot equipment system (1 armor, 1 utility, 1 weapon) with weapon-as-skill mechanics where equipped weapons transform basic attacks into skill-based attacks with properties like range, damage, crit chance.

---

## Technical Description

### Background

Item properties exist (ArmorRating, WeaponRange, WeaponBaseDamage) but no equipment slot system or inventory management. Basic attacks use simple formula without skill properties.

### The Problem Scenario

1. **Player wants to equip items**: No inventory or slot system exists
2. **Weapon variety limited**: All weapons use same basic attack formula
3. **No equipment progression**: Players can't acquire or upgrade equipment
4. **Simple combat**: Basic attacks lack depth and variety

### Equipment System Architecture

**3-Slot System:**
- **Armor Slot**: Head, Body, Hands, Legs, Feet (one item)
- **Utility Slot**: Neck, Ring, Belt (one item)
- **Weapon Slot**: Main Hand weapon (one item)

**Equipment Properties (Existing):**
- **Armor:** ArmorRating (defense boost)
- **Weapons:** WeaponRange, WeaponBaseDamage, WeaponType
- **Utility:** Various buffs and special effects

**Weapon-as-Skill System:**
```go
// When weapon equipped, basic attack becomes skill
type EquippedWeapon struct {
    ItemID       uuid.UUID
    Name         string
    WeaponType   WeaponTypes     // One/Two-Handed Melee/Ranged
    WeaponRange  int             // Attack range
    WeaponDamage int             // Base damage
    CritChance   int             // Critical chance %
    CritMultiplier int           // Critical multiplier %
    SkillEffect  *effect.Effect  // Additional skill properties
}

// Basic attack transforms to skill-based attack
func (gs *GameState) Attack(msg, req) {
    weapon := ent.EquippedWeapon
    
    if weapon != nil {
        // Use skill-based computation instead of basic formula
        gs.SkillAttack(weapon.SkillEffect, req.Target)
    } else {
        // Fallback to basic attack
        gs.BasicAttack(req.Target)
    }
}
```

**Equipment Effects:**
- **Direct Stat Bonuses:** ArmorRating adds to Defense
- **Skill Augmentation:** Weapons add properties to basic attack
- **Passive Effects:** Utility items grant ongoing benefits

### Where This Pattern Exists Today

- `upsilonbattle/battlearena/property/def/item.go` - Item properties exist
- `upsilonbattle/battlearena/property/propertyenum.go` - Property definitions
- `upsilonbattle/battlearena/ruler/rules/attack.go` - Basic attack logic

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | Item properties exist, need slot system |

---

## Recommended Fix

**Short term:** Design 3-slot equipment schema. Create equipment inventory tables. Build equip/unequip API endpoints.

**Medium term:** Implement weapon-as-skill system. Create equipment stat bonus logic. Build equipment management UI.

**Long term:** Add equipment durability and repair. Implement equipment upgrading and enchantment. Create equipment rarity tiers.

---

## References

- `V2_ARCHITECTURAL_DECISIONS.md` - Equipment system decisions
- `upsilonbattle/battlearena/property/def/item.go` - Item property definitions
- `upsilonbattle/battlearena/ruler/rules/attack.go` - Current basic attack
