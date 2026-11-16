# Pull Request: Phase 2 — Concurrent Metadata Enrichment Pipeline & Singleflight Deduplication

## 📌 PR Title
`feat(enrichment): implement concurrent metadata pipeline, singleflight deduplication, and auto-enrichment`

## 🏷️ PR Type
- [x] ✨ Feature
- [x] ⚡ Concurrency & Performance
- [x] 🛡️ Resilience & Caching
- [x] 🎨 Frontend Integration
- [x] 🧪 Automated Tests

---

## 📖 Summary of Changes

This Pull Request delivers **Phase 2** of **Recco**: a high-throughput, concurrent **Metadata Enrichment Pipeline** that turns raw external URLs into rich cultural recommendations automatically, without manual link-formatting.

### 1. High-Performance Concurrency & Resilience
- **Singleflight Request Deduplication (`golang.org/x/sync/singleflight`)**:
  - Automatically coalesces simultaneous incoming requests for the same URL into a single outbound flight.
  - Prevents cache stampedes and eliminates redundant third-party API hits under concurrent spikes.
- **Thread-Safe Generic TTL Cache (`internal/platform/cache`)**:
  - In-memory generic cache (`TTLCache[K, V]`) with lazy eviction and configurable TTL.
  - Subsequent requests for previously extracted URLs resolve in sub-millisecond time.
- **Bounded Worker Pool (`internal/platform/workerpool`)**:
  - Fixed-size goroutine worker pool with buffered job queues and graceful shutdown.
  - Prevents runaway goroutines from exhausting system file descriptors or sockets during batch operations.
- **Strict Context Timeouts**:
  - All outbound network calls are bounded by 5–7s timeouts, preventing sluggish external providers from hanging backend worker threads.

### 2. Multi-Provider Metadata Extraction (`internal/platform/enricher`)
- **`YouTubeProvider`**:
  - Leverages YouTube's official oEmbed service (`https://www.youtube.com/oembed`) requiring zero API keys.
  - Extracts canonical video ID, channel/creator name, title, and high-resolution thumbnail (`hqdefault.jpg`).
  - Graceful offline fallback: parses video ID and constructs valid thumbnails even if YouTube's oEmbed endpoint is unreachable.
- **`OpenLibraryProvider`**:
  - Extracts book metadata from `openlibrary.org` work/edition/ISBN URLs.
  - Fetches title, authors, publication details, and full-resolution cover art (`covers.openlibrary.org/b/id/...-L.jpg`).
- **`OpenGraphProvider`**:
  - Generic HTML parser powered by streaming tokenization (`golang.org/x/net/html`).
  - Limits read size to 512KB to protect backend memory from denial-of-service payloads.
  - Parses OpenGraph tags (`og:title`, `og:image`, `og:description`, `og:site_name`, `author`).
  - Intelligently classifies media types based on domain signatures:
    - TMDB, IMDb, Letterboxd ➔ `movie`
    - AniList, MyAnimeList ➔ `anime`
    - Goodreads, Open Library ➔ `book`
    - Spotify, Apple Podcasts ➔ `podcast`

### 3. Application & Transport Integration
- **Transparent Ingestion Enrichment (`ItemService.CreateItem`)**:
  - If a user provides an `external_url` with missing title, cover art, or category, `ItemService` automatically enriches and fills the fields on the fly.
- **Preview Endpoint (`GET /api/v1/items/enrich?url=...`)**:
  - Fast endpoint enabling client-side instant preview and category auto-detection.

### 4. Interactive Frontend Integration (`web/`)
- **Smart Auto-Enrichment in [AddItemModal.tsx](file:///c:/Users/Shresth/Desktop/Recco/web/src/components/AddItemModal.tsx)**:
  - When a user pastes a URL (YouTube, TMDB, Open Library, AniList, etc.), an inline spinner announces: *"Extracting metadata..."*.
  - Automatically auto-populates Title, Creator, Media Category, and Cover Art URL.
  - Renders a live **Cover Art Preview Card** directly inside the modal before saving.
  - Users can immediately see their cultural recommendation with zero manual typing required.

---

## 🧪 Testing & Concurrency Verification

All test suites pass cleanly across domain, platform, service, and HTTP layers:

1. **Singleflight Deduplication Test (`TestEnrichmentService_SingleflightDeduplication`)**:
   - Dispatches 30 simultaneous goroutines against a mock slow network provider for the identical URL.
   - Verified that exactly 1 outbound HTTP call was made, and all 30 goroutines received the identical enriched response.
2. **Bounded Concurrency Test (`TestPool_BoundedConcurrency`)**:
   - Submits 15 burst tasks to a 3-worker pool; verifies peak concurrent execution never exceeds worker limit.
3. **TTL Cache Test (`TestTTLCache_ConcurrentAccess`)**:
   - 100 concurrent reader/writer goroutines verifying thread-safety and race condition prevention.
4. **Provider Tests (`TestYouTubeProvider_CanHandleAndExtract`, `TestOpenGraphProvider_Extraction`)**:
   - Verified URL detection, fallback mechanisms, and HTML OpenGraph extraction.
5. **End-to-End HTTP Test (`TestItemHandler_HTTP`)**:
   - Verified `GET /api/v1/items/enrich?url=...` endpoint.

---

## 🚀 How to Test Locally

1. **Start Go Backend**:
   ```bash
   go run ./cmd/api/main.go
   ```
2. **Start Frontend**:
   ```bash
   cd web && npm run dev
   ```
3. **Try It Out**:
   - Click **"+ Save Item"**.
   - Paste any YouTube link (e.g. `https://www.youtube.com/watch?v=094y1Z2wpJg`) or Open Library URL.
   - Watch the title, creator, category, and cover art auto-populate instantly!
