# Issue: Unconstrained Character Upgrades (Progression Bypass)

**ID:** `20260312_character_upgrade_constraints`
**Ref:** `ISS-016`
**Date:** 2026-03-12
**Severity:** High
**Status:** Resolved
**Component:** `app/Http/Controllers/API/ProfileController.php`
**Affects:** `Character` model, Frontend profile view

---

## Summary

The `/api/v1/profile/{id}/character/{characterId}/upgrade` endpoint allows arbitrary stat increments without checking win counts or attribute-specific progression rules. This violates the core progression system defined in the ATD.

---

## Technical Description

### Background
According to `rule_progression`, characters should only gain 1 attribute point per win. Movement upgrades are further restricted to once every 5 wins.

### The Problem Scenario
The current implementation in `ProfileController::updateCharacter` directly applies the requested increments from the `stats` array payload with zero validation against the user's game history or win count.

```php
// Current flawed logic
foreach ($validated['stats'] as $stat => $value) {
    $character->$stat += $value; // No check if $value > 1 or if win available
}
```

### Where This Pattern Exists Today
- `app/Http/Controllers/API/ProfileController.php:120-144`

---

## Risk Assessment

| Factor              | Value                                                  |
| ------------------- | ------------------------------------------------------ |
| Likelihood          | High                                                   |
| Impact if triggered | High (Game balance destruction)                        |
| Detectability       | Medium — manifests as over-powered characters in arena |
| Current mitigant    | None                                                   |

---

## Recommended Fix

**Short term:** Add a check to `updateCharacter` to ensure the character has "unallocated points" (would require adding a `points_available` column to `characters`).

**Medium term:** Implement a Win Tracker that awards tokens/points on match resolution, which are then consumed by the upgrade endpoint.

**Long term:** Move progression logic to a dedicated `ProgressionService` that validates the state against `rule_progression` ATD.

---

## References

- [ProfileController.php](file:///workspace/battleui/app/Http/Controllers/API/ProfileController.php)
- [[rule_progression]](file:///workspace/docs/rule_progression.atom.md)
- [[uc_match_resolution_movement_upgrade]](file:///workspace/docs/uc_match_resolution_movement_upgrade.atom.md)
