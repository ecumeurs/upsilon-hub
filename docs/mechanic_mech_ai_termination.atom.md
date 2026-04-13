---
id: mechanic_mech_ai_termination
status: DRAFT
human_name: AI Termination Pattern
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
dependents: []
priority: 4
tags: [concurrency, ai]
parents: []
---

# New Atom

## INTENT
Ensures AI controllers terminate without blocking their main actor loop during match resolution.

## THE RULE / LOGIC
1. Buffered Communications: Channels used for end-of-game signals MUST be buffered (size 1).
2. Non-blocking Sends: Sends to termination channels MUST use selecting with default:
   ```go
   select {
   case ctl.BattleFinished <- true:
   default:
   }
   ```
3. Lifecycle Independence: Termination signals should not block the reception of subsequent ActorStop messages.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_ai_termination]]`
- **Component:** `AggressiveController`
- **Channel:** `BattleFinished (buffered)`

## EXPECTATION
- AI Controller does not hang during BattleEnd.
- Message processing continues until ActorStop is received.
- BattleFinished channel send returns immediately if unconsumed.
