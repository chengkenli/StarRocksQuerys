/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_err_stream
 *@date    2025/7/30 9:58
 */

package app

import (
	"StarRocksQuerys/tools"
	"StarRocksQuerys/util"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

func (engine *threadMap) showStreamErr(c *gin.Context) {
	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}

	//beginTime := c.GetHeader("Begintime")
	//endTime := c.GetHeader("Endtime")

	appid = setdefault(appid)
	db, err := engine.getmapConnect(appid)
	if err != nil {
		util.Loggrs.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if tools.Version(db) < 3.3 {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, sortByStartTimeDesc(getStreamErrs(appid, leader(db))))
}

func getStreamErrs(appid string, leader string) []showErr {
	restys, err := engine.getmapResty(appid)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil
	}
	return restyStream(restys, appid, leader)
}

func restyStream(restys *resty.Client, appid, leader string) []showErr {
	var errs []showErr

	uri := fmt.Sprintf(`http://%s:8030/system?path=//stream_loads`, leader)
	//发送POST请求并处理响应
	respones, err := restys.R().Get(uri)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return nil
	}
	menu, _ := htmlquery.Parse(strings.NewReader(string(respones.Body())))
	table := htmlquery.Find(menu, `//*[@id="table_id"]/tbody/tr`)
	for _, node := range table {
		tr := htmlquery.Find(node, "td")
		if td(tr[4]) != "CANCELLED" {

			continue
		}
		//从queryid中获取作业信息
		uri := fmt.Sprintf(`http://%s:8030/system?path=//stream_loads/%s`, leader, td(tr[0]))
		respones, err := restys.R().Get(uri)
		if err != nil {
			util.Loggrs.Error(err.Error())
			continue
		}
		menu, _ := htmlquery.Parse(strings.NewReader(string(respones.Body())))
		table := htmlquery.Find(menu, `//*[@id="table_id"]/tbody/tr`)
		if table == nil {
			continue
		}
		for _, node := range table {
			tr := htmlquery.Find(node, "td")
			if tr == nil {
				continue
			}
			itemData := util.StreamData{
				Label:                 td(tr[0]),
				Id:                    td(tr[1]),
				LoadId:                td(tr[2]),
				TxnId:                 td(tr[3]),
				DbName:                td(tr[4]),
				TableName:             td(tr[5]),
				State:                 td(tr[6]),
				ErrorMsg:              td(tr[7]),
				TrackingURL:           td(tr[8]),
				ChannelNum:            td(tr[9]),
				PreparedChannelNum:    td(tr[10]),
				NumRowsNormal:         td(tr[11]),
				NumRowsAbNormal:       td(tr[12]),
				NumRowsUnselected:     td(tr[13]),
				NumLoadBytes:          td(tr[14]),
				TimeoutSecond:         td(tr[15]),
				CreateTimeMs:          td(tr[16]),
				BeforeLoadTimeMs:      td(tr[17]),
				StartLoadingTimeMs:    td(tr[18]),
				StartPreparingTimeMs:  td(tr[19]),
				FinishPreparingTimeMs: td(tr[20]),
				EndTimeMs:             td(tr[21]),
				ChannelState:          td(tr[22]),
				Type:                  td(tr[23]),
				TrackingSQL:           td(tr[24]),
			}
			//marshal, err := json.Marshal(itemData)
			//if err != nil {
			//	util.Loggrs.Error(err.Error())
			//}
			label := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%d/getstream?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
				util.H.Ip, util.Read.Server.Port, appid, itemData.Label, itemData.Label)
			errs = append(errs,
				showErr{
					Starttime: itemData.CreateTimeMs,
					User:      itemData.DbName,
					Queryid:   label,
					Errmsg:    itemData.ErrorMsg,
					Errinfo:   itemData.TrackingSQL,
					Stmt:      itemData.TxnId,
					Category:  categoryModel(itemData.ErrorMsg),
				})
		}
	}

	return errs
}
