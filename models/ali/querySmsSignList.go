package aliModel

import (
	"encoding/json"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

func QuerySmsSignList(AccountID string) (signList *[]smsSign) {
	s := new([]smsSign)

	account := utilModel.Key[AccountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)
	client, _ := SmsCreateClient(tea.String(id), tea.String(secret))

	querySmsSignListRequest := &dysmsapi20170525.QuerySmsSignListRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := client.QuerySmsSignListWithOptions(querySmsSignListRequest, runtime)
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range resp.Body.SmsSignList {
		sms := new(smsSign)

		if err = json.Unmarshal([]byte(v.GoString()), sms); err != nil {
			log.Printf("查询短信签名序列化失败: %v", err)
			continue
		}

		// 更改 AuditStatus 字段值
		sms.strconv()

		*s = append(*s, *sms)
	}

	return s

}

type smsSign struct {
	SignName    string    `json:"SignName,omitempty"`
	AuditStatus string    `json:"AuditStatus,omitempty"`
	CreateDate  string    `json:"CreateDate,omitempty"`
	Reason      smsReason `json:"Reason"`
}

type smsReason struct {
	RejectSubInfo string `json:"RejectSubInfo,omitempty"`
	RejectDate    string `json:"RejectDate,omitempty"`
	RejectInfo    string `json:"RejectInfo,omitempty"`
}

func (s *smsSign) strconv() {
	switch s.AuditStatus {
	case "AUDIT_STATE_INIT":
		s.AuditStatus = "审核中"
	case "AUDIT_STATE_PASS":
		s.AuditStatus = "审核通过"
	case "AUDIT_STATE_NOT_PASS":
		s.AuditStatus = "审核未通过"
	case "AUDIT_STATE_CANCEL":
		s.AuditStatus = "取消审核"
	}
}
