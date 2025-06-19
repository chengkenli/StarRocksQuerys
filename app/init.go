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
	"strconv"
	"strings"
	"time"
)

var querycache = cache.New(5*time.Second, 10*time.Second)
var shortcache = cache.New(5*time.Minute, 10*time.Minute)
var apicache = cache.New(12*time.Hour, 24*time.Hour)
var errcache = cache.New(1*time.Hour, 2*time.Hour)
var engine *threadMap

func init() {
	// 初始化对象引擎
	engine = newPool()
	engine._init()
	engine._init2()
}

func (engine *threadMap) _init2() {
	if util.Read.Schema.Ipsystem == "" {
		return
	}
	if util.Read.Schema.Ipapp == "" {
		return
	}
	go func() {
		if util.ClientIPDec == nil {
			db, err := engine.getmapConnect(util.Read.Schema.Ipapp)
			if err != nil {
				util.Loggrs.Error(err.Error())
				util.Loggrs.Fatalf(fmt.Sprintf("校验失败！配置文件中的主集群 ipapp:%s 在元数据表%s 中匹配不到连接信息，请确保集群名称一致！", util.Read.Schema.Ipapp, util.Read.Metadb.Base))
				return
			}
			getsign(db)
		}

		ticker := time.NewTicker(time.Hour * 5)
		for {
			select {
			case <-ticker.C:
				db, err := engine.getmapConnect(util.Read.Schema.Ipapp)
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

// 假如前端传过来一个默认值，那么给他返回第一个集群
func setdefault(app string) string {
	var reals string
	var applist []string
	if app == "sr-default" {
		for _, m := range util.MetaLink {
			applist = append(applist, m["app"].(string))
		}
		reals = applist[0]
	} else {
		reals = app
	}
	return reals
}

func gettitle(db *gorm.DB, app string) string {
	if val, ok := apicache.Get(app + "title"); ok {
		return val.(string)
	}
	var nick string
	for _, m := range util.MetaLink {
		if m["app"].(string) == app {
			nick = m["nickname"].(string)
		}
	}
	var (
		f util.Fronends
		b util.Backends
		e util.Computes
	)
	db.Raw("show frontends").Scan(&f)

	var backends int
	db.Raw("show backends").Scan(&b)
	backends = len(b)
	if backends == 0 {
		db.Raw("show compute nodes").Scan(&e)
		backends = len(e)
	}
	title := fmt.Sprintf("%s FE:(%d) BE:(%d) VERSION:(%s)", nick, len(f), backends, versions(db))
	go apicache.Set(app+"title", title, cache.DefaultExpiration)
	return title
}

func versions(db *gorm.DB) string {
	type Version struct {
		Version string `bson:"version"`
	}
	/*匹配starrocks版本*/
	var v Version
	sql := fmt.Sprintf("select current_version() as version")
	db.Raw(sql).Scan(&v)
	/*end*/
	return strings.Split(v.Version, " ")[0]
}

func storage(sr string, infos util.Backends) string {
	var f, g, h, k float64
	for _, info := range infos {
		f = f + flos(info.TotalCapacity)
		k = k + flos(info.DataUsedCapacity)
		g = g + flos(info.AvailCapacity)
	}
	h = f - g
	return fmt.Sprintf("%s总存储:%0.2ftb, 目前存储:%0.2ftb(数据实际:%0.2ftb), 空闲:%0.2ftb, 百分比:%0.2f%%", sr, f, h, k, f-h, h/f*100)
}

func flos(s string) float64 {
	var maxDiskUsedPct float64
	if strings.Contains(strings.ToLower(strings.ReplaceAll(s, " ", "")), "gb") {
		m := strings.Split(s, " ")[0]
		maxDiskUsedPct, _ = strconv.ParseFloat(m, 64)
		maxDiskUsedPct = maxDiskUsedPct / 1024
		return maxDiskUsedPct
	}
	m := strings.Split(s, " ")[0]
	maxDiskUsedPct, _ = strconv.ParseFloat(m, 64)
	return maxDiskUsedPct
}
