/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_pie
 *@date    2025/7/21 15:43
 */

package app

import (
	"StarRocksQuerys/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"sort"
)

// 定义数据项结构体
type chartitem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// 定义数据结构体
type chartData struct {
	Title      string      `json:"title"`
	Subtext    string      `json:"subtext"`
	Categories []string    `json:"categories"`
	Data       []chartitem `json:"data"`
}

func (engine *threadMap) pie(c *gin.Context) {
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
	result := showPie(appid, db, beginTime, endTime)
	c.JSON(http.StatusOK, result)
}

func showPie(app string, db *gorm.DB, begin, end string) chartData {
	if util.Read.Schema.Auditops == "" {
		return chartData{}
	}
	os.Setenv("TZ", "Asia/Shanghai")

	var m []map[string]interface{}
	stmt := fmt.Sprintf("select count(*) as count,user from %s where timestamp >= '%s' and timestamp < '%s' group by user order by count desc", util.Read.Schema.Auditops, begin, end)
	util.Loggrs.Info(stmt)
	r := db.Raw(stmt).Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return chartData{}
	}
	util.Loggrs.Info("pie:", len(m))
	if len(m) == 0 {
		return chartData{
			Title:      fmt.Sprintf("StarRocks 请求指标(%s~%s)", begin, end),
			Subtext:    fmt.Sprintf("数据来源: StarRocks连接数记录(%d/%d)", 0, len(m)),
			Categories: nil,
			Data:       []chartitem{},
		}
	}
	_, statlist, err := parseData(m)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return chartData{}
	}
	ulist, slist := sortract(statlist)
	// 创建数据实例
	chart := chartData{
		Title:      fmt.Sprintf("StarRocks 请求指标(%s~%s)", begin, end),
		Subtext:    fmt.Sprintf("数据来源: StarRocks连接数记录(%d/%d)", len(ulist), len(m)),
		Categories: ulist,
		Data:       slist,
	}
	return chart
}

// 解析原始数据并返回两个结果
func parseData(m []map[string]interface{}) ([]string, []chartitem, error) {
	var users []string
	var stats []chartitem
	// 临时map用于去重和统计
	userMap := make(map[string]int)

	for _, item := range m {
		// 获取user字段
		userVal, ok := item["user"]
		if !ok {
			continue // 跳过没有user字段的项
		}
		// 类型断言确保user是string
		user, ok := userVal.(string)
		if !ok {
			continue // 跳过类型不正确的user
		}
		// 获取count字段
		countVal, ok := item["count"]
		if !ok {
			continue // 跳过没有count字段的项
		}
		// 累加统计
		userMap[user] += int(countVal.(int64))
	}
	// 生成结果
	for user, total := range userMap {
		users = append(users, user)
		stats = append(stats, chartitem{
			Name:  user,
			Value: total,
		})
	}
	return users, stats, nil
}

//// 降序排序并提取Name
//func sortract(items []chartitem) []string {
//	// 1. 按Value降序排序
//	sort.Slice(items, func(i, j int) bool {
//		return items[i].Value > items[j].Value // 降序
//	})
//	// 2. 提取Name到[]string
//	names := make([]string, 0, len(items))
//	for _, item := range items {
//		names = append(names, item.Name)
//	}
//	return names
//}

// 降序排序并截取（如果超过150条则取前130条）
// 返回：names []string, topItems []chartitem
func sortract(items []chartitem) ([]string, []chartitem) {
	// 1. 按Value降序排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})
	// 2. 如果超过150条，截取前130条
	var topItems []chartitem
	if len(items) > 150 {
		topItems = items[:130]
	} else {
		topItems = items // 不超过150条，全部保留
	}
	// 3. 提取Name到[]string
	names := make([]string, 0, len(topItems))
	for _, item := range topItems {
		names = append(names, item.Name)
	}

	return names, topItems
}
