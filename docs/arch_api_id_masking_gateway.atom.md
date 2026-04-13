---
id: arch_api_id_masking_gateway
human_name: "Architectural API ID Masking Gateway"
type: SERVICE
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 5
tags: [security, api, masking, uuid]
parents:
  - [[requirement_customer_user_id_privacy]]
dependents: []
---

# Architectural API ID Masking Gateway

## INTENT
To provide a secure translation layer between internal database identifiers (UUIDs) and public-facing semantic or masked identifiers, preventing reconnaissance and primary key enumeration.

## THE RULE / LOGIC
- **Internal vs Public Boundary:** All raw database UUIDs (User, Character) MUST be intercepted at the API Gateway (Laravel) before reaching the network.
- **Masking Mechanisms:**
  - **Boolean Flags:** Replace User IDs with `is_self: boolean` in collection views (e.g., Leaderboards, Match History).
  - **Pseudonyms:** Use persistent, non-traceable keys for long-term identification where a boolean is insufficient (e.g., `ws_channel_key`).
  - **Encoded Slugs/HashIDs:** (Future) Convert internal IDs to URL-safe alphanumeric strings for public identification of entities (Characters).
- **Inbound Validation (Ownership):** 
  - For every state-changing request (Actions, Upgrades), the Gateway MUST verify that the authenticated User owns the targeted Entity (Character/Match Participant) before proxying to the Battle Engine.
  - Formula: `authenticated_user_id == target_entity.owner_id` (via Laravel Policies).
- **Match Scoping:** Match IDs are permissible in URLs but MUST be guarded by participant-level authorization.

## TECHNICAL INTERFACE (The Bridge)
- **Laravel Resources:** Use `toArray()` to filter out `id` and inject `is_self`.
- **Middleware/Policies:** `CharacterPolicy` and `MatchParticipantPolicy` for ownership enforcement.
- **Code Tag:** `@spec-link [[arch_api_id_masking_gateway]]`

## EXPECTATION (For Testing)
- `GET /api/v1/leaderboard` -> No `id` field present; `is_self` correctly identifies the caller.
- `GET /api/v1/matchmaking/status` -> `user_id` should be omitted.
- `POST /api/v1/game/{id}/action` with an `entity_id` not owned by the user -> `403 Forbidden`.
