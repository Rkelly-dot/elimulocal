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
## Step 3 — SQLite database and file uploads
**Date:** April 2025
**Branch:** feature/initial-setup

### What we built
- Added SQLite database using `modernc.org/sqlite` package
- Created `resources` table with all fields
- Replaced in-memory slice with real database reads and writes
- Real PDF file uploads — files saved to uploads/ folder on disk
- Real file downloads — students can download uploaded PDFs
- Download counter increments every time a file is downloaded
- Seed data added automatically on first run
- Data survives server restarts permanently

### New functions added to main.go
- `initDB()` — opens database connection and creates table
- `seedDB()` — adds starter data if database is empty
- `saveResource()` — inserts a new resource row into the database
- `getResources()` — reads and filters resources from the database
- `getUniversities()` — reads unique university names from the database
- `incrementDownloads()` — updates download count when file is downloaded
- `downloadHandler()` — serves PDF files to the browser

### SQL commands learned
- `CREATE TABLE IF NOT EXISTS` — creates table safely on every startup
- `INSERT INTO ... VALUES (?, ?, ...)` — adds a new row
- `SELECT ... FROM ... WHERE ... LIKE ?` — reads and filters rows
- `SELECT DISTINCT` — returns unique values only
- `UPDATE ... SET ... WHERE` — modifies an existing row
- `SELECT COUNT(*)` — counts rows in a table
- `ORDER BY id DESC` — sorts newest first

### Key decisions
- Used `modernc.org/sqlite` — pure Go driver, no C compiler needed
- Used `?` placeholders in all SQL queries — prevents SQL injection
- Used `time.Now().UnixNano()` for unique filenames — prevents overwrites
- Used `io.Copy()` for file saving — handles large files without loading
  them entirely into memory
- Used `defer rows.Close()` and `defer file.Close()` — ensures cleanup
  always happens even if an error occurs

### Go concepts introduced
- Blank import `_ "modernc.org/sqlite"` — runs init without direct use
- `database/sql` package — Go's standard database interface
- `db.Exec()` — runs SQL that does not return rows
- `db.Query()` — runs SQL that returns multiple rows
- `db.QueryRow()` — runs SQL that returns exactly one row
- `rows.Scan()` — reads a database row into Go variables
- `&variable` — passing a pointer so Scan can write into the variable
- `defer` — guarantees cleanup code runs when a function ends
- `r.ParseMultipartForm()` — prepares Go to handle file uploads
- `r.FormFile()` — retrieves an uploaded file from the request
- `os.Create()` — creates a new file on disk
- `io.Copy()` — copies data from one stream to another
- `filepath.Join()` — builds file paths correctly for the OS
- `filepath.Ext()` — gets the file extension
- `10 << 20` — bitwise shift used to express 10MB in bytes
## Step 4 — Local fonts
**Date:** April 2025
**Branch:** feature/initial-setup

### What we built
- Downloaded Sora (Regular and SemiBold) and Space Mono font files
- Added @font-face rules to style.css to load fonts from local disk
- Removed Google Fonts dependency completely
- App now works with zero internet connection on campus LAN

### Why this matters
- Students on campus WiFi with no internet can still use the full app
- No external requests means faster load times on slow networks
- No dependency on Google's servers being available

---

## Month 1 — COMPLETE
**Date:** April 2025

### Month 1 checklist
- [x] Project setup and folder structure
- [x] HTML templates and CSS styling
- [x] Dynamic university list
- [x] SQLite database
- [x] Real file uploads
- [x] Real file downloads
- [x] Local fonts — no internet dependency

### What ElimuLocal can do at end of Month 1
- Run on a campus LAN with no internet
- Accept PDF uploads from teachers and students
- Store resources permanently in a database
- Serve files for download
- Search and filter resources
- Auto-populate university suggestions
- Look like a real product with custom fonts and styling