package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"net/http"
	"reflect"
	"runtime"
)

func ErrorHandler(f func(http.ResponseWriter, *http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := f(w, r)
		if err != nil {
			res := restutils.ErrorResponse{
				Event: runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
				Error: err.Error(),
			}
			log.Error(res.Event, res.Error)
			restutils.ResponseWithJson(w, status, res)
		}
	}
}
