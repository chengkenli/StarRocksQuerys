/*
 *@author  chengkenli
 *@project StarRocksApp
 *@package conn
 *@file    ConnectMySQL
 *@date    2024/9/27 16:08
 */

package conn

import (
	"StarRocksAPIs/util"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func ConnectMySQL() (*gorm.DB, error) {
	newLogger := logger.New(nil,
		logger.Config{
			SlowThreshold: time.Second * 1000, // 控制慢SQL阈值
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&charset=utf8mb4&loc=Local",
		util.Config.GetString("metadb.user"),
		util.Config.GetString("metadb.password"),
		util.Config.GetString("metadb.host"),
		util.Config.GetInt("metadb.port"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println(err)
	}
	return db, err
}
