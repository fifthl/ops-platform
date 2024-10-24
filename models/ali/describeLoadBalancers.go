// 查询负载均衡列表
package aliModel

import (
	"encoding/json"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/beego/beego/v2/client/orm"
	"log"
	utilModel "yw_cloud/models/util"
)

func init() {
	orm.RegisterModel(new(loadBalancer))
}

// 对加密过的 AccessKeyId AccessKeySecret 解密
func IdSecret(idEnc, secretEnc string) (id, secret string) {
	id, secret = utilModel.Decrypt(idEnc, secretEnc)
	return id, secret

}

// 请求阿里云
func reqSlbList(id, secret *string) (slbs []loadBalancer) {

	client, _ := SlbCreateClient(id, secret)

	describeLoadBalancersRequest := &slb20140515.DescribeLoadBalancersRequest{
		RegionId: tea.String("cn-beijing"),
		PageSize: tea.Int32(100),
	}
	runtime := &util.RuntimeOptions{}

	res, err := client.DescribeLoadBalancersWithOptions(describeLoadBalancersRequest, runtime)
	if err != nil {
		log.Println(err)
		return slbs
	}

	for index, LoadBalancer := range res.Body.LoadBalancers.LoadBalancer {
		s := new(loadBalancer)
		if err = json.Unmarshal([]byte(LoadBalancer.String()), s); err != nil {
			log.Println(err, index)
			continue
		}
		slbs = append(slbs, *s)
	}

	return slbs

}

// 获取负载均衡列表
func DescribeLoadBalancers(AccountID string) (result string, err error) {

	account := utilModel.Key[AccountID]
	// 解密
	id, secret := IdSecret(account.ID, account.Secret)

	// 请求
	LoadBalancers := reqSlbList(&id, &secret)

	// 循环拿出插入
	o := orm.NewOrm()
	//orm.RunSyncdb("default", false, false)

	for _, s := range LoadBalancers {
		s.AccountID = AccountID
		s.changValue()
		if _, err = o.InsertOrUpdate(&s); err != nil {
			log.Println("slb 信息更新 db 失败: ", err)
			continue
		}
	}

	return "更新成功", nil

}

type loadBalancer struct {
	LoadBalancerId          string `json:"LoadBalancerId,omitempty" orm:"column(slb_id);unique;pk;description(实例ID)"`
	Address                 string `json:"Address,omitempty" orm:"description(IP地址)"`
	LoadBalancerName        string `json:"LoadBalancerName,omitempty" orm:"description(实例名称)"`
	LoadBalancerSpec        string `json:"LoadBalancerSpec,omitempty" orm:"description(性能规格)"`
	LoadBalancerStatus      string `json:"LoadBalancerStatus,omitempty" orm:"description(实例状态)"` // 实例状态
	Bandwidth               int32  `json:"Bandwidth,omitempty" orm:"description(带宽峰值)"`
	InternetChargeType      string `json:"InternetChargeType,omitempty" orm:"description(公网类型实例付费方式)"`  // 计费方式
	InternetChargeTypeAlias string `json:"InternetChargeTypeAlias,omitempty" orm:"description(公网计费方式)"` // 公网计费方式
	MasterZoneId            string `json:"MasterZoneId,omitempty" orm:"description(主可用区ID)"`            // 地域
	SlaveZoneId             string `json:"SlaveZoneId,omitempty" orm:"description(备可用区ID)"`             // 备可用区
	VpcId                   string `json:"VpcId,omitempty" orm:"description(VPC ID)"`                   // 专有网络 ID
	AddressIPVersion        string `json:"AddressIPVersion,omitempty" orm:"description(IPV4/IPV6)"`
	AddressType             string `json:"AddressType,omitempty" orm:"description(公网负载均衡/内网负载均衡)"`
	CreateTime              string `json:"CreateTime,omitempty" orm:"description(实例创建时间)"`
	DeleteProtection        string `json:"DeleteProtection,omitempty" orm:"description(是否删除保护状态)"`
	NetworkType             string `json:"NetworkType,omitempty" orm:"description(私网负载均衡实例网络类型)"` // 私网类型
	PayType                 string `json:"PayType,omitempty" orm:"description(付费模式)"`             // 付费模式
	VSwitchId               string `json:"VSwitchId,omitempty" orm:"description(私网负载均衡实例的交换机ID)"` // 交换机 ID
	AccountID               string `json:"AccountID" orm:"column(account_id)"`
}

func (*loadBalancer) TableName() string {
	return "slb_list"
}

func (r *loadBalancer) changValue() {
	switch r.PayType {
	case "PayOnDemand":
		r.PayType = "按量付费"
	case "PrePay":
		r.PayType = "包年包月"
	}

	switch r.InternetChargeTypeAlias {
	case "paybybandwidth":
		r.InternetChargeTypeAlias = "按带宽计费"
	case "paybytraffic":
		r.InternetChargeTypeAlias = "按流量计费"
	}

	switch r.LoadBalancerStatus {
	case "inactive":
		r.LoadBalancerStatus = "实例已停止"
	case "active":
		r.LoadBalancerStatus = "实例运行中"
	case "locked":
		r.LoadBalancerStatus = "实例已锁定"
	}

	switch r.InternetChargeType {
	case "4":
		r.InternetChargeType = "按带宽计费"
	case "3":
		r.InternetChargeType = "按流量计费"
	}

	switch r.AddressType {
	case "internet":
		r.AddressType = "公网负载均衡"
	case "intranet":
		r.AddressType = "公网负载均衡"
	}
}
