/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    init
 *@date    2025/5/22 15:54
 */

package app

import (
    "StarRocksQuerys/tools"
    "StarRocksQuerys/util"
    "fmt"
    "github.com/patrickmn/go-cache"
    "gorm.io/gorm"
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
    if util.Read.Schema.Ipsystem == "" {
        return
    }
    if util.Read.Schema.Ipapp == "" {
        return
    }
    go func() {
        if util.ClientIPDec == nil {
            db, err := getmapConnect(util.Read.Schema.Ipapp)
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
                db, err := getmapConnect(util.Read.Schema.Ipapp)
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
    if util.Read.Schema.Ipsystem == "" {
        return
    }
    r := db.Raw("select * from " + util.Read.Schema.Ipsystem).Scan(&util.ClientIPDec)
    if r.Error != nil {
        util.Loggrs.Error(r.Error.Error())
        return
    }
    util.Loggrs.Info("[ok] 初始化加载ipsystem缓存", len(util.ClientIPDec))
}

func getipname(device_name string) string {
    for _, data := range util.ClientIPDec {
        if data.ComputerName == device_name {
            return fmt.Sprintf("%s (%s)(%s) %s(%s)", data.UserName, data.ComputerType, data.ComputerStatus, data.ComputerVersion, data.Brand)
        } else if tools.StrInSlice(device_name, util.Read.Schema.Gybi) {
            return fmt.Sprintf("观远bi")
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
