/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_cancel_submit
 *@date    2025/9/1 10:29
 */

package app

import (
	"StarRocksQuerys/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (engine *threadMap) cancelSubmit(c *gin.Context) {
	appid := c.GetHeader("AppID")
	if appid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "app is nil."})
		return
	}
	appid = setdefault(appid)
	db, err := engine.getmapConnect(appid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	type taskmsg struct {
		Taskname []string `json:"taskname"`
	}
	var ids taskmsg
	data, _ := c.GetRawData()
	json.Unmarshal(data, &ids)
	util.Loggrs.Info("cancel label:", strings.Join(ids.Taskname, ","))

	if ids.Taskname == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "success:nil,failed:nil",
		})
		return
	}
	r := db.Exec(fmt.Sprintf("DROP TASK %s ", ids.Taskname[0]))
	if r.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "success:nil,failed:db exec stmt is failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("success:%d,failed:%d", 1, 0),
	})
}
