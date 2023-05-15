package ip

import (
	"landau/internal/files"
	"landau/internal/rule"
)

type Rule struct {
}

func New() Rule {
	return Rule{}
}

func (r Rule) Check(f *files.Node) rule.ScanResult {
	//TODO implement me
	panic("implement me")
}
