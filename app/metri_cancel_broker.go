/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    metri_cancel_broker
 *@date    2025/9/1 8:49
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

func (engine *threadMap) cancelBroker(c *gin.Context) {
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

	type labelmsg struct {
		LabelId []string `json:"labelId"`
	}
	var ids labelmsg
	data, _ := c.GetRawData()
	json.Unmarshal(data, &ids)
	util.Loggrs.Info("cancel label:", strings.Join(ids.LabelId, ","))

	if ids.LabelId == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "success:nil,failed:nil",
		})
		return
	}
	tag := strings.Split(ids.LabelId[0], ".")
	r := db.Exec(fmt.Sprintf("cancel load from %s where label = '%s'", tag[0], tag[1]))
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
