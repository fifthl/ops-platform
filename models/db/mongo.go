/*
 * @Author: nevin
 * @Date: 2023-11-10 09:19:57
 * @LastEditTime: 2023-11-22 10:25:16
 * @LastEditors: nevin
 * @Description: mongo数据库链接
 */
package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var (
	Client *mongo.Client
)

func init() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(beego.AppConfig.String("mongodb_addr")))
	if err != nil {
		log.Fatalln(err)
		log.Fatalln("链接失败")
		return
	}

	Client = client
}

func GetCollection(collectionName string) *mongo.Collection {
	collection := Client.Database(beego.AppConfig.String("mongodb_db_name")).Collection(collectionName)
	return collection
}

type Ids[T int | int8 | int16 | int32 | int64] struct {
	IdValue    T         `bson:"id_value"`
	IdName     string    `bson:"id_name"`
	UpdateTime time.Time `bson:"update_time"`
}

/*
 * 创建自增ID
 * idName id名称
 * idValue id初始值
 */
func CreateId[T int | int8 | int16 | int32 | int64](idName string, idValue T) (T, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := GetCollection("t_ids")

	// 有就更新
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := collection.FindOneAndUpdate(ctx, bson.M{"id_name": idName}, bson.M{"$inc": bson.M{"id_value": 1}, "$set": bson.M{"update_time": time.Now()}}, options)

	var oldId Ids[T]
	res.Decode(&oldId)

	// 返回
	if res.Err() == nil {
		return oldId.IdValue, nil
	}

	// 没有就创建
	newIds := Ids[T]{
		idValue,
		idName,
		time.Now(),
	}

	_, err2 := collection.InsertOne(ctx, newIds)

	if err2 != nil {
		return newIds.IdValue, err2
	}

	// 没有就创建
	return newIds.IdValue, nil
}

func FindAll(coll *mongo.Collection, ctx context.Context, FindValue, MonthsValue string) (cursor *mongo.Cursor) {

	cursor, err := coll.Find(ctx, bson.D{{"accountid", FindValue}, {"months", MonthsValue}})
	if err != nil {
		fmt.Println("查询失败", err)
	}
	return cursor
}
