# E2E / Edge Case Failure Investigation

**Date:** 2026-04-28
**Suite run:** 25 failures / 92 tests (67 passed)
**Source of truth:** per-bot logs in `upsiloncli/tests/logs/<scenario>_Bot-NN.log`

This report categorizes each failure by **first surfaced symptom** (assertion message or 4xx/5xx body) and gives a tentative
verdict: `APP` (server bug), `TEST` (faulty scenario), or `INFRA` (test harness / runner). Several failures look like
"test was written against an older contract" — those are flagged `TEST (stale contract)` and should be the easiest wins.

> Note: tests themselves may be faulty. Where the error_key/HTTP status the test asserts no longer matches what the
> backend actually returns, the test is the most likely culprit — the backend behavior is consistent across siblings.

---

## 1. Cross-cutting issues (multiple tests share a root cause)

### 1.1 Shop-item teardown FK violation — `APP` *(or test cleanup bug)*
Affected:
- `e2e_exotic_weapon_dual_path` (Bot-01)
- `e2e_item_grants_skill` (Bot-01)

Symptom:
```
[REPLY 500]
Route admin_shop_item_delete failed: SQLSTATE[23503]: Foreign key violation:
 update or delete on table "shop_items" violates foreign key constraint
 "player_inventory_shop_item_id_foreign" on table "player_inventory"
```

Cause: both scenarios admin-create a one-off shop item, have the test player purchase it, then try to admin-delete the
item before the player is anonymized. The FK on `player_inventory.shop_item_id` has no cascade, so deletion 500s and the
script aborts.

Fix options:
- **APP:** soft-delete shop items, OR cascade-delete inventory rows, OR refuse with 409 + message.
- **TEST:** delete the player's inventory entry (or anonymize the bot account) before admin-deleting the shop item.

### 1.2 Friendly-fire / pathing test scaffolding — `TEST` (or `INFRA`)
Affected:
- `e2e_friendly_fire_prevention_with_4` (Bot-02/03/04)
- `edge_attack_friendly_fire_with_4` (Bot-01/02/04)
- `edge_attack_targeting_rules_with_4` (Bot-01/03/04)
- `edge_movement_entity_collision_with_2` (Bot-01)

Symptom (representative):
```
JS Exception: Assertion Failed: Never reached an ally to attempt friendly-fire within 60 rounds
JS Exception: Assertion Failed: Never reached enemy to test entity collision within 60 rounds
```

Cause: bot AI walks toward an ally/enemy looking for a positioning condition that doesn't fire within the 60-round
budget. Either the map seed produces unreachable positions, or the bot logic is not converging. Bot-01 in
`e2e_friendly_fire_prevention_with_4` *did* manage a 400 from `Friendly fire is not allowed` — the rule itself works,
but only one of four bots reaches the assertion point.

Fix options: lengthen the budget, or inject the targets via test setup rather than walking to them.

### 1.3 Stale `error_key` expectations — `TEST (stale contract)`
Affected (all `entity.controller.missmatch` family):
- `edge_attack_wrong_controller_with_2` (Bot-02)
- `edge_movement_wrong_controller_with_2` (Bot-02)

Symptom:
```
Assertion Failed: Expected entity.controller.missmatch (Expected: entity.controller.missmatch, Actual: <nil>)
```
Server returns **HTTP 403** with body `Forbidden: You do not own the entity specified in this action.` — i.e. the rule
is enforced, but the response carries no machine-readable `error_key`, so the test sees `<nil>`.

Fix options:
- **APP:** populate `error_key=entity.controller.missmatch` on this 403 to match the documented envelope.
- **TEST:** assert on HTTP 403 + message substring instead of `error_key`.

Other stale-contract instances:
| Test | Expected | Actual |
|---|---|---|
| `edge_movement_already_attacked` | `entity.movement.already` | `entity.path.notadjacent` |
| `edge_movement_out_of_turn_with_2` | `entity.turn.missmatch` | `arena.notfound` |
| `edge_movement_obstacle_collision` | path-family `error_key` | `undefined` |
| `edge_attack_targeting_rules_with_4` Bot-02 | `rule.friendly_fire` | `entity.attack.friendlyfire` |

### 1.4 "Acceptable status / message" assertions firing on the *expected* status — `TEST (assertion logic)`
Affected (each got the "right" 4xx but the test still failed):
- `edge_shop_unknown_item` — got 404, asserts "Error must be 404 Not Found" → assertion text is misleading; the assert is checking *something else* (probably `error_key` or message body) but fails opaquely.
- `edge_skill_template_not_found` — got 404, "Must be 404 Not Found" failed.
- `edge_skill_equip_invalid_id` — got 404, "Must be 404, 422, or 403" failed.
- `edge_skill_unequip_not_equipped` — got 422, "Must be 422 Unprocessable or 400" failed.
- `edge_skill_unowned_character_equip` — got 403, "Must be 403 Forbidden or 404" failed.
- `edge_skill_unowned_character_roll` — got 403, "Must be 403 Forbidden or 404" failed.
- `edge_unequip_empty_slot` — got 404, "Error must be 404 Not Found" failed.

Pattern: status code matches what the test claims to want, yet the assertion fails. Almost certainly the assertion
helper checks both status **and** an additional field (likely `error_key` or message) and the helper's failure message
only mentions the status. Verify the helper, or split into two assertions so the failure log identifies which limb broke.

### 1.5 `edge_equip_unowned_*` use 8-char password against a 15-char policy — `TEST (clear bug)`
Affected:
- `edge_equip_unowned_character` (Bot-01)
- `edge_equip_unowned_item` (Bot-01)

Symptom:
```
[REPLY 422] auth_register failed -- "The password field must be at least 15 characters."
JS Exception: [object Object]
```

Both scenarios hardcode `passA = passB = "Pass123!"` (8 chars). Password policy requires ≥15. Tests need to be updated
to use a compliant password (e.g. `"Password12345!"` or pull from a shared test constant).

---

## 2. Single-test failures

| Test | Verdict | Symptom / Root cause |
|---|---|---|
| `edge_auth_password_policy_full` | `TEST (stale contract)` | Asserts message contains `'least 15'`, server returned `'-- DEBUG MODE -- Validation failed'`. Real per-field detail is in `meta.errors.password[0]`. Test should read `meta.errors`, not `message`. |
| `edge_char_reroll_limit` | `TEST` | Got expected 403 + `Reroll limit reached.` but threw `JS Exception: [object Object]`. Test treats the 403 as failure — should accept it as the expected "limit hit" branch. |
| `edge_match_action_after_end_with_2` | `TEST` or `APP (race)` | Bot-01 fails with `Match ended unexpectedly`; Bot-02 then 400s on forfeit (`Game is not in progress`). Match completed before the post-end action could be issued — test assumed it could act before win. |
| `edge_match_queue_while_in_match_with_2` | `TEST` or `APP` | After 409 on re-queue (correctly enforced), test asserts "Match state should still be accessible" and that fails. The match-state read is returning something it does not expect. Worth checking whether match read returns `null` once you 409 once. |
| `edge_prog_allocation_no_wins` | `APP` | `Assertion Failed: HP changed after failed upgrade (Expected: 30, Actual: 31)`. Allocation without wins must **not** mutate state but HP went from 30 → 31. Real backend bug: failed allocation is being partially applied. |

---

## 3. Verdict summary

| Bucket | Count | Notes |
|---|---|---|
| **APP bugs** | 2–3 | shop-item delete FK (1.1), `prog_allocation_no_wins` HP leak, possibly some 1.3 entries if we treat missing `error_key` as a contract violation. |
| **Test bugs (clear)** | 9 | 1.4 (7 tests), 1.5 (2 tests). |
| **Test bugs (stale contract)** | 5 | 1.3 entries — backend renamed/removed `error_key`s without updating tests. |
| **Test scaffolding (bot AI / pathing)** | 7 | 1.2 — friendly-fire and collision walks can't converge in 60 rounds. |
| **Test bugs (other)** | 3 | `edge_auth_password_policy_full`, `edge_char_reroll_limit`, possibly `edge_match_*`. |

Recommended order of attack:
1. Quick wins: 1.5 (passwords) and 1.4 (assertion helper) — pure test edits.
2. 1.3 + `edge_auth_password_policy_full` — decide whether to fix backend contract (add `error_key`) or update tests to assert on status only. `error_key`s are part of the documented envelope, so fixing backend is the right call.
3. 1.1 — backend fix (cascading delete on shop_item) or cleanup the bot inventory before deletion.
4. `edge_prog_allocation_no_wins` — real bug, HP must not change on failed upgrade.
5. 1.2 — increase round budget or rewrite scenarios to position bots deterministically.
