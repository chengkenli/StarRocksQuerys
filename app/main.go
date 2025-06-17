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
        admin.POST("/disconnect", disconnect)
    }
    r.Run(fmt.Sprintf(":%d", util.Read.Server.Port))
}
