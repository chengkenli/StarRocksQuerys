/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package tools
 *@file    object
 *@date    2025/2/12 14:01
 */

package tools

import (
	"StarRocksAPIs/util"
	"bufio"
	"crypto/rand"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// RmDuplicaSlice /*数组去重*/
func RmDuplicaSlice(strs []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复字符串
	for _, e := range strs {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func Size(s string) float64 {
	s = strings.ToLower(s)
	//正则
	re := regexp.MustCompile(`\d+\.?\d*`)
	float, _ := strconv.ParseFloat(re.FindString(s), 64)
	if strings.Contains(s, "kb") {
		return float * 1024
	}
	if strings.Contains(s, "mb") {
		return float * 1024 * 1024
	}
	if strings.Contains(s, "gb") {
		return float * 1024 * 1024 * 1024
	}
	if strings.Contains(s, "tb") {
		return float * 1024 * 1024 * 1024 * 1024
	}
	return 0
}

func Version(db *gorm.DB) float64 {
	sql := fmt.Sprintf("select current_version() as version")
	var m map[string]interface{}
	db.Raw(sql).Scan(&m)

	arr := strings.Split(m["version"].(string), " ")[0]
	if len(arr) < 2 {
		return 0
	}
	if !strings.Contains(arr, ".") {
		return 0
	}

	version, err := strconv.ParseFloat(fmt.Sprintf("%s.%s", strings.Split(arr, ".")[0], strings.Split(arr, ".")[1]), 64)
	if err != nil {
		return 0
	}
	return version
}

// WriteFile 文件落地
func WriteFile(fname, msg string) {
	fileHandle, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriterSize(fileHandle, len(msg))

	buf.WriteString(msg)

	err = buf.Flush()
	if err != nil {
		util.Loggrs.Error(err.Error())
		return
	}
}

// StrInSlice 检查数组中是否存在某个元素
func StrInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func VToFloat(version string) float64 {
	// Split the string by '-' and take the first part
	parts := strings.Split(version, "-")
	mainPart := parts[0]
	// Split the main part by '.'
	subParts := strings.Split(mainPart, ".")
	// Construct the new version format as float
	var newVersionFloat float64
	var err error
	if len(subParts) >= 2 {
		// Combine the first two parts and convert to float
		combined := subParts[0] + "." + subParts[1]
		if len(subParts[2]) > 1 {
			combined += subParts[2][:2] // Append the second character of the third part if it exists
		} else {
			combined += subParts[2]
		}
		newVersionFloat, err = strconv.ParseFloat(combined, 64)
		if err != nil {
			fmt.Println(err.Error())
			return 0
		}
	}
	return newVersionFloat
}

// RandomPassWord 根据指定长度的生产高敏感度字符串
func RandomPassWord(n int) string {
	allowedChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		// 生成一个随机索引
		ri, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChars))))
		if err != nil {
			return ""
		}
		// 使用随机索引获取一个字符
		b[i] = allowedChars[ri.Int64()]
	}

	return string(b)
}

func GetHour(second int) (formatString string) {
	// GetHour 秒格式化
	hours := second / 3600
	minutes := (second % 3600) / 60
	secs := second % 60

	if hours >= 1 {
		return fmt.Sprintf("%02dh:%02dmin:%02ds", hours, minutes, secs)
	}
	if minutes >= 1 {
		return fmt.Sprintf("%02dmin:%02ds", minutes, secs)
	}
	return fmt.Sprintf("%02ds", secs)
}
