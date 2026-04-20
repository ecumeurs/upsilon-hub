# AI Usage Manifesto: The UpsilonBattle Experiment

## The New Frontier of Collaboration

In this repository, development is not a solitary act of humans or a purely automated push from machines. It is a **shared cognitive journey**. We have moved beyond "using AI tools" to a paradigm where **AI Agents are first-class contributors**, governed by the same strict architectural rules as their human counterparts.

This document outlines the philosophy, the social contract, and the technical bridge that makes this collaboration possible.

---

## 1. The Shared Mind: Atomic Traceable Documentation (ATD)

The greatest barrier to AI assisted development is **hallucination**—the gap between what a human *thinks* and what an AI *guesses*. We solve this with **ATD**.

Every rule, mechanic, and API endpoint is frozen into an **Atom** (`.atom.md`). These atoms are not just "documentation"; they are the **Shared Mind** of the project. 

- **Ground Truth**: When an AI (like Antigravity) approaches a task, it doesn't speculate on intent. It queries the ATD index.
- **Intent Locking**: By requiring every feature to start with a `DRAFT` atom, we ensure that the "Why" is locked before the "How" is typed.
- **Traceability**: The link between code and docs is sacred. A `@spec-link` is a handshake of correctness.

*Reference:* [CLAUDE.md](CLAUDE.md), [docs/](docs/)

---

## 2. Intent-Locked Development

We operate under a simple mandate: **No code without an atom.**

If a rule is too small to be an atom, it's a detail. If it's big enough to change the system state, it must be an atom. This forces the AI to think in single-responsibility units, preventing the "drift" that often occurs in long coding sessions.

> [!IMPORTANT]
> This structure allows the AI to "load" specific system segments into its context window, ensuring high-fidelity implementation without the noise of unrelated logic.

---

## 3. Resilient Memory: Issue Management

AI agents are ephemeral. They "die" at the end of a session and are "reborn" with a fresh context. To prevent the loss of critical technical intuition, we use **Formal Issue Management**.

Every bug discovered, every risk identified (like a potential deadlock or a security gap), and every piece of technical debt is recorded as an `ISS-NNN` document in [issues/](issues/).

- **Persistent Memory**: These files act as the "Technical Memory" of the repository, surviving across agent resets.
- **Ref-Referencing**: By referencing `ISS-NNN` in atoms and commits, we create a historical record that allows the AI to understand *why* a certain refactor happened months ago.

---

## 4. Human-AI Symbiosis

We recognize distinct roles in this ecosystem:

| Role | Responsibility | Primary Actor |
| :--- | :--- | :--- |
| **Architect** | High-level design, review of `DRAFT` atoms, final approval. | Human |
| **Implementer** | Writing code, matching specifications, establishing links. | AI Agent |
| **Auditor** | Crawling the graph, finding orphans, checking for drift. | AI Agent |
| **Validator** | Running test suites, verifying BRD compliance. | Hybrid |

---

## 5. A Living (and Evolving) System

> [!WARNING]
> **Work in Progress:** The ecosystem described here—specifically the ATD tooling and the automated traceability workflows—is a living experiment. The tools are frequently updated, the logic of "bloating factors" is being refined, and the very way we communicate intent is evolving.

Users should expect friction. Tools might fail, indices might need rebuilding, and the "Social Contract" between human and machine is subject to renegotiation as models like **Gemini** and the **MCP (Model Context Protocol)** framework advance.

---

## 6. Technical Appendix

### The AI Stack
- **The Orchestra of Models**: UpsilonBattle utilizes a multi-model strategy to balance speed and reasoning:
    - **Gemini 3 Flash**: The primary Agentic Assistant (Antigravity) for real-time code generation and high-level strategy.
    - **Llama 3.2 / 3.1 (8B)**: Handlers for atom auditing, intent extraction, and structural snapshots.
    - **DeepSeek R1 (7B)**: Specialist for code auditing, congruence checking, and complex reconciliation logic.
    - **Qwen 2.5 Coder (14B)**: Utilized for high-fidelity file dissection and surgical reconstruction.
    - **Nomic Embed Text**: The mathematical foundation for the semantic documentation index.
- **Control Interface**: [Antigravity](https://github.com/google-deepmind) (Agentic AI Coding Assistant).
- **Governance**: ATD MCP Server — Provides the AI with direct tools to `search`, `trace`, `weave`, and `update` the documentation graph.

### Key Knowledge Pillars
- **[BRD.md](BRD.md)**: The ultimate source of business truth.
- **[SSD.md](SSD.md)**: The technical bridge between business goals and implementation.
- **[issues/README.md](issues/README.md)**: The index of known system risks.

---

*Drafted by Antigravity in collaboration with the UpsilonBattle Team.*
