# Issue: Player Inventory System

**ID:** `20260423_player_inventory`
**Ref:** `ISS-075`
**Date:** 2026-04-23
**Severity:** Medium
**Status:** Open
**Component:** `battleui`, `upsilonapi`
**Affects:** Database schema, inventory management, ownership tracking

---

## Summary

Implement a proper relational inventory system for players to track all owned items (equipment, skills, etc.) with quantity tracking, purchase history, and proper foreign key relationships. This replaces JSON-based storage with normalized tables.

---

## Technical Description

### Background

Players can earn credits and purchase items, but no proper inventory system exists to track ownership. Using JSON columns in user/character tables is not scalable and doesn't support queries like "all items from all users" or "item usage statistics".

### The Problem Scenario

1. **No ownership tracking**: Can't query which items a user owns
2. **No quantity management**: Can't stack identical items
3. **No purchase history**: No record of when items were bought
4. **No statistics**: Can't track how often items are used
5. **Poor scalability**: JSON queries on large tables are slow

### Player Inventory Architecture

**Normalized Tables:**

```sql
-- Player inventory (owned items)
CREATE TABLE player_inventory (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES users(id),
    shop_item_id UUID REFERENCES shop_items(id),
    character_id UUID REFERENCES characters(id),  -- NULL = global items, set = equipped items
    quantity INTEGER DEFAULT 1,
    purchased_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(player_id, shop_item_id, character_id)  -- One row per item per character
);

-- Purchase history (audit trail)
CREATE TABLE inventory_transactions (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES users(id),
    shop_item_id UUID REFERENCES shop_items(id),
    quantity INTEGER NOT NULL,
    credits_spent INTEGER NOT NULL,
    transaction_type ENUM('purchase', 'refund', 'gift', 'admin_grant') DEFAULT 'purchase',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Item usage statistics (optional)
CREATE TABLE item_usage_stats (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES users(id),
    shop_item_id UUID NOT NULL REFERENCES shop_items(id),
    uses_total INTEGER DEFAULT 0,
    damage_dealt INTEGER DEFAULT 0,
    last_used_at TIMESTAMP
);
```

**Inventory Operations:**
- **Purchase**: Add to player_inventory, deduct credits, log transaction
- **Equip**: Update character_id on inventory row
- **Unequip**: Set character_id to NULL on inventory row
- **Stack**: Increase quantity when buying duplicate items
- **Transfer**: Move item between characters

**API Endpoints:**

**Inventory Management:**
- `GET /api/v1/player/inventory` - List all owned items
- `GET /api/v1/character/{id}/equipped` - List items equipped on specific character
- `POST /api/v1/inventory/equip/{itemId}/{characterId}` - Equip item to character
- `POST /api/v1/inventory/unequip/{itemId}/{characterId}` - Remove item from character

**Statistics:**
- `GET /api/v1/inventory/stats` - View usage statistics

**Purchase:**
- `POST /api/v1/shop/buy/{itemId}` - Buy item (global or character)
- Input: `{ "itemId": "...", "characterId": "uuid|null", "quantity": 1 }`

### Data Integrity

**Cascading Deletes:**
- When user deleted, all inventory rows cascade to delete
- When character deleted, their equipped items remain in inventory (character_id set to NULL)

**Quantity Rules:**
- Maximum 99 of any item to prevent exploit
- Stacking limited to 10 identical items
- Cannot equip more than 1 of same type per character (unless slots allow)

### Privacy & Security

- Players can only see their own inventory
- No endpoint to view other players' items
- Character equipment hidden from enemy players (properties masked)
- Usage statistics only visible to owner

### Where This Pattern Exists Today

- `battleui/app/Models/User.php` - User model (needs inventory relationship)
- `docs/entity_player_credits.atom.md` - Credit tracking system
- `ISS-074` (Simple Shop Inventory) - Shop items to purchase

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | None (no inventory system) |

---

## Recommended Fix

**Short term:** Create normalized inventory tables. Implement basic inventory API endpoints. Add purchase transaction logging.

**Medium term:** Implement item usage statistics. Add inventory filtering and sorting. Create equipment management endpoints.

**Long term:** Implement item durability and repair. Add inventory sharing between characters. Build inventory optimization for large catalogs.

---

## References

- `ISS-067` (Credit Economy) - Credits used for purchases
- `ISS-074` (Simple Shop) - Shop items to purchase
- `docs/rule_item_pricing_simple.atom.md` - Item pricing formula
- `docs/rule_credit_action_communication_layer.atom.md` - Transaction logging requirements
