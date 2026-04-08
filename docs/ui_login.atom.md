---
id: ui_login
human_name: "Player Login UI"
type: UI
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 3
tags: [ui, auth]
parents:
  - [[uc_player_login]]
dependents: []
---

# Player Login UI

## INTENT
To provide a secure, industrial-themed interface for players to authenticate themselves using their unique survivor credentials.

## THE RULE / LOGIC
1.  **Identification**: Accepts `account_name` as the primary identifier.
2.  **Authentication**: Accepts `password` for credential verification.
3.  **Aesthetics**: Must adhere to the "Neon & Rust" industrial design system.
4.  **Feedback**: Must display "Identification Failure" or similar for invalid credentials.
5.  **Persistence**: Successfully authenticated sessions must persist the JWT Bearer token in the browser's local storage.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag**: `@spec-link [[ui_login]]`
- **Route**: `/login`
- **API Call**: `POST /api/v1/auth/login`

## EXPECTATION (For Testing)
- Entering valid credentials redirects the user to the Tactical Dashboard.
- Entering invalid credentials displays a clear error state and prevents redirection.
