package aliModel

import (
	"fmt"
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/beego/beego/v2/client/orm"
	"log"
	"strings"
	utilModel "yw_cloud/models/util"
)

/*
内容分发,将源站的内容主动预热到缓存节点上
阿里云接口: 预热URL (PushObjectCache)
*/

func PushObjectCache(Urls []string) (mess string) {
	// 配置 QueryBuilder 对象查询 sql
	dp, _ := orm.NewQueryBuilder("mysql")

	dp.Select("account").
		From("cdn_domain").
		Where("domain=?")

	// 生成 sql 语句和对象
	sql, o := dp.String(), orm.NewOrm()

	for _, url := range Urls {
		var account []string
		domain := selectAccountPush(url)

		// 查询传入的domain属于哪个阿里云账号
		if _, err := o.Raw(sql, domain).QueryRows(&account); err != nil {
			fmt.Println("预热接口查询账号id失败:", err)
		}

		a := utilModel.Key[account[0]]
		id, secrt := utilModel.Decrypt(a.ID, a.Secret)

		// 使用拿出来的账号请求阿里云
		// 看情况是否要对返回的 pushid 做处理
		_ = aliyumCdnPush(id, secrt, url)

	}
	return "URL预热成功"

}

// 调用阿里云 sdk 相关代码发起请求
func aliyumCdnPush(AK, SK, URL string) (PushTaskId string) {
	client, _ := CdnCreateClient(tea.String(AK), tea.String(SK))

	pushObjectCacheRequest := &cdn20180510.PushObjectCacheRequest{
		ObjectPath: tea.String(URL),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.PushObjectCacheWithOptions(pushObjectCacheRequest, runtime)
	if err != nil {
		log.Println("预热接口调用失败:", err)
		return
	}
	return *resp.Body.PushTaskId

}

// 从mysql查找域名对应的阿里云账号
func selectAccountPush(url string) (domain string) {
	u := strings.Split(url, "/")
	domain = u[0]

	return domain
}
