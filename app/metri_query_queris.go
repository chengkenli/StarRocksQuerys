/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package app
 *@file    metri_queris
 *@date    2025/5/27 10:17
 */

package app

import (
	"StarRocksAPIs/tools"
	"StarRocksAPIs/util"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func metriQuery(c *gin.Context) {
	appid := c.GetHeader("AppID")
	// set
	if v, ok := querycache.Get(appid + "queries"); ok {
		c.JSON(http.StatusOK, v.([]util.QueryResource))
		return
	}

	db, err := getmapConnect(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	restys, err := getmapResty(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	data := uriqueries(restys, db, leader(db))
	// 3. 返回JSON数据
	// 载入缓存
	go func() {
		// set
		querycache.Set(appid+"queries", data, cache.DefaultExpiration)
	}()
	c.JSON(http.StatusOK, data)
}

func uriqueries(client *resty.Client, db *gorm.DB, fe string) []util.QueryResource {
	var resource []util.QueryResource
	uri := fmt.Sprintf(`http://%s:8030/system?path=//current_queries`, fe)
	//创建Resty客户端
	//发送POST请求并处理响应
	respones, err := client.R().Get(uri)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	menu, _ := htmlquery.Parse(strings.NewReader(string(respones.Body())))
	table := htmlquery.Find(menu, `//*[@id="table_id"]/tbody/tr`)
	for _, node := range table {
		tr := htmlquery.Find(node, "td")
		if len(tr) >= 11 {
			//var wh string
			//if len(tr) == 12 {
			//	wh = td(tr[11])
			//}
			if tools.Version(db) >= 3.3 {
				id, _ := strconv.Atoi(split(td(tr[3])))
				if id == 0 {
					continue
				}
				scanBytes := int64(tools.Size(td(tr[6])))
				memoryUsage := int64(tools.Size(td(tr[8])))
				scanRows, _ := strconv.Atoi(split(td(tr[7])))
				cpuTime, _ := strconv.ParseFloat(split(td(tr[10])), 64)
				excTime, _ := strconv.ParseFloat(split(td(tr[11])), 64)
				if scanRows == 0 && scanBytes == 0 && memoryUsage == 0 && cpuTime == 0 {
					continue
				}
				resource = append(resource,
					util.QueryResource{
						ID:          int64(id),
						User:        td(tr[5]),
						ScanBytes:   scanBytes,
						ScanRows:    int64(scanRows),
						MemoryUsage: memoryUsage,
						CPUTime:     cpuTime,
						ExecTime:    excTime,
					})
			} else {
				id, _ := strconv.Atoi(split(td(tr[2])))
				if id == 0 {
					continue
				}
				scanBytes := int64(tools.Size(td(tr[5])))
				memoryUsage := int64(tools.Size(td(tr[7])))
				scanRows, _ := strconv.Atoi(split(td(tr[6])))
				cpuTime, _ := strconv.ParseFloat(split(td(tr[9])), 64)
				excTime, _ := strconv.ParseFloat(split(td(tr[10])), 64)
				if scanRows == 0 && scanBytes == 0 && memoryUsage == 0 && cpuTime == 0 {
					continue
				}
				resource = append(resource,
					util.QueryResource{
						ID:          int64(id),
						User:        td(tr[5]),
						ScanBytes:   scanBytes,
						ScanRows:    int64(scanRows),
						MemoryUsage: memoryUsage,
						CPUTime:     cpuTime,
						ExecTime:    excTime,
					})
			}
		}
	}
	return resource
}

func td(n *html.Node) string {
	v := htmlquery.InnerText(n)
	return v
}

func split(str string) string {
	return strings.Split(str, " ")[0]
}
