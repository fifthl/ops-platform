package aliModel

import (
	"context"
	"encoding/json"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"yw_cloud/models/db"
	utilModel "yw_cloud/models/util"
)

/*
查询短信模板
*/

type Reason struct {
	RejectInfo string `json:"RejectInfo" orm:"column(rejectsubinfo)"`
	RejectDate string `json:"RejectDate" orm:"column(rejectdate)"`
}

type SmsTemplateList struct {
	AuditStatus       string `json:"AuditStatus,omitempty" orm:"column(auditstatus)"`
	CreateDate        string `json:"CreateDate,omitempty" orm:"column(createdate)"`
	OuterTemplateType int    `json:"OuterTemplateType" orm:"column(outertemplatetype)"`
	Reason            Reason `json:"Reason" orm:"column(reason)"`
	TemplateCode      string `json:"TemplateCode,omitempty" orm:"column(templatecode)"`
	TemplateContent   string `json:"TemplateContent,omitempty" orm:"column(templatecontent)"`
	TemplateName      string `json:"TemplateName,omitempty" orm:"column(templatename)"`
	AccountID         string `json:"AccountID,omitempty"`
}

func changAUDIT(s *SmsTemplateList) {
	switch s.AuditStatus {
	case "AUDIT_STATE_INIT":
		s.AuditStatus = "审核中"
	case "AUDIT_STATE_PASS":
		s.AuditStatus = "审核通过"
	case "AUDIT_STATE_NOT_PASS":
		s.AuditStatus = "审核未通过"
	case "AUDIT_STATE_CANCEL", "AUDIT_SATE_CANCEL":
		s.AuditStatus = "取消审核"
	}

}

func changType(s *SmsTemplateList) (Templatetype string) {
	switch s.OuterTemplateType {
	case 1:
		Templatetype = "验证码短信"
	case 0:
		Templatetype = "通知短信"
	case 2:
		Templatetype = "推广短信"
	case 3:
		Templatetype = "国际/港澳台短信"
	case 7:
		Templatetype = "数字短信"
	}
	return Templatetype
}

func GetSmsTemplate(AccountID string, PageIndex, PageSize int64) (stls []SmsTemplateList) {

	// 更新
	v := utilModel.Key[AccountID]
	ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
	client, _ := SmsCreateClient(tea.String(ID), tea.String(Secret))

	querySmsTemplateListRequest := &dysmsapi20170525.QuerySmsTemplateListRequest{
		PageSize: tea.Int32(50),
	}
	runtime := &util.RuntimeOptions{}

	_resp, _err := client.QuerySmsTemplateListWithOptions(querySmsTemplateListRequest, runtime)
	if _err != nil {
		log.Println("请求sms模板列表失败: ", _err)
		return nil
	}

	// 拿所有短信列表
	for _, smsList := range _resp.Body.SmsTemplateList {
		stl := new(SmsTemplateList)

		if err := json.Unmarshal([]byte(smsList.GoString()), stl); err != nil {
			log.Println("序列化结构体失败: ", err)
		}
		stl.AccountID = v.AccountId
		// 替换
		changAUDIT(stl)

		// 不返回，每行写入 db
		coll := db.GetCollection("SmsTemplate")

		filter := bson.D{{"templatecode", stl.TemplateCode}}
		opts := options.Update().SetUpsert(true)
		update := bson.D{{"$set", bson.D{{"auditstatus", stl.AuditStatus},
			{"createdate", stl.CreateDate},
			{"outertemplatetype", changType(stl)},
			{"reason", stl.Reason}, {"templatecode", stl.TemplateCode},
			{"templatecontent", stl.TemplateContent},
			{"templatename", stl.TemplateName},
			{"accountid", stl.AccountID}}}}

		_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			log.Println("插入mongo失败: ", err)
			continue
		}
		stls = append(stls, *stl)

	}

	return stls

	// 旧
	//coll := db.GetCollection("SmsTemplate")
	//
	//opts := options.Find()
	//opts.SetProjection(bson.D{{"_id", 0}})
	//
	//// 分页
	//if PageSize > 0 {
	//	opts.SetLimit(PageSize)
	//	opts.SetSkip(PageSize * PageIndex)
	//
	//}
	//
	//filter := bson.M{"accountid": AccountID}
	//
	//// 根据 account 查询对应下面的模板
	//r, err := coll.Find(context.TODO(), filter, opts)
	//if err != nil {
	//	log.Println("短信列表查找失败: ", err)
	//	return nil
	//}
	//
	//var result []SmsTemplateList
	//if _err := r.All(context.TODO(), &result); _err != nil {
	//	log.Println("短信转存结构体失败: ", _err)
	//	return nil
	//}

	//return &result
}
