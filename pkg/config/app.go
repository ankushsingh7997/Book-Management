package config

import (
	"github.com/ankush/bookstore/env"
	"github.com/ankush/bookstore/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func Connect() {

	logg := logger.NewLogger("UserService", "production")
	dbUrl := env.Get("DB_URL", "")

	d, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		panic(err)
	}
	logg.Info("Connected to Database")
	db = d
}

func GetDB() *gorm.DB {
	return db
}
