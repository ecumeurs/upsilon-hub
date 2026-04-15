# Issue: Friendly Fire Rule Enforcement Missing

**ID:** `Ref_20260415_rule_friendly_fire_not_enforced`
**Ref:** `ISS-043`
**Date:** 2026-04-15
**Severity:** High
**Status:** Open
**Component:** `docs/rule_friendly_fire.atom.md`
**Affects:** Battle engine combat logic.

---

## Summary

The `Friendly Immunity Rule` (`rule_friendly_fire`) defined in the architecture layer is not currently implemented or enforced in the battle engine. Combat actions do not check if the target is on the same team as the attacker, allowing players to damage their own allies. This issue is blocked by or closely related to the missing team management system ([ISS-003](ISS-003_20260305_upsilonbattle_missing_teams.md)).

---

## Technical Description

### Background

The `rule_friendly_fire` atom specifies that self-inflicted harm between players and characters on the same side should be prohibited. This is a core mechanic for team-based combat in Upsilon Battle.

### The Problem Scenario

1. Player A and Player B are (intended to be) on the "Red Team".
2. Player A initiates an attack command targeting Player B.
3. The battle logic processes the attack without checking team affiliation (or team affiliation is undefined).
4. Player B receives damage from Player A.

```
Attacker (Team Red) -> Attack Action -> Defender (Team Red)
                                      |
                                      V
                             Damage Applied (FAIL)
```

### Where This Pattern Exists Today

- `docs/rule_friendly_fire.atom.md` exists as STABLE but lacks `@spec-link` references in the codebase.
- Combat logic (likely in `GameController` or a core battle engine service) does not yet incorporate team-based filtering for damage application.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High — Players will notice ally damage immediately. |
| Current mitigant | None. |

---

## Recommended Fix

**Short term:** Define a basic team property on characters and add a guard clause in the combat resolution logic to check if `attacker.team == defender.team` before applying damage.

**Medium term:** Resolve [ISS-003](ISS-003_20260305_upsilonbattle_missing_teams.md) to provide a robust team management system.

**Long term:** Implement the constituent rules `rule_friendly_fire_match_type` and `rule_friendly_fire_team_validation` to handle varying friendly fire rules per match type.

---

## References

- [rule_friendly_fire.atom.md](../docs/rule_friendly_fire.atom.md)
- [ISS-003_20260305_upsilonbattle_missing_teams.md](ISS-003_20260305_upsilonbattle_missing_teams.md)
