package routes

import (
	"backend_reservation/internal/infrastructure/web/handlers"
	"net/http"
)

func AuthRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Usar la sintaxis correcta para Go 1.22+ sin prefijo
	mux.HandleFunc("POST /api/login", handlers.LoginHandler)
	mux.HandleFunc("POST /api/register", handlers.RegisterHandler)

	return mux
}
