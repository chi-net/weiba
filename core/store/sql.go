package store

import (
	"fmt"
	"github.com/chi-net/weiba/core/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() {
	db, err := gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&types.DBTableUser{})
	err = db.AutoMigrate(&types.DBTableRanking{})
	err = db.AutoMigrate(&types.DBTableKVNumbers{})
	if err != nil {
		panic("migration failed")
	}
	sql = db
}
