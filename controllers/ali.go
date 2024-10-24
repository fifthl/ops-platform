/*
 * @Author: nevin
 * @Date: 2023-11-13 14:21:29
 * @LastEditTime: 2023-11-15 14:11:13
 * @LastEditors: nevin
 * @Description: 阿里云接口
 */
package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	aliModel "yw_cloud/models/ali"
)

// 检查传入参数是否为空
func (u *AliController) checkParam(args ...interface{}) {
	for _, arg := range args {
		if arg == "" || arg == 0 {
			u.respond(1, "请求失败", "参数不能为空")
		}
	}
}

// Operations about Ali
type AliController struct {
	BaseController
}

// @Title 阿里余额
// @Description 阿里获取余额
// @Success 200 {object} models.RemainingSum
// @router /all [get]
func (u *AliController) GetAll() {
	ye := aliModel.GetRemainingSum()
	u.respond(0, "请求成功", ye)
}

// @Title 获取昨日余额
// @Description 阿里获取昨日消费
// @Success 200 {object} models.GetAccount
// @router /account [get]
func (u *AliController) GetAccount() {
	account := aliModel.GetAccount()
	u.respond(0, "请求成功", account)
}

// @Title 获取上周账单
// @Description 阿里获取上周消费
// @Success 200 {object} models.GetWeeklyBilling
// @router /weeklybill [get]
func (u *AliController) GetWeeklyBilling() {
	resp := aliModel.GetWeeklyBilling()
	u.respond(0, "请求成功", resp)
}

// @Title 获取上月账单
// @Description 阿里获取上月消费
// @Success 200 {object} models.GetMonthlyBilling
// @router /monthlybill [get]
func (u *AliController) GetMonthlyBilling() {
	resp := aliModel.GetMonthlyBilling()
	u.respond(0, "请求成功", resp)
}

// @Title 获取月分摊账单
// @Description 阿里获取月分摊账单
// @Success 200 {object} models.GetBillOverview
// @router /overviewbill [post]
func (u *AliController) GetBillOverview() {

	account := u.GetString("accountid")
	starttime := u.GetString("starttime")
	endtime := u.GetString("endtime")

	u.checkParam(account, starttime, endtime)

	if starttime == "" || endtime == "" {
		u.respond(2, "starttime 或 endtime 不能为空", "")
	}

	result, err := aliModel.BillOverview(account, starttime, endtime)
	if err != nil {
		u.respond(1, "请求失败", result)
	}
	u.respond(0, "请求成功", result)

}

// @Title 阿里云账号消费
// @Description 阿里云账号消费
// @Success 200 {object} models.GetConsume
// @router /consume [get]
func (u *AliController) GetConsume() {
	date := u.GetString("date")
	u.checkParam(date)

	resp := aliModel.GetConsume(date)
	u.respond(0, "请求成功", resp)
}

// ECS tag 费用
func (u *AliController) EcsPrice() {
	resp := aliModel.EcsPrice()
	u.respond(0, "请求成功", resp)
}

// 获取短信模板
func (u *AliController) GetSmsTemplate() {
	accountID := u.GetString("id")
	pageIndex, _ := u.GetInt64("index")
	pageSize, _ := u.GetInt64("size")
	u.checkParam(accountID, pageIndex, pageSize)

	resp := aliModel.GetSmsTemplate(accountID, pageIndex, pageSize)
	u.respond(0, "请求成功", resp)
}

// 添加短信模板
func (u *AliController) AddSmsTemplate() {
	type p struct {
		TemplateType    int32  `json:"TemplateType,omitempty"`
		TemplateName    string `json:"TemplateName,omitempty"`
		TemplateContent string `json:"TemplateContent,omitempty"`
		Remark          string `json:"Remark,omitempty"`
		AccountID       string `json:"AccountID,omitempty"`
	}
	param := new(p)

	if err := json.Unmarshal(u.Ctx.Input.RequestBody, param); err != nil {
		log.Println("结构体转换失败: ", err)
	}

	resp, err := aliModel.AddSmsTemplate(param.TemplateType, param.TemplateName, param.TemplateContent, param.Remark, param.AccountID)

	if err != nil {
		u.respond(1, "请求失败", err.Error())
	} else {
		u.respond(0, "请求成功", resp)
	}

}

// 获取短信发送条目
func (u *AliController) SendStatistics() {
	IsGlobe := u.GetString("globe")
	StartDate := u.GetString("start")
	EndDate := u.GetString("end")
	AccountID := u.GetString("id")

	u.checkParam(IsGlobe, StartDate, EndDate, AccountID)

	IntGlobe, _ := strconv.Atoi(IsGlobe)

	resp := aliModel.SendStatistics(int32(IntGlobe), StartDate, EndDate, AccountID)
	u.respond(0, "请求成功", resp)
}

// 根据接收号码查询短信发送记录
func (u *AliController) QuerySendDetails() {

	PhoneNumber := u.GetString("number")
	SendDate := u.GetString("date")
	AccountID := u.GetString("id")

	u.checkParam(PhoneNumber, SendDate, AccountID)

	resp := aliModel.QuerySendDetails(PhoneNumber, SendDate, AccountID)
	u.respond(0, "请求成功", resp)

}

// 获取 cdn 域名
func (u *AliController) GetCdnDomains() {
	account := u.GetString("id")
	u.checkParam(account)

	resp := aliModel.GetCdnDomains(account)
	u.respond(0, "请求成功", resp)
}

// 预热 URL
func (u *AliController) PushObjectCache() {
	type p struct {
		URL []string `json:"URL"`
	}

	param := new(p)
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, param); err != nil {
		log.Printf("预热序列化结构体失败: %v", err)
	}
	mess := aliModel.PushObjectCache(param.URL)
	u.respond(0, "请求成功", mess)
}

// 缓存刷新
func (u *AliController) RefreshObjectCaches() {
	type p struct {
		URL        []string `json:"URL"`
		RefresType string   `json:"RefresType"`
	}

	param := new(p)
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, param); err != nil {
		log.Printf("刷新缓存序列化结构体失败: %v", err)
	}
	mess := aliModel.RefreshObjectCaches(param.URL, param.RefresType)
	u.respond(0, "请求成功", mess)
}

// 获取cdn刷新/预热记录
func (u *AliController) DescribeRefreshTasks() {

	AccountID := u.GetString("id")
	u.checkParam(AccountID)

	RefreshTask := aliModel.DescribeRefreshTasks(AccountID)
	u.respond(0, "请求成功", RefreshTask)
}

// 查询cdn可刷新量
func (u *AliController) DescribeRefreshQuota() {

	AccountID := u.GetString("id")
	u.checkParam(AccountID)

	RefreshTask := aliModel.DescribeRefreshQuota(AccountID)
	u.respond(0, "请求成功", RefreshTask)
}

// 刷新cdn域名
func (u *AliController) RefreshCdnDomain() {
	resp := aliModel.RefreshCdnDomain()
	u.respond(0, "请求成功", resp)
}

// 更改未成功的短信模板
func (u *AliController) ModifySmsTemplate() {
	type p struct {
		AccountID       string `json:"AccountID,omitempty"`
		TemplateType    int32  `json:"TemplateType,omitempty"`
		TemplateName    string `json:"TemplateName,omitempty"`
		TemplateCode    string `json:"TemplateCode,omitempty"`
		TemplateContent string `json:"TemplateContent,omitempty"`
		Remarks         string `json:"Remarks,omitempty"`
	}

	param := new(p)
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, param); err != nil {
		log.Println("序列化更改短信内容失败: ", err)
		return
	}
	code, msg := aliModel.ModifySmsTemplate(param.AccountID, param.TemplateType, param.TemplateName, param.TemplateCode, param.TemplateContent, param.Remarks)
	u.respond(code, "请求成功", msg)
}

// 查询账号下短息签名
func (u *AliController) QuerySmsSignList() {
	account := u.GetString("id")
	u.checkParam(account)

	sign := aliModel.QuerySmsSignList(account)
	u.respond(0, "请求成功", sign)
}

// 查询账号安全评分
func (u *AliController) GetScore() {
	res := aliModel.GetScore()
	u.respond(0, "请求成功", res)
}

// 查询资产漏洞
func (u *AliController) GetVul() {
	res := aliModel.GetVul()
	u.respond(0, "请求成功", res)
}

// 更新负载均衡列表
func (u *AliController) DescribeLoadBalancers() {
	accountID := u.GetString("id")
	u.checkParam(accountID)

	if res, err := aliModel.DescribeLoadBalancers(accountID); err != nil {
		u.respond(1, "请求失败", err)
	} else {
		u.respond(0, "请求成功", res)
	}
}

// 查询负载均衡监听
func (u *AliController) DescribeLoadBalancerListeners() {
	accountID := u.GetString("id")
	slbId := u.GetString("slbId")

	u.checkParam(accountID, slbId)

	res := aliModel.DescribeLoadBalancerListeners(accountID, slbId)
	u.respond(0, "请求成功", res)
}

// 查询负载均衡后端服务器
func (u *AliController) DescribeLoadBalancerAttribute() {
	accountID := u.GetString("id")
	slbId := u.GetString("slbId")

	u.checkParam(accountID, slbId)

	if res, err := aliModel.DescribeLoadBalancerAttribute(accountID, slbId); err != nil {
		u.respond(1, "请求失败", "请检查服务日志")
	} else {
		u.respond(0, "请求成功", res)
	}
}

// 查询负载均衡列表
func (u *AliController) GetSlbList() {
	accountID := u.GetString("id")
	u.checkParam(accountID)

	res := aliModel.GetSlbList(accountID)
	u.respond(0, "请求成功", res)
}

// 弹性公网
func (u *AliController) DescribeEipAddresses() {
	accountID := u.GetString("id")
	u.checkParam(accountID)

	res := aliModel.DescribeEipAddresses(accountID)
	u.respond(0, "请求成功", res)
}

// 查询域名列表和域名信息
func (u *AliController) QueryDomainList() {
	accountID := u.GetString("id")
	u.checkParam(accountID)

	res := aliModel.QueryDomainList(accountID)
	u.respond(0, "请求成功", res)
}

func (u *AliController) GetOSS() {
	accountID := u.GetString("id")
	u.checkParam(accountID)

	res := aliModel.GetOSS(accountID)
	u.respond(0, "请求成功", res)
}

func (u *AliController) TagDetails() {
	tag := u.GetString("tag")

	res := aliModel.TagDetails(tag)
	u.respond(0, "请求成功", res)
}

// ocr-表格识别
func (u *AliController) TableOcr() {
	object := u.GetString("object")

	if object == "" {
		u.respond(1, "请求失败", "缺少参数")
		return
	}

	resp, err := aliModel.TableOcr(object)
	if err != nil {
		u.respond(1, "请求失败", err)
		return
	}

	u.respond(0, "请求成功", resp)
}

// 函数计算绑定 nas
func (u *AliController) BindingNas() {
	name := u.GetString("name")
	company := u.GetString("company")
	designCenter := u.GetString("designCenter")
	bpmId := u.GetString("bpmId")
	bpmId = strings.ToUpper(bpmId)

	if name == "" || company == "" || designCenter == "" || bpmId == "" {
		u.respond(1, "缺少失败", "缺少参数")
		return
	}
	//nasId := aliModel.BindingNas(name, company, designCenter, bpmId)

	domain, user, pass, mount, date, manDomain, err := aliModel.BindingNas(name, company, designCenter, bpmId)
	if err != nil {
		u.respond(1, "请求失败", err)
	}

	if domain == "无空闲应用" {
		u.respond(2, "请求失败", "无空闲应用")
	}

	infoResp := struct {
		SdDomain     string `json:"sdDomain,omitempty"`
		SdUser       string `json:"sdUser,omitempty"`
		SdPasswd     string `json:"sdPasswd,omitempty"`
		SdMount      string `json:"sdMount,omitempty"`
		SdDate       string `json:"sdDate"`
		SdManaDomain string `json:"sdManaDomain"`
	}{
		SdDomain:     domain,
		SdUser:       user,
		SdPasswd:     pass,
		SdMount:      mount,
		SdDate:       date,
		SdManaDomain: manDomain,
	}

	u.respond(0, "请求成功", infoResp)
}

// 获取所有函数列表
func (u *AliController) ListFunctions() {
	resp, _, err := aliModel.ListFunctions()
	if err != nil {
		u.respond(1, "请求失败", err)
		return
	}

	if len(resp) > 0 {
		u.respond(1, "部分函数更新失败", resp)
	}

	u.respond(0, "请求成功", "全部更新完成")
}

// 解绑应用
func (u *AliController) UnBind() {
	bpmId := u.GetString("bpmId")
	bpmId = strings.ToUpper(bpmId)
	err := aliModel.UnBind(bpmId)
	if err != nil {
		u.respond(1, "请求失败", err)
		return
	}

	u.respond(0, "请求成功", "解绑成功")
}

// 通过接口创建用户-ecs-nas绑定
func (u *AliController) HttpBindDir() {
	bpmId := u.GetString("bpmId")
	bpmId = strings.ToUpper(bpmId)
	str, err := aliModel.HttpBindDir(bpmId)
	if err != nil {
		u.respond(1, "请求失败", err)
		return
	}
	u.respond(0, "请求成功", str)
}

// 清理 ecs nas 绑定关系
func (u *AliController) Clear() {
	bpmId := u.GetString("bpmId")
	bpmId = strings.ToUpper(bpmId)
	resp, err := aliModel.Clear(bpmId)
	if err != nil {
		u.respond(1, "请求失败", err)
	}

	u.respond(0, "请求成功", resp)
}

//func (u *AliController) QueryFcStat() {
//	unUse, use, err := aliModel.QueryFcStat()
//	if err != nil {
//		u.respond(1, "请求失败", err)
//	}
//
//	data := struct {
//		UnUse int
//		Use   int
//	}{
//		UnUse: unUse,
//		Use:   use,
//	}
//
//	if unUse == 0 || use == 0 {
//		u.respond(1, "请求成功", data)
//	}
//
//	u.respond(0, "请求成功", data)
//}

// 用于邀约制查询是否存在账号
func (u *AliController) QueryBpmId() {
	bpmId := u.GetString("bpmId")
	msg, b, err := aliModel.QueryBpmId(bpmId)
	if err != nil {
		u.respond(1, "请求失败", err)
		return
	}

	if b {
		u.respond(0, "请求成功", msg)
	}

	u.respond(1, "请求成功", msg)
}

// 复制大模型
func (u *AliController) CopyModel() {
	bpmId := u.GetString("bpmId")
	model := u.GetString("modelName")

	if bpmId == "" || model == "" {
		u.respond(1, "请求失败", "参数不能为空")
		return
	}

	bpmId = strings.ToUpper(bpmId)
	err := aliModel.CopyLargeModel(bpmId, model)

	if err != nil {
		u.respond(1, "请求失败", err.Error())
		return
	}

	u.respond(0, "请求成功", "复制完成")
}

// 复制 lora
func (u *AliController) CopyLoraModel() {
	bpmId := u.GetString("bpmId")
	bpmId = strings.ToUpper(bpmId)
	model := u.GetString("modelName")

	err := aliModel.CopyLoraModel(bpmId, model)

	if err != nil {
		u.respond(1, "请求失败", err.Error())
		return
	}

	u.respond(0, "请求成功", "复制完成")
}

// 查询模型列表 ModelList
func (u *AliController) ModelList() {

	data := aliModel.ModelList()

	u.respond(0, "请求成功", data)
}

func (u *AliController) NasAmount() {

	aliModel.NasAmount()

	return
}

// 查询未使用的应用
func (u *AliController) FunctionUse() {

	result, err := aliModel.FunctionUse()
	if err != nil {
		u.respond(1, "请求失败", err.Error())
		return
	}

	u.respond(0, "请求成功", result)

}

// 查询使用的应用
func (u *AliController) FunctionUnUse() {

	result, err := aliModel.FunctionUnUse()
	if err != nil {
		u.respond(1, "请求失败", err.Error())
		return
	}

	u.respond(0, "请求成功", result)

}

// 查询lora列表
func (u *AliController) QueryLora() {
	result, err := aliModel.QueryLora()
	if err != nil {
		u.respond(1, "请求失败", err)
		return
	}

	u.respond(0, "请求成功", result)
}

// 复制 lora
func (u *AliController) CopyLora() {
	dir := u.GetString("dir")
	fileName := u.GetString("fileName")
	bpmId := u.GetString("bpmId")

	if fileName == "" || bpmId == "" || dir == "" {
		u.respond(1, "请求成功", "dir/fileName/bpmId 参数不能为空")
		return
	}

	msg, err := aliModel.CopyLora(dir, fileName, bpmId)
	if err != nil {
		u.respond(1, "请求失败", err)
		return
	}

	u.respond(0, "请求成功", msg)
}

// 查询controlNet列表
func (u *AliController) QueryControlNet() {
	result := aliModel.QueryControlNet()
	u.respond(0, "请求成功", result)
}
