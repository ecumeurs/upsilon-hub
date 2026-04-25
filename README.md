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
- **[trigger_all_tests.sh](trigger_all_tests.sh)**: Executes the full local test suite (E2E + Edge Cases) against the running local stack using the `--local` flag.

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
All fundamental mechanics, structural constraints, entities, and network rules that form the game are housed individually within the project-specific `docs/` folders (e.g., `upsilonapi/docs/`, `upsilonbattle/docs/`) governed by the ATD Workspace. These Atoms serve as the uncompromising basis for evaluating developer implementation logic.

## V2 Development: Tactical RPG Evolution

**Status:** In Development - Target Q3 2026

UpsilonBattle V2 represents a comprehensive evolution transforming the tactical RPG foundation with skill systems, time-based mechanics, credit economy, equipment progression, and enhanced AI.

### 🎯 Major V2 Features

**Skill System Overhaul:** Mathematical Skill Weight (SW) system with I-V grading, skill selection at creation and every 10 levels, and skill reforging mechanics.

**Time-Based Mechanics:** Channeling (pre-execution delay), temporary entities for traps/area effects, and multi-entity cell support enabling complex temporal strategies.

**Credit Economy:** 1 HP = 1 credit earning system with support credits for damage mitigation and status effects, plus shop system for skills and equipment.

**Equipment System:** 3-slot inventory (armor, utility, weapon) with weapon-as-skills transforming basic attacks into property-rich combat actions.

**AI Enhancement:** Four archetypes (Fighter, Ranger, Support, Sneak) following player progression with team composition constraints.

**Backstabbing:** 150% damage multiplier with 50% armor penetration for attacks from behind, rewarding positional tactics.

**Stat System Redesign:** x10 baseline stats (HP 30-50, Attack 10, Defense 5, Movement 3) with 100 CP point-buy system enabling meaningful percentage modifiers.

### 📋 Implementation Roadmap

- **Phase 1 (Weeks 1-4):** Foundation Systems - Skill Weight, Time-Based Mechanics, Grid Updates, Database/API
- **Phase 2 (Weeks 5-8):** Core Gameplay - Skill Selection, Channeling, Backstabbing, Credit Earning
- **Phase 3 (Weeks 9-12):** Equipment & Economy - Equipment System, Shop System, Extended Character Sheet
- **Phase 4 (Weeks 13-16):** AI Enhancement - Archetype Implementation, Progression Integration, Balancing
- **Phase 5 (Weeks 17-20):** Polish & Testing - UI Integration, Visual Feedback, Comprehensive Testing, Documentation

### 🏗️ Architectural Breakthroughs

**Unified Temporary Entity System:** All time-based mechanics (channeling, traps, area effects) represented as temporary entities with controllers.

**Skill Weight Mathematical Framework:** Net SW = 0 balance principle with precise benefit/cost calculations enabling automatic grading and pricing.

**x10 Stat Scaling:** Critical fix making percentage modifiers meaningful and character progression strategically diverse.

### 📊 Current Status

**Active Development Issues:**
- [ISS-065](issues/ISS-065_20260422_skill_weight_grading_system.md) - Skill Weight & Grading System
- [ISS-066](issues/ISS-066_20260422_time_based_mechanics.md) - Time-Based Mechanics & Temporary Entity System
- [ISS-067](issues/ISS-067_20260422_credit_economy_shop.md) - Credit Economy & Shop System
- [ISS-068](issues/ISS-068_20260422_equipment_system.md) - Equipment System & Weapon-as-Skill
- [ISS-069](issues/ISS-069_20260422_ai_archetype_enhancement.md) - AI Archetype Enhancement
- [ISS-070](issues/ISS-070_20260422_backstabbing_mechanics.md) - Backstabbing Mechanics
- [ISS-071](issues/ISS-071_20260422_starting_stats_progression.md) - Starting Stats & Character Progression Redesign
- [ISS-073](issues/ISS-073_20260423_roguelike_skill_system.md) - Roguelike Skill System: Inventory, Slots & Equipment
- [ISS-074](issues/ISS-074_20260423_simple_shop_inventory.md) - Simple Shop Inventory
- [ISS-075](issues/ISS-075_20260423_player_inventory.md) - Player Inventory System
- [ISS-076](issues/ISS-076_20260423_character_data_transfer.md) - Character Data Transfer for Battle Engine

**Comprehensive V2 Documentation:** See [v2_milestone.md](v2_milestone.md) for complete feature specifications, implementation details, and technical architecture.

## Open Issues

| Name | Date | Status | Severity | Oneliner |
|---|---|---|---|---|
| [Cross-stack error handling harmonization](issues/ISS-081_20260425_cross_stack_error_handling.md) | 2026-04-25 | Open | Medium | `error_key` is currently propagated only on the engine action paths (`POST /g... |
| [ATD for `error_key` taxonomy and possible envelope promotion](issues/ISS-080_20260425_error_key_atd_and_envelope.md) | 2026-04-25 | Open | Medium | `error_key` is now plumbed end-to-end (engine ruler → upsilonapi handler → La... |
| [Standardize cell access on Y-major layout](issues/ISS-079_20260424_cell_access_y_major_standard.md) | 2026-04-24 | Open | Medium | The tactical grid is currently serialized width-major (`cells[x][y]`) by the ... |
| [Shielding Credit Attribution System](issues/ISS-078_20260423_shielding_credit_attribution.md) | 2026-04-23 | Open | Medium | Design and implement a robust system for attributing credits earned through d... |
| [Skill Inspection Command](issues/ISS-077_20260423_skill_inspection.md) | 2026-04-23 | Open | Medium | Implement skill inspection functionality allowing players to view detailed sk... |
| [Comprehensive Item System - Shop, Inventory, Equipment & Battle Integration](issues/ISS-074_20260423_comprehensive_item_system.md) | 2026-04-23 | Open | High | Implement end-to-end item system for V2: fixed shop catalog, normalized playe... |
| [Roguelike Skill System - Inventory, Slots & Equipment](issues/ISS-073_20260423_roguelike_skill_system.md) | 2026-04-23 | Open | High | Implement comprehensive roguelike-style skill system with character skill inv... |
| [Player Choosing Facing Direction on Pass](issues/ISS-072_20260423_pass_choose_facing.md) | 2026-04-23 | Open | Medium | When a player chooses to "Pass" their turn, they must be given the option to ... |
| [Backstabbing Mechanics & Armor Penetration](issues/ISS-070_20260422_backstabbing_mechanics.md) | 2026-04-22 | Open | Medium | Implement backstabbing combat mechanic with 150% damage multiplier and 50% ar... |
| [AI Archetype Enhancement & Progression](issues/ISS-069_20260422_ai_archetype_enhancement.md) | 2026-04-22 | Open | Medium | Enhance AI system with four distinct archetypes (Fighter, Ranger, Support, Sn... |
| [Credit Economy & Shop System](issues/ISS-067_20260422_credit_economy_shop.md) | 2026-04-22 | Open | High | Implement comprehensive credit economy with multiple earning mechanisms (dama... |
| [Dead Entities Considered Obstacles](issues/ISS-059_20260420_dead_entities_obstacle_risk.md) | 2026-04-20 | Open | High | Dead entities (HP <= 0) are incorrectly treated as obstacles on both the fron... |
| [Entity Spawning Overlap](issues/ISS-058_20260420_entity_spawn_overlap.md) | 2026-04-20 | Open | Medium | In some cases, multiple entities are spawned on the same tile at the start of... |
| [Actor Message Type Validation](issues/ISS-055_20260420_actor_message_validation.md) | 2026-04-20 | Open | Low | The `Actor` implementation should validate if the target message is of the co... |
| [Game Resurrection from Board State](issues/ISS-054_20260420_game_resurrection_board_state.md) | 2026-04-20 | Open | Medium | The frontend needs a mechanism to attempt "game resurrection" from a persiste... |
| [Modernize Actor Library with Go Generics (Templates)](issues/ISS-049_20260418_actor_generics_modernization.md) | 2026-04-18 | Open | Low (Architectural Improvement) | The current Actor implementation was designed before Go 1.18 (Generics). It r... |
| [Turn Start Webhook Missing When AI Goes First](issues/ISS-048_20260418_turn_start_webhook_unicast.md) | 2026-04-18 | Open | High | The `turn.started` webhook event is intermittently missing from CI test resul... |
| [BRD Compliance CI Test Suite Blockers](issues/ISS-045_20260416_brd_compliance_ci_blockers.md) | 2026-04-16 | Open | High | The implementation of automated BRD compliance tests via specialized CLI bot ... |
| [Request Traceability Non-Compliance and Gaps](issues/ISS-042_20260415_request_traceability_gaps.md) | 2026-04-15 | Open | Medium | This issue documents the systematic non-compliance with `rule_tracing_logging... |
| [Upgradable Pawn Appearance & Model System](issues/ISS-040_20260415_pawn_appearance_system.md) | 2026-04-15 | Open | Medium | Implement an upgradable "Pawn Appearance System" that allows players to custo... |
| [Holo-Emote Procedural Reaction System](issues/ISS-039_20260415_holo_emote_system.md) | 2026-04-15 | Open | Medium | Implement a "Holo-Emote System" that triggers procedural reactions (emojis/te... |
| [Standardize Board State Naming: entities -> characters](issues/ISS-036_20260414_front_board_state_entity_naming.md) | 2026-04-14 | Open | Medium | The board state structure currently uses the term "entities" for game units. ... |
| [Ensure all logs are tagged with Request ID](issues/ISS-023_20260316_logging_tag_traceability.md) | 2026-03-16 | Open | High | The system currently lacks a strictly enforced requirement to tag every log e... |

