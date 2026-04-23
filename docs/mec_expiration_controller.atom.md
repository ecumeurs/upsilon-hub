---
id: mec_expiration_controller
human_name: Expiration Controller Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [controllers, entities, time-based]
parents:
  - [[mechanic_mech_temporary_entity_system]]
dependents: []
---

# Expiration Controller Mechanic

## INTENT
To implement the ExpirationController that manages the lifecycle of temporary entities, handling their death and cleanup when their duration expires or their effect completes.

## THE RULE / LOGIC
**Expiration Controller Purpose:**
- **Simple Cleanup:** Automatically kill temporary entities when their time comes
- **Effect Execution:** Trigger effects when entity's turn arrives
- **Duration Management:** Track and decrement entity lifespan
- **Self-Termination:** Execute effect, then die automatically

**Controller Behavior:**
```go
type ExpirationController struct {
    *controller.Controller
    TempEntity TemporaryEntity
}

func (ec *ExpirationController) ControllerNextTurn(ctx actor.NotificationContext) {
    // Execute the attached effect
    ec.ExecuteEffect()
    
    // Kill this temporary entity
    ec.Ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
        EntityID: ec.TempEntity.ID,
        ControllerID: ec.TempEntity.ControllerID,
        IsDeath: true,  // Signal this is entity death, not turn end
    }, nil))
}
```

**Area Effect Special Case:**
- **Multi-Turn Effects:** Execute effect each turn, then lose 1 HP
- **Duration Countdown:** Entity has HP counter that decrements each turn
- **Death When HP Reaches 0:** Entity dies naturally after N turns

**Channeling Special Case:**
- **One-Time Effect:** Execute effect once when channeling completes
- **Immediate Death:** Entity dies after effect execution
- **Caster Release:** Remove caster from IsCasting state

**Trap Special Case:**
- **OnStep Trigger:** Effect executes when stepped on, then dies
- **No Turn-Based Logic:** Uses movement trigger instead of turn trigger
- **One-Time Use:** Trap disappears after triggering

**Cleanup Process:**
```go
func (gs *GameState) KillEntity(ent entity.Entity) {
    // Execute OnDeath effects
    gs.OnDeath(ent)
    
    // Remove from grid
    gs.Grid.RemoveEntity(ent.Position)
    
    // Remove from turner
    gs.Turner.RemoveEntity(ent.ID)
    
    // Remove from game state
    delete(gs.Entities, ent.ID)
    
    // Release caster if applicable
    if ent.Type == Channeling {
        caster := gs.Entities[ent.CasterID]
        caster.UpdatePropertyValue(property.IsCasting, false)
        gs.Entities[caster.ID] = caster
    }
}
```

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_expiration_controller]]`
- **Related Files:** `upsilonbattle/battlearena/controller/controller.go`, `upsilonbattle/battlearena/ruler/rules/endofturn.go`
