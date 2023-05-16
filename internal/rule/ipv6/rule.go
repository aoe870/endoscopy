package ipv6

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
		Regex: regexp.MustCompile(govalidator.IP),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	for _, match := range r.Regex.FindAllStringSubmatch(string(f.Data.Data), -1) {
		if govalidator.IsIPv6(match[0]) {
			result.Total++
			result.Tags = append(result.Tags, rule.Tag{
				Content: match[0],
			})
		}
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.Ipv6
}
