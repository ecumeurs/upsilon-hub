# Issue: Turner Hands Next Turn to Recently Deceased Entity

**ID:** `ISS-046_20260416_turner_dead_entity_next_turn`
**Ref:** `ISS-046`
**Date:** 2026-04-16
**Severity:** High
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/ruler/turner`
**Affects:** `upsilonbattle/battlearena/ruler/ruler.go`, `upsilonbattle/battlearena/ruler/rules/endofturn.go`, `upsilonbattle/battlearena/ruler/rules/attack.go`

---

## Summary

When a character kills another character and then ends their turn, the `Turner` can hand the next turn to the recently-killed entity. Because the entity has already been removed from `gs.Entities` and its controller has no live entity to act with, the game stalls: `ControllerNextTurn` is never sent, or is sent to an entity whose controller cannot respond, permanently hanging the battle loop.

---

## Technical Description

### Background

`Turner` maintains a sorted `[]EntityTurn` queue. `NextTurn()` pops the head of the queue, stores it in `CurrentEntityTurn`, and returns the entity's UUID to the caller. The caller (`ruler.endOfTurn`) then looks up the entity in `gs.Entities` and sends a `ControllerNextTurn` notification to the owning controller.

When `attack.go` kills a foe (HP ≤ 0), it calls:
1. `gs.Grid.RemoveEntity(foe.Position)` — removes from grid.
2. `delete(gs.Entities, foe.ID)` — removes from entity map.
3. `gs.Turner.RemoveEntity(foe.ID)` — removes from **future** queue slots.

### The Problem Scenario

The critical issue is that `Turner.RemoveEntity` only scans and removes the entity from the `t.Turns` slice (future queue). It does **not** clear `t.CurrentEntityTurn`. Additionally, the dead entity may have been re-added to the Turner queue by `endofturn.go` in a specific ordering of events.

The most direct and reproducible path is:

```
Step 1:  Entity A (acting) attacks Entity B.
         → B's HP drops to 0.
         → attack.go:82  gs.Grid.RemoveEntity(B.Position)
         → attack.go:83  delete(gs.Entities, B.ID)
         → attack.go:84  gs.Turner.RemoveEntity(B.ID)
            ↳ Removes B from t.Turns future slots only.
            ↳ t.CurrentEntityTurn is unchanged (still = whoever
              called NextTurn() last, which may or may not be B).

Step 2:  Entity A's controller calls EndOfTurn (to pass the turn).

Step 3:  endofturn.go:87
         gs.Turner.AddEntity(req.EntityID=A, A.CurrentDelay)
         → A is re-queued into future slots. ✓

Step 4:  endofturn.go:88  gs.IncTurn()

Step 5:  ruler.go:470
         nextTurnEnt := r.GameState.Turner.NextTurn()
         → Pops from t.Turns[0].
         → IF B was the very next entity in the sorted queue at the
           time of the kill (e.g., lowest delay after A), AND the
           kill happened right after AddEntity re-queued B for its
           next turn (i.e., a previous EndOfTurn passed B's entity
           through AddEntity before B was killed this turn),
           THEN B may be at index 0 of t.Turns when NextTurn() runs.
         → NextTurn() returns B.ID and sets CurrentEntityTurn = B.ID

Step 6:  ruler.go:474
         beg, found := r.GameState.Entities[nextTurnEnt]
         → found == false (B was deleted in Step 1)
         → BeginingOfTurn is silently skipped. No error, no retry.

Step 7:  ruler.go:504-521
         The code checks `if nextTurnEnt != uuid.Nil` — B.ID is not
         Nil, so it proceeds.
         ent := r.GameState.Entities[nextTurnEnt]
         → Returns zero-value entity.Entity (Go map miss).
         ctrl, found := r.GameState.Controllers[ent.ControllerID]
         → ent.ControllerID is uuid.Nil (zero value).
         → found == false → error logged, but NO ControllerNextTurn
           notification is sent to any controller.

RESULT:  No controller receives ControllerNextTurn.
         The shot clock was started for a dead entity ID.
         When it fires, targetEnt = r.GameState.Entities[CurrentEntityTurn]
         is again a zero-value entity; the timeout EndOfTurn is filed
         against ControllerID uuid.Nil, which fails the
         CheckControllerForEntity guard → the game is permanently hung.
```

### Where This Pattern Exists Today

| File | Line(s) | Description |
|---|---|---|
| `upsilonbattle/battlearena/ruler/turner/turner.go` | 83–90 | `RemoveEntity` only purges future slots; does not clear `CurrentEntityTurn` |
| `upsilonbattle/battlearena/ruler/rules/attack.go` | 80–84 | Removes dead entity from Turner future queue after kill |
| `upsilonbattle/battlearena/ruler/rules/endofturn.go` | 87 | `AddEntity` re-queues the acting entity unconditionally before NextTurn is called |
| `upsilonbattle/battlearena/ruler/ruler.go` | 470–521 | `endOfTurn` handler calls `NextTurn()` then does a map lookup; silently misses if entity is dead, continues with zero-value entity |
| `upsilonbattle/battlearena/ruler/ruler.go` | 610–646 | `startShotClock` captures `CurrentEntityTurn` at timer fire — if it points to a dead entity, the timeout EndOfTurn will fail all guards |

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High — triggered any time a kill happens on the turn immediately before the deceased entity's scheduled next turn |
| Impact if triggered | High — the entire battle permanently hangs; no controller receives the next-turn signal, no recovery path exists |
| Detectability | Low — the only observable symptom is silence; no panic, no error returned to callers; the log shows "Controller not found" at ruler.go:513 but does not self-heal |
| Current mitigant | None — the `if !found` branch at ruler.go:510 logs an error but does not retry or skip to the next living entity |

---

## Recommended Fix

**Short term:** In `ruler.go:endOfTurn`, after calling `Turner.NextTurn()`, guard for the case where `nextTurnEnt` is not in `gs.Entities`. If the entity is absent (i.e., it was killed), immediately call `Turner.NextTurn()` again in a loop until a live entity is found or the queue is empty. This prevents handing a turn to a dead entity without changing the Turner or GameState architecture.

```go
// In ruler.go, endOfTurn handler, after line 470:
for nextTurnEnt != uuid.Nil {
    if _, alive := r.GameState.Entities[nextTurnEnt]; alive {
        break
    }
    // Entity was killed; skip to the next one in queue.
    r.RequestLogger.WithFields(logrus.Fields{
        "skippedEntityID": nextTurnEnt.String()[0:8],
    }).Warn("Next-turn entity was dead, skipping to next in queue")
    nextTurnEnt = r.GameState.Turner.NextTurn()
}
```

**Medium term:** Make `Turner.RemoveEntity` also clear `CurrentEntityTurn` to `uuid.Nil` if the removed entity is the current one. This prevents the shot clock from ever referencing a dead entity ID. Additionally, add a guard in `startShotClock` to validate that `CurrentEntityTurn` is still present in `gs.Entities` before sending the timeout EndOfTurn.

**Long term:** Introduce a `Turner.SkipDeadEntities(livingSet map[uuid.UUID]bool)` method that the ruler calls after every kill event (not just at EndOfTurn). This would let the Turner stay self-consistent rather than relying on callers to handle missing entities defensively.

---

## References

- [`upsilonbattle/battlearena/ruler/turner/turner.go`](../upsilonbattle/battlearena/ruler/turner/turner.go) — Turner queue implementation
- [`upsilonbattle/battlearena/ruler/rules/attack.go`](../upsilonbattle/battlearena/ruler/rules/attack.go) — Kill path, lines 80–84
- [`upsilonbattle/battlearena/ruler/rules/endofturn.go`](../upsilonbattle/battlearena/ruler/rules/endofturn.go) — `AddEntity` re-queue, line 87
- [`upsilonbattle/battlearena/ruler/ruler.go`](../upsilonbattle/battlearena/ruler/ruler.go) — `endOfTurn` handler, lines 470–521; `startShotClock`, lines 610–646
