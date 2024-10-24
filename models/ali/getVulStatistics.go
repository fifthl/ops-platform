package aliModel

//漏洞数量统计

import (
	sas20181203 "github.com/alibabacloud-go/sas-20181203/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"strconv"
	utilModel "yw_cloud/models/util"
)

var (
	// 固定值，每个阿里云账号组id不同
	GroupIDList = []string{"13411399", "13444169", "13458239"}
)

type AccountVul struct {
	Account           string `json:"Account,omitempty"`
	VulAsapSum        *int32 `json:"VulAsapSum,omitempty"`
	RiskInstanceCount *int32 `json:"RiskInstanceCount,omitempty"`
}

// 返回解密后的 ak as
func getIdSecret(account int) (id, secret, accountId string) {
	//因为 GroupIDList 的索引 0 1 2 与 Key 中的不匹配，所以需要加1
	account++
	accountId = strconv.Itoa(account)

	idEnc, secretEnc := utilModel.Key[accountId].ID, utilModel.Key[accountId].Secret
	id, secret = utilModel.Decrypt(idEnc, secretEnc)
	return id, secret, accountId
}

// 请求阿里云获取漏洞数量
func reqVulInfo(id, secret string, index int) (VulAsapSum *int32, err error) {

	client, _ := SasCreateClient(&id, &secret)

	getVulStatisticsRequest := &sas20181203.GetVulStatisticsRequest{

		GroupIdList: tea.String(GroupIDList[index]),
		TypeList:    tea.String("cve,sys,cms,emg,app,sca"),
	}
	runtime := &util.RuntimeOptions{}

	res, err := client.GetVulStatisticsWithOptions(getVulStatisticsRequest, runtime)
	return res.Body.VulAsapSum, err
}

// 获取存在漏洞的服务器
func reqFieldStatistics(id, secret string) (RiskInstanceCount *int32, err error) {
	client, _ := SasCreateClient(&id, &secret)
	describeFieldStatisticsRequest := &sas20181203.DescribeFieldStatisticsRequest{
		RegionId: tea.String("cn-beijing"),
	}
	runtime := &util.RuntimeOptions{}

	res, err := client.DescribeFieldStatisticsWithOptions(describeFieldStatisticsRequest, runtime)
	return res.Body.GroupedFields.RiskInstanceCount, err
}

func GetVul() (avs []AccountVul) {

	av := new(AccountVul)

	// 循环分组 ID
	for account := 0; account < 3; account++ {

		// 返回解密后的 ak as
		id, secret, acccountId := getIdSecret(account)

		// 获取漏洞数量
		VulAsapSum, err := reqVulInfo(id, secret, account)
		if err != nil {
			log.Println(err)
			continue
		}

		// 获取存在风险的资产数量
		RiskInstanceCount, err := reqFieldStatistics(id, secret)
		if err != nil {
			log.Println(err)
			continue
		}

		av.VulAsapSum = VulAsapSum
		av.RiskInstanceCount = RiskInstanceCount
		av.Account = acccountId

		// 追加
		avs = append(avs, *av)
	}
	return avs
}
