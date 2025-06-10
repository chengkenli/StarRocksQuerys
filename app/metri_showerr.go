/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_showerr
 *@date    2025/6/5 16:15
 */

package app

import (
    "StarRocksQuerys/tools"
    "StarRocksQuerys/util"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/go-resty/resty/v2"
    "github.com/patrickmn/go-cache"
    "gorm.io/gorm"
    "net/http"
    "os"
    "sort"
    "strings"
    "sync"
    "time"
)

type showErr struct {
    Starttime string `json:"starttime"`
    User      string `json:"user"`
    Queryid   string `json:"queryid"`
    Errmsg    string `json:"errmsg"`
    Errinfo   string `json:"errinfo"`
    Stmt      string `json:"stmt"`
}

func ShowErr(c *gin.Context) {
    appid := c.GetHeader("AppID")

    db, err := getmapConnect(appid)
    if err != nil {
        util.Loggrs.Error(err.Error())
        return
    }

    errs := showErrors(appid, db)
    c.JSON(http.StatusOK, sortByStartTimeDesc(errs))
}

func init_meta(app string) util.ConnectParms {
    var avg util.ConnectParms
    // connect
    if len(util.MetaLink) == 0 {
        return util.ConnectParms{}
    }
    for _, m := range util.MetaLink {
        if m["app"].(string) == app {
            var manager_access_key, manager_secret_key string
            if m["manager_access_key"] != nil {
                manager_access_key = m["manager_access_key"].(string)
            }
            if m["manager_secret_key"] != nil {
                manager_secret_key = m["manager_secret_key"].(string)
            }
            avg = util.ConnectParms{
                Host:       m["feip"].(string),
                Port:       int(m["feport"].(int32)),
                User:       m["user"].(string),
                Pass:       m["password"].(string),
                Uri:        m["address"].(string),
                AaccessKey: manager_access_key,
                SecretKey:  manager_secret_key,
            }
        }
    }
    return avg
}

func showErrors(app string, db *gorm.DB) []showErr {
    if util.Config.GetString("schema.auditops") == "" {
        return nil
    }

    avg := init_meta(app)

    var mv map[string]interface{}
    db.Raw("select current_version() as version").Scan(&mv)
    version := tools.VToFloat(mv["version"].(string))

    var queryid, starttime string
    if version >= 2.5 {
        queryid = "queryId"
        starttime = "timestamp"
    } else {
        queryid = "query_id"
        starttime = "time"
    }
    os.Setenv("TZ", "Asia/Shanghai")
    var m []map[string]interface{}

    minute := 30
    columns := fmt.Sprintf("%s as queryid,%s as starttime,queryType,clientIp,user,authorizedUser,resourceGroup,catalog,db,state,errorCode,queryTime,scanBytes,scanRows,returnRows,cpuCostNs,memCostBytes,stmtId,isQuery,feIp,stmt,digest,planCpuCosts,planMemCosts,pendingTimeMs,recordTime", queryid, starttime)
    sql := fmt.Sprintf("select %s from %s where %s >= date_sub(now(), interval %d minute) and user != 'root' and state='ERR' order by starttime desc", columns, util.Config.GetString("schema.auditops"), starttime, minute)
    r := db.Raw(sql).Scan(&m)
    if r.Error != nil {
        return nil
    }
    if len(m) == 0 {
        return nil
    }

    var errs []showErr

    if version >= 3.3 {
        switch version {
        case 3.311, 3.313:
            msgs := scan1(app, avg, m)
            errs = append(errs, msgs...)
        case 3.39:
            msgs := scan2(app, avg, m)
            errs = append(errs, msgs...)
        default:
            msgs := scan0(app, avg, m)
            errs = append(errs, msgs...)
        }
    } else {
        if app == "sr-app" && version == 3.21 {
            msgs := scan1(app, avg, m)
            errs = append(errs, msgs...)
        } else if app == "sccts" && version == 3.21 {
            msgs := scan2(app, avg, m)
            errs = append(errs, msgs...)
        } else {
            msgs := scan0(app, avg, m)
            errs = append(errs, msgs...)
        }
    }
    return errs
}

// 这是最原始的方式，支持2.x-3.x
func scan0(app string, avg util.ConnectParms, m []map[string]interface{}) []showErr {
    /*-------------------------登录manager web-------------------------------*/
    //创建Resty客户端
    Client := resty.New()
    //发送POST请求并处理响应
    result, err := Client.R().SetBody(map[string]string{
        "name":     avg.User,
        "password": avg.Pass,
    }).Post(avg.Uri + "/api/user/login")
    if err != nil {
        util.Loggrs.Warn("err:", err.Error())
        return nil
    }
    if strings.Contains(string(result.Body()), "Failed") {
        util.Loggrs.Warn("login is failed.")
        return nil
    }
    /*-------------------------登录manager web end----------------------------*/
    var errs []showErr

    ch := make(chan struct{}, 3)
    var wg sync.WaitGroup
    for _, m2 := range m {
        wg.Add(1)
        go func(m2 map[string]interface{}) {
            defer func() {
                <-ch
                wg.Done()
            }()

            ch <- struct{}{}

            // cache get
            if v, ok := errcache.Get(app + m2["queryid"].(string)); ok {
                errs = append(errs, v.(showErr))
                return
            }

            /*新的合并模块---------重新利用Cookies请求license*/
            respones3, err := Client.R().Get(avg.Uri + "/api/query/detail/" + m2["queryid"].(string))
            if err != nil {
                util.Loggrs.Warn("请求失败:", err.Error())
                return
            }
            type mm struct {
                Code int `json:"code"`
                Data struct {
                    QueryDetail struct {
                        QueryID      string   `json:"queryId"`
                        User         string   `json:"user"`
                        Status       string   `json:"status"`
                        ErrorMessage string   `json:"errorMessage"`
                        StartTime    int      `json:"startTime"`
                        EndTime      int      `json:"endTime"`
                        TimeUsed     int      `json:"timeUsed"`
                        SQL          string   `json:"sql"`
                        ExplainText  []string `json:"explainText"`
                        ProfileText  []string `json:"profileText"`
                    } `json:"queryDetail"`
                } `json:"data"`
            }

            var m mm
            err = json.Unmarshal(respones3.Body(), &m)
            if err != nil {
                util.Loggrs.Warn(err.Error())
                return
            }
            if m.Data.QueryDetail.ErrorMessage == "" {
                return
            }

            queryId := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%s/getstmtid?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
                util.H.Ip, util.Config.GetString("server.port"), app, m2["queryid"].(string), m2["queryid"].(string))

            sign := showErr{
                Starttime: m2["starttime"].(time.Time).Format("2006-01-02 15:04:05"),
                User:      strings.ReplaceAll(m2["user"].(string), "default_cluster:", ""),
                Queryid:   queryId,
                Errmsg:    m.Data.QueryDetail.ErrorMessage,
                Errinfo:   "",
                Stmt:      m.Data.QueryDetail.SQL,
            }
            if sign.Queryid == "" {
                return
            }
            // cache set
            if _, ok := errcache.Get(app + m2["queryid"].(string)); !ok {
                errcache.Set(app+m2["queryid"].(string), sign, cache.DefaultExpiration)
            }
            errs = append(errs, sign)

        }(m2)
    }
    wg.Wait()

    return errs
}

// 这里也是最原始的方式，但这里的密码是加了base64的
func scan2(app string, avg util.ConnectParms, m []map[string]interface{}) []showErr {
    /*-------------------------登录manager web-------------------------------*/
    //创建Resty客户端
    Client := resty.New()
    //发送POST请求并处理响应
    result, err := Client.R().SetBody(map[string]string{
        "name":     avg.User,
        "password": base64.StdEncoding.EncodeToString([]byte(avg.Pass)),
    }).Post(avg.Uri + "/api/user/login")
    if err != nil {
        util.Loggrs.Warn("err:", err.Error())
        return nil
    }
    if strings.Contains(string(result.Body()), "Failed") {
        util.Loggrs.Warn("login is failed.")
        return nil
    }
    /*-------------------------登录manager web end----------------------------*/
    var errs []showErr

    ch := make(chan struct{}, 3)
    var wg sync.WaitGroup
    for _, m2 := range m {
        wg.Add(1)
        go func(m2 map[string]interface{}) {
            defer func() {
                <-ch
                wg.Done()
            }()

            ch <- struct{}{}

            // cache get
            if v, ok := errcache.Get(app + m2["queryid"].(string)); ok {
                errs = append(errs, v.(showErr))
                return
            }

            /*新的合并模块---------重新利用Cookies请求license*/
            respones3, err := Client.R().Get(avg.Uri + "/api/query/detail/" + m2["queryid"].(string))
            if err != nil {
                util.Loggrs.Warn("请求失败:", err.Error())
                return
            }
            type mm struct {
                Code int `json:"code"`
                Data struct {
                    QueryDetail struct {
                        QueryID      string   `json:"queryId"`
                        User         string   `json:"user"`
                        Status       string   `json:"status"`
                        ErrorMessage string   `json:"errorMessage"`
                        StartTime    int      `json:"startTime"`
                        EndTime      int      `json:"endTime"`
                        TimeUsed     int      `json:"timeUsed"`
                        SQL          string   `json:"sql"`
                        ExplainText  []string `json:"explainText"`
                        ProfileText  []string `json:"profileText"`
                    } `json:"queryDetail"`
                } `json:"data"`
            }

            var m mm
            err = json.Unmarshal(respones3.Body(), &m)
            if err != nil {
                util.Loggrs.Warn(err.Error())
                return
            }
            if m.Data.QueryDetail.ErrorMessage == "" {
                return
            }

            queryId := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%s/getstmtid?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
                util.H.Ip, util.Config.GetString("server.port"), app, m2["queryid"].(string), m2["queryid"].(string))

            sign := showErr{
                Starttime: m2["starttime"].(time.Time).Format("2006-01-02 15:04:05"),
                User:      strings.ReplaceAll(m2["user"].(string), "default_cluster:", ""),
                Queryid:   queryId,
                Errmsg:    m.Data.QueryDetail.ErrorMessage,
                Errinfo:   "",
                Stmt:      m.Data.QueryDetail.SQL,
            }
            if sign.Queryid == "" {
                return
            }
            // cache set
            if _, ok := errcache.Get(app + m2["queryid"].(string)); !ok {
                errcache.Set(app+m2["queryid"].(string), sign, cache.DefaultExpiration)
            }

            errs = append(errs, sign)

        }(m2)
    }
    wg.Wait()
    return errs
}

func scan1(app string, avg util.ConnectParms, m []map[string]interface{}) []showErr {
    if len(m) == 0 {
        return nil
    }
    var errs []showErr
    client := resty.New()
    ch := make(chan struct{}, 3)
    var wg sync.WaitGroup
    for _, m2 := range m {
        wg.Add(1)
        go func(m2 map[string]interface{}) {
            defer func() {
                <-ch
                wg.Done()
            }()

            ch <- struct{}{}

            // cache get
            if v, ok := errcache.Get(app + m2["queryid"].(string)); ok {
                errs = append(errs, v.(showErr))
                return
            }

            queryid := m2["queryid"].(string)
            sign := getOpenApi(app, client, avg, queryid)

            if sign.Queryid == "" {
                return
            }

            // cache set
            if _, ok := errcache.Get(app + m2["queryid"].(string)); !ok {
                errcache.Set(app+m2["queryid"].(string), sign, cache.DefaultExpiration)
            }

            errs = append(errs, sign)
        }(m2)
    }
    wg.Wait()
    return errs
}

func getOpenApi(app string, client *resty.Client, avg util.ConnectParms, queryid string) showErr {
    nonce := tools.RandomPassWord(32)
    unix := time.Now().Unix()

    AuthCredential := fmt.Sprintf("%s/%d/%s", avg.AaccessKey, unix, nonce)
    AuthCredentialEncrypted := util.HexEncrypt(AuthCredential, avg.SecretKey)
    AuthContent := fmt.Sprintf(`HTTPMethod:GET
CanonicalURI:/openapi/v1/query/record/%s
CanonicalQueryString:
CanonicalForm:`, queryid)
    AuthSignature := util.HexEncrypt(AuthContent, AuthCredentialEncrypted)
    AuthHeader := fmt.Sprintf("MANAGER-HMAC-SHA256 Credential=%s,Signature=%s", AuthCredential, AuthSignature)

    uri := fmt.Sprintf("%s/openapi/v1/query/record/%s", avg.Uri, queryid)
    //发送POST请求并处理响应
    respones, err := client.
        R().
        SetHeaders(map[string]string{
            "Content-Type":  "application/x-www-form-urlencoded",
            "Authorization": AuthHeader,
        }).
        Get(strings.NewReplacer(" ", "").Replace(uri))
    if err != nil {
        util.Loggrs.Warn("报错：", err.Error())
        return showErr{}
    }
    if strings.Contains(string(respones.Body()), "404 page not found") {
        return showErr{}
    }
    var item util.HexData
    err = json.Unmarshal(respones.Body(), &item)
    if err != nil {
        util.Loggrs.Warn(err.Error())
        return showErr{}
    }

    queryId := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%s/getstmtid?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
        util.H.Ip, util.Config.GetString("server.port"), app, queryid, queryid)

    if item.Data.StartTime > 0 {
        return showErr{
            Starttime: time.Unix(int64(item.Data.StartTime), 0).Format("2006-01-02 15:04:05"),
            User:      item.Data.User,
            Queryid:   queryId,
            Errmsg:    item.Data.FailedReason,
            Errinfo:   "",
            Stmt:      item.Data.SQL,
        }
    }
    return showErr{}
}

// 按时间降序排序
func sortByStartTimeDesc(errs []showErr) []showErr {
    // 创建副本以避免修改原始切片
    sorted := make([]showErr, len(errs))
    copy(sorted, errs)

    // 自定义排序函数
    sort.Slice(sorted, func(i, j int) bool {
        ti, err1 := time.Parse("2006-01-02 15:04:05", sorted[i].Starttime)
        tj, err2 := time.Parse("2006-01-02 15:04:05", sorted[j].Starttime)

        // 如果解析错误，保持原顺序
        if err1 != nil || err2 != nil {
            return false
        }
        // 降序排序（最新时间排前面）
        return ti.After(tj)
    })
    return sorted
}
