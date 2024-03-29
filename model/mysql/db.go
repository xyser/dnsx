// Package mysql connect
package mysql

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/xyser/dnsx/pkg/config"
	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var once sync.Once
var db *gorm.DB

// Init db connect
func Init() {
	var err error
	once.Do(func() {
		// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
		dsn := config.GetString("mysql.dns")
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			fmt.Println("\033[1;30;42m[info]\033[0m db [maste] connect success")
		} else {
			fmt.Printf("\033[1;30;41m[error]\033[0m db [master] connect error: %s", err.Error())
			os.Exit(1)
		}

		sqlDB, _ := db.DB()
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.GetInt("mysql.conn_max_lifetime")))
		sqlDB.SetMaxIdleConns(config.GetInt("mysql.max_idle_conn"))
		sqlDB.SetMaxOpenConns(config.GetInt("mysql.max_open_conn"))
	})
}

// GetDB get db connect
func GetDB() *gorm.DB {
	return db
}
