/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_queris
 *@date    2025/5/27 10:17
 */

package app

import (
	"StarRocksQuerys/tools"
	"StarRocksQuerys/util"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (engine *threadMap) metriQuery(c *gin.Context) {
	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}
	// set
	if v, ok := querycache.Get(appid + "queries"); ok {
		c.JSON(http.StatusOK, v.([]util.GlobalQueries))
		return
	}
	appid = setdefault(appid)
	db, err := engine.getmapConnect(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	var data []util.GlobalQueries
	if tools.Version(db) >= 3.3 {
		data = globalqueries(db)
	} else {
		restys, err := engine.getmapResty(appid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		data = uriqueries(restys, db, leader(db))
	}
	// 3. 返回JSON数据
	// 载入缓存
	go func() {
		// set
		querycache.Set(appid+"queries", data, cache.DefaultExpiration)
	}()
	c.JSON(http.StatusOK, data)
}

func globalqueries(db *gorm.DB) []util.GlobalQueries {
	var m []util.GlobalQueries
	r := db.Raw("show proc '/global_current_queries'").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return m
	}
	return sortByScanRowsDesc(m)
}

func uriqueries(client *resty.Client, db *gorm.DB, fe string) []util.GlobalQueries {

	var resource []util.GlobalQueries
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
					util.GlobalQueries{
						StartTime:     "",
						QueryId:       "",
						ConnectionId:  int64(id),
						Database:      "",
						User:          td(tr[5]),
						ScanBytes:     fmt.Sprintf("%d", scanBytes),
						ScanRows:      fmt.Sprintf("%d", scanRows),
						MemoryUsage:   fmt.Sprintf("%d", memoryUsage),
						DiskSpillSize: "",
						CPUTime:       fmt.Sprintf("%0.1f", cpuTime),
						ExecTime:      fmt.Sprintf("%0.1f", excTime),
						Warehouse:     "",
						CustomQueryId: "",
						ResourceGroup: "",
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
					util.GlobalQueries{
						StartTime:     "",
						QueryId:       "",
						ConnectionId:  int64(id),
						Database:      "",
						User:          td(tr[5]),
						ScanBytes:     fmt.Sprintf("%d", scanBytes),
						ScanRows:      fmt.Sprintf("%d", scanRows),
						MemoryUsage:   fmt.Sprintf("%d", memoryUsage),
						DiskSpillSize: "",
						CPUTime:       fmt.Sprintf("%0.1f", cpuTime),
						ExecTime:      fmt.Sprintf("%0.1f", excTime),
						Warehouse:     "",
						CustomQueryId: "",
						ResourceGroup: "",
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

func sortByScanRowsDesc(queries []util.GlobalQueries) []util.GlobalQueries {
	// 复制切片以避免修改原数据
	sorted := make([]util.GlobalQueries, len(queries))
	copy(sorted, queries)

	// 按 ScanRows 降序排序
	sort.Slice(sorted, func(i, j int) bool {
		inum := strings.Split(sorted[i].ScanRows, " ")[0]
		jnum := strings.Split(sorted[j].ScanRows, " ")[0]
		// 将 ScanRows 从字符串转为 int64 进行比较
		valI, _ := strconv.ParseInt(inum, 10, 64)
		valJ, _ := strconv.ParseInt(jnum, 10, 64)
		return valI > valJ // 降序
	})

	return sorted
}
