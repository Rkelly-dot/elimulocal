# 📚 ElimuLocal

> A premium learning resource hub for university students — built with **Go** and designed for **Campus LAN** environments.

ElimuLocal empowers university students to share and download study materials (notes, past papers, textbooks, and lecture videos) over a local network. **No internet required. No data costs.** Just a campus WiFi network and a community of learners.

---

## 🗺️ What this project is

ElimuLocal is a robust Go web application that solves a critical problem: the high cost of data and textbooks for university students. By running on a local server (like a student's laptop or a dedicated university machine), it creates a high-speed, zero-cost haven for educational resources organized by university, course, and category.

---

## 👥 Who does what — roles and workflow

### 🖥️ The Administrator
The person who hosts the server (often a student leader or lab attendant).
1.  **Launch:** Runs the platform via `docker compose up` or `go run .`.
2.  **Broadcast:** Shares the local IP address (e.g., `http://192.168.1.100:8080`) via WhatsApp groups or campus notices.
3.  **Manage:** Oversees the server health and storage.

### 👩‍🏫 Registered Students & Teachers
- **Contribute:** Upload PDFs or Video Lectures (up to 500MB).
- **Vote:** Mark resources as "Helpful" to highlight quality materials.
- **Track:** View their own contributions and download counts.

### 👨‍🎓 Visitors
- **Browse:** Search and filter everything without needing an account.
- **Download:** Save any material directly to their device.

---

## 🚀 How to run it

### Option 1: Using Docker (Recommended)
The easiest way to get started. All dependencies are pre-configured.
```bash
# 1. Clone the repository
git clone https://github.com/Rkelly-dot/elimulocal.git && cd elimulocal

# 2. Setup environment (optional, defaults provided)
cp .env.example .env

# 3. Start the server
docker compose up -d
```

### Option 2: Running with Go
1.  **Install dependencies:** Ensure Go 1.21+ is installed.
2.  **Configuration:** Create a `.env` file based on `.env.example`.
3.  **Launch:**
    ```bash
    go run .
    ```

---

## 📁 Project Structure

```
elimulocal/
├── main.go            # Entry point & core handlers
├── auth.go            # Authentication & Session logic
├── storage.go         # Database interactions & File handling
├── migrate.go         # SQL Migration runner
├── .env               # Environment configuration
│
├── migrations/        # SQL schema versioning
├── templates/         # HTML page templates (Go html/template)
├── static/            # Static assets (CSS, local fonts, JS)
├── uploads/           # Storage for PDF and Video files
└── elimulocal.db      # SQLite database file
```

---

## 🧠 Tech Stack & Concepts

- **Backend:** Go (Golang) — Minimal, fast, and type-safe.
- **Frontend:** Pure HTML/CSS with **GSAP** for premium animations.
- **Database:** SQLite with a custom Migration System.
- **Auth:** `bcrypt` for hashing and `gorilla/sessions` for session management.
- **Offline:** No external CDNs. All fonts and icons are served locally.

---

## ✨ Features

- [x] **Zero-Data Transfers:** Share massive video files instantly over WiFi.
- [x] **Smart Search:** Filter by category, university, or course.
- [x] **Video Support:** Native HTML5 player for `.mp4`, `.mkv`, and `.webm`.
- [x] **User Auth:** Secure student accounts with private profiles.
- [x] **Success Badges:** Real-time feedback for uploads and actions.
- [x] **Premium UI:** Animated landing page and responsive mobile design.

---

## 🗓️ Build Roadmap

### Month 1 — The Foundation ✅
- Core structure, SQLite integration, and local font embedding.
- Basic search/filter functionality.

### Month 2 — Community & UX ✅
- Upvote system and download counters.
- Mobile responsiveness fixes and campus LAN testing.
- Functional success/error flash messages.

### Month 3 — Scale & Polish 🏗️
- [x] **User Authentication:** Registration and secure login.
- [x] **Video Support:** Handling large uploads (500MB+).
- [x] **Landing Page:** Premium animated first impression.
- [x] **Migrations:** Robust schema versioning.
- [x] **Dockerization:** Containerized deployment for stability.
- [ ] **Cloud Hybrid:** (Optional) Syncing with Turso/B2 for backups.
- [ ] **Multi-Tenant:** Better university isolation and management.

---

## 🌍 Design Decisions

- **Why Go Standard Library?** To minimize dependencies and ensure the project remains lightweight and easy to understand for beginners.
- **Why 500MB Limits?** Balances the need for high-quality lecture videos with the reality of local server disk space.
- **Why GSAP?** To provide a "premium" feel that inspires confidence in students, making it feel like a professional tool rather than a student project.
- **Why CGO-free SQLite?** Ensures that the project can be cross-compiled for any OS (Windows/Mac/Linux) without needing a local C compiler.

---

## 🤝 Contributing

We welcome contributions from fellow students! Whether it's fixing a CSS bug, adding a new feature, or simply testing on your local campus network.

*Built for university students, by university students. Share knowledge, uplift each other.
