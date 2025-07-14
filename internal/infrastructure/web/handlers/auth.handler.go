package handlers

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/internal/application/services"
	"backend_reservation/pkg/firmador"
	"encoding/json"
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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login failed", "error": err.Error()})
		return
	}
	data := map[string]string{
		"user_id": strconv.Itoa(int(user.ID)),
		"email":   user.Email,
		"name":    user.Name,
	}

	token, err := firmador.FirmarToken(data, 24*time.Hour)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "No se pudo firmar el token", "error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "token": token, "user": user.Name})
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Register failed", "error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Register successful", "user": user.Name})
}
