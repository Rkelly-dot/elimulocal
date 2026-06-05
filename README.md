#  ElimuLocal

> Free study materials for every Kenyan university student.

ElimuLocal is a platform where university students share and access notes,
past papers, textbooks and lecture videos — for free. It works on campus
WiFi with no internet required, and is available from anywhere via the
cloud at [elimulocal.onrender.com](https://elimulocal.onrender.com).

Built entirely in Go by Ryan Kelly — Zone01 Kisumu, Kenya.

---

##  What this project is

ElimuLocal solves three problems that every Kenyan university student faces:

- **Textbooks are expensive** — a single recommended textbook can cost more
  than a student's entire monthly budget
- **Notes are trapped in WhatsApp groups** — good materials exist but are
  buried in chat history, only accessible if you know the right person
- **Internet is unreliable** — most edtech solutions require stable internet
  which many Kenyan campuses cannot guarantee

ElimuLocal is the answer to all three — a free, offline-first platform where
students share knowledge with each other openly and without barriers.

---

##  Features

- Browse and search resources by keyword, university and category
- Upload PDF notes, past papers, textbooks, summaries and video lectures
- Preview PDFs in the browser using PDF.js — no download needed
- Stream video lectures directly in the browser
- Download any resource for offline reading
- Upvote helpful resources so the best materials rise to the top
- Download counter shows how popular each resource is
- User accounts — register, login, logout
- Edit and delete your own uploaded resources
- Works on campus WiFi with zero internet connection
- Works from anywhere via the cloud
- Animated landing page explaining the platform

---

##  Architecture

```
Students on campus WiFi          Students off campus
───────────────────────          ───────────────────
Campus LAN server                Cloud server (Koyeb)
     ↓                                   ↓
Go app (Docker)              Go app (Docker on Koyeb)
     ↓                                   ↓
Turso (cloud SQLite)         Turso (same database)
     ↓                                   ↓
Backblaze B2                 Backblaze B2 (same files)
(file storage)               (file storage)
```

Both the campus server and the cloud server share the same Turso database
and the same Backblaze B2 file storage. A resource uploaded on campus is
immediately available to students off campus and vice versa.

---

##  How to run it

### Requirements

- Go 1.25 or newer
- Docker and Docker Compose (for containerised deployment)
- A Turso account (free) for the database
- A Backblaze B2 account (free) for file storage

### Quick start with Docker

```bash
# 1. Clone the repository
git clone https://github.com/Rkelly-dot/elimulocal.git
cd elimulocal

# 2. Copy the example environment file and fill in your values
cp .env.example .env
nano .env

# 3. Start the server
docker compose up -d

# 4. Open your browser
# Go to: http://localhost:8080
```

### Run without Docker (development)

```bash
# 1. Clone the repository
git clone https://github.com/Rkelly-dot/elimulocal.git
cd elimulocal

# 2. Copy and fill in environment variables
cp .env.example .env
nano .env

# 3. Start the server
go run .

# 4. Open your browser
# Go to: http://localhost:8080
```

### Share on campus WiFi

```bash
# 1. Find your laptop's local IP address
ip addr show | grep "inet " | grep -v "127.0.0.1"    # Linux
ifconfig                                               # Mac
ipconfig                                               # Windows

# 2. Start the server
go run .

# 3. Share this link with students on the same WiFi:
# http://YOUR-IP:8080
```

---

##  Environment variables

Copy `.env.example` to `.env` and fill in your values.

| Variable | Description | Example |
|---|---|---|
| `PORT` | Server port | `8080` |
| `APP_ENV` | Environment | `development` or `production` |
| `SESSION_SECRET` | Long random secret for cookies | `your-secret-here` |
| `TURSO_URL` | Turso database URL | `libsql://db.turso.io` |
| `TURSO_TOKEN` | Turso auth token | `eyJhbGci...` |
| `B2_KEY_ID` | Backblaze B2 key ID | `your-key-id` |
| `B2_APP_KEY` | Backblaze B2 application key | `your-app-key` |
| `B2_BUCKET` | Backblaze B2 bucket name | `elimulocal-uploads` |
| `B2_ENDPOINT` | Backblaze B2 S3 endpoint | `s3.us-west-004.backblaze.com` |

---

##  Project structure

```
elimulocal/
├── main.go              ← HTTP handlers, routing, database queries
├── auth.go              ← Authentication, sessions, user management
├── storage.go           ← File storage (Backblaze B2 + local fallback)
├── migrate.go           ← Database migration runner
├── go.mod               ← Go module and dependencies
├── go.sum               ← Dependency checksums
├── Dockerfile           ← Multi-stage Docker build
├── docker-compose.yml   ← Container orchestration
├── .env.example         ← Environment variable template
├── README.md            ← You are here
├── DEVLOG.md            ← Developer log — every decision documented
│
├── migrations/          ← SQL migration files run in order
│   ├── 001_create_users.sql
│   ├── 002_create_resources.sql
│   └── 003_backfill_user_id.sql
│
├── templates/           ← HTML templates rendered by Go
│   ├── landing.html     ← Landing page (standalone, GSAP animated)
│   ├── base.html        ← Shared layout (navbar, footer)
│   ├── home.html        ← Browse and search resources
│   ├── upload.html      ← Upload a resource
│   ├── preview.html     ← In-app PDF and video viewer
│   ├── edit.html        ← Edit an existing resource
│   ├── login.html       ← Login page
│   └── register.html    ← Registration page
│
├── static/              ← Files served directly to the browser
│   ├── css/
│   │   └── style.css    ← All styling
│   ├── fonts/           ← Local font files (no internet needed)
│   └── pdfjs/           ← PDF.js viewer (not committed to Git)
│
└── uploads/             ← Local file fallback when B2 not configured
    └── .gitkeep
```

---

##  Who does what

### Administrator (person running the server)

**One-time setup:**
1. Clone the repository
2. Fill in `.env` with your Turso and Backblaze credentials
3. Run `docker compose up -d`
4. Find your local IP with `ip addr show`
5. Share `http://YOUR-IP:8080` with students on campus WiFi

**Every time you want ElimuLocal available on campus:**
1. Connect your machine to campus WiFi
2. Run `docker compose up -d`
3. The server runs in the background — no need to keep a terminal open

### Teachers and students

No installation needed. Just a browser.

1. Connect to campus WiFi (or use elimulocal.co.ke from anywhere)
2. Open any browser — Chrome, Firefox, Safari, anything
3. Go to `http://YOUR-CAMPUS-IP:8080` or `https://elimulocal.onrender.com`
4. Register for a free account
5. Browse, upload, download and upvote resources

---

##  Roadmap

### Month 1 — Core product ✓
- [x] Project setup and folder structure
- [x] HTML templates and CSS styling
- [x] SQLite database with file uploads and downloads
- [x] Search and filter by keyword, university and category
- [x] Local fonts — no internet dependency

### Month 2 — Real users ✓
- [x] Category filter and sort options
- [x] Upvote and helpful rating system
- [x] Download counter on resource cards
- [x] Success flash messages
- [x] Campus LAN deployment — tested on real Android phone
- [x] Mobile responsive navbar

### Month 3 — Polish and scale ✓
- [x] User authentication — register, login, logout
- [x] Edit and delete resources (owner only)
- [x] In-app PDF preview with PDF.js
- [x] Video upload and playback (MP4, MKV, WebM)
- [x] Video Lecture category
- [x] Animated landing page
- [x] Docker and docker-compose
- [x] Turso cloud database
- [x] Backblaze B2 file storage
- [x] Cloud deployment on Koyeb

### Month 4 — Hybrid architecture (planned)
- [ ] Campus sync service — resources shared between campus and cloud
- [ ] Auto-routing — campus LAN when on campus, cloud when off
- [ ] Premium campus tier with subdomain and analytics
- [ ] Multi-campus support
- [ ] M-Pesa billing integration

---

##  Design decisions

**Why Go?**
Go compiles to a single binary with no runtime dependencies. A campus can
run ElimuLocal by downloading one file and running it. It is also fast enough
to serve hundreds of concurrent students on modest hardware.

**Why SQLite via Turso?**
Turso gives us SQLite's simplicity with cloud persistence. All existing Go
code works unchanged — just a different connection string. No schema changes,
no ORM, no migration framework beyond our own simple runner.

**Why Backblaze B2 and not Cloudflare R2?**
R2 requires a credit card even for the free tier. B2 does not. Both are
S3-compatible so switching later requires changing only environment variables,
not code.

**Why Koyeb and not AWS or DigitalOcean?**
Koyeb's free tier requires no credit card. For a project in early growth
stage, zero infrastructure cost matters. The app can migrate to any Docker-
compatible host later with no code changes.

**Why offline-first?**
Many Kenyan university campuses have unreliable or expensive internet.
ElimuLocal was designed from day one to work on a campus LAN with no internet.
The cloud deployment is additive — it does not break the offline use case.

---

##  Tech stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Web framework | Go standard library (`net/http`) |
| Templates | Go `html/template` |
| Database | SQLite via Turso (libsql) |
| File storage | Backblaze B2 (S3 compatible) |
| Authentication | gorilla/sessions + bcrypt |
| PDF viewer | PDF.js (self-hosted) |
| Animations | GSAP 3 |
| Container | Docker (multi-stage build) |
| Hosting | Koyeb |
| DNS and SSL | Cloudflare |

---

## Learning resources

This project was built while learning Go. If you are doing the same:

- [A Tour of Go](https://go.dev/tour/) — official interactive tutorial
- [Go by Example](https://gobyexample.com/) — practical code snippets
- [Let's Go](https://lets-go.alexedwards.net/) — web apps in Go (book)

---

## Contributing

If you are a student at a Kenyan university and want to help test ElimuLocal
or contribute resources for your campus, reach out via GitHub.

---

*Built for Kenyan university students. Share knowledge, uplift each other. *

*Ryan Kelly — Zone01 Kisumu — 2025*
