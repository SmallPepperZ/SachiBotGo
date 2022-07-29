package database

import (
	"time"

	"github.com/smallpepperz/sachibotgo/api/config"
	"github.com/smallpepperz/sachibotgo/api/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Get() *gorm.DB {
	if db == nil {
		var err error
		db, err = gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
		if err != nil {
			logger.Err().Fatalln(err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			logger.Err().Fatalln(err)
		}
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(time.Hour)
		db.AutoMigrate(&PotentialInvite{})
	}

	return db
}

func GetPotentialInvite(id string) (potentialInvite *PotentialInvite, err error) {
	tx := db.First(&potentialInvite, id)
	err = tx.Error
	return
}
