package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"quiz-app-fyne/shared"
	"sort"
	"sync"
	"time"
)

type Game struct {
	Code          string
	Players       map[int]*shared.User
	Questions     []shared.Question
	Riddle        *shared.Riddle
	Scores        map[int]int
	Mode          string
	AnswerChan    chan int
	Mutex         sync.Mutex
	CurrentManche int
	// ===== MANCHE 2 : COURSE CONTRE LA MONTRE =====
	Manche2Questions     []shared.Question
	Manche2StartTime     time.Time
	Manche2Duration      time.Duration
	CurrentQuestionIndex map[int]int
	StartTimerLaunched   bool
	RiddleAnswers        map[int]bool
}

type GameManager struct {
	Games map[string]*Game
	Mutex sync.RWMutex
	Conn  *net.UDPConn
}

var Manager = &GameManager{
	Games: make(map[string]*Game),
}

func (gm *GameManager) CreateGame(host *shared.User) *Game {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	var code string
	for {
		code = fmt.Sprintf("%04d", rand.Intn(10000))
		if _, exists := gm.Games[code]; !exists {
			break
		}
	}

	game := &Game{
		Code:                 code,
		Players:              make(map[int]*shared.User),
		Scores:               make(map[int]int),
		AnswerChan:           make(chan int, 10),
		CurrentQuestionIndex: make(map[int]int),
		RiddleAnswers:        make(map[int]bool),
	}

	game.Players[host.ID] = host
	game.Scores[host.ID] = 0
	gm.Games[code] = game

	log.Printf("✅ Partie créée %s host: %s", code, host.Email)
	return game
}

func (gm *GameManager) JoinGame(code string, player *shared.User) (*Game, error) {
	gm.Mutex.Lock()
	game, exists := gm.Games[code]
	gm.Mutex.Unlock()

	if !exists {
		return nil, fmt.Errorf("partie introuvable")
	}

	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	if _, alreadyJoined := game.Players[player.ID]; alreadyJoined {
		return game, nil
	}

	game.Players[player.ID] = player
	game.Scores[player.ID] = 0

	log.Printf("✅ Joueur %s (%d) a rejoint la partie %s", player.Email, player.ID, code)
	return game, nil
}

func (gm *GameManager) StartGame(code string) error {
	gm.Mutex.RLock()
	game, exists := gm.Games[code]
	gm.Mutex.RUnlock()

	if !exists {
		return fmt.Errorf("partie introuvable")
	}

	questionsManche1, err := DB.GetRandomQuestionsForManche1()
	if err != nil {
		return fmt.Errorf("échec chargement manche 1: %v", err)
	}
	game.Questions = questionsManche1

	riddle, err := DB.GetRandomRiddle()
	if err != nil {
		log.Printf("⚠️ Aucune devinette disponible: %v", err)
	} else {
		game.Riddle = riddle
	}

	log.Printf("🚀 Partie %s démarrée avec %d joueurs", code, len(game.Players))
	return nil
}

func (gm *GameManager) RunGame(conn *net.UDPConn, code string) {
	gm.Mutex.RLock()
	game, ok := gm.Games[code]
	gm.Mutex.RUnlock()

	if !ok {
		return
	}

	// sauvegarder la connexion dans GameManager pour l'utiliser ailleurs
	gm.Conn = conn

	// Manche 1 : QCM classique
	log.Printf("🎮 Partie %s - Début Manche 1 (QCM)", code)
	for i, q := range game.Questions {
		log.Printf("📝 Question %d/%d envoyée", i+1, len(game.Questions))
		gm.sendQuestionToAll(conn, game, q)
		gm.waitForAnswersOrTimeout(game, q.ID, 10*time.Second)
	}
	if false {
		// ==================
		// Manche 2 : QCM contre-la-montre 60s
		// ==================
		log.Printf("🎮 Partie %s - Début Manche 2 (Contre-la-montre)", code)

		questionsManche2, err := DB.GetRandomQuestionsForManche2()
		if err != nil {
			log.Println("❌ Impossible de charger la manche 2")
			return
		}

		game.Mutex.Lock()
		game.Manche2Questions = questionsManche2
		game.Manche2StartTime = time.Now()
		game.Manche2Duration = 60 * time.Second
		game.CurrentManche = 2
		game.CurrentQuestionIndex = make(map[int]int)
		for id := range game.Players {
			game.CurrentQuestionIndex[id] = 0
			// envoyer directement la première question
			gm.sendNextManche2Question(conn, game, id)
		}
		game.Mutex.Unlock()

		// boucle de timer pour manche 2
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			<-ticker.C
			if time.Since(game.Manche2StartTime) >= game.Manche2Duration {
				log.Printf("⏱️ Fin Manche 2")
				break
			}
		}
	}

	// Manche 3 : Devinette
	if game.Riddle != nil {
		log.Printf("🎮 Partie %s - Début Manche 3 (Devinette)", code)
		gm.sendRiddleToAll(conn, game)
		time.Sleep(60 * time.Second)
	}

	// Mise à jour des scores et fin de partie
	log.Printf("🏁 Partie %s terminée - Mise à jour des scores", code)
	for id, score := range game.Scores {
		if err := DB.UpdateUserScore(id, score); err != nil {
			log.Printf("❌ Erreur mise à jour score utilisateur %d: %v", id, err)
		}
	}

	gm.sendGameOver(conn, game)
	go gm.cleanupGame(code)
}

// ProcessAnswer corrigé pour manche 2
func (gm *GameManager) ProcessAnswer(userID, questionID, choice int) {
	gm.Mutex.RLock()
	defer gm.Mutex.RUnlock()

	var game *Game
	for _, g := range gm.Games {
		if _, ok := g.Players[userID]; ok {
			game = g
			break
		}
	}

	if game == nil {
		log.Printf("⚠️ Partie introuvable pour l'utilisateur %d", userID)
		return
	}

	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	for _, q := range game.Questions {
		if q.ID == questionID {
			var correct bool
			switch choice {
			case 0:
				correct = q.CorrectAnswer == "A"
			case 1:
				correct = q.CorrectAnswer == "B"
			case 2:
				correct = q.CorrectAnswer == "C"
			case 3:
				correct = q.CorrectAnswer == "D"
			}

			if q.Manche == 1 {
				if correct {
					game.Scores[userID] += 15
					log.Printf("✅ Joueur %d: +15 points (manche 1)", userID)
				}
			} else if q.Manche == 2 {
				if correct {
					game.Scores[userID] += 10
				} else {
					game.Scores[userID] -= 3
				}
				// avancer l'index et envoyer la prochaine question
				game.CurrentQuestionIndex[userID]++
				gm.sendNextManche2Question(gm.Conn, game, userID)
			}

			game.AnswerChan <- questionID
			return
		}
	}
}

// sendNextManche2Question corrigé
func (gm *GameManager) sendNextManche2Question(conn *net.UDPConn, game *Game, userID int) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	if time.Since(game.Manche2StartTime) >= game.Manche2Duration {
		// envoyer un signal "Manche terminée" au joueur
		msg := shared.Message{
			Type:    shared.MsgGameOver, // ou MsgManche2Finished
			Payload: map[string]string{"message": "Manche 2 terminée"},
		}
		if player, ok := game.Players[userID]; ok && player.Addr != nil {
			SendResponse(conn, player.Addr, msg)
		}
		return
	}

	index := game.CurrentQuestionIndex[userID]
	if index >= len(game.Manche2Questions) {
		return
	}

	q := game.Manche2Questions[index]
	payload := shared.QuestionPayload{
		Question: shared.QuestionMessage{
			ID:      q.ID,
			Text:    q.QuestionText,
			Options: []string{q.ChoiceA, q.ChoiceB, q.ChoiceC, q.ChoiceD},
			Level:   q.DifficultyLevel,
		},
		Manche: 2,
	}

	msg := shared.Message{
		Type:    shared.MsgQuestion,
		Payload: payload,
	}

	if player, ok := game.Players[userID]; ok && player.Addr != nil {
		SendResponse(conn, player.Addr, msg)
	}
}

func (gm *GameManager) sendQuestionToAll(
	conn *net.UDPConn,
	game *Game,
	q shared.Question,
) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	payload := shared.QuestionPayload{
		Question: shared.QuestionMessage{
			ID:      q.ID,
			Text:    q.QuestionText,
			Options: []string{q.ChoiceA, q.ChoiceB, q.ChoiceC, q.ChoiceD},
			Level:   q.DifficultyLevel,
		},
		Manche: q.Manche,
	}

	msg := shared.Message{
		Type:    shared.MsgQuestion,
		Payload: payload,
	}

	for _, player := range game.Players {
		if player.Addr != nil {
			SendResponse(conn, player.Addr, msg)
		} else {
			log.Printf("⚠️ Adresse UDP manquante pour le joueur %s", player.Email)
		}
	}
}

func (gm *GameManager) sendRiddleToAll(conn *net.UDPConn, game *Game) {
	payload := shared.RiddlePayload{
		RiddleID: game.Riddle.ID,
		Text:     game.Riddle.RiddleText,
	}

	msg := shared.Message{
		Type:    shared.MsgRiddle,
		Payload: payload,
	}

	for _, player := range game.Players {
		if player.Addr != nil {
			SendResponse(conn, player.Addr, msg)
		}
	}
}

func (gm *GameManager) sendGameOver(conn *net.UDPConn, game *Game) {
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	results := []shared.PlayerResult{}
	for id, score := range game.Scores {
		if player, exists := game.Players[id]; exists {
			results = append(results, shared.PlayerResult{
				UserID: id,
				Email:  player.Email,
				Score:  score,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	msg := shared.Message{
		Type: shared.MsgGameOver,
		Payload: shared.GameOverPayload{
			Results: results,
		},
	}

	log.Printf("🏆 Partie %s terminée - Classement:", game.Code)
	for i, result := range results {
		log.Printf("  %d. %s: %d points", i+1, result.Email, result.Score)
	}

	for _, player := range game.Players {
		if player.Addr != nil {
			SendResponse(conn, player.Addr, msg)
		} else {
			log.Printf("⚠️ Adresse UDP manquante pour le joueur %s", player.Email)
		}
	}
}

func (gm *GameManager) SendRiddleHint(conn *net.UDPConn, userID, hintType int, addr *net.UDPAddr) {
	gm.Mutex.RLock()
	defer gm.Mutex.RUnlock()

	var game *Game
	for _, g := range gm.Games {
		if _, ok := g.Players[userID]; ok {
			game = g
			break
		}
	}

	if game == nil || game.Riddle == nil {
		return
	}

	var text string
	var cost int
	if hintType == 1 {
		text = game.Riddle.HintLevel1
		cost = 25
	} else if hintType == 2 {
		text = game.Riddle.HintLevel2
		cost = 50
	} else {
		return
	}

	game.Mutex.Lock()
	game.Scores[userID] -= cost
	game.Mutex.Unlock()

	msg := shared.Message{
		Type: shared.MsgRiddleHint,
		Payload: shared.RiddleHintPayload{
			RiddleID: game.Riddle.ID,
			Text:     text,
			Cost:     cost,
		},
	}

	SendResponse(conn, addr, msg)
}

func (gm *GameManager) ProcessRiddleAnswer(userID int, answer string) {
	gm.Mutex.RLock()
	defer gm.Mutex.RUnlock()

	var game *Game
	for _, g := range gm.Games {
		if _, ok := g.Players[userID]; ok {
			game = g
			break
		}
	}

	if game == nil || game.Riddle == nil {
		return
	}

	if answer == game.Riddle.CorrectWord {
		game.Mutex.Lock()
		game.Scores[userID] += 100
		game.Mutex.Unlock()
		log.Printf("🎉 Joueur %d a deviné correctement ! +100 points", userID)
	}
}

func (gm *GameManager) cleanupGame(code string) {
	time.Sleep(5 * time.Minute)

	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if _, exists := gm.Games[code]; exists {
		delete(gm.Games, code)
		log.Printf("🧹 Partie %s nettoyée", code)
	}
}
func (gm *GameManager) waitForAnswersOrTimeout(game *Game, questionID int, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case qid := <-game.AnswerChan:
			if qid == questionID {
				log.Printf("➡️ Question %d répondue, on passe à la suivante", questionID)
				return
			}
		case <-timer.C:
			log.Printf("⏱️ Temps écoulé pour la question %d", questionID)
			return
		}
	}
}
func (gm *GameManager) WaitAndStartGame(game *Game) {
	game.Mutex.Lock()
	if game.StartTimerLaunched {
		game.Mutex.Unlock()
		return
	}
	game.StartTimerLaunched = true
	game.Mutex.Unlock()

	go func() {
		for {
			game.Mutex.Lock()
			playerCount := len(game.Players)
			game.Mutex.Unlock()

			if playerCount >= 2 {
				log.Printf("🎮 Partie %s - 2 joueurs minimum atteints, lancement dans 30s", game.Code)
				time.Sleep(30 * time.Second)
				gm.StartGame(game.Code)
				gm.RunGame(gm.Conn, game.Code)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
func (gm *GameManager) MonitorLobby(game *Game, conn *net.UDPConn) {
	go func() {
		for {
			game.Mutex.Lock()
			playerCount := len(game.Players)
			game.Mutex.Unlock()

			if playerCount >= 2 {
				log.Printf("🕹️ Partie %s - minimum 2 joueurs atteints, lancement dans 30s", game.Code)
				time.Sleep(30 * time.Second)
				gm.StartGame(game.Code)
				go gm.RunGame(conn, game.Code) // <- utiliser conn ici
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
