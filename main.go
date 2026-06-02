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

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
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
	Upvotes     int
	UserID      int
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
	Resource     Resource
	CurrentUser  User
	LoggedIn     bool
	IsVideo      bool
	MimeType     string
}

var db *sql.DB

func initDB() {
	var err error

	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_TOKEN")

	if tursoURL != "" && tursoToken != "" {
		connStr := tursoURL + "?authToken=" + tursoToken
		db, err = sql.Open("libsql", connStr)
		if err != nil {
			log.Fatal("Could not open Turso database:", err)
		}
		fmt.Println("Connected to Turso cloud database")
		runMigrations(db)
	} else {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "elimulocal.db"
		}
		db, err = sql.Open("sqlite", dbPath)
		if err != nil {
			log.Fatal("Could not open local database:", err)
		}
		fmt.Println("Connected to local SQLite database")
		runMigrations(db)
	}
}
func seedDB() {
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
	INSERT INTO resources (title, course, university, category, description, uploaded_by, uploaded_at, file_name, downloads, upvotes, user_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, 0, ?)`

	_, err := db.Exec(query,
		r.Title,
		r.Course,
		r.University,
		r.Category,
		r.Description,
		r.UploadedBy,
		r.UploadedAt,
		r.FileName,
		r.UserID,
	)
	return err
}

func getResources(search, university, category, sort string) ([]Resource, error) {
	query := "SELECT id, title, course, university, category, description, uploaded_by, uploaded_at, file_name, downloads, upvotes, user_id FROM resources"

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
			&r.UserID,
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
func editHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, loggedIn := getSessionUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/edit/")
	var id int
	fmt.Sscan(idStr, &id)
	if id == 0 {
		http.NotFound(w, r)
		return
	}

	resource, err := getResourceByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if resource.UserID != currentUser.ID {
		http.Error(w, "You can only edit your own resources.", http.StatusForbidden)
		return
	}

	if r.Method == "GET" {
		data := newPageData(r, "Edit Resource - ElimuLocal")
		data.Resource = resource
		data.Universities = getUniversities()
		renderTemplate(w, "edit.html", data)
		return
	}

	if r.Method == "POST" {
		title := strings.TrimSpace(r.FormValue("title"))
		course := strings.TrimSpace(r.FormValue("course"))
		university := strings.TrimSpace(r.FormValue("university"))
		category := r.FormValue("category")
		description := strings.TrimSpace(r.FormValue("description"))

		if title == "" || course == "" || university == "" {
			data := newPageData(r, "Edit Resource - ElimuLocal")
			data.Message = "Please fill in all required fields."
			data.Resource = resource
			data.Universities = getUniversities()
			renderTemplate(w, "edit.html", data)
			return
		}

		fileName := resource.FileName

		r.ParseMultipartForm(500 << 20)
		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			ext := strings.ToLower(filepath.Ext(header.Filename))
			if ext == ".pdf" || ext == ".mp4" || ext == ".mkv" {
				newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
				newFilePath := filepath.Join("uploads", newFileName)
				dst, err := os.Create(newFilePath)
				if err == nil {
					defer dst.Close()
					io.Copy(dst, file)
					if fileName != "" {
						deleteFile(resource.FileName)
					}
					fileName = newFileName
				}
			}
		}

		_, err = db.Exec(
			"UPDATE resources SET title=?, course=?, university=?, category=?, description=?, file_name=? WHERE id=? AND user_id=?",
			title, course, university, category, description, fileName, id, currentUser.ID,
		)
		if err != nil {
			data := newPageData(r, "Edit Resource - ElimuLocal")
			data.Message = "Could not save changes. Please try again."
			data.Resource = resource
			data.Universities = getUniversities()
			renderTemplate(w, "edit.html", data)
			return
		}

		http.Redirect(w, r, "/?success=1", http.StatusSeeOther)
		return
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, loggedIn := getSessionUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
	var id int
	fmt.Sscan(idStr, &id)
	if id == 0 {
		http.NotFound(w, r)
		return
	}

	resource, err := getResourceByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if resource.UserID != currentUser.ID {
		http.Error(w, "You can only delete your own resources.", http.StatusForbidden)
		return
	}

	if resource.FileName != "" {
		deleteFile(resource.FileName)
	}

	db.Exec("DELETE FROM resources WHERE id = ? AND user_id = ?", id, currentUser.ID)

	http.Redirect(w, r, "/?deleted=1", http.StatusSeeOther)
}

func getResourceByID(id int) (Resource, error) {
	var r Resource
	err := db.QueryRow(
		"SELECT id, title, course, university, category, description, uploaded_by, uploaded_at, file_name, downloads, upvotes, user_id FROM resources WHERE id = ?",
		id,
	).Scan(
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
		&r.UserID,
	)
	return r, err
}

func newPageData(r *http.Request, title string) PageData {
	currentUser, loggedIn := getSessionUser(r)
	return PageData{
		Title:        title,
		Universities: getUniversities(),
		CurrentUser:  currentUser,
		LoggedIn:     loggedIn,
	}
}

func renderLanding(w http.ResponseWriter, data PageData) {
	tmpl, err := template.ParseFiles("templates/landing.html")
	if err != nil {
		http.Error(w, "Could not load page: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "landing", data)
	if err != nil {
		http.Error(w, "Could not render page: "+err.Error(), http.StatusInternalServerError)
	}
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := newPageData(r, "ElimuLocal — Free Study Materials for Kenyan University Students")
	renderLanding(w, data)
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
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
		message = "✅ Resource uploaded successfully!"
	}
	if r.URL.Query().Get("deleted") == "1" {
		message = "🗑️ Resource deleted successfully."
	}
	if r.URL.Query().Get("success") == "registered" {
		message = "🎉 Welcome to ElimuLocal! You are now registered and logged in."
	}

	data := newPageData(r, "Browse Resources - ElimuLocal")
	data.Resources = resources
	data.Universities = getUniversities()
	data.Search = search
	data.University = university
	data.Category = category
	data.Sort = sort
	data.Message = message

	renderTemplate(w, "home.html", data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, loggedIn := getSessionUser(r)

	if r.Method == "GET" {
		data := PageData{
			Title:        "Share a Resource - ElimuLocal",
			Universities: getUniversities(),
			CurrentUser:  currentUser,
			LoggedIn:     loggedIn,
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
				CurrentUser:  currentUser,
				LoggedIn:     loggedIn,
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		// allow larger uploads for video files (up to 500MB)
		err := r.ParseMultipartForm(500 << 20)
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "File too large. Maximum size is 500MB.",
				Universities: getUniversities(),
				CurrentUser:  currentUser,
				LoggedIn:     loggedIn,
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Please select a file to upload.",
				Universities: getUniversities(),
				CurrentUser:  currentUser,
				LoggedIn:     loggedIn,
			}
			renderTemplate(w, "upload.html", data)
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".pdf" && ext != ".mp4" && ext != ".mkv" && ext != ".webm" {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Only PDF and video files (MP4, MKV, WEBM) are allowed.",
				Universities: getUniversities(),
				CurrentUser:  currentUser,
				LoggedIn:     loggedIn,
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)

		err = uploadFile(fileName, file, "application/octet-stream")
		if err != nil {
			data := newPageData(r, "Share a Resource - ElimuLocal")
			data.Message = "Could not save file. Please try again."
			data.Universities = getUniversities()
			renderTemplate(w, "upload.html", data)
			return
		}
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Could not save file. Please try again.",
				Universities: getUniversities(),
				CurrentUser:  currentUser,
				LoggedIn:     loggedIn,
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
			UserID:      currentUser.ID,
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

		http.Redirect(w, r, "/browse?success=1", http.StatusSeeOther)
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

	var id int
	fmt.Sscan(idStr, &id)
	incrementDownloads(id)

	w.Header().Set("Content-Disposition", "attachment; filename="+title+".pdf")
	w.Header().Set("Content-Type", "application/pdf")

	err = serveFile(fileName, w)
	if err != nil {
		http.Error(w, "Could not retrieve file", http.StatusInternalServerError)
		return
	}
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

func previewHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/preview/")
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

	resource, err := getResourceByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if resource.FileName == "" {
		http.Error(w, "No file available for this resource yet.", http.StatusNotFound)
		return
	}
	ext := strings.ToLower(filepath.Ext(resource.FileName))
	isVideo := ext == ".mp4" || ext == ".mkv" || ext == ".webm"

	mimeType := "application/pdf"
	if ext == ".mp4" {
		mimeType = "video/mp4"
	} else if ext == ".mkv" {
		mimeType = "video/x-matroska"
	} else if ext == ".webm" {
		mimeType = "video/webm"
	}

	data := newPageData(r, resource.Title+" - ElimuLocal")
	data.Resource = resource
	data.IsVideo = isVideo
	data.MimeType = mimeType

	renderTemplate(w, "preview.html", data)
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/stream/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	var fileName string
	var title string
	err := db.QueryRow(
		"SELECT file_name, title FROM resources WHERE id = ?", idStr,
	).Scan(&fileName, &title)

	if err != nil || fileName == "" {
		http.NotFound(w, r)
		return
	}

	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".pdf":
		w.Header().Set("Content-Type", "application/pdf")
	case ".mp4":
		w.Header().Set("Content-Type", "video/mp4")
	case ".mkv":
		w.Header().Set("Content-Type", "video/x-matroska")
	case ".webm":
		w.Header().Set("Content-Type", "video/webm")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Header().Set("Content-Disposition", "inline; filename=\""+title+"\"")
	err = serveFile(fileName, w)
	if err != nil {
		http.Error(w, "Could not retrieve file", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Load .env file (silently ignored if not present)
	_ = godotenv.Load()

	initDB()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/browse", browseHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/preview/", previewHandler)
	http.HandleFunc("/stream/", streamHandler)
	http.HandleFunc("/upvote/", upvoteHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("ElimuLocal is running!")
	fmt.Printf("Open your browser and go to: http://localhost:%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server.")

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
