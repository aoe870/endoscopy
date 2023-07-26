package template

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
		//正则
		Regex: regexp.MustCompile(``),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {

	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.URL
}
