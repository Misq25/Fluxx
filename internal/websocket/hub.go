package websocket

import (
	"encoding/json"
	"log"

	// 🚨 NOUVEL IMPORT : Nous avons besoin du package store
	"fluxx/internal/store"
)

// Note: Ce package s'appelle 'websocket' et inclut Client et Hub.
// La structure Message est définie ailleurs dans ton package websocket, mais elle doit correspondre à store.Message.

// Hub maintient la liste des connexions actives et gère les canaux de messages.
type Hub struct {
	Clients    map[*Client]bool // La liste des utilisateurs connectés
	Store      *store.Store     // 🚨 NOUVEAU CHAMP : Connexion à la BDD Supabase
	Broadcast  chan Message     // Canal où les messages entrants sont envoyés (pour diffusion)
	Register   chan *Client     // Canal pour l'ajout d'un client
	Unregister chan *Client     // Canal pour la suppression d'un client
}

// 🚨 MODIFICATION DE LA SIGNATURE : NewHub accepte maintenant le Store.
func NewHub(s *store.Store) *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Store:      s, // 🚨 AFFECTER LE STORE
	}
}

// Run est la boucle principale qui écoute les canaux du Hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// ... (Gestion de l'enregistrement inchangée)
			h.Clients[client] = true

			// 💡 OPTIONNEL : Nous pourrions ajouter ici la logique pour charger l'historique
			// des messages depuis le Store et les envoyer à ce nouveau client. (Prochaine étape!)

		case client := <-h.Unregister:
			// ... (Gestion de la désinscription inchangée)
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

		case message := <-h.Broadcast:
			// 🚨 LOGIQUE BDD : Sauvegarder le message AVANT de le diffuser

			// Note: message.Sender correspond à clientID et message.Content au contenu
			if err := h.Store.SaveMessage(message.Sender, message.Content); err != nil {
				log.Printf("Erreur lors de l'enregistrement du message dans la BDD: %v", err)
				// Le chat continue, mais le message est perdu après un redémarrage.
			}

			// --- Diffusion (Broadcasting) inchangée ---

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
