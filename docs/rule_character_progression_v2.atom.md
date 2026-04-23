---
id: rule_character_progression_v2
human_name: Character Progression Rule V2
type: RULE
layer: ARCHITECTURE
version: 2.0
status: DRAFT
priority: 5
tags: [progression, character, stats]
parents:
  - [[entity_character]]
  - [[requirement_customer_player_profile]]
dependents: []
---

# Character Progression Rule V2

## INTENT
To govern character stat progression using the 100 CP point-buy system with x10 baseline stats and weighted attribute costs for V2 character development.

## THE RULE / LOGIC
**V2 Starting Stats (x10 Baseline):**
- **HP:** 30-50 (random within range)
- **Attack:** 10 (fixed baseline)
- **Defense:** 5 (fixed baseline)  
- **Movement:** 3 (fixed baseline)

**Point-Buy System:**
- **Starting Pool:** 100 Character Points (CP) to spend
- **No Random Distribution:** Players allocate points strategically
- **Weighted Costs:** Different attributes have different CP costs

**Standard Attribute Costs:**
- **HP (+1):** Cost 1 CP (Linear, cheap)
- **Attack (+1):** Cost 5 CP (Direct damage scaling is powerful)
- **Defense (+1):** Cost 5 CP (Direct damage mitigation is equally powerful)

**Exotic Attribute Costs:**
- **Critical Chance (+1%):** Cost 10 CP (High value, caps easily)
- **Critical Multiplier (+5%):** Cost 5 CP
- **Jump Height (+1):** Cost 15 CP (Drastically alters terrain navigation)

**Movement Premium:**
- **Movement (+1 cell):** Cost 30 CP (Most potent stat, expensive = natural restriction)
- **Note:** Eliminates V1 "once every 5 wins" movement hard-lock

**Progression Rules:**
- **Win Reward:** +10 CP per win (instead of +1 point)
- **Total Cap:** 100 + (total_wins × 10)
- **Unrestricted Allocation:** CP can be spent on any attribute (no movement locks)

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_character_progression_v2]]`
- **Test Names:** `TestV2CharacterCreation`, `TestV2ProgressionAllocation`, `TestMovementNaturalRestriction`
