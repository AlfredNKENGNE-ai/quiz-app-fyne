package main

import (
	"fmt"
	"log"
	"net"

	"quiz-app-fyne/server"
)

const ServerPort = 9000

func main() {
	// Initialisation des bases de données
	server.InitDatabases()

	// Création du socket UDP
	addr := net.UDPAddr{
		Port: ServerPort,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("🚀 Serveur UDP lancé sur le port", ServerPort)
	log.Println("🚀 Serveur UDP prêt et à l'écoute")

	buffer := make([]byte, 4096)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
			continue
		}
		go server.HandleMessage(conn, clientAddr, buffer[:n])
	}
}
