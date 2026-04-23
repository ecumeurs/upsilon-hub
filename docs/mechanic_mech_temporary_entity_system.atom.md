---
id: mechanic_mech_temporary_entity_system
status: DRAFT
parents: []
dependents:
  - [[mec_cell_attached_effects]]
  - [[mec_channeling_mechanic]]
  - [[mec_effect_caster_tracking]]
  - [[mec_expiration_controller]]
human_name: Temporary Entity System Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
priority: 5
tags: [time-based, entities, effects]
---

# New Atom

## INTENT
To implement the temporary entity system where time-based mechanics (channeling, traps, area effects) are represented as entities with controlled lifespans and trigger behaviors.

## THE RULE / LOGIC
The Temporary Entity System provides unified infrastructure for all time-based game mechanics:

**Core Principle:** 1 skill effect = 1 entity (simplified approach)

**Entity Types:**
- **TimeBased:** Channeling skills, delayed effects
- **Trap:** Triggers when stepped on
- **AreaEffect:** Affects multiple cells, expires after duration

**Temporary Entity Structure:**
```go
type TemporaryEntity struct {
    entity.Entity
    CasterID     uuid.UUID     // Who created this (for credits/interruption)
    TrueEffect   skill.Skill   // What to execute
    TriggerType  TriggerType   // When to execute
    Duration     int           // How many turns to live
}
```

**Trigger System:**
- **OnTurn:** Execute when entity's turn arrives (channeling, area effects)
- **OnStep:** Execute when entity stepped on (traps, quagmire)
- **OnDeath:** Execute when entity dies (explosions, cleanup)

**ExpirationController:** Simple controller that kills temporary entity when triggered, executing appropriate effects.

**Effect Caster Tracking:** All effects remember caster via CasterID until effect ends, enabling proper credit assignment and interruption mechanics.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_temporary_entity_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/beginingofturn.go`

## EXPECTATION
