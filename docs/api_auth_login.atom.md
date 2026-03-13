---
id: api_auth_login
human_name: Player Login API
type: API
version: 1.0
status: DRAFT
priority: CORE
tags: [auth, login, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents: []
---

# Player Login API

## INTENT
To authenticate existing users and provide a session token.

## THE RULE / LOGIC
**Endpoint:** `POST /api/v1/auth/login`

### Request (Wrapped in [[api_standard_envelope]])
- `email`: `string`
- `password`: `string`

### Response (Wrapped in [[api_standard_envelope]])
- `user`: `UserObject`
- `token`: `string` - JWT Bearer Token.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /api/v1/auth/login`
- **Code Tag:** `@spec-link [[api_auth_login]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestSuccessfulLogin`, `TestInvalidCredentials`

## EXPECTATION (For Testing)
- Correct credentials -> Return 200 OK with Token.
- Wrong password -> Return 401 Unauthorized.
