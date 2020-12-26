package dbservice

import (
	"fmt"
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
var WhereSymbolAndNilEma = "where symbol=$1 and ema%s is null"

func (s *DBService) FindNullEma(data *[]RawDataEntity, symbol string, timePeriod string) error {
	where := fmt.Sprintf(WhereSymbolAndNilEma, timePeriod)
	query := fmt.Sprintf("%s %s", SelectFromDailyRawData, where)
	_, err := s.client.DB.Select(data, query, symbol)
	if err != nil {
		return err
	}
	return nil
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
