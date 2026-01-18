package api

import (
	"bd_bot/internal/repository"
	"net/http"
)

func NewRouter(solRepo *repository.SolicitationRepository, userRepo *repository.UserRepository, matchRepo *repository.MatchRepository) *http.ServeMux {
	mux := http.NewServeMux()

	solHandler := &SolicitationHandler{repo: solRepo}
	authHandler := &AuthHandler{repo: userRepo}
	userHandler := &UserHandler{repo: userRepo}
	matchHandler := &MatchHandler{repo: matchRepo}

	mux.HandleFunc("/api/solicitations", solHandler.List)
	mux.HandleFunc("/api/matches", matchHandler.List)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	mux.HandleFunc("/api/auth/me", authHandler.Me)
	mux.HandleFunc("/api/auth/password", AuthMiddleware(authHandler.ChangePassword))
	mux.HandleFunc("/api/user/narrative", AuthMiddleware(userHandler.UpdateNarrative))
	mux.HandleFunc("/api/user/profile", AuthMiddleware(userHandler.UpdateProfile))
	mux.HandleFunc("/api/user/avatar", AuthMiddleware(userHandler.UploadAvatar))

	// Serve uploaded files
	fs := http.FileServer(http.Dir("uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	return mux
}
