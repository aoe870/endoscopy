package api_key

import (
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"regexp"
)

type Rule struct {
	Regex *regexp.Regexp
}

func New() Rule {
	return Rule{
		//正则匹配
		Regex: regexp.MustCompile(`API_KEY\s*=\s*["']?([A-Za-z0-9]{32})["']?`),
		//写一个正则

	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {

	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	for _, match := range r.Regex.FindAllStringSubmatch(string(f.Data.Data), -1) {
		apiKey := match[0]
		result.Total++
		result.Tags = append(result.Tags, rule.Tag{
			Content: apiKey,
		})
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.API_KEY
}
