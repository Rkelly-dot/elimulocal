package main

import (
	"net/http"

	"github.com/gorilla/sessions"

	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("elimulocal-secret-key-change-in-production"))

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

func CheckPasswordHash(password, hash string) bool {
	bytes, err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getSessionUser(r *http.Request) (User, bool) {
	session, err := store.Get(r, "elimulocal-session")
	if err != nil {
		return User{}, false
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok || userID == 0 {
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
	return session.Save(r, w)
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

		

