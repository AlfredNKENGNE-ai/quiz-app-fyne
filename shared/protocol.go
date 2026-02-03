package shared

// Types de messages UDP
const (
	MsgRegister          = "REGISTER"
	MsgLogin             = "LOGIN"
	MsgLoginOK           = "LOGIN_OK"
	MsgLoginError        = "LOGIN_ERROR"
	MsgCreateGame        = "CREATE_GAME"
	MsgJoinGame          = "JOIN_GAME"
	MsgStartGame         = "START_GAME"
	MsgGameOver          = "GAME_OVER"
	MsgQuestion          = "QUESTION"
	MsgAnswer            = "ANSWER"
	MsgScoreUpdate       = "SCORE_UPDATE"
	MsgRequestRiddleHint = "REQUEST_RIDDLE_HINT"
	MsgRiddleHint        = "RIDDLE_HINT"
	MsgRiddleAnswer      = "RIDDLE_ANSWER"
	MsgRiddle            = "RIDDLE"
)

// Message UDP générique
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// LOGIN
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginOKPayload struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

// QUESTION
type QuestionMessage struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Level   int      `json:"level"`
}
type QuestionPayload struct {
	Question QuestionMessage `json:"question"`
	Manche   int             `json:"manche"`
}

// REPONSES
type AnswerPayload struct {
	UserID     int `json:"user_id"`
	QuestionID int `json:"question_id"`
	Choice     int `json:"choice"`
}

// MULTIJOUEUR
type CreateGamePayload struct {
	UserID int `json:"user_id"`
}
type JoinGamePayload struct {
	UserID   int    `json:"user_id"`
	GameCode string `json:"game_code"`
}

// DEVINETTE
type RiddlePayload struct {
	RiddleID int    `json:"riddle_id"`
	Text     string `json:"text"`
}
type RiddleHintPayload struct {
	RiddleID int    `json:"riddle_id"`
	Text     string `json:"text"`
	Cost     int    `json:"cost"`
}
type RiddleAnswerPayload struct {
	UserID int    `json:"user_id"`
	Answer string `json:"answer"`
}

// SCORES ET RESULTATS
type PlayerResult struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Score  int    `json:"score"`
}
type GameOverPayload struct {
	Results []PlayerResult `json:"results"`
}
