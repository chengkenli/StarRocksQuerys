/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_get_errbead
 *@date    2025/6/20 13:34
 */

package app

import (
	"StarRocksQuerys/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
)

func (engine *threadMap) metriErrbead(c *gin.Context) {
	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}

	appid = setdefault(appid)

	type data struct {
		Errquery  int64 `json:"errquery"`
		Errbroker int64 `json:"errbroker"`
		Errsubmit int64 `json:"errsubmit"`
	}

	if val, ok := escache.Get(appid + "errors"); ok {
		c.JSON(http.StatusOK, val.(data))
		return
	}

	beginTime := c.GetHeader("Begintime")
	endTime := c.GetHeader("Endtime")

	db, err := engine.getmapConnect(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	d := data{
		Errquery:  queryerrs(appid, db, beginTime, endTime),
		Errbroker: brokererrs(appid, db, beginTime, endTime),
		Errsubmit: submiterrs(appid, db, beginTime, endTime),
	}
	go escache.Set(appid+"errors", d, cache.DefaultExpiration)

	c.JSON(http.StatusOK, d)
}

// 统计query erros
func queryerrs(app string, db *gorm.DB, begin, end string) int64 {
	var queryes map[string]interface{}

	var stmt string
	if len(begin) != 0 && len(end) != 0 {
		stmt = fmt.Sprintf("select count(*) as count from %s where timestamp >= '%s' and timestamp < '%s' and user != 'root' and state='ERR'", util.Read.Schema.Auditops, begin, end)
	} else {
		stmt = fmt.Sprintf("select count(*) as count from %s where timestamp >= date_sub(now(), interval 30 minute) and user != 'root' and state='ERR'", util.Read.Schema.Auditops)
	}

	r := db.Raw(stmt).Scan(&queryes)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return -1
	}
	return queryes["count"].(int64)
}

// 统计broker errors
func brokererrs(app string, db *gorm.DB, begin, end string) int64 {
	var brokeres []map[string]interface{}
	if val, ok := apicache.Get(app + "dbs"); ok {
		brokeres = val.([]map[string]interface{})
	} else {
		r := db.Raw("show proc '/dbs'").Scan(&brokeres)
		if r.Error != nil {
			util.Loggrs.Error(r.Error.Error())
			return -1
		}
	}

	var err_broker []map[string]interface{}
	done := make(chan struct{}, 5)
	var wg sync.WaitGroup
	for _, m2 := range brokeres {
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

				ok, err := isInrange(begin, end, load["CreateTime"].(string))
				if err != nil {
					util.Loggrs.Error(err.Error())
					continue
				}
				if !ok {
					continue
				}

				err_broker = append(err_broker, load)
			}
		}(m2)
	}
	wg.Wait()

	go apicache.Set(app+"dbs", brokeres, cache.DefaultExpiration)
	return int64(len(err_broker))
}

// 统计submit task errors
func submiterrs(app string, db *gorm.DB, begin, end string) int64 {
	var submites []map[string]interface{}
	r := db.Raw(fmt.Sprintf("select * from information_schema.task_runs where state='FAILED' and CREATE_TIME >='%s' and CREATE_TIME<'%s'", begin, end)).Scan(&submites)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return -1
	}
	return int64(len(submites))
}
