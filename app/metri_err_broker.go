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
)

func (engine *threadMap) showBrokerErr(c *gin.Context) {
	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}
	appid = setdefault(appid)
	db, err := engine.getmapConnect(appid)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, sortByStartTimeDesc(showbrokercancel(appid, db)))
}

func showbrokercancel(app string, db *gorm.DB) []showErr {
	var m []map[string]interface{}
	r := db.Raw("show proc '/dbs'").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return nil
	}
	var errs []showErr
	var wg sync.WaitGroup
	for _, m2 := range m {
		wg.Add(1)
		go func(m2 map[string]interface{}) {
			defer wg.Done()
			database := strings.ReplaceAll(m2["DbName"].(string), "default_cluster:", "")
			var tbLoads []map[string]interface{}
			r := db.Raw(fmt.Sprintf("show load from %s where STATE = 'CANCELLED'", database)).Scan(&tbLoads)
			if r.Error != nil {
				util.Loggrs.Error(r.Error.Error())
				return
			}
			for _, load := range tbLoads {

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
				})
			}
		}(m2)
	}
	wg.Wait()

	return errs
}
