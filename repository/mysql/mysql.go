package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string

	MaxLifeTime time.Duration
	MaxOpenConn int
	MaxIdleConn int
}

type MysqlDB struct {
	db *sql.DB
}

func New(conf MysqlConfig) (*MysqlDB, error) {
	connParameter := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", conf.Username, conf.Password, "", conf.Host, conf.Port, conf.DBName) // the empty string here define the protocol
	db, err := sql.Open("mysql", connParameter)
	if err != nil {
		panic(fmt.Errorf("database connection error: %w", err))
	}

	db.SetConnMaxLifetime(conf.MaxLifeTime)
	db.SetMaxOpenConns(conf.MaxOpenConn)
	db.SetMaxIdleConns(conf.MaxIdleConn)

	return &MysqlDB{
		db: db,
	}, nil
}
