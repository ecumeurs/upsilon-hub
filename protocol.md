# Technical Protocol - Upsilon Laravel API Gateway

This document outlines the standard sequence of API requests and responses for the Upsilon Battle system. All requests (except login/register) require a `Bearer {token}` header. All API responses follow the `api_standard_envelope` (request_id, success, message, data).

---

## 1. Authentication & Character Creation

### Step 1: Registration
**Request:** `POST /api/v1/auth/register`
```json
{
  "account_name": "HeroPlayer",
  "email": "hero@example.com",
  "password": "SecurePassword123",
  "password_confirmation": "SecurePassword123"
}
```
**Expected Result:**
- Status `201 Created`.
- `data.user`: User object with a `uuid`.
- `data.token`: Bearer token for future requests.
- **Side Effect:** 3 characters are automatically generated for the player (10 stat points each).

---

## 2. Character Review & Preparation

### Step 2: List Characters
**Request:** `GET /api/v1/profile/{user_id}/characters`
**Expected Result:**
- Array of 3 characters with stats: `hp`, `attack`, `defense`, `movement`.

### Step 3: Reroll (Optional)
**Request:** `POST /api/v1/profile/{user_id}/character/{character_id}/reroll`
**Expected Result:**
- Stats are redistributed. `reroll_count` increments (Limit: 3 total per user account).

---

## 3. Matchmaking

### Step 4: Initiate Matchmaking
**Request:** `POST /api/v1/matchmaking/join`
```json
{
  "game_mode": "1v1_PVP" 
}
```
**Side Effect:** The system automatically pulls the user's first 3 characters for the match.

### Step 5: Poll Match Status
**Request:** `GET /api/v1/matchmaking/status`
**Expected Result:**
- `data.status`: "queued", "matched", or "idle".
- `data.match_id`: `uuid` (if "matched").
- `data.expected_participants`: `int` (total needed for mode).
- `data.empty_slots`: `int` (remaining slots).
- `data.queued_at`: `datetime` (if "queued").

**Available Modes (`game_mode`):**
- `1v1_PVP`: Standard 2-player match.
- `1v1_PVE`: Instant match against AI.
- `2v2_PVP`: 4-player match (Teams: P1+P2 vs P3+P4).
- `2v2_PVE`: 2-player co-op against 2 AI opponents.

**Initial Join Result:**
- Status `200 OK`.
- Same data structure as the status polling endpoint.

---

## 4. Real-time Communication (WebSocket)

### Step 6: Connect to WebSocket
**Server:** Laravel Reverb (default port 8080 or 443).
**Primary Tactical Channel:** `private-user.{ws_channel_key}` (Requires authentication).
**Shared Interaction Channel:** `private-arena.{match_id}` (For common events like chat/emojis).
**Core Event:** `board.updated`.

**Expected Ingestion:**
- Every time a tactical action occurs (via engine webhook), a `board.updated` event is pushed to each participant's private channel.
- These events are **surgically masked**: your view contains full data for your characters and limited data for opponents.

---

## 5. Battle Interaction (Proxy)

### Step 7: Send Battle Action
**Request:** `POST /api/v1/game/{match_id}/action`
```json
{
  "player_id": "user-uuid",
  "entity_id": "character-uuid",
  "type": "Move",
  "target_coords": [2, 3]
}
```
**Expected Result:**
- Status `200 OK`.
- Forwarded result from Go engine.
- A broadcast is triggered to the WebSocket channel.

---

## 6. Game Completion

### Step 8: Final State Update
The game ends when the engine sends a final webhook.
**Webhook Ingestion:** `POST /api/webhook/upsilon`
- The `data` carries the final board state and a `is_finished: true` flag (engine dependent).
- Match result is persisted in `game_matches.game_state_cache`.

---

## 7. Progression & Level Up

### Step 9: Upgrade Character
**Request:** `POST /api/v1/profile/{user_id}/character/{character_id}/upgrade`
```json
{
  "stats": {
    "attack": 1,
    "hp": 2
  }
}
```
**Expected Result:**
- Stats increment in the database.
- Character updated object returned in `data`.
