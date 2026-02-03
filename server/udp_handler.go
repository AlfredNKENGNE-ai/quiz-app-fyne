package server

import (
	"encoding/json"
	"log"
	"net"
	"quiz-app-fyne/shared"
)

// HandleMessage traite tous les messages UDP entrants
func HandleMessage(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	var msg shared.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println("❌ JSON invalide :", err)
		return
	}

	log.Printf("📩 Message reçu de %s → %s", addr.String(), msg.Type)

	switch msg.Type {

	case shared.MsgLogin:
		payload := msg.Payload.(map[string]interface{})
		email := payload["email"].(string)
		user, err := DB.GetUserByEmail(email)
		if err != nil {
			SendResponse(conn, addr, shared.Message{Type: shared.MsgLoginError})
			return
		}
		// Stocker l'adresse UDP dans une nouvelle instance
		loggedInUser := &shared.User{
			ID:           user.ID,
			Email:        user.Email,
			Username:     user.Username,
			PasswordHash: user.PasswordHash,
			TotalScore:   user.TotalScore,
			GamesPlayed:  user.GamesPlayed,
			CreatedAt:    user.CreatedAt,
			LastLogin:    user.LastLogin,
			Addr:         addr, // ICI on assigne l'adresse
		}

		SendResponse(conn, addr, shared.Message{
			Type: shared.MsgLoginOK,
			Payload: shared.LoginOKPayload{
				UserID: loggedInUser.ID,
				Email:  loggedInUser.Email,
			},
		})

	case shared.MsgCreateGame:
		payload := msg.Payload.(map[string]interface{})
		userID := int(payload["user_id"].(float64))
		mode := payload["mode"].(string)

		user, _ := DB.GetUserByID(userID)
		user.Addr = addr // ✅ TRÈS IMPORTANT

		game := Manager.CreateGame(user)
		game.Mode = mode

		SendResponse(conn, addr, shared.Message{
			Type: shared.MsgCreateGame,
			Payload: map[string]string{
				"game_code": game.Code,
				"mode":      mode,
			},
		})

		if mode == "solo" {
			Manager.StartGame(game.Code)
			Manager.Conn = conn
			go Manager.RunGame(conn, game.Code)
		}
		if mode == "multi" {
			Manager.MonitorLobby(game, conn)
		}

	case shared.MsgJoinGame:
		payload := msg.Payload.(map[string]interface{})
		userID := int(payload["user_id"].(float64))
		gameCode := payload["game_code"].(string)
		user, err := DB.GetUserByID(userID)
		if err != nil {
			log.Println("⚠️ Utilisateur introuvable")
			return
		}
		user.Addr = addr
		game, err := Manager.JoinGame(gameCode, user)
		if err != nil {
			log.Println("⚠️ Impossible de rejoindre la partie:", err)
			return
		}
		Manager.MonitorLobby(game, conn)
		log.Printf("✅ Joueur %s a rejoint la partie %s", user.Email, gameCode)

	case shared.MsgStartGame:
		payload := msg.Payload.(map[string]interface{})
		gameCode := payload["game_code"].(string)
		err := Manager.StartGame(gameCode)
		if err != nil {
			log.Println("❌ Impossible de démarrer la partie :", err)
			return
		}
		log.Println("🚀 Partie démarrée :", gameCode)
		go Manager.RunGame(conn, gameCode)

	case shared.MsgAnswer:
		payload := msg.Payload.(map[string]interface{})
		userID := int(payload["user_id"].(float64))
		questionID := int(payload["question_id"].(float64))
		choice := int(payload["choice"].(float64))
		Manager.ProcessAnswer(userID, questionID, choice)

	case shared.MsgRequestRiddleHint:
		payload := msg.Payload.(map[string]interface{})
		userID := int(payload["user_id"].(float64))
		hintType := int(payload["hint_type"].(float64)) // 1 ou 2
		Manager.SendRiddleHint(conn, userID, hintType, addr)

	case shared.MsgRiddleAnswer:
		payload := msg.Payload.(map[string]interface{})
		userID := int(payload["user_id"].(float64))
		answer := payload["answer"].(string)
		Manager.ProcessRiddleAnswer(userID, answer)

	default:
		log.Println("⚠️ Type de message inconnu :", msg.Type)
	}

}

// SendResponse envoie un message UDP au client
func SendResponse(conn *net.UDPConn, addr *net.UDPAddr, msg shared.Message) {
	data, _ := json.Marshal(msg)
	_, err := conn.WriteToUDP(data, addr)
	if err != nil {
		log.Println("❌ Erreur envoi UDP :", err)
	}
}
