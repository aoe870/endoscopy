package email

import (
	"github.com/asaskevich/govalidator"
	"landau/internal/files"
	"landau/internal/rule"
	"regexp"
)

type Rule struct {
	Regex *regexp.Regexp
}

func New() Rule {
	return Rule{
		//匹配邮箱正则
		Regex: regexp.MustCompile(govalidator.Email),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	for _, match := range r.Regex.FindAllStringSubmatch(string(f.Data.Data), -1) {
		result.Total++
		result.Tags = append(result.Tags, rule.Tag{
			Content: match[0],
		})
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.Email
}
