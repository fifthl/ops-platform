package aliModel

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"yw_cloud/models/db"
)

/*
获取月分摊账单
查询用户某个账期内账单总览信息
*/
var (
	ctx = context.Background()
)

type overview struct {
	PretaxAmount float32 `json:"PretaxAmount"`
	ProductName  string  `json:"ProductName"`
}

type Overviews struct {
	Item      []overview
	AccountId string `json:"AccountId"`
	Months    string `json:"Months"`
}

func BillOverview(AccountID, StartTime, EndTime string) (allResult *[]Overviews, err error) {

	coll := db.GetCollection("overview")
	// mongo 区间查询
	filter := bson.M{"months": bson.M{"$gte": StartTime, "$lte": EndTime}, "accountid": AccountID}

	c, err := coll.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("find 错误", err)
	}

	Result := new([]Overviews)
	if e := c.All(context.TODO(), Result); e != nil {
		fmt.Println("存入allResult失败", e)

	}

	for _, overviews := range *Result {

		for i := 0; i < len(overviews.Item); i++ {
			v, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", overviews.Item[i].PretaxAmount), 32)
			overviews.Item[i].PretaxAmount = float32(v)

		}

	}

	return Result, nil
}
