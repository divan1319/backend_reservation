package routes

import (
	"backend_reservation/internal/infrastructure/web/middleware"
	"net/http"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Obtener las rutas de autenticación
	authRoutes := AuthRoutes()
	adminRoutes := AdminRoutes()
	// Montar las rutas de autenticación en el router principal sin prefijo
	mux.Handle("/", authRoutes)

	// Usar StripPrefix para remover "/api/admin" antes de pasar al handler de admin
	mux.Handle("/api/admin/", http.StripPrefix("/api/admin", middleware.PasetoMiddleware(middleware.AdminMiddleware(adminRoutes))))

	return mux
}
