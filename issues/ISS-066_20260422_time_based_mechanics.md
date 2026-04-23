# Issue: Time-Based Mechanics & Temporary Entity System

**ID:** `20260422_time_based_mechanics`
**Ref:** `ISS-066`
**Date:** 2026-04-22
**Severity:** High
**Status:** Open
**Component:** `upsilonbattle/battlearena/ruler`
**Affects:** `upsilonbattle/battlearena/entity`, `upsilonmapdata/grid`

---

## Summary

Implement time-based mechanics including channeling (pre-execution delay), temporary entities with expiration, effect caster tracking, and multi-entity cell support. This system enables delayed skill effects, traps, and environmental hazards.

---

## Technical Description

### Background

Current skill system executes effects immediately with simple delay costs. No support for channeling (delayed execution), temporary entities, or complex timing mechanics.

### The Problem Scenario

1. **Player wants to cast "Fireball"**: Should take 400 delay to cast, then effect happens
2. **No channeling system**: Effects happen immediately, no delay between cast and effect
3. **No temporary entities**: Cannot place traps, turrets, or environmental effects
4. **Single entity per cell**: Grid doesn't support multiple effects per location

### Time-Based Mechanics Architecture

**Cost Types:**
- **Pre-Execution Costs**: SP, MP, Channeling (delay before effect)
- **Post-Execution Costs**: Delay (delay after effect)

**Temporary Entity System:**
- 1 skill effect = 1 entity (simplified approach)
- Entity type: TimeBased, Trap, AreaEffect
- ExpirationController: Simple cleanup when entity dies
- Caster tracking: ControllerID field links to creator

**Trigger System:**
- **OnTurn**: Execute when entity's turn comes (channeling, area effects)
- **OnStep**: Execute when entity stepped on (traps, quagmire)
- **OnDeath**: Execute when entity dies (explosions, cleanup)

**Multi-Entity Cell System:**
- Grid cells support: 1 character + multiple effects
- WalkThrough property: Allows movement through non-blocking effects
- Cell-attached effects: Movement cost modifiers (quagmire = 2 movement per step)

**Effect Caster Tracking:**
```go
type Effect struct {
    Properties []property.Property
    Name       string
    CasterID   uuid.UUID  // Track creator for credits/interruption
    OriginTime time.Time  // When effect was applied
}
```

### Where This Pattern Exists Today

- `upsilonbattle/battlearena/ruler/rules/beginingofturn.go` - BeginOfTurn exists
- `upsilonbattle/battlearena/ruler/rules/endofturn.go` - EndOfTurn exists
- `upsilonbattle/battlearena/entity/entity.go` - Entity structure
- `upsilonmapdata/grid/grid.go` - Grid cell system

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Add OnDeath() trigger to GameState. Create TimeBased entity type. Implement ExpirationController.

**Medium term:** Update grid for multi-entity cells. Implement cell-attached effects. Add channeling logic to skill execution.

**Long term:** Build complex area effects and environmental hazards. Add interruption mechanics for channeling.

---

## References

- `V2_ARCHITECTURAL_DECISIONS.md` - Time-based mechanics architecture
- `upsilonbattle/battlearena/ruler/rules/beginingofturn.go` - Existing turn start logic
- `upsilonbattle/battlearena/ruler/rules/endofturn.go` - Existing turn end logic
