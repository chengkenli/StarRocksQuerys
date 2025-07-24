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
	"reflect"
	"time"
)

func (engine *threadMap) showSubmitErr(c *gin.Context) {
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
	c.JSON(http.StatusOK, sortByStartTimeDesc(showSubmitcancel(appid, db, beginTime, endTime)))
}

func showSubmitcancel(app string, db *gorm.DB, beginTime, endTime string) []showErr {
	var errs []showErr

	var m []map[string]interface{}
	r := db.Raw(fmt.Sprintf("select * from information_schema.task_runs where state='FAILED' and CREATE_TIME >='%s' and CREATE_TIME<'%s'", beginTime, endTime)).Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return nil
	}
	if m == nil {
		util.Loggrs.Info("submit task failed job is nil.")
		return nil
	}
	for _, m2 := range m {

		queryId := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%d/getsubmit?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
			util.H.Ip, util.Read.Server.Port, app, m2["QUERY_ID"].(string), m2["QUERY_ID"].(string))

		var errmsg string
		if val, ok := m2["ERROR_MESSAGE"]; ok {
			if val != nil && reflect.TypeOf(val).Kind() == reflect.String {
				errmsg = m2["ERROR_MESSAGE"].(string)
			}
		}

		errs = append(errs,
			showErr{
				Starttime: m2["CREATE_TIME"].(time.Time).Format("2006-01-02 15:04:05"),
				User:      m2["TASK_NAME"].(string),
				Queryid:   queryId,
				Errmsg:    errmsg,
				Errinfo:   "",
				Stmt:      m2["DEFINITION"].(string),
				Category:  categoryModel(errmsg),
			})
	}
	return errs
}
