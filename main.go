/*
 * @Author: nevin
 * @Date: 2023-11-13 11:14:39
 * @LastEditTime: 2024-11-13 14:11:39
 * @LastEditors: nevin
 * @Description: 入口
 */
package main

import (
	cronAliyun "yw_cloud/models/CronAliyun"
	//_ "yw_cloud/models/ali"
	_ "yw_cloud/models/db"
	_ "yw_cloud/routers"

	"github.com/astaxie/beego"
)

func main() {
	// 启动定时任务
	// 这里后期改成 goroutine，避免启动需要执行完定时才能往下走
	cronAliyun.CronWeeklyAccount()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()

}
