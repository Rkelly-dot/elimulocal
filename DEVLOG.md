# ElimuLocal — Developer Log

This file documents the building of ElimuLocal step by step.
It records what was built, why decisions were made, and what was learned.

---

## Step 1 — Project setup
**Date:** March 2025
**Branch:** feature/initial-setup

### What we built
- Created the project folder structure
- Initialised Go module with `go mod init elimulocal`
- Wrote the first `main.go` — a minimal server that responds with
  "ElimuLocal is working!"
- Set up Git and pushed to GitHub

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
- `package main` — entry point of every Go program
- `import` — bringing in built-in packages
- Functions — `func homeHandler(...)`
- `http.HandleFunc` — connecting URLs to functions
- `http.ListenAndServe` — starting the web server

---

## Step 2 — HTML templates and CSS
**Date:** March 2025
**Branch:** feature/initial-setup

### What we built
- Created three HTML templates: `base.html`, `home.html`, `upload.html`
- Wrote `style.css` with full styling for the homepage and upload form
- Updated `main.go` to render templates with dynamic data
- Added `Resource` and `PageData` structs
- Added search and university filter functionality

### Key decisions
- Used Go's `html/template` package instead of a third party library —
  keeps dependencies at zero and is built into Go
- Used `{{define "base"}}` and `{{template "content" .}}` pattern so
  navbar and footer are written once and shared across all pages
- Used `<datalist>` for university autocomplete instead of a hardcoded
  `<select>` — universities grow automatically as students submit
- Used `method="GET"` for search so searches are bookmarkable URLs
- Used `method="POST"` for uploads because they change data

### Go concepts introduced
- Structs — `type Resource struct` and `type PageData struct`
- Slices — `[]Resource` for storing multiple resources
- Maps — `map[string]bool` for deduplication in `getUniversities()`
- Range loops — `for _, r := range resources`
- String operations — `strings.Contains`, `strings.EqualFold`,
  `strings.TrimSpace`, `strings.ToLower`
- Error handling — `if err != nil`
- HTTP methods — checking `r.Method == "GET"` vs `"POST"`
- Template rendering — `template.ParseFiles` and `ExecuteTemplate`
- Serving static files — `http.FileServer` and `http.StripPrefix`

### CSS concepts introduced
- CSS variables — `:root` and `var(--name)`
- Flexbox layout — `display: flex`, `justify-content`, `align-items`
- Responsive design — `@media (max-width: 640px)`
- Hover and focus states — `:hover`, `:focus`
- CSS transitions for smooth animations
- `position: sticky` for the navbar

---

## Step 3 — SQLite database
**Date:** April 2025
**Branch:** feature/initial-setup

### What we are building
Replacing the in-memory Go slice with a real SQLite database so that
resources survive server restarts.

### Why SQLite and not PostgreSQL or MySQL
- SQLite requires zero setup — it is just a file on disk
- No separate database server process needed
- Fast enough for hundreds of concurrent students on a LAN
- The entire database is one file — easy to back up, move, or copy
- Can upgrade to PostgreSQL later if the app grows beyond one campus

### SQL commands we will use
- `CREATE TABLE` — create the resources table once at startup
- `INSERT INTO` — save a new resource when someone uploads
- `SELECT` — read resources when someone browses or searches