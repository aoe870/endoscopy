package task

import (
	"github.com/goccy/go-json"
	"landau/internal/files"
	"landau/internal/report"
	"landau/internal/rule"
	"landau/internal/rule/email"
	"landau/internal/rule/ipv4"
	"landau/internal/rule/ipv6"
	"landau/internal/rule/url"
	"path/filepath"
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

type State struct {
	ElapsedTime string      `json:"elapsedTime"`
	Info        []Detection `json:"info"`
	TaskId      string      `json:"taskId"`
	StartTime   time.Time   `json:"startTile"`
	Path        string      `json:"path"`
}

type Detection struct {
	Source  string   `json:"source"`
	Type    string   `json:"type"`
	Content []string `json:"content"`
}

func (detection *Detection) To(result rule.ScanResult) {
	detection.Source = result.Source
	detection.Type = string(result.RuleType)
	for _, v := range result.Tags {
		detection.Content = append(detection.Content, v.Content)
	}
}

func New(conf Config) {

	Rules := []rule.Rule{
		email.New(),
		ipv4.New(),
		ipv6.New(),
		url.New(),
	}
	var taskState State
	taskState.StartTime = time.Now()
	taskState.TaskId = conf.TaskID
	metadata, err := files.New(conf.InPutPath)
	if err != nil {
		return
	}

	for _, f := range metadata.Table {
		if f.Data == nil {
			continue
		}
		for _, r := range Rules {
			scan := r.Check(f)
			if scan.Total > 0 {
				var detection Detection
				detection.To(scan)
				taskState.Info = append(taskState.Info, detection)
			}
		}
	}

	taskState.ElapsedTime = time.Now().Sub(taskState.StartTime).String()
	//TODO 保存结果
	data, err := json.Marshal(taskState)
	if err != nil {
		return
	}
	output := filepath.Join(conf.OutPutPath, "landau-"+conf.TaskID+".json")
	report.New(output, data)

	//通知平台
	return
}
