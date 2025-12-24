package api

import (
	"encoding/json"
	"fluxx/internal/store"
	"net/http"
)

// AuthHandler va contenir notre Store pour pouvoir l'utiliser
type AuthHandler struct {
	repo *store.Store
}

func NewAuthHandler(s *store.Store) *AuthHandler {
	return &AuthHandler{repo: s}
}

// HandleSyncProfile reçoit les infos de l'utilisateur et les enregistre
func (h *AuthHandler) HandleSyncProfile(w http.ResponseWriter, r *http.Request) {
	var p store.Profile

	// 1. On décode le JSON envoyé par le front (ID, Username, etc.)
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// 2. On appelle ta super fonction SyncUserProfile du Store
	if err := h.repo.SyncUserProfile(p); err != nil {
		http.Error(w, "Erreur lors de la sauvegarde du profil", http.StatusInternalServerError)
		return
	}

	// 3. On répond que tout est OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
