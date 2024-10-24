package aliModel

import (
	"fmt"
	cdn20180510 "github.com/alibabacloud-go/cdn-20180510/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"time"
	utilModel "yw_cloud/models/util"
)

/*
查询最近一小时的刷新/预热记录
*/

type RefreshTasks struct {
	Status       string `json:"Status,omitempty"`
	ObjectPath   string `json:"ObjectPath,omitempty"`
	TaskId       string `json:"TaskId,omitempty"`
	CreationTime string `json:"CreationTime,omitempty"`
	Process      string `json:"Process,omitempty"`
}

func DescribeRefreshTasks(AccountID string) (Tasks []RefreshTasks) {

	startTime := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02T15:04:05Z")
	endTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	fmt.Println(endTime)
	fmt.Println(startTime)
	account := utilModel.Key[AccountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := CdnCreateClient(&id, &secret)

	describeRefreshTasksRequest := &cdn20180510.DescribeRefreshTasksRequest{
		StartTime: tea.String(startTime),
		EndTime:   tea.String(endTime),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.DescribeRefreshTasksWithOptions(describeRefreshTasksRequest, runtime)
	if err != nil {
		log.Println("查询刷新/预热记录失败")
	}

	//如果查询的时间没有刷新记录，则返回空数组

	if len(resp.Body.Tasks.CDNTask) == 0 {
		log.Println("查询刷新/预热记录为空")
		return
	}

	for _, task := range resp.Body.Tasks.CDNTask {
		t := new(RefreshTasks)

		t.TaskId = *task.TaskId
		t.CreationTime = *task.CreationTime
		t.ObjectPath = *task.ObjectPath
		t.Process = *task.Process
		t.Status = *task.Status

		t.change()
		Tasks = append(Tasks, *t)
	}

	return Tasks

}

func (r *RefreshTasks) change() {
	switch r.Status {
	case "Complete":
		r.Status = "已完成"
	case "Refreshing":
		r.Status = "刷新中"
	case "Failed":
		r.Status = "失败"

	}
}
