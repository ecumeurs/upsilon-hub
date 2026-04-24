# Preparation Document: ISS-066 Time-Based Mechanics & Temporary Entity System

**ID:** `ISS-066`
**Date:** 2026-04-23
**Status:** Preparation
**Component:** `upsilonbattle/battlearena/ruler`, `upsilonbattle/battlearena/controller`, `upsilonmapdata/grid`

---

## EXECUTIVE SUMMARY

This document outlines the complete design for implementing time-based mechanics including:
- Turner-based entities (turrets, delayed effects, walls)
- Multi-entity cell system
- Cell-attached effects with triggers
- Behavior-based AI system
- Composite controllers (OR/AND modes)

---

## PHASE 1: TURNER-BASED MECHANICS

### 1.1 New Entity Types

Extend `EntityType` enum in `entity.go`:

```go
type EntityType int

const (
    Character  EntityType = 0
    Monster    EntityType = 1
    TimeBased  EntityType = 2  // Channeling entities, delayed effects
    Trap       EntityType = 3  // Triggers on OnStep
    AreaEffect EntityType = 4  // Affects zone each turn
    Obstacle   EntityType = 5  // Walls, barriers (don't act, block movement)
    Others     EntityType = 6
)
```

### 1.2 New Properties

Add to `property.EntityProperties` in `propertyenum.go`:

```go
// Entity lifecycle
Duration       EntityProperties = "Duration"        // How many turns to live
ExpiresWithCaster EntityProperties = "ExpiresWithCaster" // Remove when caster dies
WalkThrough     EntityProperties = "WalkThrough"   // Can walk through this entity?
Invisible       EntityProperties = "Invisible"     // Not visible to clients
```

Add to `property.SkillProperties` for triggers:

```go
// Effect triggers
TriggerType    SkillProperties = "TriggerType"     // When effect fires
RemoveOnTrigger SkillProperties = "RemoveOnTrigger" // Effect dies after executing
TriggerCount    SkillProperties = "TriggerCount"    // How many times can trigger (0 = unlimited)
```

TriggerType values:
```go
type TriggerType string

const (
    TriggerOnEnter TriggerType = "OnEnter"  // When entity enters cell
    TriggerOnExit  TriggerType = "OnExit"   // When entity leaves cell
    TriggerOnStep  TriggerType = "OnStep"   // Each step through cell
    TriggerOnTurn  TriggerType = "OnTurn"   // Each turn while in cell
    TriggerOnDeath TriggerType = "OnDeath"  // When entity dies
)
```

### 1.3 Turret Entities

**Definition**: Entity that acts on own turn, may have limited lifespan, uses ranged attacks.

**Implementation**:
- Use existing `AggressiveController` with `Movement = 0`
- `AttackRange` property determines range
- Has `Duration` property for limited lifespan
- Has `ControllerID` (assigned controller makes decisions)

**Flow**:
1. Create entity with Movement=0, AttackRange=N
2. Add to Turner with CurrentDelay=0
3. On turn, AggressiveController finds nearest foe in range
4. Attacks, adds Delay, requeues
5. Each EndOfTurn, check Duration, remove if 0

### 1.4 Delayed/Channeling Effects

**Definition**: Entity that waits for delay, then executes effect once and dies.

**Implementation**:
- Create entity with `EntityType = TimeBased`
- Store the `effect.Effect` to execute when turn arrives
- Add to Turner with `CurrentDelay = channeling_value`
- Has `Duration = 1` (executes once)

**Flow**:
1. Player casts channeling skill (Fireball with Channeling 400)
2. Create TimeBased entity at caster position
3. Add to Turner with CurrentDelay=400
4. Caster pays pre-execution costs (SP, MP)
5. When entity's turn arrives:
   - Execute stored effect
   - Apply post-execution Delay to caster
   - Remove entity

**Interruption** (future):
- If caster takes damage while channeling, track Interruption
- At 100, cancel channeling, entity dies, no effect executes

### 1.5 Constructed Obstacles (Walls)

**Definition**: Entity that blocks movement, has HP, may degrade over time, cannot act.

**Implementation**:
- `EntityType = Obstacle`
- `ControllerID = uuid.Nil` OR has controller with only ExpirationBehavior
- `WalkThrough = false` (blocks movement)
- Has `HP`, `Duration` (for decay)
- Pathfinding treats as obstacle

**Behavior**:
- Gets turns in Turner but does nothing (or passes immediately)
- Could use `ExpirationBehavior` to handle decay
- When HP=0 or Duration=0, removed

**Flow**:
1. Player casts "Wall" skill
2. Create Obstacle entity at target position
3. Add to Turner (optional) or don't add if no controller
4. Blocks movement (WalkThrough=false)
5. Takes damage, decays, eventually removed

### 1.6 Expiration Logic

**Location**: `EndOfTurn()` in `gamestate.go`

```go
func (gs *GameState) EndOfTurn(msg *message.Message, req rulermethods.EndOfTurn, ent entity.Entity) (ok bool, reply *message.Message) {
    // ... existing logic ...

    // Check duration for temporary entities
    duration := ent.GetPropertyI(property.Duration).I()
    if duration > 0 {
        newDuration := duration - 1
        ent.UpdatePropertyValue(property.Duration, newDuration)

        if newDuration <= 0 {
            // Entity expired
            gs.RemoveEntity(ent.ID)
            gs.Logger.WithFields(logrus.Fields{
                "entityID": ent.ID.String()[0:8],
                "reason": "duration_expired",
            }).Info("Entity expired")
            return true, msg
        }
    }

    // ... rest of existing logic ...
}
```

**RemoveEntity** method to add to GameState:

```go
func (gs *GameState) RemoveEntity(entityID uuid.UUID) {
    // Remove from Turner
    gs.Turner.RemoveEntity(entityID)

    // Remove from Grid
    if ent, exists := gs.Entities[entityID]; exists {
        gs.Grid.RemoveEntity(ent.Position)
    }

    // Remove from Entities map
    delete(gs.Entities, entityID)

    // Cleanup PositionalEffects owned by this entity
    for pos, effectIDs := range gs.PositionalEffects {
        newEffectIDs := []uuid.UUID{}
        for _, effectID := range effectIDs {
            if effect, exists := gs.Effects[effectID]; exists {
                if effect.CasterID != entityID {
                    newEffectIDs = append(newEffectIDs, effectID)
                } else {
                    // Remove the effect itself
                    delete(gs.Effects, effectID)
                }
            }
        }
        if len(newEffectIDs) == 0 {
            delete(gs.PositionalEffects, pos)
        } else {
            gs.PositionalEffects[pos] = newEffectIDs
        }
    }
}
```

---

## PHASE 2: MULTI-ENTITY CELL SYSTEM

### 2.1 Cell Structure Changes

**File**: `upsilonmapdata/grid/cell/cell.go`

```go
type Cell struct {
    Type              CellType
    Position          position.Position

    // Entities in this cell (can be multiple)
    EntityIDs         []uuid.UUID

    // Effects attached to this cell
    EffectIDs         []uuid.UUID

    // Terrain properties (future use)
    CrossingCost      int  // Default 1, higher for difficult terrain
}
```

**Note**: EntityID changed to EntityIDs (slice).

### 2.2 Grid Changes

**File**: `upsilonmapdata/grid/grid.go`

Update methods to handle multiple entities:

```go
// MoveEntity moves an entity from one position to another
func (g *Grid) MoveEntity(from, to position.Position, entityID uuid.UUID) error {
    if !g.Contains(to) {
        return fmt.Errorf("to position %v is not in the grid", to)
    }

    // Remove from old cell
    if c, ok := g.CellAt(from); ok {
        c.EntityIDs = removeID(c.EntityIDs, entityID)
    }

    // Add to new cell
    c, ok := g.CellAt(to)
    if !ok {
        return fmt.Errorf("to position %v is not in the grid", to)
    }
    c.EntityIDs = append(c.EntityIDs, entityID)

    return nil
}

// RemoveEntity removes an entity from a cell
func (g *Grid) RemoveEntity(p position.Position, entityID uuid.UUID) {
    if !g.Contains(p) {
        return
    }
    g.Cells[p].EntityIDs = removeID(g.Cells[p].EntityIDs, entityID)
}

// GetEntitiesAt returns all entities at a position
func (g *Grid) GetEntitiesAt(p position.Position) []uuid.UUID {
    if c, ok := g.CellAt(p); ok {
        return c.EntityIDs
    }
    return nil
}

// CanMoveTo checks if movement is possible (without checking effects)
func (g *Grid) CanMoveTo(from, to position.Position, jumpHeight int, entityID uuid.UUID) (bool, error) {
    if !g.Contains(to) {
        return false, fmt.Errorf("target not in grid")
    }

    cell, ok := g.CellAt(to)
    if !ok {
        return false, fmt.Errorf("cell not found")
    }

    // Check terrain
    if cell.Type != cell.Ground {
        return false, fmt.Errorf("cell is not walkable terrain")
    }

    // Check for blocking entities
    for _, entID := range cell.EntityIDs {
        if entID == entityID {
            continue  // Self
        }

        // Need GameState to check WalkThrough property
        // This check will be done in Move rule, not here
        // Grid doesn't have access to entity data
    }

    // Check terrain cost (movement validation happens in Move rule)
    if cell.CrossingCost > 10 {  // Example threshold
        return false, fmt.Errorf("terrain impassable")
    }

    return true, nil
}

// Helper function
func removeID(ids []uuid.UUID, id uuid.UUID) []uuid.UUID {
    for i, existingID := range ids {
        if existingID == id {
            return append(ids[:i], ids[i+1:]...)
        }
    }
    return ids
}
```

### 2.3 Movement Cost Calculation

Movement cost calculation happens in the `Move` rule, combining:
1. Base cost: 1
2. Cell crossing cost: `cell.CrossingCost`
3. Effect movement costs: summed from all effects on the path

---

## PHASE 2: POSITIONAL EFFECTS SYSTEM

### 2.1 Data Structures

**File**: `upsilonbattle/battlearena/ruler/rules/gamestate.go`

Add to `GameState`:

```go
type GameState struct {
    RulerID     uuid.UUID
    Grid        *grid.Grid
    Turner      turner.Turner
    Entities    map[uuid.UUID]entity.Entity
    Controllers map[uuid.UUID]actor.Communication
    Logger      *logrus.Entry
    WinnerTeamID int
    Version     int64
    TurnIndex   uint32
    ActionIndex uint32

    // Positional effects (cell-attached effects)
    PositionalEffects map[position.Position][]uuid.UUID
    Effects           map[uuid.UUID]effect.Effect  // All effects by ID
}
```

**Note**: Cell stores EffectIDs, GameState stores the actual Effect data.

### 2.2 Effect Ownership & Cleanup

**Effect Structure** (already exists in `effect/effect.go`):

```go
type Effect struct {
    Properties []property.Property
    Name       string
    CasterID   uuid.UUID  // Who created this
}
```

**New Effect Properties** (add to `property.SkillProperties`):

```go
TriggerType     SkillProperties = "TriggerType"      // OnEnter, OnExit, OnStep, OnTurn, OnDeath
RemoveOnTrigger SkillProperties = "RemoveOnTrigger" // bool: remove after executing
TriggerCount    SkillProperties = "TriggerCount"     // int: how many times can trigger (0 = unlimited)
```

**ExpiresWithCaster Property** (add to `property.EntityProperties`):

```go
ExpiresWithCaster EntityProperties = "ExpiresWithCaster" // bool
```

**Cleanup on Entity Death**:

When an entity dies (in `OnDeath` or `RemoveEntity`):
1. Iterate all `PositionalEffects`
2. For each effect, check if `effect.CasterID == dyingEntityID`
3. If yes, check `ExpiresWithCaster` property
4. If true, remove effect from both `PositionalEffects` and `Effects` map
5. Remove effect ID from cells

### 2.3 Zone Entity Pattern

**Concept**: A zone (like poisonous fog) is created by an invisible anchor entity.

**Flow**:
1. Create invisible `TimeBased` entity at zone center
2. Entity has `Duration` and `Invisible = true`
3. Entity has no controller (or ExpirationBehavior only)
4. Calculate zone pattern (e.g., Cross(3) for 3x3 area)
5. For each cell in zone:
   - Create `effect.Effect` with appropriate properties
   - Set `effect.CasterID = anchorEntity.ID`
   - Add to `Effects` map
   - Add effect ID to cell's `EffectIDs`
6. When anchor entity dies (Duration=0):
   - RemoveEntity cleanup finds all effects with CasterID = anchor.ID
   - All zone effects are removed

**Example - Poisonous Fog**:

```go
func CreatePoisonousFog(gs *GameState, caster entity.Entity, center position.Position, duration int) {
    // Create invisible anchor entity
    anchor := entity.Entity{
        ID:           uuid.New(),
        ControllerID: uuid.Nil,  // No controller
        Type:         TimeBased,
        Position:     center,
        Name:         "Poison Fog Anchor",
        Properties: map[string]property.Property{
            property.Duration: def.MakeIntCounterProperty(property.Duration, duration, duration, ...),
            property.Invisible: def.MakeBoolProperty(property.Invisible, true, ...),
        },
    }
    gs.Entities[anchor.ID] = anchor
    gs.Turner.AddEntity(anchor.ID, 100)  // Gets turns to check duration

    // Create zone
    zone := pattern.Cross(3)
    positions := gs.Grid.SelectPositionsByPattern(center, zone)

    for _, pos := range positions {
        // Create effect
        fogEffect := effect.Effect{
            Properties: []property.Property{
                def.MakeIntProperty(property.PoisonPower, 3, ...),
                def.MakeIntProperty(property.TriggerType, string(TriggerOnEnter), ...),
                def.MakeBoolProperty(property.RemoveOnTrigger, false, ...),  // Don't remove, persists
                def.MakeBoolProperty(property.ExpiresWithCaster, true, ...),  // Die with anchor
            },
            Name:     "Poison Fog",
            CasterID: anchor.ID,  // Anchor is the "caster"
        }

        effectID := uuid.New()
        gs.Effects[effectID] = fogEffect

        // Attach to cell
        if cell, ok := gs.Grid.CellAt(pos); ok {
            cell.EffectIDs = append(cell.EffectIDs, effectID)
            gs.PositionalEffects[pos] = cell.EffectIDs
        }
    }
}
```

### 2.4 Trigger System

**Trigger Types**:
- `OnEnter`: When entity enters cell
- `OnExit`: When entity leaves cell
- `OnStep`: Each step through cell (enter + move within)
- `OnTurn`: Apply effect each turn while entity is in cell
- `OnDeath`: Apply effect when entity dies in/on this cell

**Trigger Implementation**: In `Move` rule and `BeginingOfTurn`/`EndOfTurn`.

**Movement Triggers** (OnEnter, OnExit, OnStep):

File: `upsilonbattle/battlearena/ruler/rules/move.go`

```go
func (gs *GameState) Move(msg *message.Message, req rulermethods.ControllerMove) (reply *message.Message) {
    // ... pre-checks ...

    ent := gs.Entities[req.EntityID]

    // Track cells we enter/exit/step through
    visitedCells := req.Path  // All cells we'll move through

    // Process movement
    for i, pos := range visitedCells {
        cell, ok := gs.Grid.CellAt(pos)
        if !ok {
            return msg.ReplyWithError("Cell not found", "cell.notfound")
        }

        // Get all effects on this cell
        for _, effectID := range cell.EffectIDs {
            effect, exists := gs.Effects[effectID]
            if !exists {
                continue
            }

            triggerType := getTriggerType(effect)
            shouldTrigger := false

            // Check trigger conditions
            if i == len(visitedCells) - 1 {
                // Final position - entering this cell
                if triggerType == TriggerOnEnter || triggerType == TriggerOnStep {
                    shouldTrigger = true
                }
            } else {
                // Passing through - OnStep triggers
                if triggerType == TriggerOnStep {
                    shouldTrigger = true
                }
            }

            if shouldTrigger {
                // Apply effect using effectapplicator
                affected, damaged, credits, _, _ := effectapplicator.ApplyDirectEffect(
                    gs.Logger, &ent, effect, pos, []position.Position{pos},
                    gs.Grid, []entity.Entity{ent},
                )

                // Update entity
                if len(affected) > 0 {
                    gs.Entities[affected[0].ID] = affected[0]
                    ent = affected[0]
                }

                // Check for movement cost increase
                if effect.HasProperty(property.MvtCost) {
                    mvtCost := effect.GetPropertyI(property.MvtCost).I()
                    // This will be applied to movement consumption
                    // If movement is less than cost, entity loses all remaining movement
                    // But can still ENTER the cell (doesn't block)
                }

                // Check for force stop
                if effect.HasProperty(property.ForceStopMove) {
                    ent.UpdatePropertyValue(property.HasMoved, true)
                }
                if effect.HasProperty(property.ForceEndTurn) {
                    ent.UpdatePropertyValue(property.HasActed, true)
                    // End turn immediately
                    gs.Turner.AddEntity(ent.ID, ent.CurrentDelay + 300)
                    return msg.Reply()
                }

                // Check if effect should be removed
                if getRemoveOnTrigger(effect) {
                    gs.RemovePositionalEffect(effectID, pos)
                } else {
                    // Decrement trigger count
                    triggerCount := getTriggerCount(effect)
                    if triggerCount > 0 {
                        setTriggerCount(effect, triggerCount - 1)
                        if getTriggerCount(effect) == 0 {
                            gs.RemovePositionalEffect(effectID, pos)
                        }
                    }
                }
            }
        }
    }

    // Complete the move
    gs.Grid.MoveEntity(ent.Position, req.Path[len(req.Path)-1], ent.ID)
    ent.Position = req.Path[len(req.Path)-1]

    // Calculate total movement cost (1 + cell crossing costs + effect costs)
    totalCost := len(req.Path)
    for _, pos := range req.Path {
        if cell, ok := gs.Grid.CellAt(pos); ok {
            totalCost += cell.CrossingCost
        }
    }

    // Apply movement cost
    prop := ent.GetPropertyC(property.Movement)
    actualCost := min(prop.GetValue(), totalCost)
    prop.SetValue(prop.GetValue() - actualCost)
    ent.UpdateProperty(prop)

    // Add delay
    ent.CurrentDelay += len(req.Path) * 20

    gs.Entities[req.EntityID] = ent
    gs.IncVersion()

    // Notify controllers
    // ...

    return msg.Reply()
}
```

**Turn Triggers** (OnTurn):

File: `upsilonbattle/battlearena/ruler/rules/beginingofturn.go`

```go
func (gs *GameState) BeginingOfTurn(ent entity.Entity) {
    // ... existing stun logic ...

    // Check for OnTurn triggers at current position
    cell, ok := gs.Grid.CellAt(ent.Position)
    if !ok {
        return
    }

    for _, effectID := range cell.EffectIDs {
        effect, exists := gs.Effects[effectID]
        if !exists {
            continue
        }

        if getTriggerType(effect) == TriggerOnTurn {
            // Apply effect
            affected, damaged, credits, _, _ := effectapplicator.ApplyDirectEffect(
                gs.Logger, &ent, effect, ent.Position,
                []position.Position{ent.Position}, gs.Grid, []entity.Entity{ent},
            )

            if len(affected) > 0 {
                gs.Entities[affected[0].ID] = affected[0]
                ent = affected[0]
            }

            // Handle removal/trigger count
            // ... same as in Move ...
        }
    }

    gs.Entities[ent.ID] = ent
}
```

---

## PHASE 3: BEHAVIOR-BASED AI SYSTEM

### 3.1 Separation of Concerns

**Current**: Controller = Ownership + Communication + AI Behavior

**New**: Controller = Ownership + Communication + Behavior (pluggable)

### 3.2 Behavior Interface

**File**: `upsilonbattle/battlearena/controller/behavior/behavior.go` (new)

```go
package behavior

import (
    "github.com/ecumeurs/upsilonbattle/battlearena/entity"
    "github.com/google/uuid"
)

type GameContext interface {
    GetEntity(id uuid.UUID) (entity.Entity, bool)
    GetEntities() map[uuid.UUID]entity.Entity
    GetGrid() Grid
    GetCurrentTurn() int
}

type DecisionType string

const (
    NoDecision DecisionType = "NoDecision"  // Behavior has no opinion
    Move       DecisionType = "Move"
    Attack     DecisionType = "Attack"
    Skill      DecisionType = "Skill"
    Pass       DecisionType = "Pass"
    Flee       DecisionType = "Flee"
)

type Decision struct {
    Type      DecisionType
    Target    position.Position
    SkillID   uuid.UUID
    Path      []position.Position
    Priority  int  // Higher = more important (for AND mode)
}

type Behavior interface {
    // OnTurn is called when entity's turn begins
    OnTurn(ctx GameContext, ent entity.Entity) Decision

    // OnDamaged is called when entity takes damage
    OnDamaged(ctx GameContext, ent entity.Entity, damage int, attacker entity.Entity)

    // GetName returns the behavior name for debugging
    GetName() string
}
```

### 3.3 Concrete Behaviors

**AggressiveBehavior**:

```go
type AggressiveBehavior struct {
    Name string
}

func (ab *AggressiveBehavior) OnTurn(ctx GameContext, ent entity.Entity) Decision {
    // Find nearest foe
    target, err := findNearestFoe(ctx, ent)
    if err != nil {
        return Decision{Type: Pass}
    }

    // Check if in attack range
    attackRange := ent.GetPropertyI(property.AttackRange).I()
    if ent.Position.Distance(target.Position) <= attackRange {
        return Decision{
            Type:   Attack,
            Target: target.Position,
        }
    }

    // Try to move toward target
    path := findPath(ctx, ent.Position, target.Position, ent.GetPropertyI(property.JumpHeight).I())
    if len(path) > 1 {
        movement := ent.GetPropertyI(property.Movement).I()
        if movement > 0 {
            // Move as far as we can
            limit := min(movement, len(path)-1)
            return Decision{
                Type: Move,
                Path: path[1:limit+1],  // Exclude starting position
            }
        }
    }

    return Decision{Type: Pass}
}

func (ab *AggressiveBehavior) OnDamaged(ctx GameContext, ent entity.Entity, damage int, attacker entity.Entity) {
    // Aggressive entities might get angry and focus on attacker
    // Store attacker preference for future turns
}

func (ab *AggressiveBehavior) GetName() string {
    return ab.Name
}
```

**ExpirationBehavior**:

```go
type ExpirationBehavior struct {
    Name string
}

func (eb *ExpirationBehavior) OnTurn(ctx GameContext, ent entity.Entity) Decision {
    // Check if entity has expired
    duration := ent.GetPropertyI(property.Duration).I()
    if duration <= 0 {
        // Entity will be removed by EndOfTurn
        return Decision{Type: Pass}
    }

    // Decrement duration handled by EndOfTurn
    // This behavior doesn't make decisions, just participates in AND composites
    return Decision{Type: NoDecision}
}

func (eb *ExpirationBehavior) OnDamaged(ctx GameContext, ent entity.Entity, damage int, attacker entity.Entity) {
    // Nothing to do
}

func (eb *ExpirationBehavior) GetName() string {
    return eb.Name
}
```

**SpookedBehavior**:

```go
type SpookedBehavior struct {
    Name          string
    ThresholdHP   int
    FleeDistance  int
    LastDamage    int
    LastDamageTurn int
}

func (sb *SpookedBehavior) OnTurn(ctx GameContext, ent entity.Entity) Decision {
    currentHP := ent.GetPropertyI(property.HP).I()
    currentTurn := ctx.GetCurrentTurn()

    // Check if should flee
    shouldFlee := false

    if currentHP < sb.ThresholdHP {
        shouldFlee = true
    }

    if sb.LastDamageTurn == currentTurn-1 && sb.LastDamage > 5 {
        shouldFlee = true
    }

    if shouldFlee {
        // Find safest position (furthest from all enemies)
        safestPos := findSafestPosition(ctx, ent, sb.FleeDistance)

        path := findPath(ctx, ent.Position, safestPos, ent.GetPropertyI(property.JumpHeight).I())
        if len(path) > 1 {
            movement := ent.GetPropertyI(property.Movement).I()
            limit := min(movement, len(path)-1)
            return Decision{
                Type: Flee,
                Path: path[1:limit+1],
            }
        }
    }

    return Decision{Type: NoDecision}
}

func (sb *SpookedBehavior) OnDamaged(ctx GameContext, ent entity.Entity, damage int, attacker entity.Entity) {
    sb.LastDamage = damage
    sb.LastDamageTurn = ctx.GetCurrentTurn()
}

func (sb *SpookedBehavior) GetName() string {
    return sb.Name
}
```

### 3.4 Composite Behaviors

**CompositeMode**:

```go
type CompositeMode string

const (
    CompositeOR  CompositeMode = "OR"  // First decision wins
    CompositeAND CompositeMode = "AND" // All execute, later overrides
)
```

**CompositeBehavior**:

```go
type CompositeBehavior struct {
    Name       string
    Mode       CompositeMode
    Behaviors  []BehaviorWrapper
}

type BehaviorWrapper struct {
    Behavior Behavior
    Priority int  // Higher priority runs first in OR mode
}

func (cb *CompositeBehavior) OnTurn(ctx GameContext, ent entity.Entity) Decision {
    if cb.Mode == CompositeOR {
        // Try each behavior in priority order
        for _, wrapper := range cb.Behaviors {
            decision := wrapper.Behavior.OnTurn(ctx, ent)
            if decision.Type != NoDecision {
                return decision  // First decision wins
            }
        }
        return Decision{Type: Pass}  // No behavior had an opinion
    } else {
        // AND mode: all behaviors execute, later decisions override
        final := Decision{Type: NoDecision}

        // Sort by priority (highest first)
        sorted := make([]BehaviorWrapper, len(cb.Behaviors))
        copy(sorted, cb.Behaviors)
        sort.Slice(sorted, func(i, j int) bool {
            return sorted[i].Priority > sorted[j].Priority
        })

        for _, wrapper := range sorted {
            decision := wrapper.Behavior.OnTurn(ctx, ent)
            if decision.Type != NoDecision {
                // Override with this decision
                final = decision
            }
        }

        return final
    }
}

func (cb *CompositeBehavior) OnDamaged(ctx GameContext, ent entity.Entity, damage int, attacker entity.Entity) {
    // Notify all behaviors
    for _, wrapper := range cb.Behaviors {
        wrapper.Behavior.OnDamaged(ctx, ent, damage, attacker)
    }
}

func (cb *CompositeBehavior) GetName() string {
    return cb.Name
}
```

### 3.5 Controller Refactoring

**New Controller Structure**:

File: `upsilonbattle/battlearena/controller/controller.go`

```go
type Controller struct {
    *actor.Actor
    ID             uuid.UUID
    Assigned       bool
    ControllerName string
    Ruler          actor.Communication
    Behavior       behavior.Behavior  // NEW: Pluggable behavior
    EntityID       uuid.UUID          // NEW: Which entity we control
}
```

**ControllerNextTurn with Behavior**:

```go
func (c *Controller) ControllerNextTurn(ctx actor.NotificationContext) {
    req := ctx.Msg.TargetMethod.(rulermethods.ControllerNextTurn)

    // Only respond if this is our entity
    if req.Entity.ControllerID != c.ID {
        return
    }

    // Get decision from behavior
    gameContext := NewGameContext(c.KnownEntities, c.Grid, c.TurnIndex)
    decision := c.Behavior.OnTurn(gameContext, req.Entity)

    // Execute decision
    switch decision.Type {
    case NoDecision:
        // Behavior had no opinion, pass turn
        c.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
            EntityID:     req.Entity.ID,
            ControllerID: c.ID,
        }, nil))

    case Move, Flee:
        c.ruler.SendActor(message.Create(nil, rulermethods.ControllerMove{
            ControllerID: c.ID,
            EntityID:     req.Entity.ID,
            Path:         decision.Path,
        }, nil))

    case Attack:
        c.ruler.SendActor(message.Create(nil, rulermethods.ControllerAttack{
            ControllerID: c.ID,
            EntityID:     req.Entity.ID,
            Target:       decision.Target,
        }, nil))

    case Skill:
        c.ruler.SendActor(message.Create(nil, rulermethods.ControllerUseSkill{
            ControllerID: c.ID,
            EntityID:     req.Entity.ID,
            SkillID:      decision.SkillID,
            Target:       decision.Target,
        }, nil))

    case Pass:
        c.ruler.SendActor(message.Create(nil, rulermethods.EndOfTurn{
            EntityID:     req.Entity.ID,
            ControllerID: c.ID,
        }, nil))
    }
}
```

**Refactor AggressiveController**:

The existing `AggressiveController` logic moves to `AggressiveBehavior`.
The controller becomes a thin wrapper:

```go
type AggressiveController struct {
    *Controller  // Embed base controller
}

func NewAggressiveController(id uuid.UUID, name string) *AggressiveController {
    ctrl := &AggressiveController{
        Controller: &controller.Controller{
            ID:             id,
            Actor:          actor.New(name),
            Behavior:       &behavior.AggressiveBehavior{Name: name},
        },
    }

    // Set up handlers (same as before)
    // ...

    return ctrl
}
```

---

## PHASE 4: INTEGRATION & EXAMPLES

### 4.1 Turret Example

```go
// Create a turret that acts for 5 turns
func CreateTurret(gs *GameState, owner uuid.UUID, pos position.Position) {
    turret := entity.Entity{
        ID:           uuid.New(),
        ControllerID: owner,  // Assigned to player
        Type:         TimeBased,
        Position:     pos,
        Name:         "Turret",
        Properties: map[string]property.Property{
            property.HP:           def.MakeIntCounterProperty(property.HP, 10, 10, ...),
            property.Movement:     def.MakeIntCounterProperty(property.Movement, 0, 0, ...),  // Can't move
            property.AttackRange:  def.MakeIntProperty(property.AttackRange, 5, ...),  // Ranged
            property.Duration:     def.MakeIntCounterProperty(property.Duration, 5, 5, ...),  // 5 turns
            property.WalkThrough:  def.MakeBoolProperty(property.WalkThrough, true, ...),
        },
        Skills: map[uuid.UUID]skill.Skill{
            uuid.New(): createTurretAttackSkill(),
        },
    }

    gs.Entities[turret.ID] = turret
    gs.Turner.AddEntity(turret.ID, 0)  // Immediate first turn

    // Add to grid
    gs.Grid.AddEntity(pos, turret.ID)

    // Assign controller with AggressiveBehavior
    ctrl := controllers.NewAggressiveController(turret.ControllerID, "Turret AI")
    ctrl.EntityID = turret.ID
    gs.Controllers[turret.ControllerID] = ctrl.Communication()
}
```

### 4.2 Poisonous Fog Example

```go
// Player casts "Poison Fog" - lasts 3 turns
func CastPoisonFog(gs *GameState, caster entity.Entity, center position.Position) {
    // Pre-execution costs
    cost := 10  // MP cost
    caster.UpdatePropertyValue(property.MP, caster.GetPropertyI(property.MP).I() - cost)

    // Create zone (see Section 2.3 for full implementation)
    CreatePoisonousFog(gs, caster, center, 3)
}
```

### 4.3 Wall Example

```go
// Player creates a wall with 50 HP that decays
func CreateWall(gs *GameState, owner uuid.UUID, pos position.Position) {
    wall := entity.Entity{
        ID:           uuid.New(),
        ControllerID: uuid.Nil,  // No controller
        Type:         Obstacle,
        Position:     pos,
        Name:         "Stone Wall",
        Properties: map[string]property.Property{
            property.HP:          def.MakeIntCounterProperty(property.HP, 50, 50, ...),
            property.Duration:    def.MakeIntCounterProperty(property.Duration, 10, 10, ...),  // Decays in 10 turns
            property.WalkThrough: def.MakeBoolProperty(property.WalkThrough, false, ...),  // Blocks movement
        },
    }

    gs.Entities[wall.ID] = wall
    gs.Grid.AddEntity(pos, wall.ID)

    // No Turner entry needed - walls don't act
    // EndOfTurn will handle duration cleanup
}
```

### 4.4 Trap Example

```go
// Player places a bear trap that triggers once
func PlaceBearTrap(gs *GameState, owner uuid.UUID, pos position.Position) {
    trap := entity.Entity{
        ID:           uuid.New(),
        ControllerID: uuid.Nil,
        Type:         Trap,
        Position:     pos,
        Name:         "Bear Trap",
        Properties: map[string]property.Property{
            property.Duration:     def.MakeIntCounterProperty(property.Duration, 20, 20, ...),  // 20 turn timeout
            property.WalkThrough:  def.MakeBoolProperty(property.WalkThrough, true, ...),  // Can walk over
            property.Invisible:    def.MakeBoolProperty(property.Invisible, true, ...),  // Hidden
        },
    }

    gs.Entities[trap.ID] = trap
    gs.Turner.AddEntity(trap.ID, 100)  // Just for expiration tracking

    // Create cell effect
    trapEffect := effect.Effect{
        Properties: []property.Property{
            def.MakeIntProperty(property.Damage, 15, ...),
            def.MakeIntProperty(property.TriggerType, string(TriggerOnStep), ...),
            def.MakeBoolProperty(property.RemoveOnTrigger, true, ...),  // Single use
            def.MakeBoolProperty(property.ExpiresWithCaster, true, ...),  // Die with trap entity
            def.MakeIntProperty(property.ForceStopMove, 1, ...),  // Stop movement
        },
        Name:     "Bear Trap",
        CasterID: trap.ID,
    }

    effectID := uuid.New()
    gs.Effects[effectID] = trapEffect

    if cell, ok := gs.Grid.CellAt(pos); ok {
        cell.EffectIDs = append(cell.EffectIDs, effectID)
        gs.PositionalEffects[pos] = cell.EffectIDs
    }
}
```

### 4.5 Spooked + Aggressive Entity Example

```go
// Create a cowardly archer that flees when damaged, otherwise attacks
func CreateCowardlyArcher(gs *GameState, owner uuid.UUID, pos position.Position) {
    archer := entity.Entity{
        ID:           uuid.New(),
        ControllerID: owner,
        Type:         Character,
        Position:     pos,
        Name:         "Cowardly Archer",
        Properties: map[string]property.Property{
            property.HP:           def.MakeIntCounterProperty(property.HP, 20, 20, ...),
            property.Movement:     def.MakeIntCounterProperty(property.Movement, 5, 5, ...),
            property.AttackRange:  def.MakeIntProperty(property.AttackRange, 8, ...),
            property.Attack:       def.MakeIntProperty(property.Attack, 3, ...),
        },
        Skills: map[uuid.UUID]skill.Skill{
            uuid.New(): createBowAttackSkill(),
        },
    }

    gs.Entities[archer.ID] = archer
    gs.Turner.AddEntity(archer.ID, 100)

    // Composite controller: Spooked OR Aggressive
    composite := &behavior.CompositeBehavior{
        Name:  "Cowardly Archer",
        Mode:  behavior.CompositeOR,
        Behaviors: []behavior.BehaviorWrapper{
            {
                Behavior: &behavior.SpookedBehavior{
                    Name:         "Spooked",
                    ThresholdHP:  10,  // Flee if HP < 10
                    FleeDistance: 5,
                },
                Priority: 10,  // Check this first
            },
            {
                Behavior: &behavior.AggressiveBehavior{
                    Name: "Aggressive",
                },
                Priority: 5,  // If not spooked, attack
            },
        },
    }

    ctrl := controller.NewController(archer.ControllerID, "Cowardly Archer AI")
    ctrl.Behavior = composite
    ctrl.EntityID = archer.ID
    gs.Controllers[archer.ControllerID] = ctrl.Communication()
}
```

---

## PHASE 5: TERRAIN FEATURES (FUTURE)

### 5.1 Cell Crossing Cost

Add to `Cell` struct:

```go
type Cell struct {
    Type          CellType
    EntityIDs     []uuid.UUID
    EffectIDs     []uuid.UUID
    Position      position.Position
    CrossingCost  int  // Default 1, higher for difficult terrain
}
```

**Effects on Crossing Cost**:
- Cell's own CrossingCost (terrain-based)
- Effect's MvtCost property (applied during move)
- Combined: total cost = 1 + cell.CrossingCost + sum(effect.MvtCost)

**Example - Quagmire**:
- Create cell effect at position with MvtCost = 2
- Entity moving through pays 3 movement per cell (1 base + 2 effect)
- If entity has only 2 movement, it uses all 2 and stops (but can still ENTER the cell)

### 5.2 Future Considerations (Deferred)

- **VisibilityModifier**: For fog of war, invisibility (complex, deferred)
- **Cell-attached HP**: Cells that take damage and break
- **Dynamic terrain**: Terrain that changes over time (rising water, spreading fire)

---

## IMPLEMENTATION CHECKLIST

### Phase 1: Turner-Based Mechanics
- [x] Add `EntityType` enum values: TimeBased, Trap, AreaEffect, Obstacle
- [x] Add `Duration`, `ExpiresWithCaster`, `WalkThrough`, `Invisible` properties
- [x] Add `TriggerType`, `RemoveOnTrigger`, `TriggerCount` properties
- [x] Implement `RemoveEntity()` in GameState
- [x] Add expiration check in `EndOfTurn()`
- [ ] Test turret with AggressiveController (Movement=0)
- [ ] Test delayed effect entity
- [ ] Test wall/obstacle with HP and Duration

### Phase 2: Multi-Entity Cell & Positional Effects
- [x] Update `Cell` struct: EntityID → EntityIDs, add EffectIDs, CrossingCost
- [x] Update `Grid` methods for multiple entities
- [x] Add `PositionalEffects` and `Effects` to GameState
- [ ] Implement `CreatePoisonousFog()` zone pattern
- [ ] Implement trigger logic in `Move()` rule
- [ ] Implement `OnTurn` triggers in `BeginingOfTurn()`
- [ ] Implement effect cleanup in `RemoveEntity()`
- [ ] Add effect helper functions: getTriggerType, getRemoveOnTrigger, getTriggerCount

### Phase 3: Behavior System
- [x] Create `behavior` package
- [x] Define `Behavior` interface and `GameContext`
- [x] Define `Decision` struct and `DecisionType`
- [x] Implement `AggressiveBehavior`
- [x] Implement `ExpirationBehavior`
- [x] Implement `SpookedBehavior`
- [x] Implement `CompositeBehavior` with OR/AND modes
- [x] Refactor `Controller` to have `Behavior` field
- [ ] Update `ControllerNextTurn()` to use behavior decisions
- [ ] Refactor existing `AggressiveController` to use `AggressiveBehavior`
- [ ] Update all controller tests

### Phase 4: Integration
- [ ] Test turret example
- [ ] Test poisonous fog example
- [ ] Test wall example
- [ ] Test trap example
- [ ] Test spooked + aggressive composite
- [ ] Update ATD documentation

### Phase 5: Future
- [ ] Add `CrossingCost` to Cell
- [ ] Update movement cost calculation
- [ ] Test quagmire effect
- [ ] Defer visibility system

---

## RISK ASSESSMENT

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Cell structure change breaks pathfinding | Medium | High | Comprehensive testing, A* update |
| Effect cleanup performance | Low | Medium | Periodic cleanup, use maps |
| Behavior system complex to migrate | Medium | Medium | Incremental refactor, keep old controllers initially |
| Trigger edge cases (enter/exit/step) | High | Low | Clear documentation, extensive testing |
| Composite behavior conflicts | Medium | Medium | Priority system, clear OR/AND semantics |

---

## REFERENCES

- `upsilonbattle/battlearena/entity/entity.go` - Entity structure
- `upsilonbattle/battlearena/controller/controller.go` - Controller base
- `upsilonbattle/battlearena/controller/controllers/aggressive.go` - Aggressive controller
- `upsilonbattle/battlearena/property/effect/effect.go` - Effect structure
- `upsilonbattle/battlearena/property/propertyenum.go` - Property enums
- `upsilonbattle/battlearena/ruler/rules/gamestate.go` - GameState
- `upsilonbattle/battlearena/ruler/rules/move.go` - Movement logic
- `upsilonbattle/battlearena/ruler/rules/beginingofturn.go` - Turn start
- `upsilonbattle/battlearena/ruler/rules/endofturn.go` - Turn end
- `upsilonmapdata/grid/grid.go` - Grid structure
- `upsilonmapdata/grid/cell/cell.go` - Cell structure
- `upsilonbattle/battlearena/ruler/turner/turner.go` - Turner system
- `upsilonbattle/battlearena/property/effect/effectapplicator/effectapplicator.go` - Effect application

---

## OPEN QUESTIONS FOR DISCUSSION

1. **PositionalEffects storage**: Confirmed in GameState, Cell stores EffectIDs
2. **Zone effect caster**: Confirmed as the anchor entity ID
3. **Effect cleanup**: Confirmed via CasterID + ExpiresWithCaster flag
4. **Trigger location**: Confirmed in Move rule (for movement triggers) and BeginingOfTurn (for OnTurn)
5. **WalkThrough**: Confirmed as new property
6. **Behavior refactor**: Confirmed - refactor AggressiveController
7. **Decision response**: Confirmed - use Decision struct with NoDecision fallback
8. **Movement cost timing**: Confirmed - during move execution, not before validation
9. **Effect stacking order**: Confirmed - First Come First Served
10. **Turret/Obstacle Turner**: Confirmed - both have controllers (wall has only ExpirationBehavior)
