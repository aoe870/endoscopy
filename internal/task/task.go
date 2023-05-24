package task

import (
	"endoscopy/internal/files"
	"endoscopy/internal/logs"
	"endoscopy/internal/report"
	"endoscopy/internal/rule"
	"endoscopy/internal/rule/email"
	"endoscopy/internal/rule/ipv4"
	"endoscopy/internal/rule/ipv6"
	"endoscopy/internal/rule/ssh_rsa"
	"endoscopy/internal/rule/url"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
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
	Platform   string
}

type State struct {
	ElapsedTime string                   `json:"elapsedTime"`
	StartTime   time.Time                `json:"startTile"`
	InPut       string                   `json:"inPath"`
	State       string                   `json:"state"`
	OutPut      string                   `json:"outPut"`
	Info        map[rule.RuleType][]Info `json:"infos"`
	TaskId      string
}

type Info struct {
	Source  string   `json:"source"`
	Content []string `json:"content"`
}

func (i *Info) To(table rule.ScanResult) {
	for _, v := range table.Tags {
		i.Source = table.Source
		i.Content = append(i.Content, v.Content)
	}
}

func New(conf Config) {

	logs.Info("Start Task: " + conf.TaskID)
	Rules := []rule.Rule{
		email.New(),
		ipv4.New(),
		ipv6.New(),
		url.New(),
		ssh_rsa.New(),
	}
	var taskState State
	taskState.StartTime = time.Now()
	taskState.TaskId = conf.TaskID
	taskState.State = "success"
	taskState.InPut = conf.InPutPath
	taskState.OutPut = conf.OutPutPath
	taskState.Info = make(map[rule.RuleType][]Info)
	metadata, err := files.New(conf.InPutPath)
	if err != nil {
		taskState.State = "fail"
		logs.Error(errors.Wrap(err, "task fail:"+err.Error()))
		return
	}

	for _, r := range Rules {
		for _, f := range metadata.Table {
			if f.Data == nil {
				continue
			}
			scan := r.Check(f)
			if scan.Total <= 0 {
				continue
			}
			var info Info
			info.To(scan)
			taskState.Info[r.GetRuleType()] = append(taskState.Info[r.GetRuleType()], info)
		}
	}
	taskState.ElapsedTime = time.Now().Sub(taskState.StartTime).String()
	//TODO 保存结果
	data, err := json.Marshal(taskState)
	if err != nil {
		return
	}

	logs.Info(" scan end")
	if conf.OutPutPath == "" {
		fmt.Println(string(data))
	} else {
		err := report.New(conf.OutPutPath, data)
		if err != nil {
			logs.Error(errors.Wrap(err, "save report fail"))
			return
		}
		logs.Info("save report :" + conf.OutPutPath)
	}

	return
}
