/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_errbroker
 *@date    2025/6/19 14:16
 */

package app

import (
	"StarRocksQuerys/tools"
	"StarRocksQuerys/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type submitData struct {
	data                            []util.TaskRuns
	running, pending, merged, total int
}

func (engine *threadMap) showSubmit(c *gin.Context) {
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

	var item submitData
	if val, ok := cachelimit.Get(appid + "submit"); ok {
		item = val.(submitData)
	} else {
		item = showSubmit(appid, db, beginTime, endTime)
	}
	c.JSON(http.StatusOK, gin.H{"data": sortByStartTimeDesc3(item.data), "running": item.running, "pending": item.pending, "merged": item.merged, "total": item.running + item.pending + item.merged})
}

func showSubmit(app string, db *gorm.DB, beginTime, endTime string) submitData {
	var running, pending, merged []string

	var infos []util.TaskRuns
	var m []util.TaskRuns
	//stmt := fmt.Sprintf("select * from information_schema.task_runs where state='FAILED' and CREATE_TIME >='%s' and CREATE_TIME<'%s'", beginTime, endTime)
	stmt := fmt.Sprintf("select * from information_schema.task_runs where state!='FAILED' and state !='SUCCESS' and CREATE_TIME >='%s' and CREATE_TIME<'%s'", beginTime, endTime)
	r := db.Raw(stmt).Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return submitData{}
	}
	if m == nil {
		return submitData{}
	}
	for _, m2 := range m {
		switch m2.STATE {
		case "RUNNING":
			running = append(running, m2.QUERY_ID)
		case "PENDING":
			pending = append(pending, m2.QUERY_ID)
		case "MERGED":
			merged = append(merged, m2.QUERY_ID)
		}
		queryId := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%d/getsubmit?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
			util.H.Ip, util.Read.Server.Port, app, m2.QUERY_ID, m2.QUERY_ID)

		layout := "2006-01-02T15:04:05+08:00"
		t1, _ := time.Parse(layout, m2.CREATE_TIME)
		t2, _ := time.Parse(layout, time.Now().Format(layout))

		progressInt, _ := strconv.Atoi(strings.NewReplacer("%", "", " ", "").Replace(m2.PROGRESS))

		infos = append(infos,
			util.TaskRuns{
				QUERY_ID:      queryId,
				TASK_NAME:     m2.TASK_NAME,
				CREATE_TIME:   m2.CREATE_TIME,
				FINISH_TIME:   m2.FINISH_TIME,
				STATE:         m2.STATE,
				CATALOG:       m2.CATALOG,
				WAREHOUSE:     m2.WAREHOUSE,
				DATABASE:      m2.DATABASE,
				DEFINITION:    m2.DEFINITION,
				EXPIRE_TIME:   m2.EXPIRE_TIME,
				ERROR_CODE:    m2.ERROR_CODE,
				ERROR_MESSAGE: m2.ERROR_MESSAGE,
				PROGRESS:      m2.PROGRESS,
				EXTRA_MESSAGE: m2.EXTRA_MESSAGE,
				PROPERTIES:    m2.PROPERTIES,
				// 自定义函数
				GetHour:     tools.GetHour(int(t2.Sub(t1).Seconds())),
				ProgressInt: progressInt,
				TimeMsg:     fmt.Sprintf("%s(CREATE_TIME)\n%s(EXPIRE_TIME)\n%s(FINISH_TIME)", m2.CREATE_TIME, m2.EXPIRE_TIME, m2.FINISH_TIME),
				QueryIdName: m2.QUERY_ID,
				Command:     fmt.Sprintf(`<button id="cancel-task" class="btn btn-light btn-sm ms-2" data-task="%s">❌</button>`, m2.TASK_NAME),
			})
	}

	sdata := submitData{
		data:    infos,
		running: len(running),
		pending: len(pending),
		merged:  len(merged),
		total:   len(running) + len(pending) + len(merged),
	}

	if _, ok := cachelimit.Get(app + "submit"); !ok {
		cachelimit.Set(app+"submit", sdata, cache.DefaultExpiration)
	}

	return sdata
}

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
	//stmt := fmt.Sprintf("select * from information_schema.task_runs where state='FAILED' and CREATE_TIME >='%s' and CREATE_TIME<'%s'", beginTime, endTime)
	stmt := fmt.Sprintf("select * from information_schema.task_runs where state='FAILED' and CREATE_TIME >='%s' and CREATE_TIME<'%s'", beginTime, endTime)
	r := db.Raw(stmt).Scan(&m)
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
