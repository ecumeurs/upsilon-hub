# V2 Testing Plan

**Document Version:** 1.0  
**Date:** 2026-04-23  
**Status:** Ready for Review  
**Target:** UpsilonBattle V2 Major Release

---

## EXECUTIVE SUMMARY

This comprehensive testing plan covers all V2 features across the 20-week development timeline. The plan ensures that every system meets ATD compliance, performance standards, and user experience goals before the Q3 2026 launch.

**Testing Philosophy:** Test-Driven Development with continuous verification through ATD traceability and automated CI/CD pipelines.

---

## TESTING STRUCTURE

### Test Categories

| Category | Purpose | Tools | Timeline |
|---|---|---|---|
| **Unit Tests** | Component-level verification | Go testing framework | Continuous |
| **Integration Tests** | System interaction validation | Custom test harness | End of each feature |
| **E2E Tests** | Complete user workflows | upsiloncli + scripts | End of each phase |
| **Performance Tests** | Load and stress testing | Benchmarks | Phase completions |
| **Security Tests** | Vulnerability assessment | Security scanners | Phase 4 & 5 |
| **ATD Compliance** | Documentation verification | ATD tools | Continuous |
| **User Acceptance** | Real-world validation | Beta testers | Phase 5 |

### Test Coverage Targets

- **Unit Test Coverage:** ≥85% for new code, ≥70% for modified code
- **Integration Coverage:** 100% of API endpoints and system boundaries
- **E2E Coverage:** 100% of critical user flows
- **ATD Coverage:** 100% @spec-link compliance, 90% @test-link coverage

---

## PHASE 1: FOUNDATION SYSTEMS (Weeks 1-4)

### Week 1: Skill Weight & Grading System
**Reference:** ISS-065

#### Unit Tests
```go
// @test-link [[ISS-065_skill_weight_calculator]]
func TestSkillWeightCalculator(t *testing.T) {
    // Test benefits calculation
    testCases := []struct {
        skill Skill
        expectedSW int
    }{
        // 100% damage = +100 SW
        {skill: Skill{Damage: 100}, expectedSW: 100},
        // Range 3 = +20 SW (10 per cell > 1)
        {skill: Skill{Range: 3}, expectedSW: 20},
        // +25% crit = +50 SW
        {skill: Skill{CritChance: 25}, expectedSW: 50},
        // Combined benefits
        {skill: Skill{Damage: 100, Range: 3, CritChance: 25}, expectedSW: 170},
    }
    
    for _, tc := range testCases {
        result := CalculateSkillWeight(tc.skill)
        assert.Equal(t, tc.expectedSW, result.PositiveSW)
    }
}

// @test-link [[ISS-065_skill_grading_algorithm]]
func TestSkillGrading(t *testing.T) {
    testCases := []struct {
        sw int
        grade string
    }{
        {100, "Grade I"},      // 0-150
        {250, "Grade II"},     // 151-300
        {400, "Grade III"},    // 301-500
        {600, "Grade IV"},     // 501-750
        {800, "Grade V"},      // 750+
    }
    
    for _, tc := range testCases {
        result := GradeSkill(tc.sw)
        assert.Equal(t, tc.grade, result)
    }
}

// @test-link [[ISS-065_credit_cost_formula]]
func TestCreditCostCalculation(t *testing.T) {
    testCases := []struct {
        positiveSW int
        expectedCost int
    }{
        {100, 200},  // Basic attack
        {250, 500},  // Fireball
        {400, 800},  // Meteor Swarm
    }
    
    for _, tc := range testCases {
        result := CalculateCreditCost(tc.positiveSW)
        assert.Equal(t, tc.expectedCost, result)
    }
}
```

#### Integration Tests
- Skill generation produces balanced skills (net SW = 0)
- Credit cost formula integrates with shop system
- Skill grading updates UI correctly
- Skill generator uses SW calculator

#### E2E Tests
```bash
# E2E-001: Skill Generation Verification
# 1. Generate 100 random skills
# 2. Verify all have net SW = 0
# 3. Verify grade distribution follows expected ranges
# 4. Verify credit costs = positiveSW × 2
upsiloncli test_e2e skill_generation_balance
```

### Week 2: Time-Based Mechanics Core
**Reference:** ISS-066

#### Unit Tests
```go
// @test-link [[ISS-066_timebased_entity_creation]]
func TestTimeBasedEntityCreation(t *testing.T) {
    entity := CreateTimeBasedEntity("FireballChannel", 400, casterID)
    
    assert.Equal(t, "FireballChannel", entity.Name)
    assert.Equal(t, 400, entity.ChannelDelay)
    assert.Equal(t, casterID, entity.ControllerID)
    assert.Equal(t, TimeBased, entity.Type)
}

// @test-link [[ISS-066_ondeath_trigger]]
func TestOnDeathTrigger(t *testing.T) {
    gs := NewGameState()
    entity := CreateEntityWithOnDeath(func() {
        gs.AddMessage("Explosion!")
    })
    
    entity.Die()
    assert.Contains(t, gs.Messages, "Explosion!")
}

// @test-link [[ISS-066_caster_tracking]]
func TestEffectCasterTracking(t *testing.T) {
    caster := CreateEntity("Player")
    target := CreateEntity("Enemy")
    
    effect := ApplyEffect(caster, target, "Poison", 10)
    assert.Equal(t, caster.ID, effect.CasterID)
}
```

#### Integration Tests
- TimeBased entities properly expire after delay
- OnDeath triggers execute in correct order
- Caster tracking persists through entity death
- Channeling interrupts work correctly

#### E2E Tests
```bash
# E2E-002: Channeling Mechanics
# 1. Player starts channeling Fireball (400 delay)
# 2. Verify channeling entity appears on grid
# 3. Interrupt channeling by attacking player
# 4. Verify Fireball does not execute
# 5. Test completed channeling executes effect
upsiloncli test_e2e channeling_interruption
```

### Week 3: Grid System Updates
**Reference:** ISS-066

#### Unit Tests
```go
// @test-link [[ISS-066_multi_entity_cells]]
func TestMultiEntityCells(t *testing.T) {
    grid := NewGrid(10, 10)
    character := CreateCharacter("Player")
    effect := CreateEffect("FireTrap")
    
    grid.AddEntity(character, 5, 5)
    grid.AddEntity(effect, 5, 5)
    
    cell := grid.GetCell(5, 5)
    assert.Len(t, cell.Entities, 2)
}

// @test-link [[ISS-066_walkthrough_property]]
func TestWalkThroughProperty(t *testing.T) {
    grid := NewGrid(10, 10)
    character := CreateCharacter("Player")
    effect := CreateEffectWithProperty("GasCloud", property.WalkThrough, true)
    
    grid.AddEntity(character, 5, 5)
    grid.AddEntity(effect, 5, 6)
    
    // Should be able to move through gas cloud
    assert.True(t, grid.CanMove(character, 5, 5, 5, 6))
}
```

#### Integration Tests
- Multiple entities can occupy same cell
- WalkThrough property allows movement through effects
- Cell-attached effects modify movement cost
- Collision detection handles multi-entity cells

### Week 4: Database & API Foundation
**Reference:** All V2 systems

#### Unit Tests
```go
// @test-link [[ISS-067_credit_tracking]]
func TestCreditTracking(t *testing.T) {
    character := CreateCharacter()
    
    character.AddCredits(100)
    assert.Equal(t, 100, character.GetCredits())
    
    character.SpendCredits(50)
    assert.Equal(t, 50, character.GetCredits())
}
```

#### Integration Tests
- Database schema supports new V2 fields
- API endpoints return correct data structures
- Character creation uses V2 stats (30-50/10/5/3)
- Credit tracking persists through saves

#### E2E Tests
```bash
# E2E-003: Database Migration
# 1. Create V1 character (3/1/1/1 stats)
# 2. Run migration script
# 3. Verify character converted to V2 stats
# 4. Verify credits initialized to 0
# 5. Verify skills/equipment fields created
upsiloncli test_e2e database_migration_v2
```

---

## PHASE 2: CORE GAMEPLAY (Weeks 5-8)

### Week 5: Skill Selection System
**Reference:** ISS-065

#### Unit Tests
```go
// @test-link [[ISS-065_skill_choice_creation]]
func TestSkillChoiceAtCreation(t *testing.T) {
    character := CreateCharacter()
    availableSkills := GenerateSkillChoices(character.Level, 3)
    
    assert.Len(t, availableSkills, 3)
    
    // Player selects one skill
    character.SelectSkill(availableSkills[0])
    assert.Len(t, character.Skills, 1)
}
```

#### Integration Tests
- Skill selection modal displays correct options
- Skill tooltips show SW/grade information
- Progression triggers appear at correct levels
- Skill reforging deducts credits correctly

#### E2E Tests
```bash
# E2E-004: Skill Selection Flow
# 1. Create new character
# 2. Select 1 of 3 skills offered
# 3. Verify skill added to character
# 4. Level up to 10
# 5. Verify new skill selection available
# 6. Test skill reforging for credits
upsiloncli test_e2e skill_selection_progression
```

### Week 6: Time-Based Mechanics Implementation
**Reference:** ISS-066

#### Unit Tests
```go
// @test-link [[ISS-066_channeling_skills]]
func TestChannelingSkill(t *testing.T) {
    gs := NewGameState()
    caster := CreateCharacter("Player")
    target := CreateCharacter("Enemy")
    
    skill := CreateChannelingSkill("Fireball", 400, 100)
    gs.CastSkill(caster, target, skill)
    
    // Channeling entity should be created
    channeling := gs.FindChannelingEntity(caster.ID)
    assert.NotNil(t, channeling)
    assert.Equal(t, 400, channeling.Delay)
}
```

#### Integration Tests
- Channeling skills spawn temporary entities
- Area effects affect correct zones
- Channeling interruption prevents effect
- Effect expiration cleans up properly

#### E2E Tests
```bash
# E2E-005: Area Effects
# 1. Create healing zone (3x3 area)
# 2. Place effect on grid
# 3. Move ally into zone
# 4. Verify ally healed each turn
# 5. Test zone expiration after 5 turns
upsiloncli test_e2e area_effects_zones
```

### Week 7: Backstabbing System
**Reference:** ISS-070

#### Unit Tests
```go
// @test-link [[ISS-070_backstab_detection]]
func TestBackstabDetection(t *testing.T) {
    attacker := CreateCharacter("Rogue")
    target := CreateCharacter("Enemy")
    
    // Position target facing Right
    target.SetOrientation(Right)
    target.SetPosition(5, 5)
    
    // Position attacker behind target (facing Left)
    attacker.SetPosition(5, 6)
    
    assert.True(t, attacker.IsBackstabbing(target))
}

// @test-link [[ISS-070_backstab_damage]]
func TestBackstabDamageCalculation(t *testing.T) {
    attacker := CreateCharacter("Rogue")
    target := CreateCharacter("Enemy")
    target.SetArmorRating(10)  // 50% ignored = 5
    
    baseDamage := 20
    backstabDamage := CalculateBackstabDamage(attacker, target, baseDamage)
    
    // 20 * 1.5 = 30, ignore 5 armor = 25
    assert.Equal(t, 25, backstabDamage)
}
```

#### Integration Tests
- Backstab detection uses orientation system
- 150% damage multiplier applies correctly
- 50% armor penetration excludes shields
- Weapon attacks support backstabbing
- AI avoids exposing back to enemies

#### E2E Tests
```bash
# E2E-006: Backstabbing Combat
# 1. Create rogue character with dagger
# 2. Position enemy facing away
# 3. Attack from behind
# 4. Verify 150% damage and armor penetration
# 5. Test AI backstab awareness
upsiloncli test_e2e backstabbing_combat
```

### Week 8: Credit Earning System
**Reference:** ISS-067

#### Unit Tests
```go
// @test-link [[ISS-067_damage_credits]]
func TestDamageCreditEarning(t *testing.T) {
    attacker := CreateCharacter("Player")
    target := CreateCharacter("Enemy")
    target.SetHP(100)
    
    // Deal 25 damage
    target.TakeDamage(25)
    
    // Should earn 25 credits
    assert.Equal(t, 25, attacker.GetCredits())
}

// @test-link [[ISS-067_healing_credits]]
func TestHealingCreditEarning(t *testing.T) {
    healer := CreateCharacter("Healer")
    ally := CreateCharacter("Ally")
    ally.SetHP(50)
    
    // Heal 15 HP
    ally.Heal(15)
    
    // Should earn 15 credits
    assert.Equal(t, 15, healer.GetCredits())
}
```

#### Integration Tests
- 1 HP damage = 1 credit earned
- 1 HP healing = 1 credit earned
- Shield casters earn mitigation credits
- Status effects earn SW/10 credits
- Credits persist through character save/load

#### E2E Tests
```bash
# E2E-007: Credit Earning Flow
# 1. Start match with 3 characters
# 2. Deal 50 damage, heal 25 HP, apply poison (100 SW)
# 3. Verify credits: 50 + 25 + 10 = 85
# 4. Test credit summary at match end
upsiloncli test_e2e credit_earning_multiplayer
```

---

## PHASE 3: EQUIPMENT & ECONOMY (Weeks 9-12)

### Week 9: Equipment System
**Reference:** ISS-068

#### Unit Tests
```go
// @test-link [[ISS-068_equipment_slots]]
func TestEquipmentSlots(t *testing.T) {
    character := CreateCharacter()
    
    armor := CreateArmor("Chainmail", 5)
    weapon := CreateWeapon("Sword", 10)
    utility := CreateUtility("Ring", 1)
    
    character.EquipArmor(armor)
    character.EquipWeapon(weapon)
    character.EquipUtility(utility)
    
    assert.NotNil(t, character.GetArmor())
    assert.NotNil(t, character.GetWeapon())
    assert.NotNil(t, character.GetUtility())
}

// @test-link [[ISS-068_weapon_as_skill]]
func TestWeaponAsSkill(t *testing.T) {
    character := CreateCharacter()
    weapon := CreateWeapon("Longsword", 15)
    weapon.SetRange(2)
    
    character.EquipWeapon(weapon)
    
    // Attack should use weapon properties
    damage := character.GetAttackDamage()
    assert.Equal(t, 15, damage)
    assert.Equal(t, 2, character.GetAttackRange())
}
```

#### Integration Tests
- 3-slot system limits equipment correctly
- Weapon attacks use weapon properties
- Armor bonuses apply to defense
- Equipment stat bonuses calculate correctly

#### E2E Tests
```bash
# E2E-008: Equipment Management
# 1. Create character with 3-slot equipment
# 2. Equip armor (+5 defense)
# 3. Equip weapon (+10 attack, range 2)
# 4. Equip utility (+3 HP)
# 5. Verify stat bonuses apply
upsiloncli test_e2e equipment_management
```

### Week 10: Shop System
**Reference:** ISS-067

#### Unit Tests
```go
// @test-link [[ISS-067_shop_pricing]]
func TestShopPricing(t *testing.T) {
    skill := CreateSkill("Fireball")
    skill.SetPositiveSW(250)
    
    cost := CalculateShopPrice(skill)
    assert.Equal(t, 500, cost)  // 250 SW × 2
}

// @test-link [[ISS-067_purchase_flow]]
func TestPurchaseFlow(t *testing.T) {
    character := CreateCharacter()
    character.AddCredits(1000)
    
    skill := CreateShopSkill("Lightning", 300)
    
    err := character.PurchaseSkill(skill)
    assert.NoError(t, err)
    assert.Equal(t, 700, character.GetCredits())
    assert.Len(t, character.Skills, 1)
}
```

#### Integration Tests
- Shop displays correct pricing based on SW
- Purchase checks credit balance
- Inventory updates after purchase
- Affordability indicators display correctly

#### E2E Tests
```bash
# E2E-009: Shop Purchase Flow
# 1. Character has 1500 credits
# 2. Browse shop for Grade II skills
# 3. Purchase 2 skills (500 + 500 credits)
# 4. Verify credit balance: 1500 - 1000 = 500
# 5. Test insufficient credits rejection
upsiloncli test_e2e shop_purchase_flow
```

### Week 11: Extended Character Sheet
**Reference:** All V2 systems

#### Integration Tests
- Extended character sheet displays all V2 properties
- Skill tooltips show complete information
- Equipment bonuses calculate correctly
- Buff/debuff display updates in real-time

#### E2E Tests
```bash
# E2E-010: Character Sheet Display
# 1. Create V2 character with full equipment and skills
# 2. Open character sheet
# 3. Verify all stats display correctly
# 4. Verify skill tooltips show SW/grade/cooldown
# 5. Verify equipment shows stat bonuses
upsiloncli test_e2e character_sheet_display
```

### Week 12: Economy Balancing
**Reference:** ISS-067

#### Performance Tests
```bash
# PERF-001: Economy Stress Test
# 1. Simulate 100 matches
# 2. Verify credit earning rates stay balanced
# 3. Check for economy inflation
# 4. Validate shop pricing remains fair
upsiloncli test_perf economy_stress_test
```

#### E2E Tests
```bash
# E2E-011: Economy Progression
# 1. Track character through 20 matches
# 2. Verify credit earning ~100-150 credits/match
# 3. Verify progression pace matches expectations
# 4. Test skill purchase timing balance
upsiloncli test_e2e economy_progression_pace
```

---

## PHASE 4: AI ENHANCEMENT (Weeks 13-16)

### Week 13: AI Architecture
**Reference:** ISS-069

#### Unit Tests
```go
// @test-link [[ISS-069_archetype_controllers]]
func TestArchetypeControllers(t *testing.T) {
    fighter := CreateFighterController()
    ranger := CreateRangerController()
    support := CreateSupportController()
    sneak := CreateSneakController()
    
    assert.IsType(t, &FighterController{}, fighter)
    assert.IsType(t, &RangerController{}, ranger)
    assert.IsType(t, &SupportController{}, support)
    assert.IsType(t, &SneakController{}, sneak)
}

// @test-link [[ISS-069_team_composition]]
func TestTeamComposition(t *testing.T) {
    team := CreateTeam()
    
    team.AddCharacter(CreateFighter())
    team.AddCharacter(CreateRanger())
    team.AddCharacter(CreateSupport())
    team.AddCharacter(CreateSneak())
    
    assert.True(t, team.ValidateComposition())
    
    // Adding second support should fail
    team.AddCharacter(CreateSupport())
    assert.False(t, team.ValidateComposition())
}
```

#### Integration Tests
- Four archetypes use appropriate skill pools
- AI follows player progression rules
- Team composition enforces limits
- AI skill selection uses grade-appropriate skills

### Week 14: Archetype Implementation
**Reference:** ISS-069

#### Unit Tests
```go
// @test-link [[ISS-069_fighter_behavior]]
func TestFighterBehavior(t *testing.T) {
    fighter := CreateFighterController()
    gameState := NewGameState()
    
    // Fighter should prioritize direct approach
    action := fighter.DecideAction(gameState)
    assert.Equal(t, "Attack", action.Type)
    assert.True(t, action.IsAggressive())
}
```

#### Integration Tests
- Fighter uses aggressive melee tactics
- Ranger maintains kiting distance
- Support stays near allies and prioritizes healing
- Sneak uses flanking and backstabbing

#### E2E Tests
```bash
# E2E-012: AI Archetype Behavior
# 1. Create team with each archetype
# 2. Observe AI decision-making over 20 turns
# 3. Verify each archetype behaves according to design
# 4. Test skill usage patterns
upsiloncli test_e2e ai_archetype_behavior
```

### Week 15: AI Progression Integration
**Reference:** ISS-069

#### Integration Tests
- AI stats scale with level (+10 CP per level)
- AI skill grades match level
- AI difficulty scales appropriately
- Bot filling system creates balanced teams

#### E2E Tests
```bash
# E2E-013: AI Progression Scaling
# 1. Create level 1 AI team
# 2. Create level 10 AI team
# 3. Compare stats and skills
# 4. Verify proper scaling
# 5. Test matchmaking balance
upsiloncli test_e2e ai_progression_scaling
```

### Week 16: AI Testing & Balancing
**Reference:** ISS-069

#### Performance Tests
```bash
# PERF-002: AI Decision Performance
# 1. Run 1000 AI decision cycles
# 2. Measure average decision time
# 3. Verify < 100ms per decision
# 4. Test with multiple AI types
upsiloncli test_perf ai_decision_performance
```

#### E2E Tests
```bash
# E2E-014: AI vs Player Balance
# 1. Play 50 matches vs AI
# 2. Track win/loss ratio
# 3. Verify balanced difficulty (~50% win rate)
# 4. Test AI skill effectiveness
upsiloncli test_e2e ai_player_balance
```

---

## PHASE 5: POLISH & TESTING (Weeks 17-20)

### Week 17: UI Integration
**Reference:** All V2 systems

#### Integration Tests
- ActionPanel displays skill buttons correctly
- Skill selection modal works smoothly
- Equipment management UI functions properly
- Character sheet shows all V2 properties

#### E2E Tests
```bash
# E2E-015: UI Integration Suite
# 1. Test all UI components end-to-end
# 2. Verify responsive design
# 3. Test accessibility features
# 4. Validate user flows
upsiloncli test_e2e ui_integration_suite
```

### Week 18: Visual Feedback
**Reference:** All V2 systems

#### E2E Tests
```bash
# E2E-016: Visual Feedback Verification
# 1. Test channeling indicators
# 2. Verify backstab damage display
# 3. Check credit earning animations
# 4. Validate equipment change feedback
upsiloncli test_e2e visual_feedback
```

### Week 19: Testing & Balancing
**Reference:** All V2 systems

#### ATD Compliance Tests
```bash
# ATD-001: ATD Compliance Verification
# 1. Check all new code has @spec-link tags
# 2. Verify all tests have @test-link tags
# 3. Run atd_verify() to check compliance
# 4. Fix any orphaned atoms
mcp__atd__atd_verify()
```

#### Comprehensive Balance Tests
```bash
# BALANCE-001: Skill System Balance
# 1. Generate 500 skills
# 2. Verify all have net SW = 0
# 3. Check grade distribution
# 4. Validate pricing formulas

# BALANCE-002: Economy Balance
# 1. Simulate 1000 matches
# 2. Track credit earning rates
# 3. Verify no inflation
# 4. Test shop affordability

# BALANCE-003: Combat Balance
# 1. Run 500 PvP matches
# 2. Analyze win rates by archetype
# 3. Check damage distribution
# 4. Validate time-based mechanics
```

### Week 20: Documentation & Launch Prep

#### Documentation Tests
```bash
# DOC-001: Documentation Completeness
# 1. Verify all V2 atoms have STABLE status
# 2. Check all code has @spec-link tags
# 3. Verify player guides exist
# 4. Validate API documentation
mcp__atd__atd_stats()
```

#### Launch Readiness Tests
```bash
# LAUNCH-001: Launch Readiness Verification
# 1. Run full test suite
# 2. Verify performance benchmarks
# 3. Check security scan results
# 4. Validate deployment pipeline
# 5. Test rollback procedures
```

---

## PERFORMANCE TESTING

### Performance Benchmarks

| System | Metric | Target | Test Method |
|---|---|---|---|
| **Skill Generation** | 1000 skills | < 100ms | Benchmark |
| **AI Decision** | Per decision | < 100ms | Load test |
| **Credit Calculation** | Per match | < 50ms | Benchmark |
| **Equipment Updates** | Per equip | < 10ms | Integration test |
| **Database Queries** | Character load | < 200ms | Load test |
| **API Response Time** | Average request | < 500ms | Load test |
| **Match Creation** | New match | < 2s | E2E test |

### Load Testing Scenarios

```bash
# LOAD-001: Concurrent Match Simulation
# 1. Create 100 concurrent matches
# 2. Simulate 100 turns per match
# 3. Verify < 5% error rate
# 4. Check average response times
upsiloncli test_load concurrent_matches

# LOAD-002: Database Stress Test
# 1. Simulate 1000 character creations
# 2. Concurrent skill purchases
# 3. Equipment updates
# 4. Verify database integrity
upsiloncli test_load database_stress
```

---

## SECURITY TESTING

### Security Test Categories

| Category | Tests | Tools | Timeline |
|---|---|---|---|  
| **Input Validation** | SQL injection, XSS, command injection | OWASP ZAP, custom tests | Phase 4 |
| **Authentication** | Session management, token security | Security scanners | Phase 4 |
| **Authorization** | Access control, privilege escalation | Manual testing | Phase 4 |
| **Data Validation** | Credit manipulation, stat tampering | Fuzz testing | Phase 5 |
| **API Security** | Rate limiting, request validation | API security tools | Phase 5 |

### Security Test Scenarios

```bash
# SEC-001: Credit Manipulation Tests
# 1. Attempt to modify credits directly
# 2. Test duplicate credit earning
# 3. Verify negative credit prevention
# 4. Check refund vulnerabilities

# SEC-002: Skill System Security
# 1. Test unauthorized skill acquisition
# 2. Verify skill reforging security
# 3. Check skill template access

# SEC-003: Equipment Security
# 1. Test equipment duplication exploits
# 2. Verify stat bonus manipulation prevention
# 3. Check equipment trading security
```

---

## ATD COMPLIANCE TESTING

### ATD Verification Checklist

- [ ] All new code has @spec-link tags
- [ ] All new tests have @test-link tags  
- [ ] All V2 atoms have STABLE status
- [ ] No orphaned atoms exist
- [ ] Dependencies are correctly woven
- [ ] Documentation matches implementation
- [ ] Test coverage meets targets

### ATD Testing Commands

```bash
# Verify ATD compliance
mcp__atd__atd_verify()

# Check for orphaned atoms
mcp__atd__atd_crawl(gaps=true)

# Run health check
mcp__atd__atd_stats()

# Check specific atom coverage
mcp__atd__atd_trace(atom="ISS-065_skill_weight_calculator")
```

---

## USER ACCEPTANCE TESTING

### Beta Testing Plan

| Phase | Participants | Duration | Focus |
|---|---|---|---|
| **Alpha** | Internal team | 1 week | Core functionality |
| **Beta 1** | 50 external testers | 2 weeks | Major systems |
| **Beta 2** | 200 external testers | 2 weeks | Balance and polish |
| **RC** | 1000 players | 1 week | Launch readiness |

### UAT Scenarios

```bash
# UAT-001: New Player Experience
# 1. Create new account and character
# 2. Complete character creation with V2 stats
# 3. Play tutorial match
# 4. Purchase first skill
# 5. Equip first weapon

# UAT-002: Advanced Player Experience
# 1. Import V1 character (if applicable)
# 2. Test character progression
# 3. Use all major V2 features
# 4. Complete 10 matches
# 5. Provide feedback on balance

# UAT-003: Competitive Experience
# 1. Play PvP matches
# 2. Test matchmaking
# 3. Experience all archetypes
# 4. Use advanced tactics
# 5. Report balance issues
```

---

## REGRESSION TESTING

### Regression Test Suite

```bash
# REGRESSION-001: V1 Feature Preservation
# 1. Verify all V1 features still work
# 2. Test character creation (basic)
# 3. Test combat mechanics (basic)
# 4. Test progression (basic)
# 5. Verify no V1 features broken

# REGRESSION-002: Database Compatibility
# 1. Test with V1 database
# 2. Run migration scripts
# 3. Verify data integrity
# 4. Test rollback procedures

# REGRESSION-003: API Compatibility
# 1. Test existing API endpoints
# 2. Verify response formats
# 3. Check authentication flows
# 4. Test error handling
```

---

## TEST AUTOMATION & CI/CD

### Automated Testing Pipeline

```yaml
# .github/workflows/v2_testing.yml
name: V2 Testing Pipeline

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Run unit tests
        run: make test-unit
      
      - name: Generate coverage report
        run: make test-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  integration-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    steps:
      - name: Run integration tests
        run: make test-integration
      
      - name: Test database migrations
        run: make test-migrations

  e2e-tests:
    runs-on: ubuntu-latest
    needs: [unit-tests, integration-tests]
    steps:
      - name: Run E2E tests
        run: make test-e2e
      
      - name: Test V2 features
        run: upsiloncli test_e2e v2_feature_suite

  atd-compliance:
    runs-on: ubuntu-latest
    steps:
      - name: Check ATD compliance
        run: mcp__atd__atd_verify()
      
      - name: Check for orphans
        run: mcp__atd__atd_crawl(gaps=true)

  performance-tests:
    runs-on: ubuntu-latest
    needs: [integration-tests, e2e-tests]
    steps:
      - name: Run performance benchmarks
        run: make test-performance
      
      - name: Run load tests
        run: make test-load

  security-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Run security scans
        run: make test-security
      
      - name: Check dependencies
        run: make audit-dependencies
```

---

## TEST METRICS & REPORTING

### Success Metrics

| Metric | Target | Measurement |
|---|---|---|
| **Unit Test Coverage** | ≥85% new code, ≥70% modified | Codecov |
| **Integration Coverage** | 100% APIs/systems | Test reports |
| **E2E Coverage** | 100% critical flows | Test execution |
| **ATD Compliance** | 100% @spec-link, 90% @test-link | ATD tools |
| **Performance** | All benchmarks met | Benchmark results |
| **Security** | 0 critical vulnerabilities | Security scans |
| **Defect Density** | < 5 bugs/KLOC | Bug tracker |

### Reporting

```bash
# Weekly Test Report
make test-report-weekly

# Coverage Report
make test-coverage-report

# ATD Compliance Report
make test-atd-report

# Performance Report
make test-performance-report
```

---

## ROLLBACK & CONTINGENCY PLANNING

### Rollback Triggers

- Critical bugs in production
- Performance degradation >20%
- Security vulnerabilities discovered
- User experience issues >10% negative feedback

### Rollback Procedures

```bash
# Database Rollback
# 1. Stop V2 deployment
# 2. Restore V1 database backup
# 3. Verify data integrity
# 4. Restart V1 services

# Code Rollback
# 1. Revert to last stable commit
# 2. Clear caches
# 3. Restart services
# 4. Verify functionality

# Partial Rollback
# 1. Disable specific V2 features
# 2. Maintain core functionality
# 3. Monitor system health
```

---

## CONCLUSION

This comprehensive testing plan ensures that UpsilonBattle V2 meets all quality standards, performance requirements, and user experience goals. The phased approach aligns with the 20-week development timeline, ensuring continuous testing and validation throughout the development process.

**Key Testing Principles:**
1. **Test-Driven Development** - Write tests before implementation
2. **ATD Compliance** - Every feature has documentation and tests
3. **Continuous Verification** - Automated testing in CI/CD pipeline
4. **User-Focused** - Real-world testing scenarios and beta testing
5. **Performance-First** - Benchmark and load test throughout development

**Success Criteria:**
- All unit, integration, and E2E tests passing
- 100% ATD compliance achieved
- All performance benchmarks met
- Zero critical security vulnerabilities
- Positive user feedback in beta testing

With this testing plan, UpsilonBattle V2 will launch with confidence in quality, performance, and user experience.