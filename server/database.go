package server

import (
	"database/sql"
	"log"
	"quiz-app-fyne/shared"

	_ "github.com/mattn/go-sqlite3"
)

var DB *Database

type Database struct {
	usersDB *sql.DB
	quizDB  *sql.DB
}

// INITIALISATION DB
func InitDatabases() {
	var err error
	DB, err = NewDatabase("server/databases/users.db", "server/databases/quiz_data.db")
	if err != nil {
		log.Fatalf("❌ Erreur initialisation DB : %v", err)
	}
	log.Println("✅ Bases de données initialisées")
}

func NewDatabase(usersPath, quizPath string) (*Database, error) {
	usersDB, err := sql.Open("sqlite3", usersPath)
	if err != nil {
		return nil, err
	}
	quizDB, err := sql.Open("sqlite3", quizPath)
	if err != nil {
		usersDB.Close()
		return nil, err
	}
	if err := usersDB.Ping(); err != nil {
		usersDB.Close()
		quizDB.Close()
		return nil, err
	}
	if err := quizDB.Ping(); err != nil {
		usersDB.Close()
		quizDB.Close()
		return nil, err
	}
	return &Database{usersDB: usersDB, quizDB: quizDB}, nil
}

// UTILISATEURS
func (db *Database) GetUserByEmail(email string) (*shared.User, error) {
	row := db.usersDB.QueryRow(`SELECT id, email, username, password_hash, total_score, games_played, created_at, last_login FROM users WHERE email=?`, email)
	user := &shared.User{}
	var lastLogin sql.NullTime
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.TotalScore, &user.GamesPlayed, &user.CreatedAt, &lastLogin)
	if err != nil {
		return nil, err
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	return user, nil
}

func (db *Database) GetUserByID(id int) (*shared.User, error) {
	row := db.usersDB.QueryRow(`SELECT id, email, username, password_hash, total_score, games_played, created_at, last_login FROM users WHERE id=?`, id)
	user := &shared.User{}
	var lastLogin sql.NullTime
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.TotalScore, &user.GamesPlayed, &user.CreatedAt, &lastLogin)
	if err != nil {
		return nil, err
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}
	return user, nil
}

// QUESTIONS
func (db *Database) GetQuestionsByLevelAndManche(level, manche, limit int) ([]shared.Question, error) {
	rows, err := db.quizDB.Query(`SELECT id, question_text, choice_a, choice_b, choice_c, choice_d, correct_answer, difficulty_level, manche, category
		FROM questions WHERE difficulty_level=? AND manche=? ORDER BY RANDOM() LIMIT ?`, level, manche, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []shared.Question
	for rows.Next() {
		var q shared.Question
		err := rows.Scan(&q.ID, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &q.CorrectAnswer, &q.DifficultyLevel, &q.Manche, &q.Category)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// Manche 1 à 8 QCM (4 niveau1 + 4 niveau2)
func (db *Database) GetRandomQuestionsForManche1() ([]shared.Question, error) {
	q1, err := db.GetQuestionsByLevelAndManche(1, 1, 4)
	if err != nil {
		return nil, err
	}
	q2, err := db.GetQuestionsByLevelAndManche(2, 1, 4)
	if err != nil {
		return nil, err
	}
	return append(q1, q2...), nil
}

// GetRandomQuestionsForManche2 - Questions pour la manche 2 (contre-la-montre)
func (db *Database) GetRandomQuestionsForManche2() ([]shared.Question, error) {
	return db.GetQuestionsByLevelAndManche(3, 2, 60)
}

// Manche 3 : devinette
func (db *Database) GetRandomRiddle() (*shared.Riddle, error) {
	row := db.quizDB.QueryRow(`SELECT id, riddle_text, correct_word, hint_level1, hint_level2, difficulty_level FROM riddles ORDER BY RANDOM() LIMIT 1`)
	r := &shared.Riddle{}
	err := row.Scan(&r.ID, &r.RiddleText, &r.CorrectWord, &r.HintLevel1, &r.HintLevel2, &r.DifficultyLevel)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateUserScore - Met à jour le score d'un utilisateur
func (db *Database) UpdateUserScore(userID, score int) error {
	_, err := db.usersDB.Exec(
		`UPDATE users SET total_score = total_score + ?, games_played = games_played + 1 WHERE id = ?`,
		score,
		userID,
	)
	return err
}
