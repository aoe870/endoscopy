package task

import "landau/internal/files"

//引擎返回结果
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

func New(conf Config) {

	files.New(conf.OutPutPath)
}
