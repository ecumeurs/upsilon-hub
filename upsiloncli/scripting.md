@spec-link [[script_farm]]

### **1. The `upsilon` API Reference**

Before writing a script, you need to know the tools available in your JavaScript environment.

#### Core Actions
* **`upsilon.call(route_name, {params})`**: Executes an API request. The `route_name` must match an endpoint in the registry (e.g., `auth_login`, `game_action`). Returns the parsed JSON `data` object.
* **`upsilon.waitForEvent(event_name, timeout_ms)`**: Pauses the script until the WebSocket receives the specified event (e.g., `match.found`, `board.updated`). Returns the event payload.
* **`upsilon.log(message)`**: Prints a message to the console, automatically prefixed with the agent's ID (e.g., `[Bot-1] Moving to 3,2`).
* **`upsilon.planTravelToward(entity_id, target_pos, board)`**: Calculates the optimal path for an entity toward a coordinate. It handles grid occupancy (units/obstacles) and movement credit limits. Returns an array of positions.

#### Session Context (Local to this specific Agent)
* **`upsilon.getContext(key)`**: Retrieves a value from the agent's local session (e.g., `user_id`, `match_id`).
* **`upsilon.setContext(key, value)`**: Overrides a context value.

* **`upsilon.onTeardown(function)`**: Registers a callback that is **guaranteed** to run when the script finishes, crashes, or fails an assertion (@spec-link [[mechanic_script_lifecycle]]).

#### Multi-Agent Sync (The Hive Mind @spec-link [[mechanic_shared_memory]])
* **`upsilon.setShared(key, value)`**: Writes data to a thread-safe global store visible to all concurrent agents.
* **`upsilon.getShared(key)`**: Reads data from the global store.
* **`upsilon.sleep(ms)`**: Pauses this specific agent without blocking others. (Interruptible via Ctrl+C).

#### Tactical Utilities ([[script_farm]])
* **`upsilon.myPlayer()`**: Returns the participant record for the current agent.
* **`upsilon.currentPlayer()`**: Returns the participant record for the player whose turn it is.
* **`upsilon.currentCharacter()`**: Returns the entity currently selected for the turn.
* **`upsilon.myCharacters()`**: Returns an array of entities owned by the agent.
* **`upsilon.myAllies()` / `upsilon.myAlliesCharacters()`**: Returns allies (excluding self) or their entities.
* **`upsilon.myFoes()` / `upsilon.myFoesCharacters()`**: Returns opponents or their entities.
* **`upsilon.cellContentAt(x, y)`**: Returns `{ obstacle: bool, entity: Entity|null }` for a specific grid coordinate.

---

### **2. Building a Scenario: Step-by-Step**

Let’s build a robust test script. We will create an agent that creates an account, joins a PVE match, tries to make a move, and ensures the account is deleted afterward.

#### Step 2.1: Register the Teardown
Always start by defining your cleanup sequence. This ensures you don't litter your test database with ghost accounts.

```javascript
// my_scenario.js

upsilon.log("Starting PVE Scenario Test");

upsilon.onTeardown(() => {
    upsilon.log("Running Teardown: Deleting temporary account.");
    try {
        // auth_delete removes the account from the database
        upsilon.call("auth_delete", {}); 
    } catch (e) {
        upsilon.log("Failed to clean up account: " + e.message);
    }
});



// Register teardown for robust cleanup
upsilon.onTeardown(() => {
    upsilon.log("Running Teardown: Cleanup sequence triggered.");
    
    // 1. Ensure we leave matchmaking queue if still waiting
    try {
        upsilon.call("matchmaking_leave", {});
    } catch (e) {
        // Expected if not in queue
    }

    // 2. Forfeit if currently in a match
    if (matchId) {
        try {
            upsilon.log("Forfeiting match " + matchId);
            upsilon.call("game_action", { id: matchId, type: "forfeit" });
        } catch (e) {
            // Expected if game already ended or not our turn
        }
    }

    // 3. Always delete the temporary account
    try {
        upsilon.log("Deleting temporary account: " + accountName);
        upsilon.call("auth_delete", {});
    } catch (e) {
        upsilon.log("Failed to clean up account: " + e.message);
    }
});
```

#### Step 2.2: Setup & Authentication
Next, authenticate the agent. If the call fails, the script will crash safely and run the teardown.

```javascript
// Create or login to the account
upsilon.call("auth_register", { 
    account_name: "qa_bot_01", 
    password: "secure_password_123" 
});

upsilon.log("Authentication successful.");
```

#### Step 2.3: Trigger Contextual Flow
Use the registered endpoints to flow through the application logic.

```javascript
// Join the PVE matchmaking queue
upsilon.call("matchmaking_join", { game_mode: "1v1_PVE" });

// Wait for the Reverb WebSocket to push the match data
let matchEvent = upsilon.waitForEvent("match.found", 20000);

// Ensure the event actually fired before the 20s timeout
upsilon.assert(matchEvent != null, "Matchmaking timed out!");

// Store the match ID in the session for subsequent calls
upsilon.setContext("match_id", matchEvent.match_id);
```

#### Step 2.4: Execution and Assertions
Fetch the complex state, interact with it natively in JS, and assert the results.

```javascript
// Fetch the full tactical board
let board = upsilon.call("game_state", { id: matchEvent.match_id }).game_state;

// Identify self and foes easily
let me = upsilon.myPlayer();
let foes = upsilon.myFoesCharacters();

upsilon.log("I am " + me.nickname + ". Engaging " + foes.length + " enemies.");

// Execute a move action using my first character
let myUnits = upsilon.myCharacters();
if (myUnits.length > 0) {
    upsilon.call("game_action", {
        id: matchEvent.match_id,
        entity_id: myUnits[0].id,
        type: "move",
        target_coords: "3,3"
    });
}

upsilon.log("Move executed successfully. Reaching end of script.");
// The script naturally ends here, which automatically triggers the onTeardown block.
```

---

---

### **3. Running Your Scripts**

Once your scripts are written, you execute them using the CLI coordinator. 

**Execution Options:**
*   `--timeout <seconds>`: Global execution timeout. Triggers `onTeardown()` if reached.
*   `--logs <dir>`: Save individual agent logs to a directory.
*   `--auto`: (Experimental) Full jouney autopilot.

**Single Test Run:**
```bash
./bin/upsiloncli --farm my_scenario.js --timeout 60
```

**Multi-Agent Farm (e.g., 2v2 PVP):**
```bash
./bin/upsiloncli --farm bot_alpha.js bot_beta.js ... --logs ./logs
```
