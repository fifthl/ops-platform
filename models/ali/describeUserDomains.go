package aliModel

import (
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

/*
获取阿里云 cdn 域名列表
*/

type Domain struct {
	DomainName   string `json:"DomainName,omitempty"`
	DomainStatus string `json:"DomainStatus,omitempty"`
	GmtCreated   string `json:"GmtCreated,omitempty"`
}

func GetCdnDomains(AccountID string) []Domain {
	var (
		domains = make([]Domain, 0)
	)

	account := utilModel.Key[AccountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := CdnCreateClient(tea.String(id), tea.String(secret))

	// 请求参数
	describeUserDomainsRequest := &cdn20180510.DescribeUserDomainsRequest{}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.DescribeUserDomainsWithOptions(describeUserDomainsRequest, runtime)
	if _err != nil {
		log.Println("获取cdn域名失败: ", _err)
		log.Println(account.AccountId)
		return nil
	}

	for _, domainInfo := range resp.Body.Domains.PageData {
		domain := &Domain{
			DomainName:   *domainInfo.DomainName,
			DomainStatus: *domainInfo.DomainStatus,
			GmtCreated:   *domainInfo.GmtCreated,
		}

		domains = append(domains, *domain)
	}
	log.Println("获取cdn域名成功")
	return domains
}
