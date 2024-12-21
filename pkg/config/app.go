package config

import (
	"github.com/ankush/bookstore/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func Connect() {
	logg := logger.NewLogger("UserService", "production")
	d, err := gorm.Open("mysql", "root:12345678@tcp(0.0.0.0:3306)/rest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	logg.Info("Connected to Database")
	db = d
}

func GetDB() *gorm.DB {
	return db
}
