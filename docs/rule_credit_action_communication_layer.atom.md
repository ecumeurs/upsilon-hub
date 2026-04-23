---
id: rule_credit_action_communication_layer
status: STABLE
priority: 5
layer: ARCHITECTURE
version: 2.0
tags: ["credits", "communication", "api"]
parents: []
dependents: []
type: RULE
---

# New Atom

## INTENT
To establish the communication layer protocol for all combat-related actions, ensuring traceability through request IDs, version control for compatibility, effect feedback for result confirmation, and credit tracking association with player IDs.

## THE RULE / LOGIC
**Action Message Structure:**
```go
type ActionMessage struct {
    RequestID   string     // Traceability identifier
    Version     string     // System version for compatibility
    Action      string     // Action type (attack, move, skill_use)
    EntityID    uuid.UUID   // Entity performing action
    TargetID    uuid.UUID   // Target entity (if applicable)
}
```

**Action Response Structure:**
```go
type ActionResponse struct {
    RequestID    string     // Echo for traceability
    Version      string     // Response version
    Success      bool       // Action outcome (true/false)
    Modified     Modification // Changed game state (delta)
    Credits      int         // Credits earned from this action
    PlayerID     uuid.UUID   // Credit recipient (player ID)
    Error        string     // Error message if failed
}
```

**Credit Association Rules:**
- All credit-earning actions must include player ID in response
- Credits are immediately associated with the player who performed the action
- Shield caster receives credits even if they die after shield application
- Skill caster receives credits for status effects they apply

**Version Protocol:**
- Version string in format: "v2.0.0"
- Backwards compatibility checks based on version comparison
- Mismatch versions trigger warning or rejection

**Traceability:**
- Request ID is UUIDv7 for global uniqueness
- All logs must use 8-character ref_id prefix
- Response echoes Request ID for client correlation

**Error Handling:**
- Failures include descriptive error string
- Success responses include modified game state
- Credits are always returned (0 if no credits earned)

### Communication Flow:
```
Client → Server (ActionMessage)
        ↓
Server validates action
        ↓
Server → Client (ActionResponse) { Success: true, Credits: X, Modified: {...} }
        ↓
Server credits player's account
        ↓
Server broadcasts state update (if applicable)
```

## TECHNICAL INTERFACE

## EXPECTATION
