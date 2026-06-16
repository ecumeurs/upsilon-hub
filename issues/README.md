# Issues

> **Index reconciled 2026-06-16** during the Principal Architect audit
> (`reporting/audit/`). Links and statuses below reflect the **actual files** in
> this directory. The project-root `README.md` carries an auto-generated table of
> *active* issues via `issues --update-readme`; **this file is the full index**
> (open + resolved) and is maintained by hand.

## Open / Active

| Ref | File | Severity | Status | Summary |
|---|---|---|---|---|
| ISS-099 | [ISS-099_20260513_skill_template_zone_support_gap.md](ISS-099_20260513_skill_template_zone_support_gap.md) | Medium | Open | Go engine `ZoneProperty` only supports Single/Neighbours; silently drops Circle/Square AoE |
| ISS-098 | [ISS-098_20260513_battle_engine_internal_user_id_leak.md](ISS-098_20260513_battle_engine_internal_user_id_leak.md) | High | Open | Raw internal User UUID leaked as `player_id` in battle-state DTOs |
| ISS-096 | [ISS-096_20260510_trap_trigger_enforcement.md](ISS-096_20260510_trap_trigger_enforcement.md) | Medium | Open | Traps without a `TriggerType` fail silently (no log/error) |
| ISS-095 | [ISS-095_20260507_skill_rework.md](ISS-095_20260507_skill_rework.md) | Medium | Open | Skill generation/execution framework design inconsistencies |
| ISS-094 | [ISS-094_20260501_atd_layer_testing_protocol.md](ISS-094_20260501_atd_layer_testing_protocol.md) | Low | Open — Planned Feature | ATD Layer Testing Protocol ("The Naive LLM") |
| ISS-093 | [ISS-093_20260429_admin_account_self_destruction_risk.md](ISS-093_20260429_admin_account_self_destruction_risk.md) | Critical | Open | Admin can anonymize/lock out own account; cascading CI failures |
| ISS-090 | [ISS-090_20260427_action_endpoint_segregation.md](ISS-090_20260427_action_endpoint_segregation.md) | Medium | Open | All tactical actions funneled through one endpoint; needs segregation |
| ISS-089 | [ISS-089_20260426_mechanic_random_shop_algorithm.md](ISS-089_20260426_mechanic_random_shop_algorithm.md) | Medium | Open | Deterministic daily rotating shop algorithm |
| ISS-087 | [ISS-087_20260426_grid_generator_tuning.md](ISS-087_20260426_grid_generator_tuning.md) | Medium | Open | Generated battle maps consistently mis-sized / mis-densified |
| ISS-083 | [ISS-083_20260425_automate_api_help_endpoints.md](ISS-083_20260425_automate_api_help_endpoints.md) | Medium | Open | API help endpoints rely on fragile reflection + atom parsing |
| ISS-081 | [ISS-081_20260425_cross_stack_error_handling.md](ISS-081_20260425_cross_stack_error_handling.md) | Medium | Open | `error_key` only propagated on engine action paths; harmonize cross-stack |
| ISS-080 | [ISS-080_20260425_error_key_atd_and_envelope.md](ISS-080_20260425_error_key_atd_and_envelope.md) | Medium | Open | ATD for `error_key` taxonomy; possible promotion to envelope root |
| ISS-079 | [ISS-079_20260424_cell_access_y_major_standard.md](ISS-079_20260424_cell_access_y_major_standard.md) | Medium | Open | Standardize cell access on Y-major layout via shared helper |
| ISS-078 | [ISS-078_20260423_shielding_credit_attribution.md](ISS-078_20260423_shielding_credit_attribution.md) | Medium | Open | Robust credit attribution for damage mitigation (shield caster) |
| ISS-077 | [ISS-077_20260423_skill_inspection.md](ISS-077_20260423_skill_inspection.md) | Medium | Open | Skill inspection UI/CLI for detailed skill properties |
| ISS-072 | [ISS-072_20260423_pass_choose_facing.md](ISS-072_20260423_pass_choose_facing.md) | Medium | Open | "Pass" action should let player choose facing (anti-backstab) |
| ISS-055 | [ISS-055_20260420_actor_message_validation.md](ISS-055_20260420_actor_message_validation.md) | Low | Open | Actor should validate target message type |
| ISS-049 | [ISS-049_20260418_actor_generics_modernization.md](ISS-049_20260418_actor_generics_modernization.md) | Low | Open | Modernize actor library with Go generics |
| ISS-042 | [ISS-042_20260415_request_traceability_gaps.md](ISS-042_20260415_request_traceability_gaps.md) | Medium | Open | Systematic non-compliance with `rule_tracing_logging` |
| ISS-040 | [ISS-040_20260415_pawn_appearance_system.md](ISS-040_20260415_pawn_appearance_system.md) | Medium | Open | Upgradable pawn appearance customization system |
| ISS-039 | [ISS-039_20260415_holo_emote_system.md](ISS-039_20260415_holo_emote_system.md) | Medium | Open | Holo-Emote System: procedural reactions above units |
| ISS-036 | [ISS-036_20260414_front_board_state_entity_naming.md](ISS-036_20260414_front_board_state_entity_naming.md) | Medium | Open | Board state "entities" naming inconsistency on the front end |
| ISS-023 | [ISS-023_20260316_logging_tag_traceability.md](ISS-023_20260316_logging_tag_traceability.md) | High | Open | No enforced requirement to tag every log entry with its trace |

## Resolved

| Ref | File | Severity | Status | Summary |
|---|---|---|---|---|
| ISS-097 | [ISS-097_20260511_actor_stop_panic_race.md](ISS-097_20260511_actor_stop_panic_race.md) | Critical | Resolved | `Actor.Stop()` non-idempotent → close-of-closed-channel panic |
| ISS-092 | [ISS-092_20260428_api_skill_property_sync_failure.md](ISS-092_20260428_api_skill_property_sync_failure.md) | High | Resolved | Bridge/engine out of sync on complex skill properties |
| ISS-091 | [ISS-091_20260428_engine_movement_obstacle_validation_bypass.md](ISS-091_20260428_engine_movement_obstacle_validation_bypass.md) | High | Resolved | Movement validation bypass on non-ground cells |
| ISS-088 | [ISS-088_20260426_credit_economy_payload_mismatch.md](ISS-088_20260426_credit_economy_payload_mismatch.md) | Medium | Resolved | `e2e_credit_economy.js` misread attacker vs target state |
| ISS-086 | [ISS-086_20260426_skill_item_registry_admin.md](ISS-086_20260426_skill_item_registry_admin.md) | High | Resolved | DB-backed skill/item registries with admin CRUD |
| ISS-085 | [ISS-085_20260425_extract_properties_shared_library.md](ISS-085_20260425_extract_properties_shared_library.md) | Medium | Resolved | Extract property system + skill weight to shared library |
| ISS-084 | [ISS-084_20260425_component_split_effects_plan.md](ISS-084_20260425_component_split_effects_plan.md) | Medium | Resolved | Split arena UI components + restore visual effects |
| ISS-082 | [ISS-082_20260425_frontend_playwright_test_seams.md](ISS-082_20260425_frontend_playwright_test_seams.md) | Medium | Resolved | Front-end Playwright suite + component-isolation seams |
| ISS-074 | [ISS-074_20260423_comprehensive_item_system.md](ISS-074_20260423_comprehensive_item_system.md) | High | Resolved | End-to-end item system: shop, inventory, equipment, battle |
| ISS-073 | [ISS-073_20260423_roguelike_skill_system.md](ISS-073_20260423_roguelike_skill_system.md) | High | Resolved | Roguelike skill inventory + slot progression |
| ISS-071 | [ISS-071_20260422_starting_stats_progression.md](ISS-071_20260422_starting_stats_progression.md) | High | Resolved | V2 starting stats + 100 CP point-buy progression |
| ISS-070 | [ISS-070_20260422_backstabbing_mechanics.md](ISS-070_20260422_backstabbing_mechanics.md) | Medium | Resolved | Backstabbing: 150% damage, 50% armor penetration |
| ISS-069 | [ISS-069_20260422_ai_archetype_enhancement.md](ISS-069_20260422_ai_archetype_enhancement.md) | Medium | Resolved | Four AI archetypes (Fighter/Ranger/Support/Sneak) |
| ISS-067 | [ISS-067_20260422_credit_economy_shop.md](ISS-067_20260422_credit_economy_shop.md) | High | Resolved | Credit economy with multiple earning mechanisms |
| ISS-066 | [ISS-066_20260422_time_based_mechanics.md](ISS-066_20260422_time_based_mechanics.md) | High | Resolved | Time-based mechanics: channeling, temporary entities |
| ISS-065 | [ISS-065_20260422_skill_weight_grading_system.md](ISS-065_20260422_skill_weight_grading_system.md) | High | Resolved | Mathematical Skill Weight (SW) system + grading |
| ISS-054 | [ISS-054_20260420_game_resurrection_board_state.md](ISS-054_20260420_game_resurrection_board_state.md) | Medium | Resolved | Game resurrection from persisted board state |

## Companion reports

| Doc | Description |
|---|---|
| [ISS-054_investigation_report.md](ISS-054_investigation_report.md) | Investigation: game resurrection from board state |

---

### Index integrity notes (2026-06-16 audit)
- The previous index linked dead `Ref_*` filenames (renamed to `ISS-NNN_*`) and
  several resolved issues whose files no longer exist (ISS-046/047/050–053/063),
  while omitting ~20 issues that do exist. All links above are verified present.
- **ISS-046** (turner hands turn to dead entity) was referenced by the old index
  but has **no file**; its tracking is lost. Not recreated here (per audit scope).
- The `issues` CLI builds links for the **root** README from each file's real
  name; this directory index is maintained by hand and was the stale artifact.
