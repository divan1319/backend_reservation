package middleware

import (
	"backend_reservation/pkg/handler"
	"net/http"
)

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserIDFromContext(r.Context())

		if !ok {
			handler.Error(w, r, http.StatusUnauthorized, "Unauthorized, no tienes permisos")
			return
		}

		permission, err := HasPermission(userID, "user")

		if err != nil {
			handler.Error(w, r, http.StatusUnauthorized, "Unauthorized, no se pudo obtener el rol")
			return
		}

		if !permission {
			handler.Error(w, r, http.StatusForbidden, "Forbidden, no tienes permisos de usuario")
			return
		}

		next.ServeHTTP(w, r)
	})
}
