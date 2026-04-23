---
id: mechanic_effect_caster_tracking
status: DRAFT
priority: 5
version: 2.0
parents: []
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
---

# New Atom

## INTENT
To implement effect caster tracking system where all combat effects remember their original caster until effect ends, enabling proper credit assignment, interruption mechanics, and support play credit earning.

## THE RULE / LOGIC
**Effect Caster Tracking System:**

**Core Principle:**
All combat effects (damage, healing, shields, status effects, temporary entities) must track their original caster (CasterID) for the entire duration of the effect.

**Why Caster Tracking Matters:**

**Credit Assignment:**
- **Damage Credits:** Go to damage dealer, not effect applicator
- **Healing Credits:** Go to healer, not effect recipient
- **Shield Credits:** Go to shield caster, not shielded character
- **Support Play:** Enables "healer" role where support earns credits from shield/heal effects

**Interruption Mechanics:**
- **Channeling Interruption:** Know which caster is interrupted when skill cancelled
- **Temporary Entity Control:** Know which player to reward when entity expires
- **Effect Removal:** Clean up effects tied to disconnected players

**Debugging & Balancing:**
- **Effect Origin:** Trace which skills/players caused specific effects
- **Abuse Detection:** Identify exploitative play patterns
- **Performance Analysis:** Track which skills are over/underperforming

**Caster Tracking Implementation:**

**Effect Data Structure:**
```go
type TrackedEffect struct {
    EffectID    uuid.UUID
    CasterID     uuid.UUID    // Original creator
    EffectType   EffectType // Shield, Heal, Poison, Stun, etc.
    TargetID      uuid.UUID    // Who is affected
    CreationTurn  int          // When effect was created
    Duration      int          // How many turns effect lasts
    Intensity     property.Property // Effect power (e.g., shield amount, poison damage)
}
```

**Caster Tracking in Effect Application:**

**Shield Example:**
```go
func ApplyShieldEffect(caster Character, target Character, amount int) {
    shieldEffect := TrackedEffect{
        CasterID:   caster.ID,      // Support player who cast shield
        EffectType: "Shield",
        TargetID:    target.ID,
        Intensity:   amount,
        CreationTurn: currentTurn,
        Duration:     3,
    }
    
    // Add shield to target
    target.Shield += amount
    
    // Track shield for credit assignment
    target.ActiveEffects = append(target.ActiveEffects, shieldEffect)
}
```

**Poison Example:**
```go
func ApplyPoisonEffect(caster Character, target Character, damage int, duration int) {
    poisonEffect := TrackedEffect{
        CasterID:   caster.ID,      // Caster earns poison credits
        EffectType: "Poison",
        TargetID:    target.ID,
        Intensity:   damage,        // Per-turn damage
        CreationTurn: currentTurn,
        Duration:     duration,
    }
    
    target.Poison += damage
    target.ActiveEffects = append(target.ActiveEffects, poisonEffect)
}
```

**Caster Death Handling:**

**Support Credits After Death:**
- **Shield Continues:** Shield blocks damage even after caster dies → credits go to dead caster
- **Channeling Interruption:** If caster dies during channeling → skill cancelled, no credits
- **Poison Damage:** Poison damage continues after caster death → credits go to original caster
- **Zone Effects:** Area effects continue after caster death → credits assigned to original caster

**Credit Assignment with Caster Tracking:**
```go
func AssignCombatCredits(damage int, effects []TrackedEffect) map[uuid.UUID]int {
    credits := make(map[uuid.UUID]int)
    
    // Assign damage credits
    credits[damageDealer.ID] += damage
    
    // Assign shield credits
    for _, effect := range effects {
        if effect.EffectType == "Shield" && effect.BlocksDamage > 0 {
            credits[effect.CasterID] += effect.BlocksDamage
        }
    }
    
    return credits
}
```

**Effect Duration Management:**

**Turn-by-Turn Processing:**
```go
func ProcessTrackedEffects(entities []Entity, currentTurn int) {
    for _, entity := range entities {
        for i, effect := range entity.ActiveEffects {
            // Calculate effect age
            effectAge := currentTurn - effect.CreationTurn
            
            // Check expiration
            if effectAge >= effect.Duration {
                // Effect expired - remove from entity
                entity.RemoveEffect(i)
                
                // Trigger expiration callbacks
                effect.OnExpiration()
                continue
            }
            
            // Apply effect per-turn logic
            switch effect.EffectType {
                case "Poison":
                    entity.HP -= effect.Intensity  // Poison damage
                    credits[effect.CasterID] += effect.Intensity
                case "HealZone":
                    entity.HP += effect.Intensity  // Zone healing
                    credits[effect.CasterID] += effect.Intensity
                case "Stun":
                    entity.CannotAct = true  // Stun effect
            }
        }
    }
}
```

**Caster Tracking Benefits:**

**Support Play Validation:**
- **Healer Recognition:** Healers earn credits from shield/heal effects
- **Team Coordination:** Support players are rewarded for protecting allies
- **Strategic Play:** "I'll shield my teammate and earn credits even if I die"

**Channeling Mechanism:**
- **Interruption Awareness:** Know which caster's channeling was interrupted
- **Resource Loss Tracking:** Caster loses MP/SP when channeling interrupted
- **Fair Balance:** Risky channeling has clear penalty

**Multi-Entity Effects:**
- **Zone Ownership:** Know who created poison zones, healing areas
- **Proper Credit Assignment:** Zone damage credited to original creator
- **Effect Cleanup:** Remove zone effects when duration expires

**Integration Points:**

**Temporary Entity System:**
- **CasterID Field:** Every temporary entity tracks its creator
- **Effect Execution:** Triggers assign credits to CasterID, not executor
- **Expiration:** Temporary entities know which player to notify on completion

**Credit Economy:**
- **Damage Tracking:** Simple damage dealer assignment
- **Support Tracking:** Complex effect-based credit calculation
- **Caster Resolution:** All credits eventually assigned to original casters

**Status Effects:**
- **Flat Rate Credits:** Status effects use SkillWeight/10 formula
- **Caster Assignment:** Credits go to effect applicator
- **No Per-Turn Credits:** One-time reward at application, not per-turn

**Debugging & Maintenance:**

**Effect Audit Trail:**
```go
type EffectAuditLog struct {
    EffectID      uuid.UUID
    CasterID       uuid.UUID
    AppliedTurn     int
    ExpiredTurn     int
    CreditsEarned   int
    TargetsAffected  []uuid.UUID
}

// Log all effects for balance analysis
func LogEffect(effect TrackedEffect) {
    audit := EffectAuditLog{
        EffectID:       effect.EffectID,
        CasterID:        effect.CasterID,
        AppliedTurn:     effect.CreationTurn,
        CreditsEarned:   CalculateEffectCredits(effect),
        TargetsAffected:  effect.Targets,
    }
    
    effectAuditLog.Append(audit)
}
```

**Player Experience:**

**Support Player Scenario:**
"I cast 'Shield' on my teammate for 20 HP. The enemy attacks my teammate for 15 damage—the shield blocks it all. Even though I'm not taking damage directly, I earn 15 credits for the protection I provided. Later that turn, I die to a different attack, but my shield continues protecting my teammate, and I continue earning shield credits. This makes playing support meaningful!"

**Channeling Player Scenario:**
"I start channeling 'Fireball' (400 delay). The enemy attacks me for 10 damage—my channeling is interrupted! I lose the 15 MP cost and the skill doesn't execute. The system knows I was the caster, so it properly cancels my skill and deducts the costs. Next time I'll position behind cover before channeling."

**Implementation Priority:** HIGH - Critical for Phase 2 time-based mechanics and credit economy

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[effect_caster_tracking]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entity.go` (effect tracking), `upsilonbattle/battlearena/ruler/rules/effects.go`
- **Integration:** Works with `channeling_mechanic`, `mec_credit_spending_shop`, `rule_credit_earning_support`

## EXPECTATION
