# Issue: Comprehensive Item System - Shop, Inventory, Equipment & Battle Integration

**ID:** `20260423_comprehensive_item_system`
**Ref:** `ISS-074`
**Date:** 2026-04-23
**Severity:** High
**Status:** Open
**Component:** `battleui`, `upsilonapi`, `upsilonbattle`
**Affects:** Database schema, shop UI, credit spending, battle engine, character progression

---

## Summary

Implement end-to-end item system for V2: fixed shop catalog, normalized player inventory, 3-slot equipment system (armor/utility/weapon), and buff-based battle integration. Players earn credits, purchase items, manage inventory, equip items to characters, and receive stat bonuses in battle through the existing buff system.

**This issue consolidates requirements from:**
- ISS-068 (Equipment System & Weapon-as-Skill)
- ISS-074 (Simple Shop Inventory)
- ISS-075 (Player Inventory System)
- ISS-076 (Character Data Transfer for Battle Engine)

---

## Technical Description

### Background

Credit economy system exists (ISS-067) but lacks items to purchase. No inventory system, no equipment slots, no battle integration. Item properties exist in code (`property.ItemProperties`) but no end-to-end flow from purchase to battle.

### The Problem Scenario

1. **No Shop**: Credits earned but nothing to spend on
2. **No Inventory**: No way to track owned items, purchase history, or usage stats
3. **No Equipment**: No 3-slot system (armor/utility/weapon) for equipping items
4. **No Battle Integration**: Purchased items don't affect combat
5. **No Player Progression**: Equipment doesn't contribute to character power curve

### Overall Architecture: 4-Layer Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   SHOP      │────▶│  INVENTORY  │────▶│ EQUIPMENT   │────▶│   BATTLE    │
│ (Catalog)   │     │ (Owned)     │     │  (Equipped) │     │  (Buffs)    │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
     Credits        player_inventory   character_slots   entity.Buffs
                     shop_items         3-slot system    Forever=true
```

### Layer 1: Shop System

**Fixed Item Catalog (V2.0):**
- **Armor:** "Basic Armor" - +5 Armor Rating - 200 credits
- **Weapon:** "Basic Sword" - +5 Weapon Rating (Damage) - 300 credits
- **Movement:** "Swift Boots" - +1 Movement - 150 credits

**Shop Items Table:**
```sql
CREATE TABLE shop_items (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type ENUM('armor', 'weapon', 'movement', 'utility') NOT NULL,
    slot ENUM('armor', 'weapon', 'utility') NOT NULL,
    properties JSON NOT NULL,          -- Item stat bonuses
    cost INTEGER NOT NULL,              -- Credit cost
    available BOOLEAN DEFAULT TRUE,
    version VARCHAR(10) DEFAULT '2.0'
);
```

**Shop API Endpoints:**
- `GET /api/v1/shop/items` - Browse available items
- `POST /api/v1/shop/purchase` - Purchase item (deducts credits)

### Layer 2: Inventory System

**Normalized Tables:**
```sql
-- Player inventory (owned items)
CREATE TABLE player_inventory (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shop_item_id UUID NOT NULL REFERENCES shop_items(id),
    character_id UUID REFERENCES characters(id),  -- NULL = in inventory, set = equipped
    quantity INTEGER DEFAULT 1,
    purchased_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(player_id, shop_item_id, character_id)
);

-- Purchase history (audit trail)
CREATE TABLE inventory_transactions (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES users(id),
    shop_item_id UUID NOT NULL REFERENCES shop_items(id),
    quantity INTEGER NOT NULL,
    credits_spent INTEGER NOT NULL,
    transaction_type ENUM('purchase', 'refund', 'gift', 'admin_grant') DEFAULT 'purchase',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Item usage statistics (optional, V2.1+)
CREATE TABLE item_usage_stats (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES users(id),
    shop_item_id UUID NOT NULL REFERENCES shop_items(id),
    uses_total INTEGER DEFAULT 0,
    damage_dealt INTEGER DEFAULT 0,
    last_used_at TIMESTAMP
);
```

**Inventory API Endpoints:**
- `GET /api/v1/player/inventory` - List all owned items
- `GET /api/v1/character/{id}/equipped` - List items equipped on specific character
- `POST /api/v1/inventory/equip/{itemId}/{characterId}` - Equip item to character
- `DELETE /api/v1/inventory/unequip/{itemId}/{characterId}` - Remove item from character
- `GET /api/v1/inventory/stats` - View usage statistics (V2.1+)

### Layer 3: Equipment System

**3-Slot System:**
- **Armor Slot:** One piece of armor (Head, Body, Hands, Legs, Feet)
- **Utility Slot:** One utility item (Neck, Ring, Belt, Movement)
- **Weapon Slot:** One weapon (Melee or Ranged, One or Two-Handed)

**Character Equipment Table:**
```sql
CREATE TABLE character_equipment (
    character_id UUID PRIMARY KEY REFERENCES characters(id) ON DELETE CASCADE,
    armor_item_id UUID REFERENCES player_inventory(id),
    utility_item_id UUID REFERENCES player_inventory(id),
    weapon_item_id UUID REFERENCES player_inventory(id)
);
```

**Equipment Constraints:**
- Exactly one item per slot type
- Slot independence (armor, utility, weapon distinct)
- Mutual exclusivity (equipping new item replaces old in slot)
- Two-handed weapons may restrict utility slot (V2.1+)

**Equipment API Endpoints:**
- `POST /api/v1/character/{id}/equip-armor/{itemId}` - Equip armor
- `POST /api/v1/character/{id}/equip-utility/{itemId}` - Equip utility
- `POST /api/v1/character/{id}/equip-weapon/{itemId}` - Equip weapon
- `DELETE /api/v1/character/{id}/unequip/{slot}` - Unequip item from slot

### Layer 4: Battle Integration (Buff-Based)

**IMPORTANT: Items are implemented as BUFFS with `Forever=true`**

Items should not directly modify entity base properties. Instead, they are registered as buffs with the existing buff system.

**Character Transfer Schema (from Laravel to Go Engine):**
```go
type CharacterData struct {
    // Existing fields
    ID       string   // Character ID
    PlayerID string   // Owner
    Team     int      // Team assignment
    Name     string   // Character name

    // V2 Stats
    HP       int
    MaxHP    int
    Attack   int
    Defense  int
    Move     int
    MaxMove  int

    // V2 Skill System
    EquippedSkills []uuid.UUID    // Skills selected for battle
    SkillSlots     int              // Total available slots

    // V2 Item System (NEW)
    EquippedItems []EquippedItem  // Items equipped on character
}

type EquippedItem struct {
    ItemID    uuid.UUID     // Shop item reference
    Name      string        // Item name
    Type      ItemType      // Armor, Weapon, Movement, Utility
    Slot      EquipmentSlot // armor, weapon, utility
    Properties map[string]property.Property  // Item stat bonuses
}

type EquipmentSlot string

const (
    SlotArmor   EquipmentSlot = "armor"
    SlotWeapon  EquipmentSlot = "weapon"
    SlotUtility EquipmentSlot = "utility"
)
```

**Entity Initialization with Items as Buffs:**
```go
// When creating entity for battle, load equipped items as buffs
func NewEntityFromCharacter(char Character) Entity {
    entity := NewEntity()

    // Set base stats
    entity.RepsertPropertyValue(property.HP, char.HP)
    entity.RepsertPropertyValue(property.MaxHP, char.MaxHP)
    entity.RepsertPropertyValue(property.Attack, char.Attack)
    entity.RepsertPropertyValue(property.Defense, char.Defense)
    entity.RepsertPropertyValue(property.Move, char.Move)

    // Load equipped items as buffs
    for _, item := range char.EquippedItems {
        buff := property.MakeTemporaryProperties(0)
        buff.Forever = true  // Permanent while equipped
        buff.Properties = item.Properties
        buff.OriginEntityID = item.ItemID
        entity.RegisterBuff(buff)
    }

    // Load equipped skills
    for _, skillID := range char.EquippedSkills {
        skill := GetSkillByID(skillID)
        entity.RegisterSkill(skill)
    }

    return entity
}
```

**Why Buff-Based Approach:**
- Uses existing `RegisterBuff()` / `BuffTickDown()` infrastructure
- Easy to remove on unequip: filter and remove buff by `OriginEntityID`
- Consistent with skill effects (poison, stun, etc.)
- Tracked via buff origin for debugging
- `Forever=true` means no duration management needed for equipment
- Property resolution via `entity.GetProperty()` automatically applies buff bonuses

**Unequipping (Character Sheet):**
```go
func (e *Entity) RemoveItemBuff(itemID uuid.UUID) {
    nbbuf := make([]property.TemporaryProperties, 0)
    for _, buff := range e.Buffs {
        if buff.OriginEntityID != itemID {
            nbbuf = append(nbbuf, buff)
        }
    }
    e.Buffs = nbbuf
}
```

**Enhanced Arena Start Request:**
```go
type ArenaStartRequest struct {
    MatchID     string
    CallbackURL string
    Players     []CharacterData  // Enhanced with skills and items
}
```

**Request Structure:**
```json
POST /arena/start
{
    "match_id": "...",
    "callback_url": "...",
    "players": [
        {
            "id": "char-1",
            "player_id": "player-1",
            "team": 1,
            "name": "Warrior",
            "hp": 35, "max_hp": 35,
            "attack": 12, "defense": 6, "move": 3, "max_move": 3,

            "equipped_skills": ["skill-1-id", "skill-2-id"],
            "skill_slots": 2,

            "equipped_items": [
                {
                    "item_id": "armor-1",
                    "name": "Basic Armor",
                    "type": "armor",
                    "slot": "armor",
                    "properties": {"ArmorRating": 5}
                },
                {
                    "item_id": "weapon-1",
                    "name": "Basic Sword",
                    "type": "weapon",
                    "slot": "weapon",
                    "properties": {"WeaponBaseDamage": 5}
                }
            ]
        }
    ]
}
```

### Weapon-as-Skill System (V2.1+)

Basic attacks transform when weapon is equipped:
```go
// When weapon equipped, basic attack becomes skill-based
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

### Privacy Considerations

- Player only receives information for their own characters
- Enemy character properties are hidden during match (only visible on reveal)
- Skill/item properties not exposed for unowned characters
- No API endpoint allows viewing detailed properties of characters not belonging to current user
- Inventory only visible to owner
- Character equipment hidden from enemy players (properties masked)
- Usage statistics only visible to owner

### Data Integrity Rules

**Cascading Deletes:**
- When user deleted, all inventory rows cascade to delete
- When character deleted, their equipped items remain in inventory (character_id set to NULL)

**Quantity Rules:**
- Maximum 99 of any item to prevent exploit
- Stacking limited to 10 identical items (V2.1+)
- Cannot equip more than 1 of same type per slot

**Slot Type Mapping:**
- **Armor Slot:** Accepts ItemType = Wearable (Head, Body, Hands, Legs, Feet)
- **Utility Slot:** Accepts ItemType = Wearable (Neck, Ring, Belt) or Movement items
- **Weapon Slot:** Accepts ItemType = Wearable (OneHandedMelee, TwoHandedMelee, OneHandedRanged, TwoHandedRanged)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | Item properties exist, buff system exists, need end-to-end integration |

---

## Recommended Fix

**Week 1 - Shop & Inventory:**
- Create shop_items, player_inventory, inventory_transactions tables
- Implement basic shop catalog with 3 fixed items
- Add purchase endpoint with credit deduction and transaction logging
- Create inventory view API

**Week 2 - Equipment System:**
- Create character_equipment table
- Implement 3-slot equipment system (armor/utility/weapon)
- Add equip/unequip API endpoints
- Validate slot constraints and item types
- Create equipment management UI

**Week 3 - Battle Integration:**
- Update ArenaStartRequest with EquippedItems
- Implement buff-based item stat bonuses to entity initialization
- Update Go engine's NewEntityFromCharacter()
- Test end-to-end: purchase → equip → battle

**Week 4 - Polish & V2.1 Prep:**
- Expand shop catalog with more items
- Add item usage statistics tracking
- Build inventory filtering and sorting
- Implement weapon-as-skill system (V2.1)

---

## User Experience Flow

**Player Preparation:**
1. Player earns credits from combat (ISS-067)
2. Player visits shop, sees 3 items available
3. Player purchases "Basic Armor" (200 credits), "Basic Sword" (300 credits)
4. Items added to player_inventory
5. Player equips items to character:
   - Select "Basic Armor" → equip to Armor Slot
   - Select "Basic Sword" → equip to Weapon Slot
6. Player selects 2 skills from inventory (2/3 slots used)

**Battle Init:**
1. Server receives character data with skills and items
2. Go engine creates entity with:
   - 2 registered skills
   - 2 item buffs (ArmorRating +5, WeaponBaseDamage +5)
3. Entity is ready for combat with full V2 systems

**Combat:**
1. Skills appear as action buttons (only equipped skills)
2. Attacks include weapon bonuses (via buff system)
3. All purchased items contribute to character power
4. Credits earned based on performance (including equipment contributions)

---

## References

**Parent Issues:**
- `ISS-067` (Credit Economy & Shop) - Defines credit earning system

**Superseded Issues:**
- `ISS-068` (Equipment System & Weapon-as-Skill) - 3-slot and weapon-as-skill requirements merged here
- `ISS-075` (Player Inventory System) - Inventory and transaction requirements merged here
- `ISS-076` (Character Data Transfer) - Battle integration schema merged here

**ATD Documentation:**
- `docs/rule_credit_earning_damage.atom.md` - Credit earning (base for pricing)
- `docs/entity_player_credits.atom.md` - Player credit tracking
- `docs/entity_equipment_system.atom.md` - Equipment system architecture
- `docs/mech_equipment_stat_bonuses.atom.md` - Equipment stat bonuses
- `docs/mech_three_slot_equipment_system.atom.md` - 3-slot equipment rules
- `docs/mec_weapon_as_skill_system.atom.md` - Weapon-to-skill conversion
- `docs/rule_item_pricing_simple.atom.md` - Item pricing formula

**Code References:**
- `upsilonbattle/battlearena/entity/entity.go` - Entity properties, buff system, GetProperty methods
- `upsilonbattle/battlearena/property/buff.go` - TemporaryProperties with Forever flag
- `upsilonbattle/battlearena/property/def/item.go` - Item property definitions
- `upsilonbattle/battlearena/property/propertyenum.go` - Property definitions
