package aliModel

import (
	"fmt"
	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strconv"
	"time"
	"yw_cloud/models/db"
	utilModel "yw_cloud/models/util"
)

/*
 * @Author: zhangqinbo
 * @Date: 2023-11-20 17:08:16
 * @LastEditTime: 2023-11-20 17:08:16
 * @LastEditors: zhangqinbo
 * @Description: 查询昨日费用
 */

type date struct {
	Cycle string
	Date  string
}

// GetDate 获取查询账单的时间
func getDate() *date {
	// 指定时间格式
	location, _ := time.LoadLocation("Asia/Shanghai")
	beforeDate := time.Now().In(location).AddDate(0, 0, -1).Format("2006-01-02")
	beforeCycle := time.Now().In(location).Format("2006-01")

	return &date{
		Cycle: beforeCycle,
		Date:  beforeDate,
	}
}

type Account struct {
	PaymentAmount float32 `json:"paymentAmount"`
	AccountID     string  `json:"accountID"`
}

func GetAccount() []Account {

	bacAct := []Account{}
	redisAct := Account{}

	// 循环获取所有阿里云账号的昨日费用，将汇总数据追加到 bacAct 中返回
	for _, value := range utilModel.Key {
		// 去 redis 中读取余额
		res, err := db.GetRedis(value.Name + "YesterDay")

		if err != nil {
			fmt.Println("redis 内不存在值")
			// redis中不存在键，调用阿里云接口
			bacAct = append(bacAct, getAliyun(value))

		} else {
			// 将 redis 中的值追加后直接返回
			redisAct.AccountID = value.AccountId
			PaymentAmount, _ := strconv.ParseFloat(res, 32)
			redisAct.PaymentAmount = float32(PaymentAmount)
			bacAct = append(bacAct, redisAct)

		}
	}

	return bacAct
}

// 阿里云昨日费用接口
func getAliyun(v utilModel.IdAndSecret) Account {
	aliAct := Account{}
	t := getDate()

	queryAccountBillRequest := &bssopenapi20171214.QueryAccountBillRequest{
		Granularity:  tea.String("DAILY"),
		BillingCycle: tea.String(t.Cycle),
		BillingDate:  tea.String(t.Date),
	}

	runtime := &util.RuntimeOptions{}

	// 调用阿里云api
	ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
	// 创建客户端
	client, _err := BillCreateClient(&ID, &Secret)

	if _err != nil {
		fmt.Println("err: ", _err)
		//return _err
	}

	// 发起请求
	_res, _ := client.QueryAccountBillWithOptions(queryAccountBillRequest, runtime)

	// 如果PaymentAmount长度为零，则查询日没有账单生成
	if len(_res.Body.Data.Items.Item) == 0 {
		aliAct.PaymentAmount = 0
	} else {
		aliAct.PaymentAmount = *(_res.Body.Data.Items.Item[0].PaymentAmount)
	}

	// 插入redis
	if err := db.RedisDb.Set(v.Name+"YesterDay", aliAct.PaymentAmount, utilModel.DaySeconds()).Err(); err != nil {
		fmt.Println("redis写入失败")
		return Account{}
	}
	return aliAct

}
