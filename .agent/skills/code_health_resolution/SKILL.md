---
name: code_health_resolution
description: Run code health checks and resolve errors using ATD methodology for true architectural compliance.
---

# Code Health Resolution Skill

This skill provides a systematic approach to resolving code health errors reported by `scripts/code_health_check.py`. The objective is **True Resolution**: improving the code structure and documentation traceability, rather than just suppressing warnings.

## Core Methodology: The ATD Loop

Resolving health errors must follow the Atomic Traceable Documentation (ATD) lifecycle:
1.  **Analyze**: Run the health check and identify the root cause (Complexity? Missing specs? Bloat?).
2.  **Trace**: Use `atd_trace(summary=true)` to understand the atom's ancestry.
3.  **Map**: Use `atd_discover` and `atd_map` to align code boundaries with atoms.
4.  **Refactor**: Modify code/documentation to meet health standards.
5.  **Verify**: Run `atd_check` and `scripts/code_health_check.py` again.

## Triggers
This skill should be activated when:
*   `scripts/code_health_check.py` reports `[ERROR]` or `[WARN]`.
*   A user requests to "improve code quality," "resolve linting," or "fix atd coverage."
*   After a significant refactoring or feature implementation to ensure architectural health.
*   The `atd_check` tool indicates missing links or non-compliant implementations.

## Execution Flow

### Phase 0: Baseline Verification
1.  **Run Unit Tests:** Execute `scripts/run_all_unit_tests.sh`. 
    *   **CRITICAL:** If tests fail before you start, you must fix the logic errors first or report them to the user. Never attempt health resolution on a broken baseline.

### Phase 1: Diagnostics
1.  **Run Health Check:** `python3 scripts/code_health_check.py <target_path>`
2.  **Identify Targets:** Focus on `[ERROR]` messages first.
3.  **Context Gathering:** For each failing file, run:
    *   `atd_map(file="path/to/file")` to see current link recommendations.
    *   `atd_check(file="path/to/file")` to see current coverage.

### Phase 2: Resolution

#### For Complexity/Bloat Errors:
1.  Identify the logic that can be extracted.
2.  Search for existing atoms that might cover the sub-logic: `atd_search(query="logic description", scope="docs")`.
3.  If no atom exists, propose one: `atd_discover(file="path/to/file", new=true)`.
4.  Refactor the code into smaller, documented functions.
5.  Apply `@spec-link [[atom_id]]` directly above the new function.

#### For Missing ATD Links:
1.  Run `atd_discover(file="path/to/file")`.
2.  If it finds existing atoms, use `atd_update(file="path/to/file", spec_link="atom_id", spec_link_file="path/to/file")` to inject them (or do it manually above the relevant block).
3.  If it finds no atoms, you MUST find the BUSINESS requirement first.

### Phase 3: Vertical Traceability (The Business Link)
Every implementation MUST have a parent.
1.  Check the atom's `parents` field.
2.  If empty or only points to other IMPLEMENTATION atoms:
    *   **STOP.** You are missing the "Why".
    *   Use `atd_trace(atom="id", summary=true)` to check the full graph.
    *   If no BUSINESS atom is found in the ancestry, **ASK THE USER**:
        > "I am resolving health errors in `[[atom_id]]`, but it lacks a Business Layer origin. What is the core Business Rule or Requirement this code satisfies?"
3.  Create/Link the Business atom once the user provides the context.

## Prohibited Practices
*   **Global Headers:** Never put `@spec-link` at the top of a file. It must be "surgical" (above the logic).
*   **Defaulting:** Do not use `@lint-ignore` tags unless the user explicitly approves it for a specific, justified reason.
*   **Shallow Docs:** Do not add "The function X does Y" comments. Add "Intent: X is required to ensure Y constraint is met" style comments.

## Phase 4: Final Verification
1.  **Lint Check:** Run `python3 scripts/code_health_check.py <target_path>` to confirm 0 errors.
2.  **Structural Check:** Run `atd_weave()` followed by `atd_check(semantic=true)`.
3.  **Regression Check:** Run `scripts/run_all_unit_tests.sh` again.
    *   **CRITICAL:** Ensure the health resolution (comments, refactors) did not break existing functionality.
4.  **Build Check**: Run `scripts/build_services.sh` to ensure the code compiles successfully.
5.  **Start Services**: Run `scripts/start_services.sh` to ensure the services start successfully. Run seeding using `scripts/seed_ci.sh`
6.  **Integration Check**: Run `scripts/trigger_quick_ci_tests.sh` to ensure the services integrate successfully.

## Atom Location Policy
When creating new atoms during resolution:
*   **Local First:** Create the `.atom.md` file within the specific project's `docs/` directory (e.g., `upsilonapi/docs/`) if the logic is contained within that project.
*   **Transversal Only:** Use the root `docs/` directory ONLY for features that span multiple projects or describe high-level cross-cutting concerns (e.g., communication protocols, shared data structures).

## Tool Guardrails
*   **atd_map(file=...)**: Use this to get recommendations on where to place `@spec-link` tags.
*   **atd_check(semantic=true)**: Use this to verify that the code actually implements what the atom specifies.
*   **atd_weave()**: ALWAYS run this after creating or modifying atoms to sync the dependency graph.
