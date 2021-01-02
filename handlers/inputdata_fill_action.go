package handlers

import (
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"github.com/xiaofeiqiu/data-preprocessor/services/utils"
	"math"
	"net/http"
	"strconv"
	"time"
)

const length = 30

func (api *ApiHandler) DataInputFillAction(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := NewDataInputRequest(r)
	if err != nil {
		return 400, errors.New("error creating new ema request, " + err.Error())
	}
	log.Info("DataInputFillNEma", "Valid data input ema request")

	// find entries to fill
	entries, err := api.findEntriesToFill(req)
	if err != nil {
		return 500, err
	}

	// get raw data
	rawData, err := api.DBService.FindRawData(entries)
	if err != nil {
		return 500, errors.New("FindRawData failed, " + err.Error())
	}

	SetAction(entries, rawData)

	//update to db
	ct, err := api.DBService.UpdateDataInput(entries)
	if err != nil {
		return 500, err
	}
	log.Info("DataInputFillNEma", strconv.Itoa(ct)+" data input inserted")
	restutils.ResponseWithJson(w, 200, strconv.Itoa(ct)+" updated")

	return 0, nil
}

func SetAction(entries []dbservice.DataInputEntity, rawData []dbservice.RawDataEntity) {

	for i := 0; i+length-1 < len(rawData); i++ {
		rawData[i+length-1].Day30Change = utils.PercentChange(rawData[i].Close, rawData[i+length-1].Close)
	}

	rawDataMap := RawDataArrayToMap(rawData)
	entriesMap := DataInputArrayToMap(entries)

	for _, v := range entries {
		tmp := rawDataMap[v.Date.Format(time.RFC3339)]
		if tmp != nil && tmp.Day30Change != nil {
			entriesMap[v.Date.Format(time.RFC3339)].Action = GetAction(*tmp.Day30Change)
		}
	}
}

var StrongBuy = "Strong Buy"
var StrongSell = "Strong Sell"
var Buy = "Buy"
var Sell = "Sell"
var Hold = "Hold"

func GetAction(change float64) *string {
	if math.Abs(change) <= 3 {
		return &Hold
	} else if change > 3 && change < 10 {
		return &Buy
	} else if change < 3 && change > -10 {
		return &Sell
	} else if change >= 10 {
		return &StrongBuy
	} else if change <= -10 {
		return &StrongSell
	} else {
		return nil
	}
}
