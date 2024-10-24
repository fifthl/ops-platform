package aliModel

import (
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

/*
短信发送统计
*/

type TargetList struct {
	TotalCount            int64  `json:"TotalCount"`
	RespondedSuccessCount int64  `json:"RespondedSuccessCount"`
	RespondedFailCount    int64  `json:"RespondedFailCount"`
	SendDate              string `json:"SendDate,omitempty"`
}

type SmsTotal struct {
	Success int64 `json:"Success,omitempty"`
	Fail    int64 `json:"Fail,omitempty"`
	Sum     int64 `json:"Sum,omitempty"`
}

type Total struct {
	TargetList []TargetList `json:"TargetList,omitempty"`
	SmsTotal   SmsTotal     `json:"SmsTotal"`
}

/*
IsGlobe             国内消息/国际消息
StartDate EndDate   开始时间/结束时间   格式: 20231201/20231225
AccountID           账号ID
*/
func SendStatistics(IsGlobe int32, StartDate, EndDate, AccountID string) Total {
	sm := new(SmsTotal)

	var (
		tl  TargetList
		tls []TargetList
		t   Total
	)

	account := utilModel.Key[AccountID]
	Id, Secret := utilModel.Decrypt(account.ID, account.Secret)
	client, err := SmsCreateClient(&Id, &Secret)
	if err != nil {
		log.Println("添加短信client初始化失败:", err)
	}

	querySendStatisticsRequest := &dysmsapi20170525.QuerySendStatisticsRequest{
		IsGlobe:   tea.Int32(IsGlobe),
		StartDate: tea.String(StartDate),
		EndDate:   tea.String(EndDate),
		PageIndex: tea.Int32(int32(1)),
		PageSize:  tea.Int32(int32(50)),
	}
	runtime := &util.RuntimeOptions{}
	_resp, _err := client.QuerySendStatisticsWithOptions(querySendStatisticsRequest, runtime)
	if _err != nil {
		log.Println("获取短信发送详情失败: ", _err)
	}

	// 计算查询时间段内的成功/失败/总条数
	for _, oneSms := range _resp.Body.Data.TargetList {
		sm.Success += *oneSms.RespondedSuccessCount
		sm.Fail += *oneSms.RespondedFailCount
		sm.Sum += *oneSms.TotalCount
	}

	// range获取每天短信
	for _, oneSms := range _resp.Body.Data.TargetList {
		tl.RespondedFailCount = *oneSms.RespondedFailCount
		tl.RespondedSuccessCount = *oneSms.RespondedSuccessCount
		tl.SendDate = *oneSms.SendDate
		tl.TotalCount = *oneSms.TotalCount

		tls = append(tls, tl)
	}

	t.SmsTotal = *sm
	t.TargetList = tls

	//log.Println(_resp.Body.Data.TargetList)
	return t
}
