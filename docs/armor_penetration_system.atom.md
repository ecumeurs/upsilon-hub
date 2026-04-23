---
id: temp
status: DRAFT
version: 1.0
parents: []
dependents: []
---

# New Atom

## INTENT
To implement armor penetration system where certain attacks (like backstabs) can bypass partial armor protection, reducing damage mitigation effectiveness of target's defensive bonuses.

## THE RULE / LOGIC
**Armor Penetration System:**

**Core Principle:**
Some attacks can penetrate (ignore) partial armor rating, making them more effective against heavily armored targets.

**Armor System Basics:**

**Standard Damage Calculation:**
```
Damage = Attack - Defense - ArmorRating - Shield
```

**Armor Penetration Effect:**
```
PenetratedDamage = Attack - (Defense + ArmorRating × Penetration%)
```

**Penetration Mechanic:**

**Backstab Penetration:**
- **Penetration Amount:** 50% of armor rating ignored
- **Calculation:** (ArmorRating × 0.5) is penetrated amount
- **Formula:** Damage = Attack - Defense - (ArmorRating × 0.5) - Shield
- **Example:** Target has 10 ArmorRating, 50% penetration = ignores 5 armor

**Penetration Types:**

**Percentage Penetration:**
- **Backstab:** 50% armor penetration (V2 core mechanic)
- **Critical Hit:** May have penetration bonus (future: critical attacks)
- **Weapon Specific:** Some weapons have penetration properties
- **Skill-Based:** Future skills with "Penetration" property

**Fixed Penetration:**
- **Piercing Attacks:** Ignore N points of armor flat amount
- **Magical Effects:** Magic damage that bypasses physical armor
- **Elemental Bonus:** Fire vs Ice armor, etc.
- **Shield Ignorance:** Some effects bypass shields entirely

**Penetration Formula Implementation:**

**Standard Damage:**
```go
func CalculateStandardDamage(attacker, target Entity) int {
    attackPower := attacker.GetAttackPower()
    defense := target.GetDefense()
    armorRating := target.GetArmorRating()
    shield := target.GetShield()
    
    damage := attackPower - defense - armorRating - shield
    return max(1, damage)  // Minimum 1 damage rule
}
```

**Backstab Damage with Penetration:**
```go
func CalculateBackstabDamage(attacker, target Entity, isBackstab bool) int {
    attackPower := attacker.GetAttackPower()
    defense := target.GetDefense()
    armorRating := target.GetArmorRating()
    shield := target.GetShield()
    
    damage := attackPower - defense
    
    if isBackstab {
        // 150% damage multiplier
        damage = int(float64(damage) * 1.5)
        
        // 50% armor penetration
        penetratedArmor := int(float64(armorRating) * 0.5)
        damage -= penetratedArmor
    }
    
    // Shield applies fully (not penetrated)
    damage -= shield
    
    return max(1, damage)  // Minimum 1 damage rule
}
```

**Armor Penetration Examples:**

**Light Armor (3 ArmorRating):**
- **Standard Attack:** Attack 10 - 3 Defense - 3 Armor = 4 damage
- **Backstab Attack:** Attack 15 (10×1.5) - 3 Defense - 1.5 Armor (50% pen) = 10.5 damage
- **Penetration Bonus:** +6.5 damage from armor penetration

**Heavy Armor (10 ArmorRating):**
- **Standard Attack:** Attack 10 - 3 Defense - 10 Armor = -3 (clamped to 1 minimum)
- **Backstab Attack:** Attack 15 - 3 Defense - 5 Armor (50% of 10) = 7 damage
- **Penetration Bonus:** +4 damage vs 0 damage (enables damage against tanks)

**Shield Interaction:**

**Shield Not Penetrated:**
- **Shield Applies Fully:** Shield value subtracted after penetration
- **Backstab Example:** Attack 15 - 3 Defense - 1.5 Armor (pen) - 10 Shield = 0.5 damage
- **Shield Value:** Shields counter both normal and penetrating attacks equally

**Shield Penetration:**
- **Rare Skills:** Some effects may penetrate shields partially
- **Percentage Penetration:** 25-50% shield penetration (future)
- **Shield Break:** Attacks designed to break shields completely

**Penetration Stackability:**

**Multiple Sources:**
- **Single Penetration:** Only one penetration effect applies per attack
- **Priority Order:** Highest penetration effect wins (if multiple sources)
- **Additive vs Highest:** Percentage penetrations use highest value
- **Different Types:** Percentage and fixed penetration don't stack additively

**Equipment with Penetration:**

**Weapon Properties:**
- **Piercing Weapons:** +25% armor penetration property
- **Daggers:** Enhanced backstab penetration (+10-25%)
- **Magical Weapons:** May have elemental penetration
- **Skill Enhancement:** Some skills add penetration temporarily

**Armor Penetration Items:**
- **Piercing Arrows:** Ammunition that adds penetration
- **Penetration Potions:** Temporary buff granting penetration
- **Enchantments:** Equipment enchantments for penetration
- **Debuffs:** Reduce enemy armor (soft penetration)

**Penetration Balance Considerations:**

**Anti-Tank Purpose:**
- **Heavy Armor Problem:** Heavy tanks ignore low damage attacks
- **Backstab Solution:** 50% penetration makes backstabs effective vs tanks
- **Skill Balance:** Penetration skills useful vs armor builds
- **Counter-Play:** Armor penetration creates rock-paper-scissors dynamics

**Diminishing Returns:**
- **High Armor:** Penetration less effective vs light armor
- **Low Armor:** Over-penetration doesn't cause negative armor
- **Shield Importance:** Shields remain valuable even with penetration
- **Stackability:** Single penetration prevents infinite damage scaling

**Penetration Visual Feedback:**

**Damage Display:**
```
[Attack] 10 damage
[Penetration!] Armor Rating: 5 → 2.5 (50% penetrated)
[Damage] 7.5 damage taken
```

**Attack Animation:**
- **Normal Attack:** Standard weapon swing animation
- **Penetrating Attack:** Enhanced visual effect (piercing sound, spark)
- **Backstab Animation:** Special backstab animation with penetration indicator
- **Color Coding:** Penetrating attacks use different damage number color

**Integration with Other Systems:**

**Backstab Detection:**
- **Trigger:** Backstab detection algorithm enables penetration mechanic
- **Conditional:** Penetration only applies when backstab detected
- **Exclusive:** Standard attacks don't penetrate armor

**Damage Computation Rework:**
- **Weapon-as-Skill:** Weapon attacks use full damage calculation
- **Penetration Integration:** Weapon properties include penetration values
- **Skill vs Basic:** Skills may have separate penetration rules

**Equipment System:**
- **Armor Rating:** Equipment provides armor penetration targets
- **Weapon Properties:** Some weapons add penetration bonuses
- **Stat Bonuses:** Equipment can grant penetration temporarily

**Credit Economy:**
- **Damage Calculation:** Penetration affects damage dealt → credit earning
- **Tank Counterplay:** Penetration enables damage vs high-armor builds
- **Risk/Reward:** High-armor builds vulnerable to penetration attacks

**Performance Optimization:**

**Penetration Cache:**
```go
// Cache penetration calculations
type PenetrationCache struct {
    SourceID       uuid.UUID
    Penetration     float64  // 0.0-1.0 range
    CachedTurn     int
}

func GetEffectivePenetration(source Entity) float64 {
    // Check cache for current turn
    cacheKey := source.ID
    
    if cached, exists := penetrationCache[cacheKey]; exists {
        if cached.CachedTurn == currentTurn {
            return cached.Penetration  // Use cached value
        }
    }
    
    // Calculate effective penetration
    basePenetration := source.GetBasePenetration()
    equipmentBonus := GetEquipmentPenetration(source)
    skillBonus := GetSkillPenetration(source)
    
    effectivePenetration := min(1.0, basePenetration + equipmentBonus + skillBonus)
    
    penetrationCache[cacheKey] = PenetrationCache{
        SourceID:   source.ID,
        Penetration: effectivePenetration,
        CachedTurn: currentTurn,
    }
    
    return effectivePenetration
}
```

**Player Experience:**

**Heavy Armor Encounter:**
"I have 15 ArmorRating from my equipment—standard attacks only deal 1-2 damage to me! However, backstabs penetrate 50% of my armor, so I need to watch my back. When that rogue backstabs me, I take 7 damage instead of 1. Armor penetration makes me vulnerable to backstabs!"

**Backstab Play:**
"I see that tank with 20 ArmorRating—my normal attacks deal only 1 damage to him. But backstabs penetrate 50%, so I only need to get behind him. 10 ArmorRating penetration + my backstab bonus = significant damage vs his tank build. This creates tactical choices: do I focus attack or try to flank?"

**Implementation Priority:** MEDIUM - Required for Phase 2 backstabbing mechanic damage calculation

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[armor_penetration_system]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonbattle/battlearena/entity/entity.go`
- **Integration:** Works with `backstab_detection_algorithm`, `mec_backstabbing_mechanic`, `weapon_damage_computation_rework`

## EXPECTATION
