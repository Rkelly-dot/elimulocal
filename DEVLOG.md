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

Step 7 — Improved search and filtering
Date: April 2025
Branch: feature/month2-improvements
What we built

Added category filter dropdown to the search form
(All Categories / Notes / Past Paper / Textbook / Summary / Other)
Added sort dropdown to the search form
(Newest first / Most downloaded / Most helpful / Oldest first)
Dropdowns remember their selected value after form submission
Added "✕ Clear filters" link that appears when any filter is active
Added count badge showing how many results were found
Updated getResources() to accept category and sort parameters
Updated PageData struct with Search, University, Category, Sort fields

Key decisions

Used {{if eq .Category "Notes"}}selected{{end}} pattern to keep
dropdowns showing the active filter after search — without this the
dropdown resets every time which is very annoying UX
Used a switch statement for sort order — cleaner than if/else chain
Used LIKE ? with %search% pattern for partial text matching in SQL

Go concepts introduced

switch statement — cleaner alternative to multiple if/else blocks
Building SQL queries dynamically — appending WHERE clauses conditionally
[]interface{} — slice that can hold values of any type, used for
SQL query arguments when we do not know how many there will be

SQL concepts introduced

LIKE '%keyword%' — partial text matching (% is a wildcard)
ORDER BY column ASC/DESC — sorting results
ORDER BY downloads DESC — sort by most downloaded first


Step 8 — Ratings and upvotes
Date: April 2025
Branch: feature/month2-improvements
What we built

Added upvotes column to the resources table using ALTER TABLE
Added 👍 Helpful (count) button to every resource card
Clicking the button increments the upvote count in the database
Added "Most helpful" sort option to show most upvoted resources first
Added incrementUpvotes() function in main.go
Added upvoteHandler to handle POST /upvote/{id} requests

How the upvote flow works
Student clicks 👍 Helpful
    → browser sends POST /upvote/3
    → upvoteHandler extracts ID from URL
    → incrementUpvotes(3) runs UPDATE SQL
    → redirects back to the page the student was on
    → student sees updated count
Key decisions

Used ALTER TABLE resources ADD COLUMN upvotes INTEGER DEFAULT 0
with _, _ = db.Exec(...) to safely add the column to an existing
database — ignoring the error if the column already exists
Used r.Header.Get("Referer") for redirect — sends student back to
whatever page they were on, preserving their search filters
Used POST not GET for upvoting — upvoting changes data so it must
be a POST request, not a link

Go concepts introduced

ALTER TABLE for modifying existing database tables safely
r.Header.Get("Referer") — reading the previous page URL from
the request headers for smart redirects


Step 9 — Download counter on cards
Date: April 2025
Branch: feature/month2-improvements
What we built

Display the download count on every resource card
Counter was already being tracked in the database since Month 1 —
this step simply surfaced the number in the UI

What was learned
This step was intentionally simple — it demonstrated that good data
modelling pays off. Because we stored downloads in the database from
day one, showing the count was a single line of HTML. No backend
changes needed at all.

Step 10 — Success and error flash messages
Date: April 2025
Branch: feature/month2-improvements
What we built

Green success message appears at top of homepage after upload:
"✅ Resource uploaded successfully!"
Message fades out automatically after 4 seconds using JavaScript
Added deleted confirmation message for future delete feature
Used URL query parameter ?success=1 to pass the message signal
from the upload redirect to the home page

How flash messages work
1. Upload succeeds
2. Go redirects to /?success=1
3. homeHandler reads r.URL.Query().Get("success")
4. Sets message in PageData
5. home.html shows the message div
6. JavaScript setTimeout fades it out after 4 seconds
Key decisions

Used URL query parameter instead of cookies or sessions — simpler
and works without any session management infrastructure
Used JavaScript for fade-out — CSS alone cannot remove an element
after a time delay
setTimeout with two nested calls — first fades opacity over 1
second, then sets display:none after the fade completes

JavaScript concepts introduced

setTimeout(function, milliseconds) — run code after a delay
element.style.transition — smooth CSS property changes via JS
element.style.opacity — change transparency
element.style.display = 'none' — hide element completely

Step 11 — Campus LAN deployment
Date: April 2025
Branch: feature/month2-improvements
What we did

Found the laptop's local IP address using ip addr show
Confirmed firewall was inactive (ufw status showed inactive)
Started the server and shared the link with a real device
Tested ElimuLocal on a phone over campus WiFi

Result
ElimuLocal loaded successfully on a phone at http://192.168.89.223:8080.
All features worked on mobile — search, filter, browse, upload form.
Minor navbar layout issue on small screens noted for Step 12.
How campus LAN works
Your laptop runs go run main.go
    → Go listens on :8080 (all network interfaces)
    → Campus router assigns your laptop IP 192.168.89.223
    → Student phone connects to same WiFi
    → Student opens http://192.168.89.223:8080
    → Phone sends request across WiFi to your laptop
    → Go responds with the ElimuLocal page
    → No internet involved at any point
Commands learned

ip addr show | grep "inet " — find your local IP address
sudo ufw status — check if the firewall is blocking connections

Step 12-Fix mobile navbar layout
Date: April 2025
Branch: feature/month2-improvements
what we did

I patched the stylesheet to stop the brand from breaking and prevent the button overlap on screens
Applied a small CSS patch by allowing wrapping and making the share button stack on very small screen.This helped keep the brand on one line and make the header responsive

Result
ElimuLocal ran successfully on a mobile phone and the navbar was not cramped

Month 2
Date: April 2025
Month 2 checklist

 Step 7  — Improved search with category filter and sort options
 Step 8  — Upvote / helpful rating system
 Step 9  — Download counter displayed on resource cards
 Step 10 — Success flash messages after upload
 Step 11 — Deployed on campus LAN, tested on real phone
 Step 12 — Fix mobile navbar layout
 Step 13 — Merge Month 2 into main branch

What ElimuLocal can do at end of Month 2

Filter resources by category (Notes, Past Paper, Textbook etc.)
Sort resources by newest, most downloaded, most helpful, oldest
Upvote helpful resources — community quality signal
See download counts on every resource card
Get clear feedback after uploading a resource
Run live on a campus WiFi network accessible from any device
Tested and working on a real Android phone

