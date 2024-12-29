package db

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/api/model"
	"github.com/foss-opensolace/api.opensolace.com/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Postgres *gorm.DB

func New() {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: config.DB.ConnectionString}))
	if err != nil {
		panic(err)
	}

	db.Logger = logger.Discard

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		panic(err)
	}

	if err := migrate(db); err != nil {
		panic(err)
	}

	Postgres = db
}

func migrate(db *gorm.DB) error {
	models := []any{
		model.APIKey{},
		model.APIKeyPermissions{},

		model.User{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}
