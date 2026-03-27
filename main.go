package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Resource struct {
	ID			int
	Title		string
	Course		string
	University	string
	Category	string
	Description	string
	UploadedBy	string
	UploadedAt	string
	FileName	string
	Downloads	int
}
type PageData struct {
	Title			string
	Resources		[]Resource
	Universities	[]string
	Message			string
}

var resources = []Resource{
	{
		ID:          1,
		Title:       "Introduction to Data Structures - Full Notes",
		Course:      "Computer Science",
		University:  "University of Nairobi",
		Category:    "Notes",
		Description: "Covers arrays, linked lists, stacks, queues and trees.",
		UploadedBy:  "Alice M.",
		UploadedAt:  "2025-03-01",
		FileName:    "",
		Downloads:   0,
	},
	{
		ID:          2,
		Title:       "Principles of Economics - Chapter Summaries",
		Course:      "Economics",
		University:  "Strathmore University",
		Category:    "Notes",
		Description: "Clear summaries of micro and macroeconomics fundamentals.",
		UploadedBy:  "David O.",
		UploadedAt:  "2025-03-10",
		FileName:    "",
		Downloads:   0,
	},
}

var nextID = 3

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	search := r.URL.Query().Get("search")
	university := r.URL.Query().Get("university")

	filtered := filterResources(search, university)

	data := PageData{
		Title:        "ElimuLocal - Browse Study Materials",
		Resources:    filtered,
		Universities: getUniversities(),
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
		_, _, err := r.FormFile("file")
		if err != nil {
			data := PageData{
				Title:        "Share a Resource - ElimuLocal",
				Message:      "Please select a PDF file to upload.",
				Universities: getUniversities(),
			}
			renderTemplate(w, "upload.html", data)
			return
		}
		newResource := Resource{
			ID:          nextID,
			Title:       title,
			Course:      course,
			University:  university,
			Category:    category,
			Description: description,
			UploadedBy:  uploader,
			UploadedAt:  time.Now().Format("2006-01-02"),

			FileName: "",
			Downloads: 0,
		}
		resources = append(resources, newResource)
		nextID++
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
func filterResources(search, university string) []Resource {
	if search == "" && university == "" {
		return resources
	}

	var filtered []Resource

	for _, r := range resources {

		matchesSearch := true
		matchesUniversity :=true

		if search != "" {
			s := strings.ToLower(search)
			matchesSearch = strings.Contains(strings.ToLower(r.Title), s) ||
				strings.Contains(strings.ToLower(r.Course), s) ||
				strings.Contains(strings.ToLower(r.Description), s)
		}
		if university != "" {
			matchesUniversity = strings.EqualFold(r.University, university)
		}
		if matchesSearch && matchesUniversity {
			filtered = append(filtered, r)
		}
	}

	return filtered
}
func getUniversities() []string {
	seen := make(map[string]bool)

	var unis []string

	for _, r := range resources {
		name := strings.TrimSpace(r.University)
		if name != "" && !seen[name] {
			seen[name] = true
			unis = append(unis, name)
		}
	}

	return unis
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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("ElimuLocal is running!")
	fmt.Println("Open your browser and go to: http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop the server.")

	log.Fatal(http.ListenAndServe(":8080", nil))
}



