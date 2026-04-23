---
id: temp
status: DRAFT
version: 1.0
parents: []
dependents: []
---

# New Atom

## INTENT
To implement AI team composition rules that enforce maximum limits on archetype variety per AI team, ensuring balanced team composition with appropriate support and specialist limitations.

## THE RULE / LOGIC
**AI Team Composition Rules:**

**Core Principle:**
AI teams must follow composition constraints to ensure balanced gameplay and prevent overpowered combinations.

**Composition Limits:**

**Archetype Maximums:**
- **Maximum 1 Support:** Only one support AI per team
- **Maximum 1 Sneak:** Only one sneak AI per team
- **Unlimited Fighter:** Can have multiple fighter AI
- **Unlimited Ranger:** Can have multiple ranger AI

**Team Size Configuration:**

**1v1 Match (1 AI Team):**
- **Composition:** 1 AI (any archetype allowed)
- **No Restrictions:** Single AI has no composition limits
- **Flexibility:** Best opportunity for player to learn patterns

**2v2 Match (2 AI Team):**
- **Recommended:** 1 Fighter + 1 Ranger
- **Alternative:** 2 Fighters, 1 Fighter + 1 Support
- **Restriction:** Cannot have 2 Support or 2 Sneak

**3v3 Match (3 AI Team):**
- **Recommended:** 1 Fighter + 1 Ranger + 1 Support
- **Alternative:** 2 Fighters + 1 Ranger, 1 Fighter + 1 Support + 1 Sneak
- **Restriction:** Cannot have 2 Support or 2 Sneak

**4v4 Match (4 AI Team):**
- **Recommended:** 2 Fighters + 1 Ranger + 1 Support
- **Alternative:** 1 Fighter + 1 Ranger + 1 Support + 1 Sneak
- **Restriction:** Cannot have 2 Support or 2 Sneak

**Archetype Definitions:**

**Fighter AI:**
- **Role:** Frontline combat, damage dealing
- **Stats:** High Attack/Defense focus
- **Skills:** High damage melee, charge abilities
- **Tactics:** Direct engagement, aggressive positioning

**Ranger AI:**
- **Role:** Ranged combat, area control
- **Stats:** Attack/Accuracy/Movement focus
- **Skills:** Ranged damage, trap placement, positioning
- **Tactics:** Kiting, optimal range maintenance

**Support AI:**
- **Role:** Ally protection, healing, buffing
- **Stats:** MP/SP pools, Defense focus
- **Skills:** Healing, shielding, buff/debuff application
- **Tactics:** Ally proximity, defensive positioning

**Sneak AI:**
- **Role:** Flanking, backstabbing, status effects
- **Stats:** Movement/Dodge/CritChance focus
- **Skills:** Backstab bonuses, evasion, poison/stun application
- **Tactics:** Positional advantage, stealth approach

**Composition Validation Logic:**

```go
func ValidateAIComposition(team []AIController) error {
    archetypeCounts := make(map[string]int)
    
    // Count each archetype
    for _, ai := range team {
        archetype := ai.GetArchetype()
        archetypeCounts[archetype]++
    }
    
    // Validate Support limit
    if archetypeCounts["Support"] > 1 {
        return errors.New("Cannot have more than 1 Support AI")
    }
    
    // Validate Sneak limit
    if archetypeCounts["Sneak"] > 1 {
        return errors.New("Cannot have more than 1 Sneak AI")
    }
    
    // Validate total team size
    if len(team) > 5 {
        return errors.New("AI team cannot exceed 5 members")
    }
    
    return nil
}
```

**Team Generation Algorithm:**

```go
func GenerateAITeam(teamSize int, playerLevel int) []AIController {
    team := []AIController{}
    
    // Add mandatory archetypes first
    if teamSize >= 2 {
        fighter1 := CreateFighterAI(playerLevel)
        team = append(team, fighter1)
    }
    
    if teamSize >= 3 {
        ranger1 := CreateRangerAI(playerLevel)
        team = append(team, ranger1)
    }
    
    if teamSize >= 4 {
        support := CreateSupportAI(playerLevel)
        team = append(team, support)
    }
    
    if teamSize >= 5 {
        sneak := CreateSneakAI(playerLevel)
        team = append(team, sneak)
    }
    
    // Fill remaining slots with flexible archetypes
    for len(team) < teamSize {
        flexibleAI := CreateFlexibleAI(playerLevel, team)
        team = append(team, flexibleAI)
    }
    
    // Validate final composition
    if err := ValidateAIComposition(team); err != nil {
        // Retry generation if invalid
        return GenerateAITeam(teamSize, playerLevel)
    }
    
    return team
}
```

**Flexible Archetype Selection:**

**Player Level Matching:**
- **AI Level:** Equals average player team level
- **Stat Allocation:** AI follows same CP system as players
- **Skill Access:** AI gets skills from same grade pools as players
- **Progression Sync:** AI +10 CP per win, same as players

**Archetype Priorities:**
- **Counter-Picking:** AI composition counters player team composition
- **Balanced Approach:** Mix of archetypes for tactical variety
- **Random Variation:** Some randomness in flexible slot selection
- **Difficulty Scaling:** Higher player levels = more complex AI behavior

**Composition Templates:**

**Standard Templates (2v2):**
```
Template A (Balanced):
- 1 Fighter (frontline damage)
- 1 Ranger (ranged support)

Template B (Aggressive):
- 2 Fighters (pure damage)

Template C (Tactical):
- 1 Fighter (damage + utility)
- 1 Support (healing + buffs)

Template D (Ranged Focus):
- 1 Ranger (ranged primary)
- 1 Fighter (ranged secondary)
```

**Advanced Templates (3v3):**
```
Template A (Ideal Composition):
- 1 Fighter (frontline)
- 1 Ranger (ranged support)
- 1 Support (healing + protection)

Template B (Triple Threat):
- 2 Fighters (double frontline)
- 1 Support (sustained combat)

Template C (Squad Tactics):
- 1 Fighter (damage)
- 1 Ranger (positioning)
- 1 Sneak (flanking + backstab)
```

**Specialized Compositions:**

**Anti-Meta Teams:**
- **Tank Buster:** 2 Rangers + 1 Support (kiting strategy)
- **Anti-Stealth:** Multiple detection skills, area denial
- **Anti-Burst:** High Defense composition to withstand burst damage

**Progressive Difficulty:**
```go
func IncreaseTeamDifficulty(team []AIController, playerLevel int) {
    // Analyze player performance
    playerWinRate := CalculatePlayerWinRate()
    
    // Adjust AI behavior based on player skill
    if playerWinRate > 0.7 {
        // Player winning frequently → increase AI difficulty
        for _, ai := range team {
            ai.IncreaseAggression()
            ai.ImproveTargeting()
        }
    } else if playerWinRate < 0.3 {
        // Player struggling → moderate AI
        for _, ai := range team {
            ai.ModerateAggression()
            ai.SlightlyImproveTargeting()
        }
    }
    
    return team
}
```

**Dynamic Composition Adjustment:**

**Mid-Match Substitution:**
- **Support Replacement:** If Support AI dies, may swap in Fighter/Ranger
- **Sneak Replacement:** If Sneak AI dies, may swap in Ranger
- **No Hard Replacement:** Team composition adapts rather than enforcing strict roles

**Composition Balance Metrics:**

**Archetype Distribution:**
- **Optimal Range:** 40-60% Fighters, 20-40% Rangers/Supports, 20-40% Sneaks
- **Flexibility:** Allow variation while maintaining constraints
- **Counter-Play:** Different compositions should have different playstyles

**Team Effectiveness:**
```go
func CalculateTeamCompositionScore(team []AIController) float64 {
    score := 0.0
    
    // Balance bonus for archetype variety
    archetypes := GetUniqueArchetypes(team)
    score += float64(len(archetypes)) * 10
    
    // Combat synergy bonus
    synergy := CalculateCombatSynergy(team)
    score += synergy * 15
    
    // Counter-play potential
    counterScore := CalculateCounterPotential(team, GetPlayerTeamComposition())
    score += counterScore * 20
    
    return score
}
```

**Integration with AI Progression:**

**Level Synchronization:**
- **Same CP System:** AI uses 100 CP + wins × 10 formula
- **Skill Access:** AI skill selection mirrors player level access
- **Stat Caps:** AI follows same maximums as players
- **Difficulty Scaling:** AI matches average player team level

**Archetype Scaling:**
```go
func ScaleArchetypeToLevel(archetype string, level int) AIController {
    baseStats := GetV2BaselineStats()  // HP 30-50, Attack 10, etc.
    
    // Allocate points based on level
    availableCP := 100 + (level * 10)
    stats := AllocatePointsByArchetype(archetype, availableCP)
    
    // Select skills appropriate for level
    skills := SelectSkillsByGrade(level, archetype)
    
    ai := AIController{
        Archetype: archetype,
        Stats:     stats,
        Skills:     skills,
        Level:      level,
    }
    
    return ai
}
```

**Performance Optimization:**

**Composition Caching:**
```go
type CompositionCache struct {
    TeamSize   int
    Templates   [][]AIController
    CachedTurn  int
}

func GetValidComposition(teamSize int) []AIController {
    cacheKey := fmt.Sprintf("size-%d", teamSize)
    
    if cached, exists := compositionCache[cacheKey]; exists {
        if cached.CachedTurn == currentTurn {
            // Return random template from cache
            return cached.Templates[rand.Intn(len(cached.Templates))]
        }
    }
    
    // Generate and cache new compositions
    templates := GenerateValidCompositions(teamSize)
    
    compositionCache[cacheKey] = CompositionCache{
        TeamSize:  teamSize,
        Templates:   templates,
        CachedTurn:  currentTurn,
    }
    
    return GetValidComposition(teamSize)
}
```

**Player Experience:**

**Balanced AI Team:**
"I face 2v2 match: 1 Fighter + 1 Ranger. The Fighter charges at my warrior while the Ranger keeps backing away, using ranged attacks. The team composition is balanced—damage dealer + ranged support. The AI level scaling matches our team's progression!"

**Support AI Strategy:**
"There's 1 Support AI in this 3v3 match! It stays between its Fighter teammates, ready to heal or shield. The composition rules limit teams to max 1 Support, which makes the Support AI crucial for team survival. When our team pressures, the Support keeps us in the fight!"

**Sneak Integration:**
"The 3v3 AI team includes a Sneak character! It's constantly trying to flank my team, applying poisons and seeking backstabs. Since teams can only have 1 Sneak, it's a specialized threat we need to account for. The Sneak's movement focus makes it hard to catch!"

**Implementation Priority:** MEDIUM - Required for Phase 3 AI team composition enforcement

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ai_team_composition_rules]]`
- **Related Files:** `upsilonbattle/battlearena/controller/teamcomposition.go`, matchmaking logic
- **Integration:** Works with `ai_controller_archetypes`, `mec_ai_archetype_system`, `ai_progression_matching`

## EXPECTATION
