package ci

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type Record struct {
	Filed string `json:"Filed" orm:"column(Filed)"`
	Count string `json:"Count" orm:"column(Count)"`
}

func (p Record) TableName() string {
	return "jenkins"
}

func GetDB(selectFiled, StartDate, EndDate string) []Record {
	o := orm.NewOrm()
	rds := new([]Record)

	sql := fmt.Sprintf("select %s as Filed ,COUNT(*) as Count  from jenkins  where date between '%s' and '%s'  GROUP BY  %s;", selectFiled, StartDate, EndDate, selectFiled)

	if _, err := o.Raw(sql).QueryRows(rds); err != nil {
		fmt.Println("序列化结构体失败: ", err)
		return nil
	}
	return *rds

}
