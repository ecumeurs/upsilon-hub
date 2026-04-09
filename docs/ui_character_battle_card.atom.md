---
id: ui_character_battle_card
status: STABLE
human_name: Character Battle Card UI
type: UI
layer: ARCHITECTURE
version: 1.0
priority: 5
tags: [ui, combat, character, card]
parents:
  - [[ui_team_roster_panel]]
dependents: []
---

# New Atom

## INTENT
A single character display card within a team roster panel, showing name, HP bar, movement, and stats with support for compact mode and configurable accent color.

## THE RULE / LOGIC
- **Full Mode:** Character name, HP bar (current/max with neon fill), Movement bar (current/max), Attack stat, Defense stat.
- **Compact Mode:** Character name + HP bar only.
- **Color Prop:** Accepts an accent color for team-based styling (blue, green, red, purple).
- **HP Bar:** Filled bar with neon glow effect proportional to current/max HP.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_character_battle_card]]`
- **Component:** `CharacterBattleCard.vue`
- **Props:** `character`, `compact` (bool), `accentColor` (string)

## EXPECTATION
- Full mode shows all stats.
- Compact mode shows only name + HP bar.
- Accent color applies to borders and HP bar fill.
