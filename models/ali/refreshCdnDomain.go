package aliModel

import (
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/beego/beego/v2/client/orm"
	"log"
	utilModel "yw_cloud/models/util"
)

/*
刷新 cdn 域名
*/
type cdnDomain struct {
	Domain  string `json:"DomainName,omitempty"`
	Status  string `json:"DomainStatus,omitempty"`
	Gmt     string `json:"GmtCreated,omitempty"`
	Account string `json:"AccountID"`
}

func (d *cdnDomain) TableName() string {
	return "cdn_domain"
}

var one = false

func RefreshCdnDomain() (mess string) {

	if one == false {
		orm.RegisterModel(new(cdnDomain))
		one = true
	}
	//_ = orm.RunSyncdb("default", false, false)
	o := orm.NewOrm()

	// 循环账号
	for _, account := range utilModel.Key {
		id, secret := utilModel.Decrypt(account.ID, account.Secret)
		// 循环拿 cdn 域名
		client, _ := CdnCreateClient(tea.String(id), tea.String(secret))

		// 请求参数
		describeUserDomainsRequest := &cdn20180510.DescribeUserDomainsRequest{}
		runtime := &util.RuntimeOptions{}

		resp, _err := client.DescribeUserDomainsWithOptions(describeUserDomainsRequest, runtime)
		if _err != nil {
			log.Println("获取cdn域名失败: ", _err)
			log.Println(account.AccountId)
			continue
		}

		for _, domainInfo := range resp.Body.Domains.PageData {
			domain := cdnDomain{
				Domain: *domainInfo.DomainName,
				Status: *domainInfo.DomainStatus,
				Gmt:    *domainInfo.GmtCreated,
			}
			domain.Account = account.AccountId

			if _, err := o.InsertOrUpdate(&domain); err != nil {
				log.Println("cdn 域名写入 db 失败")
				return "cdn 域名刷新失败"
			}

		}
	}
	return "cdn 域名刷新成功"
}
