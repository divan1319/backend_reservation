package middleware

import (
	"backend_reservation/pkg/firmador"
	"backend_reservation/pkg/utils"
	"context"
	"net/http"
	"strings"
)

// Tipos personalizados para las context keys (evita colisiones)
type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	NameKey   contextKey = "name"
)

func PasetoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(w, r, http.StatusUnauthorized, "No token provided")
			return
		}

		// Extraer el token del header Authorization (formato: "Bearer <token>")
		tokenStr := extractToken(authHeader)
		if tokenStr == "" {
			utils.Error(w, r, http.StatusUnauthorized, "Invalid token format")
			return
		}

		// Verificar el token usando la función VerificarToken
		token, err := firmador.VerificarToken(tokenStr)
		if err != nil {
			utils.Error(w, r, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Extraer y validar solo user_id como obligatorio
		userID, err := token.GetString("user_id")
		if err != nil || userID == "" {
			utils.Error(w, r, http.StatusUnauthorized, "Invalid token data")
			return
		}

		// Extraer email y name como opcionales (sin fallar si no están presentes)
		email, _ := token.GetString("email")
		name, _ := token.GetString("name")

		// Crear un contexto con los datos del usuario usando las keys personalizadas
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, EmailKey, email)
		ctx = context.WithValue(ctx, NameKey, name)
		// Continuar con el siguiente handler con el contexto actualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken extrae el token del header Authorization
func extractToken(authHeader string) string {
	// Verificar si tiene el prefijo "Bearer "
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Si no tiene prefijo Bearer, verificar que no esté vacío
	if strings.TrimSpace(authHeader) != "" {
		return authHeader
	}

	return ""
}

// Funciones helper para extraer datos del usuario del contexto

// GetUserIDFromContext extrae el ID del usuario del contexto
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok && userID != ""
}

// GetEmailFromContext extrae el email del usuario del contexto
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok && email != ""
}

// GetNameFromContext extrae el nombre del usuario del contexto
func GetNameFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(NameKey).(string)
	return name, ok && name != ""
}

// GetUserDataFromContext extrae todos los datos del usuario del contexto
func GetUserDataFromContext(ctx context.Context) (userID, email, name, roleID string, ok bool) {
	userID, okID := GetUserIDFromContext(ctx)
	email, _ = GetEmailFromContext(ctx)
	name, _ = GetNameFromContext(ctx)
	// Solo retorna ok=true si user_id está presente (email y name son opcionales)
	return userID, email, name, roleID, okID
}
