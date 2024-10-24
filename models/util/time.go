/*
 * @Author: nevin
 * @Date: 2023-11-13 16:30:19
 * @LastEditTime: 2023-11-13 16:31:45
 * @LastEditors: nevin
 * @Description: 时间工具
 */
package utilModel

import (
	"fmt"
	"strconv"
	"time"
)

var (
	TimeZone, _ = time.LoadLocation("Asia/Shanghai")
)

// 计算距离下周一
func TimeUntilNextWeekday() time.Duration {

	now1 := time.Now().In(TimeZone)

	//local, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(TimeZone).Format("15:04:05")

	//inputTime := "18:17:00"
	layout := "15:04:05"

	// 解析时间字符串
	parsedTime, err := time.Parse(layout, now)
	if err != nil {
		fmt.Println("解析时间错误:", err)

	}

	// 计算秒数
	seconds := parsedTime.Hour()*3600 + parsedTime.Minute()*60 + parsedTime.Second()

	overSecond := 86400 - seconds
	//fmt.Println("sec: ", overSecond)

	daysUntilNextWeekday := (int(time.Monday) - int(now1.Weekday()) + 6) % 7

	//fmt.Println("daysUntilNextWeekday: ", (daysUntilNextWeekday*60*60*24)+overSecond+3600)

	return time.Duration((daysUntilNextWeekday*60*60*24)+overSecond+3600) * time.Second
	//return time.Duration(daysUntilNextWeekday)*24*time.Hour + time.Duration(24-intHours+3)*time.Hour
}

func DaySeconds() time.Duration {
	now := time.Now().In(TimeZone).Format("15")
	intNow, _ := strconv.Atoi(now)
	second := (24 - intNow + 1) * 3600
	return time.Duration(second) * time.Second

}
