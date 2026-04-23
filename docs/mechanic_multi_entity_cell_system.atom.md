---
id: mechanic_multi_entity_cell_system
status: DRAFT
version: 2.0
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
---

# New Atom

## INTENT
To implement multi-entity cell system where multiple entities (characters, effects, traps) can occupy the same grid cell simultaneously, enabling complex tactical interactions like zones and positional effects.

## THE RULE / LOGIC
**Multi-Entity Cell System:**

**Core Principle:**
Grid cells can contain multiple entities simultaneously, enabling characters to stand in zones, walk through traps, and interact with area effects.

**Single-Entity Limitation (V1):**
- **Exclusive Occupancy:** Only one entity per cell
- **Blocking:** Character on cell prevents movement onto cell
- **Problem:** Can't have characters + effects on same tile

**Multi-Entity Solution (V2):**

**Entity Co-Existence:**
- **Character + Zone Effects:** Characters can stand in healing/poison zones
- **Character + Traps:** Traps exist on same cell until triggered
- **Multiple Effects:** Different zone effects can stack on same cell
- **Temporary + Permanent:** Characters and temporary entities share cells

**Grid Data Structure (V2):**
```go
type GridCell struct {
    // Permanent entities (max 1)
    Character    *Entity       // At most 1 character per cell
    
    // Temporary entities (unlimited)
    Effects      []TemporaryEntity   // Multiple effects allowed
    Traps        []TemporaryEntity   // Multiple traps allowed
    
    // Cell properties
    Terrain       TerrainType
    WalkThrough   bool           // Can walk through effects
}
```

**Entity Collision Rules:**

**Character Movement:**
- **Onto Empty Cell:** Normal movement, no restrictions
- **Onto Character Cell:** Blocked - can't share cells with other characters
- **Onto Effect Cell:** Allowed - walk into zones, traps, temporary entities
- **WalkThrough Property:** Some effects allow movement (poison zones), others block (walls)

**Effect Placement:**
- **On Empty Cell:** Create effect normally
- **On Character Cell:** Effect co-exists with character
- **On Effect Cell:** Multiple effects can stack (poison + healing zone)
- **Priority System:** Some effects override others (newer effects take priority)

**Movement Cost Modifications:**
- **Quagmire Effects:** Add movement cost when walking through zone
- **Beneficial Zones:** Reduce movement cost or provide bonuses
- **Terrain Interaction:** Effects combine with terrain movement costs
- **Cumulative Effects:** Multiple movement modifiers stack

**Entity Interaction Scenarios:**

**Character in Healing Zone:**
```
Grid State:
Cell (5,3): Character + HealingZone Effect

Interaction:
- Character: Can move, act, attack normally
- HealingZone: Applies +10 HP to character each turn
- Duration: HealingZone has 3 turn lifetime
- Expiration: HealingZone disappears, character remains
```

**Character Steps on Trap:**
```
Grid State (Before):
Cell (5,3): Character (approaching)

Grid State (After):
Cell (5,3): Character + Trap Entity + Damage Event

Interaction:
- Character: Takes 15 damage, trap triggers
- Trap: Instantly expires after triggering
- Character: Remains on cell (now occupied by character only)
- Credits: Awarded to trap creator
```

**Multiple Zone Effects:**
```
Grid State:
Cell (5,3): Character + PoisonZone + HealingZone + SpeedBoost

Interaction:
- PoisonZone: Applies -5 HP each turn
- HealingZone: Applies +10 HP each turn
- SpeedBoost: Movement cost reduced by 1
- Character: Net +5 HP, faster movement
- Priority: Effects processed in specific order for balance
```

**WalkThrough Property:**

**Definition:**
- **True:** Character can move through cell containing effect
- **False:** Character cannot enter cell containing effect
- **Default:** Most effects allow walkthrough (zones, terrain mods)
- **Exceptions:** Some effects are impassible (temporary walls, force fields)

**Examples:**
- **Poison Zone (WalkThrough: true):** Characters can move through, take damage
- **Ice Wall (WalkThrough: false):** Characters cannot pass through
- **Healing Zone (WalkThrough: true):** Characters can enter for benefit
- **Quagmire (WalkThrough: true):** Can move through, +2 movement cost

**Collision Detection Logic:**
```go
func CanMoveTo(character Entity, target Position) error {
    cell := grid.GetCell(target)
    
    // Check character collision
    if cell.Character != nil && cell.Character.ID != character.ID {
        return errors.New("Cell occupied by another character")
    }
    
    // Check effect collision
    for _, effect := range cell.Effects {
        if !effect.WalkThrough {
            return errors.New("Cannot pass through impassible effect")
        }
    }
    
    // Calculate movement cost
    baseCost := 1 // Normal movement
    terrainCost := GetTerrainMovementCost(cell.Terrain)
    effectCost := 0
    
    for _, effect := range cell.Effects {
        effectCost += effect.MovementModifier
    }
    
    totalCost := baseCost + terrainCost + effectCost
    if character.Movement < totalCost {
        return errors.New("Insufficient movement points")
    }
    
    return nil
}
```

**Entity Processing Order:**

**Turn-Based Processing:**
```go
func ProcessEntitiesInCell(cell GridCell, currentTurn int) {
    // Process character actions
    if cell.Character != nil {
        cell.Character.ProcessTurn(currentTurn)
    }
    
    // Process all effects
    for i, effect := range cell.Effects {
        effect.ProcessTurn(currentTurn, cell.Character)
        
        // Remove expired effects
        if effect.IsExpired(currentTurn) {
            cell.Effects = append(cell.Effects[:i], cell.Effects[i+1:]...)
        }
    }
    
    // Process traps (trigger-based)
    for i, trap := range cell.Traps {
        if trap.ShouldBeTriggered(cell.Character) {
            trap.Trigger(cell.Character)
            cell.Traps = append(cell.Traps[:i], cell.Traps[i+1:]...)
        }
    }
}
```

**Visual Representation:**

**Multi-Layer Display:**
- **Character Layer:** Shows character sprite on top
- **Effect Layer:** Shows zones, auras underneath character
- **Trap Layer:** Shows traps (hidden/revealed) on cell
- **Transparency:** Effects use alpha blending to show overlapping

**UI Indicators:**
- **Cell Hover:** Show all entities on cell when hovering
- **Effect Duration:** Show remaining turns for all active effects
- **Movement Preview:** Highlight movement cost with effect modifiers
- **Danger Indicators:** Highlight traps/dangerous zones

**Tactical Depth:**

**Strategic Positioning:**
- **Zone Control:** Players create zones to deny or control areas
- **Trap Placement:** Strategic trap placement in chokepoints
- **Zone Stacking:** Combine beneficial + harmful effects
- **Risk/Reward:** Standing in poison zone but benefit from healing zone

**Character Builds:**
- **Zone Creators:** Skills that create multi-turn area effects
- **Trap Specialists:** Characters focusing on trap placement
- **Mobile Supports:** Characters who can move through zones quickly
- **Tank Builds:** Characters who benefit from stacking defensive effects

**Integration Benefits:**

**Temporary Entity System:**
- **Shared Cells:** Characters and temporary entities coexist naturally
- **Effect Processing:** All entities on cell processed each turn
- **Clean Coordination:** Expiration handles removing effects while characters remain

**Effect Caster Tracking:**
- **Multi-Target Credits:** Zone effects earn credits for multiple targets
- **Trap Credits:** Traps award credits to creator when triggered
- **Effect Ownership:** All effects know their caster even after caster death

**Performance Optimizations:**

**Entity Lookups:**
- **O(1) Cell Access:** Get cell position directly
- **O(N) Effect Processing:** Process all effects on cell linearly
- **Batch Operations:** Update multiple cells in single operation when possible

**Player Experience:**

**Zone Control Scenario:**
"I cast 'Poison Zone' on the enemy's approach path (3 turns). I also place a 'Healing Zone' on my position. Now the enemy must walk through poison to attack me, while I stay safe in my healing zone. The grid handles multiple effects on cells smoothly—I take healing while enemy takes poison damage."

**Trap Scenario:**
"I suspect the enemy will move through this corridor, so I place a 'Bear Trap' here. The grid allows me to stand on the same cell as my trap—it sits hidden until someone walks on it. When triggered, I'll earn credits even if I'm not nearby."

**Implementation Priority:** HIGH - Required for Phase 2 complex effects and tactical depth

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[multi_entity_cell_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/movement.go`, `upsilonbattle/battlearena/grid/grid.go`
- **Integration:** Works with `mec_expiration_controller`, `effect_caster_tracking`, `channeling_mechanic`

## EXPECTATION
