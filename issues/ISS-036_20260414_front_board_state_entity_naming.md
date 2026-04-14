# Issue: Standardize Board State Naming: entities -> characters

**ID:** `20260414_front_board_state_entity_naming`
**Ref:** `ISS-036`
**Date:** 2026-04-14
**Severity:** Medium
**Status:** Open
**Component:** `battleui/resources/js`, `upsiloncli`, `battleui/app/Http/Resources`
**Affects:** `battleui/resources/js/services/tactical.js`, `upsiloncli/upsilon_log_parser.py`

---

## Summary

The board state structure currently uses the term "entities" for game units. This is "upsilonapi slang" and should be renamed to "characters" (or "character" as an object key) for better alignment with the game's domain language and better readability on the front side (CLI and Vue.js).

---

## Technical Description

### Background
The Upsilon Battle Engine refers to units as "entities". This terminology has leaked into the API response and subsequently into the CLI parser and the Vue.js frontend components.

### The Problem Scenario
Developers working on the frontend need to remember that "entities" means "characters". This adds cognitive load and is inconsistent with many UI labels that already use "Character".

### Where This Pattern Exists Today
- **Laravel Resource:** `BoardStateResource.php` (Line 55: `$player['entities']`)
- **Vue.js Service:** `tactical.js` (Line 46: `p.entities.find(...)`, Line 59: `me.entities`)
- **CLI Parser:** `upsilon_log_parser.py` (Line 126: `p.get('entities', [])`)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium — Technical debt and confusion. |
| Detectability | High — Clear naming mismatch. |
| Current mitigant | None, developers just accept the slang. |

---

## Recommended Fix

**Short term:** Update `BoardStateResource.php` to rename the key in the JSON response, then immediately update the CLI and Vue.js consumption sites.  
**Medium term:** Update the CLI `upsilon_log_parser.py` to use "characters" in its internal dictionary structure for consistency in diagnostic reports.  
**Long term:** Ensure all new frontend components use the `character` terminology exclusively.

---

## References

- [BoardStateResource.php](file:///workspace/battleui/app/Http/Resources/BoardStateResource.php)
- [tactical.js](file:///workspace/battleui/resources/js/services/tactical.js)
- [upsilon_log_parser.py](file:///workspace/upsiloncli/upsilon_log_parser.py)
