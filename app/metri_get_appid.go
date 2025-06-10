/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package app
 *@file    metri_get_appid
 *@date    2025/5/28 21:45
 */

package app

import (
	"StarRocksAPIs/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func metriGetAppid(c *gin.Context) {
	type applist struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	var meta []applist
	for _, m := range util.MetaLink {
		meta = append(meta, applist{
			Name:  m["nickname"].(string),
			Value: m["app"].(string),
		})
	}
	util.Loggrs.Info(meta)
	c.JSON(http.StatusOK, meta)
}
