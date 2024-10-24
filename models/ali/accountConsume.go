package aliModel

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"time"
	"yw_cloud/models/db"
	utilModel "yw_cloud/models/util"
)

/*
阿里云账号消费
*/
type consum struct {
	PretaxAmount float32 `json:"PaymentAmount"`
	Months       string  `json:"Months"`
	Year         string  `json:"Year"`
	AccountId    string  `json:"AccountId"`
}

func GetConsume(date string) []consum {
	c := []consum{}
	ct, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	coll := db.GetCollection("consume")
	// id字段设置成0，表示排除 id 字段
	opts := options.Find().SetProjection(bson.D{{"_id", 0}})

	for _, v := range utilModel.Key {
		cursor, err := coll.Find(ct, bson.D{{"accountid", v.AccountId}, {"year", date}}, opts)
		if err != nil {
			fmt.Println("err: ", err)
		}

		//var results []bson.D
		var results []consum
		if _err := cursor.All(ct, &results); _err != nil {
			log.Println("查询账号消费失败:", _err)
		}

		c = append(c, results...)
	}
	return c
}
