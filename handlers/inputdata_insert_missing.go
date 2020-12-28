package handlers

import (
	"errors"
	"fmt"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"net/http"
	"strconv"
)

func (api *ApiHandler) InsertMissingDataInput(w http.ResponseWriter, r *http.Request) (int, error) {

	// find missing entries
	missingEntries := []dbservice.RawDataEntity{}
	err := api.DBService.GetMissingDataInput(&missingEntries)
	if err != nil {
		return 500, errors.New("error querying missing data input entries, " + err.Error())
	}
	log.Info("InsertMissingDataInput", fmt.Sprintf("GetMissingDataInput successful, find %d entries", len(missingEntries)))

	inputEntries := GetDataInputEntries(missingEntries)
	log.Info("InsertMissingDataInput", strconv.Itoa(len(inputEntries)) + " entries set")
	restutils.ResponseWithJson(w, 200, inputEntries)
	return 0, nil
}

func GetDataInputEntries(source []dbservice.RawDataEntity) []*dbservice.DataInput {
	result := []*dbservice.DataInput{}
	for _, v := range source {
		tmp := dbservice.DataInput{
			Symbol: v.Symbol,
			Date:   v.Date,
		}
		result = append(result, &tmp)
	}
	return result
}
