---
id: ui_dashboard
human_name: Dashboard Page UI
type: UI
version: 1.0
status: REVIEW
priority: CORE
tags: [ui, secure, dashboard]
parents:
  - [[req_security]]
dependents:
  - [[req_matchmaking]]
  - [[ui_leaderboard]]
---

# Dashboard Page UI

## INTENT
To serve as the primary logged-in hub where players review their roster and initiate matchmaking.

## THE RULE / LOGIC
- Roster Display: Must present the stats (HP, Movement, Attack, Defense) of the player's 3 characters.
- Player Statistics: Must permanently display the player's total match Wins, total match Losses, and their calculated Win/Loss ratio.
- Navigation: Must provide a clear entry point to the Global Leaderboard screen.
- Queue Selection: Must present 4 distinct buttons to start a new game:
  1. 1v1 PVE
  2. 1v1 PVP
  3. 2V2 PVE
  4. 2V2 PVP
- Security: Requires a valid JWT to access.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[ui_dashboard]]`
- **Test Names:** `TestDashboardDisplaysRoster`, `TestDashboardQueueButtons`

## EXPECTATION (For Testing)
- Logged-in user visits dashboard -> Sees their 3 characters, stat counts (Wins/Losses), a link to Leaderboard, and the 4 specific queue buttons.
