package controllers

import (
	sonarqube "yw_cloud/models/soanr"
)

type SonarController struct {
	BaseController
}

// 更新sonarqube扫描的指标
func (u *SonarController) SonarCollect() {
	sonarqube.SonarCollect()
	u.respond(0, "请求成功", "更新完成")
}

// 获取单个项目的指标
func (u *SonarController) GetCollect() {
	p := u.GetString("project")

	if rep, err := sonarqube.GetCollect(p); err != nil {
		u.respond(1, "请求失败", err)
	} else {
		u.respond(0, "请求成功", rep)
	}

}
