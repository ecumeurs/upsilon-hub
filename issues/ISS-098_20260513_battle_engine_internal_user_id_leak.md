# Issue: Internal User ID Exposure in Battle State DTOs

**ID:** `Ref_20260513_battle_engine_internal_user_id_leak`
**Ref:** `ISS-098`
**Date:** 2026-05-13
**Severity:** Low
**Status:** Open-Low (hardening backlog)
**Component:** `upsilonapi/api`
**Affects:** `upsilonapi` internal bridge only; web path is masked at battleui.

---

## Summary

The Go battle engine's API bridge currently leaks internal User UUIDs via the `player_id` field in tactical `Entity` DTOs. This violates the `[[requirement_customer_user_id_privacy]]` mandate which prohibits exposing database primary keys to the client to prevent primary key enumeration and protect user privacy.

---

## Technical Description

### Background
According to `[[requirement_customer_user_id_privacy]]`, all internal User identifiers must be masked. The frontend should only receive masked identifiers (like nicknames or session-specific keys) rather than raw database UUIDs.

### The Problem Scenario
When the `upsilonapi` bridge prepares the `BoardState` DTO to send to Laravel (and subsequently to the frontend via WebSockets), it maps the engine's internal `ControllerID` directly to the JSON `player_id` field.

1. Engine executes a turn or event.
2. `upsilonapi/api/output.go` constructs a new `BoardState`.
3. `NewEntity` function is called for each tactical unit.
4. The internal `ControllerID` (a raw User UUID) is assigned to `PlayerID`.

```go
// upsilonapi/api/output.go:208
res := &Entity{
    ID:             ent.ID.String(),
    PlayerID:       ent.ControllerID.String(), // LEAK: Internal User UUID
    Team:           int(ent.Team),
    // ...
}
```

### Where This Pattern Exists Today
- [upsilonapi/api/output.go:208](file:///workspace/upsilonapi/api/output.go#L208)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High — revealed in every board state WebSocket message or API poll |
| Current mitigant | None. The DTO field is populated unconditionally. |

---

## Recommended Fix

**Short term:** Implement an `IdMasker` utility in `upsilonapi` that maps internal `ControllerID`s to a transient, match-specific `TacticalPlayerID` or simply use the Player's Team/Index if unique.

**Medium term:** Update `NewEntity` to use the masked ID instead of `ent.ControllerID.String()`. Ensure the `BoardStateResource` in Laravel also respects this masking.

**Long term:** Align with `[[arch_api_id_masking_gateway]]` to ensure all cross-service boundaries perform automatic ID transformation.

---

## References

- [requirement_customer_user_id_privacy.atom.md](file:///workspace/docs/requirement_customer_user_id_privacy.atom.md)
- [upsilonapi/api/output.go](file:///workspace/upsilonapi/api/output.go)
- [e2e_battle_starts_privacy_check.js](file:///workspace/upsiloncli/tests/scenarios/e2e_battle_starts_privacy_check.js)

---

## Resolution / Re-scoping (2026-06-16)

**Decision:** downgrade to **Low / hardening backlog**. The issue is not closed because the `docker-compose.prod.yaml` host-port binding remains; it is reclassified because the actual public-facing risk is eliminated by the masking gateway already in place.

### Why the web path is already safe

The `player_id` UUID is present on the internal `upsilonapi → battleui` hop (HTTP, inside the compose network) but is **stripped before it reaches any external client**:

- `battleui/app/Http/Resources/BoardStateResource.php` calls `unset($array['player_id'])` and `unset($array['current_player_id'])` during serialisation.
- `battleui/app/Events/BoardUpdated.php` broadcasts per-recipient board state via that Resource; the raw UUID never appears on any WebSocket channel.
- `battleui/app/Http/Controllers/GameController.php` returns HTTP polling responses through the same Resource.

`upsilonapi` is therefore a **trusted internal bridge**, not a public API. Raw IDs on the internal hop are accepted by design (documented in `[[arch_api_id_masking_gateway]]`).

### Why 8081 is not publicly reachable (in production)

- The AWS EC2 security group (`upsilonaws/scripts/setup/01_networking.sh`) authorises **only ports 22, 80, 443** from `0.0.0.0/0`. Port 8081 has no SG ingress rule.
- The nginx reverse proxy routes only `APP_PORT` (8000) and `WS_PORT` (8080); 8081 is not proxied.
- `docker-compose.ci.yaml` has **no `ports:` mapping** for the `engine` service — 8081 is compose-internal in CI.

### Remaining latent risk (Low)

`docker-compose.prod.yaml` line 101 maps `"8081:8081"` to the EC2 host. While the current SG blocks public access, this binding means that any future accidental SG rule addition would immediately expose the raw-UUID API. Recommended hardening: **remove the host-port binding** from `docker-compose.prod.yaml` and rely solely on Docker's internal network (`engine` service reachable as `http://engine:8081` from `app`). No code change required on the Go side.
