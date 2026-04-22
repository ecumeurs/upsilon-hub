---
id: ui_battle_arena
status: STABLE
type: UI
layer: ARCHITECTURE
version: 1.0
tags: [ui, combat, arena, battle]
parents:
  - [[req_trpg_game_definition]]
  - [[req_ui_look_and_feel]]
  - [[ui_board]]
human_name: Battle Arena Page UI
priority: 5
dependents:
  - [[ui_action_panel]]
  - [[ui_combat_header]]
  - [[ui_initiative_timeline]]
  - [[ui_iso_board]]
  - [[ui_team_roster_panel]]
---

# New Atom

## INTENT
The top-level Battle Arena page layout orchestrating the combat header, team rosters, isometric board, action panel, and initiative timeline into a cohesive tactical combat interface.

## THE RULE / LOGIC
- **Layout:** Three-column layout — Left Roster | Center (Header + Board + Actions) | Right Roster.
- **Left Panel:** Current player's team roster (detailed stats for own characters, compact for 2v2 ally).
- **Right Panel:** Adversary team roster (compact display, grouped by player).
- **Center Top:** CombatHeader with fighting-game HP bars, match timer, shot clock.
- **Center Middle:** Isometric board grid with character pawns.
- **Center Bottom:** Action panel (Move, Attack, Pass, Forfeit) and Initiative Timeline.
- **Security:** Requires valid JWT via TacticalLayout.
- **Aesthetic:** Must strictly follow the "Neon in the Dust" aesthetic from `[[req_ui_look_and_feel]]`.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_battle_arena]]`
- **Page:** `BattleArena.vue`
- **Route:** `/battle-arena?match_id=:id`

## EXPECTATION
- Arena renders characters and grid upon successful data fetch.
- If the Game Engine connection fails (503 Service Unavailable), a "TACTICAL LINK FAILURE" modal is displayed.
- Action panel updates in real-time based on selected character stats.
- Ghosting (renders with default 10x10) is FORBIDDEN; the UI must wait for valid engine data.
