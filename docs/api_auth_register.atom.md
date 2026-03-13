---
id: api_auth_register
human_name: Player Registration API
type: API
version: 1.0
status: DRAFT
priority: CORE
tags: [auth, register, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents:
  - [[uc_player_registration]]
---

# Player Registration API

## INTENT
To allow new users to create an account and receive an authentication token.

## THE RULE / LOGIC
**Endpoint:** `POST /api/v1/auth/register`

### Request (Wrapped in [[api_standard_envelope]])
- `account_name`: `string` - The user's displayed name.
- `email`: `string` - Must be unique and valid.
- `password`: `string` - Minimum 8 characters.
- `password_confirmation`: `string` - Must match `password`.

### Response (Wrapped in [[api_standard_envelope]])
- `user`: `UserObject`
  - `id`: `uuid`
  - `account_name`: `string`
  - `email`: `string`
- `token`: `string` - JWT Bearer Token.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /api/v1/auth/register`
- **Code Tag:** `@spec-link [[api_auth_register]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestSuccessfulRegistration`, `TestRegistrationValidationFails`

## EXPECTATION (For Testing)
- Valid data -> User created in DB -> Return 201 Created with Token.
- Duplicate email -> Return 422 Unprocessable Entity.
