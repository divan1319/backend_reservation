package middleware

import (
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/models"
	"backend_reservation/pkg/handler"
	"net/http"
	"strconv"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleID, ok := GetRoleIDFromContext(r.Context())
		if !ok {
			handler.Error(w, r, http.StatusUnauthorized, "Unauthorized, no tienes permisos")
			return
		}
		rol, err := obtenerRole(roleID)
		if err != nil {
			handler.Error(w, r, http.StatusUnauthorized, "Unauthorized, no se pudo obtener el rol")
			return
		}
		if rol != "admin" {
			handler.Error(w, r, http.StatusForbidden, "Forbidden, no tienes permisos de administrador")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func obtenerRole(roleID string) (string, error) {
	_, gormDB, err := connection.GetDB()

	if err != nil {
		return "", err
	}
	role := models.Role{}
	//conevtir roleID a int
	roleIDInt, err := strconv.Atoi(roleID)
	if err != nil {
		return "", err
	}

	resultRole := gormDB.Where("id = ?", roleIDInt).First(&role)
	if resultRole.Error != nil {
		return "", resultRole.Error
	}

	return role.Code, nil
}
