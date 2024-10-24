package aliModel

import (
	"fmt"
	"strconv"
	//cronAli "yw_cloud/models/CronAliyun"
	"yw_cloud/models/db"
	utilModel "yw_cloud/models/util"
)

/*
 * @Author: zhangqinbo
 * @Date: 2023-11-22 09:52:16
 * @LastEditTime: 2023-11-22 09:52:16
 * @LastEditors: zhangqinbo
 * @Description: 查询上周费用
 */

type weeklyBill struct {
	PaymentAmount float64 `json:"paymentAmount"`
	AccountID     string  `json:"accountID"`
}

// 从表读取周账单
func GetWeeklyBilling() []weeklyBill {
	bacWB := weeklyBill{}
	bacWBS := []weeklyBill{}

	for _, value := range utilModel.Key {
		result, err := db.RedisDb.Get(value.Name + "WeeklyBilling").Result()
		if err != nil {
			fmt.Println("Redis 中未找到 KEY. 获取周账单失败")
		}
		bacWB.AccountID = value.AccountId
		bacWB.PaymentAmount, _ = strconv.ParseFloat(result, 64)

		bacWBS = append(bacWBS, bacWB)
	}

	return bacWBS
}
