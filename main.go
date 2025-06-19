/*
 *@author  chengkenli
 *@project StarRocksQuerys
 *@package StarRocksQuerys
 *@file    main
 *@date    2025/2/12 14:02
 */

package main

import (
	"StarRocksQuerys/app"
	_ "StarRocksQuerys/init"
	"StarRocksQuerys/util"
)

func main() {
	util.Init()
	app.App()
}
