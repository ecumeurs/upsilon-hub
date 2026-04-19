# UpsilonBattle: Tactical RPG

**UpsilonBattle** is a simple, turn-based Tactical RPG (TRPG) designed for concurrent multiplayer combat and AI skirmishes. Designed entirely around the Atomic Documentation (ATD) framework, its concepts separates game logic, mechanics, and architectural boundaries into standalone, single-responsibility specifications.

## The Game At a Glance
- **Play Modes:** 1v1 (PvE or PvP) and 2v2 (PvE or PvP).
- **The Board:** A randomly generated rectangular grid (5-15 tiles per dimension, absolute minimum area of 50 tiles) containing up to 10% randomly placed impassable obstacles.
- **Victory Condition:** Eliminate all characters on the opposing team. **Friendly Fire is strictly disabled.**
- **The Roster:** Every player commands a roster of exactly 3 characters.

## Character System & Progression
- **Initial Core Roll:** New characters start with base stats (3 HP, 1 Move, 1 Attack, 1 Def) and exactly 4 additional points are randomly dispatched (Total 10 points). 
- **The Reroll Mechanic:** During account registration, players are granted an option to completely re-randomize their 3 initial character stat blocks. This reroll can be executed a strict maximum of **3 times**.
- **Stat Progression:** Securing a match victory rewards a player with 1 Attribute Point. 
  - **Constraints:** Total attributes cannot exceed `10 + total wins`. Matches result in a maximum of 1 point gain.
  - **Movement Restriction:** Upgrading the Movement attribute is heavily throttled and locked to once every 5 accumulated wins.

## Combat Mechanics
- **Initiative & Delay:** Turn order is non-linear. Characters roll a pre-initiative value ranging from `1-500`. Active turns fire when the ticker hits `0`. 
- **Action Economy:** During a turn, a character may perform a maximum of **1 Move** (`+20/tile`), **1 Attack** (`+100`), or safely **Pass** (`+300`). Performing actions accumulates a numerical "Delay Cost," mathematically extending the wait time until that character's next sequence.
- **The Shot Clock:** Active combat turns mandate a strict **30-second limit** per character. Failing to confirm an action manually results in an auto-pass forced by the server, accompanied by a penalty of `+100` (Total `+400` delay).

## Modular Architecture
The UpsilonBattle ecosystem is built as a modular multi-repo system. Each core component is maintained in its own repository and integrated into this main project as a formal **Git Submodule**.

### Repository Structure
1. **Frontend (`battleui`)**:
   - Built with Laravel, Vue.js, and Tailwind CSS.
   - Manages user sessions, JWT authentication, and player matchmaking.
   - Provides the visual interface for combat and the global leaderboard.

2. **Backend API (`upsilonapi`)**:
   - A high-performance Go JSON API.
   - Handles account management, character statistics, and match state persistence.

3. **Battle Engine (`upsilonbattle`)**:
   - The "calculating brain" that governs active combat sequences.
   - Mathematically simulates initiative, movement validation, and damage systems.

4. **Journey Explorer CLI (`upsiloncli`)**:
   - An interactive terminal tool for API exploration and verification.
   - Supports "Autopilot" sessions to simulate full player journeys.

5. **Shared Assets & Utilities**:
   - `upsilonmapdata`: Geometric board data and obstacle definitions.
   - `upsilonmapmaker`: Procedural generation tools for game boards.
   - `upsilonserializer` & `upsilontools`: Shared logic for binary/JSON serialization and common TRPG utilities.

## Setup

The Upsilon project is a complex ecosystem. For a detailed guide on how to prepare your environment, install dependencies, and configure the system, please refer to the **[Setup Documentation](Setup.md)**.

## Getting Started

### Cloning the Project
Since the project relies on submodules, you must clone recursively to fetch all components:
```bash
git clone --recursive git@github.com:ecumeurs/upsilon-hub.git
```
If you have already cloned the repository, initialize the submodules with:
```bash
git submodule update --init --recursive
```

### Development & Monitoring

#### DevContainer Environment
The project provides a pre-configured development environment via **[.devcontainer/](.devcontainer/)**. This is the recommended way to develop for Upsilon, ensuring a consistent environment across all platforms.

- **Stack:** PHP 8.4, Go, Node.js 20.
- **Tools:** Includes Composer, Postgres client, and ATD integration tools.
- **Port Forwarding:**
  - `8000`: Laravel App
  - `8080`: Reverb (WebSockets)
  - `8081`: Upsilon Engine (Go API)
  - `5173`: Vue Frontend (HMR)

#### Service Management
All services are standardized on the `main` branch. The project includes a suite of root-level scripts for service management:

- **[start_services.sh](start_services.sh)**: Launches the Laravel API, Reverb Server, Vue Frontend, and Upsilon Engine in the background. 
- **[stop_services.sh](stop_services.sh)**: Gracefully stops all tracked services.
- **[watch_services.go](watch_services.go)**: Real-time TUI dashboard for monitoring CPU/Mem and error logs.
- **[check_services.sh](check_services.sh)**: Lightweight status utility for quick health checks.

## Continuous Integration & Quality control

UpsilonBattle employs a robust CI/CD pipeline via GitHub Actions to ensure code quality, architectural integrity, and business rule compliance.

### Automated Workflows (`.github/workflows/`)
- **[Unit Tests](.github/workflows/unit-tests.yml)**: Runs Go unit tests for all backend modules and PHP unit tests for the Laravel frontend.
- **[Lint & Build](.github/workflows/lint-and-build.yml)**: Performs static analysis (Go vet) and verifies that all core components and Docker images build successfully.
- **[E2E Battle Tests](.github/workflows/e2e-battles.yml)**: Orchestrates a full ephemeral Docker stack to run integration tests. It uses specialized CLI bots to simulate real player journeys and verify complex game mechanics.

### CI Infrastructure (Docker Stack)
The project utilizes a dedicated **[docker-compose.ci.yaml](docker-compose.ci.yaml)** to spin up an ephemeral testing environment. This stack is optimized for speed and reliability.

- **Components:**
  - `db`: Postgres 18-alpine database.
  - `db-init`: Migration and seeding service.
  - `app`: The Laravel application.
  - `ws`: Reverb WebSocket server.
  - `engine`: The Upsilon Battle Engine (Go).
  - `tester`: The Upsilon CLI running in integration mode.
- **Usage:**
  ```bash
  docker compose -f docker-compose.ci.yaml up -d --wait
  ```

### Compliance & Reporting
- **BRD Compliance**: Automated scripts (located in `upsiloncli/samples/`) verify critical Business Requirements such as **Password Policy** (`[[rule_password_policy]]`) and **Character Progression** (`[[rule_progression]]`).
- **CI Reports**: Each run generates a summary report ([ci_report.sh](tests/ci_report.sh)) that is attached to the job summary, providing immediate visibility into test outcomes and compliance status.

## Specification (ATD) Maps
All fundamental mechanics, structural constraints, entities, and network rules that form the game are housed individually in `/workspace/docs/`. These Atoms serve as the uncompromising basis for evaluating developer implementation logic.

## Open Issues

| Name | Date | Status | Severity | Oneliner |
|---|---|---|---|---|
| [Implement High-Performance Manual Pagination and Search for Admin Tools](issues/ISS-053_20260419_admin_performance_pagination.md) | 2026-04-19 | Resolved | Medium | Optimized manual cursor pagination for admin registries |
| [Implement Full E2E Admin Test Suite](issues/ISS-052_20260419_admin_history_management_testing.md) | 2026-04-19 | Resolved | Medium | Full E2E Admin Test Suite implementation |
| [Implement Administrative Match History Management (Backend & Frontend)](issues/ISS-051_20260419_admin_history_management_impl.md) | 2026-04-19 | Resolved | Medium | Administrative Match History Implementation (API & UI) |
| [Modernize Actor Library with Go Generics (Templates)](issues/ISS-049_20260418_actor_generics_modernization.md) | 2026-04-18 | Open | Low (Architectural Improvement) | The current Actor implementation was designed before Go 1.18 (Generics). It r... |
| [Turn Start Webhook Missing When AI Goes First](issues/ISS-048_20260418_turn_start_webhook_unicast.md) | 2026-04-18 | Open | High | The `turn.started` webhook event is intermittently missing from CI test resul... |
| [BRD Compliance CI Test Suite Blockers](issues/ISS-045_20260416_brd_compliance_ci_blockers.md) | 2026-04-16 | Open | High | The implementation of automated BRD compliance tests via specialized CLI bot ... |
| [Request Traceability Non-Compliance and Gaps](issues/ISS-042_20260415_request_traceability_gaps.md) | 2026-04-15 | Open | Medium | This issue documents the systematic non-compliance with `rule_tracing_logging... |
| [Upgradable Pawn Appearance & Model System](issues/ISS-040_20260415_pawn_appearance_system.md) | 2026-04-15 | Open | Medium | Implement an upgradable "Pawn Appearance System" that allows players to custo... |
| [Holo-Emote Procedural Reaction System](issues/ISS-039_20260415_holo_emote_system.md) | 2026-04-15 | Open | Medium | Implement a "Holo-Emote System" that triggers procedural reactions (emojis/te... |
| [Standardize Board State Naming: entities -> characters](issues/ISS-036_20260414_front_board_state_entity_naming.md) | 2026-04-14 | Open | Medium | The board state structure currently uses the term "entities" for game units. ... |
| [Internal ID Exposure in Public APIs](issues/ISS-034_20260413_id_exposure.md) | 2026-04-13 | Open | Medium | Internal database UUIDs are currently being emitted directly to front-end and... |
| [Ensure all logs are tagged with Request ID](issues/ISS-023_20260316_logging_tag_traceability.md) | 2026-03-16 | Open | High | The system currently lacks a strictly enforced requirement to tag every log e... |
| [Security Risk: Lack of Match Participant Access Control](issues/ISS-018_20260312_match_participant_access_control.md) | 2026-03-12 | Open | Critical | Currently, any authenticated user can attempt to act or view the state of ANY... |
| [Arena not destroyed on battle end](issues/ISS-012_20260311_arena_destruction_leak.md) | 2026-03-11 | Open | Medium | Arenas are added to the `ArenaBridge.arenas` map during startup but are never... |

