package api

import (
	"bd_bot/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type AuthHandler struct {
	repo *repository.UserRepository
}

type LoginRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.FindByEmail(r.Context(), req.Email)
	if err != nil {
		// User not found
		if req.Password != "" {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Auto-provisioning (Mock SSO flow)
		if req.FullName == "" {
			req.FullName = "New User"
		}
		user, err = h.repo.Create(r.Context(), req.Email, req.FullName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// User found
		if req.Password != "" {
			if !h.repo.VerifyPassword(r.Context(), user, req.Password) {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}
		}
	}

	// Update Activity
	h.repo.UpdateLastActive(r.Context(), user.ID)

	// Set Session Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(user.ID),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Middleware should have set this
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(int)

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.NewPassword == "" {
		http.Error(w, "New password is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdatePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err.Error() == "incorrect current password" {
			http.Error(w, "Incorrect current password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Password updated successfully"}`))
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user, err := h.repo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
}

// Middleware to check session
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}