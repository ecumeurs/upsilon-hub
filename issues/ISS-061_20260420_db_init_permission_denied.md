# Issue: Permission Denied in db-init during Production Setup

**ID:** `20260420_db_init_permission_denied`
**Ref:** `ISS-061`
**Date:** 2026-04-20
**Severity:** High
**Status:** Resolved
**Component:** `battleui`
**Affects:** `docker-compose.prod.yaml` (db-init, app, ws services)

---

## Summary

The `db-init` service in the production Docker stack fails because it lacks permission to write to its own log files. This happens because the Dockerfile sets directory ownership *before* running `composer install`, which creates files as the `root` user.

---

## Technical Description

### Background
The `db-init` service is designed to automate database migrations on startup. It runs as user `www-data` (UID 33) to maintain a secure, non-root execution environment.

### The Problem Scenario
1. During `docker build`, the `battleui/Dockerfile` executes `RUN chown -R www-data:www-data storage bootstrap/cache`.
2. Afterwards, it executes `RUN composer install`.
3. Laravel's `post-autoload-dump` scripts run `@php artisan package:discover`, which boots the framework and creates `storage/logs/laravel.log` as `root` (the user running the build).
4. At runtime, the `db-init` container starts as UID 33.
5. It tries to run `php artisan migrate --force`, which attempts to append to `storage/logs/laravel.log`.
6. Result: `Failed to open stream: Permission denied`.

### Where This Pattern Exists Today
- [battleui/Dockerfile](file:///home/bastien/work/upsilon/projbackend/battleui/Dockerfile) at lines 43-49.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (100% on fresh production build) |
| Impact if triggered | High (Prevents the entire application stack from starting) |
| Detectability | High (Visible in `docker compose logs db-init`) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Move the `chown` command in the Dockerfile to a position after `composer install`.
**Medium term:** Add a more robust healthcheck or initialization script that verifies write permissions before booting the framework.
**Long term:** Standardize the use of a non-root user throughout the entire build and run lifecycle.

---

## References

- [docker-compose.prod.yaml](file:///home/bastien/work/upsilon/projbackend/docker-compose.prod.yaml)
- [battleui/Dockerfile](file:///home/bastien/work/upsilon/projbackend/battleui/Dockerfile)
