package user

import (
	"database/sql"

	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
)

type Storage struct {
	db *sql.DB
}

// To have a consistent mysql database configs, we just use the general definition of MysqlDB struct.
func New(mysqlDB *mysql.MysqlDB) *Storage {
	return &Storage{
		db: mysqlDB.GetDB(),
	}
}
