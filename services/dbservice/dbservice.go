package dbservice

import "github.com/xiaofeiqiu/data-preprocessor/lib/db"

type DBService struct {
	Client *db.DBClient
}

func (c *DBService) BulkInsertRawDataEntity () {

}
