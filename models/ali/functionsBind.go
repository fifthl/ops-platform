package aliModel

import (
	"errors"
	"fmt"
	fc20230330 "github.com/alibabacloud-go/fc-20230330/v4/client"
	nas20170626 "github.com/alibabacloud-go/nas-20170626/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
	"strings"
	"time"
	utilModel "yw_cloud/models/util"
)

func init() {
	orm.RegisterModel(new(SdDesigner))
}

type SdDesigner struct {
	Name         string    `json:"name,omitempty"`
	BpmId        string    `json:"bpmId,omitempty"`
	Company      string    `json:"company,omitempty"`
	DesignCenter string    `json:"designCenter,omitempty"`
	NasId        string    `json:"nasID,omitempty"`
	MountDomain  string    `json:"mountDomain,omitempty"`
	Date         time.Time `json:"date,omitempty" orm:"auto_now_add;type(datetime)"`
}

func (d *SdDesigner) TableName() string {
	return "sd_designer"
}

/*
BindingNas
应用和nas绑定接口
*/
func BindingNas(name, company, designCenter, bpmId string) (d, user, pass, mount, date, manDoamin string, bErr error) {
	var (
		o     = orm.NewOrm()
		nasid = []string{}
		query = fmt.Sprintf("select nas_id from sd_designer where bpm_id='%s' and name='%s' order by  date desc limit 1;", bpmId, name)
	)

	designerInfo := SdDesigner{
		Name:         name,
		Company:      company,
		DesignCenter: designCenter,
		BpmId:        bpmId,
	}

	// 查询数据库中nas和用户是否有绑定关系
	if _, err := o.Raw(query).QueryRows(&nasid); err != nil {
		log.Printf("nasId查询失败,err %v\n", err)
		return "", "", "", "", "", "", err
	}

	if len(nasid) <= 0 {
		log.Printf("%v无绑定的nas\n", name)

		newNasId, mountDomain, err := createNas(name + "-" + bpmId)
		if err != nil {
			log.Printf("nas创建失败 %v\n", err)
			return "", "", "", "", "", "", err
		}
		designerInfo.NasId = newNasId
		designerInfo.MountDomain = mountDomain

		_, insertErr := o.InsertOrUpdate(&designerInfo)
		if insertErr != nil {
			log.Printf("nas记录写入失败: %v\n", insertErr)
			return "", "", "", "", "", "", err
		}

		// 查询nas状态
		for {
			status, stauErr := checkNasStat(designerInfo.NasId)
			if stauErr != nil {
				return "", "", "", "", "", "", stauErr
			}

			log.Printf("nas 状态为 %v\n", status)
			if status == "Active" {
				break
			}
			time.Sleep(2 * time.Second)
		}

		// 绑定用户 nas 目录
		dirBind(bpmId, mountDomain, o)

	} else {
		log.Printf("%v已创建过nas\n", bpmId)
		var domain = ""
		queryMountDomain := fmt.Sprintf("select mount_domain from sd_designer where nas_id='%v';", nasid[0])
		err := o.Raw(queryMountDomain).QueryRow(&domain)
		if err != nil {
			log.Printf("查询domain失败%v\n", err)
			return "", "", "", "", "", "", err
		}
		designerInfo.NasId = nasid[0]
		designerInfo.MountDomain = domain
	}

	//判断函数应用有没有空闲
	functionIds, err := functionIsUsed(o, bpmId)
	if err != nil {
		log.Printf("函数应用状态查询失败%v\n", err)
		return "", "", "", "", "", "", err
	}

	b, existErr := functionExist(o, bpmId)
	if existErr != nil {
		log.Printf("functionExist 函数应用状态查询失败%v\n", err)
		return "", "", "", "", "", "", err
	}

	if !b {
		// 如果长度小于零，证明没有空闲应用
		if fc := len(functionIds); fc <= 0 {
			return "无空闲应用", "", "", "", "", "", err
		}

		// 否则绑定函数和 nas
		for _, functionId := range functionIds {
			log.Printf("%v nas挂载点为: %v\n", functionId, designerInfo.MountDomain)
			if strings.Contains(functionId, "__sd") {
				sdBind(functionId, designerInfo.BpmId, designerInfo.MountDomain, o)
			} else {
				Bind(functionId, designerInfo.BpmId, designerInfo.MountDomain, o)
			}

		}

	}

	// 需再次判断函数是不是有空闲
	if fc := len(functionIds); fc <= 0 {
		return "无空闲应用", "", "", "", "", "", err
	}

	// 新增一个函数查询当前在用的nas，web 账号密码等
	d, user, pass, mount, date, manDoamin, err = queryOtherInfo(bpmId, o)
	if err != nil {
		log.Printf("查询sd webui错误%v", err)
		return "", "", "", "", "", "", err
	}

	// 有可能出现在 sd_web 表中没有查询到 web 信息
	if d == "" || user == "" || pass == "" {
		return "", "", "", "", "", "", errors.New("未查询到webui地址")
	}

	// 180秒
	for i := 0; i < 36; i++ {
		httpCode := checkSdStatus(d, bpmId)

		if httpCode == 200 {
			log.Printf("%v已可用\n\n", d)
			break
		}

		time.Sleep(5 * time.Second)
	}

	//// 清理extensions（插件）目录
	//config := &ssh.ClientConfig{
	//	Timeout:         time.Second,
	//	User:            "root",
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//	Auth:            []ssh.AuthMethod{ssh.Password("Q1UiTnP966Wrlk_L6IuysFAGP")},
	//}
	//
	//addr := fmt.Sprintf("%v:%v", "172.16.197.11", "22")
	//sshClient, err := ssh.Dial("tcp", addr, config)
	//if err != nil {
	//	log.Printf("ssh 主机失败:%v\n", err)
	//	return
	//}
	//
	//defer sshClient.Close()
	//
	//// 清空和复制插件
	//sessionExtensions, err := sshClient.NewSession()
	//if err != nil {
	//	log.Printf("创建ssh session 失败", err)
	//	return
	//}
	//defer sessionExtensions.Close()
	//clearExtensionsCmd := fmt.Sprintf("cd /data/%v/fc-stable-diffusion-plus/sd/ && rm -rf extensions/* && cp -a /data/SD/extensions/* /data/%v/fc-stable-diffusion-plus/sd/extensions/", bpmId, bpmId)
	//
	//log.Println("清空插件 && 复制插件")
	//a, bindErr := sessionExtensions.CombinedOutput(clearExtensionsCmd)
	//
	//if bindErr != nil {
	//	log.Println(string(a))
	//	return
	//}
	//
	//// 复制大模型（软接）
	//sessionLager, err := sshClient.NewSession()
	//if err != nil {
	//	log.Printf("创建ssh session 失败", err)
	//	return
	//}
	//defer sessionExtensions.Close()
	//cpLagerCmd := fmt.Sprintf("cp -a /data/largeModel/* /data/%v/fc-stable-diffusion-plus/sd/models/Stable-diffusion/", bpmId)
	//
	//log.Println("复制全部大模型")
	//a, bindErr = sessionLager.CombinedOutput(cpLagerCmd)
	//
	//if bindErr != nil {
	//	log.Println(string(a))
	//	return
	//}

	return d, user, pass, mount, date, manDoamin, err

}

func queryOtherInfo(bpmId string, o orm.Ormer) (domain, user, pass, mount, date, manDomain string, sdErr error) {
	// 查询sd域名，账号密码
	querySd := fmt.Sprintf("select sd_domain,sd_user,sd_passwd from sd_web where sd_id=(select function_short from sd_app where bpm_id='%v' limit 1);", bpmId)

	// 查询nas挂载点
	queryNas := fmt.Sprintf("select mount_domain from sd_designer where bpm_id='%v' limit 1;", bpmId)

	// 查询date
	queryDate := fmt.Sprintf("select date from sd_app where bpm_id='%v' limit 1;", bpmId)

	var (
		sdDomain = []string{}
		SdUser   = []string{}
		SdPasswd = []string{}
		SdMount  = []string{}
		SdDate   = []string{}
	)

	if _, err := o.Raw(querySd).QueryRows(&sdDomain, &SdUser, &SdPasswd); err != nil {
		log.Printf("查询nas信息失败%v\n", err)
		return "", "", "", "", "", "", err
	}

	if _, err := o.Raw(queryNas).QueryRows(&SdMount); err != nil {
		log.Printf("查询sd信息失败%v\n", err)
		return "", "", "", "", "", "", err
	}

	if _, err := o.Raw(queryDate).QueryRows(&SdDate); err != nil {
		log.Printf("查询 date 失败%v\n", err)
		return "", "", "", "", "", "", err
	}
	if len(sdDomain) == 0 {
		return "", "", "", "", "", "", errors.New("未查询到sdDomain")
	}
	managerDomain := strings.Split(sdDomain[0], ".")

	return sdDomain[0], SdUser[0], SdPasswd[0], SdMount[0], SdDate[0], managerDomain[0] + "-1" + ".voglassdc.com", nil
}

/*
绑定nas和应用
*/

var upDate = time.Now().Format("2006-01-02 15:04:05")

func sdBind(functionID, bpmId, mountDomain string, o orm.Ormer) {

	// 绑定逻辑
	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := FcCreateClient(&id, &secret)

	getFunctionRequest := &fc20230330.GetFunctionRequest{}
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)

	fcResp, _ := client.GetFunctionWithOptions(tea.String(functionID), getFunctionRequest, headers, runtime)

	// 绑定函数和nas
	updateFunctionInputGPUConfig := &fc20230330.GPUConfig{
		GpuMemorySize: tea.Int32(*fcResp.Body.GpuConfig.GpuMemorySize),
		GpuType:       tea.String(*fcResp.Body.GpuConfig.GpuType),
	}

	updateFunctionInputNASConfigNASMountConfig0 := &fc20230330.NASMountConfig{
		EnableTLS:  tea.Bool(true),
		MountDir:   tea.String("/mnt/auto"),
		ServerAddr: tea.String(mountDomain + ":/fc-stable-diffusion-plus"),
	}
	updateFunctionInputNASConfig := &fc20230330.NASConfig{
		MountPoints: []*fc20230330.NASMountConfig{updateFunctionInputNASConfigNASMountConfig0},
	}
	updateFunctionInputVPCConfig := &fc20230330.VPCConfig{
		SecurityGroupId: tea.String("sg-bp1a8vtwsdkumeb0w01o"),
		VSwitchIds:      []*string{tea.String("vsw-bp18ogozj4w7x0m8xm03x")},
		VpcId:           tea.String("vpc-bp1jdhjct17tfm6bstu8f"),
	}

	updateFunctionInput := &fc20230330.UpdateFunctionInput{
		GpuConfig: updateFunctionInputGPUConfig,
		NasConfig: updateFunctionInputNASConfig,
		VpcConfig: updateFunctionInputVPCConfig,
	}
	updateFunctionRequest := &fc20230330.UpdateFunctionRequest{
		Body: updateFunctionInput,
	}

	_, _err := client.UpdateFunctionWithOptions(tea.String(functionID), updateFunctionRequest, headers, runtime)
	if _err != nil {
		log.Printf("应用更改nas实例失败%v\n", _err)
		return
	}

	// 更改数据库状态（stat，bpm_id俩字段）
	var a = []string{}
	updateStat := fmt.Sprintf("update sd_app set bpm_id='%v', stat='%v', date='%v' where function_short='%v' and stat='0';", bpmId, 1, upDate, fcShort(functionID))
	if _, err := o.Raw(updateStat).QueryRows(&a); err != nil {
		log.Printf("%v更新stat出错%v\n", functionID, err)
		return
	}
}

func Bind(functionID, bpmId, mountDomain string, o orm.Ormer) {
	// 绑定逻辑
	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := FcCreateClient(&id, &secret)

	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)

	updateFunctionInputNASConfigNASMountConfig0 := &fc20230330.NASMountConfig{
		EnableTLS:  tea.Bool(true),
		MountDir:   tea.String("/mnt/auto"),
		ServerAddr: tea.String(mountDomain + ":/fc-stable-diffusion-plus"),
	}
	updateFunctionInputNASConfig := &fc20230330.NASConfig{
		MountPoints: []*fc20230330.NASMountConfig{updateFunctionInputNASConfigNASMountConfig0},
	}

	updateFunctionInputVPCConfig := &fc20230330.VPCConfig{
		SecurityGroupId: tea.String("sg-bp1a8vtwsdkumeb0w01o"),
		VSwitchIds:      []*string{tea.String("vsw-bp18ogozj4w7x0m8xm03x")},
		VpcId:           tea.String("vpc-bp1jdhjct17tfm6bstu8f"),
	}

	updateFunctionInput := &fc20230330.UpdateFunctionInput{
		NasConfig: updateFunctionInputNASConfig,
		VpcConfig: updateFunctionInputVPCConfig,
	}
	updateFunctionRequest := &fc20230330.UpdateFunctionRequest{
		Body: updateFunctionInput,
	}

	_, _err := client.UpdateFunctionWithOptions(tea.String(functionID), updateFunctionRequest, headers, runtime)
	if _err != nil {
		log.Printf("应用更改nas实例失败%v\n", _err)
		return
	}

	// 更改数据库状态（stat，bpm_id俩字段）
	var a = []string{}
	updateStat := fmt.Sprintf("update sd_app set bpm_id='%v', stat='%v' , date='%v' where function_short='%v' and stat='0';", bpmId, 1, upDate, fcShort(functionID))
	if _, err := o.Raw(updateStat).QueryRows(&a); err != nil {
		log.Printf("%v更新stat出错%v\n", functionID, err)
		return
	}
}

func sdBindDefaultNas(functionID string) {

	// 绑定逻辑
	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := FcCreateClient(&id, &secret)

	getFunctionRequest := &fc20230330.GetFunctionRequest{}
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)

	fcResp, _ := client.GetFunctionWithOptions(tea.String(functionID), getFunctionRequest, headers, runtime)

	// 绑定函数和nas
	updateFunctionInputGPUConfig := &fc20230330.GPUConfig{
		GpuMemorySize: tea.Int32(*fcResp.Body.GpuConfig.GpuMemorySize),
		GpuType:       tea.String(*fcResp.Body.GpuConfig.GpuType),
	}

	updateFunctionInputNASConfigNASMountConfig0 := &fc20230330.NASMountConfig{
		EnableTLS:  tea.Bool(true),
		MountDir:   tea.String("/mnt/auto"),
		ServerAddr: tea.String("0650c48ec0-rcs5.cn-hangzhou.nas.aliyuncs.com:/fc-stable-diffusion-plus"),
	}
	updateFunctionInputNASConfig := &fc20230330.NASConfig{
		MountPoints: []*fc20230330.NASMountConfig{updateFunctionInputNASConfigNASMountConfig0},
	}
	updateFunctionInputVPCConfig := &fc20230330.VPCConfig{
		SecurityGroupId: tea.String("sg-bp1a8vtwsdkumeb0w01o"),
		VSwitchIds:      []*string{tea.String("vsw-bp18ogozj4w7x0m8xm03x")},
		VpcId:           tea.String("vpc-bp1jdhjct17tfm6bstu8f"),
	}

	updateFunctionInput := &fc20230330.UpdateFunctionInput{
		GpuConfig: updateFunctionInputGPUConfig,
		NasConfig: updateFunctionInputNASConfig,
		VpcConfig: updateFunctionInputVPCConfig,
	}
	updateFunctionRequest := &fc20230330.UpdateFunctionRequest{
		Body: updateFunctionInput,
	}

	_, _err := client.UpdateFunctionWithOptions(tea.String(functionID), updateFunctionRequest, headers, runtime)
	if _err != nil {
		log.Printf("应用更改默认nas实例失败%v\n", _err)
		return
	}

}

func BindDefaultNas(functionID string) {
	// 绑定逻辑
	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, _ := FcCreateClient(&id, &secret)

	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)

	updateFunctionInputNASConfigNASMountConfig0 := &fc20230330.NASMountConfig{
		EnableTLS:  tea.Bool(true),
		MountDir:   tea.String("/mnt/auto"),
		ServerAddr: tea.String("0650c48ec0-rcs5.cn-hangzhou.nas.aliyuncs.com:/fc-stable-diffusion-plus"),
	}
	updateFunctionInputNASConfig := &fc20230330.NASConfig{
		MountPoints: []*fc20230330.NASMountConfig{updateFunctionInputNASConfigNASMountConfig0},
	}

	updateFunctionInputVPCConfig := &fc20230330.VPCConfig{
		SecurityGroupId: tea.String("sg-bp1a8vtwsdkumeb0w01o"),
		VSwitchIds:      []*string{tea.String("vsw-bp18ogozj4w7x0m8xm03x")},
		VpcId:           tea.String("vpc-bp1jdhjct17tfm6bstu8f"),
	}

	updateFunctionInput := &fc20230330.UpdateFunctionInput{
		NasConfig: updateFunctionInputNASConfig,
		VpcConfig: updateFunctionInputVPCConfig,
	}
	updateFunctionRequest := &fc20230330.UpdateFunctionRequest{
		Body: updateFunctionInput,
	}

	_, _err := client.UpdateFunctionWithOptions(tea.String(functionID), updateFunctionRequest, headers, runtime)
	if _err != nil {
		log.Printf("应用更改默认nas实例失败%v\n", _err)
		return
	}

}

// 解绑
func UnBind(bpmId string) error {
	var o = orm.NewOrm()
	var t = []string{}
	var funcList = []string{}

	queryFunc := fmt.Sprintf("select function_name from sd_app where bpm_id='%v';", bpmId)
	if _, err := o.Raw(queryFunc).QueryRows(&funcList); err != nil {
		log.Printf("查询函数列表失败%v\n", err)
		return err
	}

	for _, funcL := range funcList {
		log.Printf("%v解绑后nas挂载点为: %v\n", funcL, "0650c48ec0-rcs5.cn-hangzhou.nas.aliyuncs.com:/fc-stable-diffusion-plus")
		if strings.Contains(funcL, "__sd") {
			sdBindDefaultNas(funcL)
		} else {
			BindDefaultNas(funcL)
		}
	}

	updateStat := fmt.Sprintf("update sd_app set bpm_id='', stat='0', date='' where bpm_id='%v';", bpmId)
	log.Printf("解绑%v与函数关系\n\n", bpmId)
	if _, err := o.Raw(updateStat).QueryRows(&t); err != nil {
		log.Printf("解绑应用失败%v\n", err)
		return errors.New(fmt.Sprintf("%v解绑应用失败\n", err))
	}

	return nil
}

// 清理 ecs nas db 中的数据
func Clear(bpmId string) (string, error) {
	/*
		1. ecs umount
		2. 删除 nas
		3. 删除sd_designer表中的关系
	*/

	//1. ecs umount
	config := &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password("Q1UiTnP966Wrlk_L6IuysFAGP")},
	}

	addr := fmt.Sprintf("%v:%v", "172.16.197.11", "22")
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("ssh 主机失败:%v\n", err)
		return "", err
	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Printf("创建ssh session 失败", err)
		return "", err
	}
	defer session.Close()

	umountCmd := fmt.Sprintf("umount -f /data/%v", bpmId)
	out, outErr := session.CombinedOutput(umountCmd)
	if outErr != nil {
		log.Printf("umount /data/%v失败: %v\n", bpmId, string(out))
		//return "", errors.New(string(out))
	} else {
		log.Printf("卸载/data/%v目录完成\n", bpmId)
	}

	//	删除 nas
	o := orm.NewOrm()
	var (
		nas      = []string{}
		nasMount = []string{}
	)

	queryNasSql := fmt.Sprintf("select  nas_id,mount_domain from sd_designer where bpm_id='%v';", bpmId)
	if _, err = o.Raw(queryNasSql).QueryRows(&nas, &nasMount); err != nil {
		log.Printf("查询账户关联的nas信息失败%v\n", err)
		return "", err
	}

	if len(nas) == 0 || len(nasMount) == 0 {
		log.Println("未查询到此用户关联nas信息")
		return "", errors.New("db中未查询到此用户关联nas信息")
	}

	account := utilModel.Key["1"]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)
	nasClient, _ := NasCreateClient(&id, &secret)

	deleteMountTargetRequest := &nas20170626.DeleteMountTargetRequest{
		FileSystemId:      tea.String(nas[0]),
		MountTargetDomain: tea.String(nasMount[0]),
	}
	runtime := &util.RuntimeOptions{}

	_, _err := nasClient.DeleteMountTargetWithOptions(deleteMountTargetRequest, runtime)
	if _err != nil {
		log.Printf("删除 nas 挂载点失败%v\n", _err)
		//return "", _err
	}

	time.Sleep(7 * time.Second)
	log.Printf("删除nas挂载点完成")
	deleteFileSystemRequest := &nas20170626.DeleteFileSystemRequest{
		FileSystemId: tea.String(nas[0]),
	}

	_, _err = nasClient.DeleteFileSystemWithOptions(deleteFileSystemRequest, runtime)
	if _err != nil {
		log.Printf("删除nas实例失败%v\n", _err)
		return "", _err
	}

	log.Printf("nas实例删除完成")
	delSql := fmt.Sprintf("delete from sd_designer where bpm_id='%v' and  nas_id='%v' and mount_domain='%v';", bpmId, nas[0], nasMount[0])
	var del = []string{}
	var delBind = []string{}

	if _, err = o.Raw(delSql).QueryRows(&del); err != nil {
		log.Printf("删除sd_designer数据失败%v\n", err)
		return "", err
	}

	delBindSql := fmt.Sprintf("delete  from sd_bind_info where bpm_id='%v' and  nas='%v' ;", nas[0], nasMount[0])

	if _, err = o.Raw(delBindSql).QueryRows(&delBind); err != nil {
		log.Printf("删除sd_bind_info数据失败%v\n", err)
		return "", err
	}
	log.Printf("清理db数据\n")

	return "完成", nil
}

func checkSdStatus(sdDomain, bpmId string) (code int) {
	log.Printf("检查%v绑定的%v域名\n", bpmId, sdDomain)

	resp, err := http.Get("https://" + sdDomain)
	if err != nil {
		log.Printf("检查%v失败: %v\n", "https://"+sdDomain, err)
		return
	}

	defer resp.Body.Close()

	log.Printf("返回状态码:%v\n", resp.StatusCode)
	return resp.StatusCode

}

/*
邀请制，增加一个bpm账号和数据库里是否存在的接口，如果存在可以绑定，不存在不让绑定。绑定接口由管理员操作。
*/
func QueryBpmId(bpmId string) (string, bool, error) {
	var (
		o        = orm.NewOrm()
		name     = []string{}
		querySql = fmt.Sprintf("select name from sd_designer where bpm_id='%v';", bpmId)
	)

	if _, err := o.Raw(querySql).QueryRows(&name); err != nil {
		log.Printf("查询sd_designer %v失败%v\n,", bpmId, err.Error())
		return "", false, err
	}

	if len(name) == 0 {
		log.Printf("%v不存在账号，请联系管理员创建\n", bpmId)
		return fmt.Sprintf("%v不存在账号，请联系管理员创建", bpmId), false, nil
	}

	return fmt.Sprintf("%v已存在账号", name[0]), true, nil
}

func CopyLargeModel(bpmId, model string) error {
	/*
		1. nas 绑定 ecs
		2. 复制到对应目录
	*/

	var (
		o            = orm.NewOrm()
		mountDomains = []string{}
		querySql     = fmt.Sprintf("select mount_domain from sd_designer where bpm_id='%v';", bpmId)
	)
	if _, err := o.Raw(querySql).QueryRows(&mountDomains); err != nil {
		log.Printf("查询挂载点失败%v\n", err)
		return nil
	}

	if len(mountDomains) == 0 {
		log.Printf("未查询到%v nas挂载点\n", bpmId)
		return errors.New(fmt.Sprintf("未查询到%v nas挂载点", bpmId))
	}

	dirBind(bpmId, mountDomains[0], o)

	log.Println("复制模型")
	config := &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password("Q1UiTnP966Wrlk_L6IuysFAGP")},
	}

	addr := fmt.Sprintf("%v:%v", "172.16.197.11", "22")
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("ssh 主机失败:%v\n", err)
		return err
	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Printf("创建ssh session 失败", err)
		return err
	}
	defer session.Close()

	cpCmd := fmt.Sprintf("cp -a -n /data/largeModel/%v /data/%v/fc-stable-diffusion-plus/sd/models/Stable-diffusion/", model, bpmId)
	out, cpErr := session.CombinedOutput(cpCmd)
	if cpErr != nil {
		log.Printf("复制模型数据失败: %v  %v\n", cpErr, string(out))
		return errors.New(string(out))
	}

	return nil
}

func CopyLoraModel(bpmId, modelName string) error {
	/*
		1. nas 绑定 ecs
		2. 复制到对应目录
	*/

	var (
		o            = orm.NewOrm()
		mountDomains = []string{}
		querySql     = fmt.Sprintf("select mount_domain from sd_designer where bpm_id='%v';", bpmId)
	)
	if _, err := o.Raw(querySql).QueryRows(&mountDomains); err != nil {
		log.Printf("查询挂载点失败%v\n", err)
		return nil
	}

	if len(mountDomains) == 0 {
		log.Printf("未查询到%v nas挂载点\n", bpmId)
		return errors.New(fmt.Sprintf("未查询到%v nas挂载点", bpmId))
	}

	dirBind(bpmId, mountDomains[0], o)

	log.Println("复制模型")
	config := &ssh.ClientConfig{
		Timeout:         time.Second,
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password("Q1UiTnP966Wrlk_L6IuysFAGP")},
	}

	addr := fmt.Sprintf("%v:%v", "172.16.197.11", "22")
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("ssh 主机失败:%v\n", err)
		return err
	}

	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		log.Printf("创建ssh session 失败", err)
		return err
	}
	defer session.Close()

	cpCmd := fmt.Sprintf("cp -a -n /data/loraModel/%v /data/%v/fc-stable-diffusion-plus/sd/models/Lora/", modelName, bpmId)
	out, cpErr := session.CombinedOutput(cpCmd)
	if cpErr != nil {
		log.Printf("复制模型数据失败: %v  %v\n", cpErr, string(out))
		return errors.New(string(out))
	}

	return nil
}

func ModelList() interface{} {
	log.Println("查询模型列表")
	data := struct {
		LoraModel  []string
		LargeModel []string
	}{
		LoraModel:  []string{"古风_V1.0.safetensors", "极简风_V2.0.safetensors", "浪漫法式风_V1.0.safetensors", "轻奢风_V1.0.safetensors", "奢华现代新中式风_V1.0.safetensors", "法式奶油风_V1.0.safetensors", "现代风_V1.0.safetensors"},
		LargeModel: []string{"室内设计大师_interiordesignsuperm_v2.safetensors", "麦穗真实人物_majicmixRealistic_v6.safetensors", "现实视觉_realisticVisionV60B1_v51HyperVAE.safetensors", "插画—二次薄.safetensors"},
	}

	return data
}

func TotalMonth() {
	var (
		o     = orm.NewOrm()
		month = []string{}
	)

	if _, err := o.Raw("select count(stat) from sd_app where stat='1';").QueryRows(&month); err != nil {
		log.Printf("月使用人数查询失败%v\n", err)
		return
	}

}
