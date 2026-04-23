---
id: temp
status: DRAFT
version: 1.0
parents: []
dependents: []
---

# New Atom

## INTENT
To implement individual AI controller archetypes for Fighter, Ranger, Support, and Sneak behaviors, defining specific decision trees, skill priorities, and tactical patterns for each archetype.

## THE RULE / LOGIC
**AI Controller Archetypes:**

**Core Principle:**
Each AI archetype extends base controller with specialized behavior, skill pools, and tactical priorities.

**Base Controller Interface:**
```go
type AIController interface {
    // Common methods
    DecideAction() Action
    SelectTarget() Entity
    SelectSkill() skill.Skill
    EvaluatePosition() Position
    
    // Archetype-specific methods
    GetArchetypePriority() []ArchetypePriority
    GetSkillPool() []skill.Skill
    ExecuteArchetypeBehavior() Action
}
```

**Fighter Controller Archetype:**

**Behavioral Characteristics:**
- **Aggressive Tactics:** Prioritize closing distance to enemies
- **Frontline Positioning:** Prefer positions near enemy team
- **Damage Focus:** Prioritize high-damage skills and attacks
- **Direct Approach:** Avoid stealth/flanking, favor direct confrontation

**Skill Pool:**
- **High Damage:** Skills with high damage values
- **Defensive Skills:** Shields, damage reduction, self-buffs
- **Charge Abilities:** Movement skills to close gaps quickly
- **Single Target:** Skills focused on 1v1 engagements

**Decision Tree:**
```go
func (fc *FighterController) DecideAction() Action {
    // Priority 1: Can I kill someone?
    target := FindKillableEnemy(fc.CurrentEnemy)
    if target != nil {
        return AttackAction(target, fc.GetHighestDamageSkill())
    }
    
    // Priority 2: Should I charge?
    if fc.CanCharge() {
        chargeTarget := FindChargingPosition()
        return MoveAction(chargeTarget) + AttackAction(fc.ChargeSkill)
    }
    
    // Priority 3: Attack nearest enemy
    nearestEnemy := FindNearestEnemy(fc.Position)
    return AttackAction(nearestEnemy, fc.GetBestMeleeSkill())
}
```

**Stat Allocation Pattern:**
- **Attack:** Highest priority (60% of points)
- **Defense:** High priority (25% of points)
- **HP:** Medium priority (10% of points)
- **Movement:** Low priority (5% of points)

**Ranger Controller Archetype:**

**Behavioral Characteristics:**
- **Kiting Tactics:** Maintain optimal distance from enemies
- **Positioning Focus:** Prefer high ground, chokepoints, cover
- **Precision Skills:** Prioritize accuracy and critical chance
- **Trap Placement:** Use terrain advantages, deny enemy movement

**Skill Pool:**
- **Ranged Damage:** Skills with range 4-7 cells
- **Trap Skills:** Area denial, movement penalties
- **Movement Skills:** Dash, teleport, positioning abilities
- **Precision Bonuses:** High accuracy, critical skills

**Decision Tree:**
```go
func (rc *RangerController) DecideAction() Action {
    // Priority 1: Maintain optimal kiting range
    if !rc.IsInOptimalRange() {
        optimalPos := FindOptimalRangePosition()
        return MoveAction(optimalPos)
    }
    
    // Priority 2: Place traps
    if rc.ShouldPlaceTrap() {
        trapLocation := FindTrapChokepoint()
        return SkillAction(rc.TrapSkill, trapLocation)
    }
    
    // Priority 3: Attack from safe position
    target := FindExposedEnemy(rc.Position)
    return AttackAction(target, rc.GetBestRangedSkill())
}
```

**Stat Allocation Pattern:**
- **Attack:** Medium-high priority (50% of points)
- **Accuracy:** High priority (20% of points)
- **Movement:** Medium priority (15% of points)
- **Defense:** Medium priority (10% of points)
- **HP:** Low priority (5% of points)

**Support Controller Archetype:**

**Behavioral Characteristics:**
- **Ally Proximity:** Stay near teammates
- **Defensive Focus:** Prioritize healing, shielding, protection
- **Buff Application:** Enhance ally combat capabilities
- **Positioning:** Block chokepoints, protect vulnerable allies

**Skill Pool:**
- **Healing Skills:** Restore HP to allies
- **Shield Skills:** Provide temporary damage absorption
- **Buff Skills:** Enhance ally attack/defense
- **Zone Creation:** Healing areas, protective barriers

**Decision Tree:**
```go
func (sc *SupportController) DecideAction() Action {
    // Priority 1: Heal damaged ally
    damagedAlly := FindMostDamagedAlly(sc.Allies)
    if damagedAlly != nil {
        return SkillAction(sc.HealingSkill, damagedAlly)
    }
    
    // Priority 2: Shield ally about to take damage
    threatenedAlly := FindThreatenedAlly(sc.Allies)
    if threatenedAlly != nil {
        return SkillAction(sc.ShieldSkill, threatenedAlly)
    }
    
    // Priority 3: Apply buffs to allies
    allyNeedingBuff := FindAllyNeedingBuff(sc.Allies)
    if allyNeedingBuff != nil {
        return SkillAction(sc.BuffSkill, allyNeedingBuff)
    }
    
    // Priority 4: Move to optimal support position
    optimalPos := FindSupportPosition(sc.Allies)
    return MoveAction(optimalPos)
}
```

**Stat Allocation Pattern:**
- **MP/SP Pools:** High priority (40% of points)
- **Defense:** Medium priority (25% of points)
- **HP:** Medium priority (20% of points)
- **Healing Power:** High priority (15% of points)
- **Attack:** Low priority (0% of points)

**Sneak Controller Archetype:**

**Behavioral Characteristics:**
- **Flanking Tactics:** Prioritize positioning behind enemies
- **Backstab Focus:** Seek backstab opportunities aggressively
- **Evasion Skills:** Dodge, stealth, positioning abilities
- **Status Application:** Poison, stun, debuff priorities

**Skill Pool:**
- **Backstab Skills:** Enhanced backstab damage, positioning
- **Movement Skills:** Dash, teleport, stealth
- **Poison/Stun:** Status effect application skills
- **Evasion:** Dodge, self-defense abilities

**Decision Tree:**
```go
func (sn *SneakController) DecideAction() Action {
    // Priority 1: Backstab vulnerable enemy
    target := FindBackstabableEnemy(sn.Position)
    if target != nil {
        sneakPos := FindBehindPosition(target.Position)
        return MoveAction(sneakPos) + AttackAction(target, sn.BackstabSkill)
    }
    
    // Priority 2: Apply poison/stun to priority target
    priorityTarget := FindHighestValueTarget(enemies)
    if priorityTarget != nil {
        return SkillAction(sn.PoisonSkill, priorityTarget)
    }
    
    // Priority 3: Evade incoming attacks
    if sn.IsUnderThreat() {
        evasionPos := FindEvasionPosition()
        return MoveAction(evasionPos) + SkillAction(sn.DodgeSkill, sn.Self)
    }
    
    // Priority 4: Position for future backstabs
    futurePos := FindFlankingPosition(sn.AllEnemies)
    return MoveAction(futurePos)
}
```

**Stat Allocation Pattern:**
- **Movement:** Highest priority (30% of points)
- **Attack:** Medium priority (25% of points)
- **Dodge:** High priority (20% of points)
- **Critical Chance:** Medium priority (15% of points)
- **HP/Defense:** Low priority (10% of points)

**Archetype Integration:**

**Base Controller Extension:**
```go
type FighterController struct {
    BaseController
    Archetype    "Fighter"
    SkillPool    []skill.Skill
    Personality   FighterPersonality
}

func (fc *FighterController) GetSkillPool() []skill.Skill {
    // Filter skills by archetype appropriateness
    availableSkills := GetAllSkills()
    fighterSkills := FilterByArchetype(availableSkills, "Fighter")
    
    return fighterSkills
}
```

**Skill Selection Logic:**
- **Damage Priority:** Fighter prefers highest damage skills
- **Range Priority:** Ranger prefers ranged skills with accuracy
- **Utility Priority:** Support prefers healing/shielding
- **Positional Priority:** Sneak prefers backstab/evasion skills

**Combat Engagement:**
```go
func (ai AIController) EngageCombat(enemy Entity) {
    archetype := ai.GetArchetype()
    
    switch archetype {
        case "Fighter":
            return FighterEngageBehavior(ai, enemy)
            
        case "Ranger":
            return RangerEngageBehavior(ai, enemy)
            
        case "Support":
            return SupportEngageBehavior(ai, enemy)
            
        case "Sneak":
            return SneakEngageBehavior(ai, enemy)
    }
}
```

**Team Composition Rules:**

**Archetype Limits:**
- **Maximum 1 Support:** Only one support AI per team
- **Maximum 1 Sneak:** Only one sneak AI per team
- **Unlimited Fighters/Rangers:** Can have multiple of each
- **Composition Balance:** 2 Fighters + 1 Ranger + 1 Support + 1 Sneak = 5 AI team

**Selection Algorithm:**
```go
func SelectAIComposition(teamSize int) []AIController {
    controllers := []AIController{}
    
    // Always include 1 Fighter
    controllers = append(controllers, CreateFighterController())
    
    // Always include 1 Ranger
    controllers = append(controllers, CreateRangerController())
    
    // Add 1 Support if team size permits
    if teamSize >= 3 {
        controllers = append(controllers, CreateSupportController())
    }
    
    // Add 1 Sneak if team size permits
    if teamSize >= 4 {
        controllers = append(controllers, CreateSneakController())
    }
    
    // Fill remaining slots with Fighters (most flexible)
    for len(controllers) < teamSize {
        controllers = append(controllers, CreateFighterController())
    }
    
    return controllers
}
```

**Performance Optimization:**

**Behavior Caching:**
```go
type BehaviorCache struct {
    Archetype       string
    DecisionTree    DecisionTree
    CachedTurn      int
}

func GetCachedBehavior(archetype string) DecisionTree {
    if cached, exists := behaviorCache[archetype]; exists {
        if cached.CachedTurn == currentTurn {
            return cached.DecisionTree  // Use cached behavior
        }
    }
    
    // Generate new behavior based on archetype
    behavior := GenerateArchetypeBehavior(archetype)
    
    behaviorCache[archetype] = BehaviorCache{
        Archetype:    archetype,
        DecisionTree:  behavior,
        CachedTurn:     currentTurn,
    }
    
    return behavior
}
```

**Player Experience:**

**Fighter AI:**
"The Fighter AI charges straight at my character—no flanking or positioning, just direct damage. I can tell it's a Fighter because it prioritizes closing distance and uses high-damage skills. Its stat build focuses on Attack and Defense!"

**Ranger AI:**
"The Ranger AI keeps backing away to maintain optimal range. It places traps in chokepoints and uses ranged precision attacks. I notice it's prioritizing accuracy and movement over raw damage, typical Ranger behavior!"

**Support AI:**
"The Support AI stays close to its teammates, always ready to heal or shield. It positions to protect vulnerable allies and uses buff skills strategically. The team composition rule (max 1 support) makes it a precious teammate!"

**Sneak AI:**
"The Sneak AI is constantly trying to get behind me! It uses movement skills to flank and applies poisons when I can't see it. The backstab damage hurts, and I can't easily predict where it'll appear. Tricky opponent!"

**Implementation Priority:** HIGH - Required for Phase 3 AI enhancement

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ai_controller_archetypes]]`
- **Related Files:** `upsilonbattle/battlearena/controller/controllers/fighter.go`, `upsilonbattle/battlearena/controller/controllers/ranger.go`, etc.
- **Integration:** Works with `mec_ai_archetype_system`, `ai_progression_matching`

## EXPECTATION
