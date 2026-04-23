---
id: temp
status: DRAFT
version: 1.0
parents: []
dependents: []
---

# New Atom

## INTENT
To implement AI progression matching system where AI characters follow the same point-buy progression rules as players, scaling stats, skills, and level according to player team averages.

## THE RULE / LOGIC
**AI Progression Matching:**

**Core Principle:**
AI characters follow identical progression mechanics to players, ensuring fair and challenging matches.

**Progression Synchronization:**

**Starting Characters:**
- **V2 Baseline:** HP 30-50, Attack 10, Defense 5, Movement 3
- **Point Allocation:** 100 CP to spend during AI creation
- **Stat Costs:** Same CP costs as players (HP=1, Attack=5, Defense=5, Movement=30)
- **Exotic Stats:** Same costs and restrictions (CritChance=10 CP/5 levels, etc.)

**Win Rewards:**
- **Same Formula:** +10 CP per win (identical to players)
- **Cumulative Cap:** 100 + (total_wins × 10) CP maximum
- **Skill Selection:** Every 10 levels, choose 1 of 3 skills
- **Skill Reforging:** Every 5 levels, can modify existing skills

**Level Matching:**

**Player Level Calculation:**
```go
func CalculatePlayerAverageLevel(players []Character) int {
    if len(players) == 0 {
        return 1  // Default Level 1 if no players
    }
    
    totalLevel := 0
    for _, player := range players {
        totalLevel += player.Level
    }
    
    average := totalLevel / len(players)
    return RoundToNearestInt(average)
}
```

**AI Level Assignment:**
```go
func AssignAILevels(teamSize int, playerLevel int) []int {
    levels := make([]int, teamSize)
    
    // Option A: Exact match to average
    for i := 0; i < teamSize; i++ {
        levels[i] = playerLevel
    }
    
    // Option B: Slight variation (±2 levels)
    for i := 0; i < teamSize; i++ {
        variation := rand.Intn(3) - 1  // -1 to +1
        levels[i] = playerLevel + variation
    }
    
    // Option C: Progressive matching (1v1=exact, 2v2=slight spread)
    if teamSize == 1 {
        levels[0] = playerLevel  // Exact match for 1v1
    } else if teamSize == 2 {
        levels[0] = playerLevel    // Match player level
        levels[1] = playerLevel + 1  // Slight advantage
    } else {
        // Spread levels around player average
        spread := CalculateLevelSpread(playerLevel, teamSize)
        for i := 0; i < teamSize; i++ {
            levels[i] = playerLevel + spread[i]
        }
    }
    
    return levels
}
```

**Skill Grade Access:**

**Player vs AI Skill Availability:**
- **Identical Access:** Both players and AI get same skill grades at same level
- **Grade Thresholds:** Same level ranges for skill availability
- **Pool Filtering:** AI selects from appropriate skill pools

**Grade Progression Levels:**
```
Level 1-9:     Grade I-II skills available
Level 10-19:    Grade II-III skills available  
Level 20-29:    Grade III-IV skills available
Level 30+:       Grade IV-V skills available
```

**AI Skill Selection:**
```go
func SelectAISkills(ai AIController) []skill.Skill {
    level := ai.GetLevel()
    gradeRange := GetGradeRange(level)
    
    // Filter available skills by grade
    availableSkills := GetSkillsInGradeRange(gradeRange)
    
    // Filter by archetype
    archetypeSkills := FilterByArchetype(availableSkills, ai.GetArchetype())
    
    // Select 1 of 3 random skills
    selected := SelectRandomSkills(archetypeSkills, 3)
    
    return selected
}
```

**Stat Allocation Matching:**

**AI Build Templates:**
```go
func GenerateAIStats(level int, archetype string) CharacterStats {
    baseStats := GetV2BaselineStats()  // HP 30-50, Attack 10, Defense 5, Movement 3
    availableCP := 100 + (level * 10)  // Same formula as players
    
    switch archetype {
        case "Fighter":
            return AllocateFighterStats(availableCP, baseStats)
            
        case "Ranger":
            return AllocateRangerStats(availableCP, baseStats)
            
        case "Support":
            return AllocateSupportStats(availableCP, baseStats)
            
        case "Sneak":
            return AllocateSneakStats(availableCP, baseStats)
    }
}
```

**Archetype Stat Priorities:**

**Fighter Allocation:**
- **Attack:** 60% of points (60 CP from 100)
- **Defense:** 25% of points (25 CP)
- **HP:** 10% of points (10 CP)
- **Movement:** 5% of points (5 CP)
- **Rationale:** High damage, frontline survivability

**Ranger Allocation:**
- **Attack:** 50% of points (50 CP)
- **Accuracy:** 20% of points (20 CP)
- **Movement:** 15% of points (15 CP)
- **Defense:** 10% of points (10 CP)
- **HP:** 5% of points (5 CP)
- **Rationale:** Precision, range, mobility

**Support Allocation:**
- **MP/SP:** 40% of points (40 CP)
- **Defense:** 25% of points (25 CP)
- **HP:** 20% of points (20 CP)
- **Healing Power:** 15% of points (15 CP)
- **Movement:** 0% of points (0 CP)
- **Rationale:** Support focus, defensive stats

**Sneak Allocation:**
- **Movement:** 30% of points (30 CP)
- **Attack:** 25% of points (25 CP)
- **Dodge:** 20% of points (20 CP)
- **Critical Chance:** 15% of points (15 CP)
- **HP/Defense:** 10% of points (10 CP)
- **Rationale:** Positional advantage, evasion, burst damage

**Skill Level Scaling:**

**AI Skill Progression:**
```go
func ScaleAISkillToLevel(skill skill.Skill, targetLevel int) skill.Skill {
    baseSW := CalculateSkillWeight(skill)
    
    // Scale skill power with level
    levelMultiplier := 1.0 + (targetLevel * 0.05)  // +5% per level
    scaledSW := int(float64(baseSW) * levelMultiplier)
    
    // Adjust skill properties
    scaledSkill := skill.DeepCopy()
    scaledSkill.ScaleProperties(scaledSW)
    
    return scaledSkill
}
```

**Difficulty Balancing:**

**Win Rate Adjustment:**
- **Player Win Rate > 50%:** Increase AI difficulty slightly
- **Player Win Rate < 50%:** Decrease AI difficulty slightly
- **Adjustment Range:** ±2 CP from standard allocation
- **Progressive Scaling:** Difficulty adjustments accumulate over matches

**AI Personalities:**
```go
type AIPersonality struct {
    Aggression    float64  // 0.0-1.0 (passive to aggressive)
    RiskTolerance float64  // 0.0-1.0 (safe to reckless)
    Adaptiveness   float64  // 0.0-1.0 (predictable to adaptive)
    Coordination  float64  // 0.0-1.0 (individual to team-focused)
}

func GenerateAIPersonality(difficulty string) AIPersonality {
    switch difficulty {
        case "Easy":
            return AIPersonality{Aggression: 0.3, RiskTolerance: 0.8, Adaptiveness: 0.2}
            
        case "Normal":
            return AIPersonality{Aggression: 0.5, RiskTolerance: 0.5, Adaptiveness: 0.5}
            
        case "Hard":
            return AIPersonality{Aggression: 0.7, RiskTolerance: 0.3, Adaptiveness: 0.7}
    }
}
```

**Level Difference Handling:**

**Matchmaking Considerations:**
- **Level Advantage:** AI should not have significant level advantage
- **Level Disadvantage:** AI should not be severely underpowered
- **Fair Play:** Keep level difference within ±2 for balanced matches
- **Ranked Considerations:** Higher-ranked players face higher-level AI

**Progression Tracking:**
```go
type AIProgressionData struct {
    CharacterID     uuid.UUID
    TotalWins       int
    CurrentLevel    int
    AvailableCP     int
    SkillGrades     []string
    Archetype       string
}

func TrackAIProgression(ai AIController) {
    // Record wins and CP accumulation
    if ai.HasWonMatch() {
        ai.ProgressionData.TotalWins++
        ai.ProgressionData.AvailableCP += 10  // Same reward as players
        
        // Check for skill selection/reforging
        if ai.ProgressionData.CurrentLevel >= 10 {
            GrantSkillSelection(ai)
        }
        
        if ai.ProgressionData.CurrentLevel % 5 == 0 {
            EnableSkillReforging(ai)
        }
    }
    
    // Update level and stats
    ai.RecalculateStats()
}
```

**Fair Play Enforcement:**

**No Cheating:**
- **Same Rules:** AI follows exact same progression as players
- **No Advantages:** No hidden AI bonuses beyond CP system
- **Transparency:** Players can understand AI progression mechanics
- **Balanced Challenge:** AI provides appropriate difficulty at each level

**Player Experience:**

**Fair Match:**
"The AI team seems perfectly matched to our level! They have similar stats and skill grades as us. I notice the Fighter AI has high Attack and Defense like a typical fighter build. The progression matching feels fair!"

**Skill Advantage:**
"We've been winning a lot, so the AI team is now higher level than us. Their skills are Grade III while we're still using Grade II skills. The progression system is working—they earned their +10 CP per win just like us!"

**Stat Distribution:**
"The AI Support has higher MP than our healer! This makes their shielding skills more frequent. But I notice they've spent less points on Attack, so they deal less damage directly. The stat allocation rules create interesting trade-offs!"

**Implementation Priority:** HIGH - Required for Phase 3 AI enhancement

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ai_progression_matching]]`
- **Related Files:** AI creation logic, progression tracking, matchmaking integration
- **Integration:** Works with `mec_ai_archetype_system`, `ai_controller_archetypes`, `character_stat_allocation_rules`

## EXPECTATION
