---
id: mec_cell_attached_effects
human_name: Cell Attached Effects Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [grid, effects, movement]
parents:
  - [[mechanic_mech_temporary_entity_system]]
dependents: []
---

# Cell Attached Effects Mechanic

## INTENT
To implement effects attached to grid cells that trigger on movement events (walking in/out) and modify movement costs, supporting environmental hazards like quagmire without creating multiple entities.

## THE RULE / LOGIC
**Cell Effect Attachment:**
- **Single Master Entity:** One entity controls the area effect
- **Effect Duration:** Entity loses 1 HP per turn, dies when HP reaches 0
- **OnTurn Trigger:** Entity applies effect to all entities in zone each turn

**Movement Trigger Types:**
- **OnStepIn:** Trigger when entity enters cell
- **OnStepOut:** Trigger when entity leaves cell
- **OnStepBoth:** Trigger on both entry and exit

**Movement Cost Modifiers:**
- **Example:** Quagmire increases movement cost from 1 to 2 per step
- **Implementation:** Movement cost checked during pathfinding and execution
- **Duration:** Persists as long as entity remains in affected cells

**Master Entity Approach:**
```go
// Poisonous Fog Area Effect
poisonFog := TemporaryEntity{
    Entity: entity.Entity{
        Position: {5, 5, 1},  // Center of fog
        Properties: {
            HP: IntCounter(5, 5),   // 5 turn duration
            WalkThrough: true,
        },
    },
    TrueEffect: skill.Skill{
        Zone: Cross(3),           // 3-cell cross pattern
        PoisonPower: 3,          // 3 poison damage
    },
    TriggerType: OnTurn,       // Execute every turn
}

// OnTurn: Apply poison to all entities in zone, lose 1 HP
```

**Cell vs Entity Effects:**
- **Use Entity-per-Cell:** Unique behavior per cell, individual timing
- **Use Master Entity:** Uniform effect across area, single timing
- **Simplification:** 1 skill effect = 1 entity (master entity for areas)

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_cell_attached_effects]]`
- **Related Files:** `upsilonmapdata/grid/grid.go`, `upsilonbattle/battlearena/ruler/rules/move.go`
