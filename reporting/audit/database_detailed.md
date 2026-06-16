# Database — Detailed Investigation Report

**Auditor:** Principal Systems Architect · **Date:** 2026-06-16
**Scope:** `battleui/database/migrations/**` (29 migrations), `app/Models/**`
(12 models), seeders/factories, `db.md`, match-resurrection persistence
(ISS-054), inventory/credit/skill schema (ISS-073/074/075/086)
**Method:** static analysis + `atd` CLI; PostgreSQL via Laravel migrations.

---

## 1. Schema coverage vs V2

The V2 schema is **substantially built and migrated**:

| V2 system | Migration(s) | Models |
|---|---|---|
| x10 stats / 100 CP | `extend_characters_v2_stats`, `add_spent_cp_to_characters` | `Character` |
| Credit economy | `create_credit_economy_tables`, `default_credits_1000` | `CreditTransaction` |
| Item system | `create_item_system_tables`, `add_skill_template_to_shop_items` | `ShopItem`, `PlayerInventory`, `CharacterEquipment`, `InventoryTransaction` |
| Skill system | `create_skill_templates_table`, `create_character_skills_table` | `SkillTemplate`, `CharacterSkill` |
| Reroll/roulette | `add_reroll_count_to_users`, `add_roulette_used_to_characters` | |
| Match resurrection | `add_version_to_game_matches`, `change_turn_to_big_integer` | `GameMatch` |

`db.md` is reasonably current (last updated 2026-04-27, ISS-073/086) and documents
the relational boundaries, but predates the late-April/May issue wave
(ISS-088/090/092…) and does not describe the credit/inventory normalisation in
full detail — a moderate doc-staleness, not a gap.

## 2. Match resurrection (ISS-054) — verified Resolved

The persistence substrate is real and correctly shaped:
- `game_matches`: `game_state_cache` JSON, `grid_cache` JSON, `version`
  (bigInteger, default 0 — optimistic concurrency), `turn` bigint, lifecycle
  timestamps, `game_mode` enum, `winning_team_id`.
- `upsilonapi/bridge/bridge_resurrect.go::ResurrectArena` consumes a persisted
  board state and rebuilds the arena: `validateResurrectRequest` (UUID + roster
  checks), `initResurrectedArena`, initiative-timeline + version recovery, and a
  full `buildResurrectionBoardState` response.

ISS-054's **Resolved** status is **accurate** — one of the few status labels that
matches reality. The `version` column gives resurrection a real conflict-detection
mechanism rather than last-write-wins.

## 3. ISS-098 — the masking infrastructure already exists in the DB layer

This is the most important cross-cutting database finding. A masked,
rotation-capable identifier is present and wired everywhere **except** the engine
output:
- `add_ws_channel_key_to_users` migration; `User` model auto-generates it
  (`User.php:50-51`) and regenerates on login (`AuthController.php:43`).
- It is the WS routing key: `CodeDiscoveryService.php:245`
  `'private-user.{ws_channel_key}'`, `MatchMakingController.php:183` builds
  participant ids from `ws_channel_key`, exposed via `UserResource`.

So Laravel already speaks in masked keys. The ISS-098 leak exists **solely**
because `upsilonapi/api/output.go` emits the raw engine `ControllerID` (User UUID)
instead of the masked key. **The fix needs no new schema or infra — thread the
existing `ws_channel_key` (or an entity-scoped token) through the board-state
DTO.** This materially lowers the cost of the highest-severity open issue.

## 4. Data-integrity & GDPR

**Sound**
- GDPR/privacy honoured at schema level: `add_resident_data_to_users`
  (`full_address`, `birth_date` flagged Private in `db.md`),
  `add_soft_deletes_to_users` (anonymise/soft-delete for right-to-erasure),
  matching `requirement_customer_user_id_privacy` and the GDPR E2E.
- Referential integrity: `add_cascade_to_characters_player_id`,
  `make_player_id_nullable_in_match_participants` show deliberate FK lifecycle
  handling (characters cascade; participants survive user removal).
- `add_indexes_for_admin_performance` (ISS-053) — admin pagination is indexed,
  not naive `OFFSET` scans.

**Risk / not scalable**
- **Migrations carry no ATD traceability** — schema changes are not linked to the
  `entity_*`/`rule_*` atoms they realise, so the "every change traces to a spec"
  mandate stops at the database boundary. Schema is the one layer with effectively
  zero spec/test linkage.
- **`game_state_cache` / `grid_cache` as opaque JSON blobs**: fine for
  resurrection, but there is no schema/version contract on the blob itself beyond
  the integer `version`. As the engine state evolves (V2 temporary entities,
  effects), an unversioned blob shape risks silent resurrection corruption across
  engine releases. Recommend an explicit serializer schema version inside the blob
  (ties to `upsilonserializer` + `gamestate_version.go`).
- **Credit/inventory transaction tables** exist (`CreditTransaction`,
  `InventoryTransaction`) — good for auditability — but their reconciliation rules
  (no negative balances, idempotent purchase) live in PHP/engine, not DB
  constraints; under concurrency this leans on application correctness alone.

## 5. Recommendations (no code changed here)
1. Resolve ISS-098 by threading `ws_channel_key` through `output.go` — schema
   already supports it.
2. Add a serializer schema-version field inside `game_state_cache` to protect
   resurrection across engine upgrades.
3. Add `@spec-link` tags from migrations to their `entity_*`/`rule_*` atoms to
   extend traceability to the data layer.
4. Consider DB-level guards (CHECK / unique) for credit balances and inventory
   quantities to back up the application invariants.
