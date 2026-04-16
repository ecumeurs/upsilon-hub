# CI Testing Framework

This document outlines the automated verification strategy for Upsilon Battle, ensuring that code changes remain compliant with the Business Requirement Document (BRD) and Atomic Traceable Documentation (ATD).

## CI Strategy

Verification is performed via the **Upsilon CLI**, which executes JavaScript scenarios against a live service stack. Each scenario asserts specific rules and triggers a clean teardown.

### Command Structure

Scenarios are executed using the `--farm` option:
```bash
./bin/upsiloncli --farm ./samples/<scenario>.js --timeout 60
```

## Core Scenarios

### 1. Character Progression Lifecycle
*   **Target:** [[rule_progression]]
*   **Path:** `upsiloncli/samples/progression_check.js`
*   **Description:** Simulates a player journey from registration through a match win and attempts to upgrade character stats.
*   **Assertions:**
    *   Stat gain is allowed after a win.
    *   Stat gain is rejected if it exceeds `10 + wins`.
    *   Movement gain is rejected if not on a 5-win milestone.

### 2. Authentication & Security Policy
*   **Target:** [[rule_password_policy]]
*   **Path:** `upsiloncli/samples/auth_security_check.js`
*   **Description:** Attempts various registration payloads to verify server-side validation.
*   **Assertions:**
    *   Reject passwords < 15 characters.
    *   Reject passwords without numbers/symbols.
    *   Accept compliant passwords.

---

## Running in CI/CD

To run the full suite in a headless environment:
1.  Ensure all services are running (`start_services.sh`).
2.  Run the farm coordinator:
    ```bash
    ./bin/upsiloncli --farm ./samples/progression_check.js ./samples/auth_security_check.js --timeout 120
    ```

> [!IMPORTANT]
> **Exit Codes:** The CLI exit code reflects the success or failure of the farm. Any assertion failure in a script results in a non-zero exit code, blocking the CI pipeline.
