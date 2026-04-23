---
id: temp
status: DRAFT
parents: []
dependents: []
version: 1.0
---

# New Atom

## INTENT
To implement weapon-based attack computation system where equipped weapons transform basic attacks into skill-based attacks, utilizing weapon properties as skill attributes for full damage calculation.

## THE RULE / LOGIC
**Weapon Damage Computation Rework:**

**Core Transformation:**
Equipped weapons transform basic attacks from simple formula to full skill-based damage computation.

**V1 Basic Attack (Legacy):**
```
Damage = max(1, Attack - Defense)
```
**V2 Weapon-as-Skill Attack:**
```
Damage = Full Skill Computation using Weapon Properties
```

**Weapon Property to Skill Property Mapping:**

**Weapon Properties:**
- **WeaponBaseDamage:** Base damage value (equivalent to Attack stat)
- **WeaponRange:** Attack range (overrides default range)
- **WeaponType:** Melee vs Ranged classification
- **CritChance:** Critical hit chance bonus (0-25%)
- **CritMultiplier:** Critical damage multiplier (100-200%)
- **BackstabBonus:** Enhanced backstab damage (+10-25%)
- **AttackSpeed:** Cooldown modifier (faster/slower attacks)

**Skill Properties (from Weapons):**
- **Damage:** Direct mapping from WeaponBaseDamage
- **Range:** Direct mapping from WeaponRange
- **Accuracy:** Derived from weapon type (melee = 100%, ranged = 80-95%)
- **CriticalChance:** Weapon CritChance + Character CritChance
- **CriticalMultiplier:** Weapon CritMultiplier
- **Zone:** Single cell (weapons don't have AoE)
- **TargetType:** Entity or EntityOrTile

**Weapon Attack Flow:**

**1. Attack Initiation:**
- **Player Action:** Player uses basic attack (not skill)
- **Weapon Check:** Character has equipped weapon?
- **Unarmed Fallback:** No weapon = V1 basic attack formula

**2. Skill Generation:**
```go
func CreateWeaponSkill(character Character, weapon Weapon) skill.Skill {
    // Transform weapon properties into skill
    weaponSkill := skill.Skill{
        Name:        weapon.Name + " Attack",
        Behavior:    def.Direct,      // Active skill
        Targeting: map[property.Property]property.Property{
            property.Range:    MakeIntProperty("Range", weapon.WeaponRange),
            property.Zone:     MakeIntProperty("Zone", 1), // Single target
            property.TargetType: MakeEnumProperty("TargetType", "Entity"),
            property.Accuracy: GetWeaponAccuracy(weapon),
        },
        Costs: map[property.Property]property.Property{
            property.Delay:   MakeIntProperty("Delay", GetWeaponAttackSpeed(weapon)),
        },
        Effect: effect.Effect{
            Properties: []property.Property{
                MakeIntProperty("Damage", weapon.WeaponBaseDamage),
                MakeIntProperty("CriticalChance", weapon.CritChance + character.CritChance),
                MakeIntProperty("CriticalMultiplier", weapon.CritMultiplier),
            },
        },
    }
    
    return weaponSkill
}
```

**3. Damage Calculation:**
```go
func CalculateWeaponAttack(attacker, target Entity) int {
    // Check for equipped weapon
    weapon := attacker.GetEquippedWeapon()
    
    if weapon == nil {
        // Unarmed: use V1 basic attack
        return CalculateBasicAttack(attacker, target)
    }
    
    // Weapon-as-Skill: create temporary skill
    weaponSkill := CreateWeaponSkill(attacker, weapon)
    
    // Use full skill damage computation
    damage := CalculateSkillAttack(weaponSkill, target)
    
    return damage
}
```

**Weapon Type Integration:**

**Melee Weapons (Swords, Axes, Daggers):**
- **Range:** Usually 1-2 cells
- **Accuracy:** 100% (melee default)
- **Attack Speed:** Standard cooldown (500 delay)
- **Damage:** Higher than ranged (close range advantage)
- **Examples:**
  - Sword: +5 Damage, Range 1, 100% Accuracy
  - Dagger: +3 Damage, Range 1, +25% Crit, +25% Backstab

**Ranged Weapons (Bows, Pistols):**
- **Range:** 4-7 cells
- **Accuracy:** 80-95% (distance penalty)
- **Attack Speed:** Standard cooldown (500 delay)
- **Damage:** Lower than melee (range advantage)
- **Examples:**
  - Bow: +3 Damage, Range 5, 90% Accuracy
  - Pistol: +2 Damage, Range 4, 85% Accuracy

**Two-Handed Weapons:**
- **Damage Bonus:** +2-4 additional damage (both hands committed)
- **Utility Slot Penalty:** Cannot equip utility items while using
- **Attack Speed:** Slightly slower (600 delay)
- **Examples:**
  - Two-Handed Sword: +7 Damage, Range 1, blocks utility
  - Greatsword: +6 Damage, Range 1, blocks utility

**Damage Calculation Examples:**

**V1 Unarmed Attack:**
```
Character: Attack 10, Enemy: Defense 5
Damage = max(1, 10 - 5) = 5 damage taken
```

**V2 Sword Attack:**
```
Weapon: +5 Damage, Range 1
Character: Attack 10 (base) + 5 (weapon) = 15 total
Enemy: Defense 5
Weapon Skill: Damage 15, Accuracy 100%, Crit 0%

Attack Calculation:
Base Damage: 15
Apply Defense: 15 - 5 = 10 damage
Result: Enemy takes 10 damage
```

**V2 Bow Attack:**
```
Weapon: +3 Damage, Range 5, 90% Accuracy
Character: Attack 10 + 3 = 13 total
Enemy: Defense 5 (3 cells away)
Weapon Skill: Damage 13, Accuracy 90%, Crit 5%

Attack Calculation:
Base Damage: 13
Hit Test: 90% hit chance = Hit!
Apply Defense: 13 - 5 = 8 damage
Critical Check: 5% crit chance = No crit
Result: Enemy takes 8 damage (vs 5 from unarmed)
```

**Backstab with Weapon:**
```
Dagger: +3 Damage, +25% Backstab Bonus
Character: Attack 10 + 3 = 13 total
Backstab Detected: 150% damage + Backstab Bonus
Enemy: Defense 5, ArmorRating 3

Attack Calculation:
Base Damage: 13 × 1.5 (backstab) = 19.5 (round to 19)
Apply Defense: 19 - 5 = 14 damage
Apply Backstab Bonus: 14 × 1.25 (weapon bonus) = 17.5 damage
Apply Armor Penetration: ArmorRating 3 × 0.5 = 1.5 penetrated
Final Damage: 17.5 - 1.5 (armor pen) = 16 damage
Result: Enemy takes 16 damage vs 5 from basic attack!
```

**Weapon Skill Integration:**

**Full Skill System Compatibility:**
- **Cooldown Tracking:** Weapon attacks enter cooldown after use
- **Hit Testing:** Use skill accuracy system
- **Crit System:** Weapon CritChance + Character CritChance
- **Defense Calculation:** Standard skill damage formula applies
- **Special Effects:** Some weapons add skill-like effects

**Skill Property Combinations:**

**Skill + Weapon Synergy:**
- **Weapon Base Damage:** Modified by skill damage multipliers
- **Weapon Range:** Modified by skill range properties
- **Weapon Crit:** Enhanced by skill critical bonuses
- **Future V2.2:** Skills that enhance equipped weapons

**Elemental Weapons:**
- **Fire Damage:** Extra damage vs ice armor
- **Ice Damage:** Slows movement of targets
- **Poison Damage:** Additional damage over time
- **Holy Damage:** Bonus vs undead (future content)

**UI Integration:**

**Attack Preview:**
- **Weapon Display:** Show equipped weapon and properties
- **Damage Preview:** Calculate estimated damage before attack
- **Range Indicator:** Highlight attack range on grid
- **Crit Chance Display:** Show combined character + Weapon crit chance

**Combat Log:**
```
[Sword Attack] 15 damage (Weapon +5)
[Hit] Attack connects
[Damage] 10 damage to Enemy (after 5 Defense)
```

**Equipment Bonuses Integration:**

**Stat Additions:**
- **Attack Power:** Weapon WeaponBaseDamage adds to character Attack
- **Crit Chance:** Weapon CritChance + Character CritChance
- **Attack Range:** Weapon WeaponRange overrides default range
- **Backstab Bonus:** Weapon BackstabBonus applies on backstab

**Stat Recalculation:**
```go
func RecalculateCharacterStats(character Character) {
    baseStats := character.BaseStats
    
    // Add equipment bonuses
    for _, item := range character.EquippedItems {
        weapon, ok := item.(*Weapon)
        if ok {
            baseStats.Attack += weapon.WeaponBaseDamage
            baseStats.CritChance += weapon.CritChance
        }
    }
    
    character.TotalStats = CalculateEffectiveStats(baseStats, character.ExoticStats)
}
```

**Implementation Benefits:**

**Unified Damage System:**
- **Single Formula:** Weapon attacks use same damage computation as skills
- **Balance:** Weapon properties follow same SW balance system
- **Progression:** Weapon upgrades provide same progression as skills
- **Simplification:** Remove separate basic attack vs skill attack logic

**Tactical Variety:**
- **Weapon Choice:** Players choose weapons based on playstyle
- **Range Tactics:** Ranged vs melee positioning
- **Stat Synergy:** Weapons complement character stat builds
- **Skill Integration:** Future skills can enhance weapons

**Performance Optimization:**

**Weapon Skill Caching:**
```go
// Cache weapon-as-skill generation
type WeaponSkillCache struct {
    WeaponID    uuid.UUID
    Skill      skill.Skill
    CharacterID uuid.UUID
    CachedTurn  int
}

var weaponSkillCache = make(map[string]WeaponSkillCache)

func GetWeaponSkill(character Character, weapon Weapon) skill.Skill {
    cacheKey := fmt.Sprintf("%s-%s", character.ID, weapon.ID)
    
    if cached, exists := weaponSkillCache[cacheKey]; exists {
        if cached.CachedTurn == currentTurn {
            return cached.Skill  // Use cached skill
        }
    }
    
    skill := CreateWeaponSkill(character, weapon)
    
    weaponSkillCache[cacheKey] = WeaponSkillCache{
        WeaponID:    weapon.ID,
        Skill:      skill,
        CharacterID: character.ID,
        CachedTurn:  currentTurn,
    }
    
    return skill
}
```

**Player Experience:**

**Weapon Upgrade Journey:**
"I start with basic sword (+5 Damage). As I earn credits, I upgrade to Steel Sword (+8 Damage). Now my basic attacks deal significantly more damage—I didn't need to rely on skills for damage output. But I can also use skills for special effects. The weapon-as-skill system makes my basic attacks powerful and consistent!"

**Ranged Tactics:**
"I equip bow (+3 Damage, Range 5). While melee characters have higher damage, I can attack from safety. When I upgrade to Crossbow (+5 Damage, Range 7), I can control even more battlefield. The weapon system makes my basic attacks tactical choices, not just damage calculations!"

**Implementation Priority:** HIGH - Required for Phase 2 combat mechanics unification

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[weapon_damage_computation_rework]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonbattle/battlearena/entity/skill/skill.go`
- **Integration:** Works with `mec_weapon_as_skill_system`, `armor_penetration_system`, `backstab_detection_algorithm`

## EXPECTATION
