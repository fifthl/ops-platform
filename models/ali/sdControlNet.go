package aliModel

// 接口弃用
import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

func QueryControlNet() interface{} {

	// ControlNet 列表（不含后缀）
	controlNetList := []string{
		"control_v11e_sd15_ip2p",
		"control_v11e_sd15_shuffle",
		"control_v11f1e_sd15_tile",
		"control_v11f1p_sd15_depth",
		"control_v11p_sd15_canny",
		"control_v11p_sd15_inpaint",
		"control_v11p_sd15_lineart",
		"control_v11p_sd15_mlsd",
		"control_v11p_sd15_normalbae",
		"control_v11p_sd15_openpose",
		"control_v11p_sd15s2_lineart_anime",
		"control_v11p_sd15_scribble",
		"control_v11p_sd15_seg",
		"control_v11p_sd15_softedge",
		"control_v1p_sd15_brightness",
	}

	data := struct {
		controlNetList []string
	}{
		controlNetList: controlNetList,
	}
	
	return data

}

func CopyControlNet(dir, fileNmae, bpmId string) (string, error) {
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

	copyLoraCmd := fmt.Sprintf("mkdir -p  /data/%v/fc-stable-diffusion-plus/sd/models/Lora/%v/ &&  cp -a /data/loraModel/%v/%v.*  /data/%v/fc-stable-diffusion-plus/sd/models/Lora/%v/", bpmId, dir, dir, fileNmae, bpmId, dir)
	out, smbErr := session.CombinedOutput(copyLoraCmd)
	if smbErr != nil {
		return "", errors.New(string(out))
	}

	return "复制成功", nil
}
