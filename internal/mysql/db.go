// Package mysql mysql connect
package mysql

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dingdayu/dnsx/pkg/config"
	"github.com/dingdayu/dnsx/pkg/log"

	"github.com/jinzhu/gorm"
	// mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var once sync.Once
var db *gorm.DB

// Init db connect
func Init() {
	var err error
	once.Do(func() {
		db, err = gorm.Open("mysql", config.GetString("mysql.dns"))
		if err == nil {
			fmt.Println("\033[1;30;42m[info]\033[0m db [maste] connect success")
		} else {
			fmt.Printf("\033[1;30;41m[error]\033[0m db [master] connect error: %s", err.Error())
			os.Exit(1)
		}
		db.SetLogger(log.New())
		db.LogMode(config.GetBool("mysql.log_model"))
		db.DB().SetConnMaxLifetime(time.Minute * time.Duration(config.GetInt("mysql.conn_max_lifetime")))
		db.DB().SetMaxIdleConns(config.GetInt("mysql.max_idle_conn"))
		db.DB().SetMaxOpenConns(config.GetInt("mysql.max_open_conn"))
	})
}

// GetDB get db connect
func GetDB() *gorm.DB {
	return db
}
