package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Resource struct {
	ID          int
	Title       string
	Course      string
	University  string
	Category    string
	Description string
	UploadedBy  string
	UploadedAt  string
	FileName    string
	Downloads   int
	Upvotes    int
}

type PageData struct {
	Title        string
	Resources    []Resource
	Universities []string
	Message      string
	Search       string
	University   string
	Category     string
	Sort         string
}

var db *sql.DB

func initDB() {
	var err error

	db, err = sql.Open("sqlite", "elimulocal.db")
	if err != nil {
		log.Fatal("Could not open database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS resources (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		course      TEXT NOT NULL,
		university  TEXT NOT NULL,
		category    TEXT NOT NULL,
		description TEXT,
		uploaded_by TEXT,
		uploaded_at TEXT,
		file_name   TEXT,
		downloads   INTEGER DEFAULT 0
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Could not create table:", err)
	}

	seedDB()

	fmt.Println("Database ready — elimulocal.db")
}

func seedDB() {
	_, _ = db.Exec("ALTER TABLE resources ADD COLUMN upvotes INTEGER DEFAULT 0")
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM resources").Scan(&count)
	if err != nil {
		log.Fatal("Could not count resources:", err)
	}

	if count > 0 {
		return
	}

	seeds := []Resource{
		{
			Title:       "Introduction to Data Structures - Full Notes",
			Course:      "Computer Science",
			University:  "University of Nairobi",
			Category:    "Notes",
			Description: "Covers arrays, linked lists, stacks, queues and trees.",
			UploadedBy:  "Alice M.",
			UploadedAt:  "2025-03-01",
			FileName:    "",
		},
		{
			Title:       "Principles of Economics - Chapter Summaries",
			Course:      "Economics",
			University:  "Strathmore University",
			Category:    "Notes",
			Description: "Clear summaries of micro and macroeconomics fundamentals.",
			UploadedBy:  "David O.",
			UploadedAt:  "2025-03-10",
			FileName:    "",
		},
	}

	for _, s := range seeds {
		saveResource(s)
	}

	fmt.Println("Seed data added to database.")
}

func saveResource(r Resource) error {
	query := `
	INSERT INTO resources (title, course, university, category, description, uploaded_by, uploaded_at, file_name, downloads)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0)`

	_, err := db.Exec(query,
		r.Title,
		r.Course,
		r.University,
		r.Category,
		r.Description,
		r.UploadedBy,
		r.UploadedAt,
		r.FileName,
	)
	return err
}

func getResources(search, university, category, sort string) ([]Resource, error) {
	query := "SELECT id, title, course, university, category, description, uploaded_by, uploaded_at, file_name, downloads, upvotes FROM resources"

	var args []interface{}
	var conditions []string

	if search != "" {
		conditions = append(conditions, "(title LIKE ? OR course LIKE ? OR description LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	if university != "" {
		conditions = append(conditions, "university = ?")
		args = append(args, university)
	}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	switch sort {
	case "popular":
		query += " ORDER BY downloads DESC"
	case "upvotes":
	query += " ORDER BY upvotes DESC"
	case "oldest":
		query += " ORDER BY id ASC"
	default:
		query += " ORDER BY id DESC"
	}
	
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []Resource
	for rows.Next() {
		var r Resource
		err := rows.Scan(
			&r.ID,
			&r.Title,
			&r.Course,
			&r.University,
			&r.Category,
			&r.Description,
			&r.UploadedBy,
			&r.UploadedAt,
			&r.FileName,
			&r.Downloads,
			&r.Upvotes,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, r)
	}

	return resources, nil
}

func getUniversities() []string {
	rows, err := db.Query("SELECT DISTINCT university FROM resources ORDER BY university")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var unis []string
	for rows.Next() {
		var u string
		rows.Scan(&u)
		unis = append(unis, u)
	}
	return unis
}

func incrementDownloads(id int) {
	db.Exec("UPDATE resources SET downloads = downloads + 1 WHERE id = ?", id)
}

func incrementUpvotes(id int) error {
	_, err := db.Exec("UPDATE resources SET upvotes = upvotes + 1 WHERE id = ?", id)
	return err
}

func upvoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/upvote/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}
	var id int
	fmt.Sscan(idStr, &id)
	if id == 0 {
		http.NotFound(w, r)
		return
	}

	err := incrementUpvotes(id)
	if err != nil {
		http.Error(w, "Could not upvote resource", http.StatusInternalServerError)
		return
	}

	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	search := r.URL.Query().Get("search")
	university := r.URL.Query().Get("university")
	category := r.URL.Query().Get("category")
	sort := r.URL.Query().Get("sort")

	resources, err := getResources(search, university, category, sort)
	if err != nil {
		http.Error(w, "Could not load resources: "+err.Error(), http.StatusInternalServerError)
		return
	}
	message := ""
	if r.URL.Query().Get("success") == "1" {
    message = "✅ Resource uploaded successfully! Students can now find and download it."
	}

	data := PageData{
		Title:        "ElimuLocal - Browse Study Materials",
		Resources:    resources,
		Universities: getUniversities(),
		Search:       search,
		University:   university,
		Category:     category,
		Sort:         sort,
		Message:      message,
	}

	renderTemplate(w, "home.html", data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := PageData{
			Title:        "Share a Resource - ElimuLocal",
			Universities: getUniversities(),
		}
		renderTemplate(w, "upload.html", data)
		return
	}

	if r.Method == "POST" {
		title := strings.TrimSpace(r.FormValue("title"))
		course := strings.TrimSpace(r.FormValue("course"))
		university := strings.TrimSpace(r.FormValue("university"))
		category := r.FormValue("category")
		description := strings.TrimSpace(r.FormValue("description"))
		uploader := strings.TrimSpace(r.FormValue("uploader"))

		if uploader == "" {
			uploader = "Anonymous"
		}

		if title == "" || course == "" || university == "" {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Please fill in all required fields.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "File too large. Maximum size is 10MB.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Please select a PDF file to upload.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".pdf" {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Only PDF files are allowed.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
		filePath := filepath.Join("uploads", fileName)

		dst, err := os.Create(filePath)
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Could not save file. Please try again.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Could not save file. Please try again.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		newResource := Resource{
			Title:       title,
			Course:      course,
			University:  university,
			Category:    category,
			Description: description,
			UploadedBy:  uploader,
			UploadedAt:  time.Now().Format("2006-01-02"),
			FileName:    fileName,
		}

		err = saveResource(newResource)
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Could not save resource. Please try again.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		http.Redirect(w, r, "/?success=1", http.StatusSeeOther)
		return
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/download/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	var fileName string
	var title string
	err := db.QueryRow(
		"SELECT file_name, title FROM resources WHERE id = ?", idStr,
	).Scan(&fileName, &title)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fileName == "" {
		http.Error(w, "No file available for this resource yet.", http.StatusNotFound)
		return
	}

	filePath := filepath.Join("uploads", fileName)

	var id int
	fmt.Sscan(idStr, &id)
	incrementDownloads(id)

	w.Header().Set("Content-Disposition", "attachment; filename="+title+".pdf")
	w.Header().Set("Content-Type", "application/pdf")

	http.ServeFile(w, r, filePath)
}

func renderTemplate(w http.ResponseWriter, page string, data PageData) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+page)
	if err != nil {
		http.Error(w, "Could not load page: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Could not render page: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	initDB()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/upvote/", upvoteHandler)

	fmt.Println("ElimuLocal is running!")
	fmt.Println("Open your browser and go to: http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop the server.")

	log.Fatal(http.ListenAndServe(":8080", nil))
}