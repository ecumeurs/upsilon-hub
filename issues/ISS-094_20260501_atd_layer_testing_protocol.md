# ISS-094 — ATD Layer Testing Protocol (The Naive LLM)

**Date:** 2026-05-01
**Status:** Open — Planned Feature
**Layer:** ARCHITECTURE → ATD WebUI / CLI
**Priority:** P2 (quality tooling, not blocking current dev)

---

## Context

Now that the pre-commit hook (ISS-094 companion: `scripts/hooks/pre-commit`) enforces structural traceability, the next layer of quality assurance is **semantic correctness**: does the atom say what it needs to say for its layer?

The naive LLM approach is deliberately chosen. A 7B/8B model (e.g., `llama3.1:8b`, `deepseek-r1:7b` via Ollama) cannot guess context it was not given. If it answers "MISSING" to a question about an atom, that means the atom is genuinely thin — not that the model is weak. The weakness becomes a feature.

---

## Protocol Specification

### Layer 1 — BUSINESS Layer Test ("The Synthetic PM")

**Target types:** `PERSONA`, `WORKFLOW`, `REQUIREMENT`, `USECASE`, `USER_STORY`

**Goal:** Verify the atom explains *value*, *who*, and *rules* without relying on technical jargon.

**Prompt:**
```
You are a non-technical Product Manager. Read the following business documentation.
Based strictly on the text provided, answer the following three questions.
Do not infer or invent information. If the text does not contain the answer,
reply exactly with 'MISSING'.

1. Who is the primary user or actor for this feature?
2. What is the business value or goal of this feature?
3. What is one strict rule or acceptance criteria that must be met?

[ATOM CONTENT]
```

**Rejection criteria:**
- Any answer is `MISSING` → atom is **Too Thin**
- Any answer contains table names, API ports, or struct/class names → atom is **Leaking Implementation**

---

### Layer 2 — ARCHITECTURE Layer Test ("The Synthetic Tech Lead")

**Target types:** `CONTRACT`, `MODULE`, `API`, `ENTITY`

**Goal:** Verify the atom clearly defines boundaries, data ownership, and guarantees.

**Prompt:**
```
You are an Integration Engineer joining a new team. Read the following architectural documentation.
Based strictly on the text provided, answer the following three questions.
Do not infer or invent information. If the text does not contain the answer,
reply exactly with 'MISSING'.

1. What external systems, if any, does this component communicate with?
2. What exact data or state does this component own and persist?
3. What is the primary guarantee or contract this component provides to the rest of the workspace?

[ATOM CONTENT]
```

**Rejection criteria:**
- Answer to Q2 is "It doesn't own any data" for an ENTITY atom → documentation is **Flawed**
- Q3 answer is `MISSING` → atom is **Missing Contract** (a design atom with no stated guarantee is a file directory listing, not a specification)

---

### Layer 3 — IMPLEMENTATION Layer Test (Structural Only — No LLM)

**Target types:** `MECHANIC`

**Goal:** Prove the code exists and traces upwards. Semantic content is small enough for human PR review.

**Protocol (two structural checks, no prompt):**
1. **Upward link** — the pre-commit hook already enforces that a `parents:` entry exists
2. **Downward link** — `atd verify` confirms a `@spec-link [[mechanic_id]]` exists directly above the relevant function in source code

If both links exist, the implementation atom is valid. Prose quality is a PR concern, not a tooling concern.

---

## Implementation Notes

### CLI integration
Add `atd test-layer --atom <id>` (or `atd lint --semantic`) that:
1. Reads the atom file
2. Determines its layer and type
3. Selects the correct prompt template above
4. Sends to the configured Ollama model (suggested: `llama3.1:8b` for speed)
5. Parses the response for `MISSING` occurrences and implementation leak keywords
6. Returns `PASS` / `TOO_THIN` / `LEAKING_IMPL` / `MISSING_CONTRACT`

### WebUI integration
- Add a "Test Layer" button on the atom detail view
- Results displayed inline with a badge (PASS / TOO_THIN / LEAKING / MISSING_CONTRACT)
- Can be run on-demand or as part of the REVIEW→STABLE promotion workflow

### Batch mode
`atd test-layer --all --layer BUSINESS` to audit the full shared docs pool.
Priority: run on any atom being promoted from `DRAFT` to `REVIEW`.

### Model recommendation
- BUSINESS test: `llama3.2` (speed, comprehension adequate)
- ARCHITECTURE test: `llama3.1:8b` or `deepseek-r1:7b` (needs stronger reasoning for "guarantee" concept)
- IMPLEMENTATION: no model, structural only

---

## Related

- `scripts/hooks/pre-commit` — the structural pre-commit check (Part 1, implemented)
- `docs/req_tech_debt_backlog.atom.md` — escape hatch for the pre-commit hook
- `webui_ctrlk_doc_integration.md` — WebUI integration plan
- `ISS-082` — frontend Playwright test seams
