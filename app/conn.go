/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package app
 *@file    map
 *@date    2025/6/6 16:06
 */

package app

import (
	"StarRocksAPIs/conn"
	"StarRocksAPIs/util"
	"errors"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"sync"
)

// 全局连接池（使用 sync.Map 或普通 Map + 互斥锁）
var syncmap sync.Map // key: 指标名（如 "adhoc"）, value: *gorm.DB

func setConnectMap(app string) error {
	//starrocks
	db, err := conn.StarRocks(app)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return err
	}
	syncmap.Store(app+"db", db)
	return nil
}
func setRestyMap(app, user, password string) error {
	//resty
	restys := resty.New().SetLogger(&util.CustomLogger{}).SetBasicAuth(user, password)
	syncmap.Store(app+"resty", restys)
	return nil
}

func getmapConnect(app string) (*gorm.DB, error) {
	// getConnectMap 获取指定指标的数据库连接
	if db, ok := syncmap.Load(app + "db"); ok {
		return db.(*gorm.DB), nil
	}
	return nil, errors.New("sync.Map connect db is err")
}
func getmapResty(app string) (*resty.Client, error) {
	// getConnectMap 获取指定指标的数据库连接
	if restys, ok := syncmap.Load(app + "resty"); ok {
		return restys.(*resty.Client), nil
	}
	return nil, errors.New("sync.Map connect resty is err")
}

func _init() {
	for _, m := range util.MetaLink {
		err := setConnectMap(m["app"].(string))
		if err != nil {
			util.Loggrs.Error(err.Error())
		}
		err = setRestyMap(m["app"].(string), m["user"].(string), m["password"].(string))
		if err != nil {
			util.Loggrs.Error(err.Error())
		}
	}
}
