package dao

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func InitDB(dbType string, connectionString string, logMode bool) (*gorm.DB, error) {

	var err error

	db, err = gorm.OpenPool(200, 100, dbType, connectionString)

	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}

	db.LogMode(logMode)

	return db, err
}

func GetDB() *gorm.DB {

	return db
}
