---
id: temp
status: DRAFT
version: 1.0
parents: []
dependents: []
---

# New Atom

## INTENT
To implement extended character sheet UI integration that displays all V2 properties, skills, equipment, and progression information, enabling players to understand and manage their complete character capabilities.

## THE RULE / LOGIC
**Extended Character Sheet Integration:**

**Core Principle:**
Character sheet displays comprehensive character information including V2 stats, skills, equipment, credits, and progression status.

**V1 vs V2 Character Sheet:**

**V1 Character Sheet (Legacy):**
```
Name: Alice
Level: 5
Wins: 12

Stats:
- HP: 7/7
- Attack: 2
- Defense: 1
- Movement: 1
```

**V2 Character Sheet (Extended):**
```
Name: Alice
Level: 5 (CP: 100 + 12×10 = 220)
Wins: 12
Credits: 1250

Core Stats:
- HP: 45/50 (Base 30 + 15 CP)
- Attack: 15 (Base 10 + 5 CP)  
- Defense: 8 (Base 5 + 3 CP)
- Movement: 3 (Base 3 + 0 CP)

Exotic Stats:
- Critical Chance: 8% (Base 0% + 8 CP, cost 80 CP)
- Critical Multiplier: 125% (Base 100% + 5×5% CP, cost 25 CP)
- Dodge: 5% (Base 0% + 5 CP, cost 25 CP)
- Accuracy: 100% (Base 100% + 0 CP)
- Jump Height: 3 (Base 2 + 1 CP, cost 15 CP)

Equipment:
- Armor Slot: Steel Chest (+5 Defense, +8 ArmorRating)
- Utility Slot: Ring of Accuracy (+5% CritChance, +10% CritMultiplier)
- Weapon Slot: Steel Sword (+5 Damage, Range 2)

Skills:
- Fireball (Grade III, 250 SW, 500 credits)
- Heal (Grade II, 120 SW, 240 credits)
- Shield Ally (Grade I, 80 SW, 160 credits)
```

**Character Sheet Sections:**

**1. Identity Section**
- **Character Name:** Display name and optional title
- **Level & Progression:** Current level, total wins, CP pool
- **Class/Archetype:** Player-defined role (future V2.2+)
- **Avatar Display:** Character visual representation

**2. Core Statistics Section**
- **HP Bar:** Current HP / Maximum HP with health percentage
- **Attack Power:** Total Attack value (base + equipment + buffs)
- **Defense Rating:** Total Defense value (base + equipment + buffs)
- **Movement Points:** Current Movement / Maximum Movement
- **Combat Rating:** Calculated threat level (ATK × DEF × MOV)

**3. Exotic Statistics Section**
- **Critical Chance:** Current % chance, equipment bonuses included
- **Critical Multiplier:** Current multiplier for critical hits
- **Dodge Chance:** Current % chance to evade attacks
- **Accuracy:** Current % chance to hit enemies
- **Jump Height:** Maximum tiles can jump vertically

**4. Combat Capabilities Section**
- **Attack Range:** Effective range including equipment bonuses
- **Backstab Bonus:** Enhanced backstab damage multiplier
- **Armor Penetration:** Current penetration value
- **Skill Availability:** Number and grade of owned skills
- **Cooldown Status:** Current skill cooldowns remaining

**5. Equipment Section**
- **Slot Display:** Visual representation of 3 equipment slots
- **Slot Details:**
  - **Armor Slot:** Equipped item, ArmorRating, special effects
  - **Utility Slot:** Equipped item, stat bonuses, special effects
  - **Weapon Slot:** Equipped item, damage, range, crit bonuses
- **Inventory List:** Owned but unequipped items
- **Equipment Stats:** Total stat bonuses from all equipped items

**6. Skill Inventory Section**
- **Skill List:** All owned skills with current cooldowns
- **Skill Details:** Click for expanded view showing:
  - Name and Grade
  - Damage/Healing values
  - Range and targeting
  - Costs (MP, SP, delay, cooldown)
  - Effects and duration
- **Skill Usage:** Historical statistics (times used, total damage dealt)

**7. Economy Section**
- **Credits Balance:** Current available credits
- **Skill Cost:** Total credits spent on skills
- **Equipment Cost:** Total credits spent on equipment
- **Earning Rate:** Average credits per match
- **Purchases:** Recent transactions with costs

**8. Progression Section**
- **CP Pool:** Available points for stat allocation
- **Stat History:** Changes made to base stats
- **Next Milestone:** Progress toward skill selection/reforging
- **Level Progression:** XP/points until next level
- **Unlocked Features:** Content unlocked at current level

**Data Structure:**
```go
type ExtendedCharacterSheet struct {
    // Identity
    CharacterID    uuid.UUID
    Name           string
    Level          int
    Wins           int
    Credits        int
    
    // Core Stats
    CoreStats      CharacterStats
    MaxHP          int
    CurrentHP      int
    CurrentMovement int
    MaxMovement    int
    
    // Exotic Stats
    ExoticStats    ExoticStats
    
    // Equipment
    EquipmentSlots  EquipmentSlots
    Inventory       []Equipment
    
    // Skills
    Skills         []Skill
    SkillHistory    []SkillUsage
    
    // Progression
    AvailableCP    int
    StatChanges    []StatChange
    NextLevelProgress float64
}
```

**Stat Display Logic:**
```go
func DisplayStatBlock(stats CharacterStats, equipment Equipment) StatBlock {
    // Calculate base stats
    baseAttack := stats.Attack
    baseDefense := stats.Defense
    baseMovement := stats.Movement
    
    // Add equipment bonuses
    attack := baseAttack + equipment.WeaponDamage
    defense := baseDefense + equipment.ArmorRating
    movement := baseMovement + equipment.MovementBonus
    
    // Add exotic bonuses
    critChance := stats.CritChance + equipment.CritBonus
    critMult := stats.CritMultiplier + equipment.CritMultBonus
    dodge := stats.Dodge + equipment.DodgeBonus
    
    return StatBlock{
        Attack:          attack,
        Defense:         defense,
        Movement:        movement,
        CritChance:      critChance,
        CritMultiplier:   critMult,
        Dodge:           dodge,
    }
}
```

**Skill Card Display:**
```go
func DisplaySkillCard(skill Skill, currentCooldown int) SkillCard {
    // Calculate skill power
    power := CalculateSkillWeight(skill)
    grade := DetermineGrade(power)
    
    // Format skill information
    card := SkillCard{
        Name:        skill.Name,
        Grade:       grade,
        Power:       power,
        
        // Costs
        Delay:       skill.Costs.Delay,
        MP:          skill.Costs.MP,
        SP:          skill.Costs.SP,
        Cooldown:     skill.Cooldown,
        Remaining:    currentCooldown,
        
        // Effects
        Damage:      GetDamageValue(skill),
        Range:        GetRangeValue(skill),
        Target:       GetTargetType(skill),
        
        // Visual presentation
        Icon:         GetSkillIcon(grade),
        Color:        GetSkillColor(skill.Type),
    }
    
    return card
}
```

**Equipment Slot Visualization:**

**3-Slot Layout:**
```
┌─────────────────────────────────┐
│  ARMOR     │  UTILITY   │  WEAPON    │
│             │             │             │
│ [Steel      │ [Ring of    │ [Steel      │
│  Chest]     │  Accuracy]   │  Sword]     │
│ +5 Defense  │ +5% Crit    │ +5 Damage    │
│ +8 Armor    │ +10% Crit    │ Range: 2     │
└─────────────────────────────────┘
```

**Slot State Indicators:**
- **Empty Slot:** Grayed out, "Empty" label
- **Equipped Item:** Colored by tier (Common=gray, Uncommon=green, Rare=blue)
- **Stat Bonus Preview:** Hover shows stat changes if equipped
- **Unequip Button:** Remove current item, return to inventory

**Progression Visualization:**

**CP Allocation Interface:**
```
Available CP: 50
┌─────────────────────────────────┐
│ HP (+1 CP, Cost: 1)           │
│ Attack (+5 CP, Cost: 5)         │
│ Defense (+5 CP, Cost: 5)        │
│ Movement (+30 CP, Cost: 30)      │
│ CritChance (+10 CP, Cost: 10)     │
│ Dodge (+5 CP, Cost: 5)           │
└─────────────────────────────────┘

Current Stats: HP 45/50, Attack 15, Defense 8, Movement 3
```

**Level Progress Tracking:**
```
Level 5 → Level 6 (500 XP required)
Progress: [████████░░░░░] 60%

Next Level Rewards:
- +10 CP (total: 60)
- Skill Selection: Choose 1 of 3 Grade II-III skills
```

**Interactive Features:**

**Stat Hover Tooltips:**
- **Stat Explanation:** Mouse hover shows what stat affects gameplay
- **Comparison Tool:** "Current: 8 | With Item: 12 (+4)"
- **Level Restriction:** "Available at Level 10" for exotic stats

**Skill Detail Modal:**
- **Full Stats:** Click skill to see complete breakdown
- **Targeting Preview:** Show range, AoE pattern on grid
- **Effect Preview:** Display exact damage/healing formula
- **Cost Analysis:** Resource investment vs damage output

**Equipment Comparison:**
```
Current: Wooden Sword (+2 Damage, Range 1)
New: Steel Sword (+5 Damage, Range 2)

Difference: +3 Damage, +1 Range
Upgrade Cost: 200 credits
```

**Performance Optimization:**

**Lazy Loading:**
```go
// Load character data progressively
func LoadExtendedCharacterSheet(characterID uuid.UUID) ExtendedCharacterSheet {
    // First: Load identity and core stats
    basic := LoadBasicData(characterID)
    
    // Second: Load equipment in background
    equipmentChan := make(chan []Equipment)
    go LoadEquipmentAsync(characterID, equipmentChan)
    
    // Third: Load skills with cooldowns
    skillsChan := make(chan []Skill)
    go LoadSkillsAsync(characterID, skillsChan)
    
    // Combine results
    return ExtendedCharacterSheet{
        Basic:    basic,
        Equipment: <-equipmentChan,
        Skills:    <-skillsChan,
    }
}
```

**State Caching:**
```go
// Cache frequently accessed data
type CharacterSheetCache struct {
    CharacterID    uuid.UUID
    Stats          CharacterStats
    Equipment      EquipmentSlots
    Skills         []Skill
    LastUpdated     time.Time
}

var sheetCache = make(map[string]CharacterSheetCache)

func GetCachedSheet(characterID uuid.UUID) ExtendedCharacterSheet {
    cacheKey := characterID.String()
    
    if cached, exists := sheetCache[cacheKey]; exists {
        if time.Since(cached.LastUpdated) < 5*time.Minute {
            return cached.Sheet  // Use cached data
        }
    }
    
    // Load and cache fresh data
    sheet := LoadExtendedCharacterSheet(characterID)
    sheetCache[cacheKey] = CharacterSheetCache{
        CharacterID: characterID,
        Sheet:       sheet,
        LastUpdated: time.Now(),
    }
    
    return sheet
}
```

**Player Experience:**

**Character Overview:**
"The extended character sheet shows everything about my character! I can see my V2 stats (HP 45, Attack 15), exotic stats (8% CritChance, 125% CritMultiplier), equipment (Steel Sword, Ring of Accuracy), and skills (Fireball, Heal, Shield). The CP pool shows I have 50 points left to spend!"

**Stat Planning:**
"I hover over 'Critical Chance' and it explains 'Increases chance to deal double damage'. The stat comparison shows that equipping the Ring would give me 13% instead of my current 8%. But it costs 400 credits and I only have 150. I need to save up!"

**Progression Tracking:**
"I'm Level 5 and can see I'm 60% of the way to Level 6. At Level 10, I'll get +10 CP and can choose a new skill. The progression visualization makes it clear what milestones are coming and when I'll unlock new features!"

**Tactical View:**
"Looking at my skills, I can see Fireball has 400 delay channeling, perfect for big damage but risky. Heal has no delay and targets allies. Shield Ally protects teammates but doesn't deal damage. The skill cards show me cooldowns, costs, and even SW for balance. This helps me plan my turn strategy!"

**Implementation Priority:** MEDIUM - Required for Phase 3-4 UI integration and player understanding

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[extended_character_sheet_integration]]`
- **Related Files:** Character sheet UI components, stat calculation APIs, progression tracking systems
- **Integration:** Works with `character_point_buy_system`, `entity_equipment_system`, `api_equipment_management`

## EXPECTATION
