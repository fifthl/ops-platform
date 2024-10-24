package aliModel

import (
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

func ModifySmsTemplate(AccountID string, TemplateType int32, TemplateName, TemplateCode, TemplateContent, Remark string) (code int, msg string) {
	account := utilModel.Key[AccountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)
	client, _ := SmsCreateClient(tea.String(id), tea.String(secret))

	modifySmsTemplateRequest := &dysmsapi20170525.ModifySmsTemplateRequest{
		TemplateType:    tea.Int32(TemplateType),
		TemplateName:    tea.String(TemplateName),
		TemplateCode:    tea.String(TemplateCode),
		TemplateContent: tea.String(TemplateContent),
		Remark:          tea.String(Remark),
	}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.ModifySmsTemplateWithOptions(modifySmsTemplateRequest, runtime)
	if _err != nil {
		log.Println(_err)
		return
	}

	log.Println(*resp.Body)
	if *resp.Body.Code != "OK" {
		// 失败
		log.Println("更改短信模板失败")
		return 1, *resp.Body.Message
	}

	return 0, "模板更新成功"

}
