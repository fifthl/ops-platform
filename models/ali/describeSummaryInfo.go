package aliModel

import (
	sas20181203 "github.com/alibabacloud-go/sas-20181203/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

var (
	accounts = [3]string{"1", "2", "3"}
)

type AccountScore struct {
	Account string `json:"Account,omitempty"`
	Score   *int32 `json:"Score,omitempty"`
}

// 账号安全评分
func GetScore() (SecurityScores []AccountScore) {
	as := new(AccountScore)

	for _, account := range accounts {
		id, secret := utilModel.Decrypt(utilModel.Key[account].ID, utilModel.Key[account].Secret)
		client, _ := SasCreateClient(tea.String(id), tea.String(secret))

		describeSummaryInfoRequest := &sas20181203.DescribeSummaryInfoRequest{}
		runtime := &util.RuntimeOptions{}

		res, respErr := client.DescribeSummaryInfoWithOptions(describeSummaryInfoRequest, runtime)
		if respErr != nil {
			log.Println(respErr)
			continue
		}

		as.Score = res.Body.SecurityScore
		as.Account = utilModel.Key[account].AccountId
		SecurityScores = append(SecurityScores, *as)

	}
	return SecurityScores

}
