
-- chengken.starrocks_information_connections definition

CREATE TABLE `starrocks_information_connections` (
`app` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '集群名称(英文)',
`nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '别名',
`alias` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '集群别名',
`feip` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '集群连接地址(必填)F5,VIP,CLB,FE',
`user` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '集群登录账号(必填) 建议是管理员角色的账号',
`password` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '集群登录密码(必填)',
`feport` int NOT NULL DEFAULT '9030' COMMENT '集群登录端口，默认9030',
`address` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'MANAGER地址，如果填了MANAGER地址，那么将触发定时检查LICENSE是否过期(企业级)',
`expire` int DEFAULT '30' COMMENT 'LICENSE是否过期(企业级)过期提醒倒计时，单位day',
`status` int NOT NULL DEFAULT '0' COMMENT 'LICENSE是否过期(企业级)开关,0 off, 1 on',
`fe_log_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'FE 日志目录',
`be_log_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'BE 日志目录',
`be_meta_log` varchar(200) DEFAULT NULL,
`java_udf_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'BE 日志目录',
`manager_access_key` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'manager 开发者的access key',
`manager_secret_key` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'manager 开发者的secret key',
`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='StarRocks登录配置，manager地址,(定期检查license过期日期)';


CREATE TABLE `ops_starrocks_ip_system` (
`ip` varchar(100) NOT NULL COMMENT "IP",
`user` varchar(100) NULL COMMENT "用户",
`system_name` varchar(200) NULL COMMENT "主机名称",
`console_user` varchar(100) NULL COMMENT "域",
`manufacturer` varchar(100) NULL COMMENT "制造商",
`model` varchar(200) NULL COMMENT "机型",
`operating_system` varchar(200) NULL COMMENT "系统",
`timestamp` varchar(100) NOT NULL COMMENT "数据载入时间"
) ENGINE=OLAP
PRIMARY KEY(`ip`)
COMMENT "starrocks client ip信息关系表"
DISTRIBUTED BY HASH(`ip`) BUCKETS 10
PROPERTIES (
    "compression" = "LZ4",
    "enable_persistent_index" = "true",
    "fast_schema_evolution" = "false",
    "replicated_storage" = "true",
    "replication_num" = "3"
);


CREATE TABLE `cn_asset_information` (
`ts` date NOT NULL COMMENT "分区日期",
`computer_name` varchar(65533) NOT NULL COMMENT "电脑名称",
`user_name` varchar(65533) NOT NULL COMMENT "使用者用户名称",
`computer_type` varchar(65533) NOT NULL COMMENT "电脑类型",
`computer_status` varchar(65533) NOT NULL COMMENT "电脑状态",
`ip_address` varchar(65533) NOT NULL COMMENT "电脑地址",
`serial_number` varchar(65533) NOT NULL COMMENT "电脑序列号",
`brand` varchar(65533) NOT NULL COMMENT "品牌",
`model` varchar(65533) NOT NULL COMMENT "型号",
`computer_version` varchar(65533) NOT NULL COMMENT "windoes版本",
`business_unit` varchar(65533) NOT NULL COMMENT "所属业务部门",
`business_format` varchar(65533) NOT NULL COMMENT "业态",
`data_source` varchar(65533) NOT NULL COMMENT "数据来源",
`last_update_time` date NULL COMMENT "最后更新时间",
`ai_time` date NULL COMMENT "日期",
`add_time` date NULL COMMENT "新增时间",
`last_active_time` date NULL COMMENT "最后活跃时间"
) ENGINE=OLAP
DUPLICATE KEY(`ts`, `computer_name`)
DISTRIBUTED BY RANDOM;



CREATE TABLE `ops_starrocks_schema_slowquery` (
`ts` date NOT NULL COMMENT "数据载入日期",
`app` varchar(64) NOT NULL COMMENT "集群名称",
`queryId` varchar(64) NOT NULL COMMENT "查询的唯一ID",
`origin` varchar(500) NULL COMMENT "SQL原始语句中涉及到的表",
`domain` varchar(64) NULL COMMENT "主题域简称",
`owner` varchar(64) NULL COMMENT "主题域owner",
`action` int(11) NULL COMMENT "拦截行为：⓿.状态异常停留清退,①.异常违规参数查杀,②.10分钟慢查询提醒,③.30分钟慢查询查杀,④.全表扫描亿级查杀,⑤.TB级扫描字节查杀,⑥.百亿扫描行数查杀,⑦.CATALOG违规查杀,⑧.GB级消耗内存查杀",
`timestamp` datetime NOT NULL COMMENT "查询开始时间",
`queryType` varchar(12) NULL COMMENT "查询类型（query, slow_query,connection）",
`clientIp` varchar(32) NULL COMMENT "客户端IP",
`user` varchar(64) NULL COMMENT "查询用户名",
`authorizedUser` varchar(64) NULL COMMENT "用户唯一标识，既user_identity",
`resourceGroup` varchar(64) NULL COMMENT "资源组名",
`catalog` varchar(32) NULL COMMENT "数据目录名",
`db` varchar(96) NULL COMMENT "查询所在数据库",
`state` varchar(8) NULL COMMENT "查询状态（EOF，ERR，OK）",
`errorCode` varchar(96) NULL COMMENT "错误码",
`queryTime` bigint(20) NULL COMMENT "查询执行时间（秒）",
`scanBytes` bigint(20) NULL COMMENT "查询扫描的字节数",
`scanRows` bigint(20) NULL COMMENT "查询扫描的记录行数",
`returnRows` bigint(20) NULL COMMENT "查询返回的结果行数",
`cpuCostNs` bigint(20) NULL COMMENT "查询CPU耗时（纳秒）",
`memCostBytes` bigint(20) NULL COMMENT "查询消耗内存（字节）",
`stmtId` int(11) NULL COMMENT "SQL语句增量ID",
`isQuery` tinyint(4) NULL COMMENT "SQL是否为查询（1或0）",
`feIp` varchar(32) NULL COMMENT "执行该语句的FE IP",
`stmt` varchar(1048576) NULL COMMENT "SQL原始语句",
`digest` varchar(32) NULL COMMENT "慢SQL指纹",
`planCpuCosts` double NULL COMMENT "查询规划阶段CPU占用（纳秒）",
`planMemCosts` double NULL COMMENT "查询规划阶段内存占用（字节）",
`pendingTimeMs` bigint(20) NULL COMMENT "查询在队列中等待的时间（毫秒）",
`logfile` varchar(200) NULL COMMENT "日志文件",
`optimization` int(11) NULL COMMENT "是否做过优化：0.没有优化，1~99.代表优化次数",
`optimizationItems` varchar(1048576) NULL COMMENT "优化项"
) ENGINE=OLAP
PRIMARY KEY(`ts`, `app`, `queryId`)
COMMENT "慢查询审计日志表"
PARTITION BY date_trunc('day', ts)
DISTRIBUTED BY HASH(`ts`, `queryId`)
PROPERTIES (
    "compression" = "LZ4",
    "enable_persistent_index" = "true",
    "fast_schema_evolution" = "false",
    "replicated_storage" = "true",
    "replication_num" = "3"
);



CREATE TABLE `starrocks_audit_log` (
`queryId` varchar(64) NULL COMMENT "查询的唯一ID",
`timestamp` datetime NOT NULL COMMENT "查询开始时间",
`queryType` varchar(12) NULL COMMENT "查询类型（query, slow_query,connection）",
`clientIp` varchar(32) NULL COMMENT "客户端IP",
`user` varchar(64) NULL COMMENT "查询用户名",
`authorizedUser` varchar(64) NULL COMMENT "用户唯一标识，既user_identity",
`resourceGroup` varchar(64) NULL COMMENT "资源组名",
`catalog` varchar(32) NULL COMMENT "数据目录名",
`db` varchar(96) NULL COMMENT "查询所在数据库",
`state` varchar(8) NULL COMMENT "查询状态（EOF，ERR，OK）",
`errorCode` varchar(96) NULL COMMENT "错误码",
`queryTime` bigint(20) NULL COMMENT "查询执行时间（毫秒）",
`scanBytes` bigint(20) NULL COMMENT "查询扫描的字节数",
`scanRows` bigint(20) NULL COMMENT "查询扫描的记录行数",
`returnRows` bigint(20) NULL COMMENT "查询返回的结果行数",
`cpuCostNs` bigint(20) NULL COMMENT "查询CPU耗时（纳秒）",
`memCostBytes` bigint(20) NULL COMMENT "查询消耗内存（字节）",
`stmtId` int(11) NULL COMMENT "SQL语句增量ID",
`isQuery` tinyint(4) NULL COMMENT "SQL是否为查询（1或0）",
`feIp` varchar(32) NULL COMMENT "执行该语句的FE IP",
`stmt` varchar(1048576) NULL COMMENT "SQL原���语句",
`digest` varchar(32) NULL COMMENT "慢SQL指纹",
`planCpuCosts` double NULL COMMENT "查询规划阶段CPU占用（纳秒）",
`planMemCosts` double NULL COMMENT "查询规划阶段内存占用（字节）",
`pendingTimeMs` bigint(20) NULL COMMENT "查询在队列中等待的时间（毫秒）",
`candidateMVs` varchar(65533) NULL COMMENT "候选MV列表",
`hitMvs` varchar(65533) NULL COMMENT "命中MV列表",
`recordTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "数据入库时间"
) ENGINE=OLAP
DUPLICATE KEY(`queryId`, `timestamp`, `queryType`)
COMMENT "审计日志表"
PARTITION BY RANGE(`timestamp`)
(PARTITION p20250507 VALUES [("2025-05-07 00:00:00"), ("2025-05-08 00:00:00")),
PARTITION p20250508 VALUES [("2025-05-08 00:00:00"), ("2025-05-09 00:00:00")),
PARTITION p20250509 VALUES [("2025-05-09 00:00:00"), ("2025-05-10 00:00:00")),
PARTITION p20250510 VALUES [("2025-05-10 00:00:00"), ("2025-05-11 00:00:00")),
PARTITION p20250511 VALUES [("2025-05-11 00:00:00"), ("2025-05-12 00:00:00")),
PARTITION p20250512 VALUES [("2025-05-12 00:00:00"), ("2025-05-13 00:00:00")),
PARTITION p20250513 VALUES [("2025-05-13 00:00:00"), ("2025-05-14 00:00:00")),
PARTITION p20250514 VALUES [("2025-05-14 00:00:00"), ("2025-05-15 00:00:00")),
PARTITION p20250515 VALUES [("2025-05-15 00:00:00"), ("2025-05-16 00:00:00")),
PARTITION p20250516 VALUES [("2025-05-16 00:00:00"), ("2025-05-17 00:00:00")),
PARTITION p20250517 VALUES [("2025-05-17 00:00:00"), ("2025-05-18 00:00:00")),
PARTITION p20250518 VALUES [("2025-05-18 00:00:00"), ("2025-05-19 00:00:00")),
PARTITION p20250519 VALUES [("2025-05-19 00:00:00"), ("2025-05-20 00:00:00")),
PARTITION p20250520 VALUES [("2025-05-20 00:00:00"), ("2025-05-21 00:00:00")),
PARTITION p20250521 VALUES [("2025-05-21 00:00:00"), ("2025-05-22 00:00:00")),
PARTITION p20250522 VALUES [("2025-05-22 00:00:00"), ("2025-05-23 00:00:00")),
PARTITION p20250523 VALUES [("2025-05-23 00:00:00"), ("2025-05-24 00:00:00")),
PARTITION p20250524 VALUES [("2025-05-24 00:00:00"), ("2025-05-25 00:00:00")),
PARTITION p20250525 VALUES [("2025-05-25 00:00:00"), ("2025-05-26 00:00:00")),
PARTITION p20250526 VALUES [("2025-05-26 00:00:00"), ("2025-05-27 00:00:00")),
PARTITION p20250527 VALUES [("2025-05-27 00:00:00"), ("2025-05-28 00:00:00")),
PARTITION p20250528 VALUES [("2025-05-28 00:00:00"), ("2025-05-29 00:00:00")),
PARTITION p20250529 VALUES [("2025-05-29 00:00:00"), ("2025-05-30 00:00:00")),
PARTITION p20250530 VALUES [("2025-05-30 00:00:00"), ("2025-05-31 00:00:00")),
PARTITION p20250531 VALUES [("2025-05-31 00:00:00"), ("2025-06-01 00:00:00")),
PARTITION p20250601 VALUES [("2025-06-01 00:00:00"), ("2025-06-02 00:00:00")),
PARTITION p20250602 VALUES [("2025-06-02 00:00:00"), ("2025-06-03 00:00:00")),
PARTITION p20250603 VALUES [("2025-06-03 00:00:00"), ("2025-06-04 00:00:00")),
PARTITION p20250604 VALUES [("2025-06-04 00:00:00"), ("2025-06-05 00:00:00")),
PARTITION p20250605 VALUES [("2025-06-05 00:00:00"), ("2025-06-06 00:00:00")),
PARTITION p20250606 VALUES [("2025-06-06 00:00:00"), ("2025-06-07 00:00:00")),
PARTITION p20250607 VALUES [("2025-06-07 00:00:00"), ("2025-06-08 00:00:00")),
PARTITION p20250608 VALUES [("2025-06-08 00:00:00"), ("2025-06-09 00:00:00")),
PARTITION p20250609 VALUES [("2025-06-09 00:00:00"), ("2025-06-10 00:00:00")))
DISTRIBUTED BY HASH(`queryId`) BUCKETS 3
PROPERTIES (
    "compression" = "LZ4",
    "dynamic_partition.buckets" = "3",
    "dynamic_partition.enable" = "true",
    "dynamic_partition.end" = "3",
    "dynamic_partition.history_partition_num" = "0",
    "dynamic_partition.prefix" = "p",
    "dynamic_partition.start" = "-2147483648",
    "dynamic_partition.time_unit" = "DAY",
    "dynamic_partition.time_zone" = "Asia/Shanghai",
    "fast_schema_evolution" = "false",
    "replicated_storage" = "true",
    "replication_num" = "3"
);