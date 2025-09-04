/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    main
 *@date    2025/2/12 14:03
 */

package app

import (
	"StarRocksQuerys/logs"
	"StarRocksQuerys/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func App() {

	r := gin.Default()
	//内部启动日志访问器
	go logs.Logserver(r)
	// 加载HTML模板
	r.LoadHTMLGlob(util.Read.Server.Loadhtmlglob)
	r.Static("/static", util.Read.Server.Loadstatic)
	r.POST("/api/verify-token", verify)

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index3.html", nil)
	})
	r.GET("/getstmtid", engine.metriStmtId)
	r.GET("/getbrokid", engine.metriBrokerId)
	r.GET("/getsubmit", engine.metriSubmitId)
	r.GET("/getstream", engine.metriStreamId)
	// 受保护的管理员路由
	admin := r.Group("/api")
	admin.Use(tokenAuth()) // 应用管理员认证中间件
	{
		admin.POST("/kills", engine.meticSleep)
		admin.POST("/killall", engine.metriKillall)
		admin.POST("/killw", engine.metriKillw)
		admin.POST("/killc", engine.metriKillc)
		admin.POST("/killone", engine.metriKillone)
		admin.POST("/showerr", engine.ShowErr)
		admin.POST("/broker-err", engine.showBrokerErr)
		admin.GET("/broker", engine.showBroker)
		admin.POST("/cancel-labelId", engine.cancelBroker)
		admin.GET("/submit", engine.showSubmit)
		admin.POST("/cancel-task", engine.cancelSubmit)
		admin.POST("/submit-err", engine.showSubmitErr)
		admin.GET("/query", engine.processlist)
		admin.GET("/queries", engine.metriQuery)
		admin.GET("/slowquery", engine.metriSlowquery)
		admin.GET("/appids", engine.metriGetAppid)
		admin.POST("/disconnect", engine.disconnect)
		admin.GET("/beads", engine.metriErrbead)
		admin.POST("/pie", engine.pie)
		admin.GET("/getgrafan", engine.metrigrafana)
		admin.POST("/stream-err", engine.showStreamErr)
		admin.GET("/resourcegroup", engine.resourceGroup)
	}
	r.Run(fmt.Sprintf(":%d", util.Read.Server.Port))
}
