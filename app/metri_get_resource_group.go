/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_get_resource_group
 *@date    2025/9/1 14:32
 */

package app

import (
	"StarRocksQuerys/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
)

func (engine *threadMap) resourceGroup(c *gin.Context) {
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

	var itemData []util.ResourceMsg
	if val, ok := shortcache.Get(appid + "resource"); ok {
		itemData = val.([]util.ResourceMsg)
	} else {
		itemData = resource(appid, db)
	}
	var view []util.ResourceMsg
	if len(itemData) > 20 {
		view = itemData[0:20]
	} else {
		view = itemData
	}
	c.JSON(http.StatusOK, gin.H{"data": view, "total": len(itemData)})
}

func resource(app string, db *gorm.DB) []util.ResourceMsg {
	var m []map[string]interface{}
	r := db.Raw("SHOW RESOURCE GROUPS ALL").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return nil
	}

	var srcData []util.ResourceMsg
	for _, item := range m {

		matches := regexp.MustCompile(`user=([^)]+)`).FindStringSubmatch(item["classifiers"].(string))
		var user string
		if matches == nil {
			continue
		}
		if len(matches) > 1 {
			user = matches[1]
		}

		srcData = append(srcData, util.ResourceMsg{
			Name:              item["name"].(string),
			User:              strings.Split(user, ",")[0],
			CpuWeight:         item["cpu_weight"],
			ExclusiveCpuCores: item["exclusive_cpu_cores"],
			MemLimit:          item["mem_limit"],
			ConcurrencyLimit:  item["concurrency_limit"],
			BigQueryCPU:       item["big_query_cpu_second_limit"],
			BigQueryRows:      item["big_query_scan_rows_limit"],
			BigQueryMemLimit:  item["big_query_mem_limit"],
			Command:           fmt.Sprintf(`<button id="drop-resourcegroup" class="btn btn-light btn-sm ms-2" data-resourcegroup="%s">❌</button>`, item["name"].(string)),
		})
	}

	if _, ok := shortcache.Get(app + "resource"); !ok {
		shortcache.Set(app+"resource", srcData, cache.DefaultExpiration)
	}

	return srcData
}
