package parser

import (
	"sort"

	mo "../model"
)

var (
	// RegisteredParsers put parser in here to register it
	RegisteredParsers = []Parser{}
)

type (
	// Parser is an interface of website support parser
	Parser interface {
		GetMeta() Meta

		Prepare(*mo.Task) error
		AtSubTaskDone()
		AtTaskDone()
	}

	// Meta is information of website support parser
	Meta struct {
		URLRgx   string `json:"urlrgx"`
		Priority int    `json:"priority"`

		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Author      string `json:"autor"`
		Link        string `json:"link"`
	}
)

func init() {
	register := []Parser{ // 要注册的解析器
		BaseParser{},
		EhentaiParser{},
	}

	sort.SliceStable(register, func(i, j int) bool {
		return register[i].GetMeta().Priority < register[j].GetMeta().Priority
	})

	RegisteredParsers = register
}
