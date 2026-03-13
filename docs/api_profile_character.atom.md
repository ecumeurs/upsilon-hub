---
id: api_profile_character
human_name: Character Management API
type: API
version: 1.0
status: DRAFT
priority: CORE
tags: [profile, character, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents:
  - [[entity_character]]
---

# Character Management API

## INTENT
To manage player characters, including registration-time rerolls, stat updates, and level-ups.

## THE RULE / LOGIC
**Endpoints:**
- `GET /api/v1/profile/{id}/characters`: List all characters for a user.
- `GET /api/v1/profile/{id}/character/{characterId}`: Get specific character details.
- `POST /api/v1/profile/{id}/character/{characterId}/reroll`: (New Path) Reroll stats (restricted to new accounts).
- `POST /api/v1/profile/{id}/character/{characterId}/upgrade`: (New Path) Allocate points during level up.

### Request - Upgrade (Wrapped in [[api_standard_envelope]])
- `stats`: `Object` - Key-value pair of stats to increase (e.g., `{"attack": 1}`).

### Response (Wrapped in [[api_standard_envelope]])
- `character`: `CharacterObject` (Matches [[entity_character]])

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `/api/v1/profile/*`
- **Code Tag:** `@spec-link [[api_profile_character]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestGetCharacters`, `TestRerollRestricted`, `TestLevelUpStatAllocation`

## EXPECTATION (For Testing)
- Requesting character list -> Return array of characters.
- Upgrading beyond available points -> Return 400 Bad Request.
- Rerolling after account is "Stable" -> Return 403 Forbidden.
