---
id: mechanic_character_creation_integration
status: DRAFT
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents: []
dependents: []
type: MECHANIC
---

# New Atom

## INTENT
To implement character creation and progression integration that bridges V1 legacy characters to V2 systems, handles data migration, and establishes new character creation flow with 100 CP point-buy allocation.

## THE RULE / LOGIC
**Character Creation & Progression Integration:**

**V1 to V2 Migration:**

**Legacy Character Handling:**
- **Existing Characters:** Receive automatic stat rebalancing to V2 baseline
- **Point Compensation:** Awarded difference between old random allocation and optimal 100 CP
- **Skill Compensation:** Awarded base skills based on character level
- **Equipment Compensation:** Awarded starter equipment based on progression

**Data Migration Process:**
```go
func MigrateV1Character(v1Character Character) V2Character {
    // V1 Stats: HP 3, Attack 1, Defense 1, Movement 1
    // V2 Baseline: HP 30-50, Attack 10, Defense 5, Movement 3
    
    v2Char := V2Character{
        BaseStats: V2Baseline(),
        AvailableCP: 100, // Reset to full pool
        MigratedFrom: v1Character.ID,
    }
    
    // Award skills based on level
    if v1Character.Level >= 10 {
        v2Char.GrantStarterSkills(Grade: "I-II")
    }
    
    // Award equipment based on progression
    if v1Character.Wins > 0 {
        v2Char.GrantStarterEquipment(Tier: CalculateV2Tier(v1Character))
    }
    
    return v2Char
}
```

**New Character Creation Flow:**

**Step 1: Basic Information**
- **Character Name:** Unique name validation
- **Appearance:** Visual customization (avatar, colors)
- **Starting Level:** Level 1 (future: level-based creation)

**Step 2: Stat Allocation**
- **Available CP:** Display 100 CP pool
- **Current Allocation:** Show points spent per attribute
- **Cost Display:** Show CP cost per +1 increment
- **Restriction Warnings:** Highlight unavailable exotic stats based on level
- **Validation:** Ensure exactly 100 CP spent, no minimum violations

**Step 3: Skill Selection**
- **Available Skills:** Show 3 random skills from Grade I-II pool
- **Skill Preview:** Display skill properties, effects, costs
- **Selection Required:** Must select 1 skill before proceeding
- **Skill Assignment:** Add selected skill to character's skill list

**Step 4: Equipment Selection (Optional)**
- **Starter Equipment:** Show basic tier equipment options
- **Slot Preview:** Show what equipment fills which slot
- **Stat Bonuses:** Display how equipment affects final stats
- **Skip Option:** Can proceed without selecting equipment

**Step 5: Character Confirmation**
- **Final Stats:** Display complete character sheet with V2 stats
- **Total CP:** Confirm exactly 100 CP allocated
- **Preview:** Show character appearance and equipment
- **Create Character:** Finalize character with V2 systems

**Progression Integration:**

**Post-Creation Progression:**
- **Level Up:** Awarded +10 CP (instead of +1 point)
- **CP Pool:** Unspent CP accumulates for later allocation
- **Skill Selection:** Every 10 levels, choose 1 of 3 skills
- **Skill Reforging:** Every 5 levels, modify existing skills for credits
- **Exotic Stat Access:** Unlock new exotic stat options as level increases

**Character Data Structure (V2):**
```go
type V2Character struct {
    // V1 Fields (Legacy)
    ID           uuid.UUID
    Name         string
    Level        int
    Wins         int
    
    // V2 Fields (New Systems)
    BaseStats    CharacterStats     // HP 30-50, Attack 10, Defense 5, Movement 3
    AvailableCP  int               // 100 + (wins × 10)
    ExoticStats  ExoticStats      // Crit, Dodge, Accuracy, Jump
    Skills       []Skill           // Skill inventory
    Equipment    EquipmentSlots     // 3-slot system (armor, utility, weapon)
    Credits      int               // Currency for shop purchases
    
    // Migration Fields
    MigratedFrom uuid.UUID         // Legacy character ID if migrated
    MigrationBonus int             // Additional CP from migration compensation
}
```

**Validation Rules:**

**Creation Validation:**
- **CP Allocation:** Must spend exactly 100 CP, no more, no less
- **Stat Minimums:** HP ≥30, Attack ≥10, Defense ≥5, Movement ≥3
- **Exotic Caps:** Cannot exceed maximum values for exotic attributes
- **Level Restrictions:** Respect once-every-X-levels restrictions

**Progression Validation:**
- **CP Available:** Cannot spend more CP than accumulated
- **Skill Access:** Only offer skills available at current level
- **Equipment Equipping:** Validate slot compatibility and stat effects

**Migration Compensation:**
- **Fair Transition:** Legacy players don't lose progression value
- **CP Adjustment:** Calculate difference between optimal V2 build and current V1 stats
- **Skill Compensation:** Awarded based on V1 wins/level
- **Equipment Compensation:** Awarded based on V1 progression

```go
func CalculateMigrationCompensation(v1Character Character) int {
    // Calculate optimal V2 build for same "type" of character
    optimalV2CP := CalculateOptimalV2Allocation(v1Character.Type)
    
    // Calculate current V1 equivalent CP value
    v1EquivalentCP := ConvertV1StatsToCP(v1Character.Stats)
    
    // Compensation is difference
    return optimalV2CP - v1EquivalentCP
}
```

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[character_creation_integration]]`
- **Related Files:** Character creation API, migration scripts, database schema
- **API Endpoints:** `POST /api/v1/character/create`, `GET /api/v1/character/migrate`
- **UI Components:** Character creation wizard, migration notification interface

## EXPECTATION
