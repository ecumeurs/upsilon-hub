# Database — Summary Report

**System:** PostgreSQL via Laravel migrations — 29 migrations, 12 models,
seeders/factories, resurrection persistence
**Date:** 2026-06-16 · **Full detail:** [database_detailed.md](database_detailed.md)

## Snapshot
| Dimension | Reading |
|---|---|
| V2 schema readiness | **High** — credits, items, skills, equipment, x10 stats all migrated |
| Resurrection (ISS-054) | **Resolved & verified** — `game_state_cache`/`grid_cache`/`version` + working bridge |
| GDPR/privacy at schema | **Good** — soft deletes, private fields, FK cascade discipline |
| Masking infra (ISS-098) | **Present in DB** (`ws_channel_key`) — only the engine output bypasses it |
| Schema traceability | **None** — migrations carry no ATD links |

## Top findings
1. **ISS-098 is cheaper than it looks:** the masked `ws_channel_key` already
   exists, auto-generates, rotates on login, and routes WS — Laravel speaks masked
   keys everywhere. The leak is purely the engine emitting raw `ControllerID`.
   Fix = thread the existing key; no schema work.
2. **ISS-054 status is accurate** (rare in this repo): resurrection substrate and
   bridge are genuinely implemented, with a `version` column for conflict
   detection.
3. **Unversioned state blobs:** `game_state_cache`/`grid_cache` are opaque JSON
   with only an outer integer `version` — a resurrection-corruption risk as engine
   state evolves in V2.
4. **Traceability stops at the DB:** migrations link to no atoms — the only layer
   with zero spec/test linkage.

## Architect's commentary

**What is correct / sound.** This is the most quietly competent layer in the
audit. The V2 schema is already migrated and modelled end-to-end; resurrection has
a real, versioned persistence substrate with a matching bridge; and GDPR is
handled at the schema level (soft deletes, flagged private fields, deliberate FK
cascade/null lifecycle). Admin pagination is index-backed (ISS-053), not naive
scanning. The data layer is ahead of the code that consumes it.

**What surprised me — good & bad.** *Good (the standout of the whole audit):* the
masking identifier for ISS-098 **already exists and is fully wired on the Laravel
side** — the highest-severity open issue is a one-seam engine fix, not a
cross-stack project. *Bad:* nobody connected that dot — the issue is filed as a
broad privacy violation when the infrastructure to close it has been sitting in
the schema since April. It's a communication/traceability failure more than an
engineering one.

**What is not appropriate / not scalable.** Two things. First, **opaque,
under-versioned state blobs** (`game_state_cache`) will eventually corrupt
resurrection across engine releases unless the serializer stamps a schema version
*inside* the blob — V2's temporary-entity/effect state makes this near-certain.
Second, **the traceability mandate dies at the database**: schema changes map to
no atoms, so the project's central "everything traces to a spec" claim is
literally false for the layer of record. Both are fixable cheaply now and
expensive later.
