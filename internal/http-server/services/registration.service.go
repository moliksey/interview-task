package services

import (
	"awesomeProject/internal/repo"
)

//func Registration(email string, ip string) (access string, refresh string, guid string, err error) {

//}

type RegService struct {
	repo repo.Registration
}

func NewRegService(repo repo.Registration) *RegService {
	return &RegService{repo: repo}
}
func (r *RegService) AddUser(email string) (guid string, err error) {
	return r.repo.AddUser(email)
}
