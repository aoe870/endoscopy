package rule

import (
	"endoscopy/internal/files"
)

type RuleType string

const (
	Email RuleType = "email"
	Ipv4  RuleType = "ipv4"
	Ipv6  RuleType = "ipv6"
	URL   RuleType = "url"
)

type Rule interface {
	Check(f *files.Node) ScanResult

	GetRuleType() RuleType
}

type ScanResult struct {
	Total    int    //是否检出
	Source   string //检出来源
	Tags     []Tag
	RuleType RuleType //规则类型

}

type Tag struct {
	Content string //检出内容
	Line    int    //检出行
}
