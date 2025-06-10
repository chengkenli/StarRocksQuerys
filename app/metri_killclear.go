/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package app
 *@file    metri_killc
 *@date    2025/5/26 16:34
 */

package app

import (
	"StarRocksAPIs/conn"
	"StarRocksAPIs/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

func metriKillc(c *gin.Context) {
	appid := c.GetHeader("AppID")
	db, err := getmapConnect(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	var felist []string
	var m []map[string]interface{}
	r := db.Raw("show frontends").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Warn(appid, ">", r.Error.Error())
		return
	}
	for _, item := range m {
		if item["Alive"].(string) != "true" {
			continue
		}
		felist = append(felist, item["IP"].(string))
	}

	var global sync.WaitGroup
	var fail, success []string
	for _, ip := range felist {
		global.Add(1)
		ip := ip
		go func() {
			defer global.Done()

			single, err := conn.StarRocksSingle(appid, ip)
			if err != nil {
				util.Loggrs.Warn(appid, ">", err.Error())
				return
			}

			var dbresult []map[string]interface{}
			r = single.Raw("show processlist").Scan(&dbresult)
			if r.Error != nil {
				util.Loggrs.Warn(appid, ">", r.Error.Error())
				return
			}

			done := make(chan struct{}, 10)
			var wg sync.WaitGroup
			for _, item := range dbresult {
				wg.Add(1)
				item := item
				go func() {
					defer func() {
						<-done
						wg.Done()
					}()
					done <- struct{}{}

					if item["Command"].(string) == "Sleep" {
						return
					}
					edtime := item["Time"].(int64)
					if edtime < 1500 {
						return
					}
					r := single.Exec(fmt.Sprintf("kill %d", item["Id"].(int64)))
					if r.Error != nil {
						util.Loggrs.Warn(appid, ">", r.Error.Error())
						fail = append(fail, fmt.Sprintf("%d(%s)", item["Id"].(int64), item["User"].(string)))
						return
					}
					util.Loggrs.Info(appid, ">", "处理请求成功，kill sleep:", item["Id"].(int64))
					success = append(success, fmt.Sprintf("%d(%s)", item["Id"].(int64), item["User"].(string)))
				}()
			}
			wg.Wait()

			/*每次使用完，主动关闭连接数*/
			sqlDB, err := single.DB()
			if err != nil {
				util.Loggrs.Error(appid, ">", err.Error())
				return
			}
			sqlDB.SetMaxOpenConns(10)                 //最大连接数
			sqlDB.SetMaxIdleConns(3)                  //最大空闲连接数
			sqlDB.SetConnMaxLifetime(5 * time.Second) //空闲连接最多存活时间
			sqlDB.Close()
		}()
	}
	global.Wait()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("success:%d,failed:%d", len(success), len(fail)),
	})
}
