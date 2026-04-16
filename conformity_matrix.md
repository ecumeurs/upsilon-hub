# Conformity Matrix

This document tracks the alignment between the Business Requirements (BRD), the Atomic Documentation (ATD), and the actual Implementation in the codebase.

## Core Requirements Alignment

| Requirement ID | Human Name | ATD Atom ID | Implementation Path | Status | Verification |
| :--- | :--- | :--- | :--- | :---: | :--- |
| **BRD 2.1** | User Onboarding | [[api_auth_register]] | `AuthController::register` | ✅ | `RegisterRequest.php` |
| **BRD 2.4** | Combat Engine | [[uc_combat_turn]] | `upsilonapi/handler/handler.go` | ✅ | `HandleArenaAction` |
| **BRD 2.4** | Initiative | [[mech_initiative]] | `turner.go` | ✅ | `NextTurn` |
| **BRD 2.4** | Action Economy | [[mech_action_economy]] | `turner.go` | ✅ | `EntityTurn` Delay |
| **BRD 2.4** | Turn Clock | [[rule_turn_clock]] | `ruler.go` | ✅ | `startShotClock` |
| **BRD 2.4** | Auto-Pass Penalty | [[mech_action_economy_timeout_penalty_rules]] | `ruler.go` | ✅ | +400 total delay |
| **BRD 2.4** | Friendly Fire | [[rule_friendly_fire]] | N/A | ❌ | **Blocked by [ISS-043]** |
| **BRD 2.5** | Character Progression | [[rule_progression]] | `ProfileController::updateCharacter` | ✅ | `10+wins` cap |
| **BRD 2.5** | Move Gating | [[rule_progression]] | `ProfileController::updateCharacter` | ✅ | `1 every 5 wins` |
| **BRD 3.1** | Password Policy | [[rule_password_policy]] | `RegisterRequest.php` | ✅ | 15 chars, U/N/S |
| **BRD 3.2** | GDPR Compliance | [[rule_gdpr_compliance]] | `AuthController::deleteAccount` | ✅ | Soft delete & Anonymize |
| **BRD 3.2** | Data Portability | [[api_profile_export]] | `AuthController::exportAccount` | ✅ | JSON export |

## Legend
- ✅ **Implemented:** Code follows the ATD specification.
- ⚠️ **Partial:** Implementation exists but misses specific constraints.
- ❌ **Missing:** No implementation found or blocked by issues.

## Detailed Traceability

> [!NOTE]
> For a full graph representation, run `atd trace <atom_id>` to view dependents and parents.
> The coverage ratio currently stands at ~6% due to missing @spec-link tags in legacy components. These are being backfilled.
