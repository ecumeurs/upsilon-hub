# Issue: Integer Overflow in Match Turn Tracking

**ID:** `20260415_match_turn_overflow`
**Ref:** `ISS-041`
**Date:** 2026-04-15
**Severity:** High
**Status:** Resolved
**Component:** `battleui/database/migrations`, `battleui/app/Http/Controllers/API/WebhookController.php`
**Affects:** Upsilon Battle Engine webhooks, CLI bot farm

---

## Summary

The `turn` column in the `game_matches` table is currently an `integer` (32-bit). However, the game engine's versioning strategy bit-packs turn indices into the high 32 bits of a 64-bit integer. This causes an overflow (value `4294967296` for Turn 1) when the Laravel API attempts to save the version into the `turn` column during a webhook callback.

---

## Technical Description

### Background
The [[mech_version_bit_packing]] specification dictates that the game version is an `int64` where:
`Version = (int64(TurnIndex) << 32) | int64(ActionIndex)`

### The Problem Scenario
When the game starts (Turn 1, Action 0), the engine sends `version: 4294967296`.
The Laravel `WebhookController` receives this and attempts:
```php
$match->update([
    'version' => 4294967296,
    'turn' => 4294967296,
]);
```
While the `version` column was updated to `bigInteger`, the `turn` column remains `integer`. This triggers:
`SQLSTATE[22003]: Numeric value out of range: 7 ERROR: value "4294967296" is out of range for type integer`

### Where This Pattern Exists Today
- `battleui/database/migrations/2026_03_12_081704_create_game_matches_table.php:18`
- `battleui/app/Http/Controllers/API/WebhookController.php:58`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (guaranteed on turn 1) |
| Impact if triggered | High (breaks all game state persistence) |
| Detectability | High (500 errors in logs, engine warnings) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Update the `turn` column in `game_matches` to `bigInteger` via migration.
**Medium term:** Review all database columns used for bit-packed values to ensure they are `bigInteger`.
**Long term:** Consider a more robust DTO layer that validates bit-packed ranges before database insertion.

---

## References

- [2026_03_12_081704_create_game_matches_table.php](file:///workspace/battleui/database/migrations/2026_03_12_081704_create_game_matches_table.php)
- [WebhookController.php](file:///workspace/battleui/app/Http/Controllers/API/WebhookController.php)
- [[mech_version_bit_packing]]
