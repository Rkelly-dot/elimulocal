-- Migration 002: Create resources table
CREATE TABLE IF NOT EXISTS resources (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT    NOT NULL,
    course      TEXT    NOT NULL,
    university  TEXT    NOT NULL,
    category    TEXT    NOT NULL,
    description TEXT,
    uploaded_by TEXT,
    uploaded_at TEXT,
    file_name   TEXT,
    downloads   INTEGER NOT NULL DEFAULT 0,
    upvotes     INTEGER NOT NULL DEFAULT 0,
    user_id     INTEGER NOT NULL DEFAULT 0
);
