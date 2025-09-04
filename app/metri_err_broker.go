/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_errbroker
 *@date    2025/6/19 14:16
 */

package app

import (
	"StarRocksQuerys/tools"
	"StarRocksQuerys/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type brokerData struct {
	data                                        []util.BrokerMsg
	loading, pending, queueing, prepared, total int
}

func (engine *threadMap) showBroker(c *gin.Context) {
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

	var item brokerData
	if val, ok := cachelimit.Get(appid + "broker"); ok {
		item = val.(brokerData)
	} else {
		item = showbroker(appid, db, beginTime, endTime)
	}
	c.JSON(http.StatusOK, gin.H{"data": sortByStartTimeDesc2(item.data), "loading": item.loading, "pending": item.pending, "queueing": item.queueing, "prepared": item.prepared, "total": item.total})
}

func showbroker(app string, db *gorm.DB, beginTime, endTime string) brokerData {
	var m []map[string]interface{}
	r := db.Raw("show proc '/dbs'").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return brokerData{}
	}
	var loading, pending, queueing, prepared []string
	var infos []util.BrokerMsg
	done := make(chan struct{}, 5)
	var wg sync.WaitGroup
	for _, m2 := range m {
		wg.Add(1)
		go func(m2 map[string]interface{}) {
			defer func() {
				<-done
				wg.Done()
			}()

			done <- struct{}{}
			database := strings.ReplaceAll(m2["DbName"].(string), "default_cluster:", "")
			var tbLoads, tmpmsg []util.BrokerMsg
			state := []string{"LOADING", "PENDING", "QUEUEING", "ETL"}
			for _, item := range state {

				r := db.Raw(fmt.Sprintf("show load from %s where STATE = '%s'", database, item)).Scan(&tmpmsg)
				if r.Error != nil {
					util.Loggrs.Error(r.Error.Error())
					continue
				}
				tbLoads = append(tbLoads, tmpmsg...)
			}
			for _, load := range tbLoads {
				ok, err := isInrange(beginTime, endTime, load.CreateTime)
				if err != nil {
					util.Loggrs.Error(err.Error())
					continue
				}
				if !ok {
					continue
				}
				switch load.State {
				case "LOADING":
					loading = append(loading, load.Label)
				case "PENDING":
					pending = append(pending, load.Label)
				case "QUEUEING":
					queueing = append(queueing, load.Label)
				case "ETL":
					prepared = append(prepared, load.Label)
				}
				label := fmt.Sprintf("%s.%s", database, load.Label)
				var labelname string
				if len(label) > 25 {
					labelname = label[0:25]
				} else {
					labelname = label
				}
				queryId := fmt.Sprintf(`<a href="javascript:void(0)" onclick="fetchAndRedirect('http://%s:%d/getbrokid?app=%s&id=%s')" class="plain-text-link" style="color: #333; text-decoration: none; font-weight: normal; cursor: pointer">%s</a>`,
					util.H.Ip, util.Read.Server.Port, app, label, labelname)

				var progress string
				matches := regexp.MustCompile(`LOAD:(\d+)%?`).FindStringSubmatch(load.Progress)
				if len(matches) >= 2 {
					progress = matches[1]
				}

				var jobdl util.BrokerJobDetails
				json.Unmarshal([]byte(load.JobDetails), &jobdl)
				scanInfo := fmt.Sprintf("%d(InternalTableLoadBytes)\n%d(ScanBytes)\n\n%d(InternalTableLoadRows)\n%d(ScanRows)",
					jobdl.InternalTableLoadBytes, jobdl.ScanBytes, jobdl.InternalTableLoadRows, jobdl.ScanRows)

				jobdetail := strings.NewReplacer(
					"{", "",
					"}", "",
					"[", "",
					"]", "",
					`"`, "'",
				).Replace(load.JobDetails)
				all, unfinished := detail(jobdetail)

				layout := "2006-01-02 15:04:05"
				t1, _ := time.Parse(layout, load.CreateTime)
				t2, _ := time.Parse(layout, time.Now().Format(layout))

				timeMsg := fmt.Sprintf("%s(CreateTime)\n%s(EtlStartTime)\n%s(EtlFinishTime)\n%s(LoadStartTime)\n%s(LoadFinishTime)", load.CreateTime,
					load.EtlStartTime,
					load.EtlFinishTime,
					load.LoadStartTime,
					load.LoadFinishTime)

				var backends string
				Progress, _ := strconv.Atoi(progress)
				if Progress >= 50 && len(unfinished) == 0 {
					backends = fmt.Sprintf(`all:(<strong>%d</strong>) unfinished:(<strong style="color: green;">%d</strong>)`, len(all), len(unfinished))
				} else {
					backends = fmt.Sprintf(`all:(<strong>%d</strong>) unfinished:(<strong style="color: red;">%d</strong>)`, len(all), len(unfinished))
				}

				infos = append(infos, util.BrokerMsg{
					JobId:          load.JobId,
					Label:          queryId,
					State:          load.State,
					Progress:       progress,
					Type:           load.Type,
					Priority:       load.Priority,
					ScanRows:       fmt.Sprintf("%s/%d", tools.FormatBytes(int64(jobdl.ScanBytes-jobdl.InternalTableLoadBytes)), jobdl.ScanRows-jobdl.InternalTableLoadRows),
					FilteredRows:   load.FilteredRows,
					UnselectedRows: load.UnselectedRows,
					SinkRows:       load.SinkRows,
					EtlInfo:        load.EtlInfo,
					TaskInfo:       load.TaskInfo,
					ErrorMsg:       load.ErrorMsg,
					CreateTime:     load.CreateTime,
					EtlStartTime:   load.EtlStartTime,
					EtlFinishTime:  load.EtlFinishTime,
					LoadStartTime:  load.LoadStartTime,
					LoadFinishTime: load.LoadFinishTime,
					TrackingSQL:    load.TrackingSQL,
					JobDetails: strings.NewReplacer(
						"{", "",
						"}", "",
						"[", "",
						"]", "",
						`"`, "'",
						",", ",\n",
					).Replace(load.JobDetails),
					//自定义
					Backends:    backends,
					ScanInfo:    scanInfo,
					Command:     fmt.Sprintf(`<button id="cancel-label" class="btn btn-light btn-sm ms-2" data-label="%s">❌</button>`, fmt.Sprintf("%s.%s", database, load.Label)),
					GetHour:     tools.GetHour(int(t2.Sub(t1).Seconds())),
					TimeMsg:     timeMsg,
					LabelName:   fmt.Sprintf("%s.%s", database, load.Label),
					ProgressMsg: load.Progress,
				})
			}
		}(m2)
	}
	wg.Wait()

	bdata := brokerData{
		data:     deduplicateByJobID(infos),
		loading:  len(tools.RmDuplicaSlice(loading)),
		pending:  len(tools.RmDuplicaSlice(pending)),
		queueing: len(tools.RmDuplicaSlice(queueing)),
		prepared: len(tools.RmDuplicaSlice(prepared)),
		total: len(tools.RmDuplicaSlice(loading)) +
			len(tools.RmDuplicaSlice(pending)) +
			len(tools.RmDuplicaSlice(queueing)) +
			len(tools.RmDuplicaSlice(prepared)),
	}

	if _, ok := cachelimit.Get(app + "broker"); !ok {
		cachelimit.Set(app+"broker", bdata, cache.DefaultExpiration)
	}

	return bdata
}

func detail(text string) ([]string, []string) {
	// 提取所有数字列表的模式
	re := regexp.MustCompile(`'(?:All|Unfinished) backends':[^:]+:(\d+(?:,\d+)*)`)
	matches := re.FindAllStringSubmatch(text, -1)

	var all, unfinished []string
	if len(matches) >= 1 {
		all = strings.Split(matches[0][1], ",")
	}
	if len(matches) >= 2 {
		unfinished = strings.Split(matches[1][1], ",")
	}
	return all, unfinished
}

// 结构体去重
func deduplicateByJobID(messages []util.BrokerMsg) []util.BrokerMsg {
	seen := make(map[string]bool)
	var result []util.BrokerMsg
	for _, msg := range messages {
		if _, exists := seen[msg.JobId]; !exists {
			seen[msg.JobId] = true
			result = append(result, msg)
		}
	}
	return result
}

func (engine *threadMap) showBrokerErr(c *gin.Context) {
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
	c.JSON(http.StatusOK, sortByStartTimeDesc(showbrokercancel(appid, db, beginTime, endTime)))
}

func showbrokercancel(app string, db *gorm.DB, beginTime, endTime string) []showErr {
	var m []map[string]interface{}
	r := db.Raw("show proc '/dbs'").Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return nil
	}
	var errs []showErr
	done := make(chan struct{}, 5)
	var wg sync.WaitGroup
	for _, m2 := range m {
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

				ok, err := isInrange(beginTime, endTime, load["CreateTime"].(string))
				if err != nil {
					util.Loggrs.Error(err.Error())
					continue
				}
				if !ok {
					continue
				}
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
					Category:  categoryModel(load["ErrorMsg"].(string)),
				})
			}
		}(m2)
	}
	wg.Wait()

	return errs
}

// 判断目标时间是否在开始时间和结束时间范围内（包含边界）
func isInrange(startTimeStr, endTimeStr, targetTimeStr string) (bool, error) {
	// 定义时间格式
	timeLayout := "2006-01-02 15:04:05"
	// 解析开始时间
	startTime, err := time.Parse(timeLayout, startTimeStr)
	if err != nil {
		return false, fmt.Errorf("解析开始时间失败: %v", err)
	}
	// 解析结束时间
	endTime, err := time.Parse(timeLayout, endTimeStr)
	if err != nil {
		return false, fmt.Errorf("解析结束时间失败: %v", err)
	}
	// 解析目标时间
	targetTime, err := time.Parse(timeLayout, targetTimeStr)
	if err != nil {
		return false, fmt.Errorf("解析目标时间失败: %v", err)
	}
	// 判断是否在范围内（包含边界）
	return (targetTime.Equal(startTime) || targetTime.Equal(endTime)) ||
		(targetTime.After(startTime) && targetTime.Before(endTime)), nil
}
