package api

import (
	"fluxx/internal/api/handlers"
	"fluxx/internal/websocket"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NewRouter initialise et configure toutes les routes de l'API.
// Il reçoit le Hub, car il doit le transmettre aux handlers qui en ont besoin (comme ServeWs).
func NewRouter(hub *websocket.Hub) http.Handler {
	r := chi.NewRouter()

	// Route de test simple
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Fluxx API is running!"))
	})

	// Route WebSocket pour la connexion temps réel
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Appel au handler qui va transformer la connexion HTTP en WS
		handlers.ServeWs(hub, w, r)
	})

	return r
}
