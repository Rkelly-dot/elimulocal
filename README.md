# 📚 ElimuLocal

> A learning resource hub for Kenyan university students — built with Go.

ElimuLocal lets university students share and download study materials (notes,
past papers, textbooks) over a campus WiFi network. No internet required.
No cloud. Just a laptop and a local network.

---

## 🗺️ What this project is

A beginner Go web application that solves a real problem: Kenyan university
students struggle to find and share quality study materials. ElimuLocal gives
them a simple platform to upload and browse resources organized by university,
course, and category — all running locally on a campus network.

---

## 👥 Who does what — roles and workflow

ElimuLocal has two types of users. Here is exactly what each person does.

---

### 🖥️ The administrator (the person running the server)

This is the person whose laptop or machine runs the ElimuLocal server.
On a campus this would typically be a tech-savvy student, a lab attendant,
or a teacher comfortable with a terminal. Only one person needs to do this.

**One-time setup (do this once):**
1. Install Go on your machine
2. Download the ElimuLocal project folder
3. Open a terminal and navigate into the project folder
4. Run `go run main.go` to start the server
5. Find your local IP address using `ip addr show` on Linux
6. Share the link `http://YOUR-IP:8080` with everyone on campus

**Every time you want ElimuLocal to be available:**
1. Connect your machine to the campus WiFi
2. Open a terminal and run `go run main.go`
3. Keep the terminal open — closing it shuts the server down for everyone

**To shut it down:**
- Press `Ctrl + C` in the terminal, or simply close the terminal window

> ⚠️ Important: The server only works while your machine is on, the terminal
> is open, and your machine is connected to campus WiFi. If you close your
> laptop or disconnect from WiFi, the platform goes offline for everyone.
> In Month 2 we move the server to a dedicated machine so it stays up permanently.

---

### 👩‍🏫 Teachers — sharing resources

Teachers do not need to install anything. Just a browser and campus WiFi.

**How a teacher shares a resource:**
1. Connect to campus WiFi
2. Open any browser (Chrome, Firefox, Edge — anything works)
3. Go to the ElimuLocal link: `http://192.168.x.x:8080`
4. Click **"+ Share Resource"** in the top right corner
5. Fill in the form:
   - **Title** — a clear, descriptive name for the resource
   - **Course** — e.g. "Computer Science", "Economics", "Civil Engineering"
   - **Category** — Notes, Past Paper, Textbook, Summary, or Other
   - **University** — type your university name (it auto-suggests to others next time)
   - **Description** — a short summary of what is in the resource
   - **Your name** — optional, but students appreciate knowing who shared it
6. Click **"Share Resource"** — it immediately appears for all students

**What teachers cannot do yet (coming in Month 2):**
- Upload actual PDF files (currently only resource information is saved)
- Edit or delete a resource after submitting
- See download statistics for their materials

---

### 👨‍🎓 Students — finding and downloading resources

Students also need nothing installed. Just a browser and campus WiFi.

**How a student finds study materials:**
1. Connect to campus WiFi
2. Open any browser on your phone, tablet, or laptop
3. Go to the ElimuLocal link shared by the administrator: `http://192.168.x.x:8080`
4. The homepage shows all available resources
5. Use the **search bar** to find by keyword — e.g. "organic chemistry" or "past paper"
6. Use the **university filter** to narrow results to your institution
7. Browse resource cards — each shows the title, description, university,
   course, uploader name, upload date, and download count
8. Click **"Download"** to save the file to your device

**How a student contributes back:**
Any student can upload their own notes and past papers using the same
**"+ Share Resource"** button. The more students contribute, the richer
the platform becomes for everyone.

**What to do if you cannot connect:**
- Make sure you are on the same WiFi as the server — not mobile data
- Ask the administrator to confirm the server is running
- Double-check the IP address — one wrong digit and it will not load
- If nothing works, the server machine may be off or disconnected from WiFi

---

### 🔄 The full workflow at a glance

```
Administrator                    Teachers & Students
─────────────                    ───────────────────

1. Connect to campus WiFi        1. Connect to campus WiFi
2. Run: go run main.go           2. Open any browser
3. Share the IP address link ──► 3. Go to http://192.168.x.x:8080
4. Keep terminal open            4. Browse, search, upload resources
                                 5. Download what they need
```

---

## 🚀 How to run it

### Requirements
- Go 1.21 or newer
- A Linux, Mac, or Windows machine
- (Optional) A campus WiFi network to share with other students

### Run locally (just you)
```bash
# 1. Clone or download the project
cd elimulocal

# 2. Start the server
go run main.go

# 3. Open your browser
# Go to: http://localhost:8080
```

### Share on campus WiFi (all students on same network)
```bash
# 1. Find your laptop's local IP address
ip addr show        # Linux
ifconfig            # Mac
ipconfig            # Windows

# Look for something like: 192.168.x.x

# 2. Start the server
go run main.go

# 3. Tell students to open their browser and go to:
# http://192.168.x.x:8080
# (replace x.x with your actual IP)
```

---

## 📁 Project structure

```
elimulocal/
├── main.go                  ← All Go server code lives here
├── go.mod                   ← Go module file (like package.json in Node)
├── go.sum                   ← Auto-generated dependency checksums
├── README.md                ← You are here
│
├── templates/               ← HTML page templates
│   ├── base.html            ← Shared layout (navbar, footer)
│   ├── home.html            ← Browse/search page
│   └── upload.html          ← Upload a resource page
│
├── static/                  ← Files served directly to the browser
│   ├── css/
│   │   └── style.css        ← All styling
│   └── fonts/               ← Local fonts (no internet needed)
│
└── uploads/                 ← Uploaded PDF files are saved here
    └── .gitkeep
```

---

## 🧠 Go concepts used in this project

This project is intentionally beginner-friendly. Every Go concept is
explained with comments in the code. Here is a quick map:

| Concept | Where to find it |
|---|---|
| Structs | `type Resource struct` in main.go |
| Slices | `var resources = []Resource{...}` |
| Maps | `getUniversities()` function |
| Functions | Every `func` in main.go |
| Error handling | Every `if err != nil` block |
| Range loops | `for _, r := range resources` |
| HTTP handlers | `func homeHandler(...)` |
| HTML templates | `template.ParseFiles(...)` |
| File handling | `uploadHandler` POST section |
| SQLite database | `initDB()` and query functions |

---

## ✨ Features

### Working now
- Browse all uploaded study resources
- Search by keyword (title, course, description)
- Filter by university (auto-populated from submissions)
- Upload new resources with title, course, category, description
- University list grows automatically as students submit

### Coming soon (see roadmap below)
- Real PDF file uploads and downloads
- SQLite database (resources survive server restarts)
- Download counter
- Ratings / upvotes
- Offline-ready fonts (no internet dependency)
- PDF preview in browser

---

## 🗓️ 3-month build roadmap

### Month 1 — Core product
- [x] Project setup and folder structure
- [x] HTML templates and CSS styling
- [x] Dynamic university list (no hardcoding)
- [ ] SQLite database — resources persist after restart
- [ ] Real file uploads — accept and save PDFs
- [ ] Real file downloads — serve saved PDFs
- [ ] Embed fonts locally — no internet needed

### Month 2 — Real users
- [ ] Deploy on campus LAN with real students
- [ ] Improved search (filter by year, category)
- [ ] Ratings / upvote system
- [ ] Fix bugs from real student feedback
- [ ] Download counter per resource

### Month 3 — Polish and scale
- [ ] PDF preview in browser (no download needed)
- [ ] Mobile-friendly improvements
- [ ] Support multiple campuses
- [ ] Write setup guide so others can run their own instance
- [ ] Clean up code and add more comments

---

## 🌍 Design decisions

**Why Go?**
Go is simple, fast, and produces a single binary with no dependencies.
A student can download one file and run a server. Perfect for low-resource
environments like a Raspberry Pi or old laptop on a campus network.

**Why LAN instead of the internet?**
Many Kenyan university students have limited mobile data. A campus LAN means
zero data cost for every student using the app. It also works during internet
outages, which are common.

**Why no login/accounts yet?**
Keeping it simple for Month 1. The priority is getting real materials shared
between real students. Authentication adds complexity and friction — we will
add it once the core sharing flow is solid.

**Why SQLite?**
No setup, no separate server process, one file on disk. For a campus LAN
serving hundreds of students, SQLite is more than fast enough. PostgreSQL
can come later if needed.

---

## 🤝 Contributing

This is a learning project built step by step. If you are a student at a
Kenyan university and want to help test it or add your university's materials,
reach out.

---

## 📖 Learning resources

If you are learning Go alongside this project, these are recommended:

- [A Tour of Go](https://go.dev/tour/) — official interactive tutorial
- [Go by Example](https://gobyexample.com/) — practical code snippets
- [Let's Go (book)](https://lets-go.alexedwards.net/) — web apps in Go

---

*Built for Kenyan university students. Share knowledge, uplift each other. 🇰🇪*
