# Issue: Request Traceability Non-Compliance and Gaps

**ID:** `20260415_request_traceability_gaps`
**Ref:** `ISS-042`
**Date:** 2026-04-15
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi`, `battleui`, `upsiloncli`
**Affects:** `Developers, Debugging Workflows`

---

## Summary

This issue documents the systematic non-compliance with `rule_tracing_logging` across the Upsilon stack. While the architectural rule exists theoretically, the practical implementation in both Go and Laravel backends is fragmented, making end-to-end request tracing difficult and reliant on fragmented logging patterns.

---

## Technical Description

### Background
Rule `rule_tracing_logging` defines a strict log format:
`[{YYYY-MM-DDTHH:MM:SSZ}] [{ref_id}] [{ENDPOINT_OR_HANDLER}] {Message}`
Where `ref_id` is the first 8 characters of a UUIDv7 `request_id`.

### The Problem Scenario
Investigation reveals the following gaps:

1.  **Go Backend (`upsilonapi`):**
    - **Inconsistent Format:** Logs use standard `logrus` output (`time="..." level=... msg="..."`) or Gin's default ingress log.
    - **Ref ID Logic:** The `ref_id` is sometimes manually prepended to the message as `[R ref_id]`, but it's not a global requirement in the logging engine.
    - **Ingress Gap:** HTTP requests entering Go via Gin do not log the `ref_id` or `request_id` in a reachable way unless an internal handler explicitly logs it.
    - **Webhook IDs:** Webhooks generate *new* IDs via `stdmessage.New` which are not explicitly linked to the parent `request_id` in the log stream (though they are linked via the match context).

2.  **Laravel Backend (`battleui`):**
    - **Exception Only:** The atomic log format is currently only implemented in the global exception handler in `bootstrap/app.php`.
    - **Middleware Missing:** There is no global middleware that logs successful request ingress/egress in the atomic format.
    - **Proxying:** When proxying to Go, the `request_id` is sent in the JSON body but not in the `X-Request-ID` header, which might limit some infrastructure-level tracing tools.

3.  **CLI (`upsiloncli`):**
    - **Transient Logs:** CLI logs are not persisted by default.
    - **Agent Focus:** Logs are prefixed with `[Bot-N]` instead of `[ref_id]`, creating a semantic gap when trying to correlate agent actions with backend traces.

### Where This Pattern Exists Today
- [battleui/bootstrap/app.php](file:///workspace/battleui/bootstrap/app.php#L51-65) - Exception case only.
- [upsilonapi/bridge/http_controller.go](file:///workspace/upsilonapi/bridge/http_controller.go#L108) - Generates new IDs without explicit parent linkage in logs.
- [upsilonapi/engine.log](file:///workspace/upsilonapi/engine.log) - Shows non-compliant log output.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (makes debugging production failures significantly slower) |
| Detectability | High (logs are visibly non-compliant) |
| Current mitigant | None, besides manual grepping for parts of UUIDs |

---

## Recommended Fix

**Short term:** 
- Finalize and deploy `upsilon_trace_analyzer.py` to automate the correlation of existing fragmented logs.
- Allow CLI to pipe logs to a file that the analyzer can ingest.

**Medium term:**
- Implement a logging middleware in Go (Gin) that enforces the atomic format for all ingress requests.
- Implement a similar middleware in Laravel.

**Long term:**
- Enforce the `X-Request-ID` header propagation throughout the entire proxy chain.

---

## References

- [rule_tracing_logging.atom.md](file:///workspace/docs/rule_tracing_logging.atom.md)
- [api_standard_envelope.atom.md](file:///workspace/docs/api_standard_envelope.atom.md)
- [upsilon_log_parser.py](file:///workspace/upsiloncli/upsilon_log_parser.py) (Reference for existing parser logic)
