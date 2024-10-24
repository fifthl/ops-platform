package aliModel

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

func init() {
	orm.RegisterModel(&BindInfo{})
}

type BindInfo struct {
	BpmId string `json:"bpmId"`
	Nas   string `json:"nas"`
	Dir   string `json:"dir"`
}

func (*BindInfo) TableName() string {
	return "sd_bind_info"

}
func dirBind(bpmId, nasMount string, o orm.Ormer) {
	log.Printf("开始绑定目录")

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
		return
	}

	defer sshClient.Close()

	createDirCmd := fmt.Sprintf("mkdir -p /data/%v", bpmId)

	//  创建
	err = dirCreate(sshClient, createDirCmd)
	if err != nil {
		log.Printf("创建目录失败%v\n", err)
		return
	}

	// 绑定
	err = dirSmdBind(sshClient, bpmId, nasMount)
	if err != nil {
		log.Printf("smb绑定目录失败%v\n", err)
		time.Sleep(7 * time.Second)
		_ = dirSmdBind(sshClient, bpmId, nasMount)
		return
	}

	// 写入
	mountInfo(bpmId, nasMount, o)

}

func dirCreate(sshClient *ssh.Client, createDirCmd string) error {
	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Printf("创建ssh session 失败", err)
		return err
	}
	defer session.Close()

	out, smbErr := session.CombinedOutput(createDirCmd)
	if smbErr != nil {
		return errors.New(string(out))
	}

	return nil
}

func dirSmdBind(sshClient *ssh.Client, bpmId, nasMount string) error {

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Printf("创建ssh session 失败", err)
		return err
	}
	defer session.Close()

	mountCmd := fmt.Sprintf("mount -t nfs -o vers=4,minorversion=0,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport %v:/ /data/%v", nasMount, bpmId)
	//fmt.Println(mountCmd)

	info, bindErr := session.CombinedOutput(mountCmd)
	if bindErr != nil {
		log.Println(string(info))
		log.Println("err: ", bindErr)
		return errors.New(string(info))
	}

	return nil
}

// 挂载信息写入 db
func mountInfo(bpmId string, nas string, o orm.Ormer) {
	info := BindInfo{
		BpmId: bpmId,
		Nas:   nas,
		Dir:   "/data/" + bpmId,
	}

	if _, err := o.InsertOrUpdate(&info); err != nil {
		log.Printf("写入目录与用户绑定信息失败%v\n", err)
		return
	}

}

func HttpBindDir(bpmId string) (string, error) {
	o := orm.NewOrm()
	query := fmt.Sprintf("select mount_domain from sd_designer where bpm_id='%v';", bpmId)
	var mountDomain = []string{}

	if _, err := o.Raw(query).QueryRows(&mountDomain); err != nil {
		log.Printf("HttpBindDir 查询nas挂载点失败%v\n", err)
		return "", err
	}

	a := len(mountDomain)
	if a == 0 {
		log.Printf("挂载点为空，%v用户需要先创建nas\n", bpmId)
		return "挂载点为空，用户需要先创建nas", nil
	}

	log.Println("开始nas-用户-ecs绑定")
	dirBind(bpmId, mountDomain[0], o)
	return "用户-ecs-nas绑定完成", nil
}
