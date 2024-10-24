package aliModel

import (
	"github.com/beego/beego/v2/client/orm"
	"log"
)

/*
函数状态查询
*/

func FunctionUnUse() (interface{}, error) {
	var (
		o             = orm.NewOrm()
		queryUnUseSql = "SELECT sd_id,sd_user, sd_passwd FROM sd_web WHERE sd_id IN (SELECT function_short FROM sd_app WHERE stat = '0' GROUP BY function_short);"
	)
	UnUseresult := []struct {
		SdId     string `json:"sdId,omitempty"`
		SdUser   string `json:"sdUser,omitempty"`
		SdPasswd string `json:"sdPasswd,omitempty"`
	}{}

	if _, err := o.Raw(queryUnUseSql).QueryRows(&UnUseresult); err != nil {
		log.Printf("查询未使用应用出错: %v \n", err)
		return nil, err
	}

	return &UnUseresult, nil
}

func FunctionUse() (interface{}, error) {
	var (
		o           = orm.NewOrm()
		queryUseSql = "SELECT sd_id,sd_user, sd_passwd FROM sd_web WHERE sd_id IN (SELECT function_short FROM sd_app WHERE stat = '1' GROUP BY function_short);"
	)

	Useresult := []struct {
		SdId     string `json:"sdId,omitempty"`
		SdUser   string `json:"sdUser,omitempty"`
		SdPasswd string `json:"sdPasswd,omitempty"`
	}{}

	if _, err := o.Raw(queryUseSql).QueryRows(&Useresult); err != nil {
		log.Printf("查询使用应用出错: %v \n", err)
		return nil, err
	}
	return &Useresult, nil
}
