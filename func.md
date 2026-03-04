# TRPG Functional Specification

*Generated via ATD Synthesis*

## 1. Overview
UpsilonBattle is a simple Tactical RPG focusing on fast-paced, turn-based grid combat. The goal is to provide seamless matchmaking that drops players directly into active skirmishes.

## 2. User Flows
### 2.1 Onboarding & Registration
- Users land on the promotional public landing page (`ui_landing`).
- To play, users must create an account providing only an **Account Name** and **Password** (`ui_registration`). Email addresses are prohibited by specification.
- **Reroll Phase:** Upon account generation, users are immediately granted 3 characters (`entity_player`). If unsatisfied with the randomly generated attributes, the player may trigger a full-roster reroll up to 3 maximum times (`mech_character_reroll`).

### 2.2 Dashboard & Matchmaking
- After logging in (secured by JWT), users arrive at the `ui_dashboard` displaying their roster stats, Win/Loss ratio, and a link to the `ui_leaderboard`.
- Players choose between 4 matchmaking queues (`req_matchmaking`):
  1. 1v1 PVE
  2. 1v1 PVP
  3. 2V2 PVE
  4. 2V2 PVP
- PvP selections hold the user in a `ui_waiting_room` until paring is confirmed. PvE queues instantly spawn the user onto the game board.

## 3. The Core Game Loop
### 3.1 Combat Instantiation
- Matches are strictly 1v1 or 2v2 (`spec_match_format`).
- The board generates dynamically between 5x5 and 15x15 (minimum 50 tiles) with up to 10% obstacles (`mech_board_generation`).
- Each unit rolls random pre-initiative numbers between 1 and 1000 (`mech_initiative`).

### 3.2 Action Rules & Strategy
- Characters only act when their running 'Delay Ticker' mathematically evaluates to `0`. 
- **The Shot Clock:** Active players have **30 seconds** to formulate a turn (`mech_action_economy`). If they do not act, the game forces an auto-pass with a harsh +100 delay penalty.
- **Valid Actions:** A unit may Move (once per turn), Attack (once per turn), or simply Pass. Actions incur numerical delay costs to recalculate their next turn time.
- **Friendly Fire:** Damage cannot be applied to allied units (`rule_friendly_fire`).

### 3.3 Resolution & Progression
- The match ends when an entire team has been reduced to 0 HP.
- Losing players gain nothing.
- Winning players gain **1 Attribute Point** to increase an existing character's HP, Attack, or Defense (`rule_progression`). Movement requires 5 saved wins to upgrade.
