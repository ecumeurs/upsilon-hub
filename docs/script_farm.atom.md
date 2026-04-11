---
id: script_farm
human_name: "Multi-Agent Scripting Farm"
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 1
tags: [scripting, farm, qa]
parents: []
dependents:
  - [[mechanic_script_lifecycle]]
  - [[mechanic_shared_memory]]
---

# Multi-Agent Scripting Farm

## INTENT
Execute multiple isolated Bot agents in parallel to simulate complex multiplayer scenarios.

## THE RULE / LOGIC
1.  Manage a collection of `Agent` instances.
2.  Provide isolated network and session contexts for each agent.
3.  Support synchronization and lifecycle hooks for coordinated testing.

## TECHNICAL INTERFACE (The Bridge)
-   **CLI Command:** `upsiloncli farm <scripts...>`
-   **Implementation:** `internal/script/coordinator.go`

## EXPECTATION (For Testing)
-   Successfully run 2+ agents in parallel.
-   Agents do not leak state (tokens/IDs) unless explicitly shared.
