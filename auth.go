package main

import (
	"encoding/gob"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("elimulocal-secret-key-change-in-production"))

func init() {
	gob.Register(int(0))

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 365,
		HttpOnly: true,
	}
}

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    string
}

func createUsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id		      INTEGER PRIMARY  KEY AUTOINCREMENT,
		username	  TEXT NOT NULL UNIQUE,
		email		  TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at    TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		panic("COuld not create users table: " + err.Error())
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getSessionUser(r *http.Request) (User, bool) {
	session, err := store.Get(r, "elimulocal-session")
	if err != nil {
		return User{}, false
	}

	var userID int

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

	if userID == 0 {
		return User{}, false
	}

	var u User
	err = db.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
	if err != nil {
		return User{}, false
	}

	return u, true
}

func setSessionUser(w http.ResponseWriter, r *http.Request, userID int) error {
	session, err := store.Get(r, "elimulocal-session")
	if err != nil {
		return err
	}
	session.Values["userID"] = userID
	session.Options.MaxAge = 86400 * 365
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "elimulocal-session")
	if err != nil {
		return
	}

	session.Options.MaxAge = -1
	session.Save(r, w)
}

func getUserByUsername(username string) (User, error) {
	var u User
	err := db.QueryRow(
		"SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?",
		username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

func usernameExists(username string) bool {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	return count > 0
}

func emailExists(email string) bool {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	return count > 0
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := getSessionUser(r)
	if loggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		data := PageData{Title: "Register - ElimuLocal"}
		renderTemplate(w, "register.html", data)
		return
	}
	if r.Method == "POST" {
		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		confirm := r.FormValue("confirm")

		if username == "" || email == "" || password == "" {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "Please fill in all fields.",
			}
			renderTemplate(w, "register.html", data)
			return
		}
		if len(username) < 3 {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "Username must be at least 3 characters.",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		if len(password) < 8 {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "Password must be at least 8 characters.",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		if password != confirm {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "Passwords do not match.",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		if usernameExists(username) {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "That username is already taken.",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		if emailExists(email) {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "That email is already registered.",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		hash, err := HashPassword(password)
		if err != nil {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "Something went wrong. Please try again.",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		result, err := db.Exec(
			"INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, ?)",
			username, email, hash, time.Now().Format("2006-01-02"),
		)
		if err != nil {
			data := PageData{
				Title:   "Register - ElimuLocal",
				Message: "Could not create account. Please try again.",
			}
			renderTemplate(w, "register.html", data)
			return
		}

		userID, _ := result.LastInsertId()
		setSessionUser(w, r, int(userID))

		http.Redirect(w, r, "/browse?success=registered", http.StatusSeeOther)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := getSessionUser(r)
	if loggedIn {
		http.Redirect(w, r, "/browse", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		data := PageData{Title: "Login - ElimuLocal"}
		renderTemplate(w, "login.html", data)
		return
	}

	if r.Method == "POST" {
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")

		if username == "" || password == "" {
			data := PageData{
				Title:   "Login - ElimuLocal",
				Message: "Please enter your username and password.",
			}
			renderTemplate(w, "login.html", data)
			return
		}
		user, err := getUserByUsername(username)
		if err != nil {
			data := PageData{
				Title:   "Login - ElimuLocal",
				Message: "Invalid username or password.",
			}
			renderTemplate(w, "login.html", data)
			return
		}

		if !CheckPassword(password, user.PasswordHash) {
			data := PageData{
				Title:   "Login - ElimuLocal",
				Message: "Invalid username or password.",
			}
			renderTemplate(w, "login.html", data)
			return
		}

		err = setSessionUser(w, r, user.ID)
		if err != nil {
			data := PageData{
				Title:   "Login - ElimuLocal",
				Message: "Login failed — could not create session: " + err.Error(),
			}
			renderTemplate(w, "login.html", data)
			return
		}
		http.Redirect(w, r, "/browse", http.StatusSeeOther)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
