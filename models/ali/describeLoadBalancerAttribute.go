// 查询后端服务器信息

/*
目前查询到哪台ECS，但无法查询主机名 ip 等ECS相关信息，需要从其他渠道查询
*/
package aliModel

import (
	"encoding/json"
	"fmt"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/beego/beego/v2/client/orm"
	"log"
	utilModel "yw_cloud/models/util"
)

type LoadBalancerResponse struct {
	BackendServers backendServers
}

type backendServers struct {
	ServerId         string `json:"ServerId,omitempty"`
	ServerName       string `json:"ServerName"`
	Region           string `json:"Region"`
	VpcID            string `json:"VpcID"`
	PrimaryIpAddress string `json:"PrimaryIpAddress"`
	Type             string `json:"Type,omitempty"`
	Weight           int32  `json:"Weight,omitempty"`
}

// 从db查询ecs信息
func serverInfo(ServerId string) (ServerName, Region, VpcID, PrimaryIpAddress string) {
	o := orm.NewOrmUsingDB("ecs")

	var instance []string
	var zone []string
	var vpc []string
	var ip []string

	query := fmt.Sprintf("select instance_name,zone_id,vpc_id,primary_ip_address from instances_info where instance_id='%s'", ServerId)
	if _, err := o.Raw(query).QueryRows(&instance, &zone, &vpc, &ip); err != nil {
		log.Println(err)
	}

	return instance[0], zone[0], vpc[0], ip[0]

}

// 获取后端服务器信息
func reqBackendServers(id, secret, LoadBalancerId string, bs *backendServers) (bss []backendServers, err error) {
	clinet, _ := SlbCreateClient(&id, &secret)

	describeLoadBalancerAttributeRequest := &slb20140515.DescribeLoadBalancerAttributeRequest{
		RegionId:       tea.String("cn-beijing"),
		LoadBalancerId: tea.String(LoadBalancerId),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := clinet.DescribeLoadBalancerAttributeWithOptions(describeLoadBalancerAttributeRequest, runtime)
	if err != nil {
		log.Println(err)
		return
	}

	// slb 可能绑定多个 ecs
	for _, BackendServer := range resp.Body.BackendServers.BackendServer {
		if err = json.Unmarshal([]byte(BackendServer.String()), bs); err != nil {
			log.Println(err)
			return bss, err
		}
		bss = append(bss, *bs)
	}

	return bss, nil
}

func DescribeLoadBalancerAttribute(AccountID, LoadBalancerId string) (BackendServers []backendServers, err error) {
	account := utilModel.Key[AccountID]
	bs := new(backendServers)

	// 解密
	id, secret := IdSecret(account.ID, account.Secret)

	// 请求
	BackendServers, err = reqBackendServers(id, secret, LoadBalancerId, bs)

	// 利用查询到的 ecs id 补充ip，vpc等信息
	for index, server := range BackendServers {
		ServerName, Region, VpcID, PrimaryIpAddress := serverInfo(server.ServerId)

		BackendServers[index].ServerName = ServerName
		BackendServers[index].Region = Region
		BackendServers[index].VpcID = VpcID
		BackendServers[index].PrimaryIpAddress = PrimaryIpAddress
	}

	return BackendServers, err

}
