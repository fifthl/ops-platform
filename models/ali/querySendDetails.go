package aliModel

import (
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

/*
短信发送详情
接口: 获取指定号码的发送记录
*/

func QuerySendDetails(PhoneNumber, SendDate, AccountID string) *[]SmsSendDetailDTO {
	sends := new([]SmsSendDetailDTO)

	account := utilModel.Key[AccountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	// 创建客户端
	client, _ := SmsCreateClient(tea.String(id), tea.String(secret))

	querySendDetailsRequest := &dysmsapi20170525.QuerySendDetailsRequest{
		PhoneNumber: tea.String(PhoneNumber),
		SendDate:    tea.String(SendDate),
		PageSize:    tea.Int64(50),
		CurrentPage: tea.Int64(1),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.QuerySendDetailsWithOptions(querySendDetailsRequest, runtime)
	if err != nil {
		log.Println("请求失败: ", err)
		return nil
	}

	// 循环拿出SmsSendDetailDTOs的值，调用strconvStatus更改值

	for _, SendDetail := range resp.Body.SmsSendDetailDTOs.SmsSendDetailDTO {
		smsd := new(SmsSendDetailDTO)

		smsd.Content = *SendDetail.Content
		smsd.ErrCode = *SendDetail.ErrCode
		smsd.PhoneNum = *SendDetail.PhoneNum
		smsd.SendDate = *SendDetail.SendDate
		smsd.TemplateCode = *SendDetail.TemplateCode

		smsd.strconvStatus(*SendDetail.SendStatus, *SendDetail.ErrCode) // SendStatus和ErrCode 赋值

		//fmt.Println("smsd: ", smsd)
		*sends = append(*sends, *smsd)
	}
	return sends

}

type SmsSendDetailDTO struct {
	ErrCode      string `json:"ErrCode,omitempty"`
	TemplateCode string `json:"TemplateCode,omitempty"`
	SendDate     string `json:"SendDate,omitempty"`
	PhoneNum     string `json:"PhoneNum,omitempty"`
	Content      string `json:"Content,omitempty"`
	SendStatus   string `json:"SendStatus,omitempty"`
}

func (s *SmsSendDetailDTO) strconvStatus(SendStatus int64, ErrCode string) {

	switch SendStatus {
	case int64(1):
		s.SendStatus = "等待回执"
	case int64(2):
		s.SendStatus = "发送失败"
	case SendStatus:
		s.SendStatus = "发送成功"
	}

	switch ErrCode {
	case "DELIVERED":
		s.ErrCode = "短信发送成功"
	}
}
