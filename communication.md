# Upsilon Battle: API Communication Reference

> [!TIP]
> **Live Documentation:** For the most accurate and up-to-date technical reference, including experimental endpoints and current validation rules, use the [Live API Help Endpoint](http://localhost:8000/api/v1/help).

This document provides a comprehensive reference for the communication interfaces between the Vue.js frontend, the Laravel API Gateway, and the Upsilon (Go) Battle Engine.

## 1. Shared Infrastructure

### 1.1 Standard JSON Message Envelope
**Source:** [[api_standard_envelope]]

To guarantee traceability and consistent error handling, every JSON exchange between system units (Vue, Laravel, Go) MUST conform to the following root structure:

- **Event Name:** `board.updated`
- **Payload:** Strictly follows the `[[api_standard_envelope]]` format. The tactical state is located in the `data` field of the envelope. Team-based victory is reported via `winner_team_id`. The state includes an optional `action` object providing explicit data for animations (moves, attacks, passes).
- **Envelope Example:** `{"request_id": "uuid", "success": true, "data": {"match_id": "uuid", "action": {"type": "attack", ...}, ...}}`

```json
{
  "request_id": "018f5a...", // string (UUIDv7): For chronological sortability. Rules in [[api_request_id]].
  "message": "...",         // string: Intent summary, status message, or error description.
  "success": true,          // boolean: Indicates operational success.
  "data": {},               // object|array|null: Primary payload (e.g., Resource, Collection, or Success DTO).
  "meta": {}                // object: Side information for debugging or testing (optional).
}
```

### 1.2 Request Identification
**Source:** [[api_request_id]]

The `request_id` must be a **string (UUIDv7)**. It is the responsibility of the originator (typically the Vue frontend for user actions) to generate this ID. It must be propagated across all distributed calls spanning Laravel and Go to maintain the trace defined in [[rule_tracing_logging]].

### 1.3 State Versioning & Deduplication
**Source:** [[mech_game_state_versioning]]

To ensure consistency and optimize performance during high-frequency combat, Upsilon uses a monotonic versioning system. 

1. **Versioning:** Every state mutation in the Go Engine increments a `Version` (int64).
2. **De-duplication:** The Go Internal Engine drops outgoing webhooks if the state hasn't progressed since the last transmission.
3. **Gateway Enforcement:** Laravel ignores incoming webhooks with a version lower than or equal to the current database state, effectively deduplicating the fan-out from multiple controllers. Laravel uses the `version` (int64) as the single source of truth for match progression, mapping it to both `version` and legacy `turn` columns.
4. **Broadcast Efficiency:** Clients (Vue/CLI) rely on the `version` field to ensure they are processing the latest tactical state.

### 1.4 Service Ports & Network Topology

| Service | Port (Dev) | Protocol | Role |
| :--- | :--- | :--- | :--- |
| **Laravel API** | `8000` | HTTP | External Gateway & Orchestration |
| **Reverb Server** | `8080` | WS/WSS | Tactical WebSocket Bridge |
| **Upsilon Engine** | `8081` | HTTP | Stateless Combat Engine (Internal) |
| **Vue.js (Vite)** | `5173` | HTTP | "Neon in the Dust" Frontend (Dev Only) |

> [!TIP]
> **Production Note:** In production, the Vue.js app is pre-built and served directly by the Laravel API at the same port as the web server, eliminating the need for the Vite dev server (Port 5173).

---

## 2. Laravel API (External Gateway)
**Source Module:** [[api_laravel_gateway]]  
**Base URL:** `http://localhost:8000/api/v1`  
**Authentication:** Bearer Token (Laravel Sanctum)

### 2.0 API Summary

| Verb | URI | Intent | Specification |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/register` | User Registration & Roster Creation | [[api_auth_register]] |
| `POST` | `/auth/login` | User Authentication | [[api_auth_login]] |
| `POST` | `/auth/logout` | Session Termination | [[api_auth_logout]] |
| `GET` | `/profile/characters` | List Player Roster | [[api_profile_character]] |
| `GET` | `/profile/character/{id}` | Get Character Details | [[api_profile_character]] |
| `POST` | `/profile/character/{id}/reroll` | Reset Stats (New Accounts) | [[api_profile_character]] |
| `POST` | `/profile/character/{id}/upgrade` | Attribute Point Allocation | [[api_profile_character]] |
| `POST` | `/matchmaking/join` | Enter Battle Queue | [[api_matchmaking]] |
| `GET` | `/matchmaking/status` | Poll Match Status | [[api_matchmaking]] |
| `DELETE` | `/matchmaking/leave` | Exit Battle Queue | [[api_matchmaking]] |
| `GET` | `/game/{id}` | Get Cached Board State | [[api_battle_proxy]] |
| `POST` | `/game/{id}/action` | Proxy Tactical Action to Engine | [[api_battle_proxy]] |
| `POST` | `/game/{id}/forfeit` | Standalone Forfeit Route | [[rule_forfeit_battle]] |
| `POST` | `/auth/admin/login` | Administrative Authentication (CLI/API) | [[uc_admin_login]] |
| `GET` | `/admin/dashboard` | Administrative Landing Hub | [[ui_admin_dashboard]] |
| `GET` | `/admin/users` | List Users for Auditing | [[uc_admin_user_management]] |
| `POST` | `/admin/users/{account_name}/anonymize` | GDPR Right to be Forgotten | [[uc_admin_user_management]] |
| `DELETE` | `/admin/users/{account_name}` | Administrative Soft Delete | [[uc_admin_user_management]] |
| `POST` | `/api/webhook/upsilon` | Ingest Engine State Update | [[api_go_webhook_callback]] |
| `GET` | `/leaderboard` | Global Rankings (Mode-based) | [[api_leaderboard]] |

### 2.1 Authentication

#### `POST /auth/register`
- **Specification:** [[api_auth_register]]
- **Intent:** [[uc_player_registration]]: Allow new users to create an account and receive an initial characters roster.
- **Input:**
  - `account_name`: `string`
  - `email`: `string` (must be unique)
  - `password`: `string` (minimum 15 characters)
  - `password_confirmation`: `string` (must match password)
  - `full_address`: `string` (Mandatory per [[uc_player_registration]])
  - `birth_date`: `string (ISO8601)` (Mandatory per [[uc_player_registration]])
- **Output:**
  - `user`: `UserResource` (See [[#4.4-userresource]])
  - `token`: `string` (JWT Bearer Token)

#### `POST /auth/login`
- **Specification:** [[api_auth_login]]
- **Intent:** [[uc_player_login]]: Authenticate existing users and provide a session token.
- **Input:**
  - `account_name`: `string` [Mandatory]
  - `password`: `string` [Mandatory]
- **Output:**
  - `user`: `UserResource` (See [[#4.4-userresource]])
  - `token`: `string` (JWT Bearer Token)

#### `POST /auth/admin/login`
- **Specification:** [[uc_admin_login]]
- **Intent:** Authenticate administrators for CLI or high-privilege API access.
- **Input:**
  - `account_name`: `string` [Mandatory]
  - `password`: `string` [Mandatory]
- **Validation:** Must have `Admin` role.
- **Output:**
  - `user`: `UserResource` (See [[#4.4-userresource]])
  - `token`: `string` (JWT Bearer Token)

#### `POST /auth/logout`
- **Specification:** [[api_auth_logout]]
- **Intent:** [[uc_auth_logout]]: Terminate the active session for Player or Admin and revoke the current access token.
- **Security:** Requires `auth:sanctum` middleware.
- **Output:** `null` (successful status code 200 with standard success envelope).

### 2.2 Profile & Character Management

#### `GET /profile/characters`
- **Specification:** [[api_profile_character]]
- **Intent:** List all characters associated with the authenticated player's roster.
- **Output:** `Array<CharacterResource>` (See [[#4.5-characterresource]])

#### `GET /profile/character/{characterId}`
- **Specification:** [[api_profile_character]]
- **Intent:** Retrieve detailed statistics and status for a specific character.
- **Input:**
  - `characterId`: `string (UUID)` (URL Parameter)
- **Output:** `CharacterResource` (See [[#4.5-characterresource]])

#### `POST /profile/character/{characterId}/reroll`
- **Specification:** [[api_profile_character]]
- **Intent:** [[uc_player_registration]] (Step 4): Allow fresh accounts to reroll their starting character stats.
- **Restriction:** Restricted to "New" accounts; forbidden after match participation.
- **Input:**
  - `characterId`: `string (UUID)` (URL Parameter)
- **Output:**
  - `character`: `CharacterResource` (The updated character)
  - `reroll_count`: `int` (The user's total reroll count)

#### `POST /profile/character/{characterId}/upgrade`
- **Specification:** [[api_profile_character]]
- **Intent:** [[uc_progression_stat_allocation]]: Manually allocate attribute points earned through wins.
- **Input:**
  - `characterId`: `string (UUID)` (URL Parameter)
  - `stats`: `object`
    - `hp`: `int` (increment amount, optional)
    - `attack`: `int` (increment amount, optional)
    - `defense`: `int` (increment amount, optional)
    - `movement`: `int` (increment amount, optional)
- **Validation:** Must adhere to [[rule_progression]] (Attribute Cap: `10 + wins`).
- **Output:** `CharacterResource` (The updated character)

### 2.3 Matchmaking & Queue

#### `POST /matchmaking/join`
- **Specification:** [[api_matchmaking]]
- **Intent:** [[uc_matchmaking]]: Enter the queue for a specific game mode.
- **Input:**
  - `game_mode`: `string` ("1v1_PVP", "2v2_PVP", "1v1_PVE", "2v2_PVE")
- **Output:** 
  - `status`: `string` ("queued" | "matched")
  - `match_id`: `string (UUID)|null`
  - `expected_participants`: `int`
  - `empty_slots`: `int`

#### `GET /matchmaking/status`
- **Specification:** [[api_matchmaking]]
- **Intent:** [[uc_matchmaking]]: Poll for matchmaking updates or match assignment.
- **Output:** 
  - `status`: `string` ("queued" | "matched" | "idle")
  - `match_id`: `string (UUID)|null`
  - `expected_participants`: `int|null`
  - `empty_slots`: `int|null`
  - `queued_at`: `string (ISO8601)|null`

#### `DELETE /matchmaking/leave`
- **Specification:** [[api_matchmaking]]
- **Intent:** Cancel queue entry and return to Dashboard.
- **Output:** `null` (successful status code 200 with standard success envelope).

### 2.4 Game Interaction (Proxy)

#### `GET /game/{id}`
- **Specification:** [[api_battle_proxy]]
- **Intent:** Retrieve the **cached** board state from the Laravel database.
- **Logic:** Avoids direct engine overhead by reading the last known state synced via webhook.
- **Input:**
  - `id`: `string (UUID)` (URL Parameter - Match ID)
- **Output:** `GameMatchResource` (See [[#4.6-gamematchresource]])
#### `POST /game/{id}/action`
- **Specification:** [[api_battle_proxy]]
- **Intent:** [[uc_combat_turn]]: Proxy tactical commands to the Upsilon Go Engine.
- **Input:**
  - `id`: `string (UUID)` (URL Parameter - Match ID)
  - `payload`: `ArenaActionRequest` (Note: `player_id` is automatically injected and validated by Laravel)
- **Logic:** Validates that the authenticated user owns the targeted `entity_id` before proxying the request to Upsilon `/internal/arena/:id/action` via [[api_go_battle_action]].
- **Output:** `ArenaActionResponse` (See [[#4.1-arenaactionrequest]])

#### `POST /game/{id}/forfeit`
- **Specification:** [[rule_forfeit_battle]]
- **Intent:** Allow a player to concede the match.
- **Input:**
  - `id`: `string (UUID)` (URL Parameter - Match ID)
- **Logic:** Calls standalone forfeit logic in the engine. Bypasses the need for `entity_id`.
- **Constraint:** Can only be called during a turn owned by the authenticated player (Enforced by Engine).
- **Output:** Standard Success Envelope.

### 2.5 Social & Competitive

#### `GET /leaderboard`
- **Specification:** [[api_leaderboard]]
- **Intent:** Retrieve global rankings for a specific battle mode, filtered by the current weekly cycle.
- **Input:**
  - `mode`: `string` ("1v1_PVP", "2v2_PVP", "1v1_PVE", "2v2_PVE") [Mandatory]
  - `page`: `int` (Default: 1)
- **Rules Applied:**
  - [[rule_leaderboard_score_calculation]]: Scoring formula.
  - [[rule_leaderboard_cycle]]: Temporal filter (Current Week).
- **Output:**
  - `results`: `Array<RankingResource>` (Rank, Account Name, Wins, Losses, Score) - Paginated (10 per page).
  - `self`: `RankingResource|null` (Current user's ranking context).
  - `meta`: `PaginationMeta`

---

## 3. Upsilon API (Go Internal Engine)
**Source Module:** [[api_go_battle_engine]]  
**Base URL:** `http://localhost:8081/internal`

### 3.0 API Summary

| Verb | URI | Intent | Specification |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Engine Health Check | [[api_go_health_check]] |
| `POST` | `/arena/start` | Initialize Arena Instance | [[api_go_battle_start]] |
| `POST` | `/arena/{id}/action` | Execute Combat Action | [[api_go_battle_action]] |

### 3.1 Health Check

#### `GET /health`
- **Specification:** [[api_go_health_check]]
- **Intent:** Provide a lightweight readiness probe for Docker healthchecks and CI tooling.
- **Authentication:** None (public endpoint).
- **Output:**
  - `status`: `string` ("ok")
  - `revision`: `string` (Git commit hash of the running binary)

### 3.2 Arena Life Cycle

#### `POST /arena/start`
- **Specification:** [[api_go_battle_start]]
- **Intent:** Initialize a tactical arena instance.
- **Input:** `ArenaStartRequest` (See [[#4.3-arenastartrequest]])
- **Output:** `ArenaStartResponse` (See [[#4.1-arenastartresponse]])

### 3.3 Battle Actions

#### `POST /arena/{id}/action`
- **Specification:** [[api_go_battle_action]]
- **Intent:** [[uc_combat_turn]]: Validate and execute tactical moves or attacks within a battle.
- **Input:** 
  - `id`: `string (UUID)` (URL Parameter - Arena ID)
  - `payload`: `ArenaActionRequest` (See [[#4.2-arenaactionrequest]])
- **Output:** `ArenaActionResponse` (See [[#4.1-arenaactionresponse]])

### 3.4 Asynchronous Webhook (Callback)
**Destination:** `POST /api/webhook/upsilon` (in [[api_battle_proxy]]) — Must be reachable internally from the Go Engine (e.g. `http://127.0.0.1:8000`).

#### Webhook Event Payload
- **Specification:** [[api_go_webhook_callback]]
- **Event Type (`event_type`):**
  - `game.started`: Arena initialization complete.
  - `turn.started`: New entity initiative active (starts 30s clock).
  - `board.updated`: Position or stat change (Damage/Heal/Move).
  - `game.ended`: Win condition met.
- **Data Payload:** `ArenaEvent` (See [[#4.7-arenaevent]]) which contains a `BoardState`.

---

## 4. Common Data Structures

### 4.1 Arena DTOs

#### ArenaActionResponse
- **`status`**: `string`

#### ArenaStartResponse
- **`arena_id`**: `string (UUID)`
- **`initial_state`**: `BoardState`

#### ArenaActionRequest
- **`entity_id`**: `string (UUID)`
- **`type`**: `string` ("move", "attack", "pass", "forfeit")
- **`target_coords`**: `Array<Position>`

#### ArenaStartRequest
- **`match_id`**: `string (UUID)`
- **`callback_url`**: `string` (Webhook URL - Must be reachable internally by the Go Engine)
- **`players`**: `Array<Player>`

### 4.2 Arena Components

#### BoardState
**Specification:** [[api_go_battle_engine]]

Defines the complete state of a tactical arena at a specific moment in time.

| Field | Type | Description |
| :--- | :--- | :--- |
| `players` | `Array<Player>` | Consolidated roster of participants and their live entities. |
| `grid` | `Grid` | The tactical map structure. |
| `turn` | `Array<Turn>` | Sequence of actors based on initiative. |
| `current_player_is_self` | `boolean` | **Gateway Only:** True if the current user is the acting player. Masked from Go's `current_player_id`. |
| `current_entity_id` | `string (UUID)` | ID of the entity currently acting. |
| `timeout` | `string (ISO8601)` | Timestamp when the current turn expires. |
| `start_time` | `string (ISO8601)` | Timestamp when the arena started. |
| `version` | `int64` | Monotonic sequence number for state changes. Required for deduplication. [[mech_game_state_versioning]] |
| `winner_team_id` | `int?` | ID of the winning team (if match finished). 0 or null if draw. |
| `action` | `object` | Optional. Explicit details about the last action for UI animations. |
| `version` | `int64` | State version (Bit-packed: turn << 32 \| action). |

#### Grid
- **`width`**: `int`
- **`height`**: `int`
- **`cells`**: `Array<Array<Cell>>` (2D matrix)

#### Cell
- **`entity_id`**: `string (UUID)|null`
- **`obstacle`**: `boolean`

#### Turn
- **`player_id`**: `string (UUID)`
- **`entity_id`**: `string (UUID)`
- **`delay`**: `int`

#### Position
**Specification:** [[api_go_battle_engine]]
- **`x`**: `int`
- **`y`**: `int`

### 4.3 Entity & Player

#### Entity
**Specification:** [[entity_character]]

Detailed state of a single actor.

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | `string (UUID)` | Unique identifier for the entity. |
| `is_self` | `boolean` | True if the entity belongs to the requesting player. |
| `team` | `int` | Team identifier. |
| `name` | `string` | Display name. |
| `hp` | `int` | Current Hit Points. Dead units are marked with `hp: 0`. |
| `max_hp` | `int` | Maximum Hit Points. |
| `dead` | `boolean` | True if the character has been eliminated in this session. |
| `attack` | `int` | Base offensive power. |
| `defense` | `int` | Base defensive mitigation. |
| `move` | `int` | Remaining movement range for the current turn. |
| `max_move` | `int` | Total movement range attribute. |
| `position` | `Position` | Current `{x, y}` coordinates. |

#### Player
- **`is_self`**: `boolean` (True if this is the requesting user)
- **`nickname`**: `string`
- **`team`**: `int`
- **`ia`**: `boolean` (True if controlled by engine)
- **`entities`**: `Array<Entity>`

### 4.4 UserResource
**Specification:** [[api_laravel_gateway]]

| Field | Type | Description |
| :--- | :--- | :--- |
| `account_name` | `string` | Displayed name. |
| `role` | `string` | User's role (e.g., 'Player', 'Admin'). |
| `ws_channel_key`| `string (UUID)` | Pseudonym for secure WebSocket private channel subscription. |
| `email` | `string` | User's email address. |
| `full_address` | `string` | User's residential address. |
| `birth_date` | `string (ISO8601)` | User's date of birth. |
| `total_wins` | `int` | Total career wins. |
| `total_losses` | `int` | Total career losses. |
| `ratio` | `float` | Win/Loss ratio. |
| `reroll_count` | `int` | Number of times starter stats were rerolled. |
| `characters` | `Array<CharacterResource>` | Optional: Loaded character list. |

### 4.5 CharacterResource
**Specification:** [[api_profile_character]]

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | `string (UUID)` | Character's unique identifier. |
| `name` | `string` | Display name. |
| `hp` | `int` | Character's Hit Points. |
| `attack` | `int` | Character's Attack stat. |
| `defense` | `int` | Character's Defense stat. |
| `movement` | `int` | Current movement range. |
| `initial_movement`| `int` | Base movement range (for cap calculation). |

### 4.6 GameMatchResource
**Specification:** [[api_battle_proxy]]

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | `string (UUID)` | Match unique identifier. |
| `game_mode` | `string` | e.g., "1v1_PVP". |
| `started_at` | `string (ISO8601)` | Start timestamp. |
| `concluded_at` | `string (ISO8601)\|null` | End timestamp. |
| `winner_team_id` | `int|null` | Winning team identifier. |

### 4.7 ArenaEvent
**Specification:** [[api_go_webhook_callback]]

Payload for the asynchronous engine callback.

| Field | Type | Description |
| :--- | :--- | :--- |
| `match_id` | `string (UUID)` | The Laravel Match ID. |
| `event_type` | `string` | e.g., `game.started`, `turn.started`. |
| `player_id` | `string (UUID)` | Optional: Targeted player. |
| `entity_id` | `string (UUID)` | Optional: Targeted entity. |
| `data` | `BoardState` | The current state of the board. |
| `version` | `int64` | Monotonic sequence number (synced with `data.sequence`). |
| `timeout` | `string (ISO8601)` | End of the current turn clock. |

---

## 5. Traceability Matrix

| Endpoint | Specification | Use Case (usecase.md) | Business Requirement (BRD.md) |
| :--- | :--- | :--- | :--- |
| `POST /auth/register` | [[api_auth_register]] | [[uc_player_registration]] | 2.1 User Onboarding & Identity |
| `POST /auth/login` | [[api_auth_login]] | [[uc_player_login]] | 2.1 User Onboarding & Identity |
| `POST /profile/character/{id}/reroll` | [[api_profile_character]] | [[uc_player_registration]] (Step 4) | 2.1 frictionless Entry |
| `POST /profile/character/{id}/upgrade` | [[api_profile_character]] | [[uc_progression_stat_allocation]] | 2.5 Character Progression |
| `POST /matchmaking/join` | [[api_matchmaking]] | [[uc_matchmaking]] | 2.3 Matchmaking Ecosystem |
| `POST /game/{id}/action` | [[api_battle_proxy]] | [[uc_combat_turn]] | 2.4 Combat Engine & Action Economy |
| `POST /api/webhook/upsilon` | [[api_go_webhook_callback]] | [[uc_combat_turn]] / [[uc_match_resolution]] | 2.4 Combat Engine (State Evaluation) |
| `GET /api/profile/export` | [[api_profile_export]] | Data Portability | 3.2 GDPR & Data Privacy |
| `POST /auth/logout` | [[api_auth_logout]] | [[uc_auth_logout]] | [[req_security]] |
| `Universal Envelope` | [[api_standard_envelope]] | All Interactions | 3.3 Traceability & Request ID |
| `GET /leaderboard` | [[api_leaderboard]] | [[us_leaderboard_view]] | 4 Social & Competitive |
| `POST /internal/arena/start` | [[api_go_battle_start]] | Queue to Battle transition | 2.3 PvP/PvE Matchmaking |

---

## 6. Gap Analysis: Uncovered Requirements

Based on a cross-reference with `usecase.md` and `BRD.md`, the following requirements have **no currently defined endpoints** in the API surface:

### 6.1 Administrative Management
- **Auth (UC-7):** Dedicated administrator login and high-privilege JWT exchange via `/admin/login` and role-based redirect.
- **User Controls (UC-8):** Admin can `LIST` players, `SOFT DELETE` accounts, and `ANONYMIZE` (GDPR) data via `/admin/users`.
- **Match History (UC-9):** Placeholder established in Admin Dashboard; full implementation pending log audit service logic.

### 6.2 Advanced Identity & Privacy
- **Missing Anonymization (UC-8 / BRD 3.2):** While Data Portability exists, there is no endpoint for "Right to be Forgotten" (anonymization of Address/Birth Date).
- **Missing Account Management:** No endpoint for updating `Full Address` or `Birth Date` post-registration (though registration is covered).

### 6.3 Social & Competitive
- **Missing Match History (BRD 2.3):** No endpoint for a player to view their own personal history of past matches (separate from current cached board state).
