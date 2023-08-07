package url

import (
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"fmt"
	"regexp"
	"strings"
)

type Rule struct {
	Regex *regexp.Regexp
}

func New() Rule {
	return Rule{
		//匹配url正则
		Regex: regexp.MustCompile(`(https?|ftp|rtsp|mms|wss?|HTTPS?|WSS?|FTP|RTSP)://[\S+(:\S*)?@)]+(/[a-zA-Z0-9\\.]+)?`),
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	for _, match := range r.Regex.FindAllStringSubmatch(string(f.Data.Data), -1) {
		url := strings.ReplaceAll(match[0], "\n", "")
		url = strings.ReplaceAll(url, "\r", "")
		if strings.HasPrefix(url, "http://ocsp.verisign.com") {
			fmt.Println()
		}
		url = strings.Split(url, " ")[0]
		result.Total++
		result.Tags = append(result.Tags, rule.Tag{
			Content: url,
		})
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.URL
}
