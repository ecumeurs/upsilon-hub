# Upsilon Hub — Principal Architect Audit (Index)

**Date:** 2026-06-16
**Auditor mandate:** non-destructive architectural audit across the four systems —
**Frontend, Laravel API, Go backend, Database** — verifying documentation coverage
and that minimal business requirements are met and tested. No source code was
changed. ATD/issue documentation was repaired (see §5).

**Method:** static analysis + the `atd` CLI (the ATD MCP server was down this
session — root cause and fix in §4). Test pass/fail is taken from
`reporting/CI_report.md` (commit `d29fd1b`, **2026-05-12** — stale, predates
pending commits) and reasoned against the current working tree.

---

## 1. Reports

| System | Summary | Detailed |
|---|---|---|
| Go backend (engine + bridge + libs) | [go_backend_summary.md](go_backend_summary.md) | [go_backend_detailed.md](go_backend_detailed.md) |
| Laravel API (battleui PHP) | [laravel_api_summary.md](laravel_api_summary.md) | [laravel_api_detailed.md](laravel_api_detailed.md) |
| Frontend (battleui Vue) | [frontend_summary.md](frontend_summary.md) | [frontend_detailed.md](frontend_detailed.md) |
| Database (migrations/schema) | [database_summary.md](database_summary.md) | [database_detailed.md](database_detailed.md) |

## 2. The three cross-system themes

The individual systems are, on their own, mostly well-built. The rot is in the
**connective tissue** — and it shows up identically everywhere:

1. **Implementation is traced; proof is not.** Every system has dense
   `@spec-link` tagging and near-zero `@test-link`:
   | System | spec-link | test-link |
   |---|---:|---:|
   | upsilonbattle | 125 | 101 *(only exception)* |
   | upsilonapi | 42 | 24 |
   | battleui PHP | 179 | 11 |
   | battleui Vue | 90 | 0 |
   Tests frequently **exist** (18 PHPUnit, 6 Playwright specs) but aren't linked,
   so coverage tooling cannot see them — and where they are linked, several
   don't assert their own case (engine `missmatch` typo; `e2e_credit_economy`
   reads the wrong entity).

2. **The dashboards lie.** The BRD compliance table marks CR-05/08/10/11 ❌ while
   the matching E2E scenarios pass ✅ in the same file; Playwright "passes" with
   one test while CI reports "No HTML Report found"; issues are marked Resolved
   for work that isn't done (ISS-084) and Open for work that is (ISS-091/097).
   Green/red signals no longer track reality.

3. **Status drift in the traceability layer itself.** 575 atoms, 158 Go lint
   defects, 144 DRAFT (many implemented & STABLE), a phantom High issue
   (ISS-046), and a stale `issues/README.md`. For a project whose entire thesis is
   "absolute traceability," the traceability layer is the least trustworthy
   artifact in the repo.

## 3. Consolidated risk register

| # | Risk | Sev | System | Evidence | Status |
|---|---|---|---|---|---|
| R1 | Raw User UUID leaked as `player_id` to client | **High** | Go→FE | `output.go:208/469/474`; consumed `BattleArena.vue:39/212/244` | ISS-098 Open — **fix is cheap** (masking key already in DB, §below) |
| R2 | Admin can self-anonymize and lock out admin access | **Critical** | Laravel | `AdminController::anonymize:147-150` guards only "last admin", not self | ISS-093 Open (partial) |
| R3 | BRD compliance gate decoupled from test results | **High** | Reporting | CR-05/08/10/11 ❌ vs passing E2E in same report | new finding |
| R4 | Playwright not gating CI (1 test, no report) | **High** | FE/CI | `playwright_last_run.log` 1 passed; CI "No HTML Report found" | ISS-082 mislabeled Resolved |
| R5 | `BattleArena.vue` 846 LOC, falsely "refactored" | Med | FE | > 600 ATD error; ISS-084 Resolved | ISS-084 reopen |
| R6 | Silent-fallback bugs violate Crash-Early | Med | Go | ISS-096 trap, ISS-099 zone (`skill.go:177-187`) | both Open, confirmed |
| R7 | Unversioned resurrection state blobs | Med | DB | `game_state_cache` opaque JSON, only outer `version` | new finding |
| R8 | 158 ATD lint defects + 144 DRAFT drift | Med | All | `atd lint`; `atd stats` | partially remediated §5 |
| R9 | Issue tracker integrity (phantom/stale) | Med | Process | ISS-046 missing; README dead links | remediated §5 |
| R10 | Schema changes carry no ATD links | Low | DB | migrations untagged | new finding |

**Note on R1:** the masking identifier (`ws_channel_key`) already exists, auto-
generates, rotates on login, and routes WebSockets on the Laravel side. The
highest-severity open issue is a single engine-output seam, not a cross-stack
project — see [database_detailed.md](database_detailed.md) §3.

## 4. ATD tooling recovery (done this session)

- **Root cause of MCP outage:** `~/.claude.json` launched the server from
  `/home/bastien/.local/bin/atd` (host user) while this devcontainer runs as
  `vscode`, where the binary is `/home/vscode/.local/bin/atd`. The bad path made
  `atd serve` fail on launch.
- **Fix applied:** added repo-local `/workspace/.mcp.json` pointing at the correct
  binary. Reconnect via `/mcp` or a session restart to restore the MCP tools.
- **Index rebuilt:** the semantic index (`docs/.atd_index.db`) was ~3 weeks stale
  (built 2026-05-10; newest atom 2026-05-29). Rebuilt graph (`atd crawl`) +
  semantic index (`atd index`, local `nomic-embed-text`) so coverage numbers are
  current. The CLI works regardless of the MCP server.

## 5. Documentation remediation performed

Detail in [remediation_log.md](remediation_log.md). Scope (no source code, no
new issues):

**Done (bounded, unambiguous):**
- **`issues/README.md` fully rebuilt** from the real `ISS-*.md` files — correct
  links, all 39 issues + companion report, statuses reconciled against source
  (Open vs Resolved corrected, e.g. ISS-091/092/097 → Resolved).
- **Root `README.md`** active-issues table regenerated via `issues
  --update-readme` (now 22 active, 0 dead links).
- **ATD structural lint fixes:** promoted the stray `[temp]` / "New Atom" stub to
  a valid `mech_move_validation_jump_limitations` atom; fixed two `N/A` layer
  enums (mapmaker vision→BUSINESS, contract→ARCHITECTURE); repaired 9 broken
  intra-doc alias links (`[[upsilon_vision]]`→`[[vision_upsilon_vision]]`,
  `[[upsilon_contract]]`→`[[contract_upsilon_contract]]`).

**Deliberately deferred (judgment-heavy; documented, not done):**
- ~57 atoms missing `## EXPECTATION`/`## TECHNICAL INTERFACE` — require real
  per-atom content, not linter-satisfying stubs.
- ~36 unresolved cross-project parent links — a mix of genuinely missing atoms and
  `atd` cross-project resolution limits (`atd weave` territory); needs verification
  before creating/repointing.
- **DRAFT→STABLE status drift** (144 DRAFT, many implemented) — *not* bulk-flipped:
  reclassifying status without per-atom impl+test verification would replace one
  inaccuracy with another. Recommended as a reviewed batch.

## 6. Deferred (recommended, not done)
- Full 163-graph-orphan cleanup and the broader DRAFT backlog.
- Filing new ISS-* issues for R3, R7, R10 (new findings) — left to the team.
- Any source-code/test changes (R1–R7 remediation) — out of audit scope.
- Wiring Playwright + a truthful BRD compliance generator into CI.
