package core

import (
	"encoding/gob"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&User{})

	logrus.WithFields(logrus.Fields{
		"module": "core",
		"action": "init",
	}).Info("Core module initialized")
}

func logFieldsForMethod(method string, args ...string) (fields logrus.Fields) {
	fields = logrus.Fields{
		"module": "core",
		"method": method,
	}
	return fields
}

// DatabaseMigration migrates User and Setting tables
func databaseMigration(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Setting{})
}
