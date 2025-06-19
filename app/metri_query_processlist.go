/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    analy_processlist
 *@date    2025/5/26 13:40
 */

package app

import (
	"StarRocksQuerys/conn"
	"StarRocksQuerys/meta"
	"StarRocksQuerys/tools"
	"StarRocksQuerys/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type queryResult struct {
	Id         string   `json:"id"`
	User       string   `json:"user"`
	Host       string   `json:"host"`
	Clientuser string   `json:"clientuser"`
	Db         string   `json:"db"`
	Command    string   `json:"command"`
	Time       int      `json:"time"`
	State      string   `json:"state"`
	IsPending  bool     `json:"isPending"`
	Info       string   `json:"info"`
	Warehouse  string   `json:"warehouse"`
	IsSlow     bool     `json:"isSlow"`     // 标记慢查询(600-1500ms)
	IsCritical bool     `json:"isCritical"` // 标记高耗时查询(>1500ms)
	Uri        string   `json:"uri"`
	Whitelist  bool     `json:"whitelist"`
	Feip       string   `json:"feip"`
	Shortlist  bool     `json:"shortlist"`
	Ctxip      []string `json:"ctxip"`
	Gethour    string   `json:"gethour"`
	Admin      bool     `json:"admin"`
	Point      bool     `json:"point"`
}
type dataItem struct {
	Data  []queryResult `json:"data"`
	Run   int           `json:"run"`
	Pend  int           `json:"pend"`
	Sleep int           `json:"sleep"`
	Count int           `json:"count"`
	Fe    []string      `json:"fe"`
	Title string        `json:"title"`
}

func (engine *threadMap) processlist(c *gin.Context) {
	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}
	appid = setdefault(appid)
	db, err := engine.getmapConnect(appid)
	if err != nil {
		util.Loggrs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// get
	if v, ok := querycache.Get(appid + "data"); ok {
		pend, _ := querycache.Get(appid + "pend")
		sleep, _ := querycache.Get(appid + "sleep")
		fe, _ := querycache.Get(appid + "fe")

		c.JSON(http.StatusOK, gin.H{
			"data":  sortResult(v.([]queryResult)),
			"run":   len(v.([]queryResult)),
			"pend":  pend.(int),
			"sleep": sleep.(int),
			"count": len(v.([]queryResult)) + sleep.(int),
			"fe":    fe.([]string),
		})
		util.Loggrs.Info(c.ClientIP(), " response:cache")
		return
	}

	var felist []string
	var m []map[string]interface{}
	r := db.Raw("show frontends").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Warn(r.Error.Error())
		return
	}

	var mfe []map[string]string
	for _, item := range m {
		if item["Alive"].(string) != "true" {
			continue
		}
		mfe = append(mfe, map[string]string{
			item["IP"].(string): item["Role"].(string),
		})
		felist = append(felist, item["IP"].(string))
	}

	// http sql api
	if tools.Version(db) >= 3.2 {
		c.JSON(http.StatusOK, engine.httpapi(db, appid, felist))
		util.Loggrs.Info(c.ClientIP(), " response:actual")
		return
	}
	// end

	var runs, sleep, pend []string
	var result []queryResult
	for _, ip := range felist {
		single, err := conn.StarRocksSingle(appid, ip)
		if err != nil {
			util.Loggrs.Warn(err.Error())
			continue
		}

		var dbresult []map[string]interface{}
		r = single.Raw("show processlist").Scan(&dbresult)
		if r.Error != nil {
			util.Loggrs.Warn(r.Error.Error())
			continue
		}

		var feresult []string
		for _, m := range dbresult {
			if m["Command"].(string) == "Sleep" {
				sleep = append(sleep, "1")
				continue
			}
			var isslow, iscritical bool
			if m["Time"].(int64) >= 600 && m["Time"].(int64) < 1500 {
				isslow = true
			} else if m["Time"].(int64) >= 1500 {
				iscritical = true
			}
			ok, _ := strconv.ParseBool(m["IsPending"].(string))
			if ok {
				pend = append(pend, "1")
			}
			feresult = append(feresult, "1")

			// 检查白名单标记
			whitelist := false
			if util.Read.Schema.Whitelist != "" {
				list := strings.Split(util.Read.Schema.Whitelist, ",")
				if tools.StrInSlice(m["User"].(string), list) {
					whitelist = true
				}
			}

			// 检查管理员标记
			adminsign := false
			if util.Read.Schema.Role.Admin != "" {
				list := strings.Split(util.Read.Schema.Role.Admin, ",")
				if tools.StrInSlice(m["User"].(string), list) {
					adminsign = true
				}
			}

			// 检查二级重点账号标记
			pointsign := false
			if util.Read.Schema.Role.Point != "" {
				list := strings.Split(util.Read.Schema.Role.Point, ",")
				if tools.StrInSlice(m["User"].(string), list) {
					pointsign = true
				}
			}

			clientip := strings.Split(m["Host"].(string), ":")[0]
			var domainname []string
			if v, ok := apicache.Get("ctx" + clientip); ok {
				domainname = v.([]string)
			} else {
				domainname = ctxIp(clientip)
			}
			var clientname string
			if domainname != nil {
				clientname = getipname(strings.Split(domainname[0], ".")[0])
			}

			var feip string
			for _, m2 := range mfe {
				if v, ok := m2[ip]; ok {
					role := strings.ToLower(v)
					switch role {
					case "follower":
						feip = fmt.Sprintf("😁%s[%s]", ip, strings.ToLower(v))
					case "leader":
						feip = fmt.Sprintf("😈%s[%s]", ip, strings.ToLower(v))
					default:
						feip = fmt.Sprintf("😱%s[%s]", ip, strings.ToLower(v))
					}
				}
			}

			shorts := false
			if v, ok := shortcache.Get(appid + "short"); ok {
				if v != nil {
					if tools.StrInSlice(m["User"].(string), v.([]string)) {
						shorts = true
					}
				}
			}
			labela := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%d/getstmtid?app=%s&id=%d')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%d</a>`,
				util.H.Ip, util.Read.Server.Port, appid, m["Id"].(int64), m["Id"].(int64))

			idname := fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%d">%s`, m["Id"].(int64), labela)
			// role juct
			if util.Read.Schema.Role.Admin != "" {
				role := strings.Split(util.Read.Schema.Role.Admin, ",")
				if tools.StrInSlice(m["User"].(string), role) {
					idname = fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%d">%s👑`, m["Id"].(int64), labela)
				}
			}
			if util.Read.Schema.Role.Point != "" {
				role := strings.Split(util.Read.Schema.Role.Point, ",")
				if tools.StrInSlice(m["User"].(string), role) {
					idname = fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%d">%s😏`, m["Id"].(int64), labela)
				}
			}
			result = append(result, queryResult{
				Id:         idname,
				User:       m["User"].(string),
				Host:       clientip,
				Clientuser: clientname,
				Db:         m["Db"].(string),
				Command:    fmt.Sprintf(`<button id="kill-connectid" class="btn btn-light btn-sm ms-2" data-connectid="%d">❌</button>`, m["Id"].(int64)),
				Time:       int(m["Time"].(int64)),
				State:      m["State"].(string),
				IsPending:  ok,
				Info:       m["Info"].(string),
				Warehouse:  fmt.Sprintf(`<button id="kill-disconnect" class="btn btn-light btn-sm ms-2" data-disconnect="%s">⛔</button>`, m["User"].(string)),
				IsSlow:     isslow,
				IsCritical: iscritical,
				Whitelist:  whitelist,
				Feip:       fmt.Sprintf("fe:(%s) max_connections(%s)", feip, getmaxconnections(appid, m["User"].(string), db)),
				Shortlist:  shorts,
				Ctxip:      domainname,
				Gethour:    tools.GetHour(int(m["Time"].(int64))),
				Admin:      adminsign,
				Point:      pointsign,
			})
		}

		for _, m2 := range mfe {
			if v, ok := m2[ip]; ok {
				var msg string
				role := strings.ToLower(v)
				switch role {
				case "follower":
					msg = fmt.Sprintf("😁%s[%s](%d)", ip, strings.ToLower(v), len(feresult))
				case "leader":
					msg = fmt.Sprintf("😈%s[%s](%d)", ip, strings.ToLower(v), len(feresult))
				default:
					msg = fmt.Sprintf("😱%s[%s](%d)", ip, strings.ToLower(v), len(feresult))
				}
				runs = append(runs, msg)
			}
		}
		//runs = append(runs, fmt.Sprintf("%s(%d)", msg, len(feresult)))
		/*每次使用完，主动关闭连接数*/
		sqlDB, err := single.DB()
		if err != nil {
			util.Loggrs.Error(err.Error())
			continue
		}
		sqlDB.SetMaxOpenConns(10)                 //最大连接数
		sqlDB.SetMaxIdleConns(3)                  //最大空闲连接数
		sqlDB.SetConnMaxLifetime(5 * time.Second) //空闲连接最多存活时间
		sqlDB.Close()
	}

	// 载入缓存
	go func() {
		// set
		querycache.Set(appid+"data", result, cache.DefaultExpiration)
		querycache.Set(appid+"pend", len(pend), cache.DefaultExpiration)
		querycache.Set(appid+"sleep", len(sleep), cache.DefaultExpiration)
		querycache.Set(appid+"fe", runs, cache.DefaultExpiration)

		if appid == "sr-cdp" || appid == "sr-api" || appid == "sr-ma" {
			return
		}
		if _, ok := shortcache.Get(appid + "short"); !ok {
			shortdata := meta.ShortQuery(db)
			if shortdata != nil {
				shortcache.Set(appid+"short", shortdata, cache.DefaultExpiration)
			}
		}
	}()
	util.Loggrs.Info(c.ClientIP(), " response:actual")

	c.JSON(http.StatusOK, gin.H{
		"data":  sortResult(result),
		"run":   len(result),
		"pend":  len(pend),
		"sleep": len(sleep),
		"count": len(result) + len(sleep),
		"fe":    runs,
		"title": gettitle(db, appid),
	})
}

func sortResult(results []queryResult) []queryResult {
	// 标记慢查询和关键查询
	for i := range results {
		if results[i].Time > 1500 {
			results[i].IsCritical = true
			results[i].IsSlow = true
		} else if results[i].Time > 600 {
			results[i].IsSlow = true
		}
	}
	// 按Time降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Time > results[j].Time
	})
	return results
}

func ctxIp(ip string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	domains, err := net.DefaultResolver.LookupAddr(ctx, ip)
	if err != nil {
		return nil
	}
	apicache.Set("ctx"+ip, domains, cache.DefaultExpiration)
	return domains
}

func (engine *threadMap) httpapi(db *gorm.DB, app string, felist []string) dataItem {
	restys, err := engine.getmapResty(app)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return dataItem{}
	}
	var count int
	var result []queryResult
	var running, pending, query, sleep, femsg []string
	for _, fe := range felist {
		uri := fmt.Sprintf("http://%s:8030/api/v1/catalogs/default_catalog/databases/information_schema/sql", fe)
		respones, err := restys.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"query": "show processlist",
			}).Post(uri)
		if err != nil {
			util.Loggrs.Error(err.Error())
			continue
		}
		var m util.HttpSqlApi
		err = json.Unmarshal(respones.Body(), &m)
		if err != nil {
			util.Loggrs.Error(err.Error())
			continue
		}
		for _, datum := range m.Data {
			// 判断队列堵塞
			switch datum.IsPending {
			case "true":
				pending = append(pending, datum.ID)
			case "false":
				running = append(running, datum.ID)
			}
			// 判断睡眠连接
			switch datum.Command {
			case "Sleep":
				sleep = append(sleep, datum.ID)
				continue
			case "Query":
				query = append(query, datum.ID)
			}

			// id 意图判断
			labela := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%d/getstmtid?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
				util.H.Ip, util.Read.Server.Port, app, datum.ID, datum.ID)
			idname := fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%s">%s`, datum.ID, labela)
			// role juct
			if util.Read.Schema.Role.Admin != "" {
				role := strings.Split(util.Read.Schema.Role.Admin, ",")
				if tools.StrInSlice(datum.User, role) {
					idname = fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%s">%s👑`, datum.ID, labela)
				}
			}
			if util.Read.Schema.Role.Point != "" {
				role := strings.Split(util.Read.Schema.Role.Point, ",")
				if tools.StrInSlice(datum.User, role) {
					idname = fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%s">%s😏`, datum.ID, labela)
				}
			}
			// 类型转换
			Time, _ := strconv.Atoi(datum.Time)
			isPending, _ := strconv.ParseBool(datum.IsPending)

			// 耗时计算
			var isslow, iscritical bool
			if Time >= 600 && Time < 1500 {
				isslow = true
			} else if Time >= 1500 {
				iscritical = true
			}

			// 物理地址追踪
			clientip := strings.Split(datum.Host, ":")[0]
			var domainname []string
			if v, ok := apicache.Get("ctx" + clientip); ok {
				domainname = v.([]string)
			} else {
				domainname = ctxIp(clientip)
			}
			var clientname string
			if domainname != nil {
				clientname = getipname(strings.Split(domainname[0], ".")[0])
			}

			// 白名单匹配
			whitelist := false
			if util.Read.Schema.Whitelist != "" {
				list := strings.Split(util.Read.Schema.Whitelist, ",")
				if tools.StrInSlice(datum.User, list) {
					whitelist = true
				}
			}

			// 检查管理员标记
			adminsign := false
			if util.Read.Schema.Role.Admin != "" {
				list := strings.Split(util.Read.Schema.Role.Admin, ",")
				if tools.StrInSlice(datum.User, list) {
					adminsign = true
				}
			}

			// 检查二级重点账号标记
			pointsign := false
			if util.Read.Schema.Role.Point != "" {
				list := strings.Split(util.Read.Schema.Role.Point, ",")
				if tools.StrInSlice(datum.User, list) {
					pointsign = true
				}
			}

			// 短查询匹配
			shorts := false
			if v, ok := shortcache.Get(app + "short"); ok {
				if v != nil {
					if tools.StrInSlice(datum.User, v.([]string)) {
						shorts = true
					}
				}
			}

			result = append(result, queryResult{
				Id:         idname,
				User:       datum.User,
				Host:       clientip,
				Clientuser: clientname,
				Db:         datum.Db,
				Command:    fmt.Sprintf(`<button id="kill-connectid" class="btn btn-light btn-sm ms-2" data-connectid="%s">❌</button>`, datum.ID),
				Time:       Time,
				State:      datum.State,
				IsPending:  isPending,
				Info:       datum.Info,
				Warehouse:  fmt.Sprintf(`<button id="kill-disconnect" class="btn btn-light btn-sm ms-2" data-disconnect="%s">⛔</button>`, datum.User),
				IsSlow:     isslow,
				IsCritical: iscritical,
				Whitelist:  whitelist,
				Feip:       fmt.Sprintf("fe:(%s) max_connections(%s)", fe, getmaxconnections(app, datum.User, db)),
				Shortlist:  shorts,
				Ctxip:      domainname,
				Gethour:    fmt.Sprintf("%s(%s)", datum.ConnectionStartTime, tools.GetHour(Time)),
				Admin:      adminsign,
				Point:      pointsign,
			})
		}
		count = +m.Statistics.ReturnRows
		femsg = append(femsg, fmt.Sprintf("%s (%d)(%d)(%d)", fe, len(query), len(pending), len(sleep)))
	}

	// 载入缓存
	go func() {
		// set
		querycache.Set(app+"data", result, cache.DefaultExpiration)
		querycache.Set(app+"pend", len(pending), cache.DefaultExpiration)
		querycache.Set(app+"sleep", len(sleep), cache.DefaultExpiration)
		querycache.Set(app+"fe", femsg, cache.DefaultExpiration)

		if app == "sr-cdp" || app == "sr-api" || app == "sr-ma" {
			return
		}
		if _, ok := shortcache.Get(app + "short"); !ok {
			shortdata := meta.ShortQuery(db)
			if shortdata != nil {
				shortcache.Set(app+"short", shortdata, cache.DefaultExpiration)
			}
		}
	}()

	// 统计汇报
	return dataItem{
		Data:  sortResult(result),
		Run:   len(query),
		Pend:  len(pending),
		Sleep: len(sleep),
		Count: count,
		Fe:    femsg,
		Title: gettitle(db, app),
	}
}

func getmaxconnections(app, username string, db *gorm.DB) string {
	if v, ok := apicache.Get(app + username + "max_connections"); ok {
		return v.(string)
	}
	var m map[string]interface{}
	r := db.Raw(fmt.Sprintf("SHOW PROPERTY FOR '%s' LIKE 'max_user_connections'", username)).Scan(&m)
	if r.Error != nil {
		return ""
	}
	if m["Key"] == "max_user_connections" {
		go apicache.Set(app+username+"max_connections", m["Value"], cache.DefaultExpiration)
		return m["Value"].(string)
	}
	return ""
}
