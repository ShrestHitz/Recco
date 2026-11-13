# Recco 🛋️

> **Recommended by people, not algorithms. No feed. Just friends.**

Recco is a quiet, friends-only space to share and rate movies, music, books, podcasts, YouTube videos, and anime. Designed from first principles to cultivate focus and curb distraction:

- **No Infinite Feed**: Replace algorithmic doomscrolling with a **finite inbox** of recommendations sent to you. When it's empty, you've reached Inbox Zero—the app is done.
- **Queue & Shelf**: Save accepted recommendations to your personal **Queue** (*"want to consume"*). When finished, move them to your **Shelf** with a rating and a one-line take—a taste journal first, social network second.
- **Intentional Sharing**: Every recommendation requires a personal note explaining why you're sharing it. A weekly scarcity budget (e.g., 5 shares/week) keeps quality high and stops link dumps.
- **Small Circles (10–30 people)**: Share within themed circles of like-minded friends instead of broadcasting to a public audience. No follower counts, likes, or leaderboards.
- **Closed-Loop Feedback**: When a friend finishes something you recommended, you get a quiet note with their reaction and rating—real human connection instead of algorithmic dopamine.
- **External Links**: Metadata and cover art only, opening directly in the source app (no autoplay video rabbit holes).

---

## 🛠️ Complete Engineering Action Plan

The complete architectural specification, database schema, concurrency pipelines, and phased roadmap are detailed in:  
👉 **[ACTION_PLAN.md](ACTION_PLAN.md)**

---

## 🏗️ Backend Engineering Highlights (Golang)

Recco's backend demonstrates senior-level Go engineering practices:
- **Clean / Hexagonal Architecture**: Domain-driven design with decoupled domain rules, ports, and adapters.
- **Concurrency & Goroutine Pipelines**: Metadata enrichment worker pool with `singleflight` deduplication and context deadline propagation.
- **Atomic Scarcity & Rate Limiting**: Distributed weekly budget tracking via Redis atomic counters and token buckets.
- **Transactional Outbox Pattern**: Decoupled async job dispatch, background digest generation, and closed-loop notification relays with Dead Letter Queues (DLQ).
- **Dual Transport**: High-throughput REST API (Go 1.22+ `net/http`) + gRPC for inter-service communication.
- **Production Observability**: Full OpenTelemetry distributed tracing, Prometheus metrics instrumentation, and structured `log/slog`.
- **Hermetic Testing**: Unit tests with table tests and integration testing using `testcontainers-go` for PostgreSQL and Redis.

---

## 🚀 Quick Start

See **[ACTION_PLAN.md](ACTION_PLAN.md)** for the step-by-step implementation guide.
