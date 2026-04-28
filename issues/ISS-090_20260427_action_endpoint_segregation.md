# Issue: Action Endpoint Segregation

**ID:** `20260427_action_endpoint_segregation`
**Ref:** `ISS-090`
**Date:** 2026-04-27
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi/api`, `upsilonapi/bridge`
**Affects:** `battleui`, `upsiloncli`

---

## Summary

Currently, all tactical actions (move, attack, skill, pass) are funneled through a single generic `/game/{id}/action` endpoint with a `type` field. This leads to a bloated and less type-safe DTO (`ArenaActionRequest`) where fields like `skill_id` or `target_coords` are conditionally mandatory.

---

## Technical Description

### Background
The system uses a single DTO for all actions:
```go
type ArenaActionRequest struct {
	PlayerID     string     `json:"player_id"`
	Type         string     `json:"type"`
	TargetCoords []Position `json:"target_coords"`
	EntityID     string     `json:"entity_id"`
	SkillID      string     `json:"skill_id,omitempty"`
}
```

### The Problem Scenario
As more complex actions are added (e.g., skill usage, item usage, complex movement), this single DTO will continue to grow with optional fields, making validation more difficult and the API less discoverable.

### Where This Pattern Exists Today
- `upsilonapi/api/input.go`
- `upsilonapi/bridge/bridge.go`
- `communication.md`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | Internal switch statements in the bridge |

---

## Recommended Fix

**Short term:** Maintain the single endpoint but improve validation in the bridge.  
**Medium term:** Introduce dedicated endpoints for each action type:
- `POST /game/{id}/move`
- `POST /game/{id}/attack`
- `POST /game/{id}/skill`
- `POST /game/{id}/pass`
**Long term:** Auto-generate strongly typed client libraries from these specific endpoints.

---

## References

- [communication.md](file:///workspace/communication.md)
- [bridge.go](file:///workspace/upsilonapi/bridge/bridge.go)
- [input.go](file:///workspace/upsilonapi/api/input.go)
