/*
 * @Author: nevin
 * @Date: 2023-11-13 11:18:48
 * @LastEditTime: 2023-11-13 15:46:44
 * @LastEditors: nevin
 * @Description: 阿里云
 */
package aliModel

import (
	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	domain20180129 "github.com/alibabacloud-go/domain-20180129/v4/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	fc20230330 "github.com/alibabacloud-go/fc-20230330/v4/client"
	nas20170626 "github.com/alibabacloud-go/nas-20170626/v3/client"
	ocr20210707 "github.com/alibabacloud-go/ocr-api-20210707/v2/client"
	sas20181203 "github.com/alibabacloud-go/sas-20181203/v2/client"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v6/client"
)

// 财务Client
func BillCreateClient(accessKeyId *string, accessKeySecret *string) (_result *bssopenapi20171214.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("business.aliyuncs.com")
	_result = &bssopenapi20171214.Client{}
	_result, _err = bssopenapi20171214.NewClient(config)
	return _result, _err
}

// ECS Client
func EcsCreateClient(accessKeyId *string, accessKeySecret *string, Endpoint string) (_result *ecs20140526.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String(Endpoint)
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

// 短信Client
func SmsCreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// cdn Client
func CdnCreateClient(accessKeyId *string, accessKeySecret *string) (_result *cdn20180510.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("cdn.aliyuncs.com")
	_result = &cdn20180510.Client{}
	_result, _err = cdn20180510.NewClient(config)
	return _result, _err
}

// 账号安全 Client
func SasCreateClient(accessKeyId *string, accessKeySecret *string) (_result *sas20181203.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("tds.aliyuncs.com")
	_result = &sas20181203.Client{}
	_result, _err = sas20181203.NewClient(config)
	return _result, _err
}

// SLB client
func SlbCreateClient(accessKeyId *string, accessKeySecret *string) (_result *slb20140515.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("slb.aliyuncs.com")
	_result = &slb20140515.Client{}
	_result, _err = slb20140515.NewClient(config)
	return _result, _err
}

// 专业网络 client
func VpcCreateClient(accessKeyId *string, accessKeySecret *string) (_result *vpc20160428.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("vpc.cn-beijing.aliyuncs.com")
	_result = &vpc20160428.Client{}
	_result, _err = vpc20160428.NewClient(config)
	return _result, _err
}

// 域名 client
func DomainCreateClient(accessKeyId *string, accessKeySecret *string) (_result *domain20180129.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("domain.aliyuncs.com")
	_result = &domain20180129.Client{}
	_result, _err = domain20180129.NewClient(config)
	return _result, _err
}

func OcrCreateClient(accessKeyId *string, accessKeySecret *string) (_result *ocr20210707.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("ocr-api.cn-hangzhou.aliyuncs.com")
	_result = &ocr20210707.Client{}
	_result, _err = ocr20210707.NewClient(config)
	return _result, _err
}

func NasCreateClient(accessKeyId *string, accessKeySecret *string) (_result *nas20170626.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("nas.cn-hangzhou.aliyuncs.com")
	_result = &nas20170626.Client{}
	_result, _err = nas20170626.NewClient(config)
	return _result, _err
}

func FcCreateClient(accessKeyId *string, accessKeySecret *string) (_result *fc20230330.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("1748293280207491.cn-hangzhou.fc.aliyuncs.com")
	_result = &fc20230330.Client{}
	_result, _err = fc20230330.NewClient(config)
	return _result, _err
}
