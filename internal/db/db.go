package db

import (
	"github.com/foss-opensolace/one.opensolace.com/internal/api/model"
	"github.com/foss-opensolace/one.opensolace.com/internal/config"
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

	var apiKeys int64
	if err := db.Model(&model.APIKey{}).Count(&apiKeys).Error; err != nil {
		panic(err)
	}

	if apiKeys == 0 {
		apiKey := &model.APIKey{
			Description: "Main service API key",
		}

		db.Create(apiKey)
		db.Create(&model.APIKeyPermissions{
			KeyID:            apiKey.ID,
			KeyAssign:        true,
			KeyCreate:        true,
			KeyRead:          true,
			KeyUpdate:        true,
			KeyRevoke:        true,
			KeyDelete:        true,
			Health:           true,
			Metrics:          true,
			UserAuthLogin:    true,
			UserAuthRegister: true,
			UserUpdate:       true,
			UserRead:         true,
			UserDelete:       true,
		})
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
