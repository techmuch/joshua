package api

import (
	"bd_bot/internal/repository"
	"encoding/json"
	"net/http"
)

type SolicitationHandler struct {
	repo *repository.SolicitationRepository
}

func (h *SolicitationHandler) List(w http.ResponseWriter, r *http.Request) {
	solicitations, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(solicitations)
}

func (h *SolicitationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	detail, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Solicitation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

type ClaimRequest struct {
	Type string `json:"type"` // 'interested', 'lead', 'none'
}

func (h *SolicitationHandler) Claim(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	// We need the internal integer ID for the Claim table FK.
	// GetByID fetches by source_id.
	// But UpsertClaim needs (int, int).
	// So first we must resolve the solicitation to get its PK.
	
	sol, err := h.repo.GetByID(r.Context(), idStr)
	if err != nil {
		http.Error(w, "Solicitation not found", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(int)
	
	var req ClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

		if err := h.repo.UpsertClaim(r.Context(), userID, sol.ID, req.Type); err != nil {

			http.Error(w, "Failed to update claim", http.StatusInternalServerError)

			return

		}

	

			w.WriteHeader(http.StatusOK)

	

		}

	

		

	

		type ArchiveRequest struct {

	

			Archived bool `json:"archived"`

	

		}

	

		

	

		func (h *SolicitationHandler) Archive(w http.ResponseWriter, r *http.Request) {

	

			idStr := r.PathValue("id")

	

			sol, err := h.repo.GetByID(r.Context(), idStr)

	

			if err != nil {

	

				http.Error(w, "Solicitation not found", http.StatusNotFound)

	

				return

	

			}

	

			

	

			userID := r.Context().Value("user_id").(int)

	

			var req ArchiveRequest

	

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

	

				http.Error(w, "Invalid request", http.StatusBadRequest)

	

				return

	

			}

	

		

	

			if req.Archived {

	

				err = h.repo.Archive(r.Context(), userID, sol.ID)

	

			} else {

	

				err = h.repo.Unarchive(r.Context(), userID, sol.ID)

	

			}

	

		

	

			if err != nil {

	

				http.Error(w, "Failed to update archive status", http.StatusInternalServerError)

	

				return

	

			}

	

			w.WriteHeader(http.StatusOK)

	

		}

	

		

	

		type ShareRequest struct {

	

			Email   string `json:"email"`

	

			Message string `json:"message"`

	

		}

	

		

	

		func (h *SolicitationHandler) Share(w http.ResponseWriter, r *http.Request) {

	

			idStr := r.PathValue("id")

	

			sol, err := h.repo.GetByID(r.Context(), idStr)

	

			if err != nil {

	

				http.Error(w, "Solicitation not found", http.StatusNotFound)

	

				return

	

			}

	

		

	

			userID := r.Context().Value("user_id").(int)

	

			var req ShareRequest

	

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

	

				http.Error(w, "Invalid request", http.StatusBadRequest)

	

				return

	

			}

	

		

	

			if err := h.repo.Share(r.Context(), sol.ID, userID, req.Email, req.Message); err != nil {

	

				http.Error(w, "Failed to share solicitation", http.StatusInternalServerError)

	

				return

	

			}

	

			w.WriteHeader(http.StatusOK)

	

		}

	

		

	

	type AddCommentRequest struct {

		Content string `json:"content"`

	}

	

	func (h *SolicitationHandler) AddComment(w http.ResponseWriter, r *http.Request) {

		idStr := r.PathValue("id")

		

		sol, err := h.repo.GetByID(r.Context(), idStr)

		if err != nil {

			http.Error(w, "Solicitation not found", http.StatusNotFound)

			return

		}

	

		userID := r.Context().Value("user_id").(int)

		

		var req AddCommentRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

			http.Error(w, "Invalid request", http.StatusBadRequest)

			return

		}

	

		if err := h.repo.AddComment(r.Context(), sol.ID, userID, req.Content); err != nil {

			http.Error(w, "Failed to add comment", http.StatusInternalServerError)

			return

		}

	

		w.WriteHeader(http.StatusOK)

	}

	