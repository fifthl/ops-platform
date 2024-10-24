package cronAliyun

import (
	"context"
	"encoding/json"
	"fmt"
	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/astaxie/beego"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	aliModel "yw_cloud/models/ali"
	"yw_cloud/models/db"
	utilModel "yw_cloud/models/util"
)

var timeZone, _ = time.LoadLocation("Asia/Shanghai")

type weeklyCron struct {
	PaymentAmount float32 `json:"paymentAmount"`
}

const (
	DateFormat string = "2006-01"
)

// 相当于定时任务的入口函数
func CronWeeklyAccount() {
	log.Println("定时已添加")

	// 定时任务的相关配置
	c := cron.New(cron.WithLocation(timeZone))

	//上周账单, 周账单
	_, _ = c.AddFunc("15 1 * * MON", GetBeforAccount)

	//月分摊
	_, _ = c.AddFunc("57 23 * * *", cronOverview)

	// 账号消费
	_, _ = c.AddFunc("53 23 * * *", consume)

	// 每周更新实例 tag
	_, _ = c.AddFunc("18 2 * * *", getInstancesInfo)

	// 每两小时更新一次短信模板
	_, _ = c.AddFunc("0 */2 * * *", WriteSms)

	// 获取昨日余额
	_, _ = c.AddFunc("10 3 * * *", getAliyun)

	// 获取上月费用
	_, _ = c.AddFunc("25 2 1 * *", getAliyunMonthly)
	c.Start()
}

// 用于更改jenkins表中已有的发版记录

// 写到 db
func GetBeforAccount() {
	log.Println("周账单定时执行...")

	//循环账号
	for _, value := range utilModel.Key {

		// 每次循环账号清空历史数据
		bacWC := weeklyCron{}
		bacWCS := []weeklyCron{}
		// 密钥解密
		ID, Secret := utilModel.Decrypt(value.ID, value.Secret)

		// 循环一周费用
		for i := 1; i < 8; i++ {
			beforeDate := time.Now().In(timeZone).AddDate(0, 0, -i).Format("2006-01-02")
			beforeCycle := time.Now().In(timeZone).Format(DateFormat)

			queryAccountBillRequest := &bssopenapi20171214.QueryAccountBillRequest{
				Granularity:  tea.String("DAILY"),
				BillingCycle: tea.String(beforeCycle),
				BillingDate:  tea.String(beforeDate),
			}

			runtime := &util.RuntimeOptions{}

			// 创建 client
			client, _ := aliModel.BillCreateClient(&ID, &Secret)

			// 发起请求
			_res, _err1 := client.QueryAccountBillWithOptions(queryAccountBillRequest, runtime)
			if _err1 != nil {
				log.Println("_err QueryAccountBillWithOptions: ", _err1)
			}

			// 如果PaymentAmount长度为零，则查询日没有账单生成
			if len(_res.Body.Data.Items.Item) == 0 {
				bacWC.PaymentAmount = 0
			} else {
				bacWC.PaymentAmount = *(_res.Body.Data.Items.Item[0].PaymentAmount)
			}

			bacWCS = append(bacWCS, bacWC)

		}

		log.Println("周循环完成计算总和...")
		// 取出数组，将每天值相加后总和
		total := weeklyTotal(bacWCS)

		log.Println("计算到的过期时间: ", utilModel.TimeUntilNextWeekday())

		if err := db.RedisDb.Set(value.Name+"WeeklyBilling", total, utilModel.TimeUntilNextWeekday()).Err(); err != nil {
			log.Println("redis写入失败", err)
			continue
		}
		log.Println("写入redis成功")

	}

}

func weeklyTotal(weeklytotal []weeklyCron) (total float32) {
	total = 0

	for _, value := range weeklytotal {
		total += value.PaymentAmount
		log.Println("单日费用: ", value.PaymentAmount)
	}

	log.Println("费用总和: ", total)
	return total
}

// 补充上一个与的月分摊记录

type overviewS struct {
	PretaxAmount float32 `json:"PretaxAmount"`
	ProductName  string  `json:"ProductName"`
}

type OverviewSs struct {
	Item      []overviewS
	AccountId string `json:"AccountId"`
	Months    string `json:"Months"`
}

func cronOverview() {
	log.Println("月分摊记录定时执行...")
	var (
		lastMonth = time.Now().Format(DateFormat)
		ctx       = context.Background()
	)

	type overviewS struct {
		PretaxAmount float32 `json:"PretaxAmount"`
		ProductName  string  `json:"ProductName"`
	}

	type OverviewSs struct {
		Item      []overviewS
		AccountId string `json:"AccountId"`
		Months    string `json:"Months"`
	}

	bacov := OverviewSs{}
	// 循环7个账号的上一个月
	for _, v := range utilModel.Key {
		ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
		// 初始化 client
		client, _ := aliModel.BillCreateClient(tea.String(ID), tea.String(Secret))

		queryBillOverviewRequest := &bssopenapi20171214.QueryBillOverviewRequest{
			BillingCycle: tea.String(lastMonth),
		}
		runtime := &util.RuntimeOptions{}

		resp, err := client.QueryBillOverviewWithOptions(queryBillOverviewRequest, runtime)
		if err != nil {
			log.Println("cronOver 请求出错:", err)
		}
		if uerr := json.Unmarshal([]byte(resp.Body.Data.Items.GoString()), &bacov); uerr != nil {
			log.Println("err: ", uerr.Error())
		}

		// 追加账号ID和月份
		bacov.AccountId = v.AccountId
		bacov.Months = lastMonth

		collection := db.GetCollection("overview")

		delResult, _ := collection.DeleteMany(ctx, bson.D{{"accountid", bacov.AccountId}, {"months", bacov.Months}})
		log.Println("删除的当前月分摊 id:", delResult.DeletedCount)

		//插入
		result, err := collection.InsertOne(ctx, &bacov)
		if err != nil {
			log.Println("mongo写入失败: ", err)
			return
		}
		log.Printf("cronOverview 完成, ID: %v", result.InsertedID)
	}

}

/*
账号消费
*/
type consuField struct {
	PretaxAmount float32 `json:"PretaxAmount"`
	Months       string  `json:"Months"`
	Year         string  `json:"Year"`
	AccountId    string  `json:"AccountId"`
}

// Billling -> QueryAccountBill
func consume() {
	log.Println("阿里云账号消费定时执行")
	defer log.Println("阿里云账号消费定时执行完成")

	// 当前年月
	m := time.Now().Format("01")
	y := time.Now().Format("2006")

	consumeColl := db.GetCollection("consume")

	// 循环账号
	for _, v := range utilModel.Key {

		ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
		client, _err := aliModel.BillCreateClient(&ID, &Secret)
		if _err != nil {
			log.Println("_err BillCreateClient: ", _err)
		}

		var (
			cf = consuField{}
		)

		consumeDate := y + "-" + m
		queryAccountBillRequest := &bssopenapi20171214.QueryAccountBillRequest{
			BillingCycle: tea.String(consumeDate),
		}
		runtime := &util.RuntimeOptions{}

		_res, err := client.QueryAccountBillWithOptions(queryAccountBillRequest, runtime)
		if err != nil {
			log.Println("账号消费定时查询失败: ", err)
			return
		}

		if len(_res.Body.Data.Items.Item) == 0 {
			cf.PretaxAmount = 0
		} else {
			cf.PretaxAmount = *(_res.Body.Data.Items.Item[0].PretaxAmount)
		}

		cf.Months = m
		cf.AccountId = v.AccountId
		cf.Year = y

		opts := options.Update().SetUpsert(true)
		filter := bson.D{{"year", y},
			{"months", m},
			{"accountid", cf.AccountId}}

		update := bson.D{{"$set", bson.D{
			{"accountid", cf.AccountId},
			{"months", cf.Months},
			{"pretaxamount", cf.PretaxAmount},
			{"year", cf.Year}}}}

		if _, _er := consumeColl.UpdateOne(context.TODO(), filter, update, opts); _er != nil {
			log.Println("更新账号消费记录失败:", _er)
		}
	}

}

/*
全局 mysql client
*/
func init() {
	if err := orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql_addr")); err != nil {
		log.Println("注册mysql失败: ", err)
		return
	}

	if err := orm.RegisterDataBase("ecs", "mysql", beego.AppConfig.String("mysql_addr_ecs")); err != nil {
		log.Println("注册mysql失败: ", err)
		return
	}

	orm.RegisterModel(new(Instances))

}

/*
获取 ECS 实例信息写入 mysql
*/

// 查询 ecs 实例详细信息
type Instances struct {
	InstanceId       string    `json:"InstanceId,omitempty" orm:"column(instance_id);pk;unique;description(实例 ID)"` // 实例 ID
	InstanceName     string    `json:"InstanceName" orm:"column(instance_name);description(实例名称)"`                  // 实例名称
	HostName         string    `json:"HostName" orm:"column(host_name);description(实例主机名)"`                         // 实例主机名
	OSName           string    `json:"OSName" orm:"column(os_name);description(实例的操作系统名称)"`                         // 实例的操作系统名称
	Cpu              int32     `json:"Cpu" orm:"column(cpu);description(vCPU 数)"`                                   // vCPU 数
	Memory           int32     `json:"Memory" orm:"column(memory);description(内存大小，单位MiB)"`                         // 内存大小，单位为 MiB
	GPUAmount        int32     `json:"GPUAmount" orm:"column(gpu);description(GPU 数量)"`                             // 实例规格附带的 GPU 数量
	PrimaryIpAddress string    `json:"PrimaryIpAddress" orm:"column(primary_ip_address);description(弹性网卡私有IP地址)"`   //弹性网卡私有IP地址
	PublicIpAddress  string    `json:"PublicIpAddress" orm:"column(public_ip_address);description(公网IP)"`           // 公网
	EipAddress       string    `json:"eipAddress" orm:"column(eip_address);description(弹性IP)"`                      // 弹性IP
	InstanceType     string    `json:"InstanceType,omitempty" orm:"column(instance_type);description(实例规格)"`        // 实例规格
	ZoneId           string    `json:"ZoneId" orm:"column(zone_id);description(所属地域)"`                              // 所属地域
	OSType           string    `json:"OSType" orm:"column(os_type);description(操作系统类型)"`                            // 操作系统类型
	VSwitchId        string    `json:"VSwitchId" orm:"column(vpc_attributes);description(交换机ID)"`                   //虚拟交换机 ID
	VpcId            string    `json:"VpcId" orm:"column(vpc_id);description(VPC ID)"`                              //专有网络 VPC ID
	Tags             string    `json:"Tags" orm:"column(tags);description(标签)"`
	RenewalStatus    string    `json:"InstanceChargeType" orm:"column(renewal_status)"`
	Price            float32   `json:"Price" orm:"column(price);description(续费金额)"`
	Created          time.Time `orm:"auto_now_add;type(datetime);description(更新时间)"`
}

func (*Instances) TableName() string {
	return "instances_info"
}

func getInstancesInfo() {
	log.Println("Cron/更新ECS实例信息...")
	o := orm.NewOrmUsingDB("ecs")
	_, _err := o.Raw("DELETE FROM instances_info").Exec()
	fmt.Println(_err)

	accountID := [3]string{"1", "2", "3"}

	EndpointRegionId := map[string][]string{
		//"xx": {"Endpoint", "RegionId"}
		"cn": {"ecs.cn-beijing.aliyuncs.com", "cn-beijing"},
		"us": {"ecs.us-east-1.aliyuncs.com", "us-east-1"},
		"hz": {"ecs-cn-hangzhou.aliyuncs.com", "cn-hangzhou"},
	}

	bjRegion := EndpointRegionId["cn"]
	eastRegion := EndpointRegionId["us"]
	hzRegion := EndpointRegionId["hz"]

	//循环账号
	for _, account := range accountID {

		// 创建 client
		client := newClient(account, bjRegion[0])

		// 1.获取北京地域ECS
		// 2.写入db
		reqCnAndUsAndHzInstancesInfo(client, bjRegion[1])

	}

	// 1.获取漂亮国地域ECS
	// 2.写入db
	client := newClient("1", eastRegion[0])
	reqCnAndUsAndHzInstancesInfo(client, eastRegion[1])

	hzClient := newClient("1", hzRegion[0])
	reqCnAndUsAndHzInstancesInfo(hzClient, hzRegion[1])

	log.Println("Cron/更新ECS实例信息完成")

}

// 创建 client
func newClient(account string, Endpoint string) (client *ecs20140526.Client) {
	idEnc, secretEnv := utilModel.Key[account].ID, utilModel.Key[account].Secret

	id, secret := aliModel.IdSecret(idEnc, secretEnv)

	client, _ = aliModel.EcsCreateClient(&id, &secret, Endpoint)

	// 返回 client
	return client
}

// 获取北京和弗尼及亚地域 ECS 信息
func reqCnAndUsAndHzInstancesInfo(client *ecs20140526.Client, RegionId string) {

	for Number := 1; Number < 3; Number++ {
		describeInstancesRequest := &ecs20140526.DescribeInstancesRequest{
			RegionId:   tea.String(RegionId),
			PageSize:   tea.Int32(100),
			PageNumber: tea.Int32(int32(Number)),
			Status:     tea.String("Running"),
		}
		runtime := &util.RuntimeOptions{}
		res, err := client.DescribeInstancesWithOptions(describeInstancesRequest, runtime)
		if err != nil {
			log.Println(err)
			return
		}

		// 数据需要存入db中，
		//后续的 续费接口，tag接口，负载均衡接口需要依赖这些数据
		o := orm.NewOrmUsingDB("ecs")
		instancesInfo := new(Instances)

		if len(res.Body.Instances.Instance) == 0 {
			log.Println("第二页无数据")
			continue
		}

		//err := o.Read()
		for _, instance := range res.Body.Instances.Instance {

			// 赋值
			instancesInfo.InstanceId = *instance.InstanceId
			instancesInfo.InstanceName = *instance.InstanceName
			instancesInfo.HostName = *instance.HostName
			instancesInfo.OSName = *instance.OSName
			instancesInfo.Memory = *instance.Memory / 1024
			instancesInfo.Cpu = *instance.Cpu
			instancesInfo.GPUAmount = *instance.GPUAmount
			instancesInfo.EipAddress = *instance.EipAddress.IpAddress
			instancesInfo.InstanceType = *instance.InstanceType
			instancesInfo.ZoneId = *instance.ZoneId
			instancesInfo.OSType = *instance.OSType
			instancesInfo.VSwitchId = *instance.VpcAttributes.VSwitchId
			instancesInfo.VpcId = *instance.VpcAttributes.VpcId

			instancesInfo.PrimaryIpAddress = *instance.NetworkInterfaces.NetworkInterface[0].PrimaryIpAddress
			instancesInfo.RenewalStatus = DescribeInstanceAutoRenewAttribute(client, *instance.InstanceId, RegionId) // 获取续费状态
			instancesInfo.PublicIpAddress = getPublicIpAddress(instance.PublicIpAddress)

			if instance.Tags == nil {
				fmt.Println("instance.Tags 为空")
				instancesInfo.Tags = ""
			} else {
				instancesInfo.Tags = getTag(instance.Tags)
			}

			instancesInfo.Price = GetRenewal(client, *instance.InstanceId, RegionId) // 获取续费金额

			if _, err = o.Insert(instancesInfo); err != nil {
				log.Println(err)
				continue
			}
			//fmt.Printf("insert data %v", instancesInfo)
		}
	}

}

// 获取ECS续费金额
// 需要在两个地域的ECS都更新完之后在执行
func GetRenewal(client *ecs20140526.Client, InstanceId, RegionId string) (TradePrice float32) {
	defer func() {
		if _err := recover(); _err != nil {
			log.Println("续费接口查询失败")
		}
	}()

	describeRenewalPriceRequest := &ecs20140526.DescribeRenewalPriceRequest{
		RegionId:     tea.String(RegionId),
		ResourceType: tea.String("instance"),
		PriceUnit:    tea.String("Month"),
		Period:       tea.Int32(1),
		ResourceId:   tea.String(InstanceId),
	}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.DescribeRenewalPriceWithOptions(describeRenewalPriceRequest, runtime)
	if _err != nil {
		log.Printf("查询 %v 续费接口失败: %v", InstanceId, _err)
	}

	return *resp.Body.PriceInfo.Price.TradePrice
}

// ecs 标签
func getTag(client *ecs20140526.DescribeInstancesResponseBodyInstancesInstanceTags) (tag string) {
	for _, t := range client.Tag {
		// 不保存这些标签
		if *t.TagKey == "Linux" || *t.TagKey == "Test" || *t.TagKey == "Win" || *t.TagKey == "acs:ecs:supportVtpm : true" || *t.TagKey == "acs:ecs:supportVtpm" {
			continue
		}
		return *t.TagKey
	}
	return tag
}

// ecs 公网ip
func getPublicIpAddress(client *ecs20140526.DescribeInstancesResponseBodyInstancesInstancePublicIpAddress) (IpAddress string) {
	index := len(client.IpAddress)

	if index != 0 {
		return *client.IpAddress[0]
	}
	return ""

}

// ecs 续费状态
func DescribeInstanceAutoRenewAttribute(client *ecs20140526.Client, InstanceId, RegionId string) (RenewalStatus string) {
	describeInstanceAutoRenewAttributeRequest := &ecs20140526.DescribeInstanceAutoRenewAttributeRequest{
		RegionId:   tea.String(RegionId),
		InstanceId: tea.String(InstanceId),
	}
	runtime := &util.RuntimeOptions{}

	res, err := client.DescribeInstanceAutoRenewAttributeWithOptions(describeInstanceAutoRenewAttributeRequest, runtime)
	if err != nil {
		log.Println(err)
		return "获取续费状态失败"
	}

	switch *res.Body.InstanceRenewAttributes.InstanceRenewAttribute[0].RenewalStatus {
	case "AutoRenewal":
		RenewalStatus = "自动续费"
	case "Normal":
		RenewalStatus = "手动续费"
	case "NotRenewal":
		RenewalStatus = "不再续费"
	}

	return RenewalStatus

}

/*
获取短信模板，写入 mysql
*/
type Reason struct {
	RejectInfo string `json:"RejectInfo"`
	RejectDate string `json:"RejectDate"`
}

type SmsTemplateList struct {
	AuditStatus       string `json:"AuditStatus,omitempty"`
	CreateDate        string `json:"CreateDate,omitempty"`
	OuterTemplateType int32  `json:"OuterTemplateType"`
	Reason            Reason `json:"Reason" orm:"column(reason)"`
	TemplateCode      string `json:"TemplateCode,omitempty"`
	TemplateContent   string `json:"TemplateContent,omitempty"`
	TemplateName      string `json:"TemplateName,omitempty"`
	AccountID         string `json:"AccountID,omitempty"`
}

//func (s *SmsTemplateList) TableName() string {
//	return "sms_template"
//}

// 获取阿里云短信列表，写入数据库
func WriteSms() {
	log.Println("启动更新短信模板定时...")

	for _, v := range utilModel.Key {

		ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
		client, _ := aliModel.SmsCreateClient(tea.String(ID), tea.String(Secret))

		querySmsTemplateListRequest := &dysmsapi20170525.QuerySmsTemplateListRequest{
			PageSize: tea.Int32(50),
		}
		runtime := &util.RuntimeOptions{}

		_resp, _err := client.QuerySmsTemplateListWithOptions(querySmsTemplateListRequest, runtime)
		if _err != nil {
			log.Printf("%v请求sms模板列表失败: %v", v.AccountId, _err)
			return
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
			}
		}
	}
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

/*
获取昨日余额
*/

type date struct {
	Cycle string
	Date  string
}

// GetDate 获取查询账单的时间
func getDate() *date {
	// 指定时间格式
	beforeDate := time.Now().In(timeZone).AddDate(0, 0, -1).Format("2006-01-02")
	beforeCycle := time.Now().In(timeZone).Format(DateFormat)

	return &date{
		Cycle: beforeCycle,
		Date:  beforeDate,
	}
}

type Account struct {
	PaymentAmount float32 `json:"paymentAmount"`
	AccountID     string  `json:"accountID"`
}

// 阿里云昨日费用接口
func getAliyun() {
	log.Println("获取昨日费用定时启动...")
	for _, v := range utilModel.Key {
		aliAct := Account{}
		t := getDate()

		queryAccountBillRequest := &bssopenapi20171214.QueryAccountBillRequest{
			Granularity:  tea.String("DAILY"),
			BillingCycle: tea.String(t.Cycle),
			BillingDate:  tea.String(t.Date),
		}

		runtime := &util.RuntimeOptions{}

		// 调用阿里云api
		ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
		// 创建客户端
		client, _err := aliModel.BillCreateClient(&ID, &Secret)

		if _err != nil {
			log.Println("err: ", _err)
			//return _err
		}

		// 发起请求
		_res, _ := client.QueryAccountBillWithOptions(queryAccountBillRequest, runtime)

		// 如果PaymentAmount长度为零，则查询日没有账单生成
		if len(_res.Body.Data.Items.Item) == 0 {
			aliAct.PaymentAmount = 0
		} else {
			aliAct.PaymentAmount = *(_res.Body.Data.Items.Item[0].PaymentAmount)
		}

		// 插入redis
		if err := db.RedisDb.Set(v.Name+"YesterDay", aliAct.PaymentAmount, utilModel.DaySeconds()).Err(); err != nil {
			log.Println("redis写入失败: ", err)
		}
	}

}

// 月账单
type dateMonthly struct {
	Cycle string
}

// GetDate 获取查询账单的时间
func getDateMonthly() *dateMonthly {
	// 指定时间格式
	beforeCycle := time.Now().In(timeZone).AddDate(0, -1, 0).Format(DateFormat)

	return &dateMonthly{
		Cycle: beforeCycle,
	}
}

type MonthlyBilling struct {
	PaymentAmount float32 `json:"paymentAmount"`
	AccountID     string  `json:"accountID"`
}

func getAliyunMonthly() {
	log.Println("获取月账单定时启动...")
	for _, v := range utilModel.Key {
		t := getDateMonthly()
		aliyunMb := MonthlyBilling{}

		queryAccountBillRequest := &bssopenapi20171214.QueryAccountBillRequest{
			Granularity:  tea.String("MONTHLY"),
			BillingCycle: tea.String(t.Cycle),
		}

		runtime := &util.RuntimeOptions{}

		// 调用阿里云api
		ID, Secret := utilModel.Decrypt(v.ID, v.Secret)
		// 创建客户端
		client, _err := aliModel.BillCreateClient(&ID, &Secret)
		if _err != nil {
			log.Println(_err)
		}

		_res, _ := client.QueryAccountBillWithOptions(queryAccountBillRequest, runtime)

		// 如果PaymentAmount长度为零，则查询日没有账单生成
		if len(_res.Body.Data.Items.Item) == 0 {
			aliyunMb.PaymentAmount = 0
		} else {
			aliyunMb.PaymentAmount = *(_res.Body.Data.Items.Item[0].PaymentAmount)
		}

		aliyunMb.AccountID = v.AccountId

		// 插入redis
		if err := db.RedisDb.Set(v.Name+"MonthlyBilling", aliyunMb.PaymentAmount, utilModel.TimeUntilNextWeekday()).Err(); err != nil {
			log.Println("redis写入失败", err)
		}

	}

}
