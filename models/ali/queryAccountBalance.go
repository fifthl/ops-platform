/*
 * @Author: nevin
 * @Date: 2023-11-13 14:08:16
 * @LastEditTime: 2023-11-13 15:23:23
 * @LastEditors: nevin
 * @Description: 财务相关
 */
package aliModel

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"log"
	"sort"
	"strconv"
	"strings"
	utilModel "yw_cloud/models/util"
)

/*
获取阿里云账号余额
*/

type RemainingSum struct {
	Count     float64 `json:"count"`
	AccountID string  `json:"accountID"`
}

func GetRemainingSum() []RemainingSum {
	runtime := &util.RuntimeOptions{}

	//Count: "无数值",
	bacRes := []RemainingSum{}
	rs := RemainingSum{}

	// 循环获取所有阿里云账号的余额，将汇总数据追加到 bacRes 中返回
	for _, v := range utilModel.Key {
		// 新解密
		ID, Secret := utilModel.Decrypt(v.ID, v.Secret)

		// 初始化 client
		client, _err := BillCreateClient(&ID, &Secret)
		if _err != nil {
			log.Println("余额 clinet 创建失败:", _err)
			return nil
		}
		_res, _err := client.QueryAccountBalanceWithOptions(runtime)
		if _err != nil {
			fmt.Println("请求失败: ", _err.Error())
		}

		// 获取到金额和账号id
		//
		AvailableAmountStr := strings.ReplaceAll(*_res.Body.Data.AvailableAmount, ",", "")
		availableAmount, _ := strconv.ParseFloat(AvailableAmountStr, 64)

		rs.Count = availableAmount
		rs.AccountID = v.AccountId

		// 将rs追加到返回值中
		bacRes = append(bacRes, rs)
	}
	//对余额类型转换后排序
	sort.SliceStable(bacRes, func(i, j int) bool {
		return bacRes[i].Count > bacRes[j].Count
	})
	return bacRes

}
