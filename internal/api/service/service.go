package service

import (
	"github.com/foss-opensolace/one.opensolace.com/internal/api/service/repository"
	"github.com/foss-opensolace/one.opensolace.com/internal/db"
)

var (
	APIKey repository.APIKeyRepository
	User   repository.UserRepository
)

func New() {
	db := db.Postgres

	APIKey = repository.NewAPIKeyRepository(db)
	User = repository.NewUserRepository(db)
}
