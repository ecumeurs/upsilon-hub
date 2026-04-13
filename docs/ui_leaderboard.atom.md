---
id: ui_leaderboard
human_name: Leaderboard Page UI
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[api_leaderboard]]
  - [[ui_leaderboard_data_display]]
  - [[ui_leaderboard_metrics_displayed]]
  - [[ui_leaderboard_modes]]
  - [[ui_leaderboard_primary_sorting]]
  - [[ui_leaderboard_secondary_sorting]]
  - [[ui_leaderboard_security]]
---
# Leaderboard Page UI

## INTENT
To aggregate the constituent rules of Leaderboard Page UI.

## THE RULE / LOGIC
- Displays a leaderboard to showcase player rankings and statistics.
- **Integration:** The leaderboard component must be integrated into the main Dashboard view, positioned directly below the Match Type Selector.
- Components allow switching between battle modes via tabs/buttons.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[ui_leaderboard]]`
