package aliModel

import (
	"errors"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"time"
	"yw_cloud/models/db"
	utilModel "yw_cloud/models/util"
)

/*
创建短信模板
*/
func AddSmsTemplate(TemplateType int32, TemplateName, TemplateContent, Remark string, AccountID string) (result string, err error) {

	account := utilModel.Key[AccountID]
	Id, Secret := utilModel.Decrypt(account.ID, account.Secret)
	client, _ := SmsCreateClient(&Id, &Secret)

	addSmsTemplateRequest := &dysmsapi20170525.AddSmsTemplateRequest{
		TemplateType:    tea.Int32(TemplateType),
		TemplateName:    tea.String(TemplateName),
		TemplateContent: tea.String(TemplateContent),
		Remark:          tea.String(Remark), // 申请说明
	}
	runtime := &util.RuntimeOptions{}

	_resp, _err := client.AddSmsTemplateWithOptions(addSmsTemplateRequest, runtime)
	if _err != nil {
		log.Println("添加短信模板请求失败:", _err)
		//return "添加模板失败"
		return "", _err
	}

	// 模板添加成功后 _resp.Body.Code 为 OK
	if (*_resp.Body.Code) != "OK" {
		err = errors.New(fmt.Sprintf("短信模板创建失败: %v", *_resp.Body.Message))
		return "", err
	}

	// 完成后将申请成功的短息模板信息写入 mongo 中
	coll := db.GetCollection("SmsTemplate")
	filter := bson.D{{"templatecode", _resp.Body.TemplateCode}}
	opts := options.Update().SetUpsert(true)

	// writeMongoParam 转换参数
	AuditStatus, CreateDate, OuterTemplateType, TemplateCode, reason := writeMongoParam(TemplateType, *_resp.Body.TemplateCode)

	// 拼接 mongo 串
	update := bson.D{{"$set", bson.D{{"auditstatus", AuditStatus}, {"createdate", CreateDate}, {"outertemplatetype", OuterTemplateType}, {"reason", reason}, {"templatecode", TemplateCode}, {"templatecontent", TemplateContent}, {"templatename", TemplateName}}}}

	upResult, uperr := coll.UpdateOne(context.TODO(), filter, update, opts)
	if uperr != nil {
		log.Println("插入失败: ", uperr)
	}
	fmt.Println(upResult.UpsertedID)

	log.Printf("短信模板添加成功, mongo id: $v\n", upResult.UpsertedID)

	return "成功", nil

}

// writeMongoParam 转换参数
func writeMongoParam(TemplateType int32, Code string) (AuditStatus, CreateDate, OuterTemplateType, TemplateCode string, reason Reason) {
	AuditStatus = "审核中"
	CreateDate = time.Now().Format("2006-01-02 15:04:05")
	switch TemplateType {
	case 1:
		OuterTemplateType = "验证码短信"
	case 0:
		OuterTemplateType = "审核通过"
	case 2:
		OuterTemplateType = "推广短信"
	case 3:
		OuterTemplateType = "国际/港澳台短信"
	case 7:
		OuterTemplateType = "数字短信"
	}
	TemplateCode = Code
	reason = Reason{}

	return
}
