package core

import (
	"fmt"
)

type globalHelper struct {
	Config *Config
	DB     DB
}

var GlobalHelper *globalHelper

func Init(configPath string) {
	config := NewConfig(configPath)
	dbPath := config.GetString("database.path")
	GlobalHelper = &globalHelper{
		Config: config,
		DB:     NewSqliteDB(dbPath),
	}

	fmt.Println("------init cnf success------")
}

func GetDBPath() string {
	return GlobalHelper.Config.GetString("database.path")
}
