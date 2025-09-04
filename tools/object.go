/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package tools
 *@file    object
 *@date    2025/2/12 14:01
 */

package tools

import (
	"StarRocksQuerys/util"
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"math"
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

func FormatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	if bytes <= 0 {
		return "0B"
	}
	// 计算单位索引
	unitIndex := math.Floor(math.Log(float64(bytes)) / math.Log(1024))
	// 防止超出单位数组范围
	if unitIndex >= float64(len(units)) {
		unitIndex = float64(len(units) - 1)
	}
	// 计算转换后的值
	value := float64(bytes) / math.Pow(1024, unitIndex)
	// 格式化输出，保留1位小数
	return fmt.Sprintf("%.1f%s", value, units[int(unitIndex)])
}

func formatSize(sizeMB float64) string {
	const (
		GB = 1024
		TB = 1024 * GB
	)
	switch {
	case sizeMB >= TB:
		sizeTB := math.Round(sizeMB/TB*10) / 10
		return fmt.Sprintf("%.0fTB", sizeTB)
	case sizeMB >= GB:
		sizeGB := math.Round(sizeMB/GB*10) / 10
		return fmt.Sprintf("%.0fGB", sizeGB)
	default:
		return fmt.Sprintf("%.0fMB", sizeMB)
	}
}

func Version(db *gorm.DB) float64 {
	sql := fmt.Sprintf("select current_version() as version")
	var m map[string]interface{}
	db.Raw(sql).Scan(&m)

	var arr string
	if m["version"] != nil {
		arr = strings.Split(m["version"].(string), " ")[0]
	}
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

func GetUserOperational(db *gorm.DB, app string, cachelimit *cache.Cache) {
	util.Loggrs.Info("cache ", app, " loading...")
	if util.Read.Schema.Auditops == "" {
		return
	}
	// 计算CPU
	var cpu []map[string]interface{}
	r := db.Raw(fmt.Sprintf(`SELECT
user,
SUM(cpuCostNs) / 1e9 AS total_cpu_seconds, 
(SUM(cpuCostNs) / (SELECT SUM(cpuCostNs) FROM
%s WHERE state IN ('EOF','OK') AND
timestamp >= DATE_SUB(NOW(), INTERVAL 1 DAY))) * 100 AS
cpu_usage_percentage 
FROM %s
WHERE state IN ('EOF','OK') 
AND timestamp >= '2025-09-01 00:00:00'
GROUP BY user
ORDER BY total_cpu_seconds DESC`, util.Read.Schema.Auditops, util.Read.Schema.Auditops)).Scan(&cpu)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
	}
	// 计算内存
	var mem []map[string]interface{}
	r = db.Raw(fmt.Sprintf(`SELECT
user,
MAX(memCostBytes) / (1024 * 1024) AS max_mem_mb 
FROM %s
WHERE state IN ('EOF','OK') 
AND timestamp >= DATE_SUB(NOW(), INTERVAL 1 DAY)
GROUP BY user
ORDER BY max_mem_mb DESC`, util.Read.Schema.Auditops)).Scan(&mem)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
	}
	// 计算并发
	var concurrency []map[string]interface{}
	r = db.Raw(fmt.Sprintf(`WITH UserConcurrency AS (
SELECT
user,
DATE_FORMAT(timestamp, '%%Y-%%m-%%d %%H:%%i') AS minute_bucket,
COUNT(*) AS query_concurrency
FROM %s
WHERE state IN ('EOF', 'OK')
AND timestamp >= DATE_SUB(NOW(), INTERVAL 1 DAY)
AND LOWER(stmt) LIKE '%%select%%' 
GROUP BY user, minute_bucket
HAVING query_concurrency > 1 
)
SELECT
user,
minute_bucket,
query_concurrency / 60.0 AS query_concurrency_per_second 
FROM (
SELECT
user,
minute_bucket,
query_concurrency,
ROW_NUMBER() OVER (PARTITION BY user ORDER BY query_concurrency DESC)
AS rn
FROM UserConcurrency
) ranked
WHERE rn = 1 
ORDER BY query_concurrency_per_second DESC`, util.Read.Schema.Auditops)).Scan(&concurrency)
	if r.Error != nil {
		util.Loggrs.Error(r.Error.Error())
	}

	for _, item := range combineMetrics(cpu, mem, concurrency) {
		key := app + item["user"].(string) + "operational"

		var cpu_usage_percentage, max_mem_mb, query_concurrency_per_second string
		if item["cpu_usage_percentage"] != nil {
			cpu_usage_percentage = fmt.Sprintf("<strong>%.2f%%</strong>", item["cpu_usage_percentage"].(float64))
		} else {
			cpu_usage_percentage = "0%"
		}
		if item["max_mem_mb"] != nil {
			max_mem_mb = formatSize(item["max_mem_mb"].(float64))
			if strings.Contains(max_mem_mb, "GB") {
				max_mem_mb = fmt.Sprintf(`<strong style="color: #FFB347;">%s</strong>`, max_mem_mb)
			} else if strings.Contains(max_mem_mb, "TB") {
				max_mem_mb = fmt.Sprintf(`<strong style="color: red;">%s</strong>`, max_mem_mb)
			} else {
				max_mem_mb = fmt.Sprintf(`<strong style="color: green;">%s</strong>`, max_mem_mb)
			}
		} else {
			max_mem_mb = "0"
		}
		if item["query_concurrency_per_second"] != nil {
			atoi, _ := strconv.Atoi(strings.Split(item["query_concurrency_per_second"].(string), ".")[0])
			if atoi >= 100 {
				query_concurrency_per_second = fmt.Sprintf(`<strong style="color: red;">%d</strong>`, atoi)
			} else if atoi >= 10 && atoi < 50 {
				query_concurrency_per_second = fmt.Sprintf(`<strong style="color: #FFB347;">%d</strong>`, atoi)
			} else {
				query_concurrency_per_second = fmt.Sprintf(`<strong style="color: green;">%d</strong>`, atoi)
			}
		} else {
			query_concurrency_per_second = "0"
		}
		val := fmt.Sprintf(`CPU:%v MEM:%v CON:%v`, cpu_usage_percentage, max_mem_mb, query_concurrency_per_second)
		if _, ok := cachelimit.Get(key); !ok {
			cachelimit.Set(key, val, cache.DefaultExpiration)

			marshal, _ := json.Marshal(item)
			msg := strings.NewReplacer(
				"{", "",
				"}", "",
				"[", "",
				"]", "",
				`"`, "'",
				",", ",\n",
			).Replace(string(marshal))
			cachelimit.Set(key+"msg", msg, cache.DefaultExpiration)
		}
	}
}

func combineMetrics(cpu, mem, concurrency []map[string]interface{}) []map[string]interface{} {
	// 创建一个map来存储按用户分组的结果
	result := make(map[string]map[string]interface{})
	// 处理CPU数据
	for _, entry := range cpu {
		user := entry["user"].(string)
		if _, exists := result[user]; !exists {
			result[user] = make(map[string]interface{})
			result[user]["user"] = user
		}
		result[user]["total_cpu_seconds"] = entry["total_cpu_seconds"]
		result[user]["cpu_usage_percentage"] = entry["cpu_usage_percentage"]
	}
	// 处理内存数据
	for _, entry := range mem {
		user := entry["user"].(string)
		if _, exists := result[user]; !exists {
			result[user] = make(map[string]interface{})
			result[user]["user"] = user
		}
		result[user]["max_mem_mb"] = entry["max_mem_mb"]
	}
	// 处理并发数据
	for _, entry := range concurrency {
		user := entry["user"].(string)
		if _, exists := result[user]; !exists {
			result[user] = make(map[string]interface{})
			result[user]["user"] = user
		}
		result[user]["minute_bucket"] = entry["minute_bucket"]
		result[user]["query_concurrency_per_second"] = entry["query_concurrency_per_second"]
	}
	// 将map转换为slice
	var combined []map[string]interface{}
	for _, userData := range result {
		combined = append(combined, userData)
	}

	return combined
}
