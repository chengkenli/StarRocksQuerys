/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package app
 *@file    main
 *@date    2025/2/12 14:03
 */

package app

import (
	"StarRocksAPIs/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func App() {
	r := gin.Default()
	// 加载HTML模板
	r.LoadHTMLGlob(util.Config.GetString("server.loadhtmlglob"))
	r.Static("/static", util.Config.GetString("server.loadstatic"))
	r.POST("/api/verify-token", verify)

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index3.html", nil)
	})
	r.GET("/getstmtid", metriStmtId)
	// 受保护的管理员路由
	admin := r.Group("/api")
	admin.Use(tokenAuth()) // 应用管理员认证中间件
	{
		admin.POST("/kills", meticSleep)
		admin.POST("/killall", metriKillall)
		admin.POST("/killw", metriKillw)
		admin.POST("/killc", metriKillc)
		admin.POST("/killone", metriKillone)
		admin.POST("/showerr", ShowErr)
		admin.GET("/query", processlist)
		admin.GET("/queries", metriQuery)
		admin.GET("/slowquery", metriSlowquery)
		admin.GET("/appids", metriGetAppid)
	}
	r.Run(":" + util.Config.GetString("server.port"))
}
