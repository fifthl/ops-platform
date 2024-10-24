package aliModel

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"sort"

	// don't forget this
	_ "github.com/go-sql-driver/mysql"
)

// ECS 类型费用分类/分布
type TagStruct struct {
	Price float32 `orm:"column(price)" json:"Price"` // 金额
	Tag   string  `orm:"column(tag)" json:"Tag" `    // 标签
	Num   int     `orm:"column(num)" json:"Num"`     // 数量
}

func EcsPrice() []TagStruct {
	res := []TagStruct{}
	o := orm.NewOrmUsingDB("ecs")
	a := new(TagStruct)

	var (
		price []float32
		tags  []string
		num   []int
	)

	_, err := o.Raw("SELECT SUM(price) as price, tags, COUNT(*) as num FROM instances_info WHERE tags IN ('NC', 'k8s', 'AI', 'BPM', '机电', '运维', '设计平台', '管家', '市场', '尚层家', '共享平台', '信管/大师', 'BI') GROUP BY tags;").QueryRows(&price, &tags, &num)
	if err != nil {
		fmt.Println("查询 ecs tag费用失败:", err)
	}

	for i := 0; i < len(price); i++ {
		a.Price = price[i]
		a.Tag = tags[i]
		a.Num = num[i]
		res = append(res, *a)
	}

	// 排序
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Price > res[j].Price
	})
	return res
}
