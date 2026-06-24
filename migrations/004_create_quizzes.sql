
CREATE TABLE IF NOT EXISTS quizzes (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    course      TEXT NOT NULL,
    university  TEXT NOT NULL,
    description TEXT,
    created_by  INTEGER NOT NULL,
    created_at  TEXT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS quiz_questions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_id         INTEGER NOT NULL,
    question_text   TEXT NOT NULL,
    question_type   TEXT NOT NULL,   -- 'mcq' or 'true_false'
    option_a        TEXT,
    option_b        TEXT,
    option_c        TEXT,
    option_d        TEXT,
    correct_answer  TEXT NOT NULL,   -- 'a', 'b', 'c', 'd' or 'true'/'false'
    question_order  INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id)
);

CREATE TABLE IF NOT EXISTS quiz_attempts (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_id         INTEGER NOT NULL,
    user_id         INTEGER NOT NULL,
    score           INTEGER NOT NULL DEFAULT 0,
    total_questions INTEGER NOT NULL DEFAULT 0,
    submitted_at    TEXT NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS attempt_answers (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    attempt_id      INTEGER NOT NULL,
    question_id     INTEGER NOT NULL,
    student_answer  TEXT NOT NULL,
    is_correct      INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (attempt_id) REFERENCES quiz_attempts(id),
    FOREIGN KEY (question_id) REFERENCES quiz_questions(id)
);