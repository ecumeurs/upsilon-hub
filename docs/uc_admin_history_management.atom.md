---
id: uc_admin_history_management
human_name: Administrative Match History Management Use Case
type: MODULE
version: 1.0
status: DRAFT
priority: CORE
tags: [admin, match-history]
parents: [entity_game_match]
dependents: []
---

# Administrative Match History Management Use Case

## INTENT
Provides administrators with the tools to review match results and purge old history to maintain system performance.

## THE RULE / LOGIC
- **Review Results:** Admin can list all completed matches and their outcomes.
- **History Purge:** 
  - Admin can perform a "clean" operation to remove match records older than a specific threshold (e.g., 90 days).
  - This operation is destructive and must be confirmed.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[uc_admin_history_management]]`
- **Test Names:** `TestAdminListMatches`, `TestAdminPurgeHistory`

## EXPECTATION (For Testing)
- Admin requests match list -> results shown.
- Admin triggers history clean -> matches older than threshold are deleted from DB.
