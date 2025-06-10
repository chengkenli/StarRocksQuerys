/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri
 *@date    2025/5/27 20:23
 */

package app

import (
    "StarRocksQuerys/util"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/patrickmn/go-cache"
    "net/http"
    "time"
)

type slowData struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    StartTime string `json:"starttime"`
    ExecTime  int64  `json:"exectime"`
    Type      string `json:"type"`
    State     string `json:"state"`
    Clear     int64  `json:"clear"`
    Warn      int64  `json:"warn"`
    Kill      int64  `json:"kill"`
}

func metriSlowquery(c *gin.Context) {
    appid := c.GetHeader("AppID")
    // get
    if v, ok := querycache.Get(appid + "slowquery"); ok {
        c.JSON(http.StatusOK, v.([]slowData))
        return
    }
    slow_query_data := execslow(appid)
    // 载入缓存
    go func() {
        // set
        querycache.Set(appid+"slowquery", slow_query_data, cache.DefaultExpiration)
    }()
    c.JSON(http.StatusOK, slow_query_data)
}

func execslow(app string) []slowData {
    if util.Config.GetString("schema.slowstmt") == "" {
        return nil
    }
    db, err := getmapConnect(util.Config.GetString("schema.ipapp"))
    if err != nil {
        return nil
    }
    stmt := fmt.Sprintf("select queryId,user,timestamp,queryTime,action,logfile from %s where app='%s' and timestamp>=date_sub(now(), interval 1 day) order by timestamp desc", util.Config.GetString("schema.slowstmt"), app)
    var m []map[string]interface{}
    r := db.Raw(stmt).Scan(&m)
    if r.Error != nil {
        util.Loggrs.Error(r.Error.Error())
        return nil
    }

    clears, warns, kills := execslowdata(app)

    var slows []slowData
    for _, m2 := range m {
        var comment, state string
        switch m2["action"].(int64) {
        case 0:
            comment = "异常停留"
            state = "clear"
        case 1:
            comment = "违规参数"
            state = "kill"
        case 2:
            comment = "10分钟"
            state = "warn"
        case 3:
            comment = "30分钟"
            state = "kill"
        case 4:
            comment = "全表扫描"
            state = "kill"
        case 5:
            comment = "TB字节"
            state = "kill"
        case 6:
            comment = "百亿行数"
            state = "kill"
        case 7:
            comment = "CATALOG"
            state = "kill"
        case 8:
            comment = "GB内存"
            state = "kill"
        default:
            comment = "Other"
            state = "non"
        }

        var uri string
        if m2["logfile"] != nil {
            uri = fmt.Sprintf(`<a href="%s" target="_blank" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`, m2["logfile"].(string), m2["queryId"].(string))
        } else {
            uri = m2["queryId"].(string)
        }
        slows = append(slows, slowData{
            ID:        uri,
            Name:      m2["user"].(string),
            StartTime: m2["timestamp"].(time.Time).Format("2006-01-02 15:04:05"),
            ExecTime:  m2["queryTime"].(int64),
            Type:      comment,
            State:     state,
            Clear:     clears,
            Warn:      warns,
            Kill:      kills,
        })
    }
    return slows
}

func execslowdata(app string) (int64, int64, int64) {
    if util.Config.GetString("schema.slowstmt") == "" {
        return -1, -1, -1
    }

    db, err := getmapConnect(util.Config.GetString("schema.ipapp"))
    if err != nil {
        return -1, -1, -1
    }

    stmt := fmt.Sprintf("select count(*) as count,action from %s where app='%s' and timestamp>=date_sub(now(), interval 1 day) group by action order by 1 desc", util.Config.GetString("schema.slowstmt"), app)
    var m []map[string]interface{}
    r := db.Raw(stmt).Scan(&m)
    if r.Error != nil {
        util.Loggrs.Warn(r.Error.Error())
        return -1, -1, -1
    }
    var clears, warns, kills int64
    for _, m2 := range m {
        switch m2["action"].(int64) {
        case 0:
            clears = +m2["count"].(int64)
        case 2:
            warns = +m2["count"].(int64)
        default:
            kills = +m2["count"].(int64)
        }
    }
    return clears, warns, kills
}
