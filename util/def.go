/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package util
 *@file    def
 *@date    2024/8/7 14:44
 */

package util

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Config      *viper.Viper
	Loggrs      *logrus.Logger
	P           ArvgParms
	MetaLink    []map[string]interface{}
	ClientIPDec []ClientIPData
	H           Hosts
)

type CustomLogger struct{}

func (l *CustomLogger) Errorf(format string, v ...interface{}) {}
func (l *CustomLogger) Warnf(format string, v ...interface{})  {} // 忽略 WARN
func (l *CustomLogger) Debugf(format string, v ...interface{}) {}

type ArvgParms struct {
	Help bool
}

type Hosts struct {
	Ip string
}

type ConnectParms struct {
	Host       string
	Port       int
	User       string
	Pass       string
	Base       string
	Uri        string
	AaccessKey string
	SecretKey  string
}
type ExteData struct {
	Ctime     string
	Size2gb   string
	Size2v    string
	Bucket    string
	Replica   string
	Rowcount  string
	Tablename string
}
type Message struct {
	Title   string
	Context string
	LogUrl  string
	Roboot  []string
	Err     error
}

type ClientIPData struct {
	Ip              string `json:"ip"`
	User            string `json:"user"`
	SystemName      string `json:"system_name"`
	ConsoleUser     string `json:"console_user"`
	Manufacturer    string `json:"manufacturer"`
	Model           string `json:"model"`
	OperatingSystem string `json:"operating_system"`
	Timestamp       string `json:"timestamp"`
}

type Querisign struct {
	StartTime     string `bson:"StartTime"`
	QueryId       string `bson:"QueryId"`
	ConnectionId  string `bson:"ConnectionId"`
	Database      string `bson:"Database"`
	User          string `bson:"User"`
	ScanBytes     string `bson:"ScanBytes"`
	ScanRows      string `bson:"ScanRows"`
	MemoryUsage   string `bson:"MemoryUsage"`
	DiskSpillSize string `bson:"DiskSpillSize"`
	CPUTime       string `bson:"CPUTime"`
	ExecTime      string `bson:"ExecTime"`
	Warehouse     string `bson:"Warehouse"`
}

type Queris []struct {
	StartTime     string `bson:"StartTime"`
	QueryId       string `bson:"QueryId"`
	ConnectionId  string `bson:"ConnectionId"`
	Database      string `bson:"Database"`
	User          string `bson:"User"`
	ScanBytes     string `bson:"ScanBytes"`
	ScanRows      string `bson:"ScanRows"`
	MemoryUsage   string `bson:"MemoryUsage"`
	DiskSpillSize string `bson:"DiskSpillSize"`
	CPUTime       string `bson:"CPUTime"`
	ExecTime      string `bson:"ExecTime"`
	Warehouse     string `bson:"Warehouse"`
}

// QueryResource 表示单个查询的资源消耗信息
type QueryResource struct {
	ID          int64   `json:"ID"`          // 查询ID(可能为0)
	User        string  `json:"User"`        // 执行用户
	ScanBytes   int64   `json:"ScanBytes"`   // 扫描数据量(字节)
	ScanRows    int64   `json:"ScanRows"`    // 扫描行数
	MemoryUsage int64   `json:"MemoryUsage"` // 内存使用量(字节)
	CPUTime     float64 `json:"CPUTime"`     // CPU时间(秒)
	ExecTime    float64 `json:"ExecTime"`    // 执行时间(秒)
}

type Process struct {
	Id        string `bson:"Id"`
	User      string `bson:"User"`
	Host      string `bson:"Host"`
	Cluster   string `bson:"Cluster"`
	Db        string `bson:"Db"`
	Command   string `bson:"Command"`
	Time      string `bson:"Time"`
	State     string `bson:"State"`
	Info      string `bson:"Info"`
	IsPending string `bson:"IsPending"`
	Warehouse string `bson:"Warehouse"`
}

type ShortData struct {
	Name                   string      `bson:"name"`
	Id                     int         `bson:"id"`
	CpuCoreLimit           int         `bson:"cpu_core_limit"`
	MemLimit               string      `bson:"mem_limit"`
	MaxCpuCores            interface{} `bson:"max_cpu_cores"`
	BigQueryCpuSecondLimit interface{} `bson:"big_query_cpu_second_limit"`
	BigQueryScanRowsLimit  interface{} `bson:"big_query_scan_rows_limit"`
	BigQueryMemLimit       interface{} `bson:"big_query_mem_limit"`
	ConcurrencyLimit       interface{} `bson:"concurrency_limit"`
	SpillMemLimitThreshold interface{} `bson:"spill_mem_limit_threshold"`
	Type                   string      `bson:"type"`
	Classifiers            string      `bson:"classifiers"`
}

type HexData struct {
	Status string `json:"status"`
	Data   struct {
		StartTime    int    `json:"startTime"`
		EndTime      int    `json:"endTime"`
		TimeUsedMs   int    `json:"timeUsedMs"`
		QueryID      string `json:"queryId"`
		ClientIP     string `json:"clientIp"`
		State        string `json:"state"`
		SQL          string `json:"sql"`
		User         string `json:"user"`
		CPUCostNs    int    `json:"cpuCostNs"`
		MemCostBytes int    `json:"memCostBytes"`
		ScanRows     int    `json:"scanRows"`
		ScanBytes    int    `json:"scanBytes"`
		Digest       string `json:"digest"`
		Warehouse    string `json:"warehouse"`
		FailedReason string `json:"failedReason"`
		Database     string `json:"database"`
	} `json:"data"`
}
