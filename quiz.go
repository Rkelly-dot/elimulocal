package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"html/template"
)

// -------------------------------------------------------
// STRUCTS
// -------------------------------------------------------

type Quiz struct {
	ID          int
	Title       string
	Course      string
	University  string
	Description string
	CreatedBy   int
	CreatedAt   string
}

type QuizQuestion struct {
	ID            int
	QuizID        int
	QuestionText  string
	QuestionType  string // "mcq" or "true_false"
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectAnswer string
	QuestionOrder int
}

// QuizAttempt records one student's completed attempt at a quiz.
type QuizAttempt struct {
	ID             int
	QuizID         int
	UserID         int
	Score          int
	TotalQuestions int
	SubmittedAt    string
}

// QuizPageData is passed to quiz-related templates.
// Embeds the standard session fields via newPageData pattern.
type QuizPageData struct {
	Title        string
	CurrentUser  User
	LoggedIn     bool
	Message      string
	Quiz         Quiz
	Questions    []QuizQuestion
	Quizzes      []Quiz
	Universities []string
	Attempt      QuizAttempt
	Results      []QuestionResult
}

// QuestionResult is used on the results page to show each
// question alongside the student's answer and whether it was correct.
type QuestionResult struct {
	Question      QuizQuestion
	StudentAnswer string
	IsCorrect     bool
}

// -------------------------------------------------------
// HELPER — newQuizPageData mirrors newPageData() from auth.go
// -------------------------------------------------------

func newQuizPageData(r *http.Request, title string) QuizPageData {
	currentUser, loggedIn := getSessionUser(r)
	return QuizPageData{
		Title:       title,
		CurrentUser: currentUser,
		LoggedIn:    loggedIn,
	}
}

// -------------------------------------------------------
// DATABASE HELPERS
// -------------------------------------------------------

// saveQuiz inserts a new quiz and returns its generated ID.
func saveQuiz(q Quiz) (int, error) {
	result, err := db.Exec(
		`INSERT INTO quizzes (title, course, university, description, created_by, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		q.Title, q.Course, q.University, q.Description, q.CreatedBy, q.CreatedAt,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// saveQuestion inserts one question belonging to a quiz.
func saveQuestion(qq QuizQuestion) error {
	_, err := db.Exec(
		`INSERT INTO quiz_questions
		 (quiz_id, question_text, question_type, option_a, option_b, option_c, option_d, correct_answer, question_order)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		qq.QuizID, qq.QuestionText, qq.QuestionType,
		qq.OptionA, qq.OptionB, qq.OptionC, qq.OptionD,
		qq.CorrectAnswer, qq.QuestionOrder,
	)
	return err
}

// -------------------------------------------------------
// HANDLER — createQuizHandler
// GET  shows the create-quiz form
// POST saves the quiz and its questions, then redirects
// -------------------------------------------------------

func createQuizHandler(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := getSessionUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		data := newQuizPageData(r, "Create a Quiz - ElimuLocal")
		data.Universities = getUniversities()
		renderQuizTemplate(w, "create_quiz.html", data)
		return
	}

	if r.Method == "POST" {
		currentUser, _ := getSessionUser(r)

		title := strings.TrimSpace(r.FormValue("title"))
		course := strings.TrimSpace(r.FormValue("course"))
		university := strings.TrimSpace(r.FormValue("university"))
		description := strings.TrimSpace(r.FormValue("description"))

		if title == "" || course == "" || university == "" {
			data := newQuizPageData(r, "Create a Quiz - ElimuLocal")
			data.Message = "Please fill in the title, course and university."
			data.Universities = getUniversities()
			renderQuizTemplate(w, "create_quiz.html", data)
			return
		}

		quizID, err := saveQuiz(Quiz{
			Title:       title,
			Course:      course,
			University:  university,
			Description: description,
			CreatedBy:   currentUser.ID,
			CreatedAt:   time.Now().Format("2006-01-02"),
		})
		if err != nil {
			data := newQuizPageData(r, "Create a Quiz - ElimuLocal")
			data.Message = "Could not create quiz. Please try again."
			data.Universities = getUniversities()
			renderQuizTemplate(w, "create_quiz.html", data)
			return
		}

		// Questions arrive as parallel form arrays:
		// question_text[], question_type[], option_a[], ... correct_answer[]
		// This lets one form submit any number of questions at once.
		questionTexts := r.Form["question_text[]"]
		questionTypes := r.Form["question_type[]"]
		optionsA := r.Form["option_a[]"]
		optionsB := r.Form["option_b[]"]
		optionsC := r.Form["option_c[]"]
		optionsD := r.Form["option_d[]"]
		correctAnswers := r.Form["correct_answer[]"]

		for i := range questionTexts {
			if strings.TrimSpace(questionTexts[i]) == "" {
				continue // skip blank rows
			}

			qq := QuizQuestion{
				QuizID:        quizID,
				QuestionText:  strings.TrimSpace(questionTexts[i]),
				QuestionType:  questionTypes[i],
				CorrectAnswer: strings.ToLower(strings.TrimSpace(correctAnswers[i])),
				QuestionOrder: i,
			}

			if i < len(optionsA) {
				qq.OptionA = optionsA[i]
			}
			if i < len(optionsB) {
				qq.OptionB = optionsB[i]
			}
			if i < len(optionsC) {
				qq.OptionC = optionsC[i]
			}
			if i < len(optionsD) {
				qq.OptionD = optionsD[i]
			}

			saveQuestion(qq)
		}

		http.Redirect(w, r, "/quizzes?created=1", http.StatusSeeOther)
		return
	}
}

// renderQuizTemplate mirrors renderTemplate() but for quiz pages,
// using base.html as the shared layout.
func renderQuizTemplate(w http.ResponseWriter, page string, data QuizPageData) {
	funcMap := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}

	tmpl, err := template.New("base.html").Funcs(funcMap).ParseFiles("templates/base.html", "templates/"+page)
	if err != nil {
		http.Error(w, "Could not load page: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Could not render page: "+err.Error(), http.StatusInternalServerError)
	}
}

// -------------------------------------------------------
// HANDLER — viewQuizHandler
// Shows a quiz to a student so they can answer it.
// URL pattern: /quiz/{id}
// -------------------------------------------------------

func viewQuizHandler(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := getSessionUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/quiz/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		http.NotFound(w, r)
		return
	}

	var quiz Quiz
	err = db.QueryRow(
		"SELECT id, title, course, university, description, created_by, created_at FROM quizzes WHERE id = ?",
		id,
	).Scan(&quiz.ID, &quiz.Title, &quiz.Course, &quiz.University, &quiz.Description, &quiz.CreatedBy, &quiz.CreatedAt)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	rows, err := db.Query(
		`SELECT id, quiz_id, question_text, question_type, option_a, option_b, option_c, option_d, correct_answer, question_order
		 FROM quiz_questions WHERE quiz_id = ? ORDER BY question_order ASC`,
		id,
	)
	if err != nil {
		http.Error(w, "Could not load questions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var questions []QuizQuestion
	for rows.Next() {
		var q QuizQuestion
		err := rows.Scan(
			&q.ID, &q.QuizID, &q.QuestionText, &q.QuestionType,
			&q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
			&q.CorrectAnswer, &q.QuestionOrder,
		)
		if err != nil {
			continue
		}
		questions = append(questions, q)
	}

	data := newQuizPageData(r, quiz.Title+" - ElimuLocal")
	data.Quiz = quiz
	data.Questions = questions
	renderQuizTemplate(w, "view_quiz.html", data)
}
// -------------------------------------------------------
// HANDLER — listQuizzesHandler
// Shows every quiz available, similar to browseHandler
// for resources. Students click into one to take it.
// -------------------------------------------------------

func listQuizzesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(
		`SELECT id, title, course, university, description, created_by, created_at
		 FROM quizzes ORDER BY id DESC`,
	)
	if err != nil {
		http.Error(w, "Could not load quizzes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var quizzes []Quiz
	for rows.Next() {
		var q Quiz
		err := rows.Scan(&q.ID, &q.Title, &q.Course, &q.University, &q.Description, &q.CreatedBy, &q.CreatedAt)
		if err != nil {
			continue
		}
		quizzes = append(quizzes, q)
	}

	data := newQuizPageData(r, "Quizzes - ElimuLocal")
	data.Quizzes = quizzes

	if r.URL.Query().Get("created") == "1" {
		data.Message = "✅ Quiz created successfully!"
	}

	renderQuizTemplate(w, "list_quizzes.html", data)
}
