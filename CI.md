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
