package routes

import (
	"backend_reservation/internal/infrastructure/web/handlers"
	"net/http"
)

func UserRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.GetUserDataHandler)
	return mux
}
