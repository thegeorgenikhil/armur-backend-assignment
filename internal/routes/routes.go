package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/handlers"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/middleware"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/api/register", handlers.Register)                            
	r.Post("/api/login", handlers.Login)                                 
	r.With(middleware.AuthMiddleware).Get("/api/logout", handlers.Logout) 
	r.Post("/api/refresh", handlers.RefreshToken)

	r.Post("/api/forgot", handlers.ForgotPassword)        
	r.Post("/api/reset-password", handlers.ResetPassword) 

	r.Get("/api/user-activation/verify", handlers.VerifyUser)
	r.With(middleware.AuthMiddleware).Get("/api/user-activation/check", handlers.CheckUserVerificationStatus)  
	r.With(middleware.AuthMiddleware).Get("/api/user-activation/resend", handlers.ResendUserVerificationEmail) 

	r.With(middleware.AuthMiddleware).Get("/api/user", handlers.GetUser) 

	// File Server
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/*", http.StripPrefix("/", fs))

	return r
}
