# HealthFlow

> **Event-driven, serverless healthcare booking & management platform with native apps**

[![CI/CD](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-blue)](https://github.com)
[![Go](https://img.shields.io/badge/Go-%3E%3D1.21-00ADD8)](https://go.dev)
[![Flutter](https://img.shields.io/badge/Flutter-%3E%3D3.x-02569B)](https://flutter.dev)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

---

## 1. Overview

HealthFlow is a full-stack, event-driven healthcare platform designed for appointment booking, provider discovery, payments, and real-time location awareness. The system is built with native-grade mobile apps (Flutter), serverless backend services (Go), infrastructure as code, and reliable CI/CD from day one.

### Architecture Principles

| Principle | Description |
|-----------|-------------|
| **Rapid MVP Delivery** | ≤ 2 weeks to production |
| **Horizontal Scalability** | Stateless services, auto-scaling |
| **Fault Tolerance** | Event-driven, async processing |
| **Clean Domain Boundaries** | Microservices with clear contracts |
| **Cloud-Native** | Serverless-first, IaC-managed |

---

## 2. Key Features

### End-User (Patients)

- ✅ User authentication (OTP / OAuth)
- ✅ Doctor / clinic discovery (map-based)
- ✅ Real-time availability & appointment booking
- ✅ Secure payments (UPI, Cards, NetBanking, Wallets)
- ✅ Booking history & invoices
- ✅ Location-aware search (nearby providers)

### Provider (Doctors / Clinics)

- ✅ Provider onboarding & verification
- ✅ Availability & schedule management
- ✅ Appointment lifecycle management
- ✅ Payout tracking (escrow-based)
- ✅ Analytics dashboard (MVP-lite)

### Platform

- ✅ Event-driven architecture
- ✅ Serverless compute
- ✅ Horizontal scalability
- ✅ Secure payment flows (Razorpay Escrow)
- ✅ Caching & rate limiting
- ✅ Infrastructure as Code
- ✅ Automated CI/CD

---

## 3. Technology Stack

### Mobile Apps

| Component | Technology |
|-----------|------------|
| Framework | **Flutter** (single codebase) |
| Platforms | Android, iOS |
| State Management | Riverpod / Bloc |
| Maps | Google Maps SDK |
| Location | GPS, background-safe services |
| Secure Storage | Flutter Secure Storage |

### Backend (Serverless)

| Layer | Technology |
|-------|------------|
| API Layer | **Go (≥1.21)** |
| Functions | AWS Lambda / Cloudflare Workers |
| Routing | API Gateway / Cloudflare Routes |
| Auth | JWT + OTP + OAuth |
| Async Events | EventBridge / PubSub |
| Crons | EventBridge Scheduler |
| Caching | Redis / Cloudflare KV |
| Database | PostgreSQL (RDS / Neon) |
| Search (optional) | OpenSearch |
| Storage | S3-compatible |

> [!NOTE]
> **Why Go?**
> - Predictable performance
> - Low cold start latency
> - Excellent concurrency model
> - First-class cloud tooling

### Payments

| Provider | Features |
|----------|----------|
| **Razorpay** | UPI, Cards, NetBanking, Wallets |
| Settlement | Escrow / Split settlements |
| Confirmation | Webhook-driven |
| Compliance | PCI exposure avoided (hosted checkout) |

### Infrastructure as Code (IaC)

| Tool | Purpose |
|------|---------|
| **Terraform** | Resource provisioning |
| **Terragrunt** | DRY configuration (optional) |

**Provisioned Resources:**
- API Gateway
- Lambda functions
- Event buses
- RDS / Serverless Postgres
- Redis
- Secrets Manager
- IAM (least privilege)

**Environments:**
- `dev`
- `staging`
- `prod`

### CI/CD

| Principle | Implementation |
|-----------|----------------|
| Deterministic | Reproducible builds |
| Immutable artifacts | Versioned releases |
| Zero manual steps | Fully automated |

**Tools:**
- GitHub Actions
- GoReleaser
- Flutter build pipelines
- Terraform plan/apply automation

**Pipeline Flow:**

```
PR opened → lint + tests
     ↓
Merge to main → build + package
     ↓
Infra diff (Terraform plan)
     ↓
Deploy serverless functions
     ↓
Mobile app builds (internal distribution)
```

---

## 4. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Mobile App (Flutter)                         │
│                  iOS / Android Native                           │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway                                 │
│              (Rate Limiting, Auth, Routing)                     │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Go Serverless Functions                         │
│         (Auth, Booking, Payment, Provider, Location)           │
└───────┬─────────────────┬─────────────────┬─────────────────────┘
        │                 │                 │
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────────────────┐
│  PostgreSQL   │ │  Redis Cache  │ │       Event Bus           │
│   (Primary)   │ │   (Sessions,  │ │    (EventBridge/PubSub)   │
│               │ │    Listings)  │ │                           │
└───────────────┘ └───────────────┘ └─────────┬─────────────────┘
                                              │
                          ┌───────────────────┼───────────────────┐
                          ▼                   ▼                   ▼
                   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
                   │Notifications│     │  Payments   │     │  Analytics  │
                   │   Handler   │     │   Handler   │     │   Handler   │
                   └─────────────┘     └─────────────┘     └─────────────┘
```

### Event-Driven Design

| Event | Trigger | Handlers |
|-------|---------|----------|
| `AppointmentCreated` | User books slot | Notification, Provider alert |
| `PaymentInitiated` | Checkout started | Payment gateway |
| `PaymentConfirmed` | Webhook received | Booking confirmation |
| `ProviderAssigned` | Match algorithm | User notification |
| `AppointmentCompleted` | Session ends | Review prompt, analytics |
| `PayoutReleased` | Settlement cycle | Provider notification |

> [!IMPORTANT]
> All side effects are handled asynchronously via the event bus.

---

## 5. Repository Structure

```
.
├── apps/
│   └── mobile/                   # Flutter app
│       ├── lib/
│       │   ├── core/             # Core utilities, constants
│       │   ├── features/         # Feature modules
│       │   │   ├── auth/
│       │   │   ├── booking/
│       │   │   ├── discovery/
│       │   │   ├── payments/
│       │   │   └── profile/
│       │   ├── shared/           # Shared widgets, services
│       │   └── main.dart
│       ├── android/
│       ├── ios/
│       └── pubspec.yaml
│
├── services/
│   ├── auth/                     # Auth functions (Go)
│   │   ├── cmd/
│   │   ├── internal/
│   │   └── go.mod
│   ├── booking/                  # Appointments
│   ├── payment/                  # Razorpay integration
│   ├── provider/                 # Doctors / clinics
│   ├── location/                 # Geo & proximity
│   └── notifications/            # Async handlers
│
├── infra/
│   ├── terraform/
│   │   ├── modules/
│   │   │   ├── api-gateway/
│   │   │   ├── lambda/
│   │   │   ├── database/
│   │   │   ├── cache/
│   │   │   └── events/
│   │   └── environments/
│   │       ├── dev/
│   │       ├── staging/
│   │       └── prod/
│   └── scripts/
│
├── scripts/
│   ├── local-dev.sh
│   ├── deploy.sh
│   └── seed-data.sh
│
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── cd-backend.yml
│       ├── cd-mobile.yml
│       └── infra.yml
│
├── docs/
│   ├── architecture.md
│   ├── api-reference.md
│   ├── deployment.md
│   └── onboarding.md
│
├── docker-compose.yml            # Local development
├── Makefile
└── README.md
```

---

## 6. Caching & Scalability Strategy

### Read-Heavy Optimization

| Data | Cache Layer | TTL | Invalidation |
|------|-------------|-----|--------------|
| Provider listings | Redis | 5 min | Event-driven |
| Availability snapshots | Redis | 1 min | On booking |
| Geo-based queries | Redis Geo | 5 min | Batch update |
| Session tokens | Redis | 15 min | On logout |

### Write Path

```
Client → API Gateway → Lambda → PostgreSQL (Strong Consistency)
                                      ↓
                               Event Published
                                      ↓
                              Cache Invalidated
```

### Rate Limiting

| Layer | Implementation |
|-------|----------------|
| API Gateway | Built-in throttling (requests/sec) |
| Application | Redis-based token buckets |
| Per-user | Sliding window counters |

---

## 7. Security Considerations

| Concern | Implementation |
|---------|----------------|
| Authentication | JWT with short TTL (15 min) |
| Session | Refresh tokens (7 days, rotating) |
| Secrets | Cloud secret manager (never in code) |
| Webhooks | Signature verification (HMAC) |
| Authorization | Role-based access control (RBAC) |
| Data at rest | AES-256 encryption |
| Data in transit | TLS 1.3 everywhere |
| Input validation | Schema-based validation |

> [!CAUTION]
> Never commit secrets to version control. Use `.env.example` as template.

---

## 8. Local Development

### Prerequisites

| Tool | Version |
|------|---------|
| Go | ≥1.21 |
| Flutter | ≥3.x |
| Docker | Latest |
| Terraform | ≥1.5 |
| Make | Any |
| Node.js | ≥18 (for tooling) |

### Quick Start

```bash
# Clone repository
git clone https://github.com/your-org/healthflow.git
cd healthflow

# Start infrastructure (Postgres, Redis)
docker-compose up -d

# Start backend services (local mode)
make local-backend

# In another terminal, start mobile app
cd apps/mobile
flutter pub get
flutter run
```

### Environment Setup

```bash
# Copy environment template
cp .env.example .env

# Edit with your local values
# Required: DATABASE_URL, REDIS_URL, RAZORPAY_KEY_ID, RAZORPAY_KEY_SECRET
```

### Useful Commands

```bash
make test           # Run all tests
make lint           # Run linters
make build          # Build all services
make migrate        # Run database migrations
make seed           # Seed development data
make clean          # Clean build artifacts
```

---

## 9. MVP Scope (2-Week Delivery)

### ✅ Included in MVP

| Feature | Priority | Status |
|---------|----------|--------|
| User auth (OTP) | P0 | 🔲 |
| Provider listing | P0 | 🔲 |
| Map-based discovery | P0 | 🔲 |
| Appointment booking | P0 | 🔲 |
| Razorpay payment | P0 | 🔲 |
| Basic provider dashboard | P1 | 🔲 |
| Notifications (push) | P1 | 🔲 |
| CI/CD pipelines | P0 | 🔲 |
| IaC for dev environment | P0 | 🔲 |

### ⏸️ Deferred (Post-MVP)

- Reviews & ratings
- Advanced analytics
- Multi-language support
- Offline mode
- Insurance integration
- Video consultations

---

## 10. Delivery Timeline

### Week 1

| Day | Milestone |
|-----|-----------|
| 1-2 | Infra bootstrap (Terraform, CI/CD skeleton) |
| 2-3 | Auth service + user model |
| 3-4 | Provider discovery service |
| 4-5 | Maps + location integration |
| 5-7 | Core Flutter UI (auth, discovery, booking) |

### Week 2

| Day | Milestone |
|-----|-----------|
| 8-9 | Booking flow (end-to-end) |
| 9-10 | Payments + webhooks |
| 10-11 | Notifications service |
| 11-12 | CI/CD hardening, testing |
| 13-14 | MVP release, documentation |

---

## 11. Non-Goals (Explicitly Out of Scope)

> [!WARNING]
> The following are intentionally excluded from this architecture:

- ❌ Monolithic backend
- ❌ Synchronous payment confirmation
- ❌ Stateful servers
- ❌ Heavy admin UI (MVP uses minimal dashboard)
- ❌ Manual deployments
- ❌ Vendor lock-in (cloud-agnostic where possible)

---

## 12. API Reference

### Authentication

```http
POST /auth/otp/send
POST /auth/otp/verify
POST /auth/refresh
POST /auth/logout
```

### Providers

```http
GET  /providers
GET  /providers/:id
GET  /providers/nearby?lat=&lng=&radius=
GET  /providers/:id/availability
```

### Bookings

```http
POST /bookings
GET  /bookings
GET  /bookings/:id
PUT  /bookings/:id/cancel
```

### Payments

```http
POST /payments/initiate
POST /payments/webhook
GET  /payments/:id/status
```

> Full API documentation: [docs/api-reference.md](docs/api-reference.md)

---

## 13. Contribution Guidelines

### Code Standards

| Aspect | Rule |
|--------|------|
| PR Size | Small, focused PRs only |
| Feature Flags | Required for risky changes |
| Tests | Mandatory for backend services |
| Infra Changes | Require `terraform plan` output |
| Commits | Conventional commits format |

### Branch Strategy

```
main          ← Production-ready
  └── develop ← Integration branch
        └── feature/* ← Feature branches
        └── fix/*     ← Bug fixes
```

### Review Process

1. Create feature branch from `develop`
2. Open PR with description
3. Pass CI checks (lint, test, build)
4. Get 1+ approval
5. Squash merge to `develop`

---

## 14. Monitoring & Observability

| Component | Tool |
|-----------|------|
| Logs | CloudWatch / Datadog |
| Metrics | Prometheus / CloudWatch Metrics |
| Traces | AWS X-Ray / Jaeger |
| Alerts | PagerDuty / Opsgenie |
| Dashboards | Grafana / CloudWatch Dashboards |

---

## 15. License

MIT License - see [LICENSE](LICENSE) for details.

---

## 16. Final Note

> [!TIP]
> This project is designed to scale from MVP to national deployment without architectural rewrites. The initial constraints are intentional. If something feels "over-engineered," it is because **operational debt is more expensive than code**.

---

## Quick Links

- 📖 [Architecture Deep Dive](docs/architecture.md)
- 🚀 [Deployment Guide](docs/deployment.md)
- 🔧 [API Reference](docs/api-reference.md)
- 👋 [Onboarding Guide](docs/onboarding.md)

---

**Built with ❤️ for better healthcare access**
