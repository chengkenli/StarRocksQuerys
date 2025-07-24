/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_errbroker
 *@date    2025/6/19 14:16
 */

package app

import (
	"StarRocksQuerys/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
	"time"
)

func (engine *threadMap) showBrokerErr(c *gin.Context) {
	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}

	beginTime := c.GetHeader("Begintime")
	endTime := c.GetHeader("Endtime")

	appid = setdefault(appid)
	db, err := engine.getmapConnect(appid)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, sortByStartTimeDesc(showbrokercancel(appid, db, beginTime, endTime)))
}

func showbrokercancel(app string, db *gorm.DB, beginTime, endTime string) []showErr {
	var m []map[string]interface{}
	r := db.Raw("show proc '/dbs'").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return nil
	}
	var errs []showErr
	done := make(chan struct{}, 5)
	var wg sync.WaitGroup
	for _, m2 := range m {
		wg.Add(1)
		go func(m2 map[string]interface{}) {
			defer func() {
				<-done
				wg.Done()
			}()

			done <- struct{}{}
			database := strings.ReplaceAll(m2["DbName"].(string), "default_cluster:", "")
			var tbLoads []map[string]interface{}
			r := db.Raw(fmt.Sprintf("show load from %s where STATE = 'CANCELLED'", database)).Scan(&tbLoads)
			if r.Error != nil {
				util.Loggrs.Error(r.Error.Error())
				return
			}
			for _, load := range tbLoads {

				ok, err := isInrange(beginTime, endTime, load["CreateTime"].(string))
				if err != nil {
					util.Loggrs.Error(err.Error())
					continue
				}
				if !ok {
					continue
				}
				label := fmt.Sprintf("%s.%s", database, load["Label"].(string))
				queryId := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%d/getbrokid?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
					util.H.Ip, util.Read.Server.Port, app, label, label)

				errs = append(errs, showErr{
					Starttime: load["CreateTime"].(string),
					User:      load["JobId"].(string),
					Queryid:   queryId,
					Errmsg:    load["ErrorMsg"].(string),
					Errinfo:   "",
					Stmt:      load["TaskInfo"].(string),
					Category:  categoryModel(load["ErrorMsg"].(string)),
				})
			}
		}(m2)
	}
	wg.Wait()

	return errs
}

// 判断目标时间是否在开始时间和结束时间范围内（包含边界）
func isInrange(startTimeStr, endTimeStr, targetTimeStr string) (bool, error) {
	// 定义时间格式
	timeLayout := "2006-01-02 15:04:05"
	// 解析开始时间
	startTime, err := time.Parse(timeLayout, startTimeStr)
	if err != nil {
		return false, fmt.Errorf("解析开始时间失败: %v", err)
	}
	// 解析结束时间
	endTime, err := time.Parse(timeLayout, endTimeStr)
	if err != nil {
		return false, fmt.Errorf("解析结束时间失败: %v", err)
	}
	// 解析目标时间
	targetTime, err := time.Parse(timeLayout, targetTimeStr)
	if err != nil {
		return false, fmt.Errorf("解析目标时间失败: %v", err)
	}
	// 判断是否在范围内（包含边界）
	return (targetTime.Equal(startTime) || targetTime.Equal(endTime)) ||
		(targetTime.After(startTime) && targetTime.Before(endTime)), nil
}
