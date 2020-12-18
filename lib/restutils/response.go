package restutils

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:status`
}

type ErrorResponse struct {
	Event string `json:event`
	Error string `json:error`
}

func ResponseWithJson(w http.ResponseWriter, code int, body interface{}) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling body to json"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBody)
	return
}

func Is2xxStatusCode(status int) bool {
	if status >= 200 && status < 300 {
		return true
	}
	return false
}

func Is4xxStatusCode(status int) bool {
	if status >= 400 && status < 500 {
		return true
	}
	return false
}

func Is5xxStatusCode(status int) bool {
	if status >= 500 && status < 600 {
		return true
	}
	return false
}
