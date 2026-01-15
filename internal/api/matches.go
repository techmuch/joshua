package api

import (
	"bd_bot/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type MatchHandler struct {
	repo *repository.MatchRepository
}

func (h *MatchHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	matches, err := h.repo.GetUserInbox(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch matches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}
