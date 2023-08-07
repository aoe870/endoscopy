package ipv4

import (
	"bufio"
	"bytes"
	"endoscopy/internal/files"
	"endoscopy/internal/rule"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Rule struct {
	Regex     *regexp.Regexp
	Exclusive []*regexp.Regexp
}

func New() Rule {
	return Rule{
		//ipv4正则
		Regex: regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`),
		Exclusive: []*regexp.Regexp{
			regexp.MustCompile(`<version>[a-zA-Z\d]+\.[a-zA-Z\d]+\.[a-zA-Z\d]+\.[a-zA-Z\d]+</version>`),
			regexp.MustCompile(`<project.version>[a-zA-Z\d]+\.[a-zA-Z\d]+\.[a-zA-Z\d]+\.[a-zA-Z\d]+</project.version>`),
			//匹配http://
			regexp.MustCompile(`http://\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`),
		},
	}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	var result rule.ScanResult
	result.Source = f.Path
	result.RuleType = r.GetRuleType()

	//检查是否符合规则
	reader := bytes.NewReader(f.Data.Data)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		for _, match := range r.Regex.FindAllStringSubmatch(scanner.Text(), -1) {
			if govalidator.IsIPv4(match[0]) {
				ipv4 := strings.ReplaceAll(match[0], "\n", "")
				ipv4 = strings.ReplaceAll(ipv4, "\r", "")
				//根据Keywords过滤
				if r.filtration(scanner.Text()) {
					continue
				}
				if len(ipv4) > 18 {
					continue
				}
				result.Total++
				result.Tags = append(result.Tags, rule.Tag{
					Content: ipv4,
				})
			}
		}
	}
	return result
}

func (r Rule) GetRuleType() rule.RuleType {
	return rule.Ipv4
}

func (r Rule) filtration(ipv4 string) bool {
	for _, exclusive := range r.Exclusive {
		if exclusive.FindString(ipv4) != "" {
			return true
		}
	}
	return false
}
