package task

import (
	"landau/internal/files"
	"landau/internal/rule"
	"landau/internal/rule/email"
	"landau/internal/rule/ip"
	"landau/internal/rule/phone"
	"landau/internal/rule/url"
	"time"
)

// 引擎返回结果
type Scan_State int

const (
	Scan_error   Scan_State = 4 // 扫描异常
	Scan_success Scan_State = 5 // 扫描完成
)

type Config struct {
	TaskID     string
	OutPutPath string
	InPutPath  string
}

type TaskState struct {
	ElapsedTime string        `json:"elapsedTime"`
	Detection   []Detection   `json:"detection"`
	TaskId      string        `json:"taskId"`
	StartTile   time.Duration `json:"startTile"`
}

type Detection struct {
	Source  string `json:"source"`
	Content string `json:"content"`
	Line    int    `json:"line"`
}

func (detection Detection) To(result rule.ScanResult) {
	result.Source = detection.Source
	result.Content = detection.Content
	result.Line = detection.Line
}

func New(conf Config) {

	Rules := []rule.Rule{
		email.New(),
		phone.New(),
		ip.New(),
		url.New(),
	}

	var taskState TaskState
	metadata, err := files.New(conf.InPutPath)
	if err != nil {
		return
	}

	for _, v := range metadata.Table {
		for _, r := range Rules {
			scan := r.Check(v)
			if scan.State == rule.ScanSeccess {
				//TODO 保存结果
				var detection Detection
				detection.To(scan)
				taskState.Detection = append(taskState.Detection, detection)
			}
		}
	}

}
