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

func (api *ApiHandler) CreateDataInputEntries(w http.ResponseWriter, r *http.Request) (int, error) {

	// find missing entries
	missingEntries := []dbservice.RawDataEntity{}
	err := api.DBService.GetMissingDataInput(&missingEntries)
	if err != nil {
		return 500, errors.New("error querying missing data input entries, " + err.Error())
	}
	log.Info("CreateDataInputEntries", fmt.Sprintf("GetMissingDataInput successful, find %d entries", len(missingEntries)))

	// set input entries
	inputEntries := GetDataInputEntries(missingEntries)
	log.Info("CreateDataInputEntries", strconv.Itoa(len(inputEntries))+" entries set")

	// insert input entries
	ct := api.DBService.InsertDataInputPtrArray(inputEntries)
	log.Info("CreateDataInputEntries", strconv.Itoa(ct)+" inserted")

	restutils.ResponseWithJson(w, 200, strconv.Itoa(ct)+" inserted")
	return 0, nil
}

func GetDataInputEntries(source []dbservice.RawDataEntity) []*dbservice.DataInputEntity {
	result := []*dbservice.DataInputEntity{}
	for _, v := range source {
		tmp := dbservice.DataInputEntity{
			Symbol: v.Symbol,
			Date:   v.Date,
		}
		result = append(result, &tmp)
	}
	return result
}
