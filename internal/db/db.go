package db

import (
	"context"

	"github.com/foss-opensolace/one.opensolace.com/internal/api/model"
	"github.com/foss-opensolace/one.opensolace.com/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Postgres *gorm.DB
var Mongo *mongo.Client

func New() {
	if err := initPostgres(); err != nil {
		panic(err)
	}
}

func initMongo() error {
	client, err := mongo.Connect(options.Client().ApplyURI(config.DB.MongoConnectionString))
	if err != nil {
		return err
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	Mongo = client

	return nil
}

func initPostgres() error {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: config.DB.PSQLConnectionString}))
	if err != nil {
		return err
	}

	db.Logger = logger.Discard

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}

	if err := migrate(db); err != nil {
		return err
	}

	var apiKeys int64
	if err := db.Model(&model.APIKey{}).Count(&apiKeys).Error; err != nil {
		return err
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

	return nil
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
