package database

import (
	"github.com/badaccuracyid/softeng_backend/src/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"sync"
)

var (
	database *gorm.DB
	once     sync.Once
)

func GetPostgresDatabase() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		database, err = connect()
	})

	if err != nil {
		return nil, err
	}

	return database, nil
}

func MigrateTables() error {
	db, err := GetPostgresDatabase()
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.Message{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.Conversation{})
	if err != nil {
		return err
	}

	return nil
}

func connect() (*gorm.DB, error) {
	envDsn := os.Getenv("POSTGRES_DSN")
	if envDsn == "" {
		panic("POSTGRES_DSN is not set in .env file")
	}

	db, err := gorm.Open(postgres.Open(envDsn))
	if err != nil {
		return nil, err
	}

	return db, nil
}
