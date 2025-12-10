package handlers

import (
	"fluxx/internal/websocket"
	"log"
	"net/http"

	"github.com/google/uuid" // NOUVEAU: Pour générer des IDs uniques
	gorillaWs "github.com/gorilla/websocket"
)

var upgrader = gorillaWs.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ServeWs gère l'établissement d'une connexion WebSocket.
func ServeWs(hub *websocket.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 1. Générer un ID unique pour cette session
	// En production, cet ID viendrait du jeton d'authentification.
	clientID := uuid.New().String()

	// 2. Créer le Client avec son ID
	client := &websocket.Client{
		Hub:  hub,
		Conn: conn,
		ID:   clientID, // ASSIGNATION DE L'ID
		Send: make(chan []byte, 256),
	}

	log.Printf("Client %s connected.", clientID)

	// 3. Enregistrer le Client
	client.Hub.Register <- client

	// 4. Lancer les pompes
	go client.WritePump()
	go client.ReadPump()
}
