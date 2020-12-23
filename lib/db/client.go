package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)
import "github.com/go-gorp/gorp"

const PostgresDriver = "postgres"
const PostgresDefaultPort = 5432

type DBClient struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string
	Driver   string
	Client   *gorp.DbMap
}

func NewPostgresDBClient(host, dbName, user, pwd string) (*DBClient, error) {
	c := &DBClient{
		Host:     host,
		Port:     PostgresDefaultPort,
		DBName:   dbName,
		Username: user,
		Password: pwd,
		Driver:   PostgresDriver,
		Client:   nil,
	}
	err := c.SslDisabledConnect()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *DBClient) SslDisabledConnect() error {
	conStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", c.Host, c.Port, c.Username, c.DBName, c.Password)
	conn, err := sql.Open(c.Driver, conStr)
	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		return err
	}

	c.Client = &gorp.DbMap{Db: conn, Dialect: gorp.PostgresDialect{}}
	return nil
}

// create new table if not exist, pks is primary key, auto increment disabled
func (c *DBClient) CreateTableIfNotExist(tableName string, v interface{}) error {
	c.Client.AddTableWithName(v, tableName)
	err := c.Client.CreateTablesIfNotExists()
	if err != nil {
		return err
	}
	return nil
}