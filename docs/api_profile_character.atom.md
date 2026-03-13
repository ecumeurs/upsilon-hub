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
- `POST /api/v1/profile/{id}/character/{characterId}/reroll`: Reroll stats (restricted to new accounts). **Updates `initial_movement`.**
- `POST /api/v1/profile/{id}/character/{characterId}/upgrade`: Allocate points during level up. **Validated against [[rule_progression]].**

### Request - Upgrade (Wrapped in [[api_standard_envelope]])
- `stats`: `Object` - Key-value pair of stats to increase (e.g., `{"attack": 1}`).

### Response (Wrapped in [[api_standard_envelope]])
- `character`: `CharacterObject` (Matches [[entity_character]])

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `/api/v1/profile/*`
- **Code Tag:** `@spec-link [[api_profile_character]]`
- **Related Issue:** `ISS-016`
- **Test Names:** `TestGetCharacters`, `TestRerollRestricted`, `TestLevelUpStatAllocation`

## EXPECTATION (For Testing)
- Requesting character list -> Return array of characters.
- Upgrading beyond available points (wins) -> Return 400 Bad Request.
- Upgrading that violates [[rule_progression]] (e.g., movement limit) -> Return 400 Bad Request.
- Rerolling after account is "Stable" -> Return 403 Forbidden.
