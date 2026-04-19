# Issue: Implement High-Performance Manual Pagination and Search for Admin Tools

**ID:** `20260419_admin_performance_pagination`
**Ref:** `ISS-053`
**Date:** 2026-04-19
**Severity:** Medium
**Status:** Open
**Component:** `battleui/app/Http/Controllers/API/AdminController.php`, `battleui/database/migrations/`
**Affects:** Admin UI, API Performance, Database Scalability

---

## Summary

To support extremely large tables (specifically Match History), we must avoid standard pagination patterns that rely on expensive `COUNT(*)` queries. Instead, a manual "lot-based" pagination must be implemented to ensure the Admin Dashboard remains responsive even with millions of records.

---

## Technical Description

### Background
Standard Laravel pagination typically performs a `select count(*)` before the actual data query. For tables with high volatility and millions of rows, this count becomes a performance bottleneck (the "death of large table exploration").

### Key Requirements
1. **Manual Pagination:**
   - Implement "Next/Previous" or "Load More" logic without a total record count.
   - Fetch data in batches (lots) of **50**.
   - Sort strictly by `updated_at DESC`.
2. **Search Capability:**
   - Provide a filtering option by keyword (User Handle, Match UUID).
3. **Database Optimization:**
   - Verify and add composite or individual indexes on `updated_at` and any searchable columns to ensure O(log N) lookups.
4. **Efficiency:**
   - The API should return the data and potentially a simple "next_page_exists" flag based on fetching `LIMIT 51` and checking if a 51st record was found.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (as Match History grows) |
| Impact if triggered | High (Database timeout on admin pages) |
| Detectability | Medium (Performance degrades over time) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** 
- In `AdminController.php`, replace `get()` or `paginate()` with a manual `where('updated_at', '<', $cursor)->limit(50)->orderBy('updated_at', 'desc')` implementation.
- Create a migration to add an index on `updated_at` for the `users` and `matches` tables.

**Medium term:** 
- Apply this pattern to all Admin-facing registries defined in [ISS-051](ISS-051_20260419_admin_history_management_impl.md).

---

## References

- [uc_admin_history_management.atom.md](../../docs/uc_admin_history_management.atom.md)
- [AdminController.php](../../battleui/app/Http/Controllers/API/AdminController.php)
