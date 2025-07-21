package handlers

import (
	"backend_reservation/pkg/handler"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	handler.Success(w, r, "Users", nil)
}
