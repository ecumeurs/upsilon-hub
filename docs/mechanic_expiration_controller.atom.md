---
id: mechanic_expiration_controller
status: DRAFT
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
version: 2.0
---

# New Atom

## INTENT
To implement expiration controller system that manages the lifecycle of temporary entities, cleaning up expired effects and removing entities when their duration ends naturally.

## THE RULE / LOGIC
**Expiration Controller System:**

**Core Purpose:**
Simple controller that terminates temporary entities when they expire, executing their final effects and cleaning up the grid.

**Temporary Entity Lifecycle:**

**1. Creation Phase:**
- **Source:** Skills, traps, zone effects, channeling entities
- **Duration Set:** Number of turns entity should exist
- **Spawned On:** Grid at specific position(s)
- **Controller Assigned:** Each temporary entity gets ExpirationController

**2. Active Phase:**
- **Turn-Based Processing:** Controller acts each turn to check entity age
- **Effect Application:** Some effects apply per turn (poison, healing zones)
- **Age Tracking:** Monitor how many turns entity has existed
- **Visual Feedback:** Show remaining duration to players

**3. Expiration Phase:**
- **Natural End:** Entity reaches its duration limit
- **Controller Trigger:** ExpirationController executes OnDeath trigger
- **Effect Execution:** Final effects fire (cleanup, final damage, etc.)
- **Grid Cleanup:** Entity removed from grid completely

**Expiration Controller Behavior:**

**Simple Lifecycle Management:**
```go
type ExpirationController struct {
    Controller
    EntityID       uuid.UUID    // Temporary entity being controlled
    CreationTurn   int          // When entity was created
    Duration        int          // How many turns to live
    FinalEffect     skill.Skill   // What to execute on expiration
}

func (ec *ExpirationController) Act(state *GameState) {
    // Check if entity has expired
    currentTurn := state.GetTurnIndex()
    age := currentTurn - ec.CreationTurn
    
    if age >= ec.Duration {
        // Entity expired - execute final effects
        ec.ExecuteFinalEffects()
        
        // Remove entity from grid
        state.RemoveEntity(ec.EntityID)
        
        // Cleanup controller
        ec.Terminate()
    }
}
```

**Expiration Trigger Types:**

**Natural Expiration:**
- **Duration Reached:** Entity lived its intended lifespan
- **Clean Execution:** Normal cleanup and final effect application
- **Credit Assignment:** Any final effects award credits to caster
- **Example:** Poison zone lasts 3 turns, heals everyone inside, then disappears

**Manual Cleanup:**
- **Player Action:** Some entities can be manually removed (disarming traps)
- **Immediate Expiration:** Trigger expiration immediately
- **Credit Loss:** May lose remaining potential credits from manual removal
- **Example:** Player steps on and disarms trap, trap expires instantly

**Destruction by Damage:**
- **Health-Based:** Temporary entities with HP die when HP reaches 0
- **Immediate Expiration:** OnDeath trigger instead of OnTurn trigger
- **Credit Assignment:** Any effects applied before death awarded normally
- **Example:** Channeling entity with 10 HP gets damaged for 5 HP (dies)

**Specific Temporary Entity Expiration:**

**Zone Effects (Healing/Poison):**
- **Per-Turn Processing:** Apply effect to all entities in zone each turn
- **Duration Tracking:** Count down turns remaining
- **Final Cleanup:** Remove zone, stop applying effects
- **Grid Interaction:** Zone doesn't block movement (WalkThrough property)

**Trap Entities:**
- **Trigger-Based:** Expiration only on OnStep trigger
- **Instant Effects:** Apply damage/effects immediately when triggered
- **Single-Use:** Traps expire immediately after triggering
- **Credit Assignment:** Credits awarded on trap trigger

**Channeling Entities:**
- **Delay-Based:** Wait for channeling duration before execution
- **Interruption Susceptible:** Can be destroyed/removed during channeling
- **OnTurn Trigger:** Execute channeled skill when delay completes
- **Failure State:** If interrupted, no execution, costs already consumed

**Area Effects (Fire/Ice):**
- **Damage Over Time:** Apply damage to entities in area each turn
- **Gradual Expiration:** Lose 1 HP per turn, die when HP reaches 0
- **OnDeath Trigger:** Final cleanup and removal from grid
- **Credit Tracking:** Damage credits awarded per turn to original caster

**Expiration Event Flow:**

```go
func HandleExpiration(entity TemporaryEntity) {
    // Determine expiration reason
    reason := entity.GetExpirationReason()
    
    switch reason {
        case "DurationReached":
            // Natural expiration
            entity.FinalEffect.OnDurationEnd()
            entity.Cleanup()
            
        case "Triggered":
            // Trap or channeling completion
            entity.FinalEffect.OnTrigger()
            entity.Cleanup()
            
        case "Destroyed":
            // HP-based entities killed
            entity.FinalEffect.OnDeath()
            entity.Cleanup()
    }
    
    // Assign any final credits
    credits := CalculateFinalCredits(entity)
    AssignCredits(entity.CasterID, credits)
}
```

**Integration with Other Systems:**

**Effect Caster Tracking:**
- **CasterID Preservation:** Expiration knows which player created entity
- **Credit Assignment:** Final effects award credits to original caster
- **Audit Trail:** All expirations logged with caster and target information

**Multi-Entity Cell System:**
- **Shared Expiration:** Multiple entities on same cell expire independently
- **Order Independence:** Expiration of one entity doesn't affect others
- **Grid Management:** Each expiration handled separately

**Visual Feedback System:**
- **Duration Display:** Show remaining turns to players
- **Expiration Warnings:** Highlight entities about to expire
- **Completion Indicators:** Visual cues when effects expire naturally

**Performance Considerations:**

**Batch Processing:**
- **Turn-Based Batching:** Check all expirations at once per turn
- **Efficient Cleanup:** Remove all expired entities in single grid operation
- **Memory Management:** Clean up controller references immediately

**Player Experience:**

**Zone Effect Scenario:**
"I cast 'Poison Cloud' - it appears and shows '3 turns remaining'. Next turn, it damages everyone inside (I earn credits). Turn after, shows '2 turns remaining'. Finally, it shows 'expiring...' and disappears. The system tracked the caster (me) and awarded credits for all poison damage dealt."

**Trap Scenario:**
"I place a 'Bear Trap' - it's invisible until someone steps on it. An enemy steps on it, the trap triggers instantly, dealing damage and expiring. I earn credits immediately. The trap disappears from the grid since it was single-use."

**Implementation Priority:** HIGH - Required for Phase 2 time-based mechanics cleanup

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[expiration_controller]]`
- **Related Files:** `upsilonbattle/battlearena/controller/controllers/expiration.go`, `upsilonbattle/battlearena/ruler/rules/beginingofturn.go`
- **Integration:** Works with `mechanic_mech_temporary_entity_system` and `effect_caster_tracking`

## EXPECTATION
