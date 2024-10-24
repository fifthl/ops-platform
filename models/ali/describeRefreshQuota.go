package aliModel

import (
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"log"
	utilModel "yw_cloud/models/util"
)

/*
查询 cdn 剩余刷新次数
*/

type RefreshQuota struct {
	UrlRemain string `json:"UrlRemain,omitempty"`
	DirRemain string `json:"DirRemain,omitempty"`
}

func DescribeRefreshQuota(AccountID string) (quota *RefreshQuota) {
	quota = new(RefreshQuota)

	account := utilModel.Key[AccountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := CdnCreateClient(&id, &secret)

	describeRefreshQuotaRequest := &cdn20180510.DescribeRefreshQuotaRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := client.DescribeRefreshQuotaWithOptions(describeRefreshQuotaRequest, runtime)
	if err != nil {
		log.Println("查询可刷新量失败", err)
	}

	quota.UrlRemain = *resp.Body.UrlRemain
	quota.DirRemain = *resp.Body.DirRemain

	return quota
}
