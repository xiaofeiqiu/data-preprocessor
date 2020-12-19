package restutils

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	ResponseWithJson(w,200, HealthResponse{
		Status: "ok",
	})
}
