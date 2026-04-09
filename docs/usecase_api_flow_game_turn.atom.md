---
id: usecase_api_flow_game_turn
status: DRAFT
layer: CUSTOMER
priority: 3
tags: flow,game,combat,turn,api
dependents: []
human_name: Tactical Game Turn API Flow
type: USECASE
version: 1.0
parents:
  - [[requirement_customer_api_first]]
  - [[uc_combat_turn]]
---

# New Atom

## INTENT
To detail the exact API interaction sequence required for a complete tactical turn in a match.

## THE RULE / LOGIC
### Step 1: Turn Awareness
- **Action**: `GET /api/v1/game/{match_id}`
- **Validation**: Check `current_turn.entity_id` and ensuring it maps to a character under the player's control.

### Step 2: Movement Phase
- **Action**: `POST /api/v1/game/{match_id}/action`
- **Payload**: 
  ```json
  {
    "type": "MOVE",
    "params": { "path": [[x1, y1], [x2, y2]] }
  }
  ```
- **Intent**: Reposition the combatant on the grid.

### Step 3: Skill Phase
- **Action**: `POST /api/v1/game/{match_id}/action`
- **Payload**: 
  ```json
  {
    "type": "SKILL",
    "params": { 
      "skill_id": "heavy_strike",
      "target_id": "enemy_entity_01" 
    }
  }
  ```
- **Intent**: Exhaust AP to perform a combat maneuver.

### Step 4: Turn Finalization
- **Action**: `POST /api/v1/game/{match_id}/action`
- **Payload**: `{ "type": "END_TURN" }`
- **Intent**: Commit all changes and hand over control to the next character in the initiative queue.

### Step 5: State Synchronization
- **Action**: `GET /api/v1/game/{match_id}`
- **Intent**: Verify match state has been updated and observe the result of the actions.

## TECHNICAL INTERFACE
- **Related Specs:** `[[api_go_battle_action]]`, `[[api_go_battle_engine]]`
- **Code Tag:** `@spec-link [[usecase_api_flow_game_turn]]`

## EXPECTATION
- Action sequence results in correct state transition in the engine.
- Validation errors (422) are returned for illegal moves or skills.
- Turn cycle successfully transitions to the next entity after 'END_TURN'.
