package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"os"
)

func conn() *gorm.DB {
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbdatabase := os.Getenv("DB_DATABASE")
	dbcharset := os.Getenv("DB_CHARSET")

	link := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True&loc=Local", dbuser, dbpassword, dbdatabase, dbcharset)

	db, err := gorm.Open("mysql", link)
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}

	return db
}
