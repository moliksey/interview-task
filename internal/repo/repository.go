package repo

import (
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	AuthUser(guid string, refreshToken string) (err error)
}
type Registration interface {
	AddUser(email string) (guid string, err error)
}
type Refresh interface {
	GetUser(guid string) (email string, refresh string, err error)
	UpdateRefreshToken(guid string, refresh string) error
}
type Repository struct {
	Authorization
	Registration
	Refresh
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Registration:  NewRegPostgres(db),
		Refresh:       NewRefreshPostgres(db),
	}
}
