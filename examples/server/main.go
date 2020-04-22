package main

import (
	"github.com/ledinhbao/core"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {

	server, err := core.ServerFromConfigFile("config.json")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"module": "example",
			"action": "main",
		}).Panicf("Failed to create server, %s", err.Error())
	}

	server.GET("/")

	server.Run(":9098")
}
