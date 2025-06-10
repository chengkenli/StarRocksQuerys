/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package app
 *@file    metri_stmtid
 *@date    2025/5/27 15:55
 */

package app

import (
	"StarRocksAPIs/conn"
	"StarRocksAPIs/tools"
	"StarRocksAPIs/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func metriStmtId(c *gin.Context) {
	stmtid, _ := c.GetQuery("id")
	appid, _ := c.GetQuery("app")

	db, err := getmapConnect(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	var result string
	if !regexp.MustCompile(`[a-zA-Z]`).MatchString(stmtid) {
		result = getConnectionId(appid, stmtid, db)
	} else {
		result = getQueryId(appid, stmtid, db)
	}
	c.JSON(http.StatusOK, result)
}

func getQueryId(appid, stmtid string, db *gorm.DB) string {
	if util.Config.GetString("schema.auditops") == "" {
		return ""
	}

	var m map[string]interface{}
	r := db.Raw(fmt.Sprintf("select * from %s where queryId='%s'", util.Config.GetString("schema.auditops"), stmtid)).Scan(&m)
	if r.Error != nil {
		return ""
	}

	if m == nil {
		return ""
	}
	var pendingTimeMs int64
	v, ok := m["pendingTimeMs"]
	if ok {
		pendingTimeMs = v.(int64)
	}

	msg := fmt.Sprintf(`
💬QueryId          :         %s
💬Timestamp        :         %s
💬QueryType        :         %s
💬ClientIp         :         %s
💬User             :         %s
💬AuthorizedUser   :         %s
💬ResourceGroup    :         %s
💬Catalog          :         %s
💬Db               :         %s
💬State            :         %s
💬ErrorCode        :         %s
💬QueryTime        :         %d
💬ScanBytes        :         %d
💬ScanRows         :         %d
💬ReturnRows       :         %d
💬CpuCostNs        :         %d
💬MemCostBytes     :         %d
💬StmtId           :         %d
💬IsQuery          :         %b
💬FeIp             :         %s
💬Digest           :         %s
💬PlanCpuCosts     :         %f
💬PlanMemCosts     :         %f
💬PendingTimeMs    :         %d
💬Stmt             :         
%s
`,
		m["queryId"],
		m["timestamp"],
		m["queryType"],
		m["clientIp"],
		m["user"],
		m["authorizedUser"],
		m["resourceGroup"],
		m["catalog"],
		m["db"],
		m["state"],
		m["errorCode"],
		m["queryTime"].(int64),
		m["scanBytes"].(int64),
		m["scanRows"].(int64),
		m["returnRows"].(int64),
		m["cpuCostNs"].(int64),
		m["memCostBytes"].(int64),
		m["stmtId"].(int64),
		m["isQuery"].(int64),
		m["feIp"],
		m["digest"],
		m["planCpuCosts"].(float64),
		m["planMemCosts"].(float64),
		pendingTimeMs,
		m["stmt"],
	)

	filename := fmt.Sprintf("%s/%s.%s.log", util.Config.GetString("log.path"), appid, stmtid)
	tools.WriteFile(filename, msg)
	session := "http://10.88.10.82:7890/log/" + filename
	return session
}

func getConnectionId(appid, stmtid string, db *gorm.DB) string {

	id, _ := strconv.Atoi(stmtid)

	var felist []string
	var m []map[string]interface{}
	r := db.Raw("show frontends").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Warn(r.Error.Error())
		return ""
	}
	for _, item := range m {
		if item["Alive"].(string) != "true" {
			continue
		}
		felist = append(felist, item["IP"].(string))
	}

	var session string
	var global sync.WaitGroup
	for _, ip := range felist {
		global.Add(1)
		ip := ip
		go func() {
			defer global.Done()

			single, err := conn.StarRocksSingle(appid, ip)
			if err != nil {
				util.Loggrs.Warn(err.Error())
				return
			}

			var dbresult []util.Process
			r = single.Raw("show full processlist").Scan(&dbresult)
			if r.Error != nil {
				util.Loggrs.Warn(r.Error.Error())
				return
			}
			defer func() {
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
			}()

			done := make(chan struct{}, 10)
			var wg sync.WaitGroup
			for _, m := range dbresult {
				wg.Add(1)
				m := m
				go func() {
					defer func() {
						<-done
						wg.Done()
					}()
					done <- struct{}{}

					connectionId, _ := strconv.Atoi(m.Id)
					if connectionId == id {
						filename := fmt.Sprintf("%s/%s.%s.log", util.Config.GetString("log.path"), appid, m.Id)
						tools.WriteFile(filename, fmt.Sprintf(`
💬App:            %s
💬Fe:             %s
💬ClientIP:       %s
💬Type:           %s
💬ConnectionId:   %s
💬Database:       %s
💬User:           %s
💬ExecTime:       %s
💬Stmt:
%s`, appid, ip, m.Host, m.State, m.Id, m.Db, m.User, m.Time, m.Info))
						session = "http://10.88.10.82:7890/log/" + filename
						return
					}
				}()
			}
			wg.Wait()
		}()
	}
	global.Wait()
	return session
}
