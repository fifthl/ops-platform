package aliModel

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

func QueryLora() (interface{}, error) {
	o := orm.NewOrm()

	querySql := "select * from sd_lora;"
	loraList := []struct {
		Dir      string `json:"dir,omitempty"`
		FileName string `json:"fileName,omitempty"`
	}{}

	if _, err := o.Raw(querySql).QueryRows(&loraList); err != nil {
		log.Printf("查询lora列表失败: %v\n", err)
		return nil, err
	}

	return loraList, nil
}

func CopyLora(dir, fileNmae, bpmId string) (string, error) {
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

	log.Printf("复制lora 账号: %v  lora 目录: %v 文件: %v\n\n", bpmId, dir, fileNmae)

	copyLoraCmd := fmt.Sprintf("mkdir -p  /data/%v/fc-stable-diffusion-plus/sd/models/Lora/%v/ &&  cp -a /data/loraModel/%v/%v.*  /data/%v/fc-stable-diffusion-plus/sd/models/Lora/%v/", bpmId, dir, dir, fileNmae, bpmId, dir)
	//log.Println(copyLoraCmd)
	out, smbErr := session.CombinedOutput(copyLoraCmd)
	if smbErr != nil {
		return "", errors.New(string(out))
	}

	return "复制成功", nil
}
