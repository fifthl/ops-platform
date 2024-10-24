package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["yw_cloud/controllers:AliController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:AliController"],
		beego.ControllerComments{
			Method:           "GetAccount",
			Router:           "/account",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:AliController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:AliController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/all",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:AliController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:AliController"],
		beego.ControllerComments{
			Method:           "GetConsume",
			Router:           "/consume",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:AliController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:AliController"],
		beego.ControllerComments{
			Method:           "GetMonthlyBilling",
			Router:           "/monthlybill",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:AliController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:AliController"],
		beego.ControllerComments{
			Method:           "GetBillOverview",
			Router:           "/overviewbill",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:AliController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:AliController"],
		beego.ControllerComments{
			Method:           "GetWeeklyBilling",
			Router:           "/weeklybill",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"],
		beego.ControllerComments{
			Method:           "CreateRenewal",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"],
		beego.ControllerComments{
			Method:           "DeleteRenewal",
			Router:           "/",
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"],
		beego.ControllerComments{
			Method:           "GetAllHistoryRenewals",
			Router:           `/:`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	//阿里云历史续费信息
	beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"] = append(beego.GlobalControllerRouter["yw_cloud/controllers:RenewalController"],
		beego.ControllerComments{
			Method:           "GetAliyunHistoryRenewals",
			Router:           `/aliyun/:`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
