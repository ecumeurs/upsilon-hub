---
id: mechanic_channeling_mechanic
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
To implement channeling mechanic where skills have pre-execution delay time during which the caster is vulnerable but committed to executing the skill, with interruption possibility if caster is damaged.

## THE RULE / LOGIC
**Channeling Mechanic:**

**Core Concept:**
Channeling is a pre-execution delay where the caster spends their turn "charging" a skill, becoming vulnerable to interruption but committed to execution.

**Channeling Flow:**

**1. Channeling Initiation:**
- **Player Action:** Player selects channeling skill (e.g., "Fireball" with 400 delay)
- **Commitment:** Character's turn ends immediately, skill enters channeling state
- **Vulnerability:** Character cannot move, act, or defend during channeling
- **Interruptability:** Any damage taken during channeling interrupts the skill

**2. Channeling Duration:**
- **Delay Measurement:** Delay measured in action economy units (same as standard delays)
- **Progress Tracking:** Display channeling progress (e.g., "75% charged")
- **Visual Feedback:** Character animation shows channeling effect
- **Duration:** Variable by skill (300-600 delay common range)

**3. Channeling Completion:**
- **Automatic Execution:** Skill executes automatically when channeling completes
- **No Action Required:** Skill fires at pre-selected target automatically
- **Effect Application:** Standard skill effect application rules apply
- **Cost Deduction:** All costs (MP, SP, HP) deducted at channeling start

**Channeling Skill Properties:**
```go
type ChannelingSkill struct {
    skill.Skill
    ChannelingDelay    int    // Pre-execution delay (e.g., 400)
    ChannelingInterrupt bool   // Can be interrupted by damage
    ChannelingProgress  int    // 0-100% completion tracking
    TargetWhenChanneling Position  // Target set when channeling starts
}
```

**Interruption Rules:**

**Damage Interruption:**
- **Any Damage:** Any damage taken during channeling interrupts skill
- **Interrupt Timing:** Interrupt happens immediately on damage receipt
- **Effect Cancellation:** Skill does not execute, effect cancelled
- **Cost Loss:** Resource costs (MP, SP) are consumed but skill fails

**Interruption Exceptions:**
- **Self-Inflicted Damage:** Damage from own skills/poison does NOT interrupt
- **Damage During Execution:** Once channeling completes, damage no longer interrupts
- **Status Effect Immunity:** Some channeling skills may have "Cannot Be Interrupted" flag

**Channeling vs Standard Delay:**

**Standard Delay:**
- **Action-Based:** Character performs action, delay is time before next turn
- **No Vulnerability:** Character can defend, move, or act during delay period
- **Example:** "Heavy Attack" has 800 delay but character can react during this time

**Channeling Delay:**
- **Pre-Execution:** Delay occurs before skill executes, not after
- **Maximum Vulnerability:** Character cannot act or defend during channeling
- **Commitment:** Character is locked into skill execution regardless of battlefield changes
- **Example:** "Fireball" has 400 channeling delay, then executes instantly

**Channeling Examples:**

**"Fireball" Spell:**
- **Channeling:** 400 delay
- **Damage:** 25 damage (Grade III skill, 250 SW)
- **Range:** 5 cells (AoE 3x3 area)
- **Risk:** Caster vulnerable for 400 delay if attacked

**"Healing Zone" Support Skill:**
- **Channeling:** 300 delay
- **Effect:** Creates healing zone (temporary entity) for 3 turns
- **Target:** Position-based (zone placement)
- **Risk:** Support character vulnerable during channeling

**"Snipe" Ranger Skill:**
- **Channeling:** 200 delay (shorter due to ranger archetype)
- **Damage:** 15 damage (Grade II skill, 120 SW)
- **Range:** 7 cells (long range)
- **Critical Bonus:** +50% critical chance on completion

**Risk/Reward Balance:**

**Channeling Benefits:**
- **High Power:** Channeling skills typically have higher damage/effects
- **Strategic Timing:** Player commits to powerful future action
- **Zone Effects:** Enables area denial, healing zones, traps
- **Tactical Depth:** "Should I cast this 400-delay spell knowing I could be interrupted?"

**Channeling Risks:**
- **Vulnerability:** Can't move, act, or defend during channeling
- **Interruption:** Damage cancels skill, loses resource costs
- **Predictability:** Enemies can see channeling and position accordingly
- **Mistake Cost:** Choosing wrong timing = wasted resources

**Integration with Temporary Entity System:**
Channeling skills use the unified temporary entity system:
```go
func StartChanneling(character Character, skill ChannelingSkill) TemporaryEntity {
    // Create temporary entity representing channeling process
    channelingEntity := TemporaryEntity{
        CasterID:    character.ID,
        TrueEffect:   skill,
        TriggerType:  "OnTurn",      // Execute after delay completes
        Duration:     skill.ChannelingDelay,
    }
    
    // Add character to vulnerable state
    character.State = "CHANNELING"
    
    // Create temporary entity on same tile as character
    grid.AddEntity(channelingEntity, character.Position)
    
    return channelingEntity
}
```

**Player Experience:**
"Turn 3: I start channeling 'Fireball' (400 delay). The enemy attacks me for 10 damage—my Fireball is interrupted! I lost the 15 MP cost and didn't deal any damage. Next time I need to position better before channeling."

"Turn 5: I position behind a wall, then channel 'Healing Zone' (300 delay). This time the enemy can't reach me, so the zone appears and heals my ally for the next 3 turns. Perfect timing!"

**Implementation Priority:** HIGH - Required for Phase 2 time-based mechanics

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[channeling_mechanic]]`
- **Related Files:** `upsilonbattle/battlearena/entity/skill/skill.go`, `upsilonbattle/battlearena/ruler/rules/execution.go`
- **Integration:** Works with `mechanic_mech_temporary_entity_system` and `effect_caster_tracking`

## EXPECTATION
