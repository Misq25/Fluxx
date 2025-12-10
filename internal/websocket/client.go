package websocket

import (
	"encoding/json" // NOUVEAU: Pour la conversion JSON
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait       = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client est un wrapper autour de la connexion WebSocket pour un utilisateur.
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	ID   string // NOUVEAU: Identifiant de l'utilisateur (ex: UUID)
	Send chan []byte
}

// ReadPump lit les messages du client, les décode en JSON, et les envoie au Hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Lecture des bytes bruts envoyés par le client (supposés être du JSON)
		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("error client %s: %v", c.ID, err)
			}
			break
		}

		// Le client web envoie juste le contenu dans un objet { "content": "..." }
		var msgContent struct {
			Content string `json:"content"`
		}

		// 1. Décodage du message entrant (JSON string)
		if err := json.Unmarshal(rawMessage, &msgContent); err != nil {
			log.Printf("Erreur JSON Unmarshal from client %s: %v", c.ID, err)
			continue
		}

		// 2. Création de la structure Message complète
		message := Message{
			Sender:  c.ID, // Utilise l'ID de ce client comme expéditeur
			Content: msgContent.Content,
		}

		// 3. Envoi au Hub pour diffusion
		c.Hub.Broadcast <- message
	}
}

// WritePump envoie les messages du canal Send (déjà marshallés en JSON par le Hub) au client.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(pongWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message) // Le message est déjà au format []byte JSON
			w.Close()

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(pongWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
