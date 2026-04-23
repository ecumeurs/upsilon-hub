---
id: mec_multi_entity_cell_system
human_name: Multi-Entity Cell System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [grid, entities, movement]
parents: []
dependents: []
---

# Multi-Entity Cell System Mechanic

## INTENT
To enable multiple entities to occupy the same grid cell, supporting character + effects co-location while maintaining character-vs-character collision rules.

## THE RULE / LOGIC
**Cell Structure Update:**
- **Character Slot:** Only ONE character per cell
- **Effect Slots:** Multiple effects allowed per cell
- **WalkThrough Property:** Determines if entity blocks movement

**Multi-Entity Cell Rules:**
- **Character Collision:** Cannot move into cell with another character
- **Effect Co-location:** Multiple effects can occupy same cell
- **Character + Effects:** Character can share cell with multiple effects
- **Priority:** Character collision takes precedence over effect collision

**Movement Logic:**
```go
func CanMoveTo(cell Cell, entityID uuid.UUID) bool {
    // Check character collision
    if cell.CharacterID != uuid.Nil && cell.CharacterID != entityID {
        return false  // Character collision
    }
    
    // Check effect blocking
    for _, effectID := range cell.EffectIDs {
        effect := GetEntity(effectID)
        if !effect.WalkThrough {
            return false  // Effect blocks movement
        }
    }
    
    return true  // Movement allowed
}
```

**WalkThrough Property:**
- **True:** Entity does not block movement (channeling, traps, beneficial effects)
- **False:** Entity blocks movement (walls, barriers, obstacles)
- **Default:** Temporary effects default to WalkThrough = true

**Entity Type Handling:**
- **Characters:** Always block other characters, never WalkThrough
- **TimeBased Effects:** Usually WalkThrough = true
- **Traps:** WalkThrough = true (to trigger on step)
- **Barriers:** WalkThrough = false (block movement)

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_multi_entity_cell_system]]`
- **Related Files:** `upsilonmapdata/grid/grid.go`, `upsilonbattle/battlearena/ruler/rules/move.go`
