---
id: mec_effect_caster_tracking
human_name: Effect Caster Tracking Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [effects, tracking, credits]
parents:
  - [[mechanic_mech_temporary_entity_system]]
dependents: []
---

# Effect Caster Tracking Mechanic

## INTENT
To implement effect caster tracking where all effects remember their originator (caster) until effect ends, enabling proper credit assignment and interruption mechanics even after caster death.

## THE RULE / LOGIC
**Effect Caster Tracking:**
- **All Effects Track Originator:** Every effect must remember who created it
- **CasterID Property:** Stored in effect structure, linked to creating entity
- **OriginTime Timestamp:** When effect was applied for audit and duration tracking
- **Persistent Tracking:** Credits go to original caster even if they die

**Effect Structure Enhancement:**
```go
type Effect struct {
    Properties   []property.Property
    Name         string
    CasterID     uuid.UUID  // Track creator for credit assignment
    OriginTime   time.Time  // When effect was applied
    EffectType   EffectType // Shield, Poison, Stun, Buff, etc.
}
```

**Credit Assignment Logic:**
```go
// When shield blocks damage
if shield.BlocksDamage > 0 {
    credits.Earned += shield.BlocksDamage
    credits.AssignedTo = shield.CasterID  // Original caster gets credits
}

// When effect expires naturally
if effect.OriginTime + effect.Duration < currentTime {
    effect.CasterID receives any completion bonuses
}
```

**Interruption Considerations:**
- **Caster Death:** Effect continues earning credits for original caster
- **Effect Dispel:** Credits stop accruing but earned credits remain with caster
- **Effect Expiration:** Final credit结算 based on total performance

**Temporary Entity Caster Tracking:**
- **Channeling Entities:** Track who is channeling for interruption purposes
- **Trap Entities:** Track who placed the trap for credit assignment
- **Area Effects:** Track who created the effect for damage/healing credits

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_effect_caster_tracking]]`
- **Related Files:** `upsilonbattle/battlearena/property/effect/effect.go`, `upsilonbattle/battlearena/entity/entity.go`
