package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var conn *gorm.DB

func DB() *gorm.DB {
	if conn != nil {
		return conn
	}
	var err error
	driver := os.Getenv("DB_DRIVER")
	if driver == "sqlite" {
		conn, err = gorm.Open(sqlite.Open(sqliteDsn()), &gorm.Config{})
	} else {
		conn, err = gorm.Open(mysql.Open(mysqlDsn()), &gorm.Config{})
	}
	if err != nil {
		log.Panic(err)
	}
	return conn
}

func sqliteDsn() string {
	path, _ := filepath.Abs("gorm.db")
	return path
}

func mysqlDsn() string {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DATABASE")

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
}
