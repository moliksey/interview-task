package storage

import (
	"awesomeProject/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	UserTable     string = "users"
	GuidColumn    string = "id"
	RefreshColumn string = "refresh_token"
	EmailColumn   string = "email"
)

func NewPostgresDb(config config.Storage) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Dbname, config.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil

}
