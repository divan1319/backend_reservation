package middleware

import (
	"backend_reservation/pkg/handler"
	"net/http"
)

var allowedOrigins = []string{
	"http://localhost:3000",
}

// Cors es un middleware que gestiona las políticas CORS (Cross-Origin Resource Sharing) para la API.
// Permite únicamente solicitudes desde orígenes definidos en la lista allowedOrigins.
// Si el origen no está permitido, responde con un error 403 Forbidden.
// Si la solicitud es de tipo OPTIONS (preflight), responde con los headers CORS y termina la ejecución.
// Para otros métodos, añade los headers CORS y pasa la solicitud al siguiente handler.
//
// Parámetros:
//   - next: http.Handler que representa el siguiente handler en la cadena de middlewares.
//
// Retorna:
//   - http.HandlerFunc: función que implementa la lógica CORS y delega al siguiente handler si corresponde.
func Cors(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrae el header "Origin" de la solicitud HTTP.
		origin := r.Header.Get("Origin")

		// Verifica si el origen está permitido según la lista de allowedOrigins.
		if !isAllowedOrigin(origin) {
			// Si el origen no está permitido, responde con 403 Forbidden y no continúa.
			handler.Error(w, r, http.StatusForbidden, "Forbidden")
			return
		}

		// Configura los headers CORS para permitir el origen y los métodos/headers necesarios.
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Si la solicitud es de tipo OPTIONS (preflight), responde solo con los headers y termina.
		if r.Method == http.MethodOptions {
			return
		}

		// Si no es una solicitud OPTIONS, continúa con el siguiente handler.
		next.ServeHTTP(w, r)
	}
}

// isAllowedOrigin verifica si el origen proporcionado está en la lista de orígenes permitidos.
// Parámetros:
//   - origin: string que representa el origen de la solicitud HTTP (por ejemplo, "http://localhost:3000").
//
// Retorna:
//   - bool: true si el origen está permitido, false en caso contrario.
func isAllowedOrigin(origin string) bool {
	// Itera sobre la lista de orígenes permitidos y compara con el origen recibido.
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			// Si encuentra una coincidencia, retorna true.
			return true
		}
	}
	// Si no se encuentra el origen en la lista, retorna false.
	return false
}
