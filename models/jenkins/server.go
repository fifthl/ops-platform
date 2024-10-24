package ci

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func init() {
	orm.RegisterModel(new(ReleaseInfo))
}

type ReleaseInfo struct {
	Project    string `json:"Project,omitempty"`
	Branch     string `json:"Branch,omitempty"`
	Name       string `json:"Name,omitempty"`
	Email      string `json:"Email,omitempty"`
	ChangeInfo string `json:"ChangeInfo,omitempty"`
	Rows       string `json:"Rows,omitempty"`
	BuildNum   string `json:"BuildNum"`
	GitUrl     string `json:"GitUrl"`
	Commit     string `json:"Commit"`
	//20240228添加
	Adds        string `json:"Adds"`
	Dels        string `json:"Dels"`
	JenkinsDate string `json:"JenkinsDate"`
	Date        string `json:"Date,omitempty" orm:"auto_now_add;type(datetime)"`
}

// 新增的记录会更新名称
func (receiver *ReleaseInfo) changName(o orm.Ormer) {
	query := fmt.Sprintf("select displayname from jenkins_name where gitname='%v';", receiver.Name)
	var name = []string{}

	if _, err := o.Raw(query).QueryRows(&name); err != nil {
		log.Println("名称查询出错: ", err)
	}

	if len(name) == 0 {
		log.Println(fmt.Sprintf("jenkins_name中未查询到%v用户", receiver.Name))
		return
	}

	receiver.Name = name[0]
}

func (r *ReleaseInfo) TableName() string {
	return "jenkins"
}

func NewSaveDb(r *ReleaseInfo) error {
	o := orm.NewOrm()

	preEnv := strings.Contains(r.Project, "pre-")
	testEnv := strings.Contains(r.Project, "test-")
	prodEnv := strings.Contains(r.Project, "prod-")
	devEnv := strings.Contains(r.Project, "dev-")
	releaseEnv := strings.Contains(r.Project, "release-")
	materEnv := strings.Contains(r.Project, "master-")

	preEnv1 := strings.Contains(r.Project, "PRE-")
	testEnv1 := strings.Contains(r.Project, "TEST-")
	prodEnv1 := strings.Contains(r.Project, "PROD-")
	devEnv1 := strings.Contains(r.Project, "DEV-")

	if devEnv || testEnv || preEnv || prodEnv || releaseEnv || materEnv {
		a := strings.Split(r.Project, r.Branch+"-")
		r.Project = a[1]
	}

	if preEnv1 || testEnv1 || prodEnv1 || devEnv1 {
		//fmt.Println("NewSaveDb: ", r)
		branch := ""

		if r.Branch == "master" {
			branch = "PROD"
		}

		a := strings.Split(r.Project, branch+"-")
		r.Project = a[1]
	}

	r.changName(o)

	//fmt.Println("写入db记录: ", r)

	if _, err := o.Insert(r); err != nil {
		log.Println("发版信息记录失败: ", err)
		return err
	}

	return nil
}
