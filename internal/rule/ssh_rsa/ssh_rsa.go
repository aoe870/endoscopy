package ssh_rsa

import (
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"strings"
)

type Rule struct {
	Keyword     []string
	PathKeyword []string
}

func New() Rule {
	return Rule{
		//ssh_rsa私钥 关键字
		Keyword:     []string{"-----BEGIN OPENSSH PRIVATE KEY-----"},
		PathKeyword: []string{"id_"},
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	for _, path := range r.PathKeyword {
		if !strings.HasPrefix(f.Name, path) {
			continue
		}

		for _, keyword := range r.Keyword {
			if strings.Contains(string(f.Data.Data), keyword) {
				result.Total++
				result.Tags = append(result.Tags, rule.Tag{
					Content: string(f.Data.Data),
				})
			}
		}
	}

	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.RSA
}
