---
id: ui_leaderboard_data_display
human_name: Data Display Rules
type: UI
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[ui_leaderboard]]
dependents: []
---
# Data Display Rules

## INTENT
Must present a sorted, paginated list of players globally.

## THE RULE / LOGIC
- Must present a sorted list of players globally.
- **Pagination:** Exactly 10 entries per page.
- **Current User Context:** Always display the authenticated user's position, statistics, and rank (typically pinned at the bottom or top of the view), even if they are not present on the current page of results.
- **Empty State:** If zero data is found, display a themed message: "SENSORS OFFLINE: NO DATA RECOVERED" or "AREA SCAVENGED: NO SIGNS OF LIFE".
- **Search Logic:** Searching for a user must ALWAYS provide feedback. If no matches exist (e.g. user has 0 matches in mode), display: "COMMUNICATIONS JAMMED: NO SIGNATURE FOUND" to confirm the search attempt was completed.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[ui_leaderboard_data_display]]`
