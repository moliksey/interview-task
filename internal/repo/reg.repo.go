package repo

import (
	"awesomeProject/internal/storage"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type RegPostgres struct {
	db *sqlx.DB
}

func NewRegPostgres(db *sqlx.DB) *RegPostgres {
	return &RegPostgres{db: db}
}
func (r *RegPostgres) AddUser(email string) (guid string, err error) {
	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES ($1) RETURNING %s", storage.UserTable, storage.EmailColumn, storage.GuidColumn)
	err = r.db.QueryRow(query, email).Scan(&guid)
	if err != nil {
		return "0", err
	}
	return guid, nil
}
