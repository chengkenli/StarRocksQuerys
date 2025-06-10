/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package app
 *@file    auth_token
 *@date    2025/5/26 21:20
 */

package app

import (
    "StarRocksQuerys/util"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

func tokenAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            token = c.Query("token") // 也可以从查询参数获取
        }
        if token != util.Config.GetString("server.token")+time.Now().Format("0102") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "无效或缺失访问Token",
            })
            return
        }
        c.Next()
    }
}

func verify(c *gin.Context) {
    data, err := c.GetRawData()
    if err != nil {
        util.Loggrs.Warn(err.Error())
        return
    }
    util.Loggrs.Info(string(data))
    var m map[string]interface{}
    json.Unmarshal(data, &m)

    if m["token"].(string) != util.Config.GetString("server.token")+time.Now().Format("0102") {
        c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
        return
    }
    c.JSON(http.StatusOK, gin.H{"valid": true})
}
