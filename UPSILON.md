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

## Technical Architecture Overview
The system relies on a strictly separated logic implementation:

1. **Frontend (`BattleUI` - Laravel / Vue.js / Tailwind):**
   - Operates as the user-facing client.
   - Manages Player Sessions natively, distributing and securing gameplay boundaries via stateless **JWT Authentication**.
   - Orchestrates Player Queuing (`1v1 PVE, 1v1 PVP, 2V2 PVE, 2V2 PVP`) and matches clients cleanly before instantiating the combat sequence.
   - Renders the Global Leaderboard tracking Win/Loss volumes and derived ratio metrics.

2. **Backend (`UpsilonBattle` - Go JSON API):**
   - The isolated, calculating brain behind active skirmishes.
   - Fully governs the math of active battles (HP reduction, board coordinate generation, initiative delay math, and step validation).
   - Entirely ignores matchmaking queues, interacting strictly through validated combat payloads.

3. **Journey Explorer CLI ([UpsilonCLI](file:///workspace/upsiloncli) - Go):**
   - An interactive terminal tool for rapid API exploration and verification.
   - Provides full transparency by logging the equivalent `curl` command and pretty-printed JSON response for every action.
   - Includes an **Autopilot mode** (`--auto`) to simulate and verify the complete developer journey from registration to combat cleanup.
   - Integrates real-time WebSocket monitoring for in-terminal tactical board visualization.

4. **Database (PostgreSQL):**
   - Persistent, serialized memory holding Player access credentials, individual Character state logs, match resolutions, and leaderboard calculations.

## Development & Monitoring
The project includes a suite of scripts at the root to manage and monitor the service stack during development:

- **[start_services.sh](start_services.sh)**: Launches the Laravel API, Reverb Server, Vue Frontend, and Upsilon Engine in the background. Tracks PIDs and log file mappings.
- **[stop_services.sh](stop_services.sh)**: Gracefully stops all tracked services.
- **[watch_services.go](watch_services.go)**: Real-time TUI dashboard for monitoring CPU/Mem and log errors. Run with `go run watch_services.go`.
- **[check_services.sh](check_services.sh)**: Lightweight status utility for quick health checks (useful for agents and CI).

## Specification (ATD) Maps
All fundamental mechanics, structural constraints, entities, and network rules that form the game are housed individually in `/workspace/docs/`. These Atoms serve as the uncompromising basis for evaluating developer implementation logic.

## Open Issues

| Name | Date | Status | Severity | Oneliner |
|---|---|---|---|---|
| [Unified Scripting Lifecycle and CI Testing Framework](issues/ISS-044_20260415_scripting_unified_lifecycle_and_ci.md) | 2026-04-15 | Open | High | The current scripting environment in `upsiloncli` requires boilerplate code f... |
| [Friendly Fire Rule Enforcement Missing](issues/ISS-043_20260415_rule_friendly_fire_not_enforced.md) | 2026-04-15 | Open | High | The `Friendly Immunity Rule` (`rule_friendly_fire`) defined in the architectu... |
| [Request Traceability Non-Compliance and Gaps](issues/ISS-042_20260415_request_traceability_gaps.md) | 2026-04-15 | Open | Medium | This issue documents the systematic non-compliance with `rule_tracing_logging... |
| [Upgradable Pawn Appearance & Model System](issues/ISS-040_20260415_pawn_appearance_system.md) | 2026-04-15 | Open | Medium | Implement an upgradable "Pawn Appearance System" that allows players to custo... |
| [Holo-Emote Procedural Reaction System](issues/ISS-039_20260415_holo_emote_system.md) | 2026-04-15 | Open | Medium | Implement a "Holo-Emote System" that triggers procedural reactions (emojis/te... |
| [Standardize Board State Naming: entities -> characters](issues/ISS-036_20260414_front_board_state_entity_naming.md) | 2026-04-14 | Open | Medium | The board state structure currently uses the term "entities" for game units. ... |
| [Internal ID Exposure in Public APIs](issues/ISS-034_20260413_id_exposure.md) | 2026-04-13 | Open | Medium | Internal database UUIDs are currently being emitted directly to front-end and... |
| [Ensure all logs are tagged with Request ID](issues/ISS-023_20260316_logging_tag_traceability.md) | 2026-03-16 | Open | High | The system currently lacks a strictly enforced requirement to tag every log e... |
| [Security Risk: Lack of Match Participant Access Control](issues/ISS-018_20260312_match_participant_access_control.md) | 2026-03-12 | Open | Critical | Currently, any authenticated user can attempt to act or view the state of ANY... |
| [Arena not destroyed on battle end](issues/ISS-012_20260311_arena_destruction_leak.md) | 2026-03-11 | Open | Medium | Arenas are added to the `ArenaBridge.arenas` map during startup but are never... |
| [Ruler readiness trigger enhancements](issues/ISS-010_20260311_ruler_readiness_logic.md) | 2026-03-11 | Open | Low | The current readiness trigger for the `Ruler` (the `BattleStart` notification... |
| [Ruler ownership bypass in bridge.go and public GameState](issues/ISS-009_20260311_ruler_ownership_bypass.md) | 2026-03-11 | Open | Low | In `bridge.go`'s `StartArena` function, the `Ruler`'s ownership of game resou... |

