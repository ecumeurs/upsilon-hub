# Business Requirement Document (BRD) - Upsilon Battle

## 1. Executive Summary

### 1.1 Project Vision
Upsilon Battle is a tactical RPG (TRPG) designed for competitive and cooperative turn-based combat. The project focuses on high-speed engagement, and deep tactical depth through initiative-based mechanics and character progression.

### 1.2 Business Objectives
- Provide a responsive tactical RPG experience across PvE and PvP modes.
- Maintain game balance through strict progression scaling and action economy rules.

**Source References:**  
- [[module_game]]: Defines overarching TRPG tactical flow.
- [[module_backend]]: Backend architectural goals.
- [[module_frontend]]: UI/UX design constraints.

---

## 2. Functional Requirements

### 2.1 User Onboarding & Identity
The system must support a "frictionless" entry for new players.
- Registration requires `Account Name`, `Password`, `Full Address`, and `Birth Date`.
- **Friction Balance:** While additional data is collected for compliance/external requirements, the onboarding process aims to remains streamlined.
- Success results in immediate login and redirection to the Dashboard.

**Source References:**  
- [[req_player_experience]]: Root player journey requirement.
- [[uc_player_registration]]: UC-1: Player account and roster creation.
- [[uc_player_login]]: UC-2: Player authentication and dashboard access.
- [[entity_player]]: Core player data structure with roles.

### 2.2 System Administration & Management
The system provides a dedicated role for administrators to maintain game integrity and manage the user base.
- **User Management:** Administrators can list all accounts and perform soft deletions (adhering to [[rule_gdpr_compliance]]).
- **History Management:** Administrators can review all match results and perform database maintenance (cleaning history older than 90 days).
- **Privacy Barrier:** Administrators are strictly prohibited from viewing private user data (Full Address, Birth Date).

**Source References:**  
- [[req_admin_experience]]: Root administrative requirement.
- [[uc_admin_login]]: UC-7: Administrator authentication.
- [[uc_admin_user_management]]: UC-8: Administrative account control (Soft Delete/GDPR).
- [[uc_admin_history_management]]: UC-9: Match history auditing and maintenance.
- [[rule_admin_access_restriction]]: Privacy gate for admin roles.

### 2.3 Matchmaking Ecosystem
Players must be able to quickly find games against other players or the system.
- **PvE Mode:** Instant game start against AI opponents.
- **PvP Mode:** Queue system to find human opponents.
- **Queue Selection:** Win/Loss records must be visible in the queue selection screen.
- **Lifecycle:** Players can voluntarily "Leave Queue" to return to the Dashboard at any time.

**Source References:**  
- [[uc_matchmaking]]: UC-3: PvE/PvP matchmaking and queue management.
- [[us_queue_selection]]: User story for choosing game modes.

### 2.4 Combat Engine & Turn Management
A rigid tactical engine governs the flow of battle.
- **Initiative:** Turn order is mathematically determined and displayed to players.
- **Turn Timer:** Players have a 30-second "shot clock" for actions.
- **Action Economy:** Valid actions include Move, Attack, Pass, and **Forfeit**.
- **Auto-Pass Penalty:** Timing out results in a forced Pass action and a +400 delay penalty (300 base + 100 penalty).
- **Match Resolution Loop:** Every action triggers a state evaluation. If no winner is detected, the turn passes back to the Combat Engine for the next character.
- **Integrity:** Friendly fire is strictly forbidden.
- **Action Feedback:** Every action must return a structured "Action Report" describing the full state mutation (path, damage, effects) to support rich UI visualization.

**Source References:**  
- [[uc_combat_turn]]: UC-4: Tactical turn lifecycle and action management.
- [[uc_match_resolution]]: UC-5: Win detection, forfeit handling, and reward triggers.
- [[requirement_customer_action_reporting]]: Rich visualization requirement.
- [[api_go_action_feedback]]: Technical contract for reporting.
- [[us_take_combat_turn]]: User story for combat actions and timer.
- [[mech_initiative]]: Initiative roll and requeue calculation.
- [[mech_action_economy]]: Cost of actions and timeout rules.
- [[rule_friendly_fire]]: Domain rule preventing allied hits.

### 2.5 Character Progression
Players can improve their roster through successful combat participation.
- **Post-Win Reward:** 1 attribute point allocated per game win.
- **Attribute Cap:** Total attributes must not exceed `10 + total_wins`.
- **Gated Movement:** Movement stats can only be increased every 5 accumulated wins.

**Source References:**  
- [[uc_progression_stat_allocation]]: UC-6: Manual stat adjustment from the Dashboard.
- [[rule_progression]]: Governing logic for stat allocation and limits.
- [[us_win_progression]]: Experience of allocating points after a win.

---

## 3. Non-Functional Requirements

### 3.1 Security & Access Control
- **Authentication:** All non-public endpoints MUST be protected via Laravel Sanctum using Bearer tokens.
- **Encryption:** All application traffic MUST be served over HTTPS (Self-signed certificates permitted for light/standard deployments).
- **Password Policy:** 
  - Minimum length: 15 characters.
  - Required: 1 Uppercase, 1 Number, 1 Special Symbol.

**Source References:**  
- [[req_security]]: Core security requirements.
- [[req_security_authorization]]: Access control rules.
- [[rule_password_policy]]: Specific complexity constraints.

### 3.2 GDPR & Data Privacy
The system implements a "Safe by Design" approach to user privacy.
- **Soft Deletion:** Account deletion requests trigger a "soft delete" flag rather than immediate record purging to maintain system integrity.
- **Anonymization:** Upon deletion or formal request, sensitive fields (`Full Address`, `Birth Date`) are programmatically overwritten with "ANONYMIZED" placeholders.
- **Data Portability:** Users MUST be able to request a full account data dump in a machine-readable JSON format, encompassing all personal identity data and historical records.
- **Private Data Handling:** Address and birth date are treated as strictly private and are excluded from all public-facing metrics or leaderboards.
- **Internal Identity Privacy:** Primary database UUIDs MUST NOT be exposed to the client. The system must use secure pseudonyms (Tactical IDs) and resolve identity purely via secure session context to prevent ID enumeration.

**Source References:**  
- [[rule_gdpr_compliance]]: Anonymization and soft delete logic.
- [[entity_player]]: Definition of private data fields.
- [[api_profile_export]]: Technical endpoint for data portability.
- [[requirement_customer_user_id_privacy]]: Identity/Ownership protection.
- [[arch_api_id_masking_gateway]]: ID masking mechanics.

### 3.3 Traceability & Error Management
- Every request must be traceable via a unique Request ID.
- Structured logging must be implemented for all core game loops.

**Source References:**  
- [[req_logging_traceability]]: Traceability requirements.
- [[api_request_id]]: Implementation details for tracking headers.

### 3.4 API-First & Developer Experience
- **Total Accessibility:** 100% of game-critical actions must have equivalent API endpoints.
- **Self-Discovery:** The API must expose an automated `/help` registry listing every available URI and contract.
- **Session Management:** The UI must handle session expiration gracefully via an immersive modal when neural synchronization (JWT) fails.

**Source References:**  
- [[requirement_customer_api_first]]: API playability and discovery.
- [[api_help_endpoint]]: Self-documenting registry endpoint.
- [[requirement_req_ui_session_timeout]]: Graceful session termination.

---

## 4. UI & Dashboard
The user interface must be intuitive and reflect the current state of the game and player progress.
- **Leaderboard:** Display ranked players with metrics like Total Wins and Movement.
- **Roster Management:** View character stats and allocate points.
- **Real-time Combat:** Interactive tactical board with live state updates via WebSockets.
- **Session Management:** Graceful handling of session timeouts and re-authentication.

**Source References:**  
- [[ui_dashboard]]: Main landing page after login.
- [[ui_leaderboard]]: Ranking display requirements.
- [[us_leaderboard_view]]: User story for viewing competitive standings.
- [[ui_battle_arena]]: Real-time combat interface.
- [[requirement_req_ui_session_timeout]]: Session timeout handling.

---

## 5. Implementation Status & Coverage

### 5.1 Documentation Coverage Analysis
Based on comprehensive ATD investigation (2026-04-17):

- **Total ATOMs**: 243 documentation atoms
- **True Orphans**: 8 atoms (3%) - intentionally unimplemented features
- **Implementation Coverage**: ~82% of STABLE atoms have corresponding code implementations
- **Traceability**: 421 @spec-link occurrences across 250 code files

### 5.2 Fully Implemented Features
✅ **Authentication & Identity**: Complete registration, login, logout with JWT tokens
✅ **Character Management**: Full roster system with reroll mechanics
✅ **Matchmaking**: PvE and PvP queue systems
✅ **Combat Engine**: Initiative-based turns, action economy, move/attack/validation
✅ **Progression System**: Win-based attribute point allocation
✅ **Real-time Updates**: WebSocket-based state broadcasting
✅ **Leaderboard**: Mode-based rankings with weekly cycles
✅ **Basic Admin**: User management and soft deletion

### 5.3 Partially Implemented Features
🔄 **Advanced Admin**: History management and audit trails (basic structure exists)
🔄 **GDPR Compliance**: Soft deletion implemented, full anonymization in progress
🔄 **Action Reporting**: Basic state updates, rich visualization planned
🔄 **Match History**: Basic logging, player history views pending

### 5.4 Planned Features (Not Yet Implemented)
📋 **PvP Stalemate Detection**: Draw conditions for infinite matches (ISS-029)
📋 **Rich Action Feedback**: Detailed animation data for UI effects
📋 **Advanced Privacy**: Complete "Right to be Forgotten" implementation
📋 **Personal Match History**: Player's historical match records

### 5.5 Known Issues & Gaps
- **Security**: Match participant access control needs enhancement (ISS-018)
- **Performance**: Arena lifecycle management improvements needed (ISS-012)
- **Traceability**: Request ID logging consistency (ISS-023, ISS-042)
- **Testing**: Go unit test flakiness in concurrent scenarios (ISS-047)

**Documentation References**: See `atd_investigation/` directory for detailed analysis.
