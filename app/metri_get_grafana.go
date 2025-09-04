/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_get_grafana
 *@date    2025/7/28 11:24
 */

package app

import (
	"StarRocksQuerys/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
)

type itemData struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func (engine *threadMap) metrigrafana(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "不对外开放，可自己扩展."})
	return

	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}
	appid = setdefault(appid)
	var aliasname string
	for _, m := range util.MetaLink {
		if m["app"].(string) == appid {
			aliasname = m["alias"].(string)
		}
	}
	cpu := fmt.Sprintf(`https://ss.com/api/v1/query?query=topk%%281%%2Csort_desc%%28clamp_min%%281%%20-%%20avg%%28irate%%28node_cpu_seconds_total%%7Bcluster_id%%3D%%22%s%%22%%2Cmodule%%3D%%22be%%2Fbroker%%22%%2Cmode%%3D%%22idle%%22%%7D%%5B30s%%5D%%29%%29%%20by%%20%%28instance%%29%%2C0%%29%%29%%29`, aliasname)
	mem := fmt.Sprintf(`https://ss.com/api/v1/query?query=topk%%281%%2Csort_desc%%281%%20-%%20%%28node_memory_MemAvailable_bytes%%7Bcluster_id%%3D%%22%s%%22%%2Cmodule%%3D%%22be%%2Fbroker%%22%%7D%%20%%2F%%20%%28node_memory_MemTotal_bytes%%7Bcluster_id%%3D%%22%s%%22%%2Cmodule%%3D%%22be%%2Fbroker%%22%%7D%%29%%29%%29%%29`, aliasname, aliasname)
	ios := fmt.Sprintf(`https://ss.com/api/v1/query?query=topk%%281%%2C%%20sort_desc%%28avg%%28irate%%28node_disk_io_time_seconds_total%%7Bcluster_id%%3D%%22%s%%22%%2C%%20module%%3D%%22be%%2Fbroker%%22%%7D%%5B1m%%5D%%29%%29%%20by%%20%%28instance%%29%%20%%2A%%20100%%29%%29`, aliasname)
	dik := fmt.Sprintf(`https://ss.com/api/v1/query?query=topk%%281%%2Csort_desc%%28100%%20-%%20node_filesystem_avail_bytes%%7Bcluster_id%%3D%%22%s%%22%%2Cmodule%%3D%%22be%%2Fbroker%%22%%2Cfstype%%21%%3D%%22tmpfs%%22%%2Cmountpoint%%20%%21%%7E%%20%%27.%%2A%%2Fstarrocks_dc.%%2A%%27%%7D%%2Fnode_filesystem_size_bytes%%7Bcluster_id%%3D%%22%s%%22%%2Cmodule%%3D%%22be%%2Fbroker%%22%%2Cfstype%%21%%3D%%22tmpfs%%22%%2Cmountpoint%%20%%21%%7E%%20%%27.%%2A%%2Fstarrocks_dc.%%2A%%27%%7D%%20%%2A%%20100%%29%%29`, aliasname, aliasname)

	c.JSON(http.StatusOK, gin.H{"cpu": restys(cpu, false), "memory": restys(mem, false), "io": restys(ios, true), "disk": restys(dik, true)})
}

func restys(uri string, transform bool) string {
	response, err := resty.New().R().SetHeader("Accept", "application/json").Get(uri)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return "0"
	}
	var item itemData
	err = json.Unmarshal(response.Body(), &item)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return "0"
	}
	if item.Data.Result == nil {
		util.Loggrs.Warn("解析第一层数据失败")
		util.Loggrs.Warn(item)
		return "0"
	}
	if len(item.Data.Result) == 0 {
		util.Loggrs.Warn("解析第二层数据失败")
		util.Loggrs.Warn(item)
		return "0"
	}
	if len(item.Data.Result[0].Value) < 2 {
		util.Loggrs.Warn("解析第三层数据失败")
		util.Loggrs.Warn(item)
		return "0"
	}
	if transform {
		float, _ := strconv.ParseFloat(item.Data.Result[0].Value[1].(string), 64)
		return fmt.Sprintf("%0.2f", float/100)
	}
	return item.Data.Result[0].Value[1].(string)
}
