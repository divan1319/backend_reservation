package handlers

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/internal/application/services"
	"encoding/json"
	"net/http"
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

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "user": user.Name})
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
