package ipv4

import (
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Rule struct {
	Regex *regexp.Regexp
}

func New() Rule {
	return Rule{
		//ipv4正则
		Regex: regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	for _, match := range r.Regex.FindAllStringSubmatch(string(f.Data.Data), -1) {
		if govalidator.IsIPv4(match[0]) {
			ipv4 := strings.ReplaceAll(match[0], "\n", "")
			ipv4 = strings.ReplaceAll(ipv4, "\r", "")
			if len(ipv4) > 18 {
				continue
			}
			result.Total++
			result.Tags = append(result.Tags, rule.Tag{
				Content: ipv4,
			})
		}
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.Ipv4
}
