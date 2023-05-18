package url

import (
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"github.com/asaskevich/govalidator"
	"regexp"
)

type Rule struct {
	Regex *regexp.Regexp
}

func New() Rule {
	return Rule{
		//匹配url正则
		Regex: regexp.MustCompile(`(http|https|ftp|rtsp|mms|ws|wss)://[a-zA-Z0-9\\.?]+(/[a-zA-Z0-9]+)+`),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	for _, match := range r.Regex.FindAllStringSubmatch(string(f.Data.Data), -1) {
		if govalidator.IsURL(match[0]) {
			result.Total++
			result.Tags = append(result.Tags, rule.Tag{
				Content: match[0],
			})
		}
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.URL
}
