package aliModel

import (
	"fmt"
	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"strconv"
	"time"
	"yw_cloud/models/db"
	utilModel "yw_cloud/models/util"
)

/*
 * @Author: zhangqinbo
 * @Date: 2023-11-22 15:43:16
 * @LastEditTime: 2023-11-22 15:43:16
 * @LastEditors: zhangqinbo
 * @Description: 查询上月费用
 */

type dateMonthly struct {
	Cycle string
}

// GetDate 获取查询账单的时间
func getDateMonthly() *dateMonthly {
	// 指定时间格式
	location, _ := time.LoadLocation("Asia/Shanghai")
	beforeCycle := time.Now().In(location).AddDate(0, -1, 0).Format("2006-01")

	return &dateMonthly{
		Cycle: beforeCycle,
	}
}

type MonthlyBilling struct {
	PaymentAmount float32 `json:"paymentAmount"`
	AccountID     string  `json:"accountID"`
}

func GetMonthlyBilling() []MonthlyBilling {
	bacMb := MonthlyBilling{}
	bacMbS := []MonthlyBilling{}

	for _, value := range utilModel.Key {
		res, err := db.GetRedis(value.Name + "MonthlyBilling")
		if err != nil {
			// 不为空证明 redis 中没有键，则调用阿里云
			bacMbS = append(bacMbS, getAliyunMonthly(value))
		} else {
			// 否则证明在 redis 中取到值
			bacMb.AccountID = value.AccountId

			// 类型转换
			PaymentAmount, _ := strconv.ParseFloat(res, 32)
			bacMb.PaymentAmount = float32(PaymentAmount)

			//追加
			bacMbS = append(bacMbS, bacMb)
		}
	}

	return bacMbS
}

func getAliyunMonthly(v utilModel.IdAndSecret) MonthlyBilling {
	t := getDateMonthly()
	aliyunMb := MonthlyBilling{}

	queryAccountBillRequest := &bssopenapi20171214.QueryAccountBillRequest{
		Granularity:  tea.String("MONTHLY"),
		BillingCycle: tea.String(t.Cycle),
	}

	runtime := &util.RuntimeOptions{}

	// 调用阿里云api
	ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
	// 创建客户端
	client, _err := BillCreateClient(&ID, &Secret)
	if _err != nil {
		log.Println(_err)
	}

	_res, _ := client.QueryAccountBillWithOptions(queryAccountBillRequest, runtime)

	// 如果PaymentAmount长度为零，则查询日没有账单生成
	if len(_res.Body.Data.Items.Item) == 0 {
		aliyunMb.PaymentAmount = 0
	} else {
		aliyunMb.PaymentAmount = *(_res.Body.Data.Items.Item[0].PaymentAmount)
	}

	aliyunMb.AccountID = v.AccountId

	// 插入redis
	if err := db.RedisDb.Set(v.Name+"MonthlyBilling", aliyunMb.PaymentAmount, utilModel.TimeUntilNextWeekday()).Err(); err != nil {
		fmt.Println("redis写入失败")
		return MonthlyBilling{}
	}

	return aliyunMb

}
