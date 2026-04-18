# CI Testing Framework

This document outlines the automated verification strategy for Upsilon Battle, ensuring that code changes remain compliant with the Business Requirement Document (BRD) and Atomic Traceable Documentation (ATD).

## Architecture

The CI pipeline is split into three GitHub Actions workflows with increasing scope:

| Workflow | Trigger | Purpose |
|---|---|---|
| **Lint & Build** | Push + PR | Go syntax checks, compilation |
| **Unit Tests** | Push + PR | Go + PHP isolated tests |
| **E2E Battles** | Push + PR | Full stack integration & Customer Scenarios |

### Infrastructure

| Component | File | Purpose |
|---|---|---|
| `.env.ci` | CI environment variables | Deterministic config for ephemeral stack |
| `docker-compose.ci.yaml` | CI Docker Compose | Ephemeral stack with healthchecks |
| `tests/run_all_scenarios.sh` | Scenario Runner | **Centralized discovery & execution** |
| `tests/ci_report.sh` | Report generator | Markdown summary of CI results |

---

## Customer Requirement Mapping

The following scenarios map directly to the **Conformity Matrix** and validate specific customer-facing requirements. All scripts are located in `upsiloncli/tests/scenarios/`.

| ID | Scenario Name | Primary ATOM | Script |
|---|---|---|---|
| **CR-01** | Complete New Player Onboarding | `[[uc_player_registration]]` | `e2e_customer_onboarding.js` |
| **CR-02** | Player Login & Session Management | `[[uc_player_login]]` | `e2e_customer_login.js` |
| **CR-03** | Character Reroll Mechanics | `[[us_character_reroll]]` | `e2e_character_reroll.js` |
| **CR-04** | Matchmaking Flow (PvE Instant) | `[[uc_matchmaking]]` | `e2e_matchmaking_pve_instant.js` |
| **CR-05** | Matchmaking Flow (PvP Queue) | `[[uc_matchmaking]]` | `e2e_matchmaking_pvp_queue.js` |
| **CR-06** | Combat Turn Management | `[[uc_combat_turn]]` | `e2e_combat_turn_management.js` |
| **CR-07** | Friendly Fire Prevention | `[[rule_friendly_fire]]` | `e2e_friendly_fire_prevention.js` |
| **CR-08** | Match Resolution (Standard) | `[[uc_match_resolution]]` | `e2e_match_resolution_standard.js` |
| **CR-09** | Match Resolution (Forfeit) | `[[uc_match_resolution]]` | `e2e_match_resolution_forfeit.js` |
| **CR-10** | Character Progression (Post-Win) | `[[uc_progression_stat_allocation]]` | `e2e_progression_post_win.js` |
| **CR-11** | Progression Constraints | `[[rule_progression]]` | `e2e_progression_constraints.js` |
| **CR-12** | Leaderboard Viewing | `[[us_leaderboard_view]]` | `e2e_leaderboard_viewing.js` |
| **CR-13** | Password Policy Enforcement | `[[rule_password_policy]]` | `e2e_password_policy.js` |
| **CR-14** | GDPR Data Portability | `[[api_profile_export]]` | `e2e_gdpr_portability.js` |
| **CR-15** | Admin User Management | `[[uc_admin_user_management]]` | `e2e_admin_user_management.js` |
| **CR-16** | Session Timeout Handling | `[[requirement_req_ui_session_timeout]]` | `e2e_session_timeout.js` |
| **CR-17** | API Self-Discovery | `[[requirement_customer_api_first]]` | `e2e_api_discovery.js` |

---

## E2E Testing Strategy (`e2e-battles.yml`)

Instead of individual workflow steps, all end-to-end verifications are centralized in the `upsiloncli/tests/scenarios/` directory.

### 1. Centralized Runner (`run_all_scenarios.sh`)
The runner automatically discovers all `e2e_*.js` scripts. It handles:
- **Agent Coordination**: Determines required agent count based on filename (e.g., `pvp` or `combat` triggers 2 agents).
- **Execution**: Runs the `upsiloncli --farm` command.
- **Reporting Contract**: Appends `[SCENARIO_RESULT: PASSED]` to the log file upon successful exit (code 0).

### 2. Scenario Library
All customer-facing scenarios from the **Conformity Matrix** are implemented here:
- `e2e_customer_onboarding.js` (CR-01)
- `e2e_customer_login.js` (CR-02)
- `e2e_character_reroll.js` (CR-03)
- ... (Total of 17 scenarios)

### 3. CI Report Generation
The `tests/ci_report.sh` script parses the logs in `upsiloncli/tests/logs/`. It uses a **unified detection method**: it only marks a test as `✅ PASS` if the success marker `[SCENARIO_RESULT: PASSED]` is present in the log.

---

## Adding a New CI Test

Adding a new verification scenario is a **zero-touch process** (no GitHub Actions YAML edits required):

1.  **Create the Script**: Add a new JavaScript file in `upsiloncli/tests/scenarios/`.
    - **Naming**: Use the `e2e_` prefix. If it requires 2 agents (PVP/Combat), include `pvp` or `combat` in the name.
2.  **Add Assertions**: Use the `upsilon` JS API helpers:
    - `upsilon.assert(condition, msg)`
    - `upsilon.assertEquals(actual, expected, msg)`
3.  **Tag with ATD**: Include `@spec-link [[atom_id]]` in the header for traceability.
4.  **Save & Push**: The CI runner will automatically find your script and include it in the next run.

---

## Running Locally

### Full E2E Stack
```bash
# 1. Boot Docker stack
docker compose -f docker-compose.ci.yaml up -d --wait

# 2. Build CLI
cd upsiloncli && go build -o bin/upsiloncli cmd/upsiloncli/main.go && cd ..

# 3. Run the Suite
docker compose -f docker-compose.ci.yaml exec tester /bin/sh ./tests/run_all_scenarios.sh

# 4. Generate Report
./tests/ci_report.sh > ci_report.md
```

---

## Troubleshooting

- **Mismatch between Runner and Report**: If the runner says `[PASSED]` but the report says `❌ FAIL`, check if the script is exiting with code 0 but missing the log output. The reporter relies strictly on the `[SCENARIO_RESULT: PASSED]` marker written by the runner.
- **Service Timeouts**: If healthchecks fail, check `docker compose logs`. Laravel often requires more time to boot in resource-constrained environments.
- **Ghost Tests**: Ensure your script filename starts with `e2e_`, otherwise the runner will skip it.
