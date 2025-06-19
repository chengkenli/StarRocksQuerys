/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package logs
 *@file    main
 *@date    2025/4/30 12:20
 */

package logs

import (
	"StarRocksQuerys/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

// Logserver
// 日志服务
func Logserver(r *gin.Engine) {
	r.GET("/log/*path", viewLog)
}

func viewLog(c *gin.Context) {
	logFile := c.Param("path")
	util.Loggrs.Info(fmt.Sprintf("LOGS> %s read:[%s]", time.Now().Format("2006-01-02 15:04:05"), logFile))
	fh, err := ioutil.ReadFile(logFile)
	if err != nil {
		c.String(http.StatusBadRequest, "no log file.")
	} else {
		c.String(http.StatusOK, string(fh))
	}
}
