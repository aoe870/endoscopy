package url

import (
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"regexp"

	"github.com/asaskevich/govalidator"
)

type Rule struct {
	Regex *regexp.Regexp
}

func New() Rule {
	return Rule{
		//匹配url正则
		Regex: regexp.MustCompile(govalidator.URL),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//rsa私钥正则
	regexp.MustCompile(``)

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
	return rule.URL
}
