package main

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("elimulocal-secret-key-change-in-production"))

type User struct {
	ID	   			int
	Username 		string
	Email			string
	PasswordHash	string
	CreatedAt		string
}

func createUsersTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id		      INTEGER PRIMARY  KEY AUTOINCREMENT,
		username	  TEXT NOT NULL UNIQUE,
		email		  TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at`   TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		panic("COuld not create users table: " + err.Error())
	}
}
	


