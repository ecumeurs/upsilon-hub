---
id: mechanic_character_point_buy_system
status: DRAFT
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents: []
dependents: []
---

# New Atom

## INTENT
To implement character creation point-buy system where players receive 100 Character Points (CP) to strategically allocate attributes instead of receiving 4 random points on base stats.

## THE RULE / LOGIC
**Character Point-Buy System:**

**Base Points Allocation:**
- **Starting Pool:** 100 Character Points (CP) instead of 4 random points
- **Base Stats:** HP 30-50, Attack 10, Defense 5, Movement 3 (new V2 baseline)
- **Strategic Allocation:** Players choose how to spend CP among all attributes

**Attribute Costs (CP per point increase):**
- **HP (+1):** Cost 1 CP (Linear, affordable)
- **Attack (+1):** Cost 5 CP (Powerful stat)
- **Defense (+1):** Cost 5 CP (Equally powerful mitigation)
- **Movement (+1 cell):** Cost 30 CP (Most potent stat, naturally restricted)
- **Critical Chance (+1%):** Cost 10 CP (High value, caps at reasonable maximum)
- **Critical Multiplier (+5%):** Cost 5 CP (DPS enhancement)
- **Jump Height (+1):** Cost 15 CP (Alters terrain navigation significantly)

**Cost Rationale:**
- **Linear vs Power:** HP is linear scaling (cheapest), Attack/Defense have multiplicative effects (expensive)
- **Movement Premium:** 30 CP cost creates natural restriction without hard-coded locks
- **Exotic Costs:** Higher costs for powerful tactical stats (crit, dodge, jump)

**Point Allocation Rules:**
- **Minimum Required:** Must allocate all 100 CP during character creation
- **Stat Minimums:** Cannot reduce base stats below minimum thresholds
- **Maximum Caps:** Some stats have soft caps to prevent abuse
- **Exotic Limits:** Exotic attributes (crit, jump) have lower maximums than standard stats

**Progression Impact:**
- **Starting Character:** 100 CP allocation + V2 base stats = meaningful variety
- **Win Rewards:** +10 CP per win (instead of +1 point)
- **Total Cap:** 100 + (total_wins × 10)
- **Character Power:** Dramatically higher than V1, enabling meaningful skill percentage effects

**Validation Logic:**
```go
// Character creation validation
func ValidatePointAllocation(cpAllocated int, baseStats CharacterStats) error {
    if cpAllocated > 100 {
        return errors.New("Cannot exceed 100 CP allocation")
    }
    
    // Check stat minimums
    if baseStats.HP < 30 || baseStats.Attack < 10 {
        return errors.New("Stats below minimum thresholds")
    }
    
    // Validate exotic stat caps
    if baseStats.CritChance > MAX_CRIT_CHANCE {
        return errors.New("Exceeded critical chance maximum")
    }
    
    return nil
}
```

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[character_point_buy_system]]`
- **Related Files:** Character creation logic, character database schema
- **API Endpoints:** `POST /api/v1/character/create` (with CP allocation)
- **UI Components:** Character creation form with CP spending display

## EXPECTATION
