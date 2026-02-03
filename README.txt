Projet : Quiz Battle – Application de quiz multijoueur en Go (UDP)

1. Présentation du projet

Ce projet est une application de quiz interactif développée en langage Go, composée d’un serveur UDP et d’un client graphique réalisé avec la bibliothèque Fyne.

L’application permet à plusieurs joueurs de :

* se connecter via un compte,
* créer ou rejoindre une salle de jeu grâce à un code à 4 chiffres,
* jouer ensemble de manière synchronisée,
* participer à différentes manches de jeu (QCM et devinettes),
* consulter un classement final basé sur les scores.

Le projet met l’accent sur :

* la communication réseau UDP,
* la gestion multijoueur,
* la synchronisation des joueurs,
* une sécurité minimale côté logique (identité, salles, règles de points).


2. Objectifs pédagogiques

Ce projet vise à :

* comprendre le fonctionnement d’une application client/serveur,
* manipuler les sockets UDP en Go,
* gérer un jeu multijoueur synchronisé,
* structurer un projet Go de taille moyenne,
* utiliser une interface graphique multiplateforme (Fyne),
* gérer la persistance des données avec SQLite.


3. Technologies et outils utilisés

Langage et frameworks
* Go (Golang) : logique serveur et client
* Fyne : interface graphique du client
* SQLite : base de données locale côté serveur

Outils de développement
* Visual Studio Code
* Go compiler
* Git & GitHub (gestion de versions)

Réseau
* UDP (User Datagram Protocol) pour la communication rapide client/serveur
* Messages échangés au format JSON

4. Structure du projet

quiz-app-fyne/
│
├── client/
│   ├── main.go              # Point d’entrée du client
│   ├── network.go           # Communication UDP avec le serveur
│   ├── ui_login.go          # Interface de connexion
│   ├── ui_mode.go           # Choix du mode (Solo / Multijoueur)
│   ├── ui_lobby.go          # Lobby (création / rejoindre une salle)
│   ├── ui_waiting.go        # Salle d’attente
│   ├── ui_game.go           # Interface des questions QCM
│   ├── ui_results.go        # Classement final
│
├── server/
│   ├── main.go              # Lancement du serveur UDP
│   ├── udp_handler.go       # Réception et traitement des messages UDP
│   ├── game_manager.go      # Gestion des parties, manches et scores
│   ├── database.go          # Connexion et requêtes SQLite
│
├── shared/
│   ├── message.go           # Types de messages échangés
│   ├── models.go            # Structures (User, Game, Question, etc.)
│
└── README.txt

5. Architecture générale

* Le serveur UDP :

  * gère les utilisateurs connectés,
  * crée et stocke les parties (games),
  * synchronise les joueurs d’une même salle,
  * envoie les questions et résultats,
  * calcule les scores.

* Le client Fyne :

  * affiche les interfaces,
  * envoie les actions du joueur (réponses, indices),
  * reçoit les messages du serveur,
  * met à jour l’interface en temps réel.


6. Fonctionnement de l’application

6.1 Connexion
1. Le joueur lance le client.
2. Il entre son email et son mot de passe.
3. Le serveur valide l’utilisateur et autorise l’accès.

6.2 Choix du mode
* Mode Solo : partie individuelle.
* Mode Multijoueur : accès au lobby.

6.3 Lobby multijoueur (sécurité logique)
Un joueur peut :
  * créer une partie → le serveur génère un code unique à 4 chiffres,
  * rejoindre une partie avec ce code.
Le serveur vérifie :
  * l’existence du code,
  * le nombre de joueurs.
La partie démarre 30s après que le minimum de joueurs est atteint (2 joueurs).

6.4 Déroulement du jeu

Manche 1 : QCM
* Les joueurs reçoivent les mêmes questions au même moment.
* Chaque réponse est envoyée au serveur.
* Le score est mis à jour côté serveur.

Manche 3 : Devinettes

* Le joueur saisit une réponse texte.
* Validation immédiate côté serveur.
* Système d’indices :
  * Indice 1 → −25 points
  * Indice 2 → −50 points
* Quand tous les joueurs ont validé, le serveur affiche les résultats.

6.5 Classement final
* Le serveur calcule les scores finaux.
* Le classement est envoyé à tous les joueurs.
* L’interface affiche :
  * le nom du joueur,
  * son score total.

7. Sécurité (niveau académique)

* Identification des joueurs par ID utilisateur
* Accès aux parties uniquement via code de salle
* Scores calculés uniquement côté serveur
* Le client ne peut pas modifier directement les scores
* Synchronisation centralisée par le serveur


8. Installation et exécution
Prérequis : 
* Go version supérieure à la 1.20
* Git
* SQLite
* Visual Studio Code 

Lancer le serveur: go run main.go

Lancer le client : cd client et go run main.go


