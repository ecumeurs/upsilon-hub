---
id: mec_weapon_as_skill_system
human_name: Weapon-as-Skill System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority=5
tags=[equipment, weapons, combat]
parents:
  - [[entity_equipment_system]]
dependents: []
---

# Weapon-as-Skill System Mechanic

## INTENT
To implement the weapon-as-skill system where equipped weapons transform basic attacks into skill-based attacks with properties like range, damage, critical chance, enabling weapon variety through the skill system.

## THE RULE / LOGIC
**Weapon Transformation:**
- **Basic Attack Becomes Skill:** When weapon equipped, basic attack uses skill computation
- **Weapon Properties Become Skill Properties:** WeaponRange becomes skill Range, WeaponBaseDamage becomes skill Damage
- **Attack Computation:** Use skill-based damage formula instead of simple formula

**Weapon Properties to Skill Mapping:**
```go
// Weapon Property -> Skill Property Mapping
WeaponRange -> Range (Skill Property)
WeaponBaseDamage -> Damage (Skill Property)
WeaponType -> Behavior (Melee/Ranged differentiation)
CritChance -> CriticalChance (Skill Property)
CritMultiplier -> CriticalMultiplier (Skill Property)
```

**Attack Flow with Weapon:**
```go
// When player uses basic attack
if character.EquippedWeapon != nil {
    weapon := character.EquippedWeapon
    
    // Create temporary skill from weapon properties
    weaponSkill := skill.Skill{
        Targeting: map[property.Property]{
            property.Range: weapon.WeaponRange,
        },
        Effect: effect.Effect{
            Properties: []property.Property{
                defaultproperty.MakeIntProperty(property.Damage, weapon.WeaponBaseDamage, ...),
            },
        },
    }
    
    // Use skill-based attack computation
    damage := ComputeSkillAttack(weaponSkill, target)
} else {
    // Fallback to basic attack (unarmed)
    damage := ComputeBasicAttack(character, target)
}
```

**Weapon Skill Examples:**
- **Sword (Melee):** Range 1, Damage +5, Crit +2%
- **Bow (Ranged):** Range 5, Damage +3, Crit +5%
- **Dagger (Melee):** Range 1, Damage +2, Crit +10%, Backstab +50%
- **Staff (Magic):** Range 3, Damage +1, MP Cost -5

**Weapon-Skill Synergy:**
- **Weapons define base attack behavior**
- **Skills can modify weapon effects** (future V2.2+)
- **Equipment progression** = Weapon skill progression

**Balance Considerations:**
- **Weapon Unarmed:** Default basic attack without skill properties
- **Weapon Equipped:** Always use weapon-based attack
- **Skill + Weapon:** Future synergy where skills enhance weapons

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_weapon_as_skill_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonbattle/battlearena/property/def/item.go`
