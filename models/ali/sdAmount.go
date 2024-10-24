package aliModel

import (
	"encoding/json"
	"fmt"
	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	nas20170626 "github.com/alibabacloud-go/nas-20170626/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/beego/beego/v2/client/orm"
	"log"
	"strconv"
	"sync"
	"time"
	utilModel "yw_cloud/models/util"
)

/*
计算 nas
计算应用
计算ecs
*/

func NasAmount() {

	var (
		o       = orm.NewOrm()
		results = []string{}
		nasWg   sync.WaitGroup
	)

	if _, err := o.Raw("select nas_id from sd_designer").QueryRows(&results); err != nil {
		log.Printf("查询nas列表失败%v\n", err)
		return
	}

	if len(results) == 0 {
		log.Println("sd_designer 列表为空")
	}

	// 创建nas实例
	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, err := NasCreateClient(&id, &secret)
	if err != nil {
		log.Printf("fc create client failed, err:%v\n", err)
		return
	}

	// nas 费用
	log.Println("计算器nas金额")
	nasWg.Add(len(results))
	for _, nas := range results {
		go func(nasId string) {
			defer nasWg.Done()
			requestNasAmount(nasId, client, o)
		}(nas)
	}

	nasWg.Wait()

	// 资源包费用
	requestPackageAmount(id, secret, o)

	// ECS 费用
	RenewalPrice()
}

type FileSystem struct {
	Description  string  `json:"Description" orm:"column(description)"`
	FileSystemID string  `json:"FileSystemId" orm:"column(file_system_id)"`
	MeteredSize  float64 `json:"MeteredSize" orm:"column(metered_size)"`
}

func (f *FileSystem) TableName() string {
	return "sd_amount"
}

func requestNasAmount(nasId string, client *nas20170626.Client, o orm.Ormer) {
	describeFileSystemsRequest := &nas20170626.DescribeFileSystemsRequest{
		FileSystemId:   tea.String(nasId),
		FileSystemType: tea.String("standard"),
	}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.DescribeFileSystemsWithOptions(describeFileSystemsRequest, runtime)
	if _err != nil {
		log.Printf("请求nas: %v失败 %v: \n", nasId, _err)
		return
	}

	data := new(FileSystem)
	if err := json.Unmarshal([]byte(resp.Body.FileSystems.FileSystem[0].String()), data); err != nil {
		log.Printf("序列化FileSystem失败%v\n", err)
		return
	}
	data.MeteredSize, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", data.MeteredSize/float64(1024)/float64(1024)/float64(1024)), 64)
	data.MeteredSize, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", data.MeteredSize*1.85), 64)
	fmt.Println(data)
}

type PackageData struct {
	Instance []Instance `json:"Instance"`
}
type Instance struct {
	EffectiveTime   time.Time `json:"EffectiveTime"`
	ExpiryTime      time.Time `json:"ExpiryTime"`
	InstanceID      string    `json:"InstanceId"`
	PackageType     string    `json:"PackageType"`
	RemainingAmount string    `json:"RemainingAmount"`
	Remark          string    `json:"Remark"`
	Status          string    `json:"Status"`
	TotalAmount     string    `json:"TotalAmount"`
}

func requestPackageAmount(id, secret string, o orm.Ormer) {
	log.Println("计算资源包金额")
	client, err := BillCreateClient(&id, &secret)
	if err != nil {
		log.Printf("创建client失败\n", err)
		return
	}

	queryResourcePackageInstancesRequest := &bssopenapi20171214.QueryResourcePackageInstancesRequest{
		ProductCode:    tea.String("fc"),
		PageSize:       tea.Int32(300),
		PageNum:        tea.Int32(1),
		IncludePartner: tea.Bool(true),
	}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.QueryResourcePackageInstancesWithOptions(queryResourcePackageInstancesRequest, runtime)
	if _err != nil {
		log.Printf("查询资源包剩余量失败%v\n", _err)
		return
	}

	pd := new(PackageData)
	if err = json.Unmarshal([]byte(resp.Body.Data.Instances.String()), pd); err != nil {
		log.Printf("序列化资源包信息失败%v\n", err)
		return
	}

	for _, instance := range pd.Instance {
		fmt.Printf("%v 总计: %v 剩余:%v\n", instance.Remark, instance.TotalAmount, instance.RemainingAmount)
	}

}

// 杭州地域ecs
func RenewalPrice() (price float32) {
	log.Println("计算杭州 ECS 费用")
	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := EcsCreateClient(&id, &secret, "ecs-cn-hangzhou.aliyuncs.com")

	describeRenewalPriceRequest := &ecs20140526.DescribeRenewalPriceRequest{
		RegionId:     tea.String("cn-hangzhou"),
		ResourceType: tea.String("instance"),
		PriceUnit:    tea.String("Month"),
		Period:       tea.Int32(1),
		ResourceId:   tea.String("i-bp1dpmrg52av47ay626r"),
	}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.DescribeRenewalPriceWithOptions(describeRenewalPriceRequest, runtime)
	if _err != nil {
		log.Printf("查询 %v 续费接口失败: %v", "i-bp1dpmrg52av47ay626r", _err)
	}

	fmt.Println("ecs费用: ", *resp.Body.PriceInfo.Price.TradePrice)
	return *resp.Body.PriceInfo.Price.TradePrice

}
