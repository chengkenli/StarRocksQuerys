/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package meta
 *@file    meta_shortquery
 *@date    2025/6/5 15:39
 */

package meta

import (
	"StarRocksAPIs/util"
	"gorm.io/gorm"
	"regexp"
)

func ShortQuery(db *gorm.DB) []string {
	//SHOW RESOURCE GROUP SHORT_QUERY;
	var s []util.ShortData
	r := db.Raw("SHOW RESOURCE GROUP SHORT_QUERY").Scan(&s)
	if r.Error != nil {
		return nil
	}
	if s == nil {
		return nil
	}
	var a []string
	for _, item := range s {
		user, _ := regex(item.Classifiers)
		a = append(a, user)
	}
	return a
}

func regex(str string) (string, string) {
	// 正则表达式匹配 user 和 query_type
	re := regexp.MustCompile(`user=([a-zA-Z0-9_]+), query_type in \(([^)]+)\)`)
	matches := re.FindStringSubmatch(str)

	var user, queryType string
	if len(matches) > 2 {
		user = matches[1]
		queryType = matches[2]
	}
	return user, queryType
}
