package routes

import (
	"backend_reservation/internal/infrastructure/web/handlers"
	"backend_reservation/internal/infrastructure/web/middleware"
	"net/http"
)

func AuthRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Usar la sintaxis correcta para Go 1.22+ sin prefijo
	mux.HandleFunc("POST /api/login", handlers.LoginHandler)
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler)
	mux.Handle("GET /api/profile", middleware.PasetoMiddleware(http.HandlerFunc(handlers.ProfileHandler)))
	return mux
}
