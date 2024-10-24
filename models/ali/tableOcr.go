package aliModel

import (
	"encoding/json"
	ocr20210707 "github.com/alibabacloud-go/ocr-api-20210707/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"net/http"
	"time"
	utilModel "yw_cloud/models/util"
)

const (
	ywEndpoint = "https://oss-cn-beijing.aliyuncs.com"
	ExcelDir   = "excel/"
)

type ResponseData struct {
	SubImages []SubImages `json:"SubImages"`
	ExcelUrl  string
}
type BlockDetails struct {
	BlockAngle      int    `json:"BlockAngle"`
	BlockConfidence int    `json:"BlockConfidence"`
	BlockContent    string `json:"BlockContent"`
	BlockID         int    `json:"BlockId"`
}
type BlockInfo struct {
	BlockCount   int            `json:"BlockCount"`
	BlockDetails []BlockDetails `json:"BlockDetails"`
}
type CellDetails struct {
	BlockList   []int  `json:"BlockList,omitempty"`
	CellContent string `json:"CellContent"`
	CellID      int    `json:"CellId"`
	ColumnEnd   int    `json:"ColumnEnd"`
	ColumnStart int    `json:"ColumnStart"`
	RowEnd      int    `json:"RowEnd"`
	RowStart    int    `json:"RowStart"`
}
type Footer struct {
}
type Header struct {
}
type TableDetails struct {
	CellCount   int           `json:"CellCount"`
	CellDetails []CellDetails `json:"CellDetails"`
	ColumnCount int           `json:"ColumnCount"`
	Footer      Footer        `json:"Footer"`
	Header      Header        `json:"Header"`
	RowCount    int           `json:"RowCount"`
	TableID     int           `json:"TableId"`
}
type TableInfo struct {
	TableCount   int            `json:"TableCount"`
	TableDetails []TableDetails `json:"TableDetails"`
	TableExcel   string         `json:"TableExcel"`
}
type SubImages struct {
	TableInfo TableInfo `json:"TableInfo"`
}

func TableOcr(object string) (data *ResponseData, err error) {
	account := utilModel.Key["1"]
	id, secret := IdSecret(account.ID, account.Secret)
	t := time.Now().Format("2006-01-02-15:04:04")

	excelPath := ExcelDir + t + ".xlsx"

	ossClient, err := oss.New(ywEndpoint, id, secret)
	if err != nil {
		log.Printf("创建oss client失败 %v", err)
		return nil, err
	}

	bucket, err := ossClient.Bucket("yw-private")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	signUrl, err := bucket.SignURL(object, oss.HTTPGet, 120)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ocrClient, err := OcrCreateClient(&id, &secret)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tableConfig := &ocr20210707.RecognizeAllTextRequestTableConfig{
		OutputTableExcel:   tea.Bool(true),
		IsHandWritingTable: tea.Bool(false),
		IsLineLessTable:    tea.Bool(false),
	}
	recognizeAllTextRequest := &ocr20210707.RecognizeAllTextRequest{
		Url:         tea.String(signUrl),
		Type:        tea.String("Table"),
		TableConfig: tableConfig,
	}
	runtime := &util.RuntimeOptions{}

	result, err := ocrClient.RecognizeAllTextWithOptions(recognizeAllTextRequest, runtime)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	d := new(ResponseData)
	if err = json.Unmarshal([]byte(result.Body.Data.String()), d); err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.Get(d.SubImages[0].TableInfo.TableExcel)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = bucket.PutObject(excelPath, io.Reader(res.Body)); err != nil {
		log.Printf("put excel err %v", err)
		return nil, err
	}

	sig, err := bucket.SignURL(excelPath, oss.HTTPGet, 60)
	if err != nil {
		log.Printf("get sigurl err %v", err)
		return nil, err
	}

	d.ExcelUrl = sig
	return d, nil
}
