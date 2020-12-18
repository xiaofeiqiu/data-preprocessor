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
