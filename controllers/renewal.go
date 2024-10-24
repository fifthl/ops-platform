/*
 * @Author: nevin
 * @Date: 2023-11-13 14:21:29
 * @LastEditTime: 2023-11-15 14:11:13
 * @LastEditors: nevin
 * @Description: 续费记录
 */
package controllers

import (
	"encoding/json"
	"strconv"
	"time"
	renewalModel "yw_cloud/models/renewal"
)

type RenewalController struct {
	BaseController
}

// @Title 续费
// @Description 获取阿里云历史续费
// @Success 200 {object} models.createOneRenewal
// @router / [get]
func (u *RenewalController) GetAliyunHistoryRenewals() {

	//查询参数
	var queryMap = make(map[string]interface{})
	queryMap["month"] = u.GetString("month")
	queryMap["year"] = u.GetString("year")
	queryMap["account_name"] = u.GetString("account_name")
	queryMap["account_no"] = u.GetString("account_no")
	queryMap["renewal_type"] = 1

	if u.GetString("invoice") != "" {
		invoice, _ := strconv.ParseInt(u.GetString("invoice"), 10, 8)
		queryMap["invoice"] = invoice
	}

	//分页参数
	limit, _ := strconv.ParseInt(u.GetString("limit"), 10, 64)
	if u.GetString("limit") == "" || limit <= 0 {
		limit = 100000
	}
	page, _ := strconv.ParseInt(u.GetString("page"), 10, 64)
	if u.GetString("page") == "" || page <= 0 {
		page = 1
	}

	list, con := renewalModel.GetAliyunHistoryRenewals(queryMap, limit, page)
	var res = ResList[renewalModel.AliYunRenewal]{
		List:  list,
		Count: con,
	}

	u.respond(0, "请求成功", res)

}

// @Title 续费
// @Description 获取历史续费底表数据
// @Success 200 {object} models.createOneRenewal
// @router / [get]
func (u *RenewalController) GetAllHistoryRenewals() {

	//查询参数
	var queryMap = make(map[string]interface{})
	queryMap["month"] = u.GetString("month")
	queryMap["year"] = u.GetString("year")
	queryMap["account_name"] = u.GetString("account_name")
	queryMap["account_no"] = u.GetString("account_no")
	queryMap["create_time"] = u.GetString("create_time")

	if u.GetString("invoice") != "" {
		invoice, _ := strconv.ParseInt(u.GetString("invoice"), 10, 8)
		queryMap["invoice"] = invoice
	}

	if u.GetString("renewal_type") != "" {
		invoice, _ := strconv.ParseInt(u.GetString("renewal_type"), 10, 8)
		queryMap["renewal_type"] = invoice
	}

	//分页参数
	limit, _ := strconv.ParseInt(u.GetString("limit"), 10, 64)
	if u.GetString("limit") == "" || limit <= 0 {
		limit = 100000
	}
	page, _ := strconv.ParseInt(u.GetString("page"), 10, 64)
	if u.GetString("page") == "" || page <= 0 {
		page = 1
	}

	list, con := renewalModel.GetAllHistoryRenewals(queryMap, limit, page)
	var res = ResList[renewalModel.RenewalRecord]{
		List:  list,
		Count: con,
	}

	u.respond(0, "请求成功", res)

}

// @Title 新增续费
// @Description 增加一条续费记录
// @Success 200 {object} models.createRenewal
// @router / [post]
func (u *RenewalController) CreateRenewal() {

	// 初始化数据 ，续费类型默认为其他，默认没有发票
	renewal := renewalModel.RenewalRecord{
		RenwalType: 0,
		Invoice:    0,
		CreateTime: time.Now(),
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &renewal)

	if renewal.Month == "" || renewal.Year == "" || renewal.AccountName == "" || renewal.AccountNo == "" || renewal.Money == 0 {
		u.respond(1, "数据不能为空")
		return
	}

	err := renewalModel.CreateRenewalRecord(renewal)
	if err != nil {
		u.respond(1, "创建续费记录失败")
	} else {
		u.respond(0, "创建续费记录成功")
	}

}

// // @Title 修改续费
// // @Description 修改一条续费记录
// // @Success 200 {object} models.updateRenewal
// // @router / [put]
// func (u *RenewalController) UpdateRenewal() {

// 	renewalId := u.GetString("renewalId")

// 	var ob renewalModel.RenewalRecord
// 	json.Unmarshal(u.Ctx.Input.RequestBody, &ob)

// 	err := renewalModel.UpdateRenewalInfo(renewalId, ob)
// 	if err != nil {
// 		u.respond(111, err.Error(), nil)
// 	} else {
// 		u.respond(0, "请求成功", 1)
// 	}

// }

// @Title 删除某条续费
// @Description 删除某条续费记录
// @Success 200 {object} models.deleteRenewal
// @router / [delete]
func (u *RenewalController) DeleteRenewal() {
	renewalId := u.GetString("renewalId")
	err := renewalModel.DeleteById(renewalId)
	if err != nil {
		u.respond(111, err.Error(), nil)
	} else {
		u.respond(0, "删除成功", 1)
	}
}
