/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_stmtid
 *@date    2025/5/27 15:55
 */

package app

import (
	"StarRocksQuerys/tools"
	"StarRocksQuerys/util"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

func (engine *threadMap) metriStreamId(c *gin.Context) {
	jobid, _ := c.GetQuery("id")
	appid, _ := c.GetQuery("app")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}
	db, err := engine.getmapConnect(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	restys, err := engine.getmapResty(appid)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}

	leader := leader(db)
	filename := fmt.Sprintf("%s/%s.%s.submit.log", util.Read.Log.Path, appid, jobid)
	//从queryid中获取作业信息
	uri := fmt.Sprintf(`http://%s:8030/system?path=//stream_loads/%s`, leader, jobid)
	respones, err := restys.R().Get(uri)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
	menu, _ := htmlquery.Parse(strings.NewReader(string(respones.Body())))
	table := htmlquery.Find(menu, `//*[@id="table_id"]/tbody/tr`)
	if table == nil {
		util.Loggrs.Warn("result is nil.")
		return
	}

	var wg sync.WaitGroup
	for _, node := range table {
		wg.Add(1)
		node := node
		go func() {
			defer wg.Done()

			tr := htmlquery.Find(node, "td")
			if tr == nil {
				return
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
			if jobid == itemData.Label {
				marshal, err := json.MarshalIndent(itemData, "", "  ")
				if err != nil {
					util.Loggrs.Error(err.Error())
				}
				tools.WriteFile(filename, string(marshal))
				util.Loggrs.Info("内容写入")
			}
		}()
	}
	wg.Wait()
	url := fmt.Sprintf("http://%s:%d/log/%s", util.H.Ip, util.Read.Server.Port, filename)
	c.JSON(http.StatusOK, url)
}
