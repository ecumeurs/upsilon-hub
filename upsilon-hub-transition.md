# Upsilon-Hub ATD Migration Plan: Centralized to Per-Project Documentation

## Executive Summary

This document outlines the transition plan for moving Upsilon-Hub's ATD documentation from a **centralized `docs/` folder** to **per-project `docs/` folders**. This migration aligns the documentation structure with the modular architecture of the UpsilonBattle ecosystem, enabling better organization, clearer ownership, and improved multi-project workspace support.

**Current State:** 291 ATD atoms in a single `docs/` folder serving all subprojects.
**Target State:** Each project has its own `docs/` folder with relevant atoms, using cross-project references for shared concepts.

**Migration Scope:** All 291 ATD atoms to be distributed across 9 project directories.

---

## 1. Current State Analysis

### 1.1 Project Structure

```
upsilon-hub/
├── docs/                          # Centralized (291 atoms)
├── upsilonapi/                    # Go API backend
├── upsilonbattle/                 # Game engine
├── battleui/                      # Laravel + Vue frontend
├── upsiloncli/                    # CLI tools
├── upsilonmapdata/                # Grid data
├── upsilonmapmaker/               # Board generation
├── upsilonserializer/             # Serialization
└── upsilontools/                  # Shared utilities
```

### 1.2 Atom Distribution (291 total)

| Project Prefix | Count | Examples |
|---|---|---|
| `api_*` | 29 | api_auth_login, api_matchmaking, api_profile_character |
| `mech_*`, `mechanic_*`, `mec_*` | 68 | mech_initiative, mech_combat_attack_computation, mec_skill_reforging_mechanic |
| `ui_*`, `battleui_*` | 45 | ui_login, ui_battle_arena, battleui_api_dtos |
| `entity_*` | 11 | entity_character, entity_player, entity_game_match |
| `uc_*`, `us_*` | 38 | uc_player_login, us_character_reroll, us_win_progression |
| `req_*`, `rule_*` | 34 | req_security, rule_password_policy, rule_progression |
| `module_*` | 15 | module_backend, module_frontend, module_game |
| `domain_*`, `arch_*` | 10 | domain_upsilon_engine, arch_api_id_masking_gateway |
| `data_*`, `infra_*` | 5 | data_persistence, infra_mvp_docker |
| `spec_*` | 3 | spec_match_format |
| `script_*`, `watch_*` | 4 | script_farm, watch_services |
| Other/Temp | 29 | temp, temp, temp... (needs cleanup) |

### 1.3 Cross-Project Dependencies

Based on atom naming and structure analysis, the following cross-project dependencies exist:

- **Frontend → Backend:** `ui_*` atoms reference `api_*` atoms
- **Backend → Engine:** `api_*` atoms reference `mech_*` atoms
- **All → Shared:** Business rules (`rule_*`) and requirements (`req_*`) apply across multiple projects
- **CLI → API:** Use cases (`uc_*`, `us_*`) reference API endpoints

---

## 2. Target State Design

### 2.1 Proposed Project Distribution

| Target Project | Atom Prefixes | Estimated Count | Shared Concepts |
|---|---|---|---|
| **upsilonapi** | `api_*`, `arch_*`, `domain_*`, `data_*`, `infra_*`, `module_upsilonapi` | ~55 | req_security, req_matchmaking, rule_* (business rules) |
| **upsilonbattle** | `mech_*`, `mechanic_*`, `mec_*`, `entity_*`, `module_backend_*`, `module_game` | ~95 | rule_* (game rules), req_* (game requirements) |
| **battleui** | `ui_*`, `battleui_*`, `module_frontend_*`, `module_ui_*` | ~55 | req_ui_*, req_player_experience |
| **upsiloncli** | `uc_*`, `us_*`, `script_*`, `watch_*`, `usecase_*` | ~45 | (references to api_*, req_*, rule_*) |
| **upsilonmapdata** | `entity_grid` (if applicable), map-specific atoms | ~2 | (new atoms may be needed) |
| **upsilonmapmaker** | `mech_board_generation_*` | ~5 | (references to upsilonmapdata) |
| **upsilonserializer** | (create new atoms if needed) | ~0 | (likely no atoms needed) |
| **upsilontools** | (shared utility atoms if any) | ~0 | (likely no atoms needed) |

### 2.2 Cross-Project Reference Format

```yaml
---
id: ui_login_form
parents:
  - [[upsilonapi:api_auth_login]]      # Cross-project reference
  - [[upsilonapi:api_auth_register]]
  - [[battleui:req_security_token_ttl]] # Same-project reference
type: UI
layer: ARCHITECTURE
---
```

### 2.3 Shared/Common Atoms Strategy

For business rules and requirements that apply across multiple projects:

**Option A: Duplicate with Cross-References** (Recommended)
- Place the primary atom in the most relevant project
- Other projects reference it via cross-project prefix
- Example: `rule_password_policy` lives in `upsilonapi`, referenced by `battleui` and `upsiloncli`

**Option B: Shared Workspace Folder**
- Create `docs/shared/` for truly cross-cutting atoms
- Use sparingly, only for atoms that have no clear "home" project

---

## 3. Migration Strategy

### 3.1 Migration Phases

**Phase 1: Preparation (Week 1)**
- [ ] Create migration mapping document (atom → target project)
- [ ] Set up branch protection for migration work
- [ ] Create backup of current `docs/` folder
- [ ] Update ATD tooling to support cross-project references (if needed)
- [ ] Prepare test suite to verify post-migration integrity

**Phase 2: Project-Specific Migrations (Weeks 2-4)**

| Week | Projects | Focus |
|---|---|---|
| 2 | upsilonapi, upsilonbattle | Core backend systems |
| 3 | battleui, upsiloncli | Frontend and tooling |
| 4 | upsilonmapdata, upsilonmapmaker, others | Support projects |

**Phase 3: Cross-Project Reference Updates (Week 5)**
- [ ] Update all atom references to use project prefixes where needed
- [ ] Verify `atd weave` works correctly across projects
- [ ] Update `@spec-link` tags in code if paths have changed
- [ ] Run full test suite

**Phase 4: Validation & Cleanup (Week 6)**
- [ ] Run `atd lint` on all projects
- [ ] Run `atd audit` on all projects
- [ ] Verify all `@spec-link` and `@test-link` tags resolve
- [ ] Remove old centralized `docs/` folder
- [ ] Update CI/CD pipelines
- [ ] Update documentation (README.md, CLAUDE.md)

### 3.2 Atom-by-Atom Migration Process

For each atom:

1. **Identify Target Project:**
   ```bash
   # Use naming convention as primary indicator
   # Verify by checking @spec-link tags in related code
   ```

2. **Create Project Structure:**
   ```bash
   mkdir -p upsilonapi/docs
   mkdir -p upsilonbattle/docs
   mkdir -p battleui/docs
   mkdir -p upsiloncli/docs
   # ... for each project
   ```

3. **Move Atom File:**
   ```bash
   mv docs/mech_initiative.atom.md upsilonbattle/docs/
   ```

4. **Update References:**
   - Scan atom's `parents` and `dependents` for cross-project references
   - Add project prefix to references that now live in different projects
   - Update `@spec-link` tags in code files if necessary

5. **Create/Update Project `.atd` Config:**
   ```json
   {
     "docs_path": "docs/",
     "code_paths": ["./"],
     "bloating_factor": {...},
     "model": "llama3.2"
   }
   ```

6. **Verify Integrity:**
   ```bash
   cd upsilonbattle
   atd weave
   atd lint
   atd stats
   ```

### 3.3 Handling Special Cases

**Temp/Duplicate Atoms:**
- Clean up `temp` atoms (7 instances found) before migration
- Review duplicates and merge or delete as appropriate

**Orphaned Atoms:**
- Identify atoms with no `@spec-link` tags
- Determine if they should be deleted or if links are missing

**Breakage During Migration:**
- Create a rollback branch before each phase
- Document all broken references for manual review

---

## 4. Detailed Migration Mapping

### 4.1 upsilonapi (~55 atoms)

| Pattern | Example | Count |
|---|---|---|
| `api_*` | api_auth_login, api_matchmaking, api_profile_character | 29 |
| `arch_*` | arch_api_id_masking_gateway | 1 |
| `domain_*` | domain_upsilon_engine, domain_skill_system | 6 |
| `data_*` | data_persistence | 1 |
| `infra_*` | infra_mvp_docker, infra_seed_admin | 2 |
| `module_*` | module_upsilonapi | 1 |
| Shared refs | rule_password_policy, req_security, etc. | ~15 (via references) |

### 4.2 upsilonbattle (~95 atoms)

| Pattern | Example | Count |
|---|---|---|
| `mech_*` | mech_initiative, mech_combat_attack_computation | ~45 |
| `mechanic_*` | mechanic_backstab_detection_algorithm | ~12 |
| `mec_*` | mec_ai_archetype_system, mec_channeling_mechanic | ~11 |
| `entity_*` | entity_character, entity_player, entity_game_match | 11 |
| `module_backend_*` | module_backend_action_economy, module_backend_combat_math | 5 |
| `module_game` | module_game | 1 |
| `mech_board_generation_*` | mech_board_generation_board_dimensions | 3 |
| Shared refs | rule_friendly_fire, req_matchmaking, etc. | ~7 (via references) |

### 4.3 battleui (~55 atoms)

| Pattern | Example | Count |
|---|---|---|
| `ui_*` | ui_login, ui_battle_arena, ui_character_roster | ~45 |
| `battleui_*` | battleui_api_dtos, battleui_upsilon_api_service | 2 |
| `module_frontend_*` | module_frontend_matchmaking_orchestration | 4 |
| `module_ui_*` | module_ui_tactical_layout | 1 |
| Shared refs | req_ui_look_and_feel, req_player_experience | ~3 (via references) |

### 4.4 upsiloncli (~45 atoms)

| Pattern | Example | Count |
|---|---|---|
| `uc_*` | uc_player_login, uc_matchmaking, uc_combat_turn | 8 |
| `us_*` | us_character_reroll, us_win_progression, us_leaderboard_view | 16 |
| `usecase_*` | usecase_api_flow_game_turn, usecase_api_flow_matchmaking | 2 |
| `script_*` | script_farm | 1 |
| `watch_*` | watch_services | 1 |
| Cross-refs | (all reference api_*, req_*, rule_*) | ~17 (via references) |

### 4.5 Shared/Common Atoms (~40 atoms)

These atoms apply across multiple projects and should be strategically placed:

| Atom | Primary Home | Referenced By |
|---|---|---|
| `req_security_*` | upsilonapi | battleui, upsiloncli |
| `req_matchmaking_*` | upsilonapi | upsiloncli, battleui |
| `rule_password_policy` | upsilonapi | battleui, upsiloncli |
| `rule_progression` | upsilonbattle | upsiloncli, battleui |
| `rule_friendly_fire` | upsilonbattle | upsiloncli |
| `req_player_experience` | battleui | upsiloncli |
| `spec_match_format` | upsilonbattle | upsilonapi, upsiloncli |

---

## 5. Risk Assessment & Mitigation

### 5.1 Risks

| Risk | Impact | Likelihood | Mitigation |
|---|---|---|---|
| Broken @spec-link references | High | Medium | Automated scanning + manual verification |
| Cross-project reference syntax errors | High | Medium | Implement validation in atd weave |
| Missing @spec-link tags in code | Medium | High | Pre-migration audit |
| CI/CD pipeline failures | High | Low | Update pipelines in parallel |
| Developer confusion during migration | Medium | Medium | Clear communication, freeze non-essential work |
| Loss of atom metadata | High | Low | Git version control + backup |
| Orphaned atoms after migration | Low | Medium | Post-migration orphan detection |

### 5.2 Rollback Plan

If critical issues arise during migration:

1. **Immediate Rollback:**
   ```bash
   git checkout migration-backup-point
   # Restore centralized docs/ folder
   ```

2. **Partial Rollback:**
   - Rollback only the current phase's projects
   - Keep completed phases

3. **Data Recovery:**
   - Full backup of `docs/` stored in `docs-backup-<timestamp>/`
   - Atom-to-project mapping document preserved

---

## 6. Post-Migration Benefits

### 6.1 Immediate Benefits

- **Clearer Ownership:** Each team owns their project's documentation
- **Better Organization:** Atoms live closer to the code they describe
- **Improved Workspace Support:** Native ATD workspace functionality
- **Reduced Cognitive Load:** Developers only see relevant atoms for their project

### 6.2 Long-Term Benefits

- **Easier Onboarding:** New developers can focus on their project's docs
- **Better CI/CD:** Project-specific linting and validation
- **Scalability:** Adding new projects doesn't complicate shared docs folder
- **Independent Release Cycles:** Each project's docs can evolve independently

---

## 7. Implementation Checklist

### Pre-Migration
- [ ] Create feature branch `feature/atd-per-project-migration`
- [ ] Create migration mapping document (detailed atom list)
- [ ] Backup current `docs/` folder
- [ ] Update ATD tooling to version supporting workspaces
- [ ] Notify all developers of upcoming migration
- [ ] Set up testing environment

### Migration Execution
- [ ] Phase 1: upsilonapi migration (55 atoms)
- [ ] Phase 1: upsilonbattle migration (95 atoms)
- [ ] Phase 2: battleui migration (55 atoms)
- [ ] Phase 2: upsiloncli migration (45 atoms)
- [ ] Phase 3: Support projects migration (remaining atoms)
- [ ] Update all cross-project references
- [ ] Update all @spec-link tags in code
- [ ] Run `atd weave` on all projects
- [ ] Run `atd lint` on all projects
- [ ] Run `atd audit` on all projects

### Post-Migration
- [ ] Remove old centralized `docs/` folder
- [ ] Update CI/CD pipelines
- [ ] Update README.md with new structure
- [ ] Update CLAUDE.md with workspace workflow
- [ ] Update WebUI configuration
- [ ] Train developers on new workflow
- [ ] Monitor for issues for 2 weeks
- [ ] Create final migration report

---

## 8. Appendix: Atom Reference List

### 8.1 Complete Atom List by Target Project

#### upsilonapi (55 atoms)
```
api_auth_login
api_auth_logout
api_auth_register
api_auth_user
api_battle_proxy
api_controller_methods
api_equipment_management
api_go_action_feedback
api_go_battle_action
api_go_battle_engine
api_go_battle_start
api_go_health_check
api_go_webhook_callback
api_help_endpoint
api_laravel_gateway
api_laravel_health_check
api_leaderboard
api_matchmaking
api_plan_travel_toward
api_profile_character
api_profile_export
api_request_id
api_ruler_methods
api_skill_grading_computation
api_standard_envelope
api_websocket
api_websocket_arena_updates
api_websocket_game_events
api_websocket_user_notifications
arch_api_id_masking_gateway
battleui_api_dtos
battleui_upsilon_api_service
data_persistence
domain_credit_economy
domain_ruler_state
domain_skill_system
domain_upsilon_engine
domain_upsilon_engine_domain_upsilon_engine_api_orchestration
domain_upsilon_engine_domain_upsilon_engine_combat_isolation
domain_upsilon_engine_domain_upsilon_engine_entity_stat_integration
domain_upsilon_engine_domain_upsilon_engine_expectation
domain_upsilon_engine_domain_upsilon_engine_resolution
domain_upsilon_engine_domain_upsilon_engine_technical_interface
infra_mvp_docker
infra_seed_admin
module_upsilonapi
```

#### upsilonbattle (95 atoms)
```
entity_character
entity_character_distribute_remaining_points
entity_equipment_system
entity_game_match
entity_grid
entity_match_participants
entity_player
entity_player_entity_character_rules_apply
entity_player_entity_player_initial_setup
entity_player_entity_player_registration
entity_player_entity_player_stats_tracking
entity_users
mech_action_economy
mech_action_economy_action_cost_rules
mech_action_economy_time_constraint_rules
mech_action_economy_timeout_penalty_rules
mech_actor_dispatch_loop
mech_actor_handler_context
mech_actor_lifecycle
mech_actor_pattern
mech_ai_name_generation
mech_board_generation
mech_board_generation_board_dimensions
mech_board_generation_min_area_constraint
mech_board_generation_terrain_obstacles
mech_character_reroll
mech_character_reroll_availability
mech_character_reroll_effect
mech_character_reroll_limit
mech_combat_attack_computation
mech_combat_shielding
mech_combat_standard_attack_computation
mech_controller_communication_sequence
mech_controller_handshake
mech_entity_properties
mech_entity_properties_item_properties
mech_entity_properties_skill_properties
mech_game_state_versioning
mech_initiative
mech_initiative_active_state
mech_initiative_delay_costs
mech_initiative_initiative_roll
mech_initiative_requeue_calculation
mech_matchmaking
mech_message_queue
mech_move_validation
mech_move_validation_move_validation_already_moved
mech_move_validation_move_validation_controller_mismatch
mech_move_validation_move_validation_entity_collision
mech_move_validation_move_validation_existence
mech_move_validation_move_validation_jump_limitations
mech_move_validation_move_validation_obstacle_collision
mech_move_validation_move_validation_path_adjacency
mech_move_validation_move_validation_path_length_credits
mech_move_validation_move_validation_turn_mismatch
mech_sanctum_token_renewal
mech_skill_reforging_mechanic
mech_skill_selection_progression
mech_skill_validation
mech_skill_validation_action_state_verification
mech_skill_validation_economic_cost_verification_cooldown_check
mech_skill_validation_economic_cost_verification_stat_leech
mech_skill_validation_entity_targeting_rules_verification
mech_skill_validation_existence_verification
mech_skill_validation_grid_boundaries_verification
mech_skill_validation_range_limit_verification
mech_skill_validation_turn_controller_identity_verification
mech_version_bit_packing
mech_web_catchall_router
mec_ai_archetype_system
mec_backstabbing_mechanic
mec_cell_attached_effects
mec_channeling_mechanic
mec_credit_spending_shop
mec_effect_caster_tracking
mec_equipment_stat_bonuses
mec_expiration_controller
mec_multi_entity_cell_system
mec_pre_post_execution_costs
mec_shop_inventory_system
mec_three_slot_equipment_system
mec_weapon_as_skill_system
mechanic_backstab_detection_algorithm
mechanic_cell_attached_effects
mechanic_channeling_mechanic
mechanic_character_creation_integration
mechanic_character_point_buy_system
mechanic_character_stat_allocation_rules
mechanic_effect_caster_tracking
mechanic_exotic_attribute_progression
mechanic_expiration_controller
mechanic_mech_ai_termination
mechanic_mech_arena_lifecycle
mechanic_mech_battle_startup_handshake
mechanic_mech_cli_sensitive_data_masking
mechanic_mech_frontend_auth_bridge
mechanic_mech_skill_weight_calculator
mechanic_mech_temporary_entity_system
mechanic_multi_entity_cell_system
mechanic_script_lifecycle
mechanic_shared_memory
module_backend
module_backend_action_economy
module_backend_board_generation
module_backend_combat_math
module_game
```

#### battleui (55 atoms)
```
module_frontend
module_frontend_board_ui_rendering
module_frontend_character_entity_creation
module_frontend_integration_constraint
module_frontend_matchmaking_orchestration
module_frontend_session_management
module_ui_tactical_layout
req_ui_look_and_feel
req_player_experience
ui_action_panel
ui_admin_dashboard
ui_battle_arena
ui_board
ui_character_battle_card
ui_character_pawn
ui_character_roster
ui_combat_header
ui_dashboard
ui_dashboard_match_statistics
ui_dashboard_navigation
ui_dashboard_player_statistics
ui_dashboard_profile_edit
ui_dashboard_queue_selection
ui_dashboard_roster_display
ui_dashboard_security_check
ui_holo_obstacle
ui_initiative_timeline
ui_iso_board
ui_landing
ui_leaderboard
ui_leaderboard_data_display
ui_leaderboard_metrics_displayed
ui_leaderboard_modes
ui_leaderboard_primary_sorting
ui_leaderboard_secondary_sorting
ui_leaderboard_security
ui_login
ui_modal_box
ui_registration
ui_registration_character_generation_flow
ui_registration_minimal_form_fields
ui_registration_reroll_limit
ui_registration_success_state
ui_tactical_action_report
ui_team_roster_panel
ui_theme
ui_waiting_room
```

#### upsiloncli (45 atoms)
```
req_admin_experience
req_logging_traceability
req_matchmaking
req_matchmaking_matchmaking_queue
req_matchmaking_pve_pvp_transition
req_matchmaking_transition_rules
req_security
req_security_authorization
req_security_public_access
req_security_token_exchange
req_security_token_ttl
requirement_customer_action_reporting
requirement_customer_api_first
requirement_customer_player_profile
requirement_customer_user_account
requirement_customer_user_id_privacy
requirement_req_trpg_game_definition
requirement_req_ui_session_timeout
rule_admin_access_restriction
rule_battle_readiness
rule_character_create_character
rule_character_progression_v2
rule_character_renaming
rule_credit_earning_damage
rule_credit_earning_status_effects
rule_credit_earning_support
rule_forfeit_battle
rule_friendly_fire
rule_friendly_fire_match_type
rule_friendly_fire_team_validation
rule_gdpr_compliance
rule_leaderboard_cycle
rule_leaderboard_score_calculation
rule_matchmaking_single_queue
rule_password_policy
rule_pve_winnability_balance
rule_pvp_stalemate_draw
rule_ruler_test_robustness
rule_skill_grading_system
rule_team_mechanics
rule_tracing_logging
rule_turn_atomic_selection
rule_turn_clock
script_farm
uc_admin_history_management
uc_admin_login
uc_admin_user_management
uc_auth_logout
uc_combat_turn
uc_matchmaking
uc_match_resolution
uc_player_login
uc_player_registration
uc_progression_stat_allocation
us_auth_logout
us_character_reroll
us_character_reroll_button_lockdown
us_character_reroll_create_character
us_character_reroll_reroll_button_action
us_character_reroll_reroll_counter
usecase_api_flow_game_turn
usecase_api_flow_matchmaking
us_leaderboard_view
us_leaderboard_view_auth_leaderboard
us_leaderboard_view_sort_leaderboard
us_new_player_onboard
us_queue_selection
us_queue_selection_display_win_loss_record
us_queue_selection_pve_instant_game_start
us_queue_selection_pvp_navigation
us_queue_selection_queue_buttons_visible
us_take_combat_turn
us_win_progression
us_win_progression_movement_locked
us_win_progression_progression_screen
us_win_progression_stat_reflection
us_win_progression_win_alloc_point
watch_services
```

#### Other/Shared (remaining atoms)
```
spec_match_format
spec_match_format_mode_rule
spec_match_format_team_composition_rule
temp (x7 - to be cleaned up)
```

---

## 9. Conclusion

This migration plan provides a structured approach to moving Upsilon-Hub's ATD documentation from a centralized model to a per-project model. The migration will be executed in phases with clear checkpoints, rollback plans, and validation steps to minimize risk and ensure a successful transition.

**Key Success Factors:**
1. Thorough pre-migration planning and mapping
2. Automated validation at each phase
3. Clear communication with all stakeholders
4. Comprehensive rollback plan
5. Post-migration monitoring and support

**Timeline:** 6 weeks from start to completion, with 2 weeks of post-migration monitoring.

**Resources Required:**
- 1-2 developers for execution
- ATD tooling updates (if needed)
- CI/CD pipeline updates
- Developer training materials

---

**Document Version:** 1.0
**Last Updated:** 2026-04-23
**Status:** Draft - Pending Review
