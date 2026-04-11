# Issue: Dashboard Hub Implementation

**ID:** `20260408_dashboard_hub_implementation`
**Ref:** `ISS-025`
**Date:** 2026-04-08
**Severity:** Medium
**Status:** Resolved
**Component:** `battleui/resources/js/Pages/Dashboard.vue`
**Affects:** `battleui`

---

## Summary

This issue tracks the implementation of the Dashboard Hub for the Upsilon Battle application. The dashboard serves as the central point for player actions, including match triggering, character roster management, and profile editing.

---

## Technical Description

### Background
The current dashboard documentation exists in `ui_dashboard.atom.md` and related atoms, but the full scope of features (win/loss ratio, active/waiting match counts, profile editing, and integrated character management) needs formalization and implementation.

### The Problem Scenario
A player needs a central UI to:
- Start PVP/PVE matches (1 or 2 players).
- Review character stats, available rerolls, and apply progression.
- View real-time match statistics (waiting/active matches).
- Manage their profile (address, birth date, GDPR requests, account deletion).

### Where This Pattern Exists Today
- `docs/ui_dashboard.atom.md` (Base module)
- `docs/ui_dashboard_player_statistics.atom.md` (Win/Loss)
- `docs/ui_dashboard_roster_display.atom.md` (Character display)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | Medium |
| Current mitigant | Existing siloed ATD atoms |

---

## Recommended Fix

**Short term:** Create a static mockup and define all necessary ATD atoms (match stats, profile edit).
**Medium term:** Implement the UI components in Vue.js and link them to the backend endpoints.
**Long term:** Finalize real-time updates for match statistics using WebSockets/Reverb.

---

## References

- [ui_dashboard.atom.md](file:///workspace/docs/ui_dashboard.atom.md)
- [ui_dashboard_player_statistics.atom.md](file:///workspace/docs/ui_dashboard_player_statistics.atom.md)
- [rule_gdpr_compliance.atom.md](file:///workspace/docs/rule_gdpr_compliance.atom.md)
- [us_character_reroll.atom.md](file:///workspace/docs/us_character_reroll.atom.md)
