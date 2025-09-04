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
)

var (
	Loggrs          *logrus.Logger
	P               ArvgParms
	MetaLink        []map[string]interface{}
	ClientIPDec     []ClientData
	ClientMapReduce []ClientMapReduceIP
	H               Hosts
	Read            ReadInConf
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

type ClientData struct {
	Ts              string `bson:"ts"`
	ComputerName    string `bson:"computer_name"`
	UserName        string `bson:"user_name"`
	ComputerType    string `bson:"computer_type"`
	ComputerStatus  string `bson:"computer_status"`
	IpAddress       string `bson:"ip_address"`
	SerialNumber    string `bson:"serial_number"`
	Brand           string `bson:"brand"`
	Model           string `bson:"model"`
	ComputerVersion string `bson:"computer_version"`
	BusinessUnit    string `bson:"business_unit"`
	BusinessFormat  string `bson:"business_format"`
	DataSource      string `bson:"data_source"`
	LastUpdateTime  string `bson:"last_update_time"`
	AiTime          string `bson:"ai_time"`
	AddTime         string `bson:"add_time"`
	LastActiveTime  string `bson:"last_active_time"`
}

type ClientMapReduceIP struct {
	Ip      string `bson:"ip"`
	Dev     string `bson:"dev"`
	Dept    string `bson:"dept"`
	Comment string `bson:"comment"`
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

type ReadInConf struct {
	Metadb struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Base     string `json:"base"`
	} `json:"metadb"`
	Server struct {
		Port         int    `json:"port"`
		Token        string `json:"token"`
		Loadhtmlglob string `json:"loadhtmlglob"`
		Loadstatic   string `json:"loadstatic"`
	} `json:"server"`
	Schema struct {
		Role struct {
			Admin string `json:"admin"`
			Point string `json:"point"`
		} `json:"role"`
		Ipsystem  string   `json:"ipsystem"`
		Slowstmt  string   `json:"slowstmt"`
		Auditops  string   `json:"auditops"`
		Emrwedat  string   `json:"emrwedat"`
		Ipapp     string   `json:"ipapp"`
		Whitelist string   `json:"whitelist"`
		Gybi      []string `json:"gybi"`
	} `json:"schema"`
	Category struct {
		Num0 []string `json:"0"`
		Num1 []string `json:"1"`
		Num2 []string `json:"2"`
	} `json:"category"`
	Log struct {
		Path string `json:"path"`
	} `json:"log"`
}

type GlobalQueries struct {
	StartTime     string `json:"StartTime"`
	FeIp          string `json:"FeIp"`
	QueryId       string `json:"QueryId"`
	ConnectionId  int64  `json:"ConnectionId"`
	Database      string `json:"Database"`
	User          string `json:"User"`
	ScanBytes     string `json:"ScanBytes"`
	ScanRows      string `json:"ScanRows"`
	MemoryUsage   string `json:"MemoryUsage"`
	DiskSpillSize string `json:"DiskSpillSize"`
	CPUTime       string `json:"CPUTime"`
	ExecTime      string `json:"ExecTime"`
	Warehouse     string `json:"Warehouse"`
	CustomQueryId string `json:"CustomQueryId"`
	ResourceGroup string `json:"ResourceGroup"`
}

type HttpSqlApi struct {
	Meta []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"meta"`
	Data []struct {
		ID                  string `json:"Id"`
		User                string `json:"User"`
		Host                string `json:"Host"`
		Db                  string `json:"Db"`
		Command             string `json:"Command"`
		ConnectionStartTime string `json:"ConnectionStartTime"`
		Time                string `json:"Time"`
		State               string `json:"State"`
		Info                string `json:"Info"`
		IsPending           string `json:"IsPending"`
		Warehouse           string `json:"Warehouse"`
	} `json:"data"`
	Statistics struct {
		ScanRows   int `json:"scanRows"`
		ScanBytes  int `json:"scanBytes"`
		ReturnRows int `json:"returnRows"`
	} `json:"statistics"`
}

type Fronends []struct {
	Name              string `bson:"Name"`
	IP                string `bson:"IP"`
	EditLogPort       int    `bson:"EditLogPort"`
	HttpPort          int    `bson:"HttpPort"`
	QueryPort         int    `bson:"QueryPort"`
	RpcPort           int    `bson:"RpcPort"`
	Role              string `bson:"Role"`
	IsMaster          bool   `bson:"IsMaster"`
	ClusterId         int    `bson:"ClusterId"`
	Join              string `bson:"Join"`
	Alive             string `bson:"Alive"`
	ReplayedJournalId int    `bson:"ReplayedJournalId"`
	LastHeartbeat     string `bson:"LastHeartbeat"`
	IsHelper          string `bson:"IsHelper"`
	ErrMsg            string `bson:"ErrMsg"`
	StartTime         string `bson:"StartTime"`
	Version           string `bson:"Version"`
}

type Backends []struct {
	BackendId             string `bson:"BackendId"`
	Cluster               string `bson:"Cluster"`
	IP                    string `bson:"IP"`
	HeartbeatPort         int    `bson:"HeartbeatPort"`
	BePort                int    `bson:"BePort"`
	HttpPort              int    `bson:"HttpPort"`
	BrpcPort              int    `bson:"BrpcPort"`
	LastStartTime         string `bson:"LastStartTime"`
	LastHeartbeat         string `bson:"LastHeartbeat"`
	Alive                 string `bson:"Alive"`
	SystemDecommissioned  string `bson:"SystemDecommissioned"`
	ClusterDecommissioned string `bson:"ClusterDecommissioned"`
	TabletNum             int    `bson:"TabletNum"`
	DataUsedCapacity      string `bson:"DataUsedCapacity"`
	AvailCapacity         string `bson:"AvailCapacity"`
	TotalCapacity         string `bson:"TotalCapacity"`
	UsedPct               string `bson:"UsedPct"`
	MaxDiskUsedPct        string `bson:"MaxDiskUsedPct"`
	ErrMsg                string `bson:"ErrMsg"`
	Version               string `bson:"Version"`
	Status                string `bson:"Status"`
	DataTotalCapacity     string `bson:"DataTotalCapacity"`
	DataUsedPct           string `bson:"DataUsedPct"`
}

type Computes []struct {
	ComputeNodeId         int    `bson:"ComputeNodeId"`
	IP                    string `bson:"IP"`
	HeartbeatPort         int    `bson:"HeartbeatPort"`
	BePort                int    `bson:"BePort"`
	HttpPort              int    `bson:"HttpPort"`
	BrpcPort              int    `bson:"BrpcPort"`
	LastStartTime         string `bson:"LastStartTime"`
	LastHeartbeat         string `bson:"LastHeartbeat"`
	Alive                 string `bson:"Alive"`
	SystemDecommissioned  string `bson:"SystemDecommissioned"`
	ClusterDecommissioned string `bson:"ClusterDecommissioned"`
	ErrMsg                string `bson:"ErrMsg"`
	Version               string `bson:"Version"`
	CpuCores              int    `bson:"CpuCores"`
	MemLimit              string `bson:"MemLimit"`
	NumRunningQueries     string `bson:"NumRunningQueries"`
	MemUsedPct            string `bson:"MemUsedPct"`
	CpuUsedPct            string `bson:"CpuUsedPct"`
	DataCacheMetrics      string `bson:"DataCacheMetrics"`
	HasStoragePath        string `bson:"HasStoragePath"`
	StarletPort           int    `bson:"StarletPort"`
	WorkerId              int    `bson:"WorkerId"`
	WarehouseName         string `bson:"WarehouseName"`
	TabletNum             string `bson:"TabletNum"`
}

type PieData struct {
	Title      string   `json:"title"`
	Subtext    string   `json:"subtext"`
	Categories []string `json:"categories"`
	Data       []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"data"`
}

type StreamData struct {
	Label                 string `json:"Label"`
	Id                    string `json:"Id"`
	LoadId                string `json:"LoadId"`
	TxnId                 string `json:"TxnId"`
	DbName                string `json:"DbName"`
	TableName             string `json:"TableName"`
	State                 string `json:"State"`
	ErrorMsg              string `json:"ErrorMsg"`
	TrackingURL           string `json:"TrackingURL"`
	ChannelNum            string `json:"ChannelNum"`
	PreparedChannelNum    string `json:"PreparedChannelNum"`
	NumRowsNormal         string `json:"NumRowsNormal"`
	NumRowsAbNormal       string `json:"NumRowsAbNormal"`
	NumRowsUnselected     string `json:"NumRowsUnselected"`
	NumLoadBytes          string `json:"NumLoadBytes"`
	TimeoutSecond         string `json:"TimeoutSecond"`
	CreateTimeMs          string `json:"CreateTimeMs"`
	BeforeLoadTimeMs      string `json:"BeforeLoadTimeMs"`
	StartLoadingTimeMs    string `json:"StartLoadingTimeMs"`
	StartPreparingTimeMs  string `json:"StartPreparingTimeMs"`
	FinishPreparingTimeMs string `json:"FinishPreparingTimeMs"`
	EndTimeMs             string `json:"EndTimeMs"`
	ChannelState          string `json:"ChannelState"`
	Type                  string `json:"Type"`
	TrackingSQL           string `json:"TrackingSQL"`
}

type BrokerMsg struct {
	JobId          string `bson:"JobId"`
	Label          string `bson:"Label"`
	State          string `bson:"State"`
	Progress       string `bson:"Progress"`
	Type           string `bson:"Type"`
	Priority       string `bson:"Priority"`
	ScanRows       string `bson:"ScanRows"`
	FilteredRows   string `bson:"FilteredRows"`
	UnselectedRows string `bson:"UnselectedRows"`
	SinkRows       string `bson:"SinkRows"`
	EtlInfo        string `bson:"EtlInfo"`
	TaskInfo       string `bson:"TaskInfo"`
	ErrorMsg       string `bson:"ErrorMsg"`
	CreateTime     string `bson:"CreateTime"`
	EtlStartTime   string `bson:"EtlStartTime"`
	EtlFinishTime  string `bson:"EtlFinishTime"`
	LoadStartTime  string `bson:"LoadStartTime"`
	LoadFinishTime string `bson:"LoadFinishTime"`
	TrackingSQL    string `bson:"TrackingSQL"`
	JobDetails     string `bson:"JobDetails"`
	//自定义参数
	Backends    string
	ScanInfo    string
	Command     string
	GetHour     string
	TimeMsg     string
	LabelName   string
	ProgressMsg string
}

type BrokerJobDetails struct {
	AllBackends            interface{} `json:"All backends"`
	FileNumber             int         `json:"FileNumber"`
	FileSize               int         `json:"FileSize"`
	InternalTableLoadBytes int         `json:"InternalTableLoadBytes"`
	InternalTableLoadRows  int         `json:"InternalTableLoadRows"`
	ScanBytes              int         `json:"ScanBytes"`
	ScanRows               int         `json:"ScanRows"`
	TaskNumber             int         `json:"TaskNumber"`
	UnfinishedBackends     interface{} `json:"Unfinished backends"`
}

type ResourceMsg struct {
	Name              string      `bson:"name"`
	User              string      `bson:"user"`
	CpuWeight         interface{} `bson:"cpu_weight"`
	ExclusiveCpuCores interface{} `bson:"exclusive_cpu_cores"`
	MemLimit          interface{} `bson:"mem_limit"`
	ConcurrencyLimit  interface{} `bson:"concurrency_limit"`
	BigQueryCPU       interface{} `bson:"big_query_cpu_second_limit"`
	BigQueryRows      interface{} `bson:"big_query_scan_rows_limit"`
	BigQueryMemLimit  interface{} `bson:"big_query_mem_limit"`
	// 自定义
	Command string
}

type TaskRuns struct {
	QUERY_ID      string `bson:"QUERY_ID"`
	TASK_NAME     string `bson:"TASK_NAME"`
	CREATE_TIME   string `bson:"CREATE_TIME"`
	FINISH_TIME   string `bson:"FINISH_TIME"`
	STATE         string `bson:"STATE"`
	CATALOG       string `bson:"CATALOG"`
	WAREHOUSE     string `bson:"WAREHOUSE"`
	DATABASE      string `bson:"DATABASE"`
	DEFINITION    string `bson:"DEFINITION"`
	EXPIRE_TIME   string `bson:"EXPIRE_TIME"`
	ERROR_CODE    string `bson:"ERROR_CODE"`
	ERROR_MESSAGE string `bson:"ERROR_MESSAGE"`
	PROGRESS      string `bson:"PROGRESS"`
	EXTRA_MESSAGE string `bson:"EXTRA_MESSAGE"`
	PROPERTIES    string `bson:"PROPERTIES"`
	// 自定义
	GetHour     string
	ProgressInt int
	TimeMsg     string
	QueryIdName string
	Command     string
}
