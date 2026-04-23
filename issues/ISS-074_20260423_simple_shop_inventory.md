# Issue: Simple Shop Inventory

**ID:** `20260423_simple_shop_inventory`
**Ref:** `ISS-074`
**Date:** 2026-04-23
**Severity:** High
**Status:** Open
**Component:** `battleui`, `upsilonapi`
**Affects:** Database schema, shop UI, credit spending

---

## Summary

Implement minimal shop system with fixed item catalog for V2 testing: one armor (+5 armor rating), one weapon (+5 weapon rating), one movement item (+1 movement). Items are priced based on simple formula and can be purchased by players with credits.

---

## Technical Description

### Background

Credit economy system is planned (ISS-067) but lacks concrete items to purchase. For V2 testing and initial progression, a simple fixed catalog allows players to test spending mechanics without complex procedural generation.

### The Problem Scenario

1. **Shop has no items**: Credits can be earned but nothing to spend them on
2. **No pricing model**: Items don't have defined cost
3. **No player inventory**: Players can't track what they've purchased
4. **No Go engine integration**: Items purchased but no way to use them in battle

### Simple Shop Inventory

**Fixed Item Catalog:**
- **Armor Item:** "Basic Armor" - +5 Armor Rating
- **Weapon Item:** "Basic Sword" - +5 Weapon Rating (Damage)
- **Movement Item:** "Swift Boots" - +1 Movement

**Pricing Formula:**
- Simple fixed costs (not SW-based yet)
- Armor: 200 credits
- Weapon: 300 credits
- Movement: 150 credits

**Item Properties:**
```go
type ShopItem struct {
    ID          uuid.UUID
    Name        string
    Type        ItemType      // Armor, Weapon, Movement
    Properties  map[string]property.Property
    Cost        int           // Credit cost
}
```

### Database Schema

**Player Inventory Table:**
```sql
CREATE TABLE player_inventory (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES users(id),
    item_id UUID NOT NULL REFERENCES shop_items(id),
    quantity INTEGER DEFAULT 1,
    purchased_at TIMESTAMP DEFAULT NOW()
);
```

**Shop Items Table:**
```sql
CREATE TABLE shop_items (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type ENUM('armor', 'weapon', 'movement') NOT NULL,
    properties JSON NOT NULL,
    cost INTEGER NOT NULL,
    available BOOLEAN DEFAULT TRUE
);
```

### API Endpoints

**Shop:**
- `GET /api/v1/shop/items` - Browse available items
- `POST /api/v1/shop/purchase` - Purchase item (deducts credits)

**Player Inventory:**
- `GET /api/v1/player/inventory` - View owned items
- `POST /api/v1/inventory/equip/{itemId}` - Equip item
- `DELETE /api/v1/inventory/unequip/{itemId}` - Unequip item

### Go Engine Integration

**IMPORTANT: Items are implemented as BUFFS with `Forever=true`**

Items should not directly modify entity base properties. Instead, they are registered as buffs with the existing buff system:

```go
// When creating entity for battle, load equipped items as buffs
func NewEntityFromCharacter(char Character) Entity {
    entity := NewEntity()

    // Load equipped armor as buff
    if char.ArmorItem != nil {
        buff := property.MakeTemporaryProperties(0)
        buff.Forever = true  // Permanent while equipped
        buff.Properties = char.ArmorItem.Properties
        buff.OriginEntityID = char.ArmorItem.ID
        entity.RegisterBuff(buff)
    }

    // Load equipped weapon as buff
    if char.WeaponItem != nil {
        buff := property.MakeTemporaryProperties(0)
        buff.Forever = true
        buff.Properties = char.WeaponItem.Properties
        buff.OriginEntityID = char.WeaponItem.ID
        entity.RegisterBuff(buff)
    }

    // Load equipped movement item as buff
    if char.MovementItem != nil {
        buff := property.MakeTemporaryProperties(0)
        buff.Forever = true
        buff.Properties = char.MovementItem.Properties
        buff.OriginEntityID = char.MovementItem.ID
        entity.RegisterBuff(buff)
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

**Unequipping:**
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

**Privacy Considerations:**
- Player only receives information for their own characters
- Enemy character properties are hidden during match (only visible on reveal)
- Skill/item properties not exposed for unowned characters
- No API endpoint allows viewing detailed properties of characters not belonging to current user
- "Inspect skill" command (planned) will have authentication check before showing details

### Where This Pattern Exists Today

- `docs/rule_credit_earning_damage.atom.md` - Credit earning (base for pricing)
- `docs/entity_player_credits.atom.md` - Player credit tracking
- `upsilonbattle/battlearena/entity/entity.go` - Entity properties, buff system, GetProperty methods
- `upsilonbattle/battlearena/property/buff.go` - TemporaryProperties with Forever flag
- `docs/mech_equipment_stat_bonuses.atom.md` - Equipment stat bonuses (maps to buff properties)
- `docs/entity_equipment_system.atom.md` - Equipment system architecture

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | Credit system in planning, no items to purchase |

---

## Recommended Fix

**Short term (Week 1):**
- Create shop_items and player_inventory tables
- Implement basic shop catalog with 3 fixed items
- Add purchase endpoint with credit deduction
- Create inventory view API

**Medium term (Week 2):**
- Implement item equipment/unequipment endpoints
- Add buff-based item stat bonuses to entity initialization
- Create simple shop UI with item cards

**Long term (Week 3-4):**
- Expand shop catalog with more items
- Add item rarity and upgrade systems
- Build complex inventory management (sorting, filtering)

---

## References

- `ISS-067` (Credit Economy & Shop) - Parent issue, defines credit system
- `docs/rule_credit_earning_damage.atom.md` - 1 HP = 1 credit (cost basis)
- `docs/mec_weapon_as_skill_system.atom.md` - Weapon-to-skill conversion (for weapon items)
- `upsilonbattle/battlearena/property/buff.go` - Buff system with Forever flag
