package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"quiz-app-fyne/shared"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var serverAddr *net.UDPAddr
var conn *net.UDPConn

func InitNetwork() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}
	serverAddr = addr

	conn, err = net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	go ListenServer()
}

func ListenServer() {
	buffer := make([]byte, 4096)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		var msg shared.Message
		json.Unmarshal(buffer[:n], &msg)

		switch msg.Type {

		case shared.MsgLoginOK:
			data, _ := json.Marshal(msg.Payload)
			var payload shared.LoginOKPayload
			json.Unmarshal(data, &payload)

			CurrentUser = &shared.User{
				ID:    payload.UserID,
				Email: payload.Email,
			}

			ShowModeSelectionScreen()

		case shared.MsgCreateGame:
			data, _ := json.Marshal(msg.Payload)
			var payload map[string]string
			json.Unmarshal(data, &payload)

			CurrentUser.GameCode = payload["game_code"]

			if payload["mode"] == "multi" {
				ShowLobbyWithGameCode(payload["game_code"])
			}

		case shared.MsgQuestion:
			data, _ := json.Marshal(msg.Payload)
			var qp shared.QuestionPayload
			json.Unmarshal(data, &qp)

			ShowQuestionScreen(
				qp.Question.Text,
				qp.Question.Options,
				qp.Question.ID,
			)

		case shared.MsgRiddle:
			data, _ := json.Marshal(msg.Payload)
			var rp shared.RiddlePayload
			json.Unmarshal(data, &rp)

			ShowRiddleScreen(rp.Text)

		case shared.MsgGameOver:
			data, _ := json.Marshal(msg.Payload)
			var gp shared.GameOverPayload
			json.Unmarshal(data, &gp)

			var results []string
			for _, r := range gp.Results {
				results = append(results, fmt.Sprintf("%s : %d", r.Email, r.Score))
			}
			ShowResults(results)
		case "GAME_START":
			ShowWaitingRoom()
		}

	}
}

func send(msg shared.Message) {
	data, _ := json.Marshal(msg)
	_, err := conn.Write(data)
	if err != nil {
		log.Println("Erreur envoi :", err)
	}
}

func SendLogin(email, password string) {
	send(shared.Message{
		Type: shared.MsgLogin,
		Payload: shared.LoginPayload{
			Email:    email,
			Password: password,
		},
	})
}

func SendCreateGame(userID int, mode string) {
	send(shared.Message{
		Type: shared.MsgCreateGame,
		Payload: map[string]interface{}{
			"user_id": userID,
			"mode":    mode,
		},
	})
}

func SendJoinGame(code string, userID int) {
	send(shared.Message{
		Type: shared.MsgJoinGame,
		Payload: map[string]interface{}{
			"user_id":   userID,
			"game_code": code,
		},
	})
}

func SendAnswer(questionID int, choice int) {
	send(shared.Message{
		Type: shared.MsgAnswer,
		Payload: map[string]interface{}{
			"user_id":     CurrentUser.ID,
			"question_id": questionID,
			"choice":      choice,
		},
	})
}

func SendRiddleAnswer(text string) {
	send(shared.Message{
		Type: shared.MsgRiddleAnswer,
		Payload: shared.RiddleAnswerPayload{
			UserID: CurrentUser.ID,
			Answer: text,
		},
	})
}

func RequestHint(level int) {
	send(shared.Message{
		Type: shared.MsgRequestRiddleHint,
		Payload: map[string]interface{}{
			"user_id":   CurrentUser.ID,
			"hint_type": level,
		},
	})
}

var waitLabel *widget.Label

func ShowLobbyWithGameCode(code string) {
	codeLabel := widget.NewLabel("Code de la salle : " + code)
	waitLabel = widget.NewLabel("En attente des joueurs...")

	MainWindow.SetContent(
		container.NewVBox(
			widget.NewLabel("🎮 Lobby"),
			codeLabel,
			waitLabel,
		),
	)
}
