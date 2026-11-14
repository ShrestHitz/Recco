# Pull Request: Phase 1 Frontend — Anti-Distraction UI, Finite Queue & Personal Taste Journal

## 📌 PR Title
`feat(web): implement phase 1 frontend — finite queue, queue zero, and taste journal`

## 🏷️ PR Type
- [x] ✨ Feature
- [x] 🎨 User Interface & Experience
- [x] 🔌 API Integration
- [x] 📚 Developer Tooling & Documentation

---

## 📖 Summary of Changes

This Pull Request delivers the complete **Phase 1 Web Frontend** for **Recco**, built from first principles with **React, TypeScript, Tailwind CSS, Lucide Icons, and Vite**.

The design directly reflects Recco's core philosophy: **"Recommended by people, not algorithms. No feed. Just friends."**

### 1. Intentional & Anti-Distraction User Experience
- **Finite Queue**:
  - Displays recommendations accepted by the user ("want to watch/read/listen").
  - Each item card shows rich metadata, creator attribution, cover art, and media category badges (YouTube, Movie/TV, Book, Anime, Podcast).
- **Queue Zero Celebration Screen**:
  - When the queue is empty, the interface actively protects user focus rather than serving algorithmic infinite feeds:
    > *"You've reached Queue Zero. No algorithmic feed pushing endless content. You’ve consumed what you planned. Close the tab, go for a walk, or read something offline."*
- **Personal Taste Journal (The Shelf)**:
  - Permanent record of consumed culture.
  - Displays 1–10 personal ratings and the user's authentic one-line take ("reflection").
  - Category filter pills: *All Media, Movies & TV, YouTube, Books, Anime, Podcasts*.
- **External Links Only (No Autoplay Traps)**:
  - Links open directly in source applications (YouTube, TMDB, Goodreads/Open Library, AniList) to prevent video autoplay rabbit holes.

### 2. Modals & Interactive Workflows
- **`AddItemModal` (Save Recommendation)**:
  - Smart category selector (YouTube, Movie, Book, Anime, Podcast).
  - Auto-detection helper: entering a YouTube or TMDB link automatically classifies the media type.
  - Destination routing: choose between adding to **Queue** ("want to consume") or directly recording to **Shelf** ("already finished").
- **`FinishModal` (Record to Journal)**:
  - Interactive 1–10 score selector with qualitative descriptors (e.g. `10: Masterpiece`, `8: Great`, `4: Subpar`).
  - One-line reflection input with a 280-character maximum counter.

### 3. Architecture & Developer Experience
- **Vite + React + TypeScript + Tailwind CSS**:
  - Hot Module Replacement (HMR) during local development on `http://localhost:3000`.
  - Vite API proxy forwarding `/api` and `/health` to the Go backend at `http://localhost:8080`.
- **Integrated Go Single-Binary Serving**:
  - When built (`npm run build`), the Go HTTP server (`internal/transport/http/router.go`) automatically serves the static production SPA from `web/dist/` at `http://localhost:8080/`, complete with SPA fallback routing.
- **Resilient Fallback**:
  - Connection health indicator shows whether the Go backend is online.
  - One-click **Starter Pack** seeder to immediately test the anti-distraction workflow.

---

## 📁 Key File Additions

| File | Purpose |
|---|---|
| `web/src/App.tsx` | Main application coordinator, state management, and toast notifications |
| `web/src/components/Header.tsx` | Brand header, Go API connectivity indicator, and Queue/Shelf switcher |
| `web/src/components/QueueView.tsx` | Personal queue view and the celebratory Queue Zero state |
| `web/src/components/ShelfView.tsx` | Taste journal grid with media filters and personal reflection quotes |
| `web/src/components/AddItemModal.tsx` | Cultural item registration modal with smart type auto-detection |
| `web/src/components/FinishModal.tsx` | Rating (1-10) and one-line take completion dialog |
| `web/src/components/MediaBadge.tsx` | Distinctive visual badges for Movies, YouTube, Books, Anime, and Podcasts |
| `web/src/api/client.ts` | Type-safe fetch client for the Recco Go backend |
| `internal/transport/http/router.go` | Added static Single-Page App serving with fallback to `index.html` |
| `Makefile` | Added `web-install`, `web-dev`, and `web-build` targets |

---

## 🧪 Verification & Build Results

1. **Frontend Build (`npm run build`)**:
   - TypeScript check: `PASS` (0 errors)
   - Vite bundle: `dist/index.html` (1.19 kB), CSS (22.62 kB), JS (184.04 kB)
2. **Backend Tests (`go test -v ./...`)**:
   - All domain, service, and HTTP handler test suites pass with 0 regressions.

---

## 🚀 How to Run Locally

### Option A: Fullstack Dev Server (with Vite HMR)
```bash
# Terminal 1: Start Go Backend
make run

# Terminal 2: Start Vite Frontend (proxying to :8080)
cd web && npm run dev
# Open http://localhost:3000
```

### Option B: Standalone Go Server (Serving Built Frontend)
```bash
# 1. Build frontend bundle
cd web && npm run build

# 2. Run Go server
go run ./cmd/api/main.go
# Open http://localhost:8080
```
