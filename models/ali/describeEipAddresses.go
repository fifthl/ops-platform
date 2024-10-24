package aliModel

// 弹性公网IP查询
import (
	"encoding/json"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v6/client"
	"log"
	utilModel "yw_cloud/models/util"
)

type eipAddress struct {
	Status                        string `json:"Status,omitempty"`                        //状态
	Netmode                       string `json:"Netmode,omitempty"`                       // 网路类型
	ChargeType                    string `json:"ChargeType,omitempty"`                    // 付费模式
	Description                   string `json:"Description,omitempty"`                   // 描述
	ReservationInternetChargeType string `json:"ReservationInternetChargeType,omitempty"` // 计费类型
	IpAddress                     string `json:"IpAddress,omitempty"`                     // EIP 的 IP 地址
	Bandwidth                     string `json:"Bandwidth,omitempty"`                     // 峰值
	Name                          string `json:"Name,omitempty"`                          // 名称
	InstanceRegionId              string `json:"InstanceRegionId,omitempty"`              // 地域
	InstanceId                    string `json:"InstanceId,omitempty"`                    // 当前绑定的实例 id
	ExpiredTime                   string `json:"ExpiredTime,omitempty"`                   // 到期时间
	AllocationId                  string `json:"AllocationId,omitempty"`                  // EIP 的实例 ID
	ISP                           string `json:"ISP,omitempty"`                           // 线路类型
	AllocationTime                string `json:"AllocationTime,omitempty"`                // EIP 的创建时间
}

func (e *eipAddress) changeValue() {
	switch e.Status {
	case "InUse":
		e.Status = "已分配"
	case "Available":
		e.Status = "可用"
	}

	switch e.Netmode {
	case "public":
		e.Netmode = "公网"
	}

	switch e.ChargeType {
	case "PostPaid":
		e.ChargeType = "按量计费"
	case "PrePaid":
		e.ChargeType = "包年包月"
	}

	switch e.ReservationInternetChargeType {
	case "PayByBandwidth":
		e.ReservationInternetChargeType = "固定带宽计费"
	case "PayByTraffic":
		e.ReservationInternetChargeType = "使用流量计费"
	}
}

func DescribeEipAddresses(AccountID string) (eas *[]eipAddress) {
	account := utilModel.Key[AccountID]

	// 解密
	id, secret := IdSecret(account.ID, account.Secret)

	eas = reqEipAddresses(id, secret)

	return eas

}

func reqEipAddresses(id, secret string) (eas *[]eipAddress) {
	e := new([]eipAddress)

	client, _ := VpcCreateClient(&id, &secret)
	describeEipAddressesRequest := &vpc20160428.DescribeEipAddressesRequest{
		RegionId: tea.String("cn-beijing"),
		PageSize: tea.Int32(100),
	}
	runtime := &util.RuntimeOptions{}

	res, err := client.DescribeEipAddressesWithOptions(describeEipAddressesRequest, runtime)
	if err != nil {
		log.Println(err)
		return
	}

	for _, eip := range res.Body.EipAddresses.EipAddress {
		ea := new(eipAddress)

		if err = json.Unmarshal([]byte(eip.String()), ea); err != nil {
			log.Println(err)
			continue
		}

		// 更改值类型，没有复杂逻辑
		ea.changeValue()

		*e = append(*e, *ea)
	}

	return e
}
