package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"net/http"
)

func (api *ApiHandler) ClearRawData(w http.ResponseWriter, r *http.Request) (int, error) {
	err := api.DBService.ClearDailyRawDataTable()
	if err != nil {
		return 500, err
	}
	restutils.ResponseWithJson(w, 200, "cleared")
	return 0, nil
}
