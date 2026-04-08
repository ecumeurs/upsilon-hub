---
id: ui_dashboard_match_statistics
status: DRAFT
human_name: Dashboard Match Statistics
priority: 5
tags: [matchmaking, dashboard, stats]
parents:
  - [[ui_dashboard]]
type: UI
layer: ARCHITECTURE
version: 1.0
dependents: []
---

# New Atom

## INTENT
To display real-time matchmaking activity statistics to the user on the dashboard.

## THE RULE / LOGIC
- Retrieve "Number of matches currently waiting for players".
- Retrieve "Number of active matches in progress".
- Implementation: Client-side polling every 60 seconds.
- Fallback: Manual Refresh button triggers immediate fetch.
- UI Location: Top or side panel of the Dashboard hub.

## TECHNICAL INTERFACE
- **API Endpoints:** `GET /v1/match/stats/waiting`, `GET /v1/match/stats/active`
- **Code Tag:** `@spec-link [[ui_dashboard_match_statistics]]`

## EXPECTATION
- Dashboard must display real-time counts for active and waiting matches.
- Counts must refresh at minimum every 60 seconds.
- A manual refresh button must be present in the UI.
