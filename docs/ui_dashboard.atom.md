---
id: ui_dashboard
human_name: Dashboard Page UI
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[module_frontend]]
  - [[uc_player_login]]
dependents:
  - [[ui_dashboard_player_statistics]]
  - [[ui_dashboard_queue_selection]]
  - [[ui_dashboard_navigation]]
  - [[uc_player_login]]
  - [[uc_progression_stat_allocation]]
  - [[ui_character_roster]]
  - [[module_ui_tactical_layout]]
  - [[ui_dashboard_profile_edit]]
  - [[ui_dashboard_security_check]]
  - [[ui_dashboard_roster_display]]
  - [[ui_dashboard_match_statistics]]
---
# Dashboard Page UI

## INTENT
To aggregate the constituent rules of Dashboard Page UI.

## THE RULE / LOGIC
Serves as a primary logged-in hub where players review their roster and initiate matchmaking.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[ui_dashboard]]`
