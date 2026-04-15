## Index

| Ref     | File                                                                                                                 | Severity | Status   | Summary                                                                         |
| ------- | -------------------------------------------------------------------------------------------------------------------- | -------- | -------- | ------------------------------------------------------------------------------- |
| ISS-038 | [ISS-038_20260415_action_feedback_protocol.md](ISS-038_20260415_action_feedback_protocol.md) | Medium | Resolved | Currently, when a player (AI or human) takes an action, the system broadcasts the new game state... |
| ISS-019 | [ISS-019_20260316_battleui_api_todo_list.md](ISS-019_20260316_battleui_api_todo_list.md) | Medium | Open | Consolidate and track BattleUI API TODOs and gaps. |
| ISS-018 | [20260312_match_participant_access_control.md](20260312_match_participant_access_control.md) | Critical | Open | Lack of authorization checks on battle match interactions. |
| ISS-017 | [20260312_action_player_id_usurpation.md](20260312_action_player_id_usurpation.md) | Critical | Open | Security risk: player_id in actions can be usurped by any auth user. |
| ISS-016 | [20260312_character_upgrade_constraints.md](20260312_character_upgrade_constraints.md) | High | Open | Character upgrades bypass progression win/point constraints. |
| ISS-015 | [20260312_matchmaking_trigger_failure.md](20260312_matchmaking_trigger_failure.md) | High | Open | 2v2 and PVE matchmaking fails to trigger startArena correctly. |
| ISS-010 | [ISS-010_20260311_ruler_readiness_logic.md](ISS-010_20260311_ruler_readiness_logic.md)                               | Low      | Open     | Enhance Ruler readiness trigger to verify grid and entities initialization.     |
| ISS-009 | [ISS-009_20260311_ruler_ownership_bypass.md](ISS-009_20260311_ruler_ownership_bypass.md)                             | Low      | Open     | Ruler ownership bypass in bridge.go and public GameState exposure.              |
| ISS-008 | [ISS-008_20260306_websocket_private_channel_transition.md](ISS-008_20260306_websocket_private_channel_transition.md) | High     | Open     | Transition WebSocket events from public to private channels for authentication. |
| ISS-005 | [ISS-005_20260305_laravel_websockets.md](ISS-005_20260305_laravel_websockets.md)                                     | High     | Open     | Implement real-time WebSocket communication layer.                              |
| ISS-003 | [ISS-003_20260305_upsilonbattle_missing_teams.md](ISS-003_20260305_upsilonbattle_missing_teams.md)                   | Medium   | Open     | Team management is missing from the battle logic.                               |
| ISS-001 | [ISS-001_20260305_upsilonbattle_mechanics_gap.md](ISS-001_20260305_upsilonbattle_mechanics_gap.md)                   | Medium   | Open     | Inconsistencies in battle mechanics implementation.                             |
| ISS-021 | [ISS-021_20260316_sanctum_token_sliding_ttl.md](ISS-021_20260316_sanctum_token_sliding_ttl.md)                       | Medium   | Open     | Missing 15-minute sliding TTL for Sanctum tokens.                              |
| ISS-022 | [ISS-022_20260316_battleui_error_handling.md](ISS-022_20260316_battleui_error_handling.md)                             | High     | Open     | Improper major error handling (HTML instead of JSON).                          |
| ISS-023 | [ISS-023_20260316_logging_tag_traceability.md](ISS-023_20260316_logging_tag_traceability.md) | High | Open | Ensure all logs are tagged with Request ID for cross-service tracing. |
| ISS-024 | [Ref_20260316_battleui_api_responder_inconsistency.md](Ref_20260316_battleui_api_responder_inconsistency.md) | Medium | Open | BattleUI ApiResponder Inconsistency and Underuse. |
| ISS-025 | [ISS-025_20260408_dashboard_hub_implementation.md](ISS-025_20260408_dashboard_hub_implementation.md) | Medium | Open | Dashboard Hub Implementation for Upsilon Battle. |
| ISS-026 | [ISS-026_20260409_api_journey_tester_cli.md](ISS-026_20260409_api_journey_tester_cli.md) | Medium | Open | Upsilon API Journey Explorer & Tester CLI. |
| ISS-027 | [ISS-027_20260409_upsiloncli_scripting_support.md](ISS-027_20260409_upsiloncli_scripting_support.md) | Medium | Open | UpsilonCLI Scripting & Automated Scenario Support. |
| ISS-035 | [ISS-035_20260413_websocket_boardstate_privacy_leak.md](ISS-035_20260413_websocket_boardstate_privacy_leak.md) | High | Open | WebSocket `board.updated` event leaks enemy character details. |
| ISS-036 | [ISS-036_20260414_front_board_state_entity_naming.md](ISS-036_20260414_front_board_state_entity_naming.md) | Medium | Open | Standardize naming from 'entities' to 'characters' in frontend state. |
| ISS-037 | [ISS-037_20260414_standardize_win_condition_team.md](ISS-037_20260414_standardize_win_condition_team.md) | Medium | Open | Standardize win condition: use team ID instead of 'winner_is_self'. |
