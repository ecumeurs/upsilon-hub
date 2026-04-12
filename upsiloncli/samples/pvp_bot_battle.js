// Upsilon Bot: PVP Battle Agent
// Adheres to [[rule_password_policy]]: 15+ chars, 1 uppercase, 1 digit, 1 special symbol

const botId = Math.floor(Math.random() * 100000);
const accountName = "bot_pvp_" + botId;
const password = "VeryLongBotPassword123!"; // 25 chars, satisfies policy

let matchId = ""; // Moved up for accessibility by teardown

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
            upsilon.call("game_forfeit", { id: matchId });
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

upsilon.log("Starting PVP Agent: " + accountName);

// 1. Register a new account
const regResponse = upsilon.call("auth_register", {
    account_name: accountName,
    email: accountName + "@example.com",
    nickname: "Gladiator_" + botId,
    password: password,
    password_confirmation: password,
    full_address: "Bot Street 1, Virtual Arena",
    birth_date: "1990-01-01T00:00:00Z"
});

if (regResponse && regResponse.token) {
    upsilon.log("Registration successful. JWT obtained.");
} else {
    upsilon.log("Registration failed: " + JSON.stringify(regResponse));
    throw new Error("Registration failed");
}

const myUserId = upsilon.getContext("user_id");
upsilon.log("My User ID: " + myUserId);

// 2. Join 1v1 PVP Queue
const slipMs = Math.floor(Math.random() * 3000);
upsilon.log("Preventing race condition: slipping for " + slipMs + "ms...");
upsilon.sleep(slipMs);

upsilon.log("Entering " + (upsilon.getEnv("UPSILON_GAME_MODE") || upsilon.getShared("game_mode") || "1v1_PVP") + " queue...");
upsilon.call("matchmaking_join", {
    game_mode: upsilon.getEnv("UPSILON_GAME_MODE") || upsilon.getShared("game_mode") || "1v1_PVP"
});

// 3. Wait for MatchFound
try {
    const matchData = upsilon.waitForEvent("match.found", 60000); // 60s timeout
    matchId = matchData.match_id;
    upsilon.log("Match Found! ID: " + matchId);

    // Sync check via shared memory to ensure both bots target the same match
    let sharedMatch = upsilon.getShared("current_test_match");
    if (!sharedMatch) {
        upsilon.log("First bot in. Setting shared expectation: " + matchId);
        upsilon.setShared("current_test_match", matchId);
    } else {
        upsilon.log("Second bot verified. Expected match: " + sharedMatch);
        upsilon.assert(sharedMatch === matchId, "Disjoint matches! Bot expected " + sharedMatch + " but joined " + matchId);
    }
} catch (e) {
    upsilon.log("Timed out waiting for match (match.found). Ensure another bot or player joins the 1v1_PVP queue.");
    throw e;
}

// 4. Battle Loop
upsilon.log("Entering battle loop...");
let gameOver = false;
let myResolvedPlayerId = "";

while (!gameOver) {
    try {
        // Wait for board update or turn started
        const eventData = upsilon.waitForEvent("board.updated", 60000);
        const board = eventData.data;
        
        if (board.winner_id) {
            upsilon.log("Game Over! Winner: " + board.winner_id);
            if (board.winner_id === myUserId) {
                upsilon.log("VICTORY IS MINE!");
            } else {
                upsilon.log("Defeated... perishing with honor.");
            }
            gameOver = true;
            break;
        }

        if (!myResolvedPlayerId) {
            // Resolve our UUID from participants list in initial state or update
            const stateResp = upsilon.call("game_state", { id: matchId });
            if (stateResp && stateResp.participants) {
                const me = stateResp.participants.find(p => p.nickname === accountName);
                if (me) {
                    myResolvedPlayerId = me.player_id;
                    upsilon.log("Identity Resolved! My Player UUID: " + myResolvedPlayerId);
                }
            }
        }

        if (myResolvedPlayerId && board.current_player_id === myResolvedPlayerId) {
            upsilon.log("--- My Turn! Acting with entity: " + board.current_entity_id + " ---");
            executeTacticalLogic(board, myResolvedPlayerId);
        } else {
            upsilon.log("Waiting for opponent (" + board.current_player_id + ")");
        }
    } catch (e) {
        upsilon.log("No activity detected. Checking state...");
        const stateResp = upsilon.call("game_state", { id: matchId });
        if (stateResp && stateResp.game_state) {
            const board = stateResp.game_state;
            if (board.winner_id) {
                upsilon.log("Game Over detected in poll. Winner: " + board.winner_id);
                gameOver = true;
            } else if (myResolvedPlayerId && board.current_player_id === myResolvedPlayerId) {
                executeTacticalLogic(board, myResolvedPlayerId);
            }
        }
    }
}

function executeTacticalLogic(board, myPlayerId) {
    const actingEntity = board.entities.find(e => e.id === board.current_entity_id);
    if (!actingEntity) {
        upsilon.log("Could not find my acting entity in board state.");
        return;
    }

    upsilon.log("[Unit: " + actingEntity.name + " | HP: " + actingEntity.hp + "/" + actingEntity.max_hp + " | Move: " + actingEntity.move + "]");

    const enemies = board.entities.filter(e => e.player_id !== myPlayerId && e.hp > 0);
    if (enemies.length === 0) {
        upsilon.log("No enemies left. Passing.");
        upsilon.call("game_action", { id: matchId, entity_id: actingEntity.id, type: "pass" });
        return;
    }

    let nearestEnemy = null;
    let minDistance = 1000;
    for (const enemy of enemies) {
        const dist = Math.abs(actingEntity.position.x - enemy.position.x) + Math.abs(actingEntity.position.y - enemy.position.y);
        if (dist < minDistance) {
            minDistance = dist;
            nearestEnemy = enemy;
        }
    }

    upsilon.log("Targeting nearest enemy: " + nearestEnemy.name + " at (" + nearestEnemy.position.x + "," + nearestEnemy.position.y + ")");

    let currentPos = actingEntity.position;

    // 1. Move toward enemy if not already adjacent
    if (minDistance > 1 && actingEntity.move > 0) {
        const pathSteps = upsilon.planTravelToward(actingEntity.id, nearestEnemy.position, board);

        if (pathSteps && pathSteps.length > 0) {
            const pathString = pathSteps.map(p => p.x + "," + p.y).join(";");
            
            upsilon.log("Moving " + pathSteps.length + " cells along path: " + pathString);
            upsilon.call("game_action", {
                id: matchId,
                entity_id: actingEntity.id,
                type: "move",
                target_coords: pathString
            });
            
            currentPos = pathSteps[pathSteps.length - 1];
            minDistance = Math.abs(currentPos.x - nearestEnemy.position.x) + Math.abs(currentPos.y - nearestEnemy.position.y);
        } else {
            upsilon.log("No valid path found to get closer to " + nearestEnemy.name);
        }
    }

    // 2. Attack if adjacent (either from start or after moving)
    if (minDistance === 1) {
        upsilon.log("Target in range! Attacking!");
        upsilon.call("game_action", {
            id: matchId,
            entity_id: actingEntity.id,
            type: "attack",
            target_coords: nearestEnemy.position.x + "," + nearestEnemy.position.y
        });
    }

    // 3. Always PASS to end the turn
    upsilon.log("Ending turn with pass.");
    upsilon.call("game_action", { id: matchId, entity_id: actingEntity.id, type: "pass" });
}

upsilon.log("Agent lifecycle complete.");
