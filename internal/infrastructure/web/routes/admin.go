package routes

import (
	"backend_reservation/internal/infrastructure/web/handlers"
	"net/http"
)

func AdminRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", handlers.GetUsersHandler)

	//Rutas para servicios
	mux.HandleFunc("GET /service", handlers.ObtenerServiciosHandler)
	mux.HandleFunc("POST /service", handlers.CrearServicioHandler)
	mux.HandleFunc("GET /service/{id}", handlers.ObtenerServicioHandler)
	mux.HandleFunc("PATCH /service/{id}", handlers.ActualizarServicioHandler)
	mux.HandleFunc("DELETE /service/{id}", handlers.EliminarServicioHandler)
	mux.HandleFunc("PUT /service/{id}/activate", handlers.ActivarDesactivarServicioHandler)

	return mux
}
