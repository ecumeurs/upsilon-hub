---
date: 2026-04-27
context: Local CI run after grid 3D upgrade + skill/item/economy intro (e435806)
total_failed: 58 / 92
scope_excluded: e2e_credit_economy (already under investigation by another contributor)
---

# CI Failure Investigation — Triage Summary

## High-level

Of 58 failed scenarios, **~30 are downstream of a single matchmaking-server hang**. The rest split into 6 smaller, independent clusters — mostly seed data, validation contracts, and HTTP status-code regressions.

Recommended order of attack (highest unblock value first):

1. **Cluster A — Matchmaking hang** (unblocks ~30 tests; likely `upsilonapi` matchmaker or `upsilonbattle` engine integration)
2. **Cluster B — Missing seed data** (unblocks shop/skill template flows; trivial fix)
3. **Cluster C — Admin CRUD validators** (unblocks the dependency chain into items/skills E2E)
4. **Cluster D/G — Status-code regressions on 404/403 paths** (touches a few error-handler atoms)
5. **Cluster E/H/I/J/K** — small targeted test/route fixes

---

## Cluster A — Matchmaking server hang **[BLOCKER]**

**Symptom (two flavors, same root cause suspected):**

- `POST /api/v1/matchmaking/join` → CLI times out after **30 s** with `context deadline exceeded` (no HTTP response from Laravel at all).
- Sometimes the join *does* return `200 {status: "queued"}`, then the bot waits up to **60 s** for the `match.found` WS event which never fires.

**Affected logs (sample):**
- `e2e_matchmaking_pve_instant_Bot-01.log` (1v1_PVE: hung on join)
- `e2e_matchmaking_pvp_queue_Bot-01.log` (1v1_PVP: queued, no `match.found`)
- `e2e_friendly_fire_prevention_Bot-01.log` (2v2_PVP: queued, no `match.found`)
- `e2e_match_resolution_standard`, `e2e_match_resolution_forfeit`, `e2e_combat_turn_management`, `e2e_progression_constraints`, `e2e_progression_post_win`, all `edge_attack_*`, all `edge_movement_*`, `edge_match_*`, `edge_char_reroll_post_match`, `edge_prog_attribute_cap`, `edge_ws_ping_timeout`, `edge_ws_wrong_channel`.
- **30 logs match `matchmaking/join.*context deadline` or `match.found.*timed out`.**

**Notably PVE used to be instant** (the test is literally named `pve_instant`) yet now also hangs — so this is not just a "no second player" issue. PVE bootstrap is broken.

**Likely root causes to investigate (in order):**
1. `upsilonapi` matchmaking controller — was it changed during the grid 3D / skill / item / economy upgrades? Look for blocking calls to the upsilon-engine, missing context cancellation, or a transaction that no longer commits.
2. `upsilonbattle` (Go engine on :8081) — instant-PVE used to spin up a bot match server-side; check if the bot-spawn pipeline still works after grid changes.
3. Reverb (port 8080) broadcast wiring — the `match.found` private-channel event for `private-user.{ws_channel_key}` may no longer be emitted.

---

## Cluster B — Missing seed data **[HIGH, easy fix]**

**Symptom:**
- `GET /api/v1/shop/items` → `data: []` (Asserted: "Shop catalog must not be empty")
- `GET /api/v1/skills/templates` → `data: []` (Asserted: "Must have at least 3 seeded skill templates, got 0")

**Affected:** `e2e_shop_browse_purchase`, `e2e_skill_template_browse`, `e2e_inventory_equip_battle` (also asserts armor + weapon present), and any test that depends on a non-empty catalog.

**Root cause:** `trigger_all_ci_tests.sh` only bootstraps the admin password (line 37). It does not seed shop items or skill templates. The previous test set didn't need them; the new `shop`/`skill` features do.

**Fix path:** Either
- (a) add a `php battleui/artisan db:seed --class=ShopItemsSeeder --class=SkillTemplatesSeeder` call to `trigger_all_ci_tests.sh` after the admin-bootstrap step, **or**
- (b) write a `util_seed_catalog.js` scenario equivalent to `util_purge_all.js` that uses the admin API to insert N items + N templates, and call it before the suite.
- Verify the seeders exist in `battleui/database/seeders/`. If not, they need to be created first.

---

## Cluster C — Admin CRUD validation contract drift **[MEDIUM, 4 tests]**

**Symptom — `POST /api/v1/admin/shop-items`:**
```
422 — properties field is required
```
Test sends `properties: {}`. Validator rejects empty object.

**Symptom — `POST /api/v1/admin/skill-templates`:**
```
422 — targeting / costs / effect fields required
```
Test sends `targeting: {}, costs: {}, effect: {}`. Same shape mismatch.

**Affected:** `e2e_admin_shop_item_crud`, `e2e_admin_skill_crud`, `e2e_item_grants_skill`, `e2e_exotic_weapon_dual_path` (last two depend on admin creating a template/item first).

**Decision needed (ATD)** — which is canonical?
- **Option 1 (fix scripts):** populate the four fields with realistic minimal payloads. Lowest blast radius. Need to look at the request validator FormRequest classes in `upsilonapi/app/Http/Requests/Admin/` to know minimal valid shapes.
- **Option 2 (relax validators):** allow empty objects on these schema-flexible fields. May be wrong if the new contract requires content for the engine to consume the entity.

Recommend **Option 1** unless the API atom (`api_admin_shop_item_create`, `api_admin_skill_template_create`) explicitly says `properties` is optional.

---

## Cluster D + G — Wrong HTTP status codes on not-found / unauthorized **[MEDIUM, ~6 tests]**

**Symptom:** Routes leak Laravel's `ModelNotFoundException` (which renders 500 in debug) instead of returning 404; some return 403 "unauthorized" where 404 is more accurate.

**Examples:**
- `GET /api/v1/admin/skill-templates/{uuid}` for non-existent UUID → message `No query results for model [App\Models\SkillTemplate] 0000…` (status: not 404). Test asserts "Must be 404 Not Found".
- `POST /api/v1/skills/{character_skill_id}/equip` for invalid CharacterSkill UUID → same `No query results for model` pattern.
- `POST /api/v1/skills/{character_skill_id}/equip` on unowned character → 403 "This action is unauthorized" — test expects 403 *or* 404; the assertion failure suggests the actual status code wasn't either (likely 500 from a policy that throws, or 422). Need to inspect actual REPLY status in `edge_skill_unowned_character_equip_Bot-01.log`.
- `POST /api/v1/shop/purchase` for unknown item id → message "Shop item not found.", asserted 404. Confirm status code returned.
- `POST /api/v1/profile/character//unequip/weapon` (note **double slash** — empty character ID) → "route could not be found" 404. Test expected 404 *of the resource*, not of the route. Likely a test bug: the test passes empty ID. **Subcluster H below.**

**Fix path:** Add explicit `findOrFail` → `abort(404)` or use route-model binding with a global exception handler that maps `ModelNotFoundException` to a `404` envelope (`api_standard_envelope`). One handler change covers most of these.

---

## Cluster E — `e2e_admin_user_management` last-admin lockout **[LOW]**

**Symptom:** Test calls `POST /api/v1/admin/users/{id}/anonymize` against the only admin → 400 "Cannot anonymize the last remaining administrator."

**Fix path:** Test setup needs to register a second admin (or seed one) before calling anonymize, OR target a non-admin user.

---

## Cluster F — Skill equip battle: turn events missing **[depends on A]**

**Symptom:** `e2e_skill_equip_battle` reaches a 1v1_PVE match (so the join *did* respond), then `Turn wait timed out or failed: timeout waiting for events: [board.updated game.ended]`. The follow-up `forfeit` returns "arena not found" — i.e. the engine has already discarded the arena.

**Fix path:** Likely fixed by Cluster A. If not, the engine arena lifecycle is dropping early after the grid 3D conversion. Check `upsilonbattle` logs (Engine on :8081) when this scenario runs.

---

## Cluster H — `edge_unequip_empty_slot` test bug **[LOW]**

**Symptom:** URL built is `…/character//unequip/weapon` (empty character id segment). Test passes empty string. API correctly 404s the route, but the test expected a 404 of the resource and asserts a different message.

**Fix path:** Test should resolve a real character ID first, attempt unequip on an empty slot. Or, alternatively, route handler should validate ID format and return a structured 422.

---

## Cluster I — `edge_auth_password_policy_full` regression **[LOW]**

**Symptom:** All 8 registration attempts (one per password variant) fail with 422 "Validation failed", including ones the test labels as "Valid password". Final assertion: "ERROR: Valid password registration failed".

**Possible causes:**
- Password policy was tightened (e.g. now also requires age, address, birth_date format) and the test fixture has not been updated.
- Or a non-password field in the registration payload now fails validation (e.g. `birth_date` format change).

**Fix path:** Pull the actual `meta.errors` payload from the failing log to see *which* field rejected. Then update either the test fixture or the policy atom.

---

## Cluster J — `edge_match_action_after_end` — HP +1 after blocked upgrade **[LOW]**

**Symptom:** "HP changed after failed upgrade (Expected: 30, Actual: 31)". Test attempts a progression upgrade after the match has ended, which should be rejected, but the HP changes anyway.

**Fix path:** Check the post-match progression endpoint guard (`rule_progression`). The "match-ended" gate may not be firing, or the guard is firing but the persistence call still goes through (transaction boundary issue).

---

## Cluster K — `edge_char_reroll_limit` — CLI route map missing `character_get` **[LOW]**

**Symptom:** `[INTERNAL_ERROR] unknown route: character_get`. CLI doesn't know that route name; the test was written assuming it exists.

**Fix path:** Either add `character_get` to `upsiloncli/internal/script/routes.go` (or wherever the route registry lives) pointing at `GET /api/v1/profile/character/{id}`, or rewrite the test to use whatever route name *is* registered.

---

## Suggested sonnet task split

> Each task is sized for one Sonnet agent; tasks A and B are independent and unblock the most. C/D can then run in parallel against a working matchmaker.

| # | Title | Scope | Likely files |
|---|---|---|---|
| **T-A** | Diagnose & fix matchmaking hang (PVE+PVP) | Reproduce locally, identify whether the hang is in `upsilonapi` controller, `upsilonbattle` engine RPC, or Reverb broadcast. Fix root cause. | `upsilonapi/app/Http/Controllers/MatchmakingController.php`, `upsilonbattle/internal/match*`, broadcast/event files, `upsilonapi/app/Events/Match*` |
| **T-B** | Bootstrap shop items & skill templates in CI script | Add a seed step to `trigger_all_ci_tests.sh` so `/shop/items` and `/skills/templates` return non-empty. Confirm existing seeders or add new ones. | `trigger_all_ci_tests.sh`, `battleui/database/seeders/*`, possibly new `util_seed_catalog.js` |
| **T-C** | Align admin CRUD test payloads with new validators | Inspect the FormRequest classes for `admin/shop-items` and `admin/skill-templates`, update the four scripts to send minimal valid payloads. | `e2e_admin_shop_item_crud.js`, `e2e_admin_skill_crud.js`, `e2e_item_grants_skill.js`, `e2e_exotic_weapon_dual_path.js`, plus the FormRequest source for context |
| **T-D** | Centralize 404 / 403 envelope mapping | Add (or fix) the global exception handler so `ModelNotFoundException` and `AuthorizationException` produce well-formed 404/403 in the standard envelope. Verify it fixes Cluster D + the skill cases in Cluster G. | `upsilonapi/app/Exceptions/Handler.php`, related policy classes |
| **T-E** | Misc small fixes (E, H, I, J, K) | Each is a 1–2 file change. Group into one PR or split as preferred: second admin in `admin_user_management`, character ID resolve in `unequip_empty_slot`, password fixture refresh, post-match progression guard, missing `character_get` route. | scenarios under `upsiloncli/tests/scenarios/edge_*`, `upsilonapi` progression controller, `upsiloncli/internal/script/routes.go` |

**Note:** Cluster F (`e2e_skill_equip_battle` turn-event timeout) is intentionally bundled into T-A — fixing the matchmaker likely resolves it; if not, T-A's owner is best positioned to chase it through the engine.

**Excluded from this plan:** `e2e_credit_economy` and the related synchronous-feedback / version-feedback work (already underway per `feedback_issues.md`).
