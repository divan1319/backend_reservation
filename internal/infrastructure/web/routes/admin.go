package routes

import (
	"backend_reservation/internal/infrastructure/web/handlers"
	"net/http"
)

func AdminRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", handlers.GetUsersHandler)
	mux.HandleFunc("GET /service", handlers.ObtenerServiciosHandler)
	mux.HandleFunc("POST /service", handlers.CrearServicioHandler)
	mux.HandleFunc("GET /service/{id}", handlers.ObtenerServicioHandler)

	return mux
}
