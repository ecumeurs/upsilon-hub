---
id: mechanic_backstab_detection_algorithm
status: DRAFT
dependents: []
type: MECHANIC
layer: IMPLEMENTATION
priority: 5
version: 2.0
parents: []
---

# New Atom

## INTENT
To implement backstab detection algorithm that determines when an attack originates from behind the target based on entity orientation, enabling 150% damage multiplier and 50% armor penetration bonuses.

## THE RULE / LOGIC
**Backstab Detection Algorithm:**

**Core Principle:**
Backstab occurs when attacker is positioned behind target relative to target's facing direction.

**Orientation System:**
```go
type EntityOrientation int

const (
    Up    EntityOrientation = 0
    Right  EntityOrientation = 1
    Down  EntityOrientation = 2
    Left   EntityOrientation = 3
)
```

**Backstab Detection Logic:**

**1. Calculate Attack Vector:**
- **Attacker Position:** Get attacker's coordinates
- **Target Position:** Get target's coordinates
- **Attack Direction:** Determine which direction attack comes from

**2. Determine Target Facing:**
- **Target Orientation:** Get target's current facing direction
- **Back Definition:** Back is opposite of target's facing
- **Angle Calculation:** Calculate angle between attack vector and target facing

**3. Apply Backstab Criteria:**
- **Positional Check:** Is attacker within 45° of target's back?
- **Angle Tolerance:** ±45° from perfect back angle
- **Distance Check:** Is attacker within weapon range?
- **Obstruction Check:** Is line of sight clear to target?

**Orientation-to-Back Mapping:**
```
Target Facing | Back Direction | Attack Angle Range for Backstab
    Up (0)      | Down (2)        | 135°-225°
    Right (1)     | Left (3)         | 225°-315°
    Down (2)      | Up (0)          | 315°-45°
    Left (3)       | Right (1)        | 45°-135°
```

**Backstab Detection Algorithm:**
```go
func IsBackstab(attacker, target Entity) bool {
    // Calculate attack angle from attacker to target
    attackAngle := CalculateAttackAngle(attacker.Position, target.Position)
    
    // Get target's back angle range
    backAngleMin, backAngleMax := GetBackAngleRange(target.Orientation)
    
    // Check if attack angle falls within back range
    if attackAngle >= backAngleMin && attackAngle <= backAngleMax {
        // Additional validation checks
        if IsWithinRange(attacker, target) && HasLineOfSight(attacker, target) {
            return true  // Valid backstab
        }
    }
    
    return false  // Not a backstab
}

func CalculateAttackAngle(from Position, to Position) float64 {
    // Calculate angle in degrees (0-360)
    deltaX := to.X - from.X
    deltaY := to.Y - from.Y
    angle := math.Atan2(deltaY, deltaX) * (180 / math.Pi)
    
    // Normalize to 0-360 range
    if angle < 0 { angle += 360 }
    return angle
}

func GetBackAngleRange(orientation EntityOrientation) (float64, float64) {
    switch orientation {
        case Up:    return 135, 225    // 135°-225° range
        case Right:  return 225, 315    // 225°-315° range
        case Down:  return 315, 45     // 315°-45° range
        case Left:   return 45, 135     // 45°-135° range
    }
}
```

**Edge Cases and Validation:**

**Adjacent Only:**
- **Range Check:** Backstab only works at standard weapon range (usually melee)
- **No Ranged Backstabs:** Long-range attacks don't benefit from backstab
- **Range Limit:** Backstab effective up to 2-3 cells max

**Line of Sight Required:**
- **Obstruction Check:** Walls/obstacles block backstab
- **Visibility:** Attacker must see target's back
- **Transparent Effects:** Some zones don't block backstab line of sight

**Sidestab Consideration:**
- **Partial Bonus:** Attacking from side (90° from facing) could give reduced bonus
- **Current V2:** No sidestab implemented yet (future feature)
- **Angle Threshold:** Perfect backstab vs partial backstab detection

**Multiple Attackers:**
- **Individual Calculation:** Each attacker calculated independently
- **Backstab Stacking:** Multiple attackers can backstab same target simultaneously
- **Credit Assignment:** Each backstab attacker earns full damage credits

**Target Orientation Updates:**
- **Face Toward Attacker:** When taking damage, target may auto-face attacker
- **Persistent Orientation:** Some creatures don't auto-face (stationary)
- **Exploitable Behavior:** Players can "bait" targets to turn them around

**Backstab Visual Feedback:**

**Combat Log:**
```
[Backstab!] Alice attacks Bob from behind (+50% damage, armor penetration!)
[Damage] 12 damage (8 base × 1.5) - 2 armor penetration = 10 damage taken
```

**Visual Indicators:**
- **Backstab Icon:** Special damage number with backstab symbol
- **Camera Flash:** Brief screen flash on backstab
- **Damage Animation:** Enhanced animation for backstab damage
- **Sound Effect:** Distinct backstab sound vs normal attack

**Integration with Damage System:**

**Backstab Damage Formula:**
```go
func CalculateBackstabDamage(attacker, target Entity, baseDamage int) int {
    damage := baseDamage
    
    if IsBackstab(attacker, target) {
        // 150% damage multiplier
        damage = int(float64(damage) * 1.5)
        
        // 50% armor penetration
        targetArmor := target.GetArmorRating()
        penetratedArmor := int(float64(targetArmor) * 0.5)
        
        damage -= penetratedArmor
    }
    
    // Apply shield (not penetrated)
    damage -= target.Shield
    
    return max(1, damage)  // Minimum 1 damage rule
}
```

**Scope Limitations:**

**Weapon Attacks Only:**
- **Basic Attack:** Weapon-based attacks use backstab detection
- **Skill Attacks:** Skills do NOT benefit from backstab (V2.1 limitation)
- **Future Enhancement:** Skills with "Backstab Enabled" flag (V2.2+)
- **Skill Exemption:** Magical attacks, ranged skills skip backstab logic

**All Weapons Support:**
- **Melee Weapons:** Swords, axes, daggers benefit from backstab
- **Ranged Weapons:** Bows, pistols can technically backstab at close range
- **Default Behavior:** All weapon attacks use same detection algorithm
- **Weapon Specificity:** Some weapons may have backstab bonuses (daggers +25%)

**AI Integration:**

**Backstab Awareness:**
- **Positioning:** AI attempts to avoid exposing backs to enemies
- **Facing Behavior:** AI auto-faces enemies when engaging
- **Backstab Seeking:** Sneak AI prioritizes positioning behind enemies
- **Support Behavior:** Support AI positions to protect ally backs

**Performance Optimization:**

**Angle Calculation Caching:**
```go
// Cache attack angles for performance
type AttackAngleCache struct {
    From       Position
    To         Position
    Angle       float64
    CachedTurn  int
}

var angleCache = make(map[string]AttackAngleCache)

func GetCachedAttackAngle(from, to Position) float64 {
    cacheKey := fmt.Sprintf("%s->%s", from, to)
    
    if cached, exists := angleCache[cacheKey]; exists {
        if cached.CachedTurn == currentTurn {
            return cached.Angle  // Use cached angle
        }
    }
    
    angle := CalculateAttackAngle(from, to)
    angleCache[cacheKey] = AttackAngleCache{
        From:      from,
        To:        to,
        Angle:      angle,
        CachedTurn: currentTurn,
    }
    
    return angle
}
```

**Player Experience:**

**Successful Backstab:**
"I sneak behind the enemy warrior—he hasn't faced me this turn. I attack with my dagger—'BACKSTAB!' appears, and I deal 15 damage instead of 10. The enemy's armor is reduced by half, so my attack penetrates completely. The algorithm calculated that I was 22° behind his facing, well within the 45° backstab detection range!"

**Failed Backstab Attempt:**
"I try to backstab from behind, but the enemy turned around just before my attack struck. His new facing makes me attack from his front now, so I don't get the bonus. The detection algorithm is precise—it recalculates angles each turn based on current orientation, not just initial facing!"

**Tactical Positioning:**
"I notice the enemy AI keeps turning to face whoever is closest. I use my movement advantage to get behind him while he's distracted attacking my teammate. When he does turn, I'm already in position for another backstab. The angle calculation is instant and accurate every turn!"

**Implementation Priority:** HIGH - Required for Phase 2 backstabbing mechanic

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[backstab_detection_algorithm]]`
- **Related Files:** `upsilonbattle/battlearena/entity/entity.go` (orientation), `upsilonbattle/battlearena/ruler/rules/attack.go`
- **Integration:** Works with `mec_backstabbing_mechanic`, `armor_penetration_system`

## EXPECTATION
