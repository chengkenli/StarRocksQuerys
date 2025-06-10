/*
 *@author  chengkenli
 *@project StarRocksAPIs
 *@package StarRocksAPIs
 *@file    main
 *@date    2025/2/12 14:02
 */

package main

import (
	"StarRocksAPIs/app"
	_ "StarRocksAPIs/init"
	"StarRocksAPIs/util"
)

func main() {
	util.Init()
	app.App()
}
