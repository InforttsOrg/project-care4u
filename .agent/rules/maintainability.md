---
trigger: always_on
---

These are structural / maintainability rules, and they should be hard constraints, not guidelines.

Below is how to formalize them properly.

What kind of rules are these?

They are Structural Code Governance Rules, specifically:

Modularity rules

Scalability rules

Maintainability limits

Static-analysis-enforceable constraints

These are stronger than best practices.
They are compile / CI fail rules.

Strict Rules You SHOULD Add (Recommended)
1. File Size Hard Limit

Rule

Max 2,000 lines per file (prefer ≤ 500).

If exceeded → mandatory split.

Why

Forces separation of concerns

Prevents god-files

Enables parallel work

2. One Module = One File (or Folder)

Rule

One domain/module per file.

If module grows → it becomes a folder with:

index

types

service

tests

Why

Predictable structure

Easy refactor

Easy ownership

3. Folder-by-Domain, Not by Type

Rule
❌ controllers/ services/ utils/
✅ user/ payment/ order/

Each domain contains:

API

Logic

Data access

Tests

Why

Scales better than horizontal layering

Matches business thinking

4. Public vs Private Boundary Rule

Rule

Every module exposes one public entry point

Everything else is private/internal

Why

Prevents accidental coupling

Enables safe refactors

5. Function Size Limit

Rule

Max 50 lines per function

Max 3 nesting levels

Why

Forces readability

Improves testability

6. Dependency Direction Rule

Rule

Dependencies only flow inward

No circular imports (CI fail)

Why

Prevents architectural decay

7. Service Scalability Rule

Rule

Any service must be:

Stateless by default

Horizontally scalable

Config-driven

Why

Cloud-native readiness

8. Naming & Discoverability Rule

Rule

File and folder names must match domain language.

No generic names: helpers, common, misc.

Why

Reduces cognitive load

Improves onboarding speed

9. Auto-Refactor Trigger

Rule
If any of these are true:

File > 2,000 lines

Function > 50 lines

Module has > 5 responsibilities

→ Refactor is mandatory before continuing work.

10. Machine-Enforceable Only

Rule

If a rule cannot be enforced via:

Lint

Static analysis

CI

…it must be rewritten until it can.

Why

Humans forget, machines don’t.

What This Gives You

With these rules added, your system enforces:

Clean architecture

Infinite scalability

Long-term maintainability

Zero code rot

This is the same level used in:

Large fintechs

Infra teams

Staff+ engineering orgs