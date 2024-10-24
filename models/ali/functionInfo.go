package aliModel

import (
	"encoding/json"
	"fmt"
	fc20230330 "github.com/alibabacloud-go/fc-20230330/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/beego/beego/v2/client/orm"
	"log"
	"strings"
	utilModel "yw_cloud/models/util"
)

func init() {
	orm.RegisterModel(new(Functions))
}

type RespData struct {
	Functions []Functions `json:"functions"`
}
type Functions struct {
	Stat          string `json:"stat"`
	BpmId         string `json:"bpmId"`
	FunctionShort string `json:"functionShort"`
	FunctionName  string `json:"functionName"`
	Description   string `json:"description"`
	MemorySize    int    `json:"memorySize"`
	Cpu           int    `json:"cpu"`
	Date          string `json:"date"`
}

func (f Functions) TableName() string {
	return "sd_app"
}

func getToken(currentToken string) (nextToken string, result bool, _err error) {
	defer func() {
		if err := recover(); err != nil {
			result = true
			fmt.Println("最后一页")
		}

	}()

	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := FcCreateClient(&id, &secret)
	listFunctionsRequest := &fc20230330.ListFunctionsRequest{
		Limit:     tea.Int32(100),
		NextToken: tea.String(currentToken),
	}
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)

	resp, _err := client.ListFunctionsWithOptions(listFunctionsRequest, headers, runtime)
	if _err != nil {
		log.Printf("参数配置失败 %v", _err)
		return "", false, _err
	}

	return *resp.Body.NextToken, false, nil
}

func getFc(nextToken string) *RespData {
	var (
		data = new(RespData)
	)

	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := FcCreateClient(&id, &secret)

	listFunctionsRequest := &fc20230330.ListFunctionsRequest{
		Limit:     tea.Int32(100),
		NextToken: tea.String(nextToken),
	}
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)

	resp, _err := client.ListFunctionsWithOptions(listFunctionsRequest, headers, runtime)
	if _err != nil {
		log.Printf("参数配置失败 %v", _err)
		return nil
	}

	if err := json.Unmarshal([]byte(resp.Body.GoString()), data); err != nil {
		log.Printf("反序列化失败 %v", err)
		return nil
	}
	return data
}

func ListFunctions() (fcList []string, result bool, err error) {
	var (
		nextToken = ""
		data      = []Functions{}
	)

	/*
		获取全部函数列表需要拿到 nextToken
	*/
	for {
		fc := getFc(nextToken)
		// 获取 FunctionShort
		for i := 0; i < len(fc.Functions); i++ {
			fc.Functions[i].FunctionShort = strings.Split(fc.Functions[i].FunctionName, "__")[0]
		}

		// 全部追加到一个大切片上
		data = append(data, fc.Functions...)

		nextToken, result, _ = getToken(nextToken)

		if nextToken == "" && result == true {
			break
		}
	}

	o := orm.NewOrm()
	r := []string{}

	for _, functions := range data {
		updateSql := fmt.Sprintf("INSERT IGNORE INTO sd_app (stat,bpm_id,function_short,function_name,description,memory_size,cpu,date) VALUES ('0', '', '%v','%v','%v','%v','%v','');", functions.FunctionShort, functions.FunctionName, functions.Description, functions.MemorySize, functions.Cpu)

		if _, _err := o.Raw(updateSql).QueryRows(&r); _err != nil {
			log.Printf("插入失败%v \n", _err)
			return
		}
	}
	return nil, false, nil

}

func fcShort(str string) string {
	f := strings.Split(str, "__")
	return f[0]
}

/*
functionIsUsed 判断应用是否空闲（表中的0或者1）
*/
func functionIsUsed(o orm.Ormer, bpmId string) (fcID []string, err error) {
	// 0 代表可用
	// 1 代表在用
	// 函数默认都是 0，当进来一个用户后，更改 nas 地址，在更改状态为 1
	var respStat []string
	var usedStat []string

	queryUsed := fmt.Sprintf("select function_name from sd_app where stat='1' and bpm_id='%v';", bpmId)

	if _, err = o.Raw(queryUsed).QueryRows(&usedStat); err != nil {
		return nil, err
	}

	queryStat := "select function_name from sd_app where function_short=(select function_short from sd_app where stat='0' group by function_short limit 1);"

	if _, err = o.Raw(queryStat).QueryRows(&respStat); err != nil {
		return nil, err
	}

	return respStat, nil

}

// 判断当前账号是否有在用的应用
func functionExist(o orm.Ormer, bpmId string) (bool, error) {
	var usedStat []string

	queryUsed := fmt.Sprintf("select function_name from sd_app where stat='1' and bpm_id='%v';", bpmId)

	if _, err := o.Raw(queryUsed).QueryRows(&usedStat); err != nil {
		return false, err
	}

	// 大于0 证明当前账号已经绑定过了
	if len(usedStat) > 0 {
		log.Printf("%v当前已绑定应用\n", bpmId)
		return true, nil
	}

	return false, nil

}

//func QueryFcStat() (int, int, error) {
//	o := orm.NewOrm()
//	query := "select count(stat) from sd_app where stat='0';"
//	var unUse = []int{}
//	var use = []int{}
//
//	if _, err := o.Raw(query).QueryRows(&unUse); err != nil {
//		log.Printf("查询函数可用状态失败%v", err)
//		return 0, 0, err
//	}
//
//	if _, err := o.Raw("select count(stat) from sd_app where stat='1';").QueryRows(&use); err != nil {
//		log.Printf("查询函数可用状态失败%v", err)
//		return 0, 0, err
//	}
//
//	if len(unUse) > 0 && len(use) > 0 {
//		return unUse[0] / 4, use[0] / 4, nil
//	}
//
//	return 0, 0, nil
//}
