# ATD Tooling Bug Report — Workspace-Qualified Link Resolution

**For:** ATD tooling team
**From:** Upsilon Hub architecture audit (2026-06-16)
**Component:** `atd` CLI link resolver (used by `lint`, `crawl`, and by extension
`check`/`trace`/`weave`)
**Affected build:** `atd revision a1dd918200966d06e55849ac4454f94f66f47bf7`
**Severity:** High — corrupts traceability signal and induces destructive
"fix" behavior in downstream users.

---

## 1. Summary

In a multi-project ATD workspace, the resolver **fails to resolve
workspace-qualified atom links of the form `[[project:atom_id]]`**, reporting them
as `Unresolved parent link` / `Unresolved dependent link` even when the target
atom unambiguously exists. **Bare links (`[[atom_id]]`) resolve correctly.** The
qualifier is evidently treated as part of a literal atom id rather than being
split and resolved against the project map in `.atd.workspace`.

The workspace map is present and complete (all 11 projects, name → path), so the
resolver has the data it needs; it simply does not apply it to qualified links.

## 2. Impact

- **159 of the workspace's "unresolved link" lint errors are false positives**
  from this one bug (counted across all 9 active projects). **167 qualified
  `[[project:atom]]` references** exist in the corpus; essentially all are
  mis-reported.
- **It drives destructive remediation.** Because qualified cross-project parents
  lint as "unresolved," authors delete them to obtain a clean lint. That leaves
  the atom with `parents: []`.
- **Which then trips the structural pre-commit hook.** The hook
  (`scripts/hooks/pre-commit`, "ATD Structural Integrity Check") requires every
  ARCHITECTURE/IMPLEMENTATION atom to declare ≥1 parent. The de-parented atoms are
  now rejected as orphans, blocking commits.
- **Forcing escape-hatch debt.** Authors then add `parents: [[req_tech_debt_backlog]]`
  purely to pass the hook — destroying the real lineage the qualified link encoded.

Net effect: a resolver false-negative actively **erases cross-project
traceability** and replaces it with tech-debt placeholders. For a system whose
entire value proposition is "absolute traceability," this is corrosive.

## 3. Reproduction

Minimal, using this workspace:

```bash
# The target atom exists in the 'shared' project (path '.'):
ls docs/req_tech_debt_backlog.atom.md
ls docs/requirement_req_trpg_game_definition.atom.md

# An atom whose parents use the qualified form:
#   parents:
#     - [[shared:requirement_req_trpg_game_definition]]
atd lint            # -> "Unresolved parent link: [[shared:requirement_req_trpg_game_definition]]"

# Change the same reference to the bare form:
#   parents:
#     - [[requirement_req_trpg_game_definition]]
atd lint            # -> resolves cleanly, no error
```

### Observed evidence (this build)
```
Unresolved parent link: [[shared:req_tech_debt_backlog]]                 (file exists: docs/req_tech_debt_backlog.atom.md)
Unresolved parent link: [[shared:requirement_req_trpg_game_definition]]  (file exists)
Unresolved parent link: [[upsilonbattle:entity_player]]                  (file exists: upsilonbattle/docs/entity_player.atom.md)
```
Bare `[[req_tech_debt_backlog]]`: **0** unresolved errors.

**Decisive detail:** the `shared:` prefix points at the *same* project being
linted (`shared` → `.`). It still fails. So this is not merely a cross-project
lookup miss — the resolver is not stripping/normalizing the `project:` qualifier
at all before lookup.

## 4. Root-cause hypothesis

The link parser stores the inner text of `[[ ... ]]` as the lookup key and matches
it against an index keyed by bare `id`. For `[[project:atom]]` the key becomes the
literal string `"project:atom"`, which never matches the indexed `id` (`"atom"`).
The fix is to split on the first `:` into `(project, id)`, resolve `project`
through `.atd.workspace` (treating an absent/`shared`/self qualifier as the current
project), and then match `id` within that project's atom set.

## 5. Recommended fix

1. **Normalize qualified links before lookup.** Parse `[[project:id]]` →
   `{project, id}`; resolve `project` via the `.atd.workspace` project map; fall
   back to current project when the qualifier is omitted or names the active
   project (`shared` ↔ `.`).
2. **Resolve against the target project's index**, not the current project's only.
   `crawl`/`lint`/`check`/`trace` should share one workspace-aware resolver.
3. **Make `weave` workspace-aware** so bidirectional parent/dependent linking works
   across project boundaries (today a child's qualified parent and the parent's
   dependent list cannot be reconciled across repos).
4. **Round-trip test:** add fixtures asserting that `[[p:a]]`, `[[a]]` (same
   project), and a genuine dangling `[[p:missing]]` resolve / fail as expected.
5. **Consider a distinct diagnostic** for "qualifier names an unknown project" vs
   "atom not found in resolved project" — today both collapse to the same
   misleading "Unresolved … link."

## 6. Note on the pre-commit hook (not the bug)

For completeness: the structural pre-commit hook is **not** the defective
component. It only *counts* `- [[…]]` lines under `parents:` and accepts qualified
links fine; it never resolves them. Patching the hook is not the right fix — the
resolver in the `atd` binary is. The hook merely amplifies the resolver bug by
rejecting atoms after their qualified parents have been deleted.

## 7. Current workaround in this repo (to be reverted once fixed)

De-parented atoms were unblocked with the sanctioned `[[req_tech_debt_backlog]]`
(bare) escape hatch so commits could proceed. These are placeholders, not real
lineage; once the resolver honors `[[project:atom]]`, the original cross-project
parents should be restored and the escape hatches removed. The qualified links
that were *not* yet deleted (167 references) can serve as a regression corpus.

---

### Appendix — qualified-reference distribution (this workspace)
| Prefix | Qualified refs |
|---|---:|
| `shared:` | 75 |
| `upsilonapi:` | 36 |
| `upsilonbattle:` | 24 |
| `battleui:` | 17 |
| `upsilontypes:` | 13 |
| `upsiloncli:` | 2 |
| **Total** | **167** |

Unresolved-link lint errors by project (almost all attributable to this bug):
shared 38, upsilonapi 42, upsilonbattle 32, battleui 28, upsilontypes 9,
upsiloncli 3, upsilonmapmaker 3, upsilonmapdata 2, upsilontools 2 — **159 total.**
