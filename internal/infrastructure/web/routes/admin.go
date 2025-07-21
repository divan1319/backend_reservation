package routes

import (
	"backend_reservation/internal/infrastructure/web/handlers"
	"net/http"
)

func AdminRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", handlers.GetUsersHandler)
	return mux
}
