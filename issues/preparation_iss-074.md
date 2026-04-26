# Preparation Plan — ISS-074: Comprehensive Item System

**Date:** 2026-04-26
**Status:** In progress (Phases 0-5 done, next: Phase 6)
**Parent issue:** `ISS-074_20260423_comprehensive_item_system.md`
**Related:** ISS-067 (credits), ISS-068 (equipment), ISS-071 (V2 stats), ISS-073 (skills, follow-up), ISS-075/076 (consolidated)

---

## 🔁 Handoff status (2026-04-26)

**Done:**
- ✅ **Phase 0** — ATD scaffolding. `rule_progression` v2.1 (added MP/SP at 1 CP each, promoted JumpHeight/CritChance/CritDamage from "Planned" to active, formalised Class A vs Class B). New atoms: `entity_shop_item`, `entity_player_inventory`, `entity_character_equipment`, `api_shop_browse`, `api_shop_purchase`, `api_inventory_list`, `mec_item_buff_application`, `rule_starting_credits_1000`, `rule_quantity_cap`, `rule_stat_taxonomy`, `ui_shop`, `ui_inventory`, `ui_character_equipment_panel`, `ui_character_full_stat_panel`. Amended: `api_equipment_management` (single equip endpoint per D2), `mec_credit_spending_shop` (V2.0 fixed pricing). ISS-073 updated with skill-wiring instructions (`§ Wiring instructions (post-ISS-074)`). `atd_weave` run.
- ✅ **Phase 1** — Three migrations applied + ShopItemsSeeder seeded. Verify: `select * from shop_items` returns the 3 V2.0 rows; `users.credits` default is now 1000.
  - `database/migrations/2026_04_26_000000_extend_characters_v2_stats.php` — adds `mp, sp, jump_height, crit_chance, crit_damage` (5 new columns; HP and movement already existed and are treated as max).
  - `database/migrations/2026_04_26_000100_default_credits_1000.php` — flips default + idempotent backfill.
  - `database/migrations/2026_04_26_000200_create_item_system_tables.php` — `shop_items`, `player_inventory`, `inventory_transactions`, `character_equipment`. CHECK constraints on `slot` and `transaction_type`.
  - `database/seeders/ShopItemsSeeder.php` — wired into `DatabaseSeeder::run()`.
  - `db.md` updated.
- ✅ **Phase 2** — Models + Resources + Shop endpoints + extended progression.
  - Models: `App\Models\{ShopItem, PlayerInventory, InventoryTransaction, CharacterEquipment}`. Relations on `User::inventory`, `User::inventoryTransactions`, `Character::equipment`.
  - Resources: `ShopItemResource`, `InventoryItemResource`, `CharacterEquipmentResource`. **`CharacterResource` extended with all 9 Class A stats and `equipment` (whenLoaded).**
  - Service: `App\Services\ShopService::purchase` — DB-transactional, lockForUpdate on user row, audits to `inventory_transactions` + `credit_transactions`, exception class `ShopServiceException` with `{ERR_INSUFFICIENT_CREDITS, ERR_QUANTITY_CAP, ERR_ITEM_UNAVAILABLE}`.
  - Controller: `App\Http\Controllers\API\ShopController` (`index`, `purchase`).
  - Request: `App\Http\Requests\API\Shop\PurchaseShopItemRequest`.
  - **`ProfileController::updateCharacter` extended** to handle MP/SP/JumpHeight/CritChance/CritDamage with their CP costs (1, 1, 15, 10, 5 respectively). `UpdateCharacterRequest` rules updated to whitelist all 9 Class A stats.
- ✅ **Phase 3** — Inventory + Equipment endpoints.
  - Service: `App\Services\EquipmentService::{equip, unequip}` — DB-transactional cross-character mutual exclusivity, slot inferred from `shop_item.slot`. Exception class `EquipmentServiceException` with `{ERR_SLOT_MISMATCH, ERR_INVENTORY_NOT_OWNED, ERR_SLOT_EMPTY}`.
  - Controllers: `App\Http\Controllers\API\{InventoryController, EquipmentController}`.
  - Request: `App\Http\Requests\API\Equipment\EquipItemRequest` (single field: `item_id` UUID).
  - Policy: `CharacterPolicy::{equip, unequip}` abilities added.
  - All routes wired in `routes/api.php` under the `auth:sanctum` group; verified via `php artisan route:list`.
  - `communication.md` §2.5 (Shop/Inventory/Equipment) added; subsequent sections renumbered to 2.6–2.9.

**Not yet started:**
- ✅ **Phase 4** — Engine integration (Go). Extended `upsilonapi/api/input.go` Entity with `EquippedItems []EquippedItem` and `EquippedSkills []string`. Added `RemoveBuffsByOrigin` on `Entity`. In `upsilonapi/bridge/bridge.go`, implemented loop to project equipped items into `Forever=true` buffs with property alias mapping (`ArmorRating` -> `Armor`). Updated `UpsilonEntityResource.php` to populate `equipped_items` with eager-loaded shop items. Verified with Go and PHP unit tests.
- ✅ **Phase 5** — CLI commands (upsiloncli). Implemented 6 new endpoint structs in `internal/endpoint/endpoints.go`: `shop_browse`, `shop_purchase`, `profile_inventory`, `character_equipment_list`, `character_equip`, `character_unequip`. Updated `RegisterAll` for auto-discovery and `SyncSession` to capture `credits`. Verified via `upsiloncli routes` and compilation checks.
- ⏳ Phase 6 — Frontend dashboard rebuild. Strict componentisation per the plan; theme compliance per `req_ui_look_and_feel` + `ui_theme`.
- ⏳ Phase 7 — Frontend shop UI.
- ⏳ Phase 8 — Frontend inventory & equip UI.
- ⏳ Phase 9 — E2E scenarios (CR-21, CR-22, CR-23). **Use the CLI as the test harness — do not curl-test by hand.**
- ⏳ Phase 10 — Edge-case scenarios (EC-49 through EC-55).
- ⏳ Phase 11 — ATD finalization (atom statuses, atd_audit, atd_verify).
- ⏳ Phase 12 — Doc & polish (`CI.md` test-count bump, ISS-074 status flip).

**Open notes for the next worker:**
1. **Smoke testing:** the user has confirmed CLI is the test harness for HTTP round-trips. Defer end-to-end verification to Phases 9/10 rather than curl smoke tests. Phase 4 unit tests in Go are still appropriate (`go test ./...`).
2. **Phase 4 detail:** `entity.go` already has `RegisterBuff`, `BuffTickDown`, `Forever` flag, `OriginEntityID`. The new `RemoveBuffsByOrigin(uuid.UUID)` method needs to be added — it's a simple slice filter (see `mec_item_buff_application` atom for the canonical implementation snippet). The skill registry is referenced in atoms but not yet located in code — the worker should grep for `skillregistry` / similar before wiring weapon-as-skill, and if no clean registry exists, fall back to "register buff metadata only" and defer skill activation to ISS-073 (mentioned as a Medium-likelihood risk in §7).
3. **`mec_item_buff_application.atom.md`:** the ATD tool tries to auto-prefix new MECHANIC atoms with `mechanic_`. The file at `/workspace/upsilonbattle/docs/mec_item_buff_application.atom.md` had its frontmatter id corrected by hand. If you call `atd_update` on it again, **also re-correct the id field after the call** (or the prefix will return).
4. **Layer enum lint noise:** the linter flags `BUSINESS` as invalid; this is pre-existing tooling lag (CLAUDE.md notes the CUSTOMER → BUSINESS rename was applied to atoms but not yet to the validator). The user has confirmed lint is outdated — don't waste cycles fixing it.
5. **Stat taxonomy:** Class A = 9 stats (HP, MP, SP, Attack, Defense, Movement, JumpHeight, CritChance, CritDamage); Class B = 2 stats (AttackRange, Shield). Class B is **never** CP-upgradable — `ProfileController::updateCharacter` rejects them via the whitelist in `UpdateCharacterRequest`.
6. **Componentisation rule:** the user explicitly requested "everything componentised by responsibility". Phase 6/7/8 component lists are in §4 (Phases 6, 7, 8). Don't merge components; keep one concern per file.
7. **Schema decisions are locked:** see §2 (D1-D7). In particular: `player_inventory` has NO `character_id`; `equip` is one endpoint not three; `transaction_type` and `slot` use varchar+CHECK not Postgres ENUM.
8. **Verification matrix in §5** lists the commands to run for each phase's exit criteria.

**Files to read first when picking up:**
- This file (especially §2 Decisions, §4 Phase 4 onward).
- `issues/ISS-074_20260423_comprehensive_item_system.md` (the original, but note this prep doc supersedes several of its details).
- `upsilonapi/api/input.go` (for the additive struct extension).
- `upsilonapi/bridge/bridge.go` lines 113-132 (insertion point).
- `upsilonbattle/battlearena/entity/entity.go` lines 155-214 (buff infrastructure).
- `upsilonbattle/battlearena/property/def/item.go` (ItemProperty factory + EffectProperty).

---

---

## 1. Context

ISS-074 consolidates four sub-issues (ISS-068 equipment, ISS-074 shop, ISS-075 inventory, ISS-076 character transfer) into one end-to-end item system: **shop catalog → owned inventory → 3-slot equipment → battle buffs**.

**Foundations already in place** (no architectural net-new is needed, only wiring):
- Credit ledger: `users.credits` + `credit_transactions` table (ISS-067, `database/migrations/2026_04_23_150000_create_credit_economy_tables.php:15`).
- Item property keys: `ArmorRating`, `WeaponType`, `WeaponBaseDamage`, `WeaponRange`, `ArmorType`, `Effect`, `Durability`, `Weight`, `Value`, `Stackable/StackSize`, `ItemType` — all declared in `propertyenum.go:137-153` with default factories at `property/def/item.go`.
- **`EffectProperty`** (`property/def/item.go:118`) — items already carry an optional `Effect` (a Skill). Weapon-as-skill is therefore "wire it at equip-time", not new architecture.
- Buff system: `Forever`, `OriginEntityID`, `RegisterBuff`, `BuffTickDown`, `GetProperty` resolution at `entity/entity.go:155-214` and `property/buff.go:5-14`.
- Character V2 baseline (ISS-071): HP/ATK/DEF/MOV in DB + `spent_cp`. CP costs in `ProfileController.php:108` (HP=1, ATK=5, DEF=5, MOV=30).

**Character stat taxonomy (post user clarification, 2026-04-26):**

Two classes of stats based on whether a player can level them up:

**Class A — Character-leveled (CP-upgradable, persisted on `characters`):** 9 stats
| Property | Engine default | CP cost | Status |
|---|---|---|---|
| HP | 10/10 | 1 | ✓ already in `rule_progression` |
| MP | 10/10 | **1** *(new — matches HP, resource counter)* | ✗ to add to `rule_progression` & CP table |
| SP | 10/10 | **1** *(new — matches HP, resource counter)* | ✗ to add to `rule_progression` & CP table |
| Attack | 3 | 5 | ✓ already in `rule_progression` |
| Defense | 0 | 5 | ✓ already in `rule_progression` |
| Movement | 3/3 | 30 | ✓ already in `rule_progression` |
| JumpHeight | 2 | 15 | ✓ in `rule_progression` (was "Planned") — promote to active |
| CritChance | 0 | 10 (per +1%) | ✓ in `rule_progression` (was "Planned") — promote to active |
| CritDamage | 0 | 5 (per +5%) | ✓ in `rule_progression` (was "Planned") — promote to active |

**Class B — Effective-only (granted by items/buffs only, never CP-upgradable):** 2 stats
| Property | Engine default | Source |
|---|---|---|
| AttackRange | 1 | Items / buffs only (e.g. ranged weapon WeaponRange) |
| Shield | 0/0 | Items / buffs only |

**Implications for ISS-074 scope:**
1. **`rule_progression` atom needs amending** (add MP/SP at 1 CP each; promote the three "Planned" exotics to active). Done in Phase 0 alongside the other ATD updates.
2. **`characters` schema needs 7 new columns** (MP, MaxMP, SP, MaxSP, JumpHeight, CritChance, CritDamage) — done in Phase 1 alongside the inventory tables. *(This was previously flagged as ISS-074.A follow-up; user's clarification pulls it back into scope.)*
3. **`ProfileController::upgradeCharacter` extended** with the 5 new CP-costed stats (MP/SP/JumpHeight/CritChance/CritDamage) — done in Phase 2.
4. **Dashboard renders an 11-row stat panel** (9 Class A + 2 Class B), with Class B labeled "Item-granted only" when no equipment contributes.

---

## 2. Decisions locked (post Q&A)

| # | Decision | Origin |
|---|---|---|
| **D1** | Drop `character_id` from `player_inventory`. Equipment slot binding lives **only** in `character_equipment` (one row per character, three FK columns to `player_inventory.id`). Inventory = ownership; equipment = active binding. | OQ-1 |
| **D2** | **Single equip endpoint:** `POST /v1/profile/character/{id}/equip` with body `{ item_id }`. Slot inferred from `shop_items.slot`. `DELETE /v1/profile/character/{id}/unequip/{slot}` for the inverse. | OQ-2 |
| **D3** | **1000 starting credits is a permanent V2 design decision.** Captured in `rule_starting_credits_1000` ATD atom; baked into the `users.credits` migration default and registration narrative. | OQ-3 |
| **D4** | **Cut from this issue:** stacking quantity > 1 with UI affordance, two-handed weapon → utility slot restriction. **Keep:** weapon-as-skill / item `Effect` resolution (architecturally trivial since `EffectProperty` already exists). **Drop:** `item_usage_stats` table (the user clarified this is statistics-as-metrics, not stats-as-properties — confusion in the issue text). | OQ-5 + clarification |
| **D5** | **Add `equipped_skills` field to ArenaStartRequest payload** alongside `equipped_items`, but leave it **empty/nil**. Pre-wires the structure for ISS-073 (skill system) so that issue is purely data-side. ISS-073 must be updated with a comment block explaining: where to plug the resolver (`bridge.go`'s entity-bootstrap loop), how to source skills from a future `character_skills` table, and that `Entity.RegisterSkill()` is already the registration call to use. | OQ-4 |
| **D6** | API contract change: **approved**. New endpoints listed in §4 Phase 2/3; ArenaStartRequest gains `equipped_items` + `equipped_skills` (additive); `communication.md` updated in same PR. | OQ-6 |
| **D7** | Persist equipment slot enums as **varchar + CHECK** constraint, not Postgres ENUM (matches `match_participants.status` style; sqlite test parity). | engineering call |

---

## 3. Critical files

### upsilonapi (Laravel — `battleui/` repo)
- `app/Models/User.php:33` — `credits` fillable; add `inventory()` / `transactions()` relations
- `app/Models/Character.php:26-69` — add `equipment()` HasOne relation
- `app/Http/Controllers/API/MatchMakingController.php:142-148` — payload assembly
- `app/Http/Resources/API/Upsilon/UpsilonEntityResource.php:18-39` — **inject `equipped_items` + `equipped_skills:[]` here**
- `app/Http/Resources/CharacterResource.php` — extend with equipment summary
- `app/Policies/CharacterPolicy.php:16,24,32` — add `equip`/`unequip` abilities
- `app/Traits/ApiResponder.php:18-26` — envelope helper (reuse)
- `database/migrations/2026_04_23_150000_create_credit_economy_tables.php:15` — flip default 0 → 1000 (or new follow-up migration if seeds in flight)
- `routes/api.php:26-86` — new routes inside `auth:sanctum` `/v1` group

### upsilonbattle (Go engine)
- `upsilonapi/api/input.go:20-32` — extend `Entity` with `EquippedItems []EquippedItem` and `EquippedSkills []string` (string UUIDs; left empty for ISS-074)
- `upsilonapi/bridge/bridge.go:113-132` — **insert buff-loading + Effect-skill registration loop here** before `AddEntity`
- `battlearena/entity/entity.go:190-214` — add `RemoveBuffsByOrigin(uuid.UUID)` helper (does not exist today)
- `battlearena/property/def/item.go:118-176` — `EffectProperty` exists; reuse for weapon-as-skill resolution at equip-time
- `battlearena/property/def/item.go:9-39` — reuse `ArmorRating()`, `WeaponBaseDamage()`, `WeaponRange()` factories when materializing properties from JSON
- `battlearena/property/def/entity.go:91-107` — `PropertiesForCharacter()` is the canonical 9-stat list; UI must render against this

### battleui (Vue / Inertia frontend)
- `routes/web.php:27-33` — register `/shop` and `/inventory` Inertia pages
- `resources/js/Components/CharacterRoster.vue:154-174` — **rebuild stat panel** (9 rows + equipment column)
- `resources/js/Components/IdentitySection.vue:68-77` — credits already shown
- `resources/js/services/game.js` — pattern for new `services/shop.js`, `services/inventory.js`
- `resources/js/Components/Modals/ConfirmModal.vue` — reuse for purchase / equip confirms
- `resources/js/Components/TacticalHeader.vue` — add Shop / Inventory nav

### upsiloncli (Go CLI)
- `internal/endpoint/endpoints.go:100-108` — `CharacterUpgrade` template for new `ShopList`, `ShopPurchase`, `InventoryList`, `Equip`, `Unequip`, `CharacterEquipmentShow`
- `internal/api/client.go:19-80` — envelope client unchanged
- `internal/script/bridge.go:88-111` — `upsilon.call()` is generic; **no JS-bridge changes**

### Tests
- `upsiloncli/tests/scenarios/e2e_credit_economy.js` — E2E template
- `upsiloncli/tests/scenarios/edge_prog_allocation_no_wins.js` — try/catch edge template
- `tests/edge_case_report.sh`, `tests/ci_report.sh` — `check_brd` / `check_edge` mappings

---

## 4. Phase plan

Phases are sequential where there's a hard dependency, parallel-friendly otherwise. Each phase ends with a concrete verification.

### Phase 0 — ATD scaffolding, rule_progression amend, cross-issue notes (≈1h)
- **Amend `docs/rule_progression.atom.md`** (currently STABLE):
  - Add **MP (+1): 1 CP** and **SP (+1): 1 CP** to the standard cost table.
  - Promote CritChance / CritMultiplier / JumpHeight from "Planned" to active (drop the `(Planned)` qualifier).
  - Note explicitly that **AttackRange and Shield are not CP-upgradable** — they are only granted via items or buff effects.
  - Bump version 2.0 → 2.1.
- Create DRAFT atoms via `mcp__atd__atd_update`:
  - **ENTITY:** `entity_shop_item`, `entity_player_inventory`, `entity_character_equipment`
  - **API:** `api_shop_browse`, `api_shop_purchase`, `api_inventory_list`, `api_character_equipment_show`, `api_character_equip`, `api_character_unequip`
  - **MECHANIC:** `mech_item_buff_application`, `mech_item_effect_skill_registration`
  - **RULE:** `rule_equipment_slot_validation`, `rule_starting_credits_1000`, `rule_quantity_cap`, `rule_stat_taxonomy_classA_classB`
  - **UI:** `ui_shop`, `ui_inventory`, `ui_character_equipment_panel`, `ui_character_full_stat_panel`
- Run `mcp__atd__atd_weave`.
- **Update ISS-073** (`issues/ISS-073_*.md`) with a "Wiring instructions for skills (post-ISS-074)" block: ArenaStartRequest carries an empty `equipped_skills []string` field today; ISS-073 needs to (a) add a `character_skills` join table, (b) populate the field in `UpsilonEntityResource.php`, (c) extend `bridge.go` entity-bootstrap to call `entity.RegisterSkill()` per ID, (d) drop the `// reserved for ISS-073` comment in `api/input.go`.
- **Verify:** `mcp__atd__atd_query(field="id", search="shop")` returns the new atoms; `atd_lint` clean; `rule_progression` reads correctly.

### Phase 1 — Database & migrations (Laravel) (≈1.5h)
- Migration `*_extend_characters_v2_stats.php`: add columns to `characters` —
  - `mp` int default 10, `max_mp` int default 10
  - `sp` int default 10, `max_sp` int default 10
  - `jump_height` int default 2
  - `crit_chance` int default 0  *(percent, 0-100)*
  - `crit_damage` int default 0  *(percent multiplier, 0-N)*
  - Backfill defaults for existing rows; update `Character::generateInitialRoster()` and `Character::rerollStats()` accordingly.
- Migration `*_create_item_system_tables.php`:
  - `shop_items` (UUID PK, name, type, slot varchar+CHECK in {armor,utility,weapon}, properties JSON, cost int, available bool default true, version)
  - `player_inventory` (UUID PK, FK user CASCADE, FK shop_item, quantity default 1, purchased_at; **no `character_id`**)
  - `inventory_transactions` (UUID PK, FK user, FK shop_item, quantity, credits_spent, transaction_type varchar+CHECK)
  - `character_equipment` (character_id PK + FK CASCADE, armor_item_id / utility_item_id / weapon_item_id nullable FK to player_inventory ON DELETE SET NULL)
- Migration `*_default_credits_1000.php`: change `users.credits` default to 1000 + idempotent backfill `update users set credits = 1000 where credits = 0`.
- Seeder `ShopItemsSeeder`: 3 V2.0 items with deterministic UUIDs:
  - Basic Armor — slot=armor, properties={ArmorRating:5}, cost=200
  - Basic Sword — slot=weapon, properties={WeaponBaseDamage:5, WeaponType:"One-Handed Melee", WeaponRange:1}, cost=300
  - Swift Boots — slot=utility, properties={Movement:1}, cost=150
- Update `db.md` with new tables + ER diagram.
- **Verify:** `php artisan migrate:fresh --seed` clean; `select * from shop_items` shows 3 rows; new registration → `users.credits=1000`.

### Phase 2 — Shop endpoints + extended progression (Laravel) (≈2.5h)
- **Extend `ProfileController::upgradeCharacter`** with the 5 new CP-costed stats: MP (1 CP), SP (1 CP), JumpHeight (15 CP), CritChance (10 CP per +1%), CritDamage (5 CP per +5%). Update validation request, CP cost calculation, and persistence. Reference `rule_progression` 2.1.
- Eloquent models: `ShopItem`, `PlayerInventory`, `InventoryTransaction`, `CharacterEquipment` with relations.
- `ShopService::purchase(User, ShopItem, qty=1)` — DB-transactional: check balance, debit credits, upsert inventory row, log to `inventory_transactions` and `credit_transactions` (source=`shop_purchase`).
- `ShopController`:
  - `GET /v1/shop/items` → `ShopItemResource::collection`
  - `POST /v1/shop/purchase` → `{shop_item_id, quantity?}`, returns `{credits, inventory_item}`
- Resources: `ShopItemResource`, `InventoryItemResource`. Extend `CharacterResource` with the 7 new stat columns.
- `communication.md` §2 updated.
- **Verify:** Pest/PHPUnit feature tests for `ShopService::purchase` (insufficient funds 422, success, exact-balance, quantity cap 99 → 422); upgrade endpoint accepts new stats and rejects on cap overrun; `curl` round-trips.

### Phase 3 — Inventory & equipment endpoints (Laravel) (≈2h)
- `InventoryController`:
  - `GET /v1/profile/inventory` → list owned items + which character (if any) has them equipped (LEFT JOIN `character_equipment`)
  - `GET /v1/profile/character/{id}/equipment` → 3-slot view
- `EquipmentController` (single endpoint per D2):
  - `POST   /v1/profile/character/{id}/equip` body `{item_id}`, slot inferred from `shop_items.slot`
  - `DELETE /v1/profile/character/{id}/unequip/{slot}` slot ∈ {armor,utility,weapon}
- Service-layer validation:
  - Ownership: user owns character (Policy) + user owns inventory row
  - Slot match: item's `slot` matches request (or any slot, since inferred)
  - **Mutual exclusivity (cross-character):** if item is equipped on character A and re-equipped on character B, A's slot is cleared atomically
  - Quantity cap: 99 per item per user
- Add `equip` / `unequip` Policy abilities.
- `communication.md` §2 updated.
- **Verify:** Feature tests for swap-armor / cross-character move / 403 on cross-user / 422 on slot mismatch.

### Phase 4 — Engine integration (Go) (≈2h)
- `upsilonapi/api/input.go`: extend `Entity` with:
  ```go
  EquippedItems  []EquippedItem `json:"equipped_items"`
  EquippedSkills []string       `json:"equipped_skills"` // reserved for ISS-073
  ```
  And new `EquippedItem { ItemID string; Name string; Slot string; Properties map[string]any }`.
- `upsilonbattle/battlearena/entity/entity.go`: add `func (e *Entity) RemoveBuffsByOrigin(originID uuid.UUID)` — filters `e.Buffs` keeping only those with different `OriginEntityID`.
- `upsilonapi/bridge/bridge.go:130-131`: after base properties set, for each `entity.EquippedItems`:
  1. Build `TemporaryProperties{Forever:true, OriginEntityID:item.ItemID}` populated via `def.ItemProperty(...)` factory per JSON key (handles ArmorRating, WeaponBaseDamage, WeaponRange, AttackRange, JumpHeight, Movement, etc.).
  2. `RegisterBuff` on entity.
  3. **If the item's properties include `Effect`** (a skill ID), resolve via the existing skill registry and call `entity.RegisterSkill(skill)` — this is the weapon-as-skill path. (Falls under `mech_item_effect_skill_registration`.)
- Laravel `UpsilonEntityResource.php:18-39`: include `equipped_items` (eager-load `Character::equipment.armorItem.shopItem` etc.) and `equipped_skills:[]`.
- Engine unit tests: arena init with armor + sword → `entity.GetProperty("Armor").CValue` reflects +5; `RemoveBuffsByOrigin` strips them; item with `Effect` → skill registered.
- **Verify:** `cd upsilonbattle && go test ./...` and `cd upsilonapi && go test ./bridge/...` green.

### Phase 5 — CLI commands (≈1.5h)
- New endpoint structs (`upsiloncli/internal/endpoint/endpoints.go` style): `ShopList`, `ShopPurchase`, `InventoryList`, `CharacterEquipmentShow`, `Equip`, `Unequip`. Register in `RegisterAll`.
- No JS-bridge changes (`upsilon.call()` is generic).
- **Verify:** REPL: login → `shop_list` → `shop_purchase` → `inventory_list` → `equip` → `character_equipment_show`; balance/state reflect each step.

### Phase 6 — Frontend dashboard: 11-stat panel + equipment (≈3h)
**This is the substantive UI lift the user flagged. Strict componentisation by responsibility.**

Theme compliance (per `req_ui_look_and_feel` + `ui_theme`):
- All panels use `bg-upsilon-gunmetal/30` + `backdrop-blur-md`, 1px cyan/30 or magenta/30 border, 2px corner accents (`border-t-2 border-l-2`).
- Titles: Orbitron `uppercase tracking-[0.3em]`.
- Hover: increase border-opacity + `shadow-glow-cyan` or `shadow-glow-magenta`.
- Sci-fi/post-apoc terminology in copy ("Acquire", "Hardwire", "Cache", "Link terminated").

New components (single-responsibility):
- `Components/Character/StatRow.vue` — one stat: label, base, contribution(s), effective. Props: `{ label, base, contributions[], effective, cpCost?, classB? }`.
- `Components/Character/CharacterStatPanel.vue` — composes 11 `StatRow` (Class A: HP, MP, SP, ATK, DEF, MOV, JumpHeight, CritChance, CritDamage; Class B: AttackRange, Shield).
- `Components/Character/CpEconomySummary.vue` — spent / max CP bar (`100 + total_wins*10`).
- `Components/Character/EquipmentSlotPill.vue` — one slot: icon, item name, "Empty" placeholder, click handler.
- `Components/Character/CharacterEquipmentPanel.vue` — composes 3 `EquipmentSlotPill` (armor/utility/weapon).
- `Components/Character/CharacterCard.vue` — composes name, `CharacterStatPanel`, `CpEconomySummary`, `CharacterEquipmentPanel`. **Replaces the existing inline character block in `CharacterRoster.vue`.**
- `composables/useCharacterStats.js` — pure function: `(character, equipment) → { class_a_rows, class_b_rows, contribution_breakdown }` so the component stays presentation-only.

`CharacterRoster.vue` becomes a thin layout that v-fors `CharacterCard`. The broken V1 `level` math is removed.

**Verify:** browser smoke test: character with no items shows 11 stats with Class B labelled "Item-granted only"; equipping basic sword surfaces a +5 WeaponBaseDamage contribution and an effective ATK delta; unequipping reverts; CP bar reflects post-ISS-071 cap math.

### Phase 7 — Frontend shop UI (≈2.5h)
Theme: same neon-in-the-dust panels.

New components:
- `Components/Shop/ShopItemCard.vue` — one item: name, slot icon, properties summary, cost, action button. Props: `{ item, ownedQty, canAfford }`. Emits: `@purchase`.
- `Components/Shop/ShopGrid.vue` — composes `ShopItemCard` v-for, handles loading state.
- `Components/Shop/PurchaseConfirmModal.vue` — wraps `ConfirmModal` with shop-specific copy ("Acquire {item}? -{cost} credits").
- `Pages/Shop.vue` — composes `ShopGrid` + `PurchaseConfirmModal`, calls `services/shop.js`, optimistic credit decrement on success.
- `services/shop.js` — `listItems()`, `purchase(itemId, qty=1)`. Mirrors `services/game.js` envelope handling.
- Nav link added in `TacticalHeader.vue`.

**Verify:** end-to-end click-through in browser; credits decrement; item appears in `IdentitySection` credit balance and inventory page.

### Phase 8 — Frontend inventory & equip UI (≈3h)
New components:
- `Components/Inventory/InventoryRow.vue` — one inventory entry: name, qty, slot, equipped-on (character name or "—"). Emits: `@equip`, `@unequip`.
- `Components/Inventory/InventoryTabs.vue` — slot tabs (All / Armor / Utility / Weapon). Props: `{ activeTab }`. Emits: `@change`.
- `Components/Inventory/InventoryList.vue` — composes `InventoryTabs` + filtered `InventoryRow`s.
- `Components/Inventory/EquipDrawer.vue` — slide-out: shows character + 3 slots + compatible items, equip/unequip buttons. Self-contained equip flow.
- `Pages/Inventory.vue` — composes `InventoryList` + `EquipDrawer`, calls `services/inventory.js`.
- `services/inventory.js` — `listInventory()`, `getEquipment(charId)`, `equip(charId, itemId)`, `unequip(charId, slot)`.
- The slot pill in `Components/Character/EquipmentSlotPill.vue` (Phase 6) emits `@click` → opens `EquipDrawer` for that character/slot.

**Verify:** equip basic sword → `CharacterCard` panel reflects +5 WeaponBaseDamage row + effective stat update; unequip reverts; equipping on character B atomically frees character A's slot (no stale state).

### Phase 9 — E2E customer scenarios (≈2h)
Add to `upsiloncli/tests/scenarios/`:
- `e2e_starting_credits_1000.js` (CR-21) — `[[rule_starting_credits_1000]]`: register → assert `user.credits == 1000`
- `e2e_shop_browse_purchase.js` (CR-22) — `[[api_shop_purchase]]`: register → list → buy Basic Armor → credits=800 + inventory contains armor
- `e2e_inventory_equip_battle.js` (CR-23) — `[[api_character_equip]]` + `[[mech_item_buff_application]]`: register → buy + equip armor + sword → start arena → assert entity buffs include both ItemIDs as `Forever=true` origins
- All end with `[SCENARIO_RESULT: PASSED]`.
- Update `tests/ci_report.sh` `check_brd` mapping; bump `CI.md` table.
- **Verify:** `docker compose -f docker-compose.ci.yaml exec tester /bin/sh ./tests/run_all_scenarios.sh` shows 3 new CRs PASSED.

### Phase 10 — Edge-case scenarios (≈2h)
Add (try/catch-for-expected-failure pattern):
- `edge_shop_insufficient_credits.js` (EC-49) — buy Basic Sword with 100 credits → 422
- `edge_shop_unknown_item.js` (EC-50) — purchase random UUID → 404
- `edge_equip_wrong_slot.js` (EC-51) — equip armor item via API expecting weapon path → 422 (covers slot-mismatch validation)
- `edge_equip_unowned_item.js` (EC-52) — equip another user's inventory row → 403
- `edge_equip_unowned_character.js` (EC-53) — equip own item to someone else's char → 403
- `edge_unequip_empty_slot.js` (EC-54) — unequip a slot with nothing → 404
- `edge_quantity_cap_99.js` (EC-55) — buy 100th of same item → 422
- Update `tests/edge_case_report.sh` and `CI.md` table.
- **Verify:** all 7 ECs end with `[SCENARIO_RESULT: PASSED]`.

### Phase 11 — ATD finalization (≈30 min)
- For each new atom: confirm `@spec-link` coverage in code + `@test-link` from scenarios.
- Flip atoms `DRAFT → STABLE`.
- `mcp__atd__atd_audit` and `atd_verify` clean.
- **Verify:** `mcp__atd__atd_stats` orphan count not increased.

### Phase 12 — Doc & polish (≈45 min)
- `communication.md` API summary table updated with all new endpoints.
- `db.md` ER diagram + table specs updated.
- `CI.md` test count incremented (3 CRs + 7 ECs).
- `issues/ISS-074_*.md` → status `Resolved`.
- `issues/ISS-073_*.md` updated with the "wiring instructions for skills" comment block (per Phase 0).
- **Verify:** `git diff --stat` shows expected files only.

---

## 5. Verification matrix (end-to-end)

| Check | Command |
|---|---|
| Migrations clean | `cd battleui && php artisan migrate:fresh --seed` |
| Engine unit tests | `cd upsilonbattle && go test ./...` |
| Bridge buff loading | `cd upsilonapi && go test ./bridge/...` |
| Laravel feature tests | `cd battleui && php artisan test --filter=Item` |
| CLI smoke | `upsiloncli` REPL: login → shop_list → shop_purchase → equip |
| E2E suite | `docker compose -f docker-compose.ci.yaml exec tester /bin/sh ./tests/run_all_scenarios.sh` |
| Reports | `./tests/ci_report.sh > ci_report.md && ./tests/edge_case_report.sh > edge_case_report.md` |
| ATD health | `mcp__atd__atd_stats`, `atd_lint`, `atd_audit` |

---

## 6. Adjacent findings & follow-ups

- **CP cost extension is now in scope** (per user clarification 2026-04-26): `rule_progression` is amended in Phase 0 (add MP/SP at 1 CP each, promote JumpHeight/CritChance/CritDamage from "Planned" to active); `characters` schema gains 7 new columns in Phase 1 (MP, MaxMP, SP, MaxSP, JumpHeight, CritChance, CritDamage); `ProfileController::upgradeCharacter` extended in Phase 2 to score the 5 new stats. **AttackRange and Shield remain item/buff-only — never CP-upgradable.**
- **`equipped_skills` field is reserved**, populated to `[]` by Laravel and ignored by the Go bridge. Wiring instructions added to ISS-073.
- **Item `Effect` resolution** (weapon-as-skill) is implemented in this issue, but no V2.0 catalog item carries an `Effect` yet. First exercised when a weapon-with-skill item is added to the seeder (V2.1 content).

---

## 7. Risk register

| Risk | Likelihood | Mitigation |
|---|---|---|
| Buff resolution doesn't sum across multiple buffs of same property | Low | Verified in `entity.go:155-165` — `GetBuffsFor` + `ApplyBuff` per buff. Engine unit test in Phase 4 covers it. |
| Mass-update of `users.credits=1000` surprises real users | Low (no real users yet) | Backfill is idempotent (`where credits = 0`) and documented in migration filename. |
| `UpsilonEntityResource` change breaks an in-flight match if deployed mid-game | Medium | Change is additive (new fields). Old engine ignores unknown fields per Go JSON unmarshal default. Smoke test in Phase 4. |
| Frontend stat panel rebuild breaks existing dashboard layout | Medium | Roster card stays in same `lg:col-span-3` slot; only inner content changes. No layout shift. |
| `Effect` resolution at equip-time depends on a skill registry I haven't traced | Medium | Phase 4 includes a confirmation read of the skill registry before wiring. If the registry isn't trivially callable, we fall back to "register Effect as buff metadata only" and defer skill-effect activation to ISS-073. |
| `RemoveBuffsByOrigin` interaction with `BuffTickDown` (Forever flag) | Low | `Forever=true` already exempts buffs from tick-down; `RemoveBuffsByOrigin` is unconditional removal. Test coverage in Phase 4. |

---

## 8. Cross-cutting principles (apply to every phase)

- **Componentise by responsibility.** Each component / class / function owns one concern. No "god components". Frontend in particular: prefer many small Vue SFCs (one stat row = one component, one slot pill = one component, one shop card = one component) composed by container components. Backend: services per concern (`ShopService` purchases, `EquipmentService` slot bindings, etc.); controllers stay thin.
- **Theme compliance** for every Vue surface added/touched: `req_ui_look_and_feel` ("Neon in the Dust") + `ui_theme` (gunmetal panels, cyan/magenta borders, corner accents, Orbitron uppercase tracking, sparing glow, sci-fi/post-apoc copy). Audit checklist: panel bg, border, corner accent, title typography, hover state, terminology.
- **Crash early** (per CLAUDE.md): no silent fallbacks for missing fields. Insufficient credits → 422 with `meta.reason`. Slot mismatch → 422. Unknown item → 404. Frontend surfaces errors verbatim from the standard envelope.
- **Bidirectional ATD links**: every new function/component carries `@spec-link [[atom_id]]`; every new test carries `@test-link [[atom_id]]`. Run `atd_trace` per atom in Phase 11.
- **Additive engine contract**: ArenaStartRequest fields are additive; old engine versions tolerate unknown fields (Go default JSON unmarshal). No breaking change to the `[[api_standard_envelope]]` shape.

---

## 9. Out of scope (explicit, with parking)

- Item usage statistics tracking (the issue's `item_usage_stats` table) — **dropped** as a misnomer per user clarification.
- Stacking quantity > 1 with UI affordance — V2.1 content polish.
- Two-handed weapon → utility slot restriction — V2.1 rules expansion.
- ISS-073 skill wiring (handled via the reserved `equipped_skills:[]` field + ISS-073 comment block).
- Extending `characters` schema with SP/MP/Shield/AttackRange/JumpHeight + their CP costs — recommended new ISS-074.A follow-up.
- Inventory filtering / sorting beyond category tabs — V2.1 polish.
