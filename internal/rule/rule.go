package rule

import "landau/internal/files"

type ScanState int

const (
	ScanSeccess ScanState = 1
	ScanFail    ScanState = 2
)

type Rule interface {
	// 检查是否符合规则
	Check(f *files.Node) ScanResult
}

type ScanResult struct {
	State   ScanState //是否检出
	Source  string    //检出来源
	Content string    //检出内容
	Line    int       //检出行
}
