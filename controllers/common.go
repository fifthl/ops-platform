/*
 * @Author: nevin
 * @Date: 2023-11-15 13:34:03
 * @LastEditTime: 2023-11-15 15:06:36
 * @LastEditors: nevin
 * @Description:
 */
package controllers

import (
	"github.com/astaxie/beego"
)

// 分页数据
type ResList[T any] struct {
	Count int64 `json:"count"`
	List  []T   `json:"list"`
}

type BaseController struct {
	beego.Controller
}

// 给结构体增加返回值封装函数
func (c *BaseController) respond(code int, message string, data ...interface{}) {
	c.Ctx.Output.SetStatus(200) // 设置错误码

	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	c.Data["json"] = struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Code:    code,
		Message: message,
		Data:    d,
	}
	c.ServeJSON()
}
