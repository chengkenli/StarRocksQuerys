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
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/patrickmn/go-cache"
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
}

func processlist(c *gin.Context) {
    appid := c.GetHeader("AppID")
    //_, db := connect(appid)
    db, err := getmapConnect(appid)
    if err != nil {
        c.JSON(http.StatusInternalServerError, nil)
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
            return
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

            whitelist := false
            if util.Config.GetString("schema.whitelist") != "" {
                list := strings.Split(util.Config.GetString("schema.whitelist"), ",")
                if tools.StrInSlice(m["User"].(string), list) {
                    whitelist = true
                }
            }

            clientip := strings.Split(m["Host"].(string), ":")[0]
            clientname := getipname(clientip)

            var ctxip []string
            if v, ok := apicache.Get("ctx" + clientip); ok {
                ctxip = v.([]string)
            } else {
                ctxip = ctxIp(clientip)
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
            labela := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%s/getstmtid?app=%s&id=%d')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%d</a>`,
                util.H.Ip, util.Config.GetString("server.port"), appid, m["Id"].(int64), m["Id"].(int64))

            idname := fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%d">%s`, m["Id"].(int64), labela)
            // role juct
            if util.Config.GetString("schema.role.admin") != "" {
                role := strings.Split(util.Config.GetString("schema.role.admin"), ",")
                if tools.StrInSlice(m["User"].(string), role) {
                    idname = fmt.Sprintf(`<input type="checkbox" class="query-checkbox" data-id="%d">%s👑`, m["Id"].(int64), labela)
                }
            }
            if util.Config.GetString("schema.role.point") != "" {
                role := strings.Split(util.Config.GetString("schema.role.point"), ",")
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
                Warehouse:  m["Warehouse"].(string),
                IsSlow:     isslow,
                IsCritical: iscritical,
                Whitelist:  whitelist,
                Feip:       feip,
                Shortlist:  shorts,
                Ctxip:      ctxip,
                Gethour:    tools.GetHour(int(m["Time"].(int64))),
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
            return
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
