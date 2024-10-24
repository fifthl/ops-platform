package sonarqube

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	ProjectURL = "http://10.32.12.14:9000/api/projects/search?ps=200"
)

type responseProject struct {
	Components []components `json:"components"`
}

type components struct {
	Name string `json:"name,omitempty"`
}

func getProject(ch chan string) {
	//func getProject() (projects []string) {

	client := http.Client{}
	payload := strings.NewReader(``)

	req, err := http.NewRequest("GET", ProjectURL, payload)
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

	rs := new(responseProject)
	if err = json.Unmarshal(body, rs); err != nil {
		log.Println("序列化失败", err)
	}

	for _, c := range rs.Components {
		//projects = append(projects, c.Name)
		ch <- c.Name
	}

	//return projects
}
