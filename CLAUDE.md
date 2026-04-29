# Upsilon Hub: Agent Reference Index

Welcome to Upsilon Hub. To ensure consistent and high-quality development, follow the specialized rulesets linked below in order of priority.

## 1. Project Map & Infrastructure ([UPSILON.md](.agent/rules/UPSILON.md))
**MANDATORY READING FOR GROUNDING.**
Contains the "Who's Who" of the project, including the service architecture (Laravel, Go, CLI), folder organization, port mappings, and the testing toolkit. Read this first to understand the landscape.

## 2. Development Governance ([ATD.md](.agent/rules/ATD.md))
**CORE WORKFLOW.**
Defines the Atomic Traceable Documentation (ATD) lifecycle. Follow this for all feature development and bug fixes:
- **Atom Blueprint**: Strict structure for `.atom.md` files.
- **9 Consolidated Types**: From REQUIREMENT to MECHANIC.
- **Lifecycle Phases**: Discovery → Specification → Implementation → Verification.

## 3. Operational Guards & Standards ([COMMON.md](.agent/rules/COMMON.md))
**SAFETY & PROTOCOL.**
Standard rules for communication (API envelopes), error handling (Crash Early), and strict testing tool usage. Contains the "don'ts" of the project.

## 4. Issue Management ([issues.md](.agent/rules/issues.md) / [issue_management skill](.agent/skills/issue_management/SKILL.md))
**MAINTENANCE.**
Protocol for filing, tracking, and resolving technical debt, bugs, and risks in the `/workspace/issues/` directory using the system-wide `issues` command.

---

*Always prioritize these rules over generic assumptions. When in doubt, check UPSILON.md for location and ATD.md for intent.*