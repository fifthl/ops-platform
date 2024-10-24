package aliModel

// 查看监听列表

import (
	"encoding/json"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

type listener struct {
	BackendServerPort int    `json:"BackendServerPort"` // 后端监听端口
	Bandwidth         int    `json:"Bandwidth"`         // 监听的带宽峰值
	Description       string `json:"Description"`       // 监听描述
	ListenerPort      int    `json:"ListenerPort"`      // 监听端口
	ListenerProtocol  string `json:"ListenerProtocol"`  // 监听协议
	LoadBalancerId    string `json:"LoadBalancerId"`    // 负载均衡实例的 ID
	Status            string `json:"Status"`            // 监听的状态
}

// 请求阿里云获取监听信息
func reqListenList(Id, Secret, LoadBalancerId *string) (lsListen []*listener) {
	l := new(listener)

	client, _ := SlbCreateClient(Id, Secret)

	describeLoadBalancerListenersRequest := &slb20140515.DescribeLoadBalancerListenersRequest{
		RegionId:       tea.String("cn-beijing"),
		LoadBalancerId: []*string{tea.String(*LoadBalancerId)},
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.DescribeLoadBalancerListenersWithOptions(describeLoadBalancerListenersRequest, runtime)
	if err != nil {
		log.Println(err)
		return
	}

	for _, listen := range resp.Body.Listeners {
		if err = json.Unmarshal([]byte(listen.String()), l); err != nil {
			log.Println(err)
			continue
		}
		lsListen = append(lsListen, l)
	}
	return lsListen

}

func DescribeLoadBalancerListeners(AccountId, LoadBalancerId string) (lsbListen []*listener) {
	account := utilModel.Key[AccountId]

	// 解密
	id, secret := IdSecret(account.ID, account.Secret)

	// 请求
	lsbListen = reqListenList(&id, &secret, &LoadBalancerId)

	return lsbListen

}
