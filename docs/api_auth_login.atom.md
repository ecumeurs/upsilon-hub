---
id: api_auth_login
human_name: Player Login API
type: API
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: [auth, login, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents:
  - [[uc_player_login]]
---
# Player Login API

## INTENT
To authenticate a survivor by verifying credentials and issuing a secure access token.

## THE RULE / LOGIC
- **URI:** `/api/v1/auth/login`
- **Verb:** `POST`
- **Intent:** Identity Authentication
- **Fully Detailed Input:**
  - `email`: (string) [Mandatory] The registered email address.
  - `password`: (string) [Mandatory] The survivor's secret credential.
- **Fully Detailed Output:**
  - `user`: (object) Profile data (id, account_name, email).
  - `token`: (string) JWT Bearer token for session authorization.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /api/v1/auth/login`
- **Code Tag:** `@spec-link [[api_auth_login]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestSuccessfulLogin`, `TestInvalidCredentials`

## EXPECTATION (For Testing)
- Correct credentials -> Return 200 OK with Token.
- Wrong password -> Return 401 Unauthorized.
