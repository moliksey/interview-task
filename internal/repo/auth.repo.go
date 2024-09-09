package repo

import (
	"awesomeProject/internal/storage"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
func (r *AuthPostgres) AuthUser(guid string, refreshToken string) (err error) {
	query := fmt.Sprintf("SELECT %s FROM %s where %s=$1", storage.EmailColumn, storage.UserTable, storage.GuidColumn)
	row := r.db.QueryRow(query, guid)
	var email string
	err = row.Scan(&email)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("UPDATE %s SET %s=$1 where %s=$2", storage.UserTable, storage.RefreshColumn, storage.GuidColumn)
	row = r.db.QueryRow(query, refreshToken, guid)
	return row.Err()
}
