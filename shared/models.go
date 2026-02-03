package shared

import (
	"net"
	"time"
)

// =====================
// UTILISATEUR
// =====================
type User struct {
	ID           int          `json:"id"`
	Email        string       `json:"email"`
	Username     string       `json:"username"`
	PasswordHash string       `json:"password_hash"`
	TotalScore   int          `json:"total_score"`
	GamesPlayed  int          `json:"games_played"`
	CreatedAt    time.Time    `json:"created_at"`
	LastLogin    *time.Time   `json:"last_login"`
	Addr         *net.UDPAddr `json:"-"`         // Adresse UDP du joueur (non sérialisée en JSON)
	GameCode     string       `json:"game_code"` // Code de la partie en cours
}

// =====================
// QUESTION (QCM)
// =====================
type Question struct {
	ID              int    `json:"id"`
	QuestionText    string `json:"question_text"`
	ChoiceA         string `json:"choice_a"`
	ChoiceB         string `json:"choice_b"`
	ChoiceC         string `json:"choice_c"`
	ChoiceD         string `json:"choice_d"`
	CorrectAnswer   string `json:"correct_answer"` // "A", "B", "C", "D"
	DifficultyLevel int    `json:"difficulty_level"`
	Manche          int    `json:"manche"`
	Category        string `json:"category"`
}

// =====================
// DEVINETTE (Riddle)
// =====================
type Riddle struct {
	ID              int    `json:"id"`
	RiddleText      string `json:"riddle_text"`
	CorrectWord     string `json:"correct_word"`
	HintLevel1      string `json:"hint_level1"`
	HintLevel2      string `json:"hint_level2"`
	DifficultyLevel int    `json:"difficulty_level"`
}
