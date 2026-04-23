# Issue: Backstabbing Mechanics & Armor Penetration

**ID:** `20260422_backstabbing_mechanics`
**Ref:** `ISS-070`
**Date:** 2026-04-22
**Severity:** Medium
**Status:** Open
**Component:** `upsilonbattle/battlearena/ruler/rules`
**Affects**: `upsilonbattle/battlearena/entity`, `battleui`

---

## Summary

Implement backstabbing combat mechanic with 150% damage multiplier and 50% armor penetration (applies to shield, not armor). Backstabbing works with weapon attacks only and requires attacking from behind the target.

---

## Technical Description

### Background

Facing system exists (EntityOrientation: Up, Right, Down, Left) but no damage modifiers based on orientation. All attacks deal same damage regardless of positioning.

### The Problem Scenario

1. **Positional tactics matter**: Players maneuver for tactical advantage
2. **No reward for flanking**: Sneaking behind enemies provides no damage bonus
3. **Combat lacks depth**: Frontal and rear attacks are identical
4. **Sneak archetype underpowered**: No mechanical benefit to positioning

### Backstabbing Mechanics

**Detection Algorithm:**
```go
func (attacker *Entity) IsBackstabbing(target *Entity) bool {
    // Target is backstabbed if attacker faces opposite direction
    oppositeDirection := getOppositeDirection(target.Orientation)
    
    // Check if attacker is behind target (within 45° of opposite direction)
    angleToTarget := attacker.Position.AngleTo(target.Position)
    targetAngle := getDirectionAngle(target.Orientation)
    backAngle := (targetAngle + 180) % 360
    
    // Backstab if attacker is within 45° of back angle
    angleDifference := abs(angleToTarget - backAngle)
    return angleDifference <= 45 || angleDifference >= 315
}

func getOppositeDirection(orientation EntityOrientation) EntityOrientation {
    switch orientation {
    case Up: return Down
    case Down: return Up
    case Left: return Right
    case Right: return Left
    }
}
```

**Damage Calculation:**
```go
func (gs *GameState) ApplyBackstabDamage(attacker, target entity.Entity, baseDamage int) int {
    if !attacker.IsBackstabbing(target) {
        return baseDamage  // No backstab bonus
    }
    
    // 150% damage multiplier
    backstabDamage := baseDamage * 1.5
    
    // Get target defenses
    armorRating := target.GetPropertyI(property.ArmorRating).I()
    shield := target.GetPropertyC(property.Shield).GetValue()
    
    // Ignore 50% armor rating
    effectiveArmor := armorRating * 0.5
    
    // Shield still applies fully
    finalDamage := max(1, backstabDamage - effectiveArmor)
    
    // Apply shield first
    if shield > 0 {
        if shield >= finalDamage {
            target.UpdatePropertyValue(property.Shield, shield - finalDamage)
            return 0  // All damage absorbed by shield
        } else {
            finalDamage -= shield
            target.UpdatePropertyValue(property.Shield, 0)
        }
    }
    
    return finalDamage
}
```

**Scope Limitations:**
- **Weapon Attacks Only**: Skills do not benefit from backstabbing (future expansion)
- **All Weapons**: Bows, pistols, melee weapons all support backstabbing
- **No Skill Synergy**: Future skills may have "Backstab Enabled" property (300% skillpower, etc.)

### Where This Pattern Exists Today

- `upsilonbattle/battlearena/entity/entity.go` - EntityOrientation exists
- `upsilonbattle/battlearena/entity/entity.go` - FaceToward() method exists
- `upsilonbattle/battlearena/ruler/rules/attack.go` - Current attack logic
- `upsilonbattle/battlearena/property/propertyenum.go` - ArmorRating property exists

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | Facing system exists, need damage modifier logic |

---

## Recommended Fix

**Short term:** Implement back detection algorithm using existing orientation system. Add 150% damage multiplier. Create 50% armor penetration logic.

**Medium term:** Integrate backstabbing with weapon attack system. Add visual feedback for backstabs. Update AI to avoid backstabs.

**Long term:** Expand to skills with backstab-enabled properties. Create backstab bonuses from equipment. Implement advanced positioning tactics (sidestab, flanking).

---

## References

- `V2_ARCHITECTURAL_DECISIONS.md` - Backstabbing decisions
- `upsilonbattle/battlearena/entity/entity.go` - EntityOrientation and FaceToward
- `upsilonbattle/battlearena/ruler/rules/attack.go` - Current attack computation
- `upsilonbattle/battlearena/property/def/item.go` - ArmorRating property
