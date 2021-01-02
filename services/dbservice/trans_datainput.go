package dbservice

import (
	"fmt"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
)

// Select======================
var getMissingDataInput = "SELECT * FROM daily_raw_data t1 WHERE NOT EXISTS (SELECT * FROM data_input t2 WHERE t1.symbol = t2.symbol and t1.dt = t2.dt)"

func (s *DBService) GetMissingDataInput(data *[]RawDataEntity) error {
	_, err := s.client.DB.Select(data, getMissingDataInput)
	if err != nil {
		return err
	}
	return nil
}

var SelectFromDataInput = "select * from " + dataInput

func (s *DBService) FindNullDataInput(data *[]DataInputEntity, symbol string, colName string) error {
	where := fmt.Sprintf(WhereSymbolAndNilEma, colName)
	query := fmt.Sprintf("%s %s", SelectFromDataInput, where)
	_, err := s.client.DB.Select(data, query, symbol)
	if err != nil {
		return err
	}
	return nil
}

// insert ==
func (s *DBService) InsertDataInputPtrArray(inputData []*DataInputEntity) int {
	count := 0
	for _, v := range inputData {
		err := s.client.DB.Insert(v)
		if err == nil {
			count++
		} else {
			log.Error("InsertDataInputPtrArray", err, "")
		}
	}
	return count
}

// update
func (s *DBService) UpdateDataInput(data []DataInputEntity) (int, error) {

	count := 0
	for _, v := range data {
		_, err := s.client.DB.Update(&v)
		if err == nil {
			count++
		}
	}

	return count, nil
}
