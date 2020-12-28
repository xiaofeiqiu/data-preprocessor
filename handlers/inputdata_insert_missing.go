package handlers

import (
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"net/http"
)

func (api *ApiHandler) InsertMissingDataInput(w http.ResponseWriter, r *http.Request) (int, error) {
	missingEntries := []dbservice.RawDataEntity{}
	err := api.DBService.GetMissingDataInput(&missingEntries)
	if err != nil {
		return 500, errors.New("error querying missing data input entries, " + err.Error())
	}

	restutils.ResponseWithJson(w, 200, missingEntries)
	return 0, nil
}
