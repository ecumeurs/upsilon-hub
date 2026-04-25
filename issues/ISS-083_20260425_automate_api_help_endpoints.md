# Issue: Automate API Help Endpoints using Postman and Validation Tools

**ID:** `20260425_automate_api_help_endpoints`
**Ref:** `ISS-083`
**Date:** 2026-04-25
**Severity:** Medium
**Status:** Open
**Component:** `battleui/app/Services/CodeDiscoveryService.php`
**Affects:** `battleui/app/Http/Controllers/API/HelpController.php`

---

## Summary

The current API help endpoints in the Laravel `battleui` component rely on a hybrid of self-reflection and manual parsing of `.atom.md` files and doc comments. This process is fragile and requires manual updates to keep the documentation in sync with the actual implementation and external tools like Postman. We need to automate this process by integrating the existing Postman collection and using proper validation tools to ensure the documentation accurately reflects the API contract.

---

## Technical Description

### Background

The `CodeDiscoveryService` currently iterates through Laravel routes and uses PHP Reflection to extract information from controllers and methods. It looks for `@spec-link` tags to link endpoints to ATD atoms and `@api-output` tags for DTO links. It also parses a specific DTO atom file (`docs/battleui_api_dtos.atom.md`) to build a DTO registry.

### The Problem Scenario

The current documentation workflow is partially manual and lacks a single source of truth:
1. Developers must maintain `@spec-link` and other tags in code.
2. The `CodeDiscoveryService` logic is complex and might miss changes if not updated.
3. The Postman collection (`Upsilon_Battle.postman_collection.json`) is maintained separately from the `help` endpoint output.
4. There is no automated validation between the code, the ATD atoms, and the Postman collection.

```
Current Flow:
Code (Reflexion) --\
                    >--> HelpController -> /api/v1/help
Docs (.atom.md)  --/

External:
Postman Collection <--- Manually Updated
```

### Where This Pattern Exists Today

- `battleui/app/Services/CodeDiscoveryService.php`
- `battleui/app/Http/Controllers/API/HelpController.php`
- `Upsilon_Battle.postman_collection.json`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium — Outdated documentation leading to integration issues |
| Detectability | Medium — Manifests as discrepancies between /help output and actual API behavior |
| Current mitigant | Manual updates and reflection-based discovery |

---

## Recommended Fix

**Short term:** Update `CodeDiscoveryService` to optionally read from the Postman collection to supplement endpoint information.

**Medium term:** Implement a validation tool that compares the `CodeDiscoveryService` output with the Postman collection and flags discrepancies.

**Long term:** Establish a single source of truth (e.g., OpenAPI/Swagger spec or the Postman collection) and generate both the `help` endpoint output and the Postman collection from it, or vice versa, ensuring full automation and validation.

---

## References

- [CodeDiscoveryService.php](file:///workspace/battleui/app/Services/CodeDiscoveryService.php)
- [HelpController.php](file:///workspace/battleui/app/Http/Controllers/API/HelpController.php)
- [Upsilon_Battle.postman_collection.json](file:///workspace/Upsilon_Battle.postman_collection.json)
