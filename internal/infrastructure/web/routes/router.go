package routes

import (
	"net/http"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Obtener las rutas de autenticación
	authRoutes := AuthRoutes()

	// Montar las rutas de autenticación en el router principal sin prefijo
	mux.Handle("/", authRoutes)

	return mux
}
