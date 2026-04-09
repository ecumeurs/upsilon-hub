---
id: api_help_endpoint
status: DRAFT
tags: api,discovery,meta
human_name: API Help & Discovery Endpoint
layer: ARCHITECTURE
version: 1.0
parents:
  - [[requirement_customer_api_first]]
dependents: []
type: API
priority: 3
---

# New Atom

## INTENT
To provide a programmatic, machine-readable index of all available API endpoints and their associated contracts.

## THE RULE / LOGIC
- **URI:** `/api/v1/help`
- **Verb:** `GET`
- **Intent:** Metadata Discovery
- **Fully Detailed Input:**
  - `scope`: (string) [Optional] Filter by category (auth, game, profile).
- **Fully Detailed Output:**
  - `version`: (string) Current API documentation version.
  - `endpoints`: (array) 
    - `uri`: (string) The full route path.
    - `verb`: (string) HTTP method (GET, POST, etc).
    - `intent`: (string) One-line purpose.
    - `input`: (array) parameters (name, type, description, mandatoriness, conditions).
    - `output`: (object) JSON structure of the successful response.

## TECHNICAL INTERFACE
- **API Endpoint:** `GET /api/v1/help`
- **Code Tag:** `@spec-link [[api_help_endpoint]]`
- **Test Names:** `TestHelpEndpointStructure`, `TestHelpEndpointDiscovery`

## EXPECTATION
- Request returns 200 OK JSON.
- Response contains a non-empty array of objects mapping to existing API atoms.
- Every entry includes URI, Verb, Intent, Input (params), and Output structures.
