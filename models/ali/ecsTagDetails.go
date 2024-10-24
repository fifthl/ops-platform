package aliModel

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"log"
)

type instancesInfo struct {
	InstanceId       string `json:"InstanceId,omitempty"`
	InstanceName     string `json:"InstanceName,omitempty"`
	HostName         string `json:"HostName,omitempty"`
	OsName           string `json:"OsName,omitempty"`
	OsType           string `json:"OsType,omitempty"`
	Cpu              string `json:"Cpu,omitempty"`
	Memory           string `json:"Memory,omitempty"`
	Gpu              string `json:"Gpu,omitempty"`
	PrimaryIpAddress string `json:"PrimaryIpAddress,omitempty"`
	PublicIpAddress  string `json:"PublicIpAddress"`
	EipAddress       string `json:"EipAddress"`
	InstanceType     string `json:"InstanceType,omitempty"`
	ZoneId           string `json:"ZoneId,omitempty"`
	VpcAttributes    string `json:"VpcAttributes,omitempty"`
	VpcId            string `json:"VpcId,omitempty"`
	Tags             string `json:"Tags,omitempty"`
	RenewalStatus    string `json:"RenewalStatus,omitempty"`
	Price            string `json:"Price,omitempty"`
}

// 根据 tag 批量查询
func TagDetails(Tags string) (info *[]instancesInfo) {

	iis := new([]instancesInfo)

	o := orm.NewOrmUsingDB("ecs")
	var query string

	if Tags == "" {
		query = fmt.Sprintf("select * from view_instances_info where tags like '%s';", "%")
	} else {
		query = fmt.Sprintf("select * from view_instances_info where tags='%s';", Tags)
	}

	if _, err := o.Raw(query).QueryRows(iis); err != nil {
		log.Println(err)
		return nil
	}

	return iis
}
