# Pull Request: Phase 1 - Core Foundation, Universal Item Catalog & Personal Taste Journal (Shelf & Queue)

## 📌 PR Title
`feat(core): implement phase 1 core foundation, item catalog, and personal shelf/queue`

## 🏷️ PR Type
- [x] ✨ Feature
- [x] 🏗️ Architectural Foundation
- [x] 🧪 Automated Tests
- [x] 📚 Documentation

---

## 📖 Summary of Changes

This Pull Request establishes the production-grade architectural foundation for **Recco** (Phase 1) following Hexagonal/Clean Architecture and Domain-Driven Design principles in Go:

### 1. Domain Modeling (`internal/domain`)
- **Universal Item Model (`Item`)**: Multi-platform media abstraction supporting `movie`, `youtube`, `book`, `anime`, and `podcast` with extensible `metadata` JSONB.
- **Personal Taste Journal (`ShelfEntry`)**: State machine managing item lifecycles (`queued` -> `finished` / `dropped`).
- **Intentionality Constraints**: Enforced 1–10 rating bounds, 280-character limit on personal reflections ("take"), and strict state transition rules.
- **Sentinel Errors**: Typed errors (`ErrNotFound`, `ErrInvalidInput`, `ErrConflict`, `ErrUnauthorized`).

### 2. Relational Schema & Database Migrations (`migrations/`)
- `000001_create_users_table`: Identity table with indexing on `username` and `email`.
- `000002_create_items_table`: Universal catalog with GIN indexing on `metadata`, source provider IDs, and external URLs.
- `000003_create_shelf_entries_table`: Queue and shelf persistence with uniqueness constraint on `(user_id, item_id)` and check constraints on ratings.

### 3. High-Performance Persistence Layer (`internal/repository`)
- **PostgreSQL Adapter (`internal/repository/postgres`)**:
  - Connection pooling with `jackc/pgx/v5/pgxpool`.
  - Idempotent upserts for queue and shelf entries (`ON CONFLICT (user_id, item_id) DO UPDATE ...`).
  - Single-query JOIN queries hydrating item details for fast shelf responses.
- **In-Memory Thread-Safe Adapter (`internal/repository/memory`)**:
  - `sync.RWMutex`-protected implementation for sub-millisecond unit/integration testing and dev fallback.

### 4. Application & Domain Services (`internal/service`)
- `ItemService`: URL and Source ID deduplication logic, catalog registration, and filtered listing.
- `ShelfService`: Personal queue placement, completion workflow (with rating & reflection validation), dropping, and taste journal pagination.

### 5. Transport & HTTP Endpoints (`internal/transport/http`)
- Method-based routing using modern Go 1.22+ `http.NewServeMux()`.
- Production middleware pipeline:
  - `RequestLogger`: Injects unique `X-Request-ID` and logs latency, status, and path using `log/slog`.
  - `Recovery`: Catches panics and safely returns standard HTTP 500 JSON without crashing the server.
  - `CORS`: Standard cross-origin resource sharing headers.
- **Endpoints**:
  - `GET  /health` - Health probe
  - `POST /api/v1/items` - Register item
  - `GET  /api/v1/items` - List catalog items (with type filter and pagination)
  - `GET  /api/v1/items/{id}` - Retrieve item by ID
  - `POST /api/v1/shelf/queue` - Add to personal queue
  - `POST /api/v1/shelf/finish` - Finish item (move to rated taste journal)
  - `POST /api/v1/shelf/drop` - Mark item dropped
  - `GET  /api/v1/shelf/queue` - Fetch active queue (includes `queue_zero: bool` indicator)
  - `GET  /api/v1/shelf` - Fetch rated taste journal
  - `DELETE /api/v1/shelf/{id}` - Remove entry

### 6. Production API Entrypoint & Lifecycle (`cmd/api/main.go`)
- Graceful shutdown handling OS signals (`SIGINT`, `SIGTERM`) with 10s drain deadline.
- Automatic fallback to in-memory persistence in development when local PostgreSQL is not running.

---

## 🧪 Testing & Verification

- **Automated Test Suites**:
  - `internal/domain`: Domain validation tests for items and shelf entries.
  - `internal/service`: Unit tests verifying queue transitions, completion workflows, and URL deduplication.
  - `internal/transport/http/handler`: End-to-end HTTP integration tests via `net/http/httptest`.
- **Result**: `PASS` across all suites (zero failures, zero leaks).

```bash
=== RUN   TestItem_Validate
--- PASS: TestItem_Validate (0.00s)
=== RUN   TestShelfEntry_Validate
--- PASS: TestShelfEntry_Validate (0.00s)
=== RUN   TestItemService_CreateItem
--- PASS: TestItemService_CreateItem (0.00s)
=== RUN   TestShelfService_QueueAndFinishWorkflow
--- PASS: TestShelfService_QueueAndFinishWorkflow (0.00s)
=== RUN   TestItemHandler_HTTP
--- PASS: TestItemHandler_HTTP (0.02s)
=== RUN   TestShelfHandler_HTTP
--- PASS: TestShelfHandler_HTTP (0.01s)
PASS
ok  	github.com/recco-app/recco/internal/domain
ok  	github.com/recco-app/recco/internal/service
ok  	github.com/recco-app/recco/internal/transport/http/handler
```

---

## 🚀 How to Test Locally

1. **Run the server**:
   ```bash
   go run ./cmd/api/main.go
   ```
2. **Health Check**:
   ```bash
   curl http://localhost:8080/health
   ```
3. **Run All Tests**:
   ```bash
   go test -v ./...
   ```
