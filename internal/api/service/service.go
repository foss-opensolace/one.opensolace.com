package service

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/api/service/repository"
	"github.com/foss-opensolace/api.opensolace.com/internal/db"
)

var (
	User repository.UserRepository
)

func New() {
	User = repository.NewUserRepository(db.Postgres)
}
