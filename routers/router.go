/*
 * @Author: nevin
 * @Date: 2024-1-11 15:27:26
 * @LastEditTime: 2024-1-11 15:27:26
 * @LastEditors: nevin
 * @Description: 路由
 */
// @APIVersion 1.0.0
// @Title 云资产api
// @Description 阿里云等对接的api
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"yw_cloud/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/ali",
			beego.NSRouter("/ecsprice", &controllers.AliController{}, "get:EcsPrice"),
			beego.NSRouter("/taginfo", &controllers.AliController{}, "get:TagDetails"),
			beego.NSRouter("/score", &controllers.AliController{}, "get:GetScore"),
			beego.NSRouter("/vul", &controllers.AliController{}, "get:GetVul"),

			// 短信
			beego.NSNamespace("/sms",
				beego.NSRouter("/list", &controllers.AliController{}, "get:GetSmsTemplate"),
				beego.NSRouter("/add", &controllers.AliController{}, "post:AddSmsTemplate"),
				beego.NSRouter("/sendlist", &controllers.AliController{}, "get:SendStatistics"),
				beego.NSRouter("/sendinfo", &controllers.AliController{}, "get:QuerySendDetails"),
				beego.NSRouter("/modify", &controllers.AliController{}, "put:ModifySmsTemplate"),
				beego.NSRouter("/sign", &controllers.AliController{}, "get:QuerySmsSignList"),
			),

			// cdn
			beego.NSNamespace("/cdn",
				beego.NSRouter("/domain", &controllers.AliController{}, "get:GetCdnDomains"),
				beego.NSRouter("/preheat", &controllers.AliController{}, "put:PushObjectCache"),
				beego.NSRouter("/refresh", &controllers.AliController{}, "put:RefreshObjectCaches"),
				beego.NSRouter("/refreshtask", &controllers.AliController{}, "get:DescribeRefreshTasks"),
				beego.NSRouter("/quota", &controllers.AliController{}, "get:DescribeRefreshQuota"),
				beego.NSRouter("/refredomain", &controllers.AliController{}, "get:RefreshCdnDomain"),
			),

			// 负载均衡
			beego.NSNamespace("/slb",
				beego.NSRouter("/update", &controllers.AliController{}, "get:DescribeLoadBalancers"),
				beego.NSRouter("/listen", &controllers.AliController{}, "get:DescribeLoadBalancerListeners"),
				beego.NSRouter("/backend", &controllers.AliController{}, "get:DescribeLoadBalancerAttribute"),
				beego.NSRouter("/list", &controllers.AliController{}, "get:GetSlbList"),
			),

			beego.NSNamespace("/oss",
				beego.NSRouter("/list", &controllers.AliController{}, "get:GetOSS"),
			),

			// 弹性公网
			beego.NSNamespace("/eip",
				beego.NSRouter("/list", &controllers.AliController{}, "get:DescribeEipAddresses"),
			),

			// 域名
			beego.NSNamespace("/domain",
				beego.NSRouter("/list", &controllers.AliController{}, "get:QueryDomainList"),
			),

			// ocr(表格识别)
			beego.NSNamespace("/ocr",
				beego.NSRouter("/recognize", &controllers.AliController{}, "get:TableOcr")),

			// 原始几个项目
			beego.NSInclude(
				&controllers.AliController{},
			),
		),
		//账号续费
		beego.NSNamespace("/renewal",
			beego.NSInclude(
				&controllers.RenewalController{},
			),
		),
		// jenkins 发版记录
		beego.NSNamespace("/ci",
			//beego.NSRouter("/savedb", &controllers.JenkinsController{}, "post:NewSaveDb"),
			beego.NSRouter("/savedb1", &controllers.JenkinsController{}, "post:NewSaveDb"),
			beego.NSRouter("/history", &controllers.JenkinsController{}, "get:GetDB"),
		),

		beego.NSNamespace("/sonar",
			beego.NSRouter("/info", &controllers.SonarController{}, "put:SonarCollect"),
			beego.NSRouter("/list", &controllers.SonarController{}, "get:GetCollect"),
		),

		beego.NSNamespace("/fc",
			beego.NSRouter("/binding", &controllers.AliController{}, "get:BindingNas"),
			beego.NSRouter("/list", &controllers.AliController{}, "get:ListFunctions"),
			beego.NSRouter("/unbind", &controllers.AliController{}, "get:UnBind"),
			beego.NSRouter("/bindecs", &controllers.AliController{}, "get:HttpBindDir"),
			beego.NSRouter("/clear", &controllers.AliController{}, "get:Clear"),
			//beego.NSRouter("/queryfc", &controllers.AliController{}, "get:QueryFcStat"),
			beego.NSRouter("/querybpm", &controllers.AliController{}, "get:QueryBpmId"),
			beego.NSRouter("/copymodel", &controllers.AliController{}, "get:CopyModel"),
			beego.NSRouter("/copylora", &controllers.AliController{}, "get:CopyLoraModel"),
			beego.NSRouter("/modellist", &controllers.AliController{}, "get:ModelList"),
			beego.NSRouter("/amount", &controllers.AliController{}, "get:NasAmount"),
			beego.NSRouter("/use", &controllers.AliController{}, "get:FunctionUse"),
			beego.NSRouter("/unuse", &controllers.AliController{}, "get:FunctionUnUse"),
			beego.NSRouter("/lora", &controllers.AliController{}, "get:QueryLora"),
			beego.NSRouter("/cplora", &controllers.AliController{}, "get:CopyLora"),
		),
	)
	beego.AddNamespace(ns)
}
