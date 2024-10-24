package sonarqube

import (
	"github.com/beego/beego/v2/client/orm"
	"log"
)

func GetCollect(project string) (result ResponseCollect, err error) {
	o := orm.NewOrm()
	result = ResponseCollect{Name: project}
	if err = o.Read(&result); err != nil {
		log.Printf("读取%s扫描失败: %v\n", project, err)
		return ResponseCollect{}, err
	}
	//log.Println(r)
	return result, nil
}
