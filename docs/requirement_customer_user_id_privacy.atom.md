---
id: requirement_customer_user_id_privacy
status: STABLE
type: REQUIREMENT
layer: CUSTOMER
priority: 3
tags: security,privacy,auth
version: 1.1
parents: []
dependents:
  - [[entity_player]]
human_name: User ID Privacy Policy
---

# New Atom

## INTENT
To ensure that internal database user IDs are never exposed to the frontend, protecting the system from primary key enumeration and enhancing user privacy.

## THE RULE / LOGIC
- The primary UUID (database ID) of a user MUST NOT be sent to the client (frontend).
- WebSocket private channels MUST be keyed using a secure, persistent pseudonym (`ws_channel_key`) generated on the backend and exposed via the UserResource.
- Identity resolution for API requests MUST be handled purely by backend session or JWT token processing, never by a client-provided user ID.

## TECHNICAL INTERFACE
- **PHP Resource:** `UserResource` (must exclude `id` field)
- **Frontend Utility:** `tactical_id.js`
- **Code Tag:** `@spec-link [[requirement_customer_user_id_privacy]]`

## EXPECTATION
