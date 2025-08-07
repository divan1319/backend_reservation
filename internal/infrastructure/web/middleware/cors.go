package middleware

import (
	"backend_reservation/pkg/handler"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
	AllowAllOrigins  bool //Para permitir todos los orígenes , true para desarrollo, false para producción
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept", "Origin",
			"Cache-Control",
			"X-File-Name",
		},
		AllowCredentials: true,
		MaxAge:           84600, // 24 horas en segundos
		AllowAllOrigins:  false,
	}
}

func CorsWithConfig(config CORSConfig) func(next http.Handler) http.HandlerFunc {
	methods := make([]string, len(config.AllowedMethods))

	for i, method := range config.AllowedMethods {
		methods[i] = strings.ToUpper(method)
	}

	methodsStr := strings.Join(methods, ", ")
	headersStr := strings.Join(config.AllowedHeaders, ", ")

	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin == "" {
				handler.Error(w, r, http.StatusForbidden, "No origin provided")
				return
			}

			allowedOrigin := ""

			if config.AllowAllOrigins {
				allowedOrigin = origin
			} else {
				for _, allowed := range config.AllowedOrigins {
					if isOriginAllowed(origin, allowed) {
						allowedOrigin = allowed
						break
					}
				}
			}

			if allowedOrigin == "" && !config.AllowAllOrigins {
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusForbidden)
					return
				}

				handler.Error(w, r, http.StatusForbidden, "Origin not allowed")
				return
			}

			// Configurar headers CORS
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", methodsStr)
			w.Header().Set("Access-Control-Allow-Headers", headersStr)

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if config.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
			}

			// Manejar solicitudes OPTIONS (preflight)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		}
	}

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
	return CorsWithConfig(DefaultCORSConfig())(next)
}

// isOriginAllowed verifica si un origen está permitido
// Soporta coincidencias exactas y patrones básicos con *
func isOriginAllowed(origin, allowed string) bool {
	if allowed == "*" {
		return true
	}

	if allowed == origin {
		return true
	}

	// Soporte para wildcards (ej: *.example.com)
	if strings.Contains(allowed, "*") {
		pattern := strings.ReplaceAll(allowed, "*", ".*")
		match, _ := regexp.MatchString(pattern, origin)
		return match
	}

	return false
}

// Para desarrollo: permite todos los orígenes
func CorsAllowAll(next http.Handler) http.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowAllOrigins = true
	return CorsWithConfig(config)(next)
}
