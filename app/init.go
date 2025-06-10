/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package app
 *@file    init
 *@date    2025/5/22 15:54
 */

package app

import (
	"StarRocksAPIs/util"
	"fmt"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"strings"
	"time"
)

var querycache = cache.New(5*time.Second, 10*time.Second)
var shortcache = cache.New(5*time.Minute, 10*time.Minute)
var apicache = cache.New(12*time.Hour, 24*time.Hour)
var errcache = cache.New(1*time.Hour, 2*time.Hour)

func init() {
	_init()
	_init2()
}

func _init2() {
	if util.Config.GetString("schema.ipsystem") == "" {
		return
	}
	if util.Config.GetString("schema.ipapp") == "" {
		return
	}
	go func() {
		if util.ClientIPDec == nil {
			db, err := getmapConnect(util.Config.GetString("schema.ipapp"))
			if err != nil {
				util.Loggrs.Error(err.Error())
				return
			}
			getsign(db)
		}

		ticker := time.NewTicker(time.Hour * 5)
		for {
			select {
			case <-ticker.C:
				db, err := getmapConnect(util.Config.GetString("schema.ipapp"))
				if err != nil {
					util.Loggrs.Error(err.Error())
					return
				}
				getsign(db)
			}
		}
	}()
}

func getsign(db *gorm.DB) {
	if util.Config.GetString("schema.ipsystem") == "" {
		return
	}
	r := db.Raw("select * from " + util.Config.GetString("schema.ipsystem")).Scan(&util.ClientIPDec)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return
	}
	util.Loggrs.Info("[ok] 初始化加载ipsystem缓存", len(util.ClientIPDec))
}
func getusername(ip string) string {
	var message string
	for _, data := range util.ClientIPDec {
		if data.Ip == ip {
			message = fmt.Sprintf("%s(%s)", ip, data.User)
			break
		}
	}
	if message == "" {
		message = ip
	}
	return message
}

func getipname(ip string) string {
	for _, data := range util.ClientIPDec {
		if data.Ip == ip {
			return strings.ToLower(data.User)
		}
	}
	return ""
}

func leader(db *gorm.DB) string {
	var m []map[string]interface{}
	r := db.Raw("show frontends").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Warn(r.Error.Error())
		return ""
	}
	var leaderip string
	for _, item := range m {
		if item["Alive"].(string) != "true" {
			continue
		}
		if item["Role"].(string) == "LEADER" {
			leaderip = item["IP"].(string)
			break
		}
	}
	return leaderip
}
