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
内容分发(CDN),刷新节点上的文件内容。被刷新的文件缓存将立即失效，新的请求将回源获取最新的文件
阿里云接口: 刷新缓存 (RefreshObjectCaches)
*/

func RefreshObjectCaches(Urls []string, RefresType string) (mess string) {
	//var m []string

	// 配置 QueryBuilder 对象查询 sql
	dp, _ := orm.NewQueryBuilder("mysql")

	dp.Select("account").
		From("cdn_domain").
		Where("domain=?")

	// 生成 sql 语句和对象,用于判断刷新的域名属于哪个账号
	sql, o := dp.String(), orm.NewOrm()

	for _, url := range Urls {
		var account []string
		domain := selectAccountRefres(url)

		// 查询传入的domain属于哪个阿里云账号
		if _, err := o.Raw(sql, domain).QueryRows(&account); err != nil {
			log.Println("请刷新后再次请求:", err)
			continue
		}

		a := utilModel.Key[account[0]]
		id, secrt := utilModel.Decrypt(a.ID, a.Secret)

		// 使用拿出来的账号请求阿里云
		// 看情况是否要对返回的 pushid 做处理
		mess = aliyumCdnRefres(id, secrt, url, RefresType)
		if mess != "刷新/预热成功" {
			break
		}
		//mess = append(mess, m)

	}
	//fmt.Println("mess:", m)
	return mess
}

// 调用阿里云 sdk 相关代码发起请求
func aliyumCdnRefres(AK, SK, URL string, RefresType string) (mess string) {
	client, _ := CdnCreateClient(tea.String(AK), tea.String(SK))

	refreshObjectCachesRequest := &cdn20180510.RefreshObjectCachesRequest{
		ObjectType: tea.String(RefresType),
		ObjectPath: tea.String(URL),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.RefreshObjectCachesWithOptions(refreshObjectCachesRequest, runtime)
	if err != nil {
		log.Println("缓存刷新接口调用失败:", err)
		return
	}

	if *resp.Body.RefreshTaskId == "" {
		return "刷新/预热失败"
	}

	return "刷新/预热成功"
	//return *resp.Body.RefreshTaskId

}

// 从mysql查找域名对应的阿里云账号
func selectAccountRefres(url string) (domain string) {
	d := strings.Split(url, "//")
	d1 := strings.Split(d[1], "/")
	domain = d1[0]
	fmt.Println(domain)
	return domain
}
