---
id: api_profile_character
human_name: Character Management API
type: API
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 5
tags: [profile, character, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents: []
---
# Character Management API

## INTENT
To manage player characters, including registration-time rerolls, stat updates, and level-ups.

## THE RULE / LOGIC
**Endpoints:**
- `GET /api/v1/profile/characters`: List all characters for the authenticated user.
- `GET /api/v1/profile/character/{characterId}`: Get specific character details.
- `POST /api/v1/profile/character/{characterId}/reroll`: Reroll stats (restricted to new accounts). 
- `POST /api/v1/profile/character/{characterId}/upgrade`: Allocate points. Validated against [[rule_progression]].

### CharacterResource (Common Response)
- `id`: `string (UUID)`
- `name`: `string`
- `hp`: `int`
- `attack`: `int`
- `defense`: `int`
- `movement`: `int`
- `initial_movement`: `int`

### Request - Upgrade (Wrapped in [[api_standard_envelope]])
- `stats`: `object` - Increments for stats (e.g., `{"attack": 1, "hp": 2}`).

### Response - Reroll (Wrapped in [[api_standard_envelope]])
- `character`: `CharacterResource`
- `reroll_count`: `int`

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
