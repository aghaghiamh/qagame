package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`

	MaxLifeTime time.Duration `mapstructure:"max_life_time"`
	MaxOpenConn int           `mapstructure:"max_open_conn"`
	MaxIdleConn int           `mapstructure:"max_idle_conn"`
}

type MysqlDB struct {
	db *sql.DB
}

func New(conf Config) (*MysqlDB, error) {
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

func (mysql *MysqlDB) GetDB() *sql.DB {
	return mysql.db
}
