# Upsilon Battle — Edge Case Test Report

| Field | Value |
|---|---|
| **Date** | 2026-04-28 08:18 UTC |
| **Commit** | `97b18b0` |
| **Branch** | `main` |

## Edge Case Test Results

| EC ID | Test Name | Result | ATD Atom |
|---|---|---|---|
### Phase 1: Movement & Authentication

| **EC-01** | Movement on Obstacle Tiles | ❌ FAIL | `[[mech_move_validation_move_validation_obstacle_collision]]` |
| **EC-02** | Movement on Entity Collision | ❌ FAIL | `[[mech_move_validation_move_validation_entity_collision]]` |
| **EC-03** | Movement Already Attacked | ❌ FAIL | `[[mech_move_validation_move_validation_already_moved]]` |
| **EC-04** | Movement Path Too Long | ✅ PASS | `[[mech_move_validation_move_validation_path_length_credits]]` |
| **EC-05** | Movement Path Not Adjacent | ✅ PASS | `[[mech_move_validation_move_validation_path_adjacency]]` |
| **EC-06** | Movement Out of Turn | ✅ PASS | `[[mech_move_validation_move_validation_turn_mismatch]]` |
| **EC-07** | Movement Wrong Controller | ❌ FAIL | `[[mech_move_validation_move_validation_controller_mismatch]]` |
| **EC-08** | Movement Grid Boundaries | ✅ PASS | `[[mech_skill_validation_grid_boundaries_verification]]` |
| **EC-09** | Movement Jump Limitations | ✅ PASS | `[[mech_move_validation_move_validation_jump_limitations]]` |
| **EC-20** | Password Policy Full Coverage | ❌ FAIL | `[[rule_password_policy]]` |
| **EC-21** | Invalid Credentials | ✅ PASS | `[[api_auth_login]]` |
| **EC-22** | Session Timeout / Expired Token | ✅ PASS | `[[requirement_req_ui_session_timeout]]` |
| **EC-23** | Missing Token | ✅ PASS | `[[req_security_authorization]]` |
| **EC-24** | Admin Non-Admin Access | ✅ PASS | `[[uc_admin_login]]` |

### Phase 2: Attack Validation

| **EC-10** | Attack Out of Turn | ✅ PASS | `[[mech_skill_validation_turn_controller_identity_verification]]` |
| **EC-11** | Attack Wrong Controller | ❌ FAIL | `[[mech_skill_validation_turn_controller_identity_verification]]` |
| **EC-12** | Attack Friendly Fire | ❌ FAIL | `[[rule_friendly_fire]]` |
| **EC-13** | Attack Target Not in Range | ✅ PASS | `[[mech_skill_validation_range_limit_verification]]` |
| **EC-14** | Attack Target Out of Grid | ✅ PASS | `[[mech_skill_validation_grid_boundaries_verification]]` |
| **EC-16** | Attack No Entity | ✅ PASS | `[[mech_combat_attack_computation]]` |
| **EC-17** | Attack Already Acted | ✅ PASS | `[[mech_skill_validation_action_state_verification]]` |
| **EC-18** | Attack Skill Cooldown | ✅ PASS | `[[mech_skill_validation_economic_cost_verification_cooldown_check]]` |
| **EC-19** | Attack Targeting Rules | ❌ FAIL | `[[mech_skill_validation_entity_targeting_rules_verification]]` |

### Phase 3: Character & Matchmaking

| **EC-25** | Character Reroll Limit | ❌ FAIL | `[[mech_character_reroll_limit]]` |
| **EC-26** | Reroll After Match | ✅ PASS | `[[mech_character_reroll_limit]]` |
| **EC-27** | Progression Without Wins | ❌ FAIL | `[[rule_progression]]` |
| **EC-28** | Progression Attribute Cap | ✅ PASS | `[[rule_progression]]` |
| **EC-29** | Progression Movement Gate | ✅ PASS | `[[rule_progression]]` |
| **EC-30** | Progression Negative Value | ✅ PASS | `[[rule_progression]]` |
| **EC-31** | Queue While Already Queued | ✅ PASS | `[[rule_matchmaking_single_queue]]` |
| **EC-32** | Queue While in Match | ❌ FAIL | `[[rule_matchmaking_single_queue]]` |
| **EC-33** | Invalid Game Mode | ✅ PASS | `[[api_matchmaking]]` |
| **EC-34** | Leave Queue Not Queued | ✅ PASS | `[[api_matchmaking]]` |

### Phase 4: Match Resolution

| **EC-35** | Forfeit Out of Turn | ❌ FAIL | `[[rule_forfeit_battle]]` |
| **EC-36** | Action After Match End | ❌ FAIL | `[[uc_match_resolution]]` |

### Phase 5: API & Communication

| **EC-37** | Missing Request ID | ✅ PASS | `[[api_request_id]]` |
| **EC-38** | Invalid UUID Format | ✅ PASS | `[[api_standard_envelope]]` |
| **EC-39** | Malformed JSON | ✅ PASS | `[[api_standard_envelope]]` |
| **EC-40** | 5xx Error Handling | ✅ PASS | `[[mechanic_mech_frontend_auth_bridge]]` |

### Phase 6: Leaderboard

| **EC-41** | Invalid Game Mode | ✅ PASS | `[[api_leaderboard]]` |
| **EC-42** | Over Pagination | ✅ PASS | `[[api_leaderboard]]` |

### Phase 7: Admin

| **EC-43** | Admin View Private Data | ✅ PASS | `[[rule_admin_access_restriction]]` |
| **EC-44** | Anonymize Non-Existent | ✅ PASS | `[[uc_admin_user_management]]` |
| **EC-45** | Soft Delete Non-Existent | ✅ PASS | `[[uc_admin_user_management]]` |

### Phase 8: WebSocket

| **EC-46** | Connection Without Token | ✅ PASS | `[[api_websocket]]` |
| **EC-47** | Wrong Channel | ✅ PASS | `[[api_websocket]]` |
| **EC-48** | Ping/pong Timeout | ✅ PASS | `[[api_websocket]]` |

## Summary Statistics

