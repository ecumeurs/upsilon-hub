# UpsilonBattle V2: Tactical RPG Evolution

**Version:** 2.0  
**Status:** Development Plan  
**Target Release:** Q3 2026  
**Development Timeline:** 20 Weeks  

---

## EXECUTIVE SUMMARY

UpsilonBattle V2 represents a comprehensive evolution from the basic tactical RPG foundation into a rich, progression-driven combat experience. This major update introduces skill systems, time-based mechanics, credit economy, equipment progression, AI enhancement, and combat depth improvements.

**Core Philosophy:** Build upon the excellent existing architecture with extensible systems that enable infinite tactical possibilities while maintaining balance and performance.

---

## 🎯 MAJOR FEATURE OVERVIEW

### **1. Skill System Overhaul** ⭐⭐⭐⭐⭐
**Impact:** Transforms simple attacks into complex tactical choices

**Key Components:**
- **Skill Weight (SW) System:** Mathematical balance framework where benefits = costs
- **Skill Grading:** I-V tiers based on power (0-150, 151-300, 301-500, 501-750, 750+ SW)
- **Skill Selection:** Choose 1 of 3 skills at creation, every 10 levels
- **Skill Reforging:** Modify skills every 5 levels for credits

**Player Experience:** Characters evolve from basic attacks to specialized skill sets, creating unique playstyles and tactical depth.

**Technical Foundation:** 30+ skill properties (Damage, Heal, Range, Zone, etc.) with granular control over effects.

---

### **2. Time-Based Mechanics** ⭐⭐⭐⭐⭐
**Impact:** Adds temporal strategy and risk/reward decisions to combat

**Key Components:**
- **Channeling:** Skills with pre-execution delay (400 delay to cast Fireball)
- **Temporary Entities:** Unified system for channeling, traps, area effects
- **Trigger System:** OnTurn, OnStep, OnDeath event handlers
- **Effect Caster Tracking:** Effects remember originators for credits/interruption
- **Multi-Entity Cells:** Multiple effects can occupy same location

**Player Experience:** Players must decide between immediate actions vs. powerful delayed effects. "Do I cast this 400-delay channeling spell, knowing I could be interrupted?"

**Technical Innovation:** Brilliant architecture where 1 skill effect = 1 entity, enabling infinite temporal possibilities through simple Turner integration.

---

### **3. Credit Economy** ⭐⭐⭐⭐⭐
**Impact:** Provides clear progression goals and reward systems

**Key Components:**
- **Base Earning:** 1 HP damage = 1 credit (healing also earns credits)
- **Support Credits:** 1 HP mitigated = 1 credit (for shield caster)
- **Status Credits:** SkillWeight/10 credits per application
- **Shop System:** Purchase skills (SW × 2 credits) and equipment
- **Effect Caster Tracking:** Credits go to original caster even if they die

**Player Experience:** All combat actions contribute to progression. Damage dealers, healers, and shield users all earn meaningfully.

**Economic Balance:** Transparent 1:1 HP-to-credit ratio with multiple earning paths supporting different playstyles.

---

### **4. Equipment System** ⭐⭐⭐⭐
**Impact:** Adds character customization and strategic equipment choices

**Key Components:**
- **3-Slot System:** 1 armor, 1 utility, 1 weapon
- **Weapon-as-Skill:** Equipped weapons transform basic attacks into skill-based attacks
- **Stat Bonuses:** Equipment provides direct attribute modifications
- **Equipment Properties:** ArmorRating, WeaponRange, CritChance, special effects

**Player Experience:** Characters become distinct through equipment choices. "Should I equip this +5 Attack sword or +3 Defense armor?"

**Strategic Depth:** Weapon variety (melee, ranged, two-handed) creates diverse combat approaches.

---

### **5. AI Enhancement** ⭐⭐⭐⭐
**Impact:** Provides challenging, varied AI opponents

**Key Components:**
- **4 Archetypes:** Fighter (aggressive melee), Ranger (ranged/kiting), Support (healing/buffs), Sneak (backstabbing/positioning)
- **Player Progression:** AI follows same point system and level scaling
- **Team Composition:** Max 1 support, max 1 sneak per AI team
- **Archetype Skills:** Each AI type uses appropriate skill pools

**Player Experience:** Fighting different AI archetypes requires different strategies. No more "all enemies act the same."

**Balanced Challenge:** AI power scales with player level while maintaining team composition constraints.

---

### **6. Backstabbing** ⭐⭐⭐⭐⭐
**Impact:** Rewards positional tactics and sneak playstyles

**Key Components:**
- **150% Damage:** Backstab multiplier for rear attacks
- **50% Armor Penetration:** Ignores half armor rating
- **Shield Application:** Shield still applies fully (not penetrated)
- **Weapon Scope:** All weapons support backstabbing (skills excluded for now)

**Player Experience:** Positioning matters immensely. "Get behind that enemy for massive damage!"

**Sneak Archetype:** Dedicated backstabbing mechanics make the Sneak AI/playstyle viable and rewarding.

---

### **7. Stat System Redesign** ⭐⭐⭐⭐⭐
**Impact:** Enables meaningful skill percentages and character variety

**Key Components:**
- **x10 Baseline:** HP 30-50, Attack 10, Defense 5, Movement 3
- **100 CP Point-Buy:** Strategic point allocation vs. random 4-point distribution
- **Weighted Costs:** Attack/Defense = 5 CP, HP = 1 CP, Movement = 30 CP (natural restriction)
- **Exotic Attributes:** Crit Chance, Crit Multiplier, Jump Height with specific costs

**Player Experience:** Characters become meaningfully different. "I'll build a tank with high HP and Defense" vs. "I'll build a glass cannon with high Attack and Crit."

**Skill Viability:** Percentage modifiers now work meaningfully. 120% of 10 Attack = 12 damage (vs. useless 120% of 1 Attack = 1 damage).

---

## 🏗️ ARCHITECTURAL BREAKTHROUGHS

### **Unified Temporary Entity System**
**Genius Insight:** All time-based mechanics = temporary entities with controllers
- **Channeling** = Entity with OnTurn trigger, 400 delay
- **Traps** = Entity with OnStep trigger, kills itself
- **Area Effects** = Master entity losing 1 HP per turn, affects zone

**Impact:** One system replaces schedulers, queues, and complex timing logic. Elegant, extensible, maintainable.

### **Skill Weight Mathematical Framework**
**Perfect Balance:** Net SW = 0 (benefits = costs)
- **Benefits Table:** Damage, Crit, Range, AoE, Stun, Poison all have precise SW costs
- **Payments Table:** Delay, Channeling, MP/SP, HP, Cooldown have precise SW reductions
- **Grading & Pricing:** Derived from Total Positive SW, automatic balance

**Impact:** Mathematical skill design, automatic balance checking, clear upgrade paths.

### **x10 Stat Scaling**
**Critical Fix:** Makes percentage modifiers meaningful
- **V1 Problem:** 120% of 1 Attack = 1.2 damage → useless skills
- **V2 Solution:** 120% of 10 Attack = 12 damage → viable skills
- **Progression Impact:** 100 CP system enables meaningful character variety

**Impact:** Skills work as intended, character progression creates diversity, combat math is sensible.

---

## 📊 IMPLEMENTATION ROADMAP

### **Phase 1: Foundation Systems (Weeks 1-4)**
**Focus:** Core infrastructure that enables everything else

**Week 1: Skill Weight & Grading**
- Implement SkillWeight calculator and grading algorithm
- Update skill generator to use SW system
- Build credit cost formula (SW × 2)
- Create skill template system

**Week 2: Time-Based Mechanics Core**
- Add TimeBased entity type and ExpirationController
- Implement OnDeath() trigger
- Add caster tracking to effects
- Update cost system for Channeling/Delay separation

**Week 3: Grid System Updates**
- Enable multi-entity cells (character + effects)
- Update collision logic for WalkThrough property
- Create cell-attached effects system
- Implement movement cost modifiers (quagmire)

**Week 4: Database & API Foundation**
- Create skill_templates, character skills/equipment fields
- Implement credits tracking and shop inventory table
- Build core API endpoints for new systems
- Update existing data structures for V2 stats

### **Phase 2: Core Gameplay (Weeks 5-8)**
**Focus:** Major gameplay features

**Week 5: Skill Selection System**
- Implement skill choice at character creation (1 of 3)
- Create skill selection modal and progression triggers
- Add skill tooltips with SW/grade display
- Build skill reforging interface

**Week 6: Time-Based Mechanics Implementation**
- Implement channeling skills and temporary entity spawning
- Build area effect system (poisonous fog, healing zones)
- Add channeling interruption mechanics
- Implement effect expiration and cleanup

**Week 7: Backstabbing System**
- Implement back detection algorithm using orientation system
- Add 150% damage multiplier and 50% armor penetration
- Update weapon attack logic and AI backstab awareness
- Create visual feedback for backstabs

**Week 8: Credit Earning System**
- Implement 1 HP = 1 coin (damage & healing)
- Create damage mitigation credit system (shield caster tracking)
- Add status effect credit earning (SW/10 per application)
- Build credit tracking UI and match summary

### **Phase 3: Equipment & Economy (Weeks 9-12)**
**Focus:** Equipment system and shop

**Week 9: Equipment System**
- Design 3-slot system (armor, utility, weapon)
- Create equipment inventory schema and equip/unequip mechanics
- Build equipment stat bonus system
- Implement weapon-as-skill integration

**Week 10: Shop System**
- Build shop interface with skill/equipment purchasing
- Implement inventory management and credit spending
- Add skill/equipment filters and search
- Create affordability indicators and purchase confirmations

**Week 11: Extended Character Sheet**
- Create comprehensive character stats UI (HP, MP, SP, etc.)
- Display equipment bonuses and skill lists with cooldowns
- Implement buff/debuff display and effect tracking
- Build character progression visualization

**Week 12: Economy Balancing**
- Balance credit earning rates and shop prices
- Test progression pacing and skill availability
- Adjust skill reforging costs and equipment pricing
- Implement economy caps and inflation controls if needed

### **Phase 4: AI Enhancement (Weeks 13-16)**
**Focus:** AI archetypes and behavior

**Week 13: AI Architecture**
- Create four archetype controllers extending base controller
- Implement archetype-specific skill pools and decision trees
- Add team composition enforcement and validation
- Build AI skill selection logic by grade

**Week 14: Archetype Implementation**
- Implement Fighter, Ranger, Support, Sneak controllers
- Add archetype-specific stat allocation algorithms
- Build skill usage integration into AI decision making
- Create AI targeting and positioning priorities

**Week 15: AI Progression Integration**
- Connect AI to player progression rules (+10 CP per win)
- Implement AI skill grade progression and level matching
- Create AI difficulty scaling and team composition logic
- Build bot filling system for matchmaking

**Week 16: AI Testing & Balancing**
- Test AI archetype behaviors and team compositions
- Balance AI vs player matchups and skill effectiveness
- Implement AI tactical improvements and positioning logic
- Performance testing and optimization

### **Phase 5: Polish & Testing (Weeks 17-20)**
**Focus:** UI polish, testing, refinement

**Week 17: UI Integration**
- Update ActionPanel for skill buttons and skill selection
- Create equipment management UI and shop interface
- Implement extended character sheet with all new properties
- Build visual feedback systems (channeling, backstabs, credits)

**Week 18: Visual Feedback**
- Add channeling indicators and skill effect visualizations
- Create backstabbing feedback and damage multipliers
- Implement credit earning animations and purchase effects
- Build equipment change visual feedback

**Week 19: Testing & Balancing**
- Comprehensive skill system testing and balance verification
- Time-based mechanics testing and edge case handling
- Credit economy balance testing and progression pacing validation
- Equipment system testing and AI behavior verification

**Week 20: Documentation & Launch Prep**
- Update ATOM documentation with all V2 systems
- Create player guides and admin documentation
- Prepare launch announcements and marketing materials
- Final system testing and deployment preparation

---

## 🎨 PLAYER EXPERIENCE EVOLUTION

### **V1 Experience (Current)**
- **Character Creation:** 4 random points on 3/1/1/1 base stats
- **Combat:** Basic attacks, simple positioning
- **Progression:** +1 point per win, hard movement lock
- **Strategy:** Limited tactical depth, mostly positioning

### **V2 Experience (Planned)**
- **Character Creation:** 100 CP strategic allocation on 30-50/10/5/3 base stats
- **Combat:** Skills, channeling, backstabbing, equipment bonuses
- **Progression:** +10 CP per win, skill selection, equipment upgrades
- **Strategy:** Deep tactical decisions with risk/reward trade-offs

### **V2 Sample Gameplay Scenario**

**Turn 1:**
- Player character (Level 5) chooses "Fireball" skill (Grade II, 250 SW, 500 credits)
- Enemy Fighter charges with "Power Strike" skill

**Turn 3:**
- Player starts channeling "Healing Zone" (400 delay, benefits Support AI ally)
- Enemy Ranger uses "Precision Shot" (Grade I, 120 SW) on player
- Player's "Healing Zone" will complete in 400 delay, healing ally

**Turn 5:**
- Player's "Healing Zone" activates, healing ally for 15 HP
- Player earns 15 credits from healing
- Enemy attempts to interrupt but player positioned safely

**Strategic Depth:** Multiple timing layers, resource management, positional tactics, credit economy all working together.

---

## 🔧 TECHNICAL IMPROVEMENTS

### **Code Architecture**
- **Unified Systems:** Temporary entities handle all timing mechanics
- **Extensible Properties:** 30+ skill properties enable infinite skill variety
- **Clean Separation:** Pre/post execution costs, caster tracking, effect origins
- **Modular Controllers:** AI archetypes extend base controller cleanly

### **Database Schema**
- **Character Expansion:** Skills, equipment, credits fields added
- **Skill Templates:** Predefined and procedurally generated skills
- **Equipment Library:** Comprehensive item definitions with properties
- **Credit Tracking:** Per-character credit balances and earning history

### **API Extensions**
- **Skill Management:** Selection, reforging, and inventory endpoints
- **Shop System:** Browsing, purchasing, and affordability checking
- **Character Stats:** Extended character information with all V2 properties
- **Progression APIs:** Level-based skill availability and CP allocation

### **Frontend Enhancements**
- **Skill Interface:** Selection modals, tooltips, cooldown displays
- **Equipment UI:** Slot management, stat bonuses, preview systems
- **Shop Interface:** Browsing, purchasing, inventory management
- **Visual Feedback:** Channeling indicators, backstab highlights, credit animations

---

## 📈 EXPECTED IMPACT

### **Player Engagement**
- **Deeper Progression:** Skills and equipment provide long-term goals
- **Strategic Variety:** Multiple playstyles become viable and rewarding
- **Replay Value:** Different character builds and equipment combinations

### **Combat Depth**
- **Temporal Strategy:** Channeling vs. immediate actions
- **Positional Tactics:** Backstabbing, flanking, area control
- **Resource Management:** SP/MP decisions, equipment choices, skill selection

### **Social Features**
- **Economy:** Credit-based progression enables trading potential
- **Competition:** Skill grades and equipment create power hierarchies
- **Cooperation:** Support credits reward team play and healing

### **Technical Excellence**
- **Architecture:** Unified systems reduce maintenance complexity
- **Performance:** Efficient temporary entity system scales well
- **Balance:** Mathematical frameworks enable fair and consistent design

---

## 🚀 SUCCESS METRICS

### **Quantitative Goals**
- **Skill Diversity:** 50+ unique skills across all grades by launch
- **Player Retention:** 30% increase in session length due to progression
- **Match Variety:** 200% increase in tactical scenarios
- **AI Engagement:** 50% increase in AI challenge satisfaction

### **Qualitative Goals**
- **Strategic Depth:** Players report meaningful tactical decisions
- **Balance Perception:** Players feel progression is fair and rewarding
- **System Clarity:** New systems are intuitive and well-explained
- **Performance:** No degradation in match speed or responsiveness

---

## 🎯 TARGET AUDIENCE

### **Existing Players**
- **V1 Veterans:** Will appreciate massive depth expansion
- **Progression Seekers:** Love skill trees and equipment upgrades
- **Tactical Players:** Excited by channeling and backstabbing

### **New Players**
- **Progression Clarity:** Clear credit earning and spending
- **Build Variety:** Multiple paths to character development
- **Learnability:** Graded skill system provides natural difficulty curve

### **Competitive Players**
- **Balance Framework:** Mathematical SW system ensures fair competition
- **Meta Evolution:** Skill reforging allows adaptation to balance changes
- **Strategic Depth:** Complex systems enable high-skill play

---

## 📚 DOCUMENTATION UPDATES

### **New ATOM Documentation**
- 15+ new ATOM atoms covering all V2 systems
- Updated existing atoms for V2 compatibility
- Comprehensive technical references for all new mechanics

### **Player Guides**
- Skill selection and progression guide
- Equipment and stat allocation strategies
- Combat tactics for new mechanics
- Credit earning and spending optimization

### **Technical Documentation**
- Architecture documentation for new systems
- API references for all new endpoints
- Database schema updates and migration guides
- Performance characteristics and optimization notes

---

## 🏆 V2 VISION STATEMENT

**UpsilonBattle V2 transforms a solid tactical RPG foundation into a rich, progression-driven experience. By building upon excellent existing architecture with extensible systems, we enable infinite tactical possibilities while maintaining balance and performance.**

**The unified temporary entity system, mathematical skill framework, and x10 stat scaling are architectural breakthroughs that will serve UpsilonBattle for years to come. Players will experience deeper combat, more meaningful progression, and diverse strategic options.**

**V2 isn't just an update—it's an evolution that establishes UpsilonBattle as a premier tactical RPG with systems capable of supporting years of content and community growth.**

---

**Status:** Ready for Implementation  
**Next Steps:** Begin Phase 1 implementation with Skill Weight Calculator  
**Contact:** Development Team for detailed technical specifications