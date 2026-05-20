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

#Month 3 — IN PROGRESS
Branch: feature/month3

##Step 15 — User authentication
Date: May 2026
Branch: feature/month3
###       What we built

New auth.go file — all authentication code separated from main.go
User registration — students create an account with username, email, password
User login — students sign in with username and password
User logout — session cleared, cookie deleted
Session management — Go remembers who is logged in between requests
Navbar updates dynamically — shows username and logout when logged in,
shows login and register buttons when logged out
newPageData() helper function — ensures every page always has the
correct session state passed to templates automatically
Welcome message shown after successful registration

###New files created

auth.go — contains all auth functions, handlers and session store
templates/register.html — registration form page
templates/login.html — login form page

###New packages added

golang.org/x/crypto/bcrypt — for hashing passwords safely
github.com/gorilla/sessions — for managing login sessions
github.com/gorilla/securecookie — pulled in automatically by gorilla/sessions

##Key decisions
###Why a separate auth.go file?
Keeping all authentication code in one file makes it easier to find,
read, and modify. main.go was getting long — splitting by concern
keeps both files focused and readable.
###Why bcrypt and not MD5 or SHA256?
MD5 and SHA256 are fast — which is bad for passwords. An attacker can
try billions of guesses per second. bcrypt is deliberately slow
(cost factor 12) — it takes ~250ms per hash making brute force
attacks practically impossible. It also automatically handles salting
which prevents rainbow table attacks.
###Why show "Invalid username or password" for both wrong username and wrong password?
If we said "username not found" vs "wrong password" separately, an
attacker could use the login form to discover which usernames exist
in the system. The generic message gives away no information.

###Why newPageData() helper function?
Every template needs LoggedIn and CurrentUser to render the navbar
correctly. Without the helper, every single handler had to manually
call getSessionUser(r) and add those fields — easy to forget and
caused the bug where the upload page showed the wrong navbar state.
The helper makes it impossible to forget.

##Challenges encountered and how they were solved
###Challenge 1 — `go run main.go` stopped working
After creating auth.go, running go run main.go gave errors:
undefined: User, undefined: createUsersTable etc.
######Root cause: `go run main.go `only compiles one file. When the project
has multiple Go files, all of them must be included.
######Solution: Use `go run .` instead — this compiles ALL Go files in the
current directory automatically. From this point forward, always use
`go run . `not `go run main.go.`

##Challenge 2 — Login worked but navbar still showed "Log in / Register"
After logging in successfully, the homepage showed the correct logged-in
navbar, but the upload page still showed the logged-out state.
######Root cause: uploadHandler was building PageData without including
LoggedIn and CurrentUser fields. Only homeHandler had those fields.
######Solution: Created the newPageData(r, title) helper function that
automatically reads the session and fills in LoggedIn and CurrentUser
for every handler. Updated all handlers to use it.

##Challenge 3 — Session not being read correctly after login
Login completed without error and users were saved to the database, but
after redirect the server could not read the session — treating the user
as logged out on every page.
######Root cause: Two issues combined:

gorilla/sessions uses encoding/gob to encode the session cookie,
but the int type was not registered with gob — causing silent
decode failure on every subsequent request
Session cookie options (Path, MaxAge, HttpOnly) were not set
explicitly, causing the cookie to expire immediately or be rejected

######Solution:

Added gob.Register(int(0)) in the init() function to register
the int type with gob encoder
Set explicit cookie options in init():
Path: "/", MaxAge: 86400 * 365, HttpOnly: true
Updated getSessionUser() to handle multiple possible decoded types
using a type switch — handles int, int64, and float64
Cleared browser cookies to remove stale broken session cookies

 # Step 15 — Go Concepts and SQL Reference

This document covers every new Go concept and SQL statement introduced
during Step 15 of the ElimuLocal project — user authentication.

---

## Go Concepts Introduced

### `encoding/gob`

Go's built-in binary encoder. `gorilla/sessions` uses it internally to
encode and decode the session cookie. Without registering your types with
`gob`, it silently fails to decode values — which was the root cause of
the session bug in Challenge 3.

### `gob.Register(int(0))`

Tells the `gob` encoder that the `int` type exists and should be handled
correctly. You call this in `init()` before the server starts so `gob` is
ready before any session is ever created.

### `init()`

A special Go function that runs automatically before `main()`. You cannot
call it manually — Go calls it for you. We use it to configure the session
store options so they are set correctly from the very first request.

```go
func init() {
    gob.Register(int(0))

    store.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   86400 * 365,
        HttpOnly: true,
    }
}
```

### Type Switch — `switch v := value.(type)`

Checks the actual runtime type of an interface value and lets you handle
each type differently. We use this in `getSessionUser()` because
`gorilla/sessions` can return the `userID` as `int`, `int64`, or `float64`
depending on how it was encoded. The type switch handles all three cases
safely.

```go
switch v := session.Values["userID"].(type) {
case int:
    userID = v
case int64:
    userID = int(v)
case float64:
    userID = int(v)
default:
    return User{}, false
}
```

### `bcrypt.GenerateFromPassword()`

Takes a plain text password and returns a hash — a scrambled string that
cannot be reversed. The cost factor of `12` makes it deliberately slow
(roughly 250ms per hash) so brute force attacks are impractical. This hash
is what gets stored in the database, never the real password.

```go
hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
```

### `bcrypt.CompareHashAndPassword()`

Checks whether a plain password matches a stored hash. It returns `nil` if
they match, and an error if they do not. Used every time a student tries to
log in.

```go
err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
if err == nil {
    // password is correct
}
```

### `sessions.NewCookieStore()`

Creates a session store backed by signed browser cookies. The secret key is
used to sign every cookie — if anyone tampers with the cookie value, the
signature check fails and the session is rejected entirely.

```go
var store = sessions.NewCookieStore([]byte("your-secret-key"))
```

### `session.Values["key"]`

Reads or writes a value inside the session. It works like a map — you can
store any value under any string key. We store the logged-in user's ID
under the key `"userID"`.

```go
// writing
session.Values["userID"] = user.ID

// reading
userID := session.Values["userID"]
```

### `session.Save(r, w)`

Must be called after every change to a session. Without it, changes are
lost — the updated cookie is never sent to the browser. Forgetting this
call was part of the root cause of Challenge 3.

```go
err = session.Save(r, w)
if err != nil {
    // handle error
}
```

### `result.LastInsertId()`

Returns the auto-generated ID of the row that was just inserted into the
database. We use it immediately after registration to log the new user in
without making them go through the login page separately.

```go
result, err := db.Exec("INSERT INTO users ...")
userID, _ := result.LastInsertId()
setSessionUser(w, r, int(userID))
```

### `go run .`

Compiles and runs all Go files in the current directory. Once the project
has more than one `.go` file — like `main.go` and `auth.go` — you must use
this instead of `go run main.go`, which only compiles a single named file.

```bash
# wrong — only compiles main.go
go run main.go

# correct — compiles all .go files in the folder
go run .
```

---

## New SQL Introduced

### Creating the Users Table

```sql
CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TEXT NOT NULL
);
```

`UNIQUE` on `username` and `email` means the database itself rejects
duplicate registrations — we do not have to rely only on our Go checks.
`NOT NULL` means these fields can never be left empty. `AUTOINCREMENT`
means SQLite assigns the next available ID number automatically so we never
manage IDs manually.

### Checking If a Username or Email Is Already Taken

```sql
SELECT COUNT(*) FROM users WHERE username = ?
SELECT COUNT(*) FROM users WHERE email = ?
```

Returns `0` if the value does not exist yet, `1` if it does. We check this
before trying to insert a new user so we can show a friendly error message
to the student rather than letting the database throw a raw `UNIQUE`
constraint error.

### Inserting a New User

```sql
INSERT INTO users (username, email, password_hash, created_at)
VALUES (?, ?, ?, ?)
```

Notice there is no `password` column — only `password_hash`. The real
password never touches the database. `bcrypt` converts it to a hash in Go
before this query ever runs.

### Looking Up a User at Login

```sql
SELECT id, username, email, password_hash, created_at
FROM users WHERE username = ?
```

Returns the full user record including the stored hash. Go then uses
`bcrypt.CompareHashAndPassword()` to check whether the submitted password
matches that hash — the comparison always happens in Go, never in SQL.

Step 18 — Video Upload and Playback
Date: May 2026
Branch: feature/month3
What we built

Upload size limit increased from 10MB to 500MB to accommodate
lecture videos
MP4, MKV and WebM added to the list of accepted file types alongside
PDF
Video Lecture added as a category option in the upload form, edit
form, and search filter dropdown
Server timeouts increased — ReadTimeout and WriteTimeout both set
to 5 minutes so large uploads are not cut off mid-transfer
File input in upload.html updated to show accepted formats and the
500MB size limit as a hint to the user
HTML5 video player renders automatically in preview.html when the
file extension is .mp4, .mkv, or .webm
map[string]bool used for file type validation — cleaner and more
maintainable than a long boolean expression

Key decisions
Why 500MB and not unlimited?
An unlimited upload size means a single student could accidentally or
maliciously upload a file that fills the entire server disk. 500MB is
large enough for a full lecture recording at reasonable quality while
still protecting against runaway disk usage.
Why add Video Lecture as a separate category?
A student filtering for Notes or Past Papers should not see video files
in their results — the format is completely different. A dedicated category
keeps filtering accurate. Students looking for a quick text summary do not
want a 400MB video appearing in their results.
Why increase server timeouts?
A 400MB upload on a slow campus WiFi connection can take several minutes.
Go's default timeouts would cut the connection before the upload finishes,
leaving the student with a confusing error after waiting several minutes.
Go concepts introduced

map[string]bool as a set for allowed file types — cleaner than a long
boolean expression, adding a new type means adding one line to the map
500 << 20 — bitwise left shift, Go idiom for expressing file sizes in
bytes, 1 << 20 equals exactly one megabyte so 500 << 20 is exactly
500 megabytes
&http.Server{} — explicit server struct allowing custom ReadTimeout,
WriteTimeout and IdleTimeout values essential for large file handling


Landing Page
Date: May 2026
Branch: feature/month3
What we built

Complete standalone landing page at / — the new first impression of
ElimuLocal for students, administrators and investors
GSAP-powered loading animation — ElimuLocal splits open with a
university photo expanding between the two halves of the name
Dark hero section with animated headline, navigation links and two
call-to-action buttons — Browse Resources and Register Free
Mission statement section explaining the three problems ElimuLocal
solves in large editorial typography
Three problem and solution cards — expensive textbooks, notes trapped
in WhatsApp, unreliable internet — each paired with the ElimuLocal fix
How it works section — three numbered steps
About section with author name, role at Zone01 Kisumu, and a vision
quote
Green call-to-action section at the bottom
Dark footer with GitHub link
Old / browse route moved to /browse — the landing page is now the
entry point for new visitors
Logged-in users see Browse Resources in the hero nav instead of
Get Started Free — the page adapts to session state
renderLanding() helper renders the landing template independently
from base.html since the landing page has its own full HTML
structure, fonts and scripts

Key decisions
Why a standalone template and not base.html?
The landing page has a completely different visual style — dark hero,
GSAP animations, its own nav structure. Forcing it into base.html
would require overriding almost everything base.html provides. A
standalone template is cleaner and easier to maintain independently.
Why move / to /browse?
A returning student who bookmarks the browse page should land directly
on resources without clicking through the landing page every time. The
landing page serves new visitors. /browse serves regular users.
Why show different nav buttons based on login state?
A logged-in student seeing Register Free is confusing — they already have
an account. Showing Browse Resources instead is contextually correct and
routes them directly to what they came to do.
Why GSAP for the animation?
CSS animations alone cannot achieve the letter-splitting reveal with a
photo expanding between the halves — that requires JavaScript to
coordinate multiple elements simultaneously with precise timing. GSAP
is the industry standard for this type of animation and loads from a
CDN with no installation needed.
