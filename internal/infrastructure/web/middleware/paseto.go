package middleware

import (
	"backend_reservation/pkg/firmador"
	"context"
	"net/http"
	"strings"
)

func PasetoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Extraer el token del header Authorization (formato: "Bearer <token>")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			// Si no tiene el prefijo "Bearer ", asumir que es solo el token
			tokenStr = authHeader
		}

		// Verificar el token usando la función VerificarToken
		token, err := firmador.VerificarToken(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Extraer los datos del usuario del token y agregarlos al contexto
		userID, _ := token.GetString("user_id")
		email, _ := token.GetString("email")
		name, _ := token.GetString("name")

		// Crear un contexto con los datos del usuario
		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "name", name)

		// Continuar con el siguiente handler con el contexto actualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Funciones helper para extraer datos del usuario del contexto

// GetUserIDFromContext extrae el ID del usuario del contexto
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}

// GetEmailFromContext extrae el email del usuario del contexto
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("email").(string)
	return email, ok
}

// GetNameFromContext extrae el nombre del usuario del contexto
func GetNameFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value("name").(string)
	return name, ok
}

// GetUserDataFromContext extrae todos los datos del usuario del contexto
func GetUserDataFromContext(ctx context.Context) (userID, email, name string, ok bool) {
	userID, okID := GetUserIDFromContext(ctx)
	email, okEmail := GetEmailFromContext(ctx)
	name, okName := GetNameFromContext(ctx)

	return userID, email, name, okID && okEmail && okName
}
