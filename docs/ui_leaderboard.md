---
id: ui_leaderboard
human_name: Leaderboard Page UI
type: UI
version: 1.0
status: REVIEW
priority: SECONDARY
tags: [ui, secure, rankings]
parents:
  - [[ui_dashboard]]
  - [[req_security]]
dependents: []
---

# Leaderboard Page UI

## INTENT
To foster competitive engagement by displaying the global rankings of players based on their match statistics.

## THE RULE / LOGIC
- Data Display: Must present a sorted, paginated list of players globally.
- Primary Sorting: The default sorting is based on total number of Wins (descending). 
- Secondary Sorting: If Wins are tied, the secondary sorting prioritizes the highest Win/Loss ratio.
- Metrics Displayed: Each row must display the Player Account Name, Total Wins, Total Losses, and Win/Loss Ratio.
- Security: Requires a valid JWT to access.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[ui_leaderboard]]`
- **Test Names:** `TestLeaderboardSorting`, `TestLeaderboardPagination`

## EXPECTATION (For Testing)
- Logged-in user visits leaderboard -> Sees list of players perfectly sorted by Wins then Ratio.
- Unauthenticated user attempts access -> Redirected per `req_security`.
