/*
 *@author  chengkenli
 *@project healthyreport
 *@package init
 *@file    init
 *@date    2025/1/17 15:22
 */

package init

import (
	"StarRocksAPIs/conn"
	"StarRocksAPIs/util"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	printStarRocks()

	ticker := time.NewTicker(time.Second)
	var i int
	for {
		select {
		case <-ticker.C:
			if util.Config.GetString("metadb.host") == "" ||
				util.Config.GetString("metadb.port") == "" ||
				util.Config.GetString("metadb.user") == "" ||
				util.Config.GetString("metadb.password") == "" ||
				util.Config.GetString("metadb.base") == "" {
				if i > 5 {
					util.Loggrs.Warn("元数据库配置异常，5 second后将退出!")
					os.Exit(-1)
				}
				i++
			} else {
				util.Loggrs.Info("元数据库读取成功!")
				goto loop
			}
		}
	}
loop:
	connect, err := conn.ConnectMySQL()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	r := connect.Raw(fmt.Sprintf("select * from %s where status >= 1", util.Config.GetString("metadb.base"))).Scan(&util.MetaLink)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
}

func printStarRocks() {
	c := color.New()
	s := c.Add(color.FgHiGreen).Sprint(`
            ____  _             ____            _                      
           / ___|| |_ __ _ _ __|  _ \ ___   ___| | _____               
           \___ \| __/ __ | ___| |_) / _ \ / __| |/ / __|              
            ___) | || (_| | |  |  _ < (_) | (__|   <\__ \              
           |____/ \__\____|_|  |_| \_\___/ \___|_|\_\___/              `)
	v := fmt.Sprintf("%s %s %s\n%s %s %s\n%s %s %s\n%s %s %s\n%s %s %s",
		strings.Repeat("\t", 6),
		c.Add(color.FgHiMagenta).Sprint(filepath.Base(os.Args[0])),
		c.Add(color.FgHiYellow).Sprint(getVersion()),
		strings.Repeat("\t", 6),
		c.Add(color.FgHiWhite).Sprint("md5sum"),
		c.Add(color.FgHiYellow).Sprint("@"+getHash()),
		strings.Repeat("\t", 6),
		c.Add(color.FgHiWhite).Sprint("author"),
		c.Add(color.FgHiYellow).Sprint("@chengken li"),
		strings.Repeat("\t", 6),
		c.Add(color.FgHiWhite).Sprint("system"),
		c.Add(color.FgHiYellow).Sprint("@"+getUser()),
		strings.Repeat("\t", 6),
		c.Add(color.FgHiWhite).Sprint("update"),
		c.Add(color.FgHiYellow).Sprint("@"+getStat()),
	)
	fmt.Println(fmt.Sprintf("%s\n\n%s\n", s, v))
}

func getHash() string {
	var md5Sum string
	if v, err := os.Executable(); err == nil {
		file, _ := os.Open(v)
		defer file.Close()
		hash := md5.New()
		if _, err := io.Copy(hash, file); err == nil {
			md5Sum = hex.EncodeToString(hash.Sum(nil))
		}
	}
	return md5Sum
}

func getStat() string {
	var ts string
	if v, err := os.Executable(); err == nil {
		s, _ := os.Stat(v)
		ts = s.ModTime().Format("2006-01-02 15:04:05")
	}
	return ts
}

func getUser() string {
	u, _ := user.Current()
	return u.Username
}

func getVersion() string {
	file := "/tmp/.v." + filepath.Base(os.Args[0])
	v, _ := ioutil.ReadFile(file)
	if v != nil {
		return "v" + strings.NewReplacer("\n", "").Replace(string(v))
	}
	return ""
}
