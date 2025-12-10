package websocket

import (
	"encoding/json"
	"log"
)

// Note: Ce package s'appelle 'websocket' et inclut Client et Hub.

// Hub maintient la liste des connexions actives et gère les canaux de messages.
type Hub struct {
	Clients    map[*Client]bool // La liste des utilisateurs connectés
	Broadcast  chan Message     // Canal où les messages entrants sont envoyés (pour diffusion)
	Register   chan *Client     // Canal pour l'ajout d'un client
	Unregister chan *Client     // Canal pour la suppression d'un client
}

// NewHub crée et retourne une nouvelle instance de Hub.
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// Run est la boucle principale qui écoute les canaux du Hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// Ajout d'un nouveau client
			h.Clients[client] = true
		case client := <-h.Unregister:
			// Suppression d'un client déconnecté ou en erreur
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			payload, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}
			// Diffusion du message à tous les clients
			for client := range h.Clients {
				select {
				case client.Send <- payload:
					// Envoi réussi
				default:
					// Si l'envoi bloque (le client n'arrive pas à traiter) : on déconnecte
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
