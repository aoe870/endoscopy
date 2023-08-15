package email

import (
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"github.com/asaskevich/govalidator"
	"regexp"
	"strings"
)

type Rule struct {
	Regex *regexp.Regexp
}

func New() Rule {
	return Rule{
		//匹配邮箱正则
		Regex: regexp.MustCompile(`[\w\\._-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+`),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	for _, match := range r.Regex.FindAllStringSubmatch(string(f.Data.Data), -1) {
		if govalidator.IsEmail(match[0]) {
			suffix := strings.Split(match[0], "@")[1]
			if !strings.Contains(suffix, ".") {
				continue
			}
			email := strings.ReplaceAll(match[0], "\n", "")
			email = strings.ReplaceAll(match[0], "\r", "")
			result.Total++
			result.Tags = append(result.Tags, rule.Tag{
				Content: email,
			})
		}
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.Email
}
