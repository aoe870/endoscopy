package task

import (
	"bytes"
	"endoscopy/internal/files"
	"endoscopy/internal/report"
	"endoscopy/internal/rule"
	"endoscopy/internal/rule/email"
	"endoscopy/internal/rule/ipv4"
	"endoscopy/internal/rule/ipv6"
	"endoscopy/internal/rule/url"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"net/http"
	u "net/url"
	"path/filepath"
	"strconv"
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
	Info        map[rule.RuleType][]Info `json:"infos"`
	TaskId      string                   `json:"taskId"`
	StartTime   time.Time                `json:"startTile"`
	InPut       string                   `json:"inPath"`
	State       string                   `json:"state"`
	OutPut      string                   `json:"outPut"`
	Msg         string
	Platform    string
}

type Info struct {
	Source  string   `json:"source"`
	Content []string `json:"content"`
}

func (i *Info) To(table []rule.ScanResult) {

	for _, result := range table {
		for _, v := range result.Tags {
			i.Source = result.Source
			i.Content = append(i.Content, v.Content)
		}
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
	taskState.State = "success"
	taskState.InPut = conf.InPutPath
	taskState.OutPut = conf.OutPutPath
	taskState.Platform = conf.Platform
	taskState.Info = make(map[rule.RuleType][]Info)
	defer func() {
		if len(taskState.Platform) <= 0 {
			return
		}
		err := TaskCallBack(taskState)
		if err != nil {
			return
		}
	}()
	metadata, err := files.New(conf.InPutPath)
	if err != nil {
		taskState.State = "fail"
		taskState.Msg = err.Error()
		return
	}

	for _, r := range Rules {
		scanTable := make([]rule.ScanResult, 0)
		for _, f := range metadata.Table {
			if f.Data == nil {
				continue
			}
			scan := r.Check(f)
			if scan.Total <= 0 {
				continue
			}

			scanTable = append(scanTable, scan)
		}

		if len(scanTable) <= 0 {
			continue
		}
		var info Info
		info.To(scanTable)
		taskState.Info[r.GetRuleType()] = append(taskState.Info[r.GetRuleType()], info)
	}

	taskState.ElapsedTime = time.Now().Sub(taskState.StartTime).String()
	//TODO 保存结果
	data, err := json.Marshal(taskState)
	if err != nil {
		return
	}
	output := filepath.Join(conf.OutPutPath, "endoscopy-"+conf.TaskID+".json")
	report.New(output, data)

	//通知平台

	return
}

// 将分析结果发送到平台
func TaskCallBack(state State) error {

	params := u.Values{}
	URL, err := u.Parse(state.Platform)
	if err != nil {
		return err
	}
	URL.RawQuery = params.Encode()
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	type platform struct {
		TaskID     int64  `json:"taskId"`
		ReportPath string `json:"reportPath"`
		Msg        string `json:"msg"`
		Code       int    `json:"code"`
	}
	p := platform{
		ReportPath: state.OutPut,
		Msg:        state.Msg,
	}
	i, err := strconv.ParseInt(state.TaskId, 10, 64)
	if err != nil {
		err = errors.New("TaskID:" + state.TaskId + "TaskID转换失败")
		return err
	}
	p.TaskID = i

	byteData, err := json.Marshal(p)
	if err != nil {
		return err
	}
	response, err := client.Post(URL.String(), "application/json;charset=utf-8", bytes.NewBuffer(byteData))
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	return errors.New("TaskID:" + state.TaskId + " task request fail")
}
