# Issue: Antigravity Permission Denied in DevContainer

**ID:** `20260412_antigravity_permission_denied`
**Ref:** `ISS-030`
**Date:** 2026-04-12
**Severity:** High
**Status:** Resolved
**Component:** `.devcontainer/`
**Affects:** `.devcontainer/docker-compose.yaml`

---

## Summary

The Antigravity language server fails to start inside the devcontainer due to a `permission denied` error when trying to access state files (e.g., `installation_id`) in `/home/vscode/.gemini/antigravity/`. This is caused by mounting only a single file (`mcp_config.json`) into that directory, leading Docker to create the parent directories as `root`.

---

## Technical Description

### Background
The Antigravity agent stores its configuration and state in `~/.gemini/antigravity/`. In the devcontainer, the `vscode` user (UID 1001) is used. The host user `bastien` also has UID 1001.

### The Problem Scenario
In `docker-compose.yaml`, only one file is mounted:
```yaml
- $HOME/.gemini/antigravity/mcp_config.json:/home/vscode/.gemini/antigravity/mcp_config.json:cached
```
When Docker sets up this mount, if `/home/vscode/.gemini/antigravity/` does not exist in the image, it creates the directory structure owned by `root`. 
The language server, running as `vscode`, tries to create or open other files in this directory (like `installation_id`). Since the directory is owned by `root`, it fails with `permission denied`.

### Where This Pattern Exists Today
- [docker-compose.yaml](file:///home/bastien/work/upsilon/projbackend/docker-compose.yaml#L14)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High — Agent cannot start in devcontainer |
| Detectability | High — Language server error logs |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Update `docker-compose.yaml` to mount the entire `$HOME/.gemini/antigravity` directory instead of just the single file. This ensures the directory exists and inherits the host's permissions (UID 1001).

**Medium term:** Consider ensuring the directory exists in the `Dockerfile` with correct ownership, but mounting the whole directory is more robust for state persistence.

---

## References

- [.devcontainer/docker-compose.yaml](file:///home/bastien/work/upsilon/projbackend/docker-compose.yaml)
- [devcontainer.json](file:///home/bastien/work/upsilon/projbackend/.devcontainer/devcontainer.json)
