package handler

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func Success(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	status, messageStatus := statusSuccess(r)
	if message == "Login successful" {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(status)
	}

	if message == "" {
		message = messageStatus
	}

	json.NewEncoder(w).Encode(Response{Message: message, Data: data})
}

func Error(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	_, messageStatus := statusError(status)
	w.WriteHeader(status)

	if message == "" {
		message = messageStatus
	}

	json.NewEncoder(w).Encode(Response{Message: message, Data: nil})
}

func statusSuccess(r *http.Request) (int, string) {
	switch r.Method {
	case http.MethodGet:
		return http.StatusOK, "Ok"
	case http.MethodPost:
		return http.StatusCreated, "Item created"
	case http.MethodPut:
		return http.StatusOK, "Item updated"
	case http.MethodDelete:
		return http.StatusOK, "Item deleted"
	case http.MethodPatch:
		return http.StatusOK, "Item patched"
	default:
		return http.StatusMethodNotAllowed, "Method not allowed"
	}
}

func statusError(status int) (int, string) {
	switch status {
	case http.StatusNotFound:
		return http.StatusNotFound, "Not found"
	case http.StatusBadRequest:
		return http.StatusBadRequest, "Bad request"
	case http.StatusUnauthorized:
		return http.StatusUnauthorized, "Unauthorized"
	case http.StatusForbidden:
		return http.StatusForbidden, "Forbidden"
	case http.StatusInternalServerError:
		return http.StatusInternalServerError, "Internal server error"
	default:
		return status, "Error"
	}
}
