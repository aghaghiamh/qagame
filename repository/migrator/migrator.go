package migrator

import (
	"database/sql"
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

// TODO: This will create the gorp_migrations in the users table; rename it to the migrations.

type Migrator struct {
	dbConfig   mysql.Config
	dialect    string
	migrations *migrate.FileMigrationSource
}

func New(dialect string, dbConfig mysql.Config) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}
	return Migrator{dialect: dialect, migrations: migrations, dbConfig: dbConfig}
}

func (m Migrator) Up() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
		m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DBName))

	if err != nil {
		panic(fmt.Errorf("database connection error: %w", err))
	}

	defer db.Close()

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can't apply migrations: %v", err))
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
		m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DBName))

	if err != nil {
		panic(fmt.Errorf("database connection error: %w", err))
	}

	defer db.Close()

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("can't rollback migrations: %v", err))
	}
	fmt.Printf("Applied %d rollbacks!\n", n)
}
