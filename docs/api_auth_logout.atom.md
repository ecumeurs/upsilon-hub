---
id: api_auth_logout
human_name: "Player Logout API"
type: API
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 3
tags: [auth, logout]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
  - [[uc_auth_logout]]
dependents: []
---

# Player Logout API

## INTENT
To securely terminate the active authentication session by revoking the user's current access token.

## THE RULE / LOGIC
1. **Authentication**: Requester must provide a valid Bearer Token (Sanctum).
2. **Token Revocation**: The system must identify the specific token used and delete it from the database.
3. **Redirection**: Following successful revocation, the client should treat the session as invalidated.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /api/v1/auth/logout`
- **Code Tag:** `@spec-link [[api_auth_logout]]`
- **Security:** Middleware `auth:sanctum` mandatory.

## EXPECTATION (For Testing)
1. Return `success: true` in the standard envelope.
2. Subsequent requests with the same token must return `401 Unauthorized`.
