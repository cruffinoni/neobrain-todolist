package database

import (
	"fmt"
	"github.com/cruffinoni/neobrain-todolist/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func connectToDatabase(config *config.Database) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Username, config.Password, config.Host, config.Port, config.Database),
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func NewDB(config *config.Database) (*DB, error) {
	dbConnection, err := connectToDatabase(config)
	if err != nil {
		return nil, err
	}
	return &DB{db: dbConnection}, nil
}
