# Upsilon Hub: Project Map & Infrastructure

## 1. Project Background
Upsilon Hub is a multi-stack ecosystem designed for high-performance battle engine simulation and management. It bridges a Go-based core logic with a Laravel management interface and a flexible CLI for automation and testing.

## 2. Who's Who (Service Architecture)

| Component | Stack | Role | Default Port |
|---|---|---|---|
| **battleui** | Laravel / Vue | Management UI & Database Orchestration | `8080` |
| **upsilonapi** | Go | The "Bridge" - provides API/WS access to the engine | `8081` |
| **upsilonbattle** | Go | Core Battle Engine (BattleArena) | N/A (Embedded) |
| **upsiloncli** | JS / Python | Scripting, E2E Testing, and Bot Orchestration | N/A |
| **upsilontypes** | Go | Shared type definitions and serialization | N/A |

## 3. Folder Organization
- `/battleui`: The Laravel web application.
- `/upsilonapi`: The Go bridge server.
- `/upsilonbattle`: Core engine logic (BattleArena).
- `/upsiloncli`: Command-line tools and E2E test scenarios.
- `/upsilontypes`: Shared Go structures and engine types.
- `/scripts`: Operational shell scripts for dev, CI, and deployment.
- `/docs`: Shared ATD documentation (Atoms).
- `/issues`: Tracked project risks and technical debt.

## 4. Testing Toolkit

### Core Scripts
- `scripts/trigger_one_ci_test.sh <path>`: Run a specific E2E test scenario (e.g. `tests/scenarios/edge_movement_entity_collision.js`).
- `scripts/run_all_unit_tests.sh`: Executes all Go and JS unit tests.
- `scripts/check_services.sh`: Health check for running docker/local services.

### Manual Verification
- `upsiloncli --local --farm`: Starts a local match with automated bot players.
- `scripts/upsilon_log_parser.sh`: Parses and colorizes engine logs for debugging.

## 5. Environments
- **Dev**: Local execution via `go run` or `php artisan serve`.
- **DevContainer**: Fully containerized environment (Standard for development).
- **CI**: Automated testing via GitHub Actions and `docker-compose.ci.yaml`.
- **Prod**: High-availability deployment via `docker-compose.prod.yaml`.
