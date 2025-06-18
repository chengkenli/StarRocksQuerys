/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package analysis
 *@file    Analysis_Schema
 *@date    2025/2/12 14:10
 */

package meta

import (
    "StarRocksQuerys/conn"
    "StarRocksQuerys/tools"
    "StarRocksQuerys/util"
    "fmt"
    "regexp"
    "strings"
    "sync"
    "time"
)

func SessionExtract1(stmt string) []string {
    data, err := sessionExtractSQL(stmt)
    if err != nil {
        util.Loggrs.Error(err.Error())
        return nil
    }
    return data
}
func SessionExtract2(app, stmt string) []util.ExteData {
    data, err := sessionExtractSQL(stmt)
    if err != nil {
        util.Loggrs.Error(err.Error())
        return nil
    }
    db, err := conn.StarRocks(app)
    if err != nil {
        util.Loggrs.Error(err.Error())
        return nil
    }
    for _, d := range data {
        util.Loggrs.Info(fmt.Sprintf("%-5s %-4s", "", d))
    }

    var exterData []util.ExteData
    var wg sync.WaitGroup
    for i, i2 := range data {
        wg.Add(1)
        go func(i int, i2 string) {
            defer wg.Done()
            schema := strings.Split(i2, ".")
            //提取创建日期
            var m map[string]interface{}
            sql := fmt.Sprintf("select `TABLE_CATALOG`,`TABLE_SCHEMA`,`TABLE_NAME`,`TABLE_TYPE`,`ENGINE`,`CREATE_TIME`,`TABLE_COMMENT`,`TABLE_ROWS` from information_schema.tables where TABLE_SCHEMA='%s' and TABLE_NAME='%s'", schema[0], schema[1])
            r := db.Raw(sql).Scan(&m)
            if r.Error != nil {
                util.Loggrs.Warn(r.Error.Error())
                return
            }
            if m == nil {
                return
            }
            var maxDataSize float64
            var mbucket, msize, mreplicaCount, mrowCount string
            if m["TABLE_TYPE"].(string) == "BASE TABLE" {
                //提取表总容量和行数
                var cm []map[string]interface{}
                r = db.Raw(fmt.Sprintf("show data from %s.%s", schema[0], schema[1])).Scan(&cm)
                if r.Error != nil {
                    util.Loggrs.Warn(r.Error.Error())
                    return
                }
                msize = cm[0]["Size"].(string)
                mreplicaCount = cm[0]["ReplicaCount"].(string)
                mrowCount = cm[0]["RowCount"].(string)
                //提取分桶
                var mb map[string]interface{}
                r = db.Raw(fmt.Sprintf("show partitions from %s.%s order by LastConsistencyCheckTime,DataSize desc limit 1", schema[0], schema[1])).Scan(&mb)
                if r.Error != nil {
                    util.Loggrs.Warn(r.Error.Error())
                    return
                }
                mbucket = mb["Buckets"].(string)
                //提取最大的分区容量
                var n []map[string]interface{}
                r = db.Raw(fmt.Sprintf("show partitions from %s.%s", schema[0], schema[1])).Scan(&n)
                if r.Error != nil {
                    return
                }
                var Max []float64
                for _, m := range n {
                    Max = append(Max, tools.Size(m["DataSize"].(string)))
                }
                maxDataSize = Max[0]
                for i := 0; i < len(Max); i++ {
                    if Max[i] >= maxDataSize {
                        maxDataSize = Max[i]
                    }
                }
            }

            exterData = append(exterData, util.ExteData{
                Ctime:     m["CREATE_TIME"].(time.Time).Format("2006-01-02 15:04:05"),
                Bucket:    mbucket,
                Size2gb:   fmt.Sprintf("%0.2fGB", maxDataSize/1024/1024/1024),
                Size2v:    msize,
                Replica:   mreplicaCount,
                Rowcount:  mrowCount,
                Tablename: i2,
            })

        }(i, i2)
    }
    wg.Wait()
    return exterData
}

func sessionExtractSQL(stmt string) ([]string, error) {
    database := []string{
        "adhoc", "ads", "ads_dev", "ads_dev_secure", "ads_rt", "ads_rt_dev", "ads_rt_dev_secure", "ads_rt_secure", "ads_secure",
        "algo", "algo_dev", "ap_secure", "audit", "bi_item", "bi_realty", "bi_realty_secure", "bi_sams_secure", "bi_scm", "bi_sc_secure",
        "cdp", "cdp_api", "cloud_fcst_dm", "cn_backup_secure", "cn_chilled_data", "cn_core_dim_vm", "cn_di_data", "cn_ec_bi_secure",
        "cn_ec_wmdj_user_action", "cn_mdse_dm_dl_tables", "cn_po_home_system", "cn_po_home_system_dev", "cn_pricing_dl_tables",
        "cn_sams_dl_secure", "cn_wc_highsecure", "cn_wc_mb_secure", "cn_wc_mb_vm", "cn_wc_repl_vm", "cn_wc_vm", "cn_wid_dl_secure",
        "cn_wm_mb_secure", "cn_wm_mb_vm", "cn_wm_repl_vm", "cn_wm_vm", "data_test", "demo", "dim", "dim_dev", "dim_dev_secure", "dim_rt",
        "dim_rt_dev", "dim_rt_dev_secure", "dim_rt_secure", "dim_secure", "dm", "dm_dev", "dm_dev_secure", "dm_secure", "dw", "dwd", "dwd_dev",
        "dwd_dev_secure", "dw_dev", "dw_dev_secure", "dwd_rt", "dwd_rt_dev", "dwd_rt_dev_secure", "dwd_rt_secure", "dwd_secure", "dw_rt",
        "dws", "dws_dev", "dws_dev_secure", "dw_secure", "dws_rt", "dws_rt_dev", "dws_rt_dev_secure", "dws_rt_secure", "dws_secure",
        "euclid_scn_forecast_prod", "finance_kettle", "fin_sox", "fin_sox_dev", "flash_report", "flash_report_dev", "flash_report_sit",
        "hyper_bi_secure", "hyper_ec_secure", "hyper_mdse_dm_secure", "information_schema", "ma_test", "mbrship_secure", "mcfc_report",
        "mcfc_report_dev", "o2o_datacubes_secure", "ods", "ods_app_dev", "ods_app_dev_secure", "ods_app_test", "ods_app_test_secure",
        "ods_archive", "ods_dev", "ods_dev_secure", "ods_gray", "ods_migration_td_gray", "ods_rt", "ods_rt_dev", "ods_rt_dev_secure",
        "ods_rt_secure", "ods_secure", "ods_secure_rt", "ods_sox", "ods_sox_app_dev", "ods_sox_app_test", "ods_sox_dev", "ods_sox_test",
        "ods_test", "ods_test_secure", "ops", "pro_dgtmkt_data", "pro_scct_dev", "sams_finance", "scct_inv_monitor", "scct_logis", "scct_logis_dev",
        "scm_dcqe_secure", "scm_network_secure", "scm_secure", "scm_uihealth", "starrocks_monitor", "_statistics_", "supply_kettle", "svccn_logis",
        "svccn_logis_query", "svcdordgtmkt", "sys", "wm_ad_hoc", "wm_cn_util", "wm_common_vm", "ww_core_dim_vm"}

    if strings.Contains(stmt, "hadoop") && !strings.Contains(strings.ToLower(stmt), "outfile") {
        return nil, nil
    }
    var schema []string
    /*schema.table*/
    re := regexp.MustCompile(`([a-zA-Z][^\s=,'.]+)\.([^\s=,'.]+)`)
    var result []string
    for _, s := range re.FindAllString(stmt, -1) {
        b := regexp.MustCompile(`[\\/\(\),:|+><~!@#%^&*='";?-]`).FindString(s) != ""
        if !b {
            data := strings.Split(s, ".")
            for _, s2 := range database {
                if data[0] == s2 {
                    result = append(result, s)
                }
            }
        }
    }
    /*catalog.schema.table*/
    re2 := regexp.MustCompile(`([a-zA-Z][^\s=,'.]+)\.([^\s=,'.]+)\.([^\s=,'.]+)`)
    var result2 []string
    for _, s := range re2.FindAllString(stmt, -1) {
        b := regexp.MustCompile(`[\\/\(\),:|+><~!@#%^&*='";?-]`).FindString(s) != ""
        if !b {
            data := strings.Split(s, ".")
            for _, s2 := range database {
                if data[1] == s2 {
                    result2 = append(result2, s)
                }
            }
        }
    }

    schema = append(schema, result...)
    schema = append(schema, result2...)

    return tools.RmDuplicaSlice(schema), nil
}
