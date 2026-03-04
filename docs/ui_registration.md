---
id: ui_registration
human_name: Registration Page UI
type: UI
version: 1.0
status: REVIEW
priority: CORE
tags: [ui, public, registration]
parents:
  - [[ui_landing]]
dependents:
  - [[mech_character_reroll]]
  - [[entity_player]]
---

# Registration Page UI

## INTENT
To allow new users to create an account and perform their initial character roster stat rolls.

## THE RULE / LOGIC
- Form Fields: Requires strictly minimal information (Account Name, Password). No email or additional personal data is collected.
- Character Generation Flow:
  - Upon submitting valid credentials, the user is presented with 3 randomly generated characters.
  - The UI must expose a "Reroll" action for the provided characters.
  - The reroll action is hard-capped at 3 total uses per account creation session.
- Success State: Completing this flow logs the user in (issues a JWT) and redirects to the Dashboard.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[ui_registration]]`
- **Test Names:** `TestRegistrationFormMinimal`, `TestRegistrationRerollLimit`

## EXPECTATION (For Testing)
- User generates account -> Sees 3 characters -> Can push reroll up to 3 times -> Reroll button locks after 3rd use -> Accepts roster and proceeds to Dashboard.
