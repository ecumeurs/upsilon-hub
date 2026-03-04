---
id: req_security
human_name: JWT Security Requirement
type: REQUIREMENT
version: 1.0
status: REVIEW
priority: CORE
tags: [security, network]
parents: []
dependents:
  - [[ui_dashboard]]
  - [[ui_waiting_room]]
  - [[ui_board]]
---

# JWT Security Requirement

## INTENT
Ensures that all gameplay and private user endpoints are authenticated strictly using JSON Web Tokens (JWT).

## THE RULE / LOGIC
- Public Access: Only the Landing Page and Registration features are fully exempt from authorization.
- Token Exchange: A successful login or registration immediately issues a JWT to the client.
- Authorization: Every other UI page (Dashboard, Waiting Room, Board Page) and all backend API calls require a valid JWT. If missing or expired, the user is redirected to the Landing Page.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[req_security]]`
- **Test Names:** `TestJWTMissingRedirect`, `TestJWTValidAccess`

## EXPECTATION (For Testing)
- Request to Dashboard without JWT -> HTTP 401 or redirect to Landing page.
- Request to Dashboard with valid JWT -> OK.
