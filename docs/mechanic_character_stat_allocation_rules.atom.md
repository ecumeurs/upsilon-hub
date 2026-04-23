---
id: mechanic_character_stat_allocation_rules
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
To establish character stat allocation rules defining point costs, maximum caps, and level restrictions for both standard and exotic attributes in V2 progression system.

## THE RULE / LOGIC
**Stat Allocation Rules:**

**Standard Attribute Rules:**

**HP (Health Points):**
- **Cost:** +1 HP = 1 CP
- **Range:** 30-50 base, unlimited cap via progression
- **Rationale:** Linear scaling, cheap to encourage tank builds
- **Special:** No maximum cap, but diminishing returns on very high HP

**Attack:**
- **Cost:** +1 Attack = 5 CP
- **Range:** 10 base, no hard cap (soft caps via game balance)
- **Rationale:** Expensive due to multiplicative damage scaling
- **Impact:** Every +1 Attack significantly increases damage output

**Defense:**
- **Cost:** +1 Defense = 5 CP
- **Range:** 5 base, soft cap at Attack value (prevent over-tanking)
- **Rationale:** Expensive due to damage mitigation effectiveness
- **Balance:** Defense + Attack both cost same CP for balance

**Movement:**
- **Cost:** +1 cell Movement = 30 CP (Natural restriction)
- **Range:** 3 base, maximum 5-6 cells (game balance)
- **Rationale:** Movement is most powerful stat, naturally expensive
- **Level Restriction:** Can only increase Movement once every 5 levels

**Exotic Attribute Rules:**

**Critical Chance:**
- **Cost:** +1% Critical Chance = 10 CP
- **Range:** 0% base, maximum 25-30% (balance cap)
- **Level Restriction:** Can only increase Critical Chance once every 5 levels
- **Rationale:** High value stat with exponential power increase

**Critical Multiplier:**
- **Cost:** +5% Critical Multiplier = 5 CP
- **Range:** 100% base (1x), maximum 150-200% (2x)
- **Level Restriction:** Can increase every level (no restriction)
- **Rationale:** DPS enhancement, reasonable cost

**Accuracy:**
- **Cost:** +2% Accuracy = 3 CP (if implemented in V2)
- **Range:** 100% base, maximum 100-120%
- **Level Restriction:** Can increase every 3 levels
- **Rationale:** Skill-based enhancement

**Dodge:**
- **Cost:** +1% Dodge = 5 CP (if implemented in V2)
- **Range:** 0% base, maximum 25-30%
- **Level Restriction:** Can increase every 3 levels
- **Rationale:** Defensive counter to accuracy

**Jump Height:**
- **Cost:** +1 Jump Height = 15 CP
- **Range:** 2 base, maximum 4-5 (map balance)
- **Level Restriction:** Can increase Jump Height once every 5 levels
- **Rationale:** Drastically alters terrain navigation and tactics

**Progression Rules:**

**Level-Based Restrictions:**
- **Once Every 5 Levels:** Movement, Critical Chance, Jump Height
- **Every 3 Levels:** Accuracy, Dodge (if implemented)
- **No Restriction:** Attack, Defense, Critical Multiplier, HP

**Stat Interactions:**
- **Attack vs Defense:** Higher Attack benefits more from +Attack than +Defense
- **HP Scaling:** HP benefits linearly from +1, no diminishing returns
- **Movement Premium:** 30 CP cost restricts high-movement builds naturally
- **Exotic Caps:** Critical and Dodge have maximums to prevent invincible builds

**Point Spending Validation:**
```go
func ValidateStatPurchase(currentStats CharacterStats, statType string, amount int) error {
    cpCost := GetStatCost(statType, amount)
    
    if currentStats.AvailableCP < cpCost {
        return errors.New("Insufficient Character Points")
    }
    
    // Check level restrictions
    if IsStatRestricted(statType, currentStats.Level) {
        lastIncreased := currentStats.GetLastIncreased(statType)
        levelsSinceIncrease := currentStats.Level - lastIncreased
        
        if levelsSinceIncrease < GetMinimumLevels(statType) {
            return errors.New("Stat restricted until higher level")
        }
    }
    
    // Check maximum caps
    newStatValue := currentStats.GetStat(statType) + amount
    if newStatValue > GetMaximum(statType) {
        return errors.New("Exceeded maximum stat cap")
    }
    
    return nil
}
```

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[character_stat_allocation_rules]]`
- **Related Files:** Character progression logic, stat validation system
- **UI Components:** Stat allocation interface with cost displays and restriction warnings

## EXPECTATION
