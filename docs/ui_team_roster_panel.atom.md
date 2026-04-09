---
id: ui_team_roster_panel
status: STABLE
human_name: Team Roster Panel UI
type: UI
layer: ARCHITECTURE
version: 1.0
tags: [ui, combat, roster, team]
dependents:
  - [[ui_character_battle_card]]
priority: 5
parents:
  - [[ui_battle_arena]]
---

# New Atom

## INTENT
A side panel displaying all characters for a team grouped by player nickname, supporting both detailed and compact display modes.

## THE RULE / LOGIC
- **Grouping:** Characters are grouped under their owning player's nickname as a section header.
- **Current Player (Left, Primary):** Full detailed display for own characters (name, HP, movement, attack, defense).
- **2v2 Ally (Left, Secondary):** Compact display (name, HP only) with a distinct accent color (green).
- **Adversaries (Right):** All enemy characters in compact mode, grouped by player, with enemy accent colors (red/purple).
- **Props:** Accepts `players[]`, `isDetailedForPlayer` (player_id to show detailed), `side` ('left'|'right').

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_team_roster_panel]]`
- **Component:** `TeamRosterPanel.vue`

## EXPECTATION
- Own characters display full stats.
- Ally and enemy characters display in compact mode.
- Player nicknames appear as section headers.
