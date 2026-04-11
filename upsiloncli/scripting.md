@spec-link [[script_farm]]

### **1. The `upsilon` API Reference**

Before writing a script, you need to know the tools available in your JavaScript environment.

#### Core Actions
* **`upsilon.call(route_name, {params})`**: Executes an API request. The `route_name` must match an endpoint in the registry (e.g., `auth_login`, `game_action`). Returns the parsed JSON `data` object.
* **`upsilon.waitForEvent(event_name, timeout_ms)`**: Pauses the script until the WebSocket receives the specified event (e.g., `match.found`, `board.updated`). Returns the event payload.
* **`upsilon.log(message)`**: Prints a message to the console, automatically prefixed with the agent's ID (e.g., `[Bot-1] Moving to 3,2`).

#### Session Context (Local to this specific Agent)
* **`upsilon.getContext(key)`**: Retrieves a value from the agent's local session (e.g., `user_id`, `match_id`).
* **`upsilon.setContext(key, value)`**: Overrides a context value.

* **`upsilon.onTeardown(function)`**: Registers a callback that is **guaranteed** to run when the script finishes, crashes, or fails an assertion (@spec-link [[mechanic_script_lifecycle]]).

#### Multi-Agent Sync (The Hive Mind @spec-link [[mechanic_shared_memory]])
* **`upsilon.setShared(key, value)`**: Writes data to a thread-safe global store visible to all concurrent agents.
* **`upsilon.getShared(key)`**: Reads data from the global store.
* **`upsilon.sleep(ms)`**: Pauses this specific agent without blocking others.

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
let gameState = upsilon.call("game_state", { id: matchEvent.match_id });
upsilon.assert(gameState.participants.length === 2, "Expected exactly 2 participants in a 1v1 match.");

// Find our specific entity using the user_id stored during auth
let myUserId = upsilon.getContext("user_id");
let myEntity = gameState.participants.find(p => p.user_id === myUserId);

upsilon.log("Attempting to move entity " + myEntity.id);

// Execute a move action
upsilon.call("game_action", {
    id: matchEvent.match_id,
    entity_id: myEntity.id,
    type: "move",
    target_coords: "3,3"
});

upsilon.log("Move executed successfully. Reaching end of script.");
// The script naturally ends here, which automatically triggers the onTeardown block.
```

---

### **3. Running Your Scripts**

Once your scripts are written, you execute them using the CLI coordinator we outlined. 

**Single Test Run:**
```bash
./bin/upsiloncli --farm my_scenario.js
```

**Multi-Agent Farm (e.g., 2v2 PVP):**
```bash
./bin/upsiloncli --farm bot_alpha.js bot_beta.js bot_gamma.js bot_delta.js
```
