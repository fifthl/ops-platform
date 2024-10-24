package sonarqube

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

func init() {
	orm.RegisterModel(new(ResponseCollect))

}

const (
	CollectURL = "http://10.32.12.14:9000/api/measures/component"
	METHOD     = "GET"
	TOKEN      = "squ_baa66c914603090e95b42bde01766d37eb58d588"
)

type param = string

const (
	vulnerabilities      param = "vulnerabilities"        // 漏洞
	bug                  param = "bugs"                   // bug
	securityHotspots     param = "security_hotspots"      // 安全热点
	codeSmells           param = "code_smells"            // 异味
	duplicatedLines      param = "duplicated_lines"       // 重复行
	ncloc                param = "ncloc"                  // 代码行数
	commentLinesDensity  param = "comment_lines_density"  // 注释行-百分比
	cognitiveComplexity  param = "cognitive_complexity"   // 认知复杂度
	newLines             param = "new_lines"              // 新增代码
	securityRating       param = "security_rating"        // 漏洞评级
	reliabilityRating    param = "reliability_rating"     // bug评级
	securityReviewRating param = "security_review_rating" // 安全热点评级
	sqaleRating          param = "sqale_rating"           // 异味评级

)

/*
存储 sonarqube响应的结构体
*/
type RequestCollect struct {
	Component component `json:"component"`
}
type component struct {
	Name     string     `json:"name"`
	Measures []measures `json:"measures"`
}
type measures struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
}

func NewRequest() (req *RequestCollect) {
	req = new(RequestCollect)
	return req
}

// strings 将sonar响应的原始格式化，取出指定值后赋值给新的结构体
func (r *RequestCollect) strings() (vulner, bugs, SecurityHotspots, CodeSmells, DuplicatedLines, Ncloc, CommentLinesDensity,
	CognitiveComplexity, NewLines, SecurityRating, ReliabilityRating, SecurityReviewRating, SqaleRating, Project, AverageRating string) {

	// 查询索引
	vulnerabilitiesIndex := searchIndex(r.Component.Measures, "cognitive_complexity")
	bugIndex := searchIndex(r.Component.Measures, "bugs")
	SecurityHotspotsIndex := searchIndex(r.Component.Measures, "security_hotspots")
	CodeSmellsIndex := searchIndex(r.Component.Measures, "code_smells")
	DuplicatedLinesIndex := searchIndex(r.Component.Measures, "duplicated_lines")
	NclocIndex := searchIndex(r.Component.Measures, "ncloc")
	CommentLinesDensityIndex := searchIndex(r.Component.Measures, "comment_lines_density")
	CognitiveComplexityIndex := searchIndex(r.Component.Measures, "cognitive_complexity")
	NewLinesIndex := searchIndex(r.Component.Measures, "new_lines")
	SecurityRatingIndex := searchIndex(r.Component.Measures, "security_rating")
	ReliabilityRatingIndex := searchIndex(r.Component.Measures, "reliability_rating")
	SecurityReviewRatingIndex := searchIndex(r.Component.Measures, "security_review_rating")
	SqaleRatingIndex := searchIndex(r.Component.Measures, "sqale_rating")

	// 通过索引赋值
	vulner = r.Component.Measures[vulnerabilitiesIndex].Value
	bugs = r.Component.Measures[bugIndex].Value
	SecurityHotspots = r.Component.Measures[SecurityHotspotsIndex].Value
	CodeSmells = r.Component.Measures[CodeSmellsIndex].Value
	DuplicatedLines = r.Component.Measures[DuplicatedLinesIndex].Value
	Ncloc = r.Component.Measures[NclocIndex].Value
	CommentLinesDensity = r.Component.Measures[CommentLinesDensityIndex].Value
	CognitiveComplexity = r.Component.Measures[CognitiveComplexityIndex].Value
	SecurityRating = r.Component.Measures[SecurityRatingIndex].Value
	ReliabilityRating = r.Component.Measures[ReliabilityRatingIndex].Value
	SecurityReviewRating = r.Component.Measures[SecurityReviewRatingIndex].Value
	SqaleRating = r.Component.Measures[SqaleRatingIndex].Value
	Project = r.Component.Name

	// 如果没有新增代码，索引为 -1
	if NewLinesIndex != -1 {
		NewLines = r.Component.Measures[NewLinesIndex].Value
	}

	AverageRating = ConvertString(SecurityRating, ReliabilityRating, SecurityReviewRating, SqaleRating)
	return

}

func ConvertString(SecurityRating, ReliabilityRating, SecurityReviewRating, SqaleRating string) (i string) {
	float64s, _ := strconv.ParseFloat(SecurityRating, 1)
	float64r, _ := strconv.ParseFloat(ReliabilityRating, 1)
	float64sr, _ := strconv.ParseFloat(SecurityReviewRating, 1)
	float64sr1, _ := strconv.ParseFloat(SqaleRating, 1)

	i = fmt.Sprintf("%.2f", 6-(float64s+float64r+float64sr+float64sr1)/4)
	return i
}

// 查询索引
func searchIndex(s []measures, Norm string) (index int) {
	i := slices.IndexFunc(s, func(m measures) bool {
		return m.Metric == Norm
	})

	return i
}

// ResponseCollect 响应给前端或者写入db的结构体
type ResponseCollect struct {
	Name                 string    `json:"name,omitempty" orm:"column(name)"`
	Vulnerabilities      string    `json:"vulnerabilities,omitempty" orm:"column(vulnerabilities)"`
	Bug                  string    `json:"bug,omitempty" orm:"column(bug)"`
	SecurityHotspots     string    `json:"security_hotspots,omitempty" orm:"column(security_hotspots)"`
	CodeSmells           string    `json:"code_smells,omitempty" orm:"column(code_smells)"`
	DuplicatedLines      string    `json:"duplicated_lines,omitempty" orm:"column(duplicated_lines)"`
	Ncloc                string    `json:"ncloc,omitempty" orm:"column(ncloc)"`
	CommentLinesDensity  string    `json:"comment_lines_density,omitempty" orm:"column(comment_lines_density)"`
	CognitiveComplexity  string    `json:"cognitive_complexity,omitempty" orm:"column(cognitive_complexity)"`
	NewLines             string    `json:"new_lines,omitempty" orm:"column(new_lines)"`
	SecurityRating       string    `json:"security_rating,omitempty" orm:"column(security_rating)"`
	ReliabilityRating    string    `json:"reliability_rating,omitempty" orm:"column(reliability_rating)"`
	SecurityReviewRating string    `json:"security_review_rating,omitempty" orm:"column(security_review_rating)"`
	SqaleRating          string    `json:"sqale_rating,omitempty" orm:"column(sqale_rating)"`
	AverageRating        string    `json:"average_rating" orm:"column(average_rating)"`
	CollectTime          time.Time `json:"collect_time" orm:"auto_now_add;type(datetime);column(collect_time)"`
}

// 赋值
func (r *ResponseCollect) responseInit(Vulnerabilities, Bug, SecurityHotspots, CodeSmells, DuplicatedLines, Ncloc, CommentLinesDensity, CognitiveComplexity,
	NewLines, SecurityRating, ReliabilityRating, SecurityReviewRating, SqaleRating, Name, AverageRating string) {

	r.Vulnerabilities = Vulnerabilities
	r.Bug = Bug
	r.SecurityHotspots = SecurityHotspots
	r.CodeSmells = CodeSmells
	r.DuplicatedLines = DuplicatedLines
	r.Ncloc = Ncloc
	r.CommentLinesDensity = CommentLinesDensity
	r.CognitiveComplexity = CognitiveComplexity
	r.NewLines = NewLines
	r.SecurityRating = SecurityRating
	r.ReliabilityRating = ReliabilityRating
	r.SecurityReviewRating = SecurityReviewRating
	r.SqaleRating = SqaleRating
	r.Name = Name
	r.AverageRating = AverageRating

	// 日期没有通过参数传递

}
func (*ResponseCollect) TableName() string {
	return "sonar_collect"
}

func NewResponse() (res *ResponseCollect) {
	res = new(ResponseCollect)
	return res
}

// getCollect 发起请求
func getCollect(Project string) {
	apiURL := fmt.Sprintf("%s?component=%s&metricKeys=%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
		CollectURL, Project, vulnerabilities, bug, securityHotspots, codeSmells, duplicatedLines, ncloc, commentLinesDensity, cognitiveComplexity, newLines,
		securityRating, reliabilityRating, securityReviewRating, sqaleRating)

	client := http.Client{}
	payload := strings.NewReader(``)

	req, err := http.NewRequest(METHOD, apiURL, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(TOKEN, "")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	r := NewRequest()
	if err = json.Unmarshal(body, r); err != nil {
		fmt.Println("序列化失败:", err)
	}

	// 从响应中返回指定字段
	vulner, bugs, SecurityHotspots, CodeSmells, DuplicatedLines, Ncloc, CommentLinesDensity,
		CognitiveComplexity, NewLines, SecurityRating, ReliabilityRating, SecurityReviewRating, SqaleRating, Project, AverageRating := r.strings()

	// 将参数传入响应结构体后返回
	resp := NewResponse()
	resp.responseInit(vulner, bugs, SecurityHotspots, CodeSmells, DuplicatedLines, Ncloc, CommentLinesDensity,
		CognitiveComplexity, NewLines, SecurityRating, ReliabilityRating, SecurityReviewRating, SqaleRating, Project, AverageRating)

	o := orm.NewOrm()
	if _, err = o.Insert(resp); err != nil {
		log.Println("sonar 扫描结果写入db失败")
	}
}

/*
更新sonarqube扫描的指标
*/
func SonarCollect() {
	ch := make(chan string, 10)
	exit := make(chan bool)

	defer close(ch)

	//获取有哪些项目
	go getProject(ch)

	go func() {
		for {
			select {
			case project := <-ch:
				go getCollect(project)
			case <-time.After(1 * time.Second):
				exit <- true
				return
			}
		}
	}()

	<-exit
}
