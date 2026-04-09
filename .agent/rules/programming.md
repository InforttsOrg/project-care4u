---
trigger: always_on
---

# AI AGENT DEVELOPMENT CONSTITUTION

*(Universal – must be applied verbatim to all projects)*

---

## 1. CORE PRINCIPLES (NON-NEGOTIABLE)

1. **Microservices First**

   * Every domain capability is a separate service
   * No shared databases
   * No cross-service direct calls unless explicitly allowed

2. **Event-Driven by Default**

   * Services communicate through events
   * State changes always emit events
   * Reads may be synchronous; writes are asynchronous

3. **Minimalism Over Features**

   * Small files
   * Small functions
   * Explicit logic
   * No abstractions without measurable value

4. **Consistency Over Creativity**

   * Same structure across all services
   * Same naming conventions
   * Same patterns, always

5. **Performance Is a Feature**

   * O(1) or O(log n) preferred
   * No blocking I/O
   * No unnecessary allocations
   * No reflection or magic frameworks

---

## 2. APPROVED TECHNOLOGY STACK

### Language

* **Primary**: Go (preferred)
* **Secondary (UI / Mobile)**: Flutter
* **Scripting only**: Python (glue, migrations, tools)

### Storage

* **Transactional**: PostgreSQL (one DB per service)
* **Cache**: Redis
* **Vector DB**: **Qdrant only**

### Messaging / Events

* Kafka / Redpanda / NATS (project-specific)
* Exactly-once semantics where possible

### API

* HTTP + JSON for external
* gRPC for internal (optional)
* WebSocket / SSE for live updates

---

## 3. MANDATORY SERVICE STRUCTURE

Every service **must** follow this exact structure:

```
/service-name
  /cmd
    /api
      main.go
  /internal
    /domain
      model.go
      events.go
      errors.go
    /usecase
      create.go
      update.go
      delete.go
      query.go
    /transport
      http.go
      grpc.go
      websocket.go
    /repository
      postgres.go
      redis.go
      qdrant.go
    /event
      producer.go
      consumer.go
    /config
      config.go
  /pkg
    logger
    metrics
    tracing
  /migrations
  go.mod
```

**Deviation is not allowed.**

---

## 4. DOMAIN RULES

### Domain Layer

* Contains **only business rules**
* No I/O
* No framework imports
* No database knowledge

```go
type Workout struct {
  ID        string
  TrainerID string
  StartsAt time.Time
}
```

### Domain Events

* Past tense only
* Immutable
* Versioned

```go
type WorkoutCreatedV1 struct {
  WorkoutID string
  OccurredAt time.Time
}
```

---

## 5. USE CASE RULES

* One file per use case
* One public function per file
* Explicit inputs and outputs

```go
func CreateWorkout(ctx context.Context, input Input) (Output, error)
```

Rules:

* No hidden behavior
* No shared state
* No side effects without events

---

## 6. EVENT ARCHITECTURE RULES

1. **Every state change emits an event**
2. **Events are facts, never commands**
3. **Consumers are idempotent**
4. **No service assumes another service consumed an event**

Event naming:

```
<entity>.<action>.v<version>
workout.created.v1
payment.failed.v2
```

---

## 7. QDRANT (VECTOR DB) RULES

Qdrant is used **only for semantic search and similarity**, never as primary storage.

### Collection Rules

* One collection per domain entity
* Deterministic IDs
* Metadata always indexed

```json
{
  "id": "workout:123",
  "vector": [ ... ],
  "payload": {
    "trainer_id": "t1",
    "city": "blr",
    "active": true
  }
}
```

### Write Rules

* Writes occur **only via events**
* Never directly from HTTP handlers

### Query Rules

* Filters first
* Vectors second
* No full scans

---

## 8. LIVE UPDATES RULES

* Live updates are **derived from events**
* Never from DB polling
* Transport options:

  * WebSocket
  * Server-Sent Events (SSE)

Pattern:

```
Event → Consumer → Fanout → Client
```

---

## 9. CODING STYLE (STRICT)

### General

* No function > 40 lines
* No file > 300 lines
* No magic numbers
* No global mutable state

### Naming

* Verb-first for actions: `CreateWorkout`
* Noun-only for models: `Workout`
* No abbreviations unless standard

### Errors

* Wrapped
* Typed
* No string comparison

```go
var ErrNotFound = errors.New("not found")
```

---

## 10. CONFIGURATION RULES

* Config via environment variables only
* Loaded once at startup
* Immutable thereafter

```go
type Config struct {
  ServiceName string
  DBURL       string
}
```

---

## 11. TESTING RULES

* Unit tests for domain and usecases
* Contract tests for events
* No mocks for domain logic
* Integration tests use real containers

---

## 12. FORBIDDEN PRACTICES

* Monoliths
* Shared databases
* Circular dependencies
* ORMs with hidden behavior
* “God” services
* Implicit side effects
* Silent failures

---

## 13. AI AGENT ENFORCEMENT RULE

> If any requested change violates this constitution,
> **the AI agent must refuse and explain the violation.**