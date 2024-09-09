package services

import (
	"awesomeProject/internal/repo"
)

type Authorization interface {
	AuthUser(guid string, ip string) (access string, refresh string, err error)
}
type Registration interface {
	AddUser(email string) (guid string, err error)
}
type Refresh interface {
	RefreshToken(accessOld string, refreshOld string, ip string) (access string, refresh string, err error)
}
type Service struct {
	Authorization
	Registration
	Refresh
}

func NewService(repos *repo.Repository) *Service {
	return &Service{
		Registration:  NewRegService(repos.Registration),
		Authorization: NewAuthService(repos.Authorization),
		Refresh:       NewRefreshService(repos.Refresh),
	}
}
