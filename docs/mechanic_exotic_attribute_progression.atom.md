---
id: mechanic_exotic_attribute_progression
status: DRAFT
priority: 5
version: 2.0
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To implement exotic attribute progression mechanics for V2 characters, defining how attributes like Critical Chance, Critical Multiplier, Dodge, Accuracy, and Jump Height increase through level-based progression and character point allocation.

## THE RULE / LOGIC
**Exotic Attribute Progression System:**

**Exotic Attributes Definition:**
Exotic attributes are non-standard stats that provide tactical advantages beyond basic HP, Attack, Defense, and Movement.

**Critical Chance:**
- **Base Value:** 0% (no guaranteed criticals)
- **Increment:** +1% per 1 point allocation
- **Maximum Cap:** 25-30% (game balance to prevent guaranteed crits)
- **Level Restriction:** Can allocate points once every 5 character levels
- **Point Cost:** 10 CP per +1%
- **Tactical Value:** High damage spikes, but unreliable

**Critical Multiplier:**
- **Base Value:** 100% (1x damage on crit)
- **Increment:** +5% per 1 point allocation
- **Maximum Cap:** 150-200% (1.5x-2x damage on crit)
- **Level Restriction:** No restriction (can allocate every level)
- **Point Cost:** 5 CP per +5%
- **Tactical Value:** Increases average damage, synergizes with high Crit Chance

**Dodge:**
- **Base Value:** 0% (no guaranteed dodges)
- **Increment:** +1% per 1 point allocation
- **Maximum Cap:** 25-30% (prevents untouchable characters)
- **Level Restriction:** Can allocate points once every 3 character levels
- **Point Cost:** 5 CP per +1%
- **Tactical Value:** Defensive counter to Accuracy, reduces incoming damage

**Accuracy:**
- **Base Value:** 100% (standard hit chance)
- **Increment:** +2% per 1 point allocation
- **Maximum Cap:** 100-120% (no misses)
- **Level Restriction:** Can allocate points once every 3 character levels
- **Point Cost:** 3 CP per +2%
- **Tactical Value:** Counter to Dodge, ensures skill-based attacks hit

**Jump Height:**
- **Base Value:** 2 (current V1 baseline)
- **Increment:** +1 per 1 point allocation
- **Maximum Cap:** 4-5 (map balance, prevents flying over entire board)
- **Level Restriction:** Can allocate points once every 5 character levels
- **Point Cost:** 15 CP per +1
- **Tactical Value:** Terrain navigation advantage, accessing high-ground positions

**Progression Mechanics:**

**Level-Based Availability:**
- **Every Level:** Attack, Defense, HP, Critical Multiplier
- **Every 3 Levels:** Accuracy, Dodge
- **Every 5 Levels:** Critical Chance, Jump Height, Movement

**Stat Interaction Rules:**
- **Crit vs Dodge:** High Crit Chance benefits from Accuracy increase
- **Accuracy vs Dodge:** Both stats compete for combat effectiveness
- **Jump Height Synergy:** High Jump Height enables better positioning for Crit/Accuracy
- **Movement vs Jump:** Both cost expensive (30 CP vs 15 CP), natural restriction

**Exotic Attribute Example Builds:**

**Glass Cannon Build:**
- High Attack (expensive CP)
- High Critical Chance + Critical Multiplier
- Low Defense + HP
- High Accuracy to maximize crit hits

**Evasive Tank Build:**
- High HP + Defense
- High Dodge
- High Jump Height for positioning
- Lower Attack, rely on survival

**Precision DPS Build:**
- Moderate Attack
- High Critical Chance + Multiplier
- High Accuracy (every 3 levels)
- Balanced Defense

**Progression Timeline:**
- **Level 1:** Can allocate 1 Crit Chance (every 5 levels restriction)
- **Level 3:** Can allocate Accuracy + Dodge (every 3 levels restriction)
- **Level 5:** Can allocate Jump Height + Movement + Critical Chance (every 5 levels restriction)
- **Level 10:** Multiple exotic options available, build differentiation emerges

**Balance Considerations:**
- **Expensive Costs:** All exotic attributes cost significant CP
- **Level Restrictions:** Prevents min-maxing at early levels
- **Maximum Caps:** Ensures no stat becomes game-breakingly powerful
- **Tactical Trade-offs:** Players must choose between offense (crit) and defense (dodge)

```go
func GetExoticStatLevelRestriction(statType string) int {
    switch statType {
        case "CriticalChance": return 5
        case "JumpHeight": return 5
        case "Movement": return 5
        case "Accuracy": return 3
        case "Dodge": return 3
        default: return 1 // Standard stats: no restriction
    }
}
```

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[exotic_attribute_progression]]`
- **Related Files:** Character progression system, stat validation logic
- **UI Components:** Character sheet with exotic attributes display and upgrade buttons

## EXPECTATION
