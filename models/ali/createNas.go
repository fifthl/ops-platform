package aliModel

import (
	"encoding/json"
	nas20170626 "github.com/alibabacloud-go/nas-20170626/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	utilModel "yw_cloud/models/util"
)

/*
函数计算创建 nas
*/

// nas 挂载点
const (
	FileSystemType = "standard"
	ChargeType     = "PayAsYouGo"
	StorageType    = "Performance"
	ProtocolType   = "NFS"
	ZoneId         = "cn-hangzhou-f"
)

// nas id
const (
	NetworkType     = "VPC"
	VpcId           = "vpc-bp1jdhjct17tfm6bstu8f"
	VSwitchId       = "vsw-bp18ogozj4w7x0m8xm03x"
	AccessGroupName = "DEFAULT_VPC_GROUP_NAME"
)

/*
1. 创建nas实例
2. 添加nas挂载点
*/
func createNas(description string) (NasId, MountDomain string, e error) {

	// 创建nas实例
	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, err := NasCreateClient(&id, &secret)
	if err != nil {
		log.Printf("fc create client failed, err:%v\n", err)
		return "", "", err
	}

	createFileSystemRequest := &nas20170626.CreateFileSystemRequest{
		FileSystemType: tea.String(FileSystemType),
		ChargeType:     tea.String(ChargeType),
		StorageType:    tea.String(StorageType),
		ProtocolType:   tea.String(ProtocolType),
		ZoneId:         tea.String(ZoneId),
		Description:    tea.String(description),
	}
	runtime := &util.RuntimeOptions{}

	resp, _err := client.CreateFileSystemWithOptions(createFileSystemRequest, runtime)
	if _err != nil {
		log.Printf("CreateFileSystemWithOptions failed, err:%v\n", _err)
		return "", "", _err
	}

	// 添加nas挂载点
	createMountTargetRequest := &nas20170626.CreateMountTargetRequest{
		FileSystemId:    tea.String(*resp.Body.FileSystemId),
		NetworkType:     tea.String(NetworkType),
		VpcId:           tea.String(VpcId),
		VSwitchId:       tea.String(VSwitchId),
		AccessGroupName: tea.String(AccessGroupName),
	}
	mountResp, mountErr := client.CreateMountTargetWithOptions(createMountTargetRequest, runtime)
	if mountErr != nil {
		log.Printf("创建挂载点失败%v", mountErr)
		return "", "", mountErr
	}

	return *resp.Body.FileSystemId, *mountResp.Body.MountTargetDomain, nil

}

// 查询nas状态
type MountTarget struct {
	Status string `json:"Status"`
}
type MountTargets struct {
	MountTarget []MountTarget `json:"MountTarget"`
}

func checkNasStat(nasId string) (status string, nasErr error) {
	tag := new(MountTargets)

	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := NasCreateClient(&id, &secret)

	describeMountTargetsRequest := &nas20170626.DescribeMountTargetsRequest{
		FileSystemId: tea.String(nasId),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.DescribeMountTargetsWithOptions(describeMountTargetsRequest, runtime)
	if err != nil {
		log.Println("nas状态查询失败%v\n", err)
		return "", err
	}

	if err = json.Unmarshal([]byte(resp.Body.MountTargets.String()), tag); err != nil {
		log.Printf("nas状态序列化失败%v\n", err)
		return "", err
	}

	return tag.MountTarget[0].Status, nil

}
