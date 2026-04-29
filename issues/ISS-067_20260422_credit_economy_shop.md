# Issue: Credit Economy & Shop System

**ID:** `20260422_credit_economy_shop`
**Ref:** `ISS-067`
**Date:** 2026-04-22
**Severity:** High
**Status:** Resolved
**Component:** `upsilonapi/api`, `battleui`
**Affects:** `upsilonbattle/battlearena`, database schema

---

## Summary

Implement comprehensive credit economy with multiple earning mechanisms (damage, healing, support, status effects) and shop system for purchasing skills and equipment. Credits are earned through combat performance and spent on character progression.

---

## Technical Description

### Background

No economy system exists. Players have no way to earn or spend currency. Skill selection and equipment acquisition lack progression mechanics.

### The Problem Scenario

1. **Combat Ends**: Player deals damage, heals allies, mitigates damage
2. **No Reward**: No credit earning system exists
3. **No Spending**: No shop or purchase mechanisms
4. **No Progression**: Players can't acquire new skills or equipment

### Communication Layer

**Action Message Structure:**
All combat actions must include:
- **Request ID:** For traceability
- **Version:** Current system version for compatibility
- **Effect Back:** Result of the action (success/failure, modified values)
- **Credit Earned:** Credits earned from this action (associated with player ID)

**Credit Association:**
```go
type ActionResponse struct {
    RequestID  string     // Traceability
    Version     string     // System version
    Success     bool       // Action outcome
    Modified    Modification // Changed game state
    Credits     int         // Credits earned this action
    PlayerID    uuid.UUID   // Credit recipient
}
```

**Database Schema Updates:**
```sql
-- Users table: Add credits balance
ALTER TABLE users ADD COLUMN credits INTEGER DEFAULT 0;

-- Character table: Track character-specific credits (optional)
ALTER TABLE characters ADD COLUMN credits INTEGER DEFAULT 0;

-- Credit transactions audit (optional for debugging)
CREATE TABLE credit_transactions (
    id UUID PRIMARY KEY,
    player_id UUID REFERENCES users(id),
    character_id UUID REFERENCES characters(id),
    amount INTEGER NOT NULL,
    source VARCHAR(50), -- 'damage', 'healing', 'mitigation', 'status_effect'
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Credit Earning System

**Base Rule:** 1 HP damage = 1 credit (healing also earns credits)
- Credits associated with player ID in database
- Real-time balance updates after each action

**Support Credits:**
- Damage mitigation: 1 HP mitigated = 1 credit
- Shield caster earns credits when shield blocks damage
- Effect must track caster for proper credit assignment
- Credits credited to original caster's player ID

**Status Effect Credits (Flat Rate):**
- Poison/Stun/Buff: SkillWeight/10 credits per application
- 100 SW poison skill = 10 credits per poison application
- Applies at moment of effect application, not per-turn
- Credits credited to skill caster's player ID

**Effect Caster Tracking:**
```go
type Effect struct {
    Properties []property.Property
    Name       string
    CasterID   uuid.UUID  // Track creator for credit assignment
    OriginTime time.Time  // When effect was applied
}
```
- Effects remember caster until effect ends
- Credits go to original caster's player ID even if they die later
- Critical for shield/healing credit assignment

### Shop System

**V2.0 Simple Shop (ISS-074 - Comprehensive Item System):**
- **Fixed Catalog:** 3 items (Armor +5 def, Weapon +5 dmg, Movement +1 move)
- **Pricing:** Armor 200 credits, Weapon 300 credits, Movement 150 credits
- **Reference:** `docs/rule_item_pricing_simple.atom.md`
- **Full System:** Shop → Inventory → Equipment (3-slot) → Battle (buff-based)

**V2.1 Full Shop (Future):**
- **Skills:** Available based on player level/grade
- **Equipment:** Armor, Utility, Weapon categories
- **Reference:** `docs/rule_skill_grading_system.atom.md`

**Player Inventory (ISS-074):**
- **Proper Tables:** Normalized inventory, transactions, usage stats
- **Ownership Tracking:** quantity, purchase history, equipment slots
- **Equipment Management:** equip/unequip to characters (3-slot system)

**Pricing Formula (Reference):**
- Skills: Credit Cost = Total Positive SW × 2
- Items: Fixed costs (V2.0) or SW-based (V2.1+)

**Shop System Architecture:**
- Credit spending: Deduct from appropriate balance
- Purchase validation: Cannot exceed available credits
- Inventory assignment: Global items or character-specific

### API Endpoints (Updated)

**Credit & Profile:**
- `GET /api/v1/profile/credits` - Get player credit balance (existing endpoint, updated)
- `GET /api/v1/character/{id}/credits` - Get character-specific credits
- `GET /api/v1/credits/history` - Get credit transaction history

**Shop:**
- `GET /api/v1/shop/skills` - Browse purchasable skills
- `GET /api/v1/shop/equipment` - Browse purchasable equipment
- `POST /api/v1/shop/purchase` - Purchase skill/equipment (deducts credits)

**CLI Integration:**
- `credits balance` - Display current credit balance
- `credits history` - Show transaction history
- Output format: Human-readable with timestamps and sources

### UI Integration

**Dashboard Display:**
- Credit balance prominent on player dashboard
- Real-time updates when credits earned/spent
- Transaction history panel with filters (earned, spent, by skill)
- Visual indicators for credit earning sources (damage, healing, support)

**Where This Pattern Exists Today**

- `upsilonapi/api/input.go` - Player structure
- `battleui/app/Models/User.php` - User model (needs credits column)
- `docs/rule_progression.atom.md` - Current progression rules
- `docs/entity_player_credits.atom.md` - Credit system entity
- `docs/rule_credit_action_communication_layer.atom.md` - Action message protocol (NEW)
- `ISS-074` (Comprehensive Item System) - Shop, inventory, equipment, battle integration (NEW)
- `docs/rule_item_pricing_simple.atom.md` - Item pricing model (NEW)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Implement base credit earning (1 HP = 1 coin). Create credits field in character schema. Build basic shop UI.

**Medium term:** Add support credit earning (damage mitigation). Implement status effect credit formula. Create skill/equipment purchasing system.

**Long term:** Build advanced shop features (filters, recommendations). Implement credit economy balancing tools. Add credit caps and inflation controls.

---

## References

- `third_party_reply.md` - Credit earning discussion
- `V2_ARCHITECTURAL_DECISIONS.md` - Credit economy decisions
- `docs/entity_users.atom.md` - User database entity
