package dbservice

import (
	"github.com/xiaofeiqiu/data-preprocessor/lib/db"
)

// table names
const dailyRawData = "daily_raw_data"
const dataInput = "data_input"

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
	s.client.DB.AddTableWithName(DataInput{}, dataInput)
	err := s.client.DB.CreateTablesIfNotExists()
	if err != nil {
		return err
	}
	return nil
}
