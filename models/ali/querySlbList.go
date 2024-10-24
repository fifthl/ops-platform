package aliModel

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"log"
)

type responseSlbList struct {
	SlbID                   string `json:"SlbID,omitempty" orm:"column(slb_id)"`
	Address                 string `json:"Address,omitempty"`
	LoadBalancerName        string `json:"LoadBalancerName,omitempty"`
	LoadBalancerStatus      string `json:"LoadBalancerStatus,omitempty"`
	Bandwidth               string `json:"Bandwidth,omitempty"`
	InternetChargeTypeAlias string `json:"InternetChargeTypeAlias,omitempty"`
	PayType                 string `json:"PayType,omitempty"`
}

func GetSlbList(accountID string) (SlbList *[]responseSlbList) {
	o := orm.NewOrm()
	sl := new([]responseSlbList)

	query := fmt.Sprintf("select slb_id,address,load_balancer_name,load_balancer_status,bandwidth,internet_charge_type_alias,pay_type from slb_list where account_id = '%s'", accountID)

	if _, err := o.Raw(query).QueryRows(sl); err != nil {
		log.Println(err)
		return nil
	}

	return sl
}
