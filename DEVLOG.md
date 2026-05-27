# ElimuLocal — Developer Log

This file documents the building of **ElimuLocal** step by step. It records what was built, why decisions were made, and what was learned along the way.

---

## Step 1 — Project Setup
**Date:** March 2025  
**Branch:** `feature/initial-setup`

### What we built
- Created the project folder structure.
- Initialized Go module with `go mod init elimulocal`.
- Wrote the first `main.go` — a minimal server that responds with "ElimuLocal is working!".
- Set up Git and pushed to GitHub.

### Folder structure created
```
elimulocal/
├── main.go
├── go.mod
├── README.md
├── templates/
├── static/css/
├── static/fonts/
└── uploads/
```

### Go concepts introduced
- `package main` — Entry point of every Go program.
- `import` — Bringing in built-in packages.
- **Functions** — Defining behavior like `func homeHandler(...)`.
- `http.HandleFunc` — Connecting URLs to functions.
- `http.ListenAndServe` — Starting the web server.

---

## Step 2 — HTML Templates and CSS
**Date:** March 2025  
**Branch:** `feature/initial-setup`

### What we built
- Created three HTML templates: `base.html`, `home.html`, `upload.html`.
- Wrote `style.css` with full styling for the homepage and upload form.
- Updated `main.go` to render templates with dynamic data.
- Added `Resource` and `PageData` structs.
- Added search and university filter functionality.

### Key decisions
- **Standard Library:** Used Go's `html/template` package instead of a third-party library to keep dependencies at zero.
- **Template Composition:** Used `{{define "base"}}` and `{{template "content" .}}` pattern so navbar and footer are written once and shared.
- **University Autocomplete:** Used `<datalist>` for university suggestions instead of a hardcoded `<select>` so the list grows automatically.
- **HTTP Methods:** Used `GET` for search (bookmarkable URLs) and `POST` for uploads (data modification).

### Go concepts introduced
- **Structs** — Data modeling with `type Resource struct`.
- **Slices** — Dynamic arrays like `[]Resource`.
- **Maps** — Key-value pairs for deduplication in `getUniversities()`.
- **Range loops** — Iterating over collections.
- **Strings package** — Operations like `strings.Contains` and `strings.ToLower`.
- **Template rendering** — Using `ParseFiles` and `ExecuteTemplate`.

---

## Step 3 — SQLite Database and File Uploads
**Date:** April 2025  
**Branch:** `feature/initial-setup`

### What we built
- Added SQLite database using the `modernc.org/sqlite` (pure Go) package.
- Created `resources` table to persist data across restarts.
- Implemented real PDF file uploads saved to the `uploads/` folder.
- Added file download functionality with a download counter.
- Automatic database seeding on first run.

### SQL commands learned
- `CREATE TABLE IF NOT EXISTS` — Safe table initialization.
- `INSERT INTO ... VALUES` — Adding new records.
- `SELECT ... WHERE ... LIKE` — Filtering data.
- `UPDATE ... SET` — Modifying existing rows.

### Key decisions
- **Pure Go Driver:** Used `modernc.org/sqlite` to avoid needing a C compiler (CGO_ENABLED=0).
- **Security:** Used `?` placeholders in SQL queries to prevent SQL injection.
- **Unique Filenames:** Used `time.Now().UnixNano()` to prevent file name collisions.
- **Efficiency:** Used `io.Copy()` for file saving to handle large files without high memory usage.

---

## Step 4 — Local Fonts
**Date:** April 2025  
**Branch:** `feature/initial-setup`

### What we built
- Downloaded **Sora** and **Space Mono** font files.
- Added `@font-face` rules to `style.css` to load fonts from local disk.
- Removed Google Fonts dependency entirely.

### Why this matters
- **Offline Capability:** The app now works on a campus LAN with zero internet connection.
- **Performance:** Faster load times on slow networks as no external requests are made.

---

## Month 1 — COMPLETE ✅
**Status:** Architecture and Core Features Stable

- [x] Project setup and folder structure
- [x] HTML templates and CSS styling
- [x] Dynamic university list
- [x] SQLite database integration
- [x] Real file uploads and downloads
- [x] Local fonts (internet-independent)

---

## Step 7 — Improved Search and Filtering
**Date:** April 2025  
**Branch:** `feature/month2-improvements`

### What we built
- Added category filter dropdown (Notes, Past Paper, Textbook, etc.).
- Added sort options (Newest, Most Downloaded, Most Helpful).
- Persistent filters after form submission.
- Added "Clear filters" link and a results counter badge.

### Go concepts introduced
- **Switch statements** — Cleaner logic for sort ordering.
- **Dynamic SQL** — Building queries conditionally using `[]interface{}` for arguments.

---

## Step 8 — Ratings and Upvotes
**Date:** April 2025  
**Branch:** `feature/month2-improvements`

### What we built
- Added `upvotes` column to the `resources` table.
- Implemented a "Helpful" button that increments the vote count via `POST /upvote/{id}`.
- Smart redirects using the `Referer` header to keep users on their current search page.

---

## Step 10 — Success and Error Flash Messages
**Date:** April 2025  
**Branch:** `feature/month2-improvements`

### What we built
- Added a success flash message that appears after an upload.
- Used URL query parameters (`?success=1`) for simple state signaling.
- Implemented a JavaScript fade-out effect for the message div.

---

## Step 11 — Campus LAN Deployment
**Date:** April 2025  
**Branch:** `feature/month2-improvements`

### What we did
- Identified the server's local IP address.
- Verified firewall settings.
- Successfully accessed ElimuLocal on mobile devices via campus WiFi.

---

## Step 12 — Mobile Navbar Fixes
**Date:** April 2025  
**Branch:** `feature/month2-improvements`

### What we did
- Patched CSS to prevent branding text from breaking on small screens.
- Implemented flex-wrapping for the navbar so buttons stack appropriately on mobile.

---

## Month 2 — COMPLETE ✅
**Status:** UX Improvements and Real-world Testing

- [x] Enhanced search with categories and sorting
- [x] Upvote/Rating system
- [x] Success flash messages
- [x] Mobile responsiveness fixes
- [x] Verified deployment on campus LAN

---

## Step 15 — User Authentication
**Date:** May 2026  
**Branch:** `feature/month3`

### What we built
- Created `auth.go` to isolate authentication logic.
- Implemented User Registration, Login, and Logout.
- Added session management using `github.com/gorilla/sessions`.
- Dynamic navbar that shows the current user and session state.
- Password hashing with `bcrypt`.

### Challenges & Solutions
- **Compilation:** Switched from `go run main.go` to `go run .` to compile multiple files.
- **Session Decoding:** Registered `int` type with `encoding/gob` to fix silent session failures.
- **Inconsistent Navbar:** Created `newPageData()` helper to ensure session state is passed to all templates.

### SQL Reference
- **Users Table:** Created with `UNIQUE` constraints on username and email.
- **Password Safety:** Only `password_hash` is stored; raw passwords never touch the DB.

---

## Step 18 — Video Upload and Playback
**Date:** May 2026  
**Branch:** `feature/month3`

### What we built
- Increased upload limit from **10MB to 500MB** for lecture videos.
- Added support for `.mp4`, `.mkv`, and `.webm` files.
- Implemented an HTML5 video player in the resource preview page.
- Increased server `ReadTimeout` and `WriteTimeout` to 5 minutes to support slow uploads.

---

## Step 19 — Animated Landing Page
**Date:** May 2026  
**Branch:** `feature/month3`

### What we built
- Created a standalone, premium landing page at `/`.
- Integrated **GSAP** for a high-end "letter-splitting" reveal animation on load.
- Added sections for "Mission Statement", "Problems & Solutions", and "Mission Vision".
- Moved the main resource listing to `/browse`.
- Dynamic Call-to-Action (CTA) buttons that change based on login status.

### Design Decisions
- **Standalone Template:** Used a separate template for the landing page to allow for a unique visual style (dark mode, custom animations) without bloating `base.html`.
- **GSAP:** Chosen for its ability to coordinate complex multi-element animations with high performance.

---

## Step 20 — Environment Configuration & Dockerization
**Date:** May 2026  
**Branch:** `feature/deployment-ready`

### What we built
- Standardized environment variables using a `.env` file.
- Added support for configuring Port, Database path, Session secrets, and Upload limits via environment.
- Created a `Dockerfile` for containerized deployment.
- Added `docker-compose.yml` for easy local development setup.

### Benefits
- **Security:** Secrets are no longer hardcoded in the source.
- **Portability:** The app can now be run anywhere with `docker compose up`.

---

## Step 21 — Database Migration System
**Date:** May 2026  
**Branch:** `feature/migrations`

### What we built
- Created a robust migration runner in `migrate.go`.
- Tracks applied migrations in a `schema_migrations` table.
- Migrations are stored as plain SQL files in the `migrations/` directory.
- Ensuring migrations run in a transaction to prevent partial state on failure.

### Current Migrations
1. `001_create_users.sql` — Initial user table schema.
2. `002_create_resources.sql` — Main resources table with audit fields.
3. `003_backfill_user_id.sql` — Linking existing resources to user accounts.

---

## Month 3 — IN PROGRESS 🏗️
**Focus:** Scalability and Production Readiness

- [x] User Authentication system
- [x] Video support (500MB)
- [x] Premium Landing Page
- [x] Environment configuration
- [x] Database migrations
- [ ] Multi-tenant support (University separation)
- [ ] Cloud storage integration (Turso/B2)
