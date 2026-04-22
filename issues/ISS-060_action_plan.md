# Action Plan: Database Schema and Documentation Drift (ISS-060)

## Executive Summary

After analyzing the actual PostgreSQL schema against the documented schema in `db.md`, I've identified several critical discrepancies that need immediate attention. The documentation is mostly accurate but contains some outdated field references and missing critical details about nullable fields.

## Critical Discrepancies Found

### 1. **Users Table - Missing/Incorrect Fields**
**Severity: Medium**

#### Documented but Not Implemented:
- `email_verified_at` - Documented as "Timestamp, Nullable" but **does not exist** in the actual schema
- This is a standard Laravel field for email verification, but the project doesn't use this feature

#### Implemented but Missing Documentation:
- `role` field exists (default 'Player') but missing from db.md
- This is a critical field for RBAC (Role-Based Access Control)

### 2. **Match Participants Table - Nullable Field Mismatch**
**Severity: High**

- `player_id` field is **nullable** in the actual schema
- db.md does not document this as nullable
- This is critical for AI/bot participants or future playerless matches

### 3. **System Tables Not Documented**
**Severity: Low**

The following Laravel-specific tables exist but are not documented:
- `cache`, `cache_locks` - Laravel caching system
- `jobs`, `job_batches`, `failed_jobs` - Laravel queue system
- `password_reset_tokens`, `personal_access_tokens`, `sessions` - Laravel auth system

These are infrastructure tables that don't need detailed game logic documentation but should be acknowledged.

### 4. **Constraint and Index Details Missing**
**Severity: Low**

- `match_participants.status` has CHECK constraint (WIN/LOSS) - not documented
- `game_matches.game_mode` has CHECK constraint for allowed values - not documented
- Additional indexes for performance (e.g., `users_updated_at_index`) - not documented

## Action Plan

### Phase 1: Update Core Documentation (Immediate)

**Priority: HIGH**

1. **Update `db.md` to match actual schema:**
   - Remove `email_verified_at` from users table documentation
   - Add `role` field to users table with default value
   - Mark `match_participants.player_id` as nullable
   - Add CHECK constraints documentation

2. **Update Mermaid ERD:**
   - Add `role` field to USERS entity
   - Remove `email_verified_at` from USERS entity
   - Make `player_id` nullable in MATCH_PARTICIPANTS entity
   - Add CHECK constraints as notes in the diagram

### Phase 2: Update ATD Atoms (Short-term)

**Priority: MEDIUM**

1. **Create new ATD atom:**
   - `entity_users` - Comprehensive user entity specification including `role` field

2. **Update existing ATD atoms:**
   - `entity_player` - Add role field and remove email_verified_at reference
   - `entity_game_match` - Document the participant relationship with nullable player_id

3. **Update `data_persistence` atom:**
   - Add explicit mention of system tables (cache, jobs, etc.) as infrastructure
   - Clarify the role of Laravel-specific tables vs. game logic tables

### Phase 3: Verification & Testing (Medium-term)

**Priority: MEDIUM**

1. **Create ATD atom:**
   - `rule_database_schema_validation` - Schema validation rules

2. **Add integration tests:**
   - Verify database schema matches documentation
   - Test nullable player_id in match_participants (for AI/bot matches)

### Phase 4: Automation (Long-term)

**Priority: LOW**

1. **Implement automated schema validation:**
   - Script to compare actual schema with db.md
   - Run as part of CI/CD pipeline
   - Fail build if discrepancies detected

2. **Enhance documentation generation:**
   - Consider generating db.md from migrations automatically
   - Maintain manual overrides for business logic documentation

## Proposed Changes to db.md

### Users Table Updates:
```markdown
### 1. `users` (formerly `players`)
Stores authentication identity and tracks top-level metrics for the generic Leaderboard (`ui_leaderboard`).
* `id` (UUID, Primary Key)
* `account_name` (Varchar, Unique, Not Null)
* `email` (Varchar, Unique, Not Null)
* `password_hash` (Varchar, Not Null)
* `role` (Varchar, Default 'Player') - *For RBAC: Player/Admin*
* `remember_token` (Varchar, Nullable)
* `full_address` (Text, Nullable) - *Private: GDPR protected*
* `birth_date` (Date, Nullable) - *Private: GDPR protected*
* `total_wins` (Int, Default 0)
* `total_losses` (Int, Default 0)
* `reroll_count` (Int, Default 0)
* `ratio` (Numeric, Default 0)
* `ws_channel_key` (UUID, Unique, Nullable) - *For secure WebSocket subscriptions*
* `created_at` (Timestamp)
* `updated_at` (Timestamp)
* `deleted_at` (Timestamp, Nullable) - *Soft delete for GDPR compliance*

**Indexes**: `account_name`, `email`, `ws_channel_key`, `updated_at`
**Constraints**: `role` must be 'Player' or 'Admin'
```

### Match Participants Table Updates:
```markdown
### 4. `match_participants`
Mapping table defining which Users (or AI agents) competed in a specific historical or active match.
* `id` (BigInt, Primary Key, Auto-increment)
* `match_id` (UUID, Foreign Key -> `game_matches.id`)
* `player_id` (UUID, Foreign Key -> `users.id`, Nullable) - *Nullable for AI/bot participants*
* `team` (Int) - *Team 1 or Team 2*
* `status` (Varchar, Nullable) - *'WIN', 'LOSS' with CHECK constraint*
* `created_at` (Timestamp)
* `updated_at` (Timestamp)

**Indexes**: `match_id`, `player_id`
**Constraints**: `status` must be 'WIN' or 'LOSS'
**Note**: `player_id` is nullable to support AI/bot participants in PvE modes
```

## Impact Assessment

### Development Impact:
- **Low risk** - Documentation updates only, no schema changes required
- **Medium benefit** - Improved developer onboarding and reduced confusion
- **Critical benefit** - Prevents future bugs from schema misunderstanding

### Production Impact:
- **No impact** - Schema matches current implementation
- Documentation changes are non-breaking

### Testing Impact:
- **New tests needed** - Validate nullable player_id functionality
- **Existing tests** - Should pass as they match current schema

## Dependencies

- None - this is documentation-only work
- However, Phase 3 testing should coordinate with E2E test development
- Phase 4 automation should be integrated with existing CI/CD improvements

## Timeline Estimate

- **Phase 1**: 2-3 hours (immediate priority)
- **Phase 2**: 4-6 hours (next sprint)
- **Phase 3**: 8-12 hours (following sprint)
- **Phase 4**: 16-24 hours (future enhancement)

## Success Criteria

1. ✅ db.md accurately reflects actual PostgreSQL schema
2. ✅ Mermaid ERD matches actual table structures and constraints
3. ✅ ATD atoms properly document all game entities and relationships
4. ✅ Tests validate schema-documentation alignment
5. ✅ CI/CD pipeline includes schema validation (Phase 4)

## Conclusion

The database schema drift issue is primarily a documentation problem rather than a code problem. The actual implementation is correct and functional, but the documentation contains outdated references and missing critical details. Updating the documentation will prevent future confusion and ensure consistency across the development team.

The most critical fixes are:
1. Removing the non-existent `email_verified_at` field
2. Adding the missing `role` field
3. Correcting the nullable status of `player_id` in match_participants

These changes align with the ATD philosophy of ensuring documentation and implementation remain synchronized.
