package core

import (
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB interface {
	GetDB() *gorm.DB
}

type SqliteDB struct {
	db *gorm.DB
}

func NewSqliteDB(path string) *SqliteDB {
	dbFile := path + string(filepath.Separator) + "server.db"
	DB, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &SqliteDB{
		db: DB,
	}
}

func (db *SqliteDB) GetDB() *gorm.DB {
	return db.db
}
