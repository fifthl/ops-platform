package controllers

import (
	"encoding/json"
	"fmt"
	ci "yw_cloud/models/jenkins"
)

type JenkinsController struct {
	BaseController
}

//func (u *JenkinsController) SaveDB() {
//	c := new(ci.ReleaseHistory)
//
//	if err := json.Unmarshal(u.Ctx.Input.RequestBody, c); err != nil {
//		fmt.Println("err: ", err)
//	}
//
//	if err := ci.SaveDB(c); err != nil {
//		u.respond(100, "写入失败", err.Error())
//	} else {
//		u.respond(0, "写入成功", err)
//	}
//}

func (u *JenkinsController) NewSaveDb() {
	r := new(ci.ReleaseInfo)

	// 序列化
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, r); err != nil {
		fmt.Println("err: ", err)
	}

	// 判断参数
	if r.Name == "" || r.Email == "" || r.ChangeInfo == "" || r.BuildNum == "" || r.GitUrl == "" || r.Commit == "" {
		u.respond(1, "请求失败", "参数不能为空")
		a := fmt.Sprintf("Project:%v Branch:%v Name:%v Email:%v ChangeInfo:%v BuildNum:%v GitUrl:%v Commit:%v", r.Project, r.Branch, r.Name, r.Email, r.ChangeInfo, r.BuildNum, r.GitUrl, r.Commit)
		fmt.Println("失败记录: ", a)
		return
	}

	if r.Name == "1468588829" || r.Name == "zhangqinbo" {
		// 不保留流水线
		fmt.Println("info为流水线信息，不保留db中")
		return
	}

	// 调用
	if err := ci.NewSaveDb(r); err != nil {
		u.respond(100, "写入失败", err.Error())
	} else {
		u.respond(0, "请求成功", "写入成功")
	}
}

func (u *JenkinsController) GetDB() {
	filed := u.GetString("filed")
	start := u.GetString("start")
	end := u.GetString("end")

	if filed == "" || start == "" || end == "" {
		u.respond(100, "请示失败", "参数值为空")
	}

	history := ci.GetDB(filed, start, end)
	u.respond(0, "请求成功", history)

}
