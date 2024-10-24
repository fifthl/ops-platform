package aliModel

// 查询域名列表

import (
	"encoding/json"
	"fmt"
	domain20180129 "github.com/alibabacloud-go/domain-20180129/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

type Data struct {
	DomainName             string `json:"DomainName,omitempty"`
	RegistrationDate       string `json:"RegistrationDate,omitempty"`
	ExpirationDate         string `json:"ExpirationDate,omitempty"`
	DomainStatus           string `json:"DomainStatus,omitempty"`
	ExpirationCurrDateDiff int32  `json:"ExpirationCurrDateDiff,omitempty"`
	DomainAuditStatus      string `json:"DomainAuditStatus,omitempty"`
	Remark                 string `json:"Remark,omitempty"`
}

func (d *Data) changValue() {

	switch d.DomainAuditStatus {
	case "FAILED":
		d.DomainAuditStatus = "实名认证失败"
	case "SUCCEED":
		d.DomainAuditStatus = "实名认证成功"
	case "NONAUDIT":
		d.DomainAuditStatus = "未实名认证"
	case "AUDITING":
		d.DomainAuditStatus = "审核中"
	}

	switch d.DomainStatus {
	case "1":
		d.DomainStatus = "急需续费"
	case "2":
		d.DomainStatus = "急需赎回"
	case "3":
		d.DomainStatus = "正常"
	}
}

func QueryDomainList(accountID string) (domainList []Data) {
	d := new(Data)

	account := utilModel.Key[accountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := DomainCreateClient(&id, &secret)
	queryDomainListRequest := &domain20180129.QueryDomainListRequest{
		PageSize: tea.Int32(100),
		PageNum:  tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}

	res, err := client.QueryDomainListWithOptions(queryDomainListRequest, runtime)
	if err != nil {
		log.Println(err)
		return
	}

	for _, domain := range res.Body.Data.Domain {
		if err = json.Unmarshal([]byte(domain.String()), d); err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(*d)
		d.changValue()
		domainList = append(domainList, *d)
	}

	return domainList
}
