# Documentation Remediation Log

**Date:** 2026-06-16 · **Scope guardrails:** no source-code changes, no new issue
files. Only documentation/tracking artifacts and ATD metadata were touched.

---

## 1. Issue tracking

| Action | Detail |
|---|---|
| Rebuilt `issues/README.md` | The hand-maintained directory index had dead `Ref_*` links, omitted ~20 real issues, and carried stale statuses. Replaced with a verified full index (Open / Resolved / companion report) sourced from the actual `ISS-*.md` files. |
| Regenerated root `README.md` table | Ran `issues --update-readme` (the canonical generator, which manages the **root** README, not the directory index). Result: 22 active issues, correct filenames, 0 dead links. |
| Status reconciliation | Corrected statuses to match file bodies: ISS-091, ISS-092, ISS-097, ISS-088 are **Resolved** (previously implied Open); the V2 wave ISS-065/066/067/069/070/071/073/074/082/084/085/086 confirmed Resolved. |
| Surfaced, not fixed | **ISS-046** (turner hands turn to dead entity, High) is referenced by the old index but its file does **not exist** — tracking is lost. Not recreated (new-issue embargo). Root README's hand-written roadmap also links phantom ISS-068/075/076. |

## 2. ATD structural lint

Baseline (post index rebuild): shared **85**, upsilonbattle 102, upsilonapi 36,
upsilontypes 20, battleui 47.

| Fix | Files | Result |
|---|---|---|
| Promote stray `[temp]` stub | `docs/mech_move_validation_jump_limitations.atom.md` (was `id: temp`, `# New Atom`, empty sections, but real rule content) → valid atom with proper frontmatter + sections derived from its own body | `[temp]` error class eliminated |
| Fix `N/A` layer enums | `docs/vision_mapmaker_vision.atom.md` → `BUSINESS`; `docs/contract_mapmaker_contract.atom.md` → `BUSINESS` | `Invalid layer enum` cleared on both |
| Normalize CONTRACT/VISION invariant | `docs/contract_mapmaker_contract`, `docs/contract_upsilon_contract` | CONTRACT and VISION atoms are BUSINESS-layer roots with **no dependents** (and no parents) — stripped dependents and the (incorrectly added) parent; fixed stray `# New Atom` headings. Submodule-level contracts (`contract_api/tools/mapdata/mapmaker`) carry the same violation and are flagged for a follow-up pass. |
| Repair broken alias links | 9 edits across `contract_upsilon_contract`, `vision_upsilon_vision`, `mechanic_{math_core_utils,randomization_helpers,message_queue,message_queue_management,spatial_distance_calculations}` | `[[upsilon_vision]]`→`[[vision_upsilon_vision]]`, `[[upsilon_contract]]`→`[[contract_upsilon_contract]]`; "Unresolved parent link: [[upsilon_vision]]" eliminated |

Shared lint after fixes: **83** (the residual are content/cross-project defects
below). Verified: source code referenced neither alias form (0 hits), so these
edits carry no code impact.

## 3. Deferred (documented, intentionally not done)

| Item | Count | Why deferred |
|---|---:|---|
| Missing `## EXPECTATION` / `## TECHNICAL INTERFACE` | ~57 (shared) | Each needs real authored content; linter-satisfying stubs would be dishonest coverage. |
| Unresolved cross-project parent links | ~36 (shared) | Targets like `[[battleui:req_player_experience]]` (10×), `[[upsilonapi:domain_credit_economy]]` (4×) are a mix of genuinely missing atoms and `atd` cross-project resolution limits. Needs `atd weave` + per-link verification before creating or repointing. |
| Traceability gaps (`@spec-link`, 0 `@test-link`) | 5 (shared) | Require writing/linking real tests — a code/test change, out of audit scope. |
| DRAFT → STABLE status drift | 144 DRAFT total | Not bulk-flipped. Reclassifying without per-atom impl+test verification would swap one inaccuracy for another. Recommend a reviewed batch using `atd check` to confirm impl+test before promotion. |
| Graph orphans (no parent/dependent) | 163 | Large backlog; structural weaving exercise (`atd weave`) for a dedicated pass. |

## 4. How to verify
```
# from /workspace
atd lint                       # shared structural defects (now 83, was 85)
atd stats                      # coverage / DRAFT / orphan counts
issues --status open           # active issues (reads files directly)
ls issues/ISS-*.md             # every README link resolves to a real file
```
Restoring the MCP server (after the `/workspace/.mcp.json` fix) and reconnecting
exposes the same data via the `mcp__atd__*` tools.
