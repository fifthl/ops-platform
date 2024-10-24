package aliModel

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"sync"
	"time"
	utilModel "yw_cloud/models/util"
)

var wg sync.WaitGroup

const (
	Endpoint = "oss-cn-beijing.aliyuncs.com"
)

type bucketProperties struct {
	Name         string    `json:"Name,omitempty"`
	Storage      string    `json:"Storage,omitempty"`
	StorageClass string    `json:"StorageClass,omitempty"`
	Region       string    `json:"Region,omitempty"`
	Domain       string    `json:"Domain"`
	DomainStatus string    `json:"DomainStatus"`
	CreationDate time.Time `json:"CreationDate"`
}

func (p *bucketProperties) changValue() {
	switch p.StorageClass {
	case "Standard":
		p.StorageClass = "标准存储"
	case "Infrequent Access":
		p.StorageClass = "低频访问存储"
	}
}

func GetOSS(AccountID string) (ossList []bucketProperties) {
	account := utilModel.Key[AccountID]
	id, secret := utilModel.Decrypt(account.ID, account.Secret)

	client, err := oss.New(Endpoint, id, secret)
	if err != nil {
		log.Println(err)
		return
	}

	res, err := client.ListBuckets()
	if err != nil {
		log.Println(err)
		return
	}

	var os bucketProperties
	var oss []bucketProperties

	for _, bucket := range res.Buckets {
		os.Name = bucket.Name
		os.StorageClass = bucket.StorageClass
		os.Region = bucket.Region
		os.CreationDate = bucket.CreationDate
		os.Storage = queryStorage(client, bucket.Name)

		domain, status := queryCname(client, bucket.Name)
		os.Domain = domain
		os.DomainStatus = status
		os.changValue()

		oss = append(oss, os)
	}

	return oss
}

// 获取Bucket的总存储量 字节单位
func queryStorage(Client *oss.Client, Bucket string) (storage string) {
	stat, err := Client.GetBucketStat(Bucket)
	if err != nil {
		log.Println(err)
	}

	storage = fmt.Sprintf("%.2f", float32(stat.Storage)/1024/1024/1024)
	return storage + "G"

}

func queryCname(Client *oss.Client, Bucket string) (Domain, Status string) {
	cnResult, err := Client.ListBucketCname(Bucket)
	if err != nil {
		log.Println("Error:", err)
	}

	for _, cnames := range cnResult.Cname {
		if cnames.Domain == "" {
			return "", ""
		}

		Domain, Status = cnames.Domain, cnames.Status

	}

	return Domain, Status
}
