---
id: mechanic_cell_attached_effects
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
To implement cell-attached effects system where effects are bound to grid cells rather than entities, enabling zone-based mechanics like poison clouds, healing areas, movement modifiers, and environmental traps.

## THE RULE / LOGIC
**Cell-Attached Effects System:**

**Core Concept:**
Effects exist on grid cells rather than being attached to entities, enabling zone-based mechanics that affect multiple targets simultaneously.

**Effect Types:**

**1. Zone Effects (Duration-Based)**
- **Poison Clouds:** Damage entities entering/inside zone each turn
- **Healing Areas:** Restore HP to entities in zone each turn
- **Movement Modifiers:** Quagmire (slow), ice (slippery), speed zones
- **Duration:** Effects last N turns, then expire naturally

**2. Trap Effects (Trigger-Based)**
- **Damage Traps:** Deal damage when entity steps on cell
- **Status Traps:** Apply stun/poison when triggered
- **Movement Traps:** Apply movement penalties/damage when walking through
- **Single-Use:** Triggered traps expire immediately

**3. Environmental Effects (Permanent/Duration)**
- **Terrain Mods:** Temporary bridges, walls, barriers
- **Weather Effects:** Rain (movement penalty), fog (visibility reduction)
- **Interactive Objects:** Crates (destroyable), switches (toggle effects)
- **Duration:** Can be permanent map features or temporary

**Cell Effect Data Structure:**
```go
type CellEffect struct {
    EffectID        uuid.UUID
    CasterID        uuid.UUID    // Creator for credit tracking
    EffectType      EffectType  // Zone, Trap, Environmental
    
    // Zone/Trap properties
    ZoneType        string       // Poison, Healing, Quagmire, etc.
    TriggerCondition  TriggerType // OnStep, OnTurn, OnEntry, OnExit
    TargetType      TargetType  // Entity, Player, Team, All
    
    // Effect definition
    Intensity        property.Property // Damage amount, healing amount, etc.
    Duration         int           // How many turns effect lasts
    WalkThrough      bool          // Can entities pass through effect
    IsVisible       bool          // Traps can be hidden
    IsDestructible  bool          // Can effect be destroyed?
    HP              int           // If destructible, HP threshold
    
    // Visual properties
    Color            string        // Effect color tint
    Animation         string        // Effect animation to play
    Size             int           // Effect area size (radius, pattern)}
```

**Trigger Types:**

**OnStep Trigger:**
```go
func (ce *CellEffect) OnStep(entity Entity) {
    if ce.TriggerCondition == "OnStep" {
        // Apply effect when entity steps on cell
        ce.ApplyEffect(entity)
        
        // Some traps expire after triggering
        if ce.IsSingleUse {
            ce.Expire()
        }
    }
}
```

**OnTurn Trigger:**
```go
func (ce *CellEffect) OnTurn(entities []Entity) {
    if ce.TriggerCondition == "OnTurn" {
        // Apply effect to all entities in zone each turn
        for _, entity := range entities {
            if ce.IsEntityInZone(entity) {
                ce.ApplyEffect(entity)
            }
        }
    }
}
```

**OnEntry/OnExit Triggers:**
```go
func (ce *CellEffect) OnEntry(entity Entity) {
    if ce.TriggerCondition == "OnEntry" {
        ce.ApplyEffect(entity)
    }
}

func (ce *CellEffect) OnExit(entity Entity) {
    if ce.TriggerCondition == "OnExit" {
        ce.RemoveEffect(entity)
    }
}
```

**Specific Cell Effects:**

**Poison Cloud (Zone):**
```go
poisonCloud := CellEffect{
    EffectType:      "Zone",
    ZoneType:        "Poison",
    TriggerCondition:  "OnTurn",
    TargetType:      "All",              // Affects everyone
    Intensity:        MakeIntProperty("Damage", 5), // Per turn
    Duration:         3,                 // Lasts 3 turns
    WalkThrough:      true,
    IsVisible:       true,               // Green fog visual
    Color:           "#00FF00",
}

// Behavior: Damages everyone inside 3x3 area for 3 turns
```

**Healing Zone (Beneficial):**
```go
healingZone := CellEffect{
    EffectType:      "Zone",
    ZoneType:        "Healing",
    TriggerCondition:  "OnTurn",
    TargetType:      "Team",             // Only allies
    Intensity:        MakeIntProperty("Heal", 10), // Per turn
    Duration:         2,                 // Lasts 2 turns
    WalkThrough:      true,
    IsVisible:       true,               // Blue glow visual
    Color:           "#0000FF",
}

// Behavior: Heals allies in 3x3 area for 2 turns
```

**Quagmire (Movement Modifier):**
```go
quagmire := CellEffect{
    EffectType:      "Environmental",
    ZoneType:        "Quagmire",
    TriggerCondition:  "OnStep",           // Applies on movement through
    TargetType:      "Entity",
    Intensity:        MakeIntProperty("MovementCost", +2), // +2 cost
    Duration:         5,                 // Lasts 5 turns
    WalkThrough:      true,
    IsVisible:       true,               // Brown muddy visual
    Color:           "#8B4513",
}

// Behavior: Moving through costs +2 extra movement points
```

**Bear Trap (Damage Trap):**
```go
bearTrap := CellEffect{
    EffectType:      "Trap",
    ZoneType:        "Damage",
    TriggerCondition:  "OnStep",           // Triggers when stepped on
    TargetType:      "Enemy",            // Only enemies trigger
    Intensity:        MakeIntProperty("Damage", 20),
    Duration:         0,                 // Instant, expires after trigger
    WalkThrough:      true,
    IsVisible:       false,              // Hidden trap
    IsDestructible:  true,
    HP:              10,                 // Can be destroyed
}
```

**Effect Application Rules:**

**Zone Size and Shape:**
- **Single Cell:** Effect affects 1x1 cell area
- **Radius-Based:** Effect affects circle of N cells radius
- **Pattern-Based:** Effect affects specific pattern (line, cone, cross)
- **Multi-Cell:** Effect can span multiple cells (3x3 poison zone)

**Target Selection:**
- **Self Only:** Effect only affects caster
- **Team Only:** Effect affects allies of caster
- **Enemy Only:** Effect affects enemies of caster
- **All Entities:** Effect affects everyone in zone
- **No Target:** Environmental effects affect the cell itself

**Priority and Stacking:**
- **Effect Age:** Newer effects override older same-type effects
- **Effect Intensity:** Stronger effects override weaker ones
- **Different Types:** Poison + Healing = both apply
- **Stack Limit:** Maximum N effects of same type per cell

**Movement Cost Effects:**
- **Additive Costs:** Multiple movement modifiers stack
- **Terrain Combination:** Cell terrain + movement effect = total cost
- **Minimum Cost:** Total movement cost never below 1
- **Display:** Show modified cost to player before movement

**Credit Assignment:**
```go
func AssignCellEffectCredits(effect CellEffect, affectedEntities []Entity) {
    credits := 0
    
    switch effect.EffectType {
        case "Zone":
            // Credits based on effect intensity
            for _, entity := range affectedEntities {
                credits += effect.Intensity.CalculateCredits(entity)
            }
            
        case "Trap":
            // Credits awarded on trigger
            credits += effect.Intensity.GetValue()
            
        case "Environmental":
            // Environmental effects don't earn credits
            credits = 0
    }
    
    // Assign to effect caster
    AssignCredits(effect.CasterID, credits)
}
```

**Integration with Multi-Entity Cells:**
- **Shared Cells:** Character + Zone Effect + Trap can coexist
- **Independent Processing:** Each effect processed separately per turn
- **Visual Layering:** Effects rendered under character layer
- **Collision Logic:** Effects define WalkThrough independently

**Visual Feedback:**
- **Zone Rendering:** Colored overlay showing effect area
- **Duration Display:** Countdown timer showing remaining turns
- **Trap Indication:** Hidden traps not shown, revealed after trigger
- **Effect Animation:** Particle effects for damage/healing/etc.

**Tactical Depth:**

**Zone Control:**
- **Area Denial:** Poison zones prevent safe positioning
- **Buff Zones:** Healing areas enable safe spots
- **Movement Control:** Quagmire slows enemy approaches
- **Zoning Strategy:** Create overlapping effects for complex battlefield control

**Trap Placement:**
- **Chokepoint Control:** Place traps in narrow passages
- **Predictive Positioning:** Anticipate enemy movement paths
- **Risk/Reward:** Traps cost resources but can be avoided
- **Information Control:** Revealed traps change enemy behavior

**Environmental Interaction:**
- **Terrain Modification:** Temporary walls/blocks alter line of sight
- **Weather Effects:** Global or local movement/visibility changes
- **Interactive Objects:** Destroyable cover, toggleable mechanics
- **Map Dynamics:** Static maps become dynamic battlefields

**Performance Considerations:**

**Effect Processing Optimization:**
```go
func ProcessAllCellEffects(grid Grid, entities []Entity) {
    // Batch process all effects
    for _, cell := range grid.GetAllCellsWithEffects() {
        // Cache entity positions for zone lookups
        entitiesInZone := FilterEntitiesInZone(cell, entities)
        
        // Process effect once per cell
        cell.Effect.ProcessTurn(entitiesInZone)
    }
}
```

**Player Experience:**

**Zone Combat:**
"I cast 'Poison Cloud' on the enemy team's approach path. Now they must choose: go through and take damage each turn, or find another route. Meanwhile, I stand safely behind my 'Healing Zone'. The grid handles both zones existing simultaneously, and I earn credits for poison damage dealt to multiple enemies!"

**Trap Strategy:**
"I notice the enemy always approaches through this narrow corridor, so I place a 'Bear Trap' here. The grid allows me to have multiple effects on cells—my character stands next to the trap while it's hidden. When the enemy walks on it, I earn 20 credits instantly, and the trap is removed. The grid seamlessly handles the cell becoming character-occupied after triggering."

**Implementation Priority:** HIGH - Required for Phase 2 zone-based combat and environmental interactions

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[cell_attached_effects]]`
- **Related Files:** `upsilonbattle/battlearena/grid/celleffects.go`, `upsilonbattle/battlearena/ruler/rules/movement.go`
- **Integration:** Works with `mechanic_mech_multi_entity_cell_system`, `expiration_controller`, `effect_caster_tracking`

## EXPECTATION
