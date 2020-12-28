package dbservice

import (
	"fmt"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
)

// insert ==================
func (s *DBService) BulkInsertRawDataEntity(data []*RawDataEntity) error {
	trans, err := s.client.DB.Begin()
	if err != nil {
		return err
	}

	for _, v := range data {
		err = trans.Insert(v)
		if err != nil {
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *DBService) InsertRawDataEntityIgnoreError(data []*RawDataEntity) int {
	count := 0
	for _, v := range data {
		err := s.client.DB.Insert(v)
		if err == nil {
			count++
		}
	}
	return count
}

// Select======================

// daily raw data query
var SelectFromDailyRawData = "select * from " + dailyRawData
var WhereSymbolAndNilEma = "where symbol=$1 and %s is null"

func (s *DBService) FindNullColEntries(data *[]RawDataEntity, symbol string, colName string) error {
	where := fmt.Sprintf(WhereSymbolAndNilEma, colName)
	query := fmt.Sprintf("%s %s", SelectFromDailyRawData, where)
	_, err := s.client.DB.Select(data, query, symbol)
	if err != nil {
		return err
	}
	return nil
}

var WhereSymbolAndDT = "where symbol=$1 and dt between $2 and $3 order by dt desc"

func (s *DBService) FindRawData(data []DataInputEntity) ([]RawDataEntity, error) {
	result := []RawDataEntity{}

	from, to := GetMinAndMaxDate(data)
	log.Info("FindRawData", "Querying data from "+from+" to "+to)

	where := WhereSymbolAndDT
	query := fmt.Sprintf("%s %s", SelectFromDailyRawData, where)
	_, err := s.client.DB.Select(&result, query, data[0].Symbol, from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// update ==================
func (s *DBService) UpdateEntries(data []RawDataEntity) (int, error) {

	count := 0
	for _, v := range data {
		_, err := s.client.DB.Update(&v)
		if err == nil {
			count++
		}
	}

	return count, nil
}

func toInterfaceArray(data []RawDataEntity) []interface{} {
	result := make([]interface{}, len(data))
	for i, s := range data {
		result[i] = s
	}
	return result
}

// delete
var deleteBySymbol = "delete from %s where symbol=$1"

func (s *DBService) DeleteBySymbol(table, symbol string) error {
	_, err := s.client.DB.Exec(fmt.Sprintf(deleteBySymbol, table), symbol)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBService) ClearDailyRawDataTable() error {
	var deleteAll = "delete from %s where symbol!=$1"
	_, err := s.client.DB.Exec(fmt.Sprintf(deleteAll, dailyRawData), "-")
	if err != nil {
		return err
	}
	return nil
}
