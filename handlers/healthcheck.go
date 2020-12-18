package handlers

import (
	"github.com/xiaofeiqiu/mlstock/lib/restutils"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	restutils.ResponseWithJson(w,200, restutils.HealthResponse{
		Status: "ok",
	})
}