package handlers

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/internal/application/services"
	"backend_reservation/internal/infrastructure/web/middleware"
	"backend_reservation/pkg/firmador"
	"backend_reservation/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	loginDto := dto.LoginDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := services.Login(&loginDto)

	if err != nil {
		utils.Error(w, r, http.StatusNotFound, "Login failed")
		return
	}
	data := map[string]string{
		"user_id": strconv.Itoa(int(user.ID)),
		"email":   user.Email,
		"name":    user.Name,
	}

	token, err := firmador.FirmarToken(data, 35*time.Second) //token valido por 30 segundos

	if err != nil {
		utils.Error(w, r, http.StatusInternalServerError, "No se pudo firmar el token")
		return
	}
	dataUser := map[string]string{"token": token, "user": user.Name}

	utils.Success(w, r, "Login successful", dataUser)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	registerDto := dto.RegisterDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
	}

	user, err := services.Register(&registerDto)

	if err != nil {
		utils.Error(w, r, http.StatusBadRequest, "Register failed")
		return
	}

	dataUser := map[string]interface{}{"user": user}
	utils.Success(w, r, "Register successful", dataUser)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {

	userID, email, name, ok := middleware.GetUserDataFromContext(r.Context())

	if !ok {
		utils.Error(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Crear respuesta defensiva (manejar campos opcionales)
	userProfile := map[string]string{
		"user_id": userID,
	}

	// Solo agregar email y name si están presentes
	if email != "" {
		userProfile["email"] = email
	}
	if name != "" {
		userProfile["name"] = name
	}

	utils.Success(w, r, "Profile", userProfile)
}
