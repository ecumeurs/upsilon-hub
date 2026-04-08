---
id: requirement_customer_user_id_privacy
status: DRAFT
type: REQUIREMENT
layer: CUSTOMER
priority: 3
tags: security,privacy,auth
version: 1.0
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
- The frontend MUST use a locally generated, persistent pseudonym (Tactical ID) for display purposes if an identifier is required in the UI.
- Identity resolution for API requests MUST be handled purely by backend session or JWT token processing, never by a client-provided user ID.

## TECHNICAL INTERFACE
- **PHP Resource:** `UserResource` (must exclude `id` field)
- **Frontend Utility:** `tactical_id.js`
- **Code Tag:** `@spec-link [[requirement_customer_user_id_privacy]]`

## EXPECTATION
