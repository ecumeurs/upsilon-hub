---
id: mechanic_mech_battle_engine_stress_testing
status: DRAFT
tags: ["testing","performance","stress-test"]
dependents: []
human_name: "Battle Engine Stress Testing Infrastructure"
parents:
  - [[shared:requirement_req_trpg_game_definition]]
version: 1.0
type: MECHANIC
layer: IMPLEMENTATION
priority: 3
---

# New Atom

## INTENT
To provide a robust infrastructure for long-running, multi-match stress tests of the battle engine to identify memory leaks, race conditions, and performance regressions.

## THE RULE / LOGIC
1. **Discovery:** Identifies running PIDs for `upsilonapi` and `upsilonbattle` (if running as separate process) to monitor resource usage.
2. **Orchestration:** Spawns multiple `upsiloncli` instances in parallel, each running a fast-paced bot battle.
3. **Monitoring:** Polls system metrics (CPU, Memory, FDs) at regular intervals (default 10s).
4. **Resilience:** Monitors match processes and respawns them if they terminate before the test duration expires.
5. **Consolidation:** Aggregates logs from all match instances and parses them to extract tactical metrics (actions, deaths, errors).
6. **Reporting:** Generates both machine-readable JSON and human-readable Markdown reports.

## TECHNICAL INTERFACE
- **Script:** `scripts/stress_test.py`
- **Output:** `/workspace/stress_test_report.json`, `/workspace/stress_test_report.md`
- **Code Tag:** `@spec-link [[mech_battle_engine_stress_testing]]`

## EXPECTATION
- The stress test can run for 10 minutes with 10+ concurrent matches without memory leaks in the Go engine.
- The resulting JSON report contains non-zero action counts and accurate outcome statistics.
- Service PIDs are correctly identified and monitored.
