package dbservice

import (
	"github.com/xiaofeiqiu/data-preprocessor/lib/db"
)

// table names
const dailyRawData = "daily_raw_data"

type DBService struct {
	client *db.DBClient
}

func NewDBService(client *db.DBClient) *DBService {
	return &DBService{
		client: client,
	}
}

func (s *DBService) InitDBTableMapping() error {
	s.client.DB.AddTableWithName(RawDataEntity{}, dailyRawData)
	err := s.client.DB.CreateTablesIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

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
