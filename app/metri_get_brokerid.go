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
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (engine *threadMap) metriBrokerId(c *gin.Context) {
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
	sign := strings.Split(jobid, ".")
	var m map[string]interface{}
	r := db.Raw(fmt.Sprintf("show load from %s where label='%s'", sign[0], sign[1])).Scan(&m)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
		return
	}
	marshal, _ := json.MarshalIndent(m, "", "  ")
	filename := fmt.Sprintf("%s/%s.%s.log", util.Read.Log.Path, appid, jobid)
	tools.WriteFile(filename, string(marshal))
	uri := fmt.Sprintf("http://%s:%d/log/%s", util.H.Ip, util.Read.Server.Port, filename)
	c.JSON(http.StatusOK, uri)
}
