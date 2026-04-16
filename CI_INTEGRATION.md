# CI Integration Checklist

> **Objective:** Implement a GitHub Actions CI pipeline that validates BRD compliance, runs unit tests, and executes full end-to-end battle simulations.

## Phase 1: Infrastructure

- [x] Add `/health` endpoint to Go Engine (`upsilonapi/main.go`)
- [x] Update `communication.md` with the new `/health` endpoint
- [x] Create/update ATD atom for the health endpoint (`docs/api_go_health_check.atom.md`)
- [x] Create `.env.ci` (deterministic CI environment)
- [x] Create `docker-compose.ci.yaml` (ephemeral Docker stack with healthchecks)
- [x] Update `battleui/Dockerfile` to install composer dependencies

## Phase 2: GitHub Actions Workflows

- [x] Create `.github/workflows/lint-and-build.yml` (fast compilation check)
- [x] Create `.github/workflows/unit-tests.yml` (Go + PHP unit tests)
- [x] Create `.github/workflows/e2e-battles.yml` (full stack E2E)

## Phase 3: CI Reporting & Tooling

- [x] Create `tests/ci_report.sh` (markdown summary generator)
- [x] Update `CI.md` with complete CI documentation

## Phase 4: Housekeeping

- [x] Update `.gitignore` with CI-specific entries
- [x] Verify Go build passes with health endpoint
- [ ] Verify Docker compose builds and health checks work (requires Docker)
- [ ] Push and validate workflows on GitHub

---

## BRD Compliance Test Coverage (Current → Target)

These are the BRD requirements that should be validated by CI bot scripts. Current status shows what exists today; target shows what we need to build.

| BRD § | Requirement | ATD Atom | CI Script | Status |
|---|---|---|---|---|
| 2.1 | User Onboarding | [[uc_player_registration]] | `onboard_and_match.js` | ⚠️ Partial |
| 2.2 | Admin Management | [[uc_admin_user_management]] | — | ❌ Not started |
| 2.3 | Matchmaking | [[uc_matchmaking]] | `pvp_bot_battle.js` / `pve_bot_battle.js` | ⚠️ Implicit |
| 2.4 | Combat Engine | [[uc_combat_turn]] | `run_all_battles.sh` | ✅ 4 modes |
| 2.4 | Turn Clock / Auto-Pass | [[rule_turn_clock]] | `slow_bot_battle.js` | ⚠️ Stress only |
| 2.4 | Friendly Fire | [[rule_friendly_fire]] | — | ❌ Blocked (ISS-043) |
| 2.4 | Match Resolution | [[uc_match_resolution]] | `run_all_battles.sh` | ✅ Implicit |
| 2.5 | Progression | [[rule_progression]] | `progression_check.js` | ✅ |
| 3.1 | Password Policy | [[rule_password_policy]] | `auth_security_check.js` | ✅ |
| 3.2 | GDPR Compliance | [[rule_gdpr_compliance]] | — | ❌ Not started |
| 3.2 | Data Portability | [[api_profile_export]] | — | ❌ Not started |
| 3.4 | API-First / Help | [[requirement_customer_api_first]] | — | ❌ Not started |
| 4 | Leaderboard | [[api_leaderboard]] | — | ❌ Not started |

> [!NOTE]
> Expanding the BRD compliance test suite is a separate effort tracked here for visibility. The initial CI pipeline will run the existing tests (✅ and ⚠️).
