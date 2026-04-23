---
id: rule_progression
human_name: Character Progression Rule
type: RULE
layer: ARCHITECTURE
version: 2.0
status: STABLE
priority: 5
tags: [progression, character]
parents:
  - [[entity_character]]
  - [[requirement_customer_player_profile]]
dependents:
  - [[uc_progression_stat_allocation]]
---
# Character Progression Rule

## INTENT
Governs how character attributes improve after participating in a successful game using the V2 CP Point-Buy System, enforcing mathematical balance.

## THE RULE / LOGIC
- **Post-Win Reward:** After each game win, the player's account gains exactly +10 Character Points (CP) to its allowed progression cap.
- **Point-Buy System Constraints:** 
  - **Global Cap:** A character's total `spent_cp` on upgrades MUST NOT exceed `100 + (total_wins * 10)`.
  - **Non-Negativity:** No attribute is allowed to have a negative value.
- **Attribute Costs (Standard):**
  - HP (+1): Costs 1 CP
  - Attack (+1): Costs 5 CP
  - Defense (+1): Costs 5 CP
  - Movement (+1): Costs 30 CP
- **Attribute Costs (Exotic - Planned):**
  - Critical Chance (+1%): 10 CP
  - Critical Multiplier (+5%): 5 CP
  - Jump Height (+1): 15 CP
- **Unrestricted Spend:** The old V1 movement restriction (once every 5 wins) is entirely removed because Movement inherently self-balances via its 30 CP cost.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_progression]]`
- **Test Names:** `TestPostWinStatAllocation`, `TestMovementProgressionRestriction`, `TestGlobalAttributeCap`

## EXPECTATION (For Testing)
- Character wins a game -> Gains 10 CP to their maximum allowed progression cap.
- Player assigns 1 Attack to a character -> The upgrade costs 5 CP; operation is successful if `spent_cp + 5 <= 100 + (total_wins * 10)`.
- Player attempts to increase Movement -> The upgrade costs 30 CP; operation succeeds if CP cap allows it. No win limits per movement apply.
- Player tries to upgrade attributes exceeding the allowed total CP cap -> Operation rejected.
