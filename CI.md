# CI Testing Framework

This document outlines the automated verification strategy for Upsilon Battle, ensuring that code changes remain compliant with the Business Requirement Document (BRD) and Atomic Traceable Documentation (ATD).

## Architecture

The CI pipeline is split into three GitHub Actions workflows with increasing scope and cost:

| Workflow | Trigger | Duration | Purpose |
|---|---|---|---|
| **Lint & Build** | Push to `main` + PRs | ~2 min | Go syntax checks, compilation |
| **Unit Tests** | Push to `main` + PRs | ~5 min | Go + PHP isolated tests |
| **E2E Battles** | Push to `main` + PRs | ~15 min | Full stack integration |

### Infrastructure

| Component | File | Purpose |
|---|---|---|
| `.env.ci` | CI environment variables | Deterministic config for ephemeral stack |
| `docker-compose.ci.yaml` | CI Docker Compose | Ephemeral stack with healthchecks |
| `tests/ci_report.sh` | Report generator | Markdown summary of CI results |

The ephemeral Docker stack uses:
- **PostgreSQL** with `tmpfs` (RAM disk) for fast I/O
- **Docker healthchecks** on all services (`--wait` flag)
- **`GET /health`** endpoint on the Go engine for readiness probing

---

## CI Strategy

### Tier 1: Lint & Build (`lint-and-build.yml`)

Runs `go vet ./...` and builds both `upsilonapi` and `upsiloncli` to ensure compilation succeeds.

### Tier 2: Unit Tests (`unit-tests.yml`)

**Go Tests:** Runs `go test -json ./...` across all workspace modules. Results are uploaded as artifacts.

**PHP Tests:** Builds the `battleui` Docker image, then runs PHPUnit inside the container using SQLite in-memory. Tests requiring the Go engine are excluded (group `engine-required`) and deferred to E2E.

### Tier 3: E2E Battles (`e2e-battles.yml`)

Boots the full Docker stack, seeds the database, builds the CLI, then executes:

1. **BRD Compliance Scenarios** (bot scripts with assertions)
2. **Battle Mode Battery** (4 game modes via `run_all_battles.sh`)
3. **CI Report Generation** (markdown summary)

---

## Verification Scenarios

### BRD Compliance Scripts

Each scenario targets a specific BRD requirement and its linked ATD atom.

#### 1. Character Progression Lifecycle
*   **BRD:** 2.5 — Character Progression
*   **Target:** [[rule_progression]]
*   **Path:** `upsiloncli/samples/progression_check.js`
*   **Description:** Simulates a player journey from registration through a match win and attempts to upgrade character stats.
*   **Assertions:**
    *   Stat gain is allowed after a win.
    *   Stat gain is rejected if it exceeds `10 + wins`.
    *   Movement gain is rejected if not on a 5-win milestone.

#### 2. Authentication & Security Policy
*   **BRD:** 3.1 — Password Policy
*   **Target:** [[rule_password_policy]]
*   **Path:** `upsiloncli/samples/auth_security_check.js`
*   **Description:** Attempts various registration payloads to verify server-side validation.
*   **Assertions:**
    *   Reject passwords < 15 characters.
    *   Reject passwords without numbers/symbols.
    *   Accept compliant passwords.

### Battle Engine Battery

Executed via `upsiloncli/tests/run_all_battles.sh`:

| Mode | Agents | Script | Verification |
|---|---|---|---|
| `1v1_PVE` | 1 | `pvp_bot_battle.js` | Natural conclusion + log parser |
| `2v2_PVE` | 2 | `pvp_bot_battle.js` | Natural conclusion + log parser |
| `1v1_PVP` | 2 | `pvp_bot_battle.js` | Natural conclusion + log parser |
| `2v2_PVP` | 4 | `pvp_bot_battle.js` | Natural conclusion + log parser |

---

## Running Locally

### Full E2E Stack
```bash
# Prepare environment
cp .env.ci .env

# Boot Docker stack (waits for health checks)
docker compose -f docker-compose.ci.yaml up -d --wait

# Seed database
docker compose -f docker-compose.ci.yaml exec -T app php artisan migrate:fresh --seed --force

# Build CLI
cd upsiloncli && go build -o ../bin/upsiloncli ./cmd/upsiloncli && cd ..

# Run scenarios
cd upsiloncli
UPSILON_BASE_URL=http://localhost:8000 REVERB_HOST=127.0.0.1:8080 REVERB_APP_KEY=ci_app_key \
  ../bin/upsiloncli --farm samples/auth_security_check.js --timeout 60

# Run all battles
./tests/run_all_battles.sh

# Teardown
cd .. && docker compose -f docker-compose.ci.yaml down -v
```

### Unit Tests Only
```bash
# Go tests
go test ./... -count=1 -timeout 120s

# PHP tests (requires Docker)
docker build -t battleui-ci ./battleui
docker run --rm -e APP_ENV=testing -e DB_CONNECTION=sqlite -e DB_DATABASE=:memory: \
  battleui-ci sh -c "php artisan key:generate --force && php artisan test"
```

---

## Adding a New CI Test

1. **Create a bot script** in `upsiloncli/samples/` following the `bootstrapBot` pattern (see `scripting.md`).
2. **Add `@spec-link` tags** to the script header linking to the target ATD atom.
3. **Add a step** to `.github/workflows/e2e-battles.yml` with descriptive BRD reference.
4. **Update this document** and `CI_INTEGRATION.md` with the new scenario.

---

## Troubleshooting CI Failures

### Service Won't Start
Check Docker logs: `docker compose -f docker-compose.ci.yaml logs <service>`

### Health Check Timeout
- **app:** Laravel may need longer `start_period` (migrations, composer autoload).
- **engine:** Verify `GET /health` responds: `curl http://localhost:8081/health`
- **ws:** Reverb may fail if Reverb keys are missing from `.env`.

### Bot Script Timeout
- Increase `--timeout` value in the workflow step.
- Check if the game mode requires more agents than specified.
- Verify WebSocket connectivity: bots need `REVERB_HOST` and `REVERB_APP_KEY`.

> [!IMPORTANT]
> **Exit Codes:** The CLI exit code reflects the success or failure of the farm. Any assertion failure in a script results in a non-zero exit code, blocking the CI pipeline.
