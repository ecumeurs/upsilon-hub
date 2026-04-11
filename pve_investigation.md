# Investigation Wrap-up: 1v1 PvE Battle Failures

## Summary
The investigation into the `1v1_PVE` test hang-ups and failures identified two major technical hurdles and one remaining logic bug in the bot tactics. Bots are now capable of pathfinding and moving correctly, but further refinement is needed for PvE-specific enemy targeting.

---

## 1. Resolved: Grid Layout Confirmation
**Standard:** The system uses **column-major [x, y]** indexing.
- **Verification:** Confirmed via `upsilonapi/api/output.go` where `NewBoardState` explicitly maps `Cells[x][y]`.
- **Status:** `internal/script/pathfinder.go` has been aligned to this standard. Previous "inverted" successes likely resulted from symmetrical grids or consistent transposition across the bot/engine pair.

## 2. Resolved: JS Bridge Case Sensitivity
**Issue:** Go's `goja` VM was exposing exported struct fields to JavaScript using their exact Go names (e.g., `X` and `Y`). However:
- The `pvp_bot_battle.js` script expects lowercase `x` and `y`.
- WebSocket events are parsed into raw `map[string]interface{}` (lowercase keys).
**Impact:** Bots were receiving `undefined` for coordinates returned by `planTravelToward`, leading to move attempts at `(0,0)` and subsequent `400 Bad Request` errors ("Entity not adjacent").
**Fix:** Updated `internal/script/agent.go` to use `goja.TagFieldNameMapper("json", true)`, ensuring all Go structs passed to JS respect their `json` tags (lowercase).

## 3. Current Status: PvE Targeting Logic
**Status:** **[Partially Resolved]**
Bots now successfully:
1. Register and join the matchmaking queue.
2. Receive the `MatchFound` event and initialize the arena.
3. Calculate valid paths and execute move actions (no more `400` errors on moves).

**Remaining Blocker:**
In `1v1_PVE` mode, the bot script occasionally reports **"No enemies left. Passing."** and ends its turn even when enemies are present.
- **Root Cause Hypothesis:** The `filter` logic in `samples/pvp_bot_battle.js` (line 157) relies on `e.player_id !== myPlayerId`. In PvE, AI enemies often have `player_id: null` or a specific internal ID that might not be correctly compared against the bot's resolved ID.
- **Observed Behavior:** The identity resolution log shows `Identity Resolved! My Player UUID: [UUID]`, but the turn concludes prematurely.

---

## Technical Debt / Next Steps
1. **Refine Targeting:** Update `pvp_bot_battle.js` to handle PvE enemies by checking team IDs instead of just `player_id`.
2. **Engine TTL:** Observed `400 "arena not found"` errors during long debugging sessions, suggesting the engine cleans up inactive matches faster than the CLI-farm's full registration/handshake cycle under load.
3. **Verification:** All services (Laravel, Reverb, Engine) are currently operational and verified via `check_services.sh`.

## Relevant Logs
- Initial failure: `tests/logs/verification_1v1_PVE.log`
- Current behavior: `tests/logs/final_verification_1v1_PVE.log`
