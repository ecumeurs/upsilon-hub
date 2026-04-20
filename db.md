# TRPG Database Schema

*Generated via ATD Synthesis (`data_persistence`)*
*Last Updated: 2026-04-17 based on implementation analysis*

This document outlines the relational boundaries required in the PostgreSQL implementation to support the TRPG Specifications. 

## Tables Summary

### 1. `players`
Stores authentication identity and tracks top-level metrics for the generic Leaderboard (`ui_leaderboard`).
* `id` (UUID, Primary Key)
* `account_name` (Varchar, Unique, Not Null)
* `email` (Varchar, Unique, Not Null)
* `password_hash` (Varchar, Not Null)
* `full_address` (Varchar, Not Null) - *Private: GDPR protected*
* `birth_date` (Date, Not Null) - *Private: GDPR protected*
* `role` (Enum: 'Player', 'Admin', Default 'Player')
* `total_wins` (Int, Default 0)
* `total_losses` (Int, Default 0)
* `reroll_count` (Int, Default 0)
* `ratio` (Decimal calculated/derived)
* `ws_channel_key` (UUID, Unique) - *For secure WebSocket subscriptions*
* `created_at` (Timestamp)
* `updated_at` (Timestamp)
* `deleted_at` (Timestamp, Nullable) - *Soft delete for GDPR compliance*

**Indexes**: `account_name`, `email`, `role`, `total_wins` (for leaderboard queries)

### 2. `characters`
Stores the individual entities generated via the `entity_character` limits. Linked exclusively to the Player.
* `id` (UUID, Primary Key)
* `player_id` (UUID, Foreign Key -> `players.id`)
* `name` (Varchar)
* `hp` (Int, Min 3)
* `movement` (Int, Min 1)
* `attack` (Int, Min 1)
* `defense` (Int, Min 1)
* `initial_movement` (Int) - *For progression cap calculations*
* `created_at` (Timestamp)
* `updated_at` (Timestamp)

**Constraints**: Each player limited to exactly 3 characters (enforced at application level)
**Indexes**: `player_id`, `name`

### 3. `game_matches`
Stores active and historical match data, including cached board state from the Go engine.
* `id` (UUID, Primary Key)
* `game_mode` (Enum: '1v1_PVE', '1v1_PVP', '2v2_PVE', '2v2_PVP')
* `started_at` (Timestamp)
* `concluded_at` (Timestamp, Nullable)
* `winner_team_id` (Int, Nullable) - *0 or null for draw*
* `board_state` (JSONB) - *Cached tactical state from Go engine*
* `version` (BigInt) - *Monotonic version for state deduplication*
* `turn` (Int) - *Legacy turn counter for compatibility*
* `created_at` (Timestamp)
* `updated_at` (Timestamp)

**Indexes**: `game_mode`, `started_at`, `winner_team_id`, `version`
**JSONB Structure**: Contains `players`, `grid`, `turn`, `current_entity_id`, `timeout` fields per [[api_go_battle_engine]]

### 4. `match_participants`
Mapping table defining which Players competed in a specific historical or active match.
* `id` (BigInt, Primary Key, Auto-increment)
* `match_id` (UUID, Foreign Key -> `game_matches.id`)
* `player_id` (UUID, Foreign Key -> `players.id`)
* `team` (Int) - *Team 1 or Team 2*
* `status` (Enum: 'WIN', 'LOSS', 'DRAW', 'IN_PROGRESS')
* `created_at` (Timestamp)

**Indexes**: `match_id`, `player_id`, `status`
**Unique Constraint**: `(match_id, player_id)` - Prevent duplicate participation

### 5. `matchmaking_queue`
Active queue entries for players seeking matches.
* `id` (BigInt, Primary Key, Auto-increment)
* `player_id` (UUID, Foreign Key -> `players.id`)
* `game_mode` (Enum: '1v1_PVE', '1v1_PVP', '2v2_PVE', '2v2_PVP')
* `status` (Enum: 'QUEUED', 'MATCHED', 'CANCELLED', 'EXPIRED')
* `match_id` (UUID, Nullable, Foreign Key -> `game_matches.id`)
* `queued_at` (Timestamp)
* `matched_at` (Timestamp, Nullable)
* `expires_at` (Timestamp) - *Queue entry expiration*

**Indexes**: `player_id`, `game_mode`, `status`, `expires_at`

---

## Entity Relationship Diagram

```mermaid
erDiagram
    PLAYERS ||--o{ CHARACTERS : "owns (exactly 3)"
    PLAYERS ||--o{ MATCH_PARTICIPANTS : "competes as"
    PLAYERS ||--o{ MATCHMAKING_QUEUE : "queues as"
    GAME_MATCHES ||--|{ MATCH_PARTICIPANTS : "includes"
    GAME_MATCHES ||--o{ MATCHMAKING_QUEUE : "fulfills"
    
    PLAYERS {
        uuid id PK
        varchar account_name UK
        varchar email UK
        varchar password_hash
        varchar full_address "Private: GDPR"
        date birth_date "Private: GDPR"
        varchar role "Player/Admin"
        int total_wins
        int total_losses
        int reroll_count
        decimal ratio
        uuid ws_channel_key UK
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "Soft delete"
    }
    
    CHARACTERS {
        uuid id PK
        uuid player_id FK
        varchar name
        int hp "Min 3"
        int movement "Min 1"
        int attack "Min 1"
        int defense "Min 1"
        int initial_movement
        timestamp created_at
        timestamp updated_at
    }
    
    GAME_MATCHES {
        uuid id PK
        enum game_mode "1v1/2v2 PVE/PVP"
        timestamp started_at
        timestamp concluded_at
        int winner_team_id
        jsonb board_state "Cached Go state"
        bigint version "Monotonic"
        int turn "Legacy counter"
        timestamp created_at
        timestamp updated_at
    }

    MATCH_PARTICIPANTS {
        bigint id PK
        uuid match_id FK
        uuid player_id FK
        int team "Team 1/2"
        enum status "WIN/LOSS/DRAW/IN_PROGRESS"
        timestamp created_at
    }
    
    MATCHMAKING_QUEUE {
        bigint id PK
        uuid player_id FK
        enum game_mode "1v1/2v2 PVE/PVP"
        enum status "QUEUED/MATCHED/CANCELLED"
        uuid match_id FK
        timestamp queued_at
        timestamp matched_at
        timestamp expires_at
    }
```
