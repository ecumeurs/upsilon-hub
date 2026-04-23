---
id: temp
status: DRAFT
version: 1.0
parents: []
dependents: []
---

# New Atom

## INTENT
To implement equipment shop inventory system where players can browse, filter, and purchase equipment items based on character level, with prices determined by equipment tier and rarity.

## THE RULE / LOGIC
**Equipment Shop Inventory System:**

**Core Principle:**
Shop provides filtered access to equipment based on character level, with pricing determined by tier multipliers.

**Shop Categories:**

**Armor Equipment:**
- **Head:** Helmets, crowns, masks (ArmorRating +0-5)
- **Body:** Chestplates, tunics, robes (ArmorRating +5-10)
- **Hands:** Gauntlets, gloves (ArmorRating +1-3)
- **Legs:** Greaves, pants (ArmorRating +2-6)
- **Feet:** Boots, shoes (ArmorRating +1-4)

**Utility Equipment:**
- **Neck:** Amulets, necklaces (Special effects, stat bonuses)
- **Rings:** Rings with various enchantments (Stat bonuses, MP/SP)
- **Belt:** Belts with capacity or effects (Inventory size, HP)
- **Special:** Consumables, temporary buffs (future V2.2+)

**Weapons:**
- **One-Handed Melee:** Swords, axes, maces (Damage +0-8, Range 1-2)
- **Two-Handed Melee:** Greatswords, hammers (Damage +6-12, Range 1-2)
- **One-Handed Ranged:** Bows, pistols, crossbows (Damage +0-6, Range 4-7)
- **Two-Handed Ranged:** Longbows, heavy crossbows (Damage +4-10, Range 6-10)

**Equipment Properties:**

**Common Items (Tier 1):**
- **Leather Armor:** +1-3 ArmorRating, 50 credits
- **Iron Weapons:** +2-4 Damage, 100-150 credits
- **Simple Rings:** +5% Crit, +10% CritMultiplier, 200 credits

**Uncommon Items (Tier 2):**
- **Chain Armor:** +3-6 ArmorRating, 200 credits
- **Steel Weapons:** +5-7 Damage, 300-450 credits
- **Enchanted Rings:** +10% Crit, +20% CritMultiplier, +2 MP, 400 credits

**Rare Items (Tier 3):**
- **Plate Armor:** +6-10 ArmorRating, 500-800 credits
- **Mithral Weapons:** +8-12 Damage, special effects, 700-1000 credits
- **Legendary Jewelry:** +20% Crit, +50% CritMultiplier, +10 MP, 1200 credits

**Shop Filtering:**

**Character Level Filter:**
```go
func FilterShopItems(characterLevel int, allItems []Equipment) []Equipment {
    availableItems := []Equipment{}
    
    for _, item := range allItems {
        // Check level requirement
        if item.MinimumLevel > characterLevel {
            continue  // Item too high level
        }
        
        // Check skill access requirement
        if item.RequiresSkill && !character.HasSkill(item.RequiresSkill) {
            continue  // Prerequisite skill not owned
        }
        
        availableItems = append(availableItems, item)
    }
    
    return availableItems
}
```

**Category Filters:**
- **Armor Only:** Show defensive equipment (ArmorRating +0)
- **Weapons Only:** Show attack equipment (WeaponDamage +0)
- **Utility Only:** Show special items (effects, stat bonuses)
- **Affordable Only:** Show items within current credit balance
- **All Items:** No filtering, show complete inventory

**Stat Filter:**
- **Defense Focus:** Armor +3+, Shield items
- **Attack Focus:** Weapons +5+, Crit items
- **HP Focus:** Armor +2+, HP bonus items
- **MP/SP Focus:** High MP pool items, regen items
- **Movement Focus:** Movement bonus items, Jump Height items

**Sort Options:**
- **Price: Low to High, High to Low
- **Power:** weakest to strongest (Total SW rating)
- **Tier:** Common to Legendary
- **Recent:** Newly added items first
- **Alphabetical:** A to Z order

**Purchase System:**

**Transaction Flow:**
```go
func PurchaseEquipment(character Character, equipmentID uuid.UUID) error {
    // Get equipment details
    item := GetShopItem(equipmentID)
    
    // Validate affordability
    if character.Credits < item.Cost {
        return errors.New("Insufficient credits")
    }
    
    // Validate slot compatibility
    if !CanEquipInSlot(character, item.Slot) {
        return errors.New("Cannot equip in this slot")
    }
    
    // Process transaction
    character.Credits -= item.Cost
    character.Inventory.AddItem(item)
    
    // Log purchase
    LogTransaction(character.ID, "equipment_purchase", item.ID, item.Cost)
    
    return nil
}
```

**Pricing Formula:**

**Base Cost Calculation:**
- **Item Power:** Total positive SW of equipment properties
- **Tier Multiplier:** Based on item rarity/tier
  - Common (Tier 1): 1.0×
  - Uncommon (Tier 2): 1.5×
  - Rare (Tier 3): 2.0×
  - Epic (Tier 4): 3.0×
  - Legendary (Tier 5): 5.0×
- **Final Cost:** Item Power × Tier Multiplier

**Equipment Power Calculation:**
```go
func CalculateEquipmentPower(item Equipment) int {
    power := 0
    
    // Add offensive properties
    if item.WeaponDamage > 0 {
        power += item.WeaponDamage * 10  // Damage SW cost
    }
    if item.CritChance > 0 {
        power += item.CritChance * 2   // Crit SW cost
    }
    
    // Add defensive properties
    if item.ArmorRating > 0 {
        power += item.ArmorRating * 8   // Armor SW cost
    }
    if item.ShieldPower > 0 {
        power += item.ShieldPower * 10  // Shield SW cost
    }
    
    // Add utility properties
    if item.MP > 0 {
        power += item.MP  // MP SW cost
    }
    
    return power
}
```

**Inventory Management:**

**Owned Equipment:**
- **Equipped Items:** Currently in 3 slots (Armor, Utility, Weapon)
- **Unequipped Items:** In character inventory, available to equip
- **Sorting:** Default sort by power, then tier, then category
- **Filtering:** Filter by slot type, stat bonuses, affordability

**Unequipping System:**
```go
func UnequipItem(character Character, slot EquipmentSlot) error {
    // Check if slot is empty
    currentItem := character.GetEquippedItem(slot)
    if currentItem == nil {
        return errors.New("Slot is already empty")
    }
    
    // Remove item from slot
    character.Slots[slot] = nil
    
    // Add to inventory
    character.Inventory.AddItem(currentItem)
    
    // Recalculate stats
    character.RecalculateStats()
    
    // Log unequip
    LogTransaction(character.ID, "unequip", currentItem.ID, 0)
    
    return nil
}
```

**Equip Validation:**
```go
func EquipItem(character Character, item Equipment, slot EquipmentSlot) error {
    // Validate slot type
    if item.Slot != slot {
        return errors.New("Item cannot be equipped in this slot type")
    }
    
    // Check slot availability
    if character.GetEquippedItem(slot) != nil {
        return errors.New("Slot is already occupied")
    }
    
    // Remove existing item (if any)
    if existing := character.GetEquippedItem(slot); existing != nil {
        UnequipItem(character, slot)
    }
    
    // Equip new item
    character.Slots[slot] = item
    character.Inventory.RemoveItem(item.ID)
    
    // Apply stat bonuses
    character.RecalculateStats()
    
    return nil
}
```

**Shop Refresh System:**

**Dynamic Inventory:**
- **Daily Rotation:** 20% of items replaced daily
- **Weekly Events:** Special limited-time items appear
- **Seasonal Updates:** New tiers added based on player progress
- **Player Level Scaling:** Higher-level players see more items

**Special Offers:**
- **Skill Bundles:** Packages with skill + equipment
- **Discount Items:** Limited-time offers (20-50% off)
- **Trade-In System:** Exchange old equipment for credit value
- **Mystery Boxes:** Random equipment with guaranteed minimum tier

**UI Integration:**

**Shop Interface:**
- **Category Tabs:** Separate sections for Armor, Weapons, Utility
- **Filter Controls:** Level filter, stat filter, price range slider
- **Item Cards:** Display item properties, cost, and comparison
- **Affordability:** Show green for affordable, red for too expensive
- **Preview System:** Show stat changes if item equipped

**Equipment Comparison:**
```go
func DisplayEquipmentComparison(character Character, newItem Equipment) {
    currentItem := character.GetEquippedItem(newItem.Slot)
    
    // Calculate stat differences
    currentStats := character.GetStatsWithItem(currentItem)
    newStats := character.GetStatsWithItem(newItem)
    
    differences := CalculateStatDifferences(currentStats, newStats)
    
    // Display comparison
    ShowComparisonCard(currentItem, newItem, differences)
}
```

**Transaction History:**

**Purchase Tracking:**
```go
type ShopTransaction struct {
    TransactionID   uuid.UUID
    CharacterID     uuid.UUID
    ItemID         uuid.UUID
    ItemName        string
    Cost            int
    PurchaseType    string  // "equipment", "skill", "reforge"
    Timestamp       time.Time
}

func GetTransactionHistory(characterID uuid.UUID) []ShopTransaction {
    return shopDatabase.Where("character_id = ?", characterID).OrderBy("timestamp DESC")
}
```

**Player Experience:**

**Shopping Journey:**
"I'm Level 5 and have 250 credits. The shop shows me items up to Level 7 - mostly Common and some Uncommon. I see this Iron Sword (+5 Damage, Range 2) for 200 credits—it would replace my current Wooden Sword (+2 Damage). The tier system and level filtering make it easy to find upgrades!"

**Stat-Based Filtering:**
"I want to build a tank character, so I filter by Armor Rating and Defense. I see this Plate Chest (+8 Defense, +10 ArmorRating) for 800 credits—pretty expensive! But it would make me much more durable. The power calculation seems balanced compared to weapons."

**Strategic Planning:**
"I'm saving up for a Rare tier item. They cost 500-1000 credits but have significant bonuses. I notice the Legendary items are 5× the base cost—really expensive! I wonder if I'll ever earn enough credits for those. The shop system makes high-tier items feel like major achievements!"

**Implementation Priority:** HIGH - Required for Phase 3 economy integration

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[equipment_shop_inventory_system]]`
- **Related Files:** Shop management UI, equipment database, character inventory logic
- **Integration:** Works with `mec_credit_spending_shop`, `entity_equipment_system`, `api_equipment_management`

## EXPECTATION
