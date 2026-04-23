---
id: mec_equipment_stat_bonuses
human_name: Equipment Stat Bonuses Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [equipment, stats, bonuses]
parents:
  - [[entity_equipment_system]]
dependents: []
---

# Equipment Stat Bonuses Mechanic

## INTENT
To implement equipment stat bonus system where equipped items provide direct attribute modifications to characters, with armor adding defense, weapons adding attack power, and utility items providing special effects.

## THE RULE / LOGIC
**Stat Bonus Types:**
- **Direct Attribute Modification:** Armor adds to Defense, Weapons add to Attack
- **Percentage Modifiers:** Items can provide percentage bonuses to base stats
- **New Property Introduction:** Utility items can grant entirely new properties (CritChance, JumpHeight)
- **Temporary Stat Changes:** Some equipment provides conditional bonuses

**Armor Stat Bonuses:**
- **ArmorRating:** Direct addition to Defense stat
- **Special Effects:** Some armor provides elemental resistance, movement penalties
- **Slot-Specific:** Different armor slots provide different bonus types

**Weapon Stat Bonuses:**
- **WeaponBaseDamage:** Direct addition to Attack stat (used in weapon-as-skill system)
- **WeaponRange:** Sets attack range (overrides default range)
- **Special Properties:** Weapons can add CritChance, Backstab bonuses, special damage types

**Utility Stat Bonuses:**
- **MP/SP Pools:** Increase maximum mana or stamina
- **Movement Enhancements:** Improve jump height or movement cost
- **Special Properties:** Grant entirely new capabilities (flight, stealth, etc.)

**Stat Bonus Calculation:**
```go
// When character properties are requested
func (c Character) GetProperty(prop interface{}) property.Property {
    baseProp := c.BaseProperties[prop]  // Character's base stats
    
    // Add equipment bonuses
    equipmentBonus := 0
    for _, item := range c.EquippedItems {
        if item.StatBonuses[prop] != 0 {
            equipmentBonus += item.StatBonuses[prop]
        }
    }
    
    // Apply bonus to base property
    baseProp.SetValue(baseProp.GetValue() + equipmentBonus)
    return baseProp
}
```

**Equipment Bonus Examples:**
- **Iron Armor:** Defense +3, ArmorRating +5
- **Steel Sword:** Attack +4, Range +1
- **Magic Ring:** MP +20, CritChance +2%
- **Boots of Speed:** Movement +1

**Stacking Rules:**
- **Same Property Bonuses:** Stack additively from different equipment
- **Different Equipment Types:** Armor + Weapon + Utility bonuses all apply
- **Maximum Caps:** Some bonuses may have maximum values to prevent abuse
- **Diminishing Returns:** Some bonuses have reduced effectiveness at high values

**Temporary vs Permanent:**
- **Permanent Bonuses:** Equipment stats apply as long as item is equipped
- **Conditional Bonuses:** Some items provide bonuses only under specific conditions
- **Charges:** Some equipment has limited charges for special effects

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_equipment_stat_bonuses]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entity.go`, `upsilonbattle/battlearena/property/def/item.go`
