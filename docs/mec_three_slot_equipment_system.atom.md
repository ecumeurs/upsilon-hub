---
id: mec_three_slot_equipment_system
human_name: Three-Slot Equipment System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [equipment, slots, inventory]
parents:
  - [[entity_equipment_system]]
dependents: []
---

# Three-Slot Equipment System Mechanic

## INTENT
To implement the simplified 3-slot equipment system with exactly 1 armor slot, 1 utility slot, and 1 weapon slot per character, providing focused equipment progression without inventory management complexity.

## THE RULE / LOGIC
**Equipment Slots:**
- **Armor Slot:** One piece of armor equipment (Head, Body, Legs, etc.)
- **Utility Slot:** One utility item (Ring, Amulet, Belt, etc.)
- **Weapon Slot:** One weapon (Melee or Ranged, One or Two-Handed)

**Slot Constraints:**
- **Exactly One Per Slot:** No character can have multiple items in same slot type
- **Slot Independence:** Each slot can have exactly one item equipped
- **Mutual Exclusivity:** Equipping new item in slot replaces old item
- **Slot Types:** Armor, Utility, Weapon are distinct categories

**Equipment Equipping Logic:**
```go
// Equip new item
func (c *Character) EquipItem(item Equipment, slot EquipmentSlot) {
    // Validate slot type matches item type
    if item.Type != slot.AllowedType {
        return error("Item cannot be equipped in this slot")
    }
    
    // Unequip current item in slot if any
    if c.Slots[slot] != nil {
        c.Inventory.AddItem(c.Slots[slot])
    }
    
    // Equip new item
    c.Slots[slot] = item
    
    // Apply stat bonuses immediately
    c.RecalculateStats()
}
```

**Slot Type Mapping:**
- **Armor Slot:** Accepts ItemType = Wearable (Head, Body, Hands, Legs, Feet)
- **Utility Slot:** Accepts ItemType = Wearable (Neck, Ring, Belt) or special utility items
- **Weapon Slot:** Accepts ItemType = Wearable (OneHandedMelee, TwoHandedMelee, OneHandedRanged, TwoHandedRanged)

**Two-Handed Handling:**
- **Two-Handed Weapons:** Occupy weapon slot only, but restrict utility slot
- **Two-Handed Weapons:** May prevent utility item usage while equipped
- **Balance Consideration:** Two-handed weapons typically provide higher damage to compensate

**Inventory Management:**
- **Unequipped Items:** Stored in character inventory
- **No Carrying Capacity:** Unlimited inventory for unequipped items
- **Quick Swapping:** Can swap items between slots and inventory instantly

**Slot Benefits:**
- **Simplified UI:** 3 slots = simple inventory interface
- **Clear Progression:** One item per slot = clear upgrade path
- **Reduced Complexity:** No inventory management or weight systems
- **Tactical Clarity:** Easy to understand what each character has equipped

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_three_slot_equipment_system]]`
- **Related Files:** `upsilonapi/api/input.go` (Character structure), equipment database schema
