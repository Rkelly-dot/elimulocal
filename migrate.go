package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// runMigrations applies all SQL files in the migrations/ directory that have
// not yet been recorded in the schema_migrations table. Files are run in
// lexicographic (numeric) order so that 001_, 002_, 003_, … are always
// applied in the correct sequence.
func runMigrations(db *sql.DB) {
	// Ensure the migrations tracking table exists.
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			filename   TEXT    NOT NULL UNIQUE,
			applied_at TEXT    NOT NULL DEFAULT (datetime('now'))
		);
	`)
	if err != nil {
		log.Fatalf("migrations: could not create schema_migrations table: %v", err)
	}

	// Read migration files.
	entries, err := os.ReadDir("migrations")
	if err != nil {
		// No migrations directory — nothing to do.
		if os.IsNotExist(err) {
			return
		}
		log.Fatalf("migrations: could not read migrations directory: %v", err)
	}

	// Collect and sort .sql files.
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		// Check if already applied.
		var count int
		err = db.QueryRow(
			"SELECT COUNT(*) FROM schema_migrations WHERE filename = ?", name,
		).Scan(&count)
		if err != nil {
			log.Fatalf("migrations: could not query schema_migrations: %v", err)
		}
		if count > 0 {
			continue // already applied
		}

		// Read SQL file.
		path := filepath.Join("migrations", name)
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("migrations: could not read %s: %v", name, err)
		}

		// Execute inside a transaction for atomicity.
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("migrations: could not begin transaction for %s: %v", name, err)
		}

		if _, err = tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			log.Fatalf("migrations: failed to apply %s: %v", name, err)
		}

		if _, err = tx.Exec(
			"INSERT INTO schema_migrations (filename) VALUES (?)", name,
		); err != nil {
			_ = tx.Rollback()
			log.Fatalf("migrations: could not record %s: %v", name, err)
		}

		if err = tx.Commit(); err != nil {
			log.Fatalf("migrations: could not commit %s: %v", name, err)
		}

		fmt.Printf("migrations: applied %s\n", name)
	}

	fmt.Println("migrations: all up to date")
}
