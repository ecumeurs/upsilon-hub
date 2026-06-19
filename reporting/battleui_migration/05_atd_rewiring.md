# 05 — ATD Rewiring

> Added per user note: *"we might need to rewire ATD in this as well."* Confirmed — the rewrite
> touches ATD in three distinct places. Skipping this would silently break the project's
> traceability graph (Atoms ⇄ code).

## 1. What is wired today

battleui participates in the Atomic Traceable Documentation system as a **registered ATD
project**:

- **Workspace registration** — `.atd.workspace` (repo root) lists `battleui` with `path: battleui`.
- **Project config** — `battleui/.atd` sets `docs_path: docs/`, similarity thresholds, and the
  local LLM/Ollama provider settings used by the ATD MCP tooling.
- **The Atom corpus** — `battleui/docs/` holds **63 `.atom.md`** files (e.g. `ui_battle_arena.atom.md`,
  `battleui_upsilon_api_service.atom.md`, `mech_sanctum_token_renewal`, `api_standard_envelope`).
- **Code anchors** — **200 `@spec-link [[atom]]` annotations across 60 PHP files** binding code to
  Atoms. Some are **cross-project**: `[[upsilonapi:api_shop_purchase]]`,
  `[[upsilonbattle:mec_three_slot_equipment_system]]`, `[[upsilonbattle:entity_character_equipment]]`.

The Go services already use the **same link syntax** in comments (`spec-link [[entity_character]]`,
`[[api_character_skill_inventory]]`, …), so the anchoring convention is uniform across PHP and Go —
this is what makes rewiring mechanical rather than inventive.

## 2. The three rewiring tasks

### Task 1 — Re-anchor code links (the bulk: 200 links)
Every `@spec-link [[x]]` in the deleted PHP must reappear on the **equivalent Go construct**
(handler, repo func, DTO, hub publisher). Treatment is per-atom-type, not blind copy:

- **API / endpoint atoms** (`api_standard_envelope`, `api_matchmaking`, `api_battle_proxy`,
  `api_go_webhook_callback`, …) → onto the Gin handler / middleware that now satisfies them.
- **Mechanic atoms** (`mech_game_state_versioning`, `mech_sanctum_token_renewal`,
  `mech_frontend_test_seams`) → onto the Go implementation. **Some atom *content* changes**:
  `mech_sanctum_token_renewal` describes Sanctum specifics; if auth changes (doc 02 §4) the Atom
  text must be revised, not just relinked. Same for any Reverb/Pusher-specific realtime atoms once
  the transport is chosen (doc 03).
- **Entity / rule atoms** (`entity_game_match`, `rule_leaderboard_score_calculation`,
  `rule_character_skill_slots`, …) → onto Go models/domain funcs; content is transport-agnostic,
  so these are pure relinks.
- **Cross-project links** (`[[upsilonapi:...]]`, `[[upsilonbattle:...]]`) → carry over **verbatim**;
  targets are unaffected by the rewrite. The typed Go engine client is actually a *better* anchor
  for `[[upsilonapi:*]]` links than the PHP HTTP wrapper was.

### Task 2 — Migrate the Atom corpus & project registration
- Decide the new module's directory (e.g. `upsilonhub/`) and **move/retarget `docs/`** there, or
  keep `battleui/docs/` if the dir name is retained. Update `battleui/.atd` `docs_path` accordingly.
- Update **`.atd.workspace`**: rename/repath the `battleui` project entry to the new module path
  (and keep its `name` stable if you want existing `[[battleui:*]]` inbound links from other
  projects to keep resolving — otherwise update those too).
- **Audit inbound links from other projects.** Other services may link to battleui atoms; a
  rename breaks them. Grep the whole workspace for `[[battleui:` and for the bare atom slugs
  before finalising the new name.

### Task 3 — Reconcile UI atoms with the (kept) frontend
The ~30 `ui_*.atom.md` and `module_frontend_*.atom.md` atoms describe the Vue app, which **stays**.
Their *code anchors* may currently point at Blade/Inertia/Laravel-side seams (e.g.
`mech_frontend_test_seams` maps to the `__test/*` web routes). When those routes move into Go (or
to static serving), the anchors move with them; the atom **content** mostly survives because the
Vue components don't change.

## 3. Tooling to drive it

The ATD MCP server is available in this workspace (`mcp__atd__*`). Use it rather than hand-editing:

- `atd_workspace_list` / `atd_workspace_use` — confirm/repath the project registration.
- `atd_map`, `atd_trace`, `atd_query`, `atd_search` — enumerate which atoms are anchored where in
  the PHP before deletion (build the authoritative "links to re-home" list).
- `atd_test_links` / `atd_audit` / `atd_check` / `atd_lint` — after each phase, verify no link
  dangles and code↔atom congruence holds.
- `atd_heatmap` / `atd_heatmap_code` — spot under-anchored areas in the new Go code.

> Note: `reporting/audit/atd_workspace_link_resolution_bug.md` (a now-deleted audit doc in this
> repo's history) flagged a workspace **link-resolution** issue. Re-validate cross-project link
> resolution *before* relying on `atd_test_links` as the migration gate — confirm that bug's
> status first so green results are trustworthy.

## 4. Sequencing within the migration

Fold ATD work into each code phase rather than as a big-bang at the end:

1. **Before deleting any PHP** (start of each phase): run `atd_map`/`atd_trace` to capture the
   atoms anchored in the files that phase replaces — this is the checklist.
2. **As Go code lands:** add the `[[atom]]` spec-link comments on the new constructs (cheap when
   writing, costly later — same logic as the OTel argument in doc 04).
3. **End of each phase:** `atd_test_links` + `atd_check` must be green for that slice.
4. **At cutover (Phase 6):** finalise `.atd.workspace` repath, fix inbound `[[battleui:*]]` links
   workspace-wide, retire obsolete Reverb/Sanctum-specific atom *content*, and run a full
   `atd_audit`.

## 5. Bottom line

ATD rewiring is **mechanical but non-optional and broad** (200 anchors, 63 atoms, 1 workspace
entry, plus cross-project inbound links). It is low *intellectual* risk but real *bookkeeping*
load, and it must be sequenced **per phase** — re-anchor as you port, verify with the ATD MCP
tools, and treat a few Atoms (auth, realtime transport) as content rewrites rather than pure
relinks because the underlying mechanism genuinely changes.
