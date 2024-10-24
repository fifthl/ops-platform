package renewal

import (
	"context"
	"fmt"
	"log"
	"time"
	"yw_cloud/models/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	aliYunRenewalDB = "aliyunRenewal" // 阿里云账号续费
	renewalRecordDB = "renewalRecord" // 续费底表
)

type AliYunRenewal struct {
	// ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" `
	AccountName string  `json:"account_name" bson:"account_name" `
	AccountNo   string  `json:"account_no" bson:"account_no" `
	Year        string  `json:"year"  bson:"year" `
	Month       string  `json:"month"  bson:"month" `
	Money       float64 `json:"money"  bson:"money" `
	// CreateTime  time.Time `json:"create_time"  bson:"create_time" `
}

type RenewalRecord struct {
	// RenewalID   string    `json:"renweal_id" bson:"renewal_id"`
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" `
	AccountNo   string             `json:"account_no" bson:"account_no" `
	AccountName string             `json:"account_name" bson:"account_name" `
	RenwalType  int8               `json:"renewal_type" bson:"renewal_type" ` //0:其他 1:阿里云续费 2:钉钉续费 3：电信网续费 4：小程序续费 5：萤石账号续费
	Year        string             `json:"year"  bson:"year" `
	Month       string             `json:"month"  bson:"month" `
	Money       float64            `json:"money"  bson:"money" `
	Invoice     int8               `json:"invoice"  bson:"invoice" ` //是否收到发票 0：未收到	1：已收到
	Backup      string             `json:"backup"  bson:"backup" `
	CreateTime  time.Time          `json:"create_time"  bson:"create_time" `
}

// group by
// 获取阿里云历史续费信息-读取底表数据
func GetAliyunHistoryRenewals(queryMap map[string]interface{}, limit, page int64) (list []AliYunRenewal, con int64) {
	var renewals []AliYunRenewal

	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := bson.M{}
	//循环query，账号，账号no进行模糊查询，其他精确查询
	for key, value := range queryMap {
		if value != "" {
			if key == "account_name" || key == "account_no" {
				query[key] = primitive.Regex{Pattern: fmt.Sprintf("%s", value), Options: "i"}
			} else {
				query[key] = value
			}
		}
	}

	collection := db.GetCollection(renewalRecordDB)
	//----------
	// group by
	fileds := bson.D{{Key: "account_no", Value: "$account_no"}, {Key: "account_name", Value: "$account_name"}, {Key: "year", Value: "$year"}, {Key: "month", Value: "$month"}, {Key: "renewal_type", Value: "$renewal_type"}, {Key: "invoice", Value: "$invoice"}}
	//需要显示的字段
	project := bson.D{{Key: "_id", Value: 0}, {Key: "account_no", Value: "$_id.account_no"}, {Key: "account_name", Value: "$_id.account_name"}, {Key: "year", Value: "$_id.year"}, {Key: "month", Value: "$_id.month"}, {Key: "renewal_type", Value: "$_id.renewal_type"}, {Key: "invoice", Value: "$_id.invoice"}, {Key: "money", Value: "$money"}}
	//排序-年月倒序
	sort := bson.D{{Key: "year", Value: -1}, {Key: "month", Value: -1}}
	pline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: fileds}, {Key: "money", Value: bson.D{{Key: "$sum", Value: "$money"}}}}}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: sort}},
		{{Key: "$match", Value: query}},
		{{Key: "$limit", Value: limit}},
		{{Key: "$skip", Value: limit * (page - 1)}},
	}

	cur, err := collection.Aggregate(ctx, pline)
	//----------

	var count int64 //TODO 查询不到总数，凑合返回个count
	if err != nil {
		//打印err
		fmt.Print(err)
		count = 0
		renewals = nil
		return renewals, count
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		count++

		var result AliYunRenewal
		err := cur.Decode(&result)

		if err != nil {
			// log.Fatal(err)
			count = 0
			renewals = nil
			return renewals, count
		}
		// 写到返回数据
		renewals = append(renewals, result)
	}

	if err := cur.Err(); err != nil {
		count = 0
		renewals = nil
		return renewals, count
	}

	return renewals, count
}

// // 获取阿里云历史续费信息 - 暂时保留
// func GetAliyunHistoryRenewals(queryMap map[string]interface{}, limit, page int64) (list []AliYunRenewal, con int64) {
// 	var renewals []AliYunRenewal

// 	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	query := bson.M{}
// 	//循环query，账号，账号no进行模糊查询，其他精确查询
// 	for key, value := range queryMap {
// 		if value != "" {
// 			if key == "account_name" || key == "account_no" {
// 				query[key] = primitive.Regex{Pattern: fmt.Sprintf("%s", value), Options: "i"}
// 			} else {
// 				query[key] = value
// 			}
// 		}
// 	}

// 	collection := db.GetCollection(aliYunRenewalDB)
// 	//设置分页查询
// 	var findOptions *options.FindOptions = &options.FindOptions{}
// 	if limit > 0 {
// 		findOptions.SetLimit(limit)
// 		findOptions.SetSkip(limit * (page - 1))
// 	}
// 	count, _ := collection.CountDocuments(ctx, query)

// 	cur, err := collection.Find(ctx, query, findOptions)
// 	if err != nil {
// 		//打印err
// 		fmt.Print(err)
// 		count = 0
// 		renewals = nil
// 		return renewals, count
// 	}
// 	defer cur.Close(ctx)

// 	for cur.Next(ctx) {
// 		var result AliYunRenewal
// 		err := cur.Decode(&result)

// 		if err != nil {
// 			// log.Fatal(err)
// 			count = 0
// 			renewals = nil
// 			return renewals, count
// 		}
// 		// 写到返回数据
// 		renewals = append(renewals, result)
// 	}

// 	if err := cur.Err(); err != nil {
// 		// log.Fatal(err)
// 		count = 0
// 		renewals = nil
// 		return renewals, count
// 	}

// 	return renewals, count
// }

// 获取所有续费类型历史信息
func GetAllHistoryRenewals(queryMap map[string]interface{}, limit, page int64) (list []RenewalRecord, con int64) {
	var renewals []RenewalRecord

	var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := bson.M{}
	//循环query，账号，账号no进行模糊查询，其他精确查询
	for key, value := range queryMap {
		if value != "" {
			if key == "account_name" || key == "account_no" {
				query[key] = primitive.Regex{Pattern: fmt.Sprintf("%s", value), Options: "i"}
			} else {
				query[key] = value
			}
		}
	}

	collection := db.GetCollection(renewalRecordDB)
	//设置分页查询
	var findOptions *options.FindOptions = &options.FindOptions{}
	if limit > 0 {
		findOptions.SetLimit(limit)
		findOptions.SetSkip(limit * (page - 1))
	}
	count, _ := collection.CountDocuments(ctx, query)

	cur, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		//打印err
		fmt.Print(err)
		count = 0
		renewals = nil
		return renewals, count
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result RenewalRecord
		err := cur.Decode(&result)

		if err != nil {
			// log.Fatal(err)
			count = 0
			renewals = nil
			return renewals, count
		}
		// 写到返回数据
		renewals = append(renewals, result)
	}

	if err := cur.Err(); err != nil {
		// log.Fatal(err)
		count = 0
		renewals = nil
		return renewals, count
	}

	return renewals, count
}

// 新增一条Or Update 阿里云续费数据 -- 暂时保留
// func updateAliyunRenewal(renewal RenewalRecord, isDelete bool) (err error) {

// 	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	collection := db.GetCollection(aliYunRenewalDB)
// 	// collection2 := db.GetCollection(renewalRecordDB)

// 	//查询历史月份是否有记录，如果有累加金额
// 	//构建查询条件
// 	queryMap := map[string]interface{}{
// 		"month":      renewal.Month,
// 		"year":       renewal.Year,
// 		"account_no": renewal.AccountNo,
// 	}

// 	// //获取底表 阿里云续费相同年月账号记录的金额总和
// 	var finalMoney float64
// 	if isDelete {
// 		finalMoney = 0 - renewal.Money
// 	} else {
// 		finalMoney = renewal.Money
// 	}
// 	options := options.Update().SetUpsert(true)
// 	_, err = collection.UpdateOne(ctx, queryMap, bson.M{"$inc": bson.M{"money": finalMoney}, "$set": bson.M{"create_time": time.Now(), "account_name": renewal.AccountName}}, options)

// 	return err
// }

// 续费底表 新增一条记录
func CreateRenewalRecord(renewal RenewalRecord) (err error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(renewalRecordDB)

	_, err1 := collection.InsertOne(ctx, renewal)
	if err1 != nil {
		return err1
	}

	// //如果续费类型是 阿里云账号续费，处理aliyunRenewal表数据
	// if renewal.RenwalType == 1 {
	// 	err := updateAliyunRenewal(renewal, false)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

/*
 * 删除一条续费记录
 */
func DeleteById(renewalId string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := db.GetCollection(renewalRecordDB)

	objID, _ := primitive.ObjectIDFromHex(renewalId)

	filter := bson.D{{Key: "_id", Value: objID}}

	result := collection.FindOneAndDelete(ctx, filter)
	if result.Err() != nil {
		log.Println("Find error: ", result.Err())
		return result.Err()
	}

	var renewal RenewalRecord
	err = result.Decode(&renewal)
	if err != nil {
		return err
	}
	//如果是 阿里云账号续费，处理aliyunRenewal表数据
	// if renewal.RenwalType == 1 {
	// 	err := updateAliyunRenewal(renewal, true)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
