---
id: rule_progression
human_name: Character Progression Rule
type: RULE
version: 1.1
status: STABLE
priority: CORE
tags: [progression, character]
parents:
  - [[entity_character]]
dependents: []
---

# Character Progression Rule

## INTENT
Governs how character attributes improve after participating in a successful game and enforces upper bounds on power.

## THE RULE / LOGIC
- **Post-Win Reward:** After each game win, the player can allocate exactly 1 attribute point to a character in their roster.
- **Attribute Constraints:** 
  - **Global Cap:** The sum of all attributes (HP + Attack + Defense + Movement) MUST NOT exceed `10 + total_wins`.
  - **Non-Negativity:** No attribute is allowed to have a negative value.
- **Movement Restriction:** A point can only be allocated to the Movement attribute once every 5 accumulated wins.
  - `Current Movement <= Initial Movement + floor(Total Wins / 5)`

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_progression]]`
- **Test Names:** `TestPostWinStatAllocation`, `TestMovementProgressionRestriction`, `TestGlobalAttributeCap`

## EXPECTATION (For Testing)
- Character wins a game -> Gains 1 point.
- Player assigns point to HP -> HP increases by 1 IF total sum <= 10 + total_wins.
- Player tries to assign point to Movement after 3 wins -> Operation rejected.
- Player tries to assign point that exceeds `10 + total_wins` -> Operation rejected.
