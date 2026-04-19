# Upsilon Production Stack Setup Guide

This document details the architecture, configuration, and operation of the Upsilon production environment.

## 1. Architecture Overview

The Upsilon system is orchestrated using Docker Compose (Project Name: `upsilon-prod`) across six specialized services:

| Service | Responsibility | Healthcheck |
|---|---|---|
| `db` | PostgreSQL 18 database for all system state. | `pg_isready` |
| `db-init` | Lifecycle service that runs migrations on startup. | (One-shot) |
| `app` | Laravel 11 / Vue.js 3 frontend and REST API. | `http://localhost/up` |
| `ws` | Laravel Reverb WebSocket server for real-time battle updates. | `http://localhost:8080` |
| `engine` | Go-based battle engine for damage computation and logic. | `http://localhost:8081/health` |
| `cli` | Interactive CLI runner for debug and administrative tasks. | (Active) |

## 2. Shared Secrets Management

To ensure security and consistency, the stack uses a root-level `.env` file generated from `env.example`.

### Automatic Generation
The `scripts/setup_prod.sh` script automates the following requirements:
1. **Consistency**: It ensures `REVERB_APP_KEY` is the same in `app` and `ws`, and that `VITE_REVERB_APP_KEY` matches for the frontend build.
2. **Entropy**: It generates high-entropy random strings for:
   - `APP_KEY` (AES-256 encryption key)
   - `REVERB_APP_ID`
   - `REVERB_APP_KEY`
   - `REVERB_APP_SECRET`

### Propagation
All services share the `.env` file via the `env_file` directive in `docker-compose.prod.yaml`. This avoids duplicating secret definitions across the YAML file.

## 3. Operations Procedure

### First Start (Scratch)
Run the following commands from the project root:
```bash
chmod +x scripts/setup_prod.sh
./scripts/setup_prod.sh
docker compose -f docker-compose.prod.yaml up -d --wait
```

### Persistence & Restarts
The stack is designed for **Persistence-First** operation:
- **Data Safety**: PostgreSQL data is stored in the `db_data` named volume. This survives `docker compose down`, container deletions, and host reboots.
- **Restarts**: `docker compose restart` or `up -d` will not result in data loss. The `db-init` service will verify migrations but will not destructive-seed (it uses `migrate --force`).

### Port Mapping
| Host Port | Service Port | Scope |
|---|---|---|
| `5434` | `5432` | Postgres (Alternate to avoid host conflict) |
| `8000` | `80` | Laravel App / WebUI |
| `8080` | `8080` | WebSocket Server |
| `8081` | `8081` | Go Engine |

## 4. CLI Usage & Scripting

The `cli` service provides a headless environment for executing automation scripts and performing system checks.

### Basic CLI Commands
From the project root on the host, you can interact with the containerized CLI:

*   **System Status**: Verify reachability of Laravel API, Go Engine, and WebSockets.
    ```bash
    docker compose -f docker-compose.prod.yaml exec cli upsiloncli status
    ```
*   **Help**: List all available flags and commands.
    ```bash
    docker compose -f docker-compose.prod.yaml exec cli upsiloncli --help
    ```

### Introduction to Scripting
The CLI includes a "Farm" coordination engine that executes JavaScript scenarios. Each agent identifies itself via the `upsilon` global object.

#### Core API Reference
*   `upsilon.bootstrapBot(name, password)`: Handles registration, login, and automatic cleanup.
*   `upsilon.joinWaitMatch(mode)`: Joins the queue (e.g., `1v1_PVE`) and blocks until a match is found.
*   `upsilon.waitNextTurn()`: Blocks until it's the bot's turn to act.
*   `upsilon.call(action, params)`: Executes a direct game action (move, attack, pass).
*   `upsilon.planTravelToward(id, target, board)`: Calculates tactical paths toward a target.

#### Sample Test Script: PvE 1v1 Battle
You can run an automated tactical battle between a bot and the engine's internal AI.

1. Create a script named `test_pve.js`:
```javascript
const botName = "prod_tester_" + Math.floor(Math.random() * 1000);
upsilon.bootstrapBot(botName, "SecurePassword123!");

const match = upsilon.joinWaitMatch("1v1_PVE");
upsilon.log("Entered match: " + match.match_id);

while (true) {
    const board = upsilon.waitNextTurn();
    if (!board) break; // Game finished

    const me = upsilon.currentCharacter();
    const enemies = upsilon.myFoesCharacters();
    
    // Simplistic Logic: Move toward the first enemy found
    if (enemies.length > 0) {
        const path = upsilon.planTravelToward(me.id, enemies[0].position, board);
        upsilon.call("game_action", { 
            id: match.match_id, 
            entity_id: me.id, 
            type: "move", 
            target_coords: path.map(p => p.x + "," + p.y).join(";") 
        });
    }
}
```

2. Execute the script in the production stack:
```bash
docker compose -f docker-compose.prod.yaml exec cli upsiloncli --farm test_pve.js
```

## 5. Verification
... (rest of the content)
- If services fail to start, check logs: `docker compose -f docker-compose.prod.yaml logs -f`
- To reset everything (destructive): `docker compose -f docker-compose.prod.yaml down -v`
