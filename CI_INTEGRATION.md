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
| 2.1 | User Onboarding | [[uc_player_registration]] | `onboard_and_match.js` | ❌ Blocked (ISS-045) |
| 2.2 | Admin Management | [[uc_admin_user_management]] | — | ❌ Not started |
| 2.3 | Matchmaking | [[uc_matchmaking]] | `test_match_resolution.sh` | ❌ Blocked (ISS-045) |
| 2.4 | Combat Engine | [[uc_combat_turn]] | `run_all_battles.sh` | ✅ 4 modes |
| 2.4 | Turn Clock / Auto-Pass | [[rule_turn_clock]] | `slow_bot_battle.js` | ❌ Blocked (ISS-045) |
| 2.4 | Friendly Fire | [[rule_friendly_fire]] | `friendly_fire_check.js` | ❌ Blocked (ISS-045) |
| 2.4 | Match Res: Standard 1v1 | [[uc_match_resolution]] | `test_match_resolution.sh` | ❌ Blocked (ISS-045) |
| 2.4 | Match Res: Standard 2v2 | [[uc_match_resolution]] | `match_resolution_2v2.js` | ❌ Blocked (ISS-045) |
| 2.4 | Match Res: 1v1 Forfeit | [[uc_match_resolution]] | `match_resolution_forfeit.js` | ❌ Blocked (ISS-045) |
| 2.4 | Match Res: 2v2 Forfeit | [[uc_match_resolution]] | `match_resolution_forfeit_2v2.js` | ❌ Blocked (ISS-045) |
| 2.5 | Prog: Upgrade Post-Win | [[rule_progression]] | `progression_check.js` | ❌ Blocked (ISS-045) |
| 2.5 | Prog: Upgrade Pre-Win Lock | [[rule_progression]] | `progression_check.js` | ❌ Blocked (ISS-045) |
| 2.5 | Prog: Character Reroll | [[us_character_reroll]] | `reroll_check.js` | ❌ Blocked (ISS-045) |
| 3.1 | Password Policy Edge Cases | [[rule_password_policy]] | `auth_security_check.js` | ✅ |
| 3.2 | GDPR Compliance | [[rule_gdpr_compliance]] | — | ❌ Blocked (ISS-045) |
| 3.2 | Data Portability | [[api_profile_export]] | — | ❌ Blocked (ISS-045) |
| 3.4 | API-First / Help | [[requirement_customer_api_first]] | — | ❌ Not started |
| 4 | Leaderboard | [[api_leaderboard]] | — | ❌ Blocked (ISS-045) |

> [!NOTE]
> Expanding the BRD compliance test suite is a separate effort tracked here for visibility. The initial CI pipeline will run the existing tests (✅ and ⚠️).
