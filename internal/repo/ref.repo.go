package repo

import (
	"awesomeProject/internal/storage"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type RefreshPostgres struct {
	db *sqlx.DB
}

func NewRefreshPostgres(db *sqlx.DB) *RefreshPostgres {
	return &RefreshPostgres{db: db}
}
func (r *RefreshPostgres) GetUser(guid string) (email string, refresh string, err error) {
	query := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s=($1)", storage.EmailColumn, storage.RefreshColumn, storage.UserTable, storage.GuidColumn)
	row := r.db.QueryRow(query, guid)
	err = row.Scan(&email, &refresh)
	if err != nil {
		return "", "", err
	}
	return email, refresh, nil
}
func (r *RefreshPostgres) UpdateRefreshToken(guid string, refresh string) error {
	query := fmt.Sprintf("UPDATE %s SET %s=$1 where %s=$2", storage.UserTable, storage.RefreshColumn, storage.GuidColumn)
	row := r.db.QueryRow(query, refresh, guid)
	return row.Err()
}
