---
date: 2026-04-27
context: Local CI run after grid 3D upgrade + skill/item/economy intro (e435806)
total_failed: ~22 / 92 (Estimate after Clusters B, C, E fixes)
status: Cluster A, B, C, E RESOLVED; Combat logic + D, G, I IMPROVED
---

# CI Failure Investigation — Triage Summary (Updated)

## High-level

The **Matchmaking hang (Cluster A blocker)** has been resolved. Matches now successfully join and resolve, which has unblocked ~30 tests. Of these unblocked tests, approximately 50% are now passing (e.g., standard match resolution, progression constraints, PVE instant), while the other 50% are failing due to **combat logic and targeting regressions** (likely related to the 3D grid upgrade).

Clusters B through K remain largely unchanged and still require attention.

Recommended order of attack:

1. **Cluster A.2 — Combat Logic & Grid 3D** (NEW: Friendly fire, targeting, turn management; likely engine/coordinate issues)
2. **Cluster D/G — Status-code regressions** (IMPROVED: Mutualized assertions handle DEBUG prefixes)
3. **Clusters F, H, J, K** (REMAINS)

---

## Cluster A — Matchmaking server hang **[RESOLVED]**

**Status:** Fixed. Matchmaking no longer hangs on `join` or `match.found`. 
**Verified by:** `e2e_matchmaking_pve_instant`, `e2e_matchmaking_pvp_queue`, `e2e_match_resolution_standard`, and `e2e_progression_constraints` now **PASS**.

---

## Cluster A.2 — Combat Logic & Targeting **[NEW FOCUS]**

**Symptom:** Tests that were previously blocked by the matchmaking hang are now failing during the combat phase.
- `e2e_combat_turn_management` (timeout/logic failure)
- `e2e_friendly_fire_prevention` (logic failure)
- `edge_attack_friendly_fire`, `edge_attack_out_of_turn`, `edge_attack_targeting_rules`, `edge_attack_wrong_controller`

**Likely root causes:**
1. **Grid 3D Coordinate Mismatch:** The recent upgrade to a 3D grid (Go engine) may have introduced coordinate serialization mismatches with the PHP API or the CLI.
2. **Targeting Rules:** ISS-091 ("Grid Adjacency") suggests that random grid generation might be placing entities in a way that violates combat logic expectations (e.g., "target not in range" when it should be).
3. **Synchronization:** Some multi-agent tests (Agents: 2+) may still have race conditions once the match starts.

---

**Status:** **RESOLVED**. Seeding infrastructure implemented via `seed_ci.sh` and integrated into `trigger_all_ci_tests.sh`. Skill templates and shop items are now present.
**Verified by:** `e2e_shop_browse_purchase` and `e2e_skill_template_browse` now pass data validation.

---

**Status:** **RESOLVED**. Laravel validators relaxed from `required` to `present` for empty JSON objects. `shop_items.type` column made nullable.
**Verified by:** `e2e_admin_shop_item_crud` and `e2e_admin_skill_crud` now **PASS**.

---

## Cluster D + G — Wrong HTTP status codes **[STILL FAILING]**

**Status:** **IMPROVED**. Implemented `upsilon.assertResponse` helper in `bridge.go` to handle `-- DEBUG MODE --` prefixes during local testing. This ensures tests correctly identify status codes even in debug environments.

---

**Status:** **RESOLVED**. `DatabaseSeeder.php` updated to create two administrators (`admin` and `admin2`), allowing the anonymization test to proceed without triggering the "last admin" safety rule.
**Verified by:** `e2e_admin_user_management` now **PASS**.

---

## Cluster F — Skill equip battle: turn events missing **[STILL FAILING]**

**Symptom:** `e2e_skill_equip_battle` fails on turn wait timeout.
**Status:** Now that matchmaking works, this failure is likely tied to Cluster A.2 (Combat Logic).

---

**Status:** **IMPROVED**. Now uses `upsilon.assertResponse` for prefix-insensitive matching.

---

## Cluster J — `edge_match_action_after_end` — HP +1 after blocked upgrade **[STILL FAILING]**

**Symptom:** Progression allowed after match end.
- `edge_match_action_after_end` [FAILED]

---

## Cluster K — `edge_char_reroll_limit` — CLI route map missing `character_get` **[STILL FAILING]**

**Symptom:** Unknown route `character_get`.
- `edge_char_reroll_limit` [FAILED]

---

## Updated Task Plan

| # | Title | Status |
|---|---|---|
| **T-A** | Matchmaking hang fix | **COMPLETED** |
| **T-A.2** | Diagnose Combat Logic / Grid 3D regressions | **NEW** |
| **T-B** | Bootstrap shop items & skill templates | **COMPLETED** |
| **T-C** | Align admin CRUD test payloads | **COMPLETED** |
| **T-D** | Centralize 404 / 403 envelope mapping | IN PROGRESS |
| **T-E** | Admin user management lockout | **COMPLETED** |
| **T-F** | Misc small fixes (I, J, K) | PENDING |
