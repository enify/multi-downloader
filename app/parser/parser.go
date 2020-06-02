package parser

import (
	"sort"

	mo "github.com/enify/multi-downloader/app/model"
	"github.com/enify/multi-downloader/app/request"
)

var (
	// RegisteredParsers put parser in here to register it
	RegisteredParsers = []Parser{}
)

type (
	// Parser is an interface of website support parser
	Parser interface {
		GetMeta() Meta

		Prepare(task *mo.Task, client *request.HTTPClient) error
		AtSubTaskDone(subtask *mo.SubTask) error
		AtTaskDone(task *mo.Task) error
	}

	// Meta is information of website support parser
	Meta struct {
		URLRgx   string `json:"urlrgx"`
		Priority int    `json:"priority"`

		Name         string `json:"name"`
		InternalName string `json:"internal_name"`
		Version      string `json:"version"`
		Description  string `json:"description"`
		Author       string `json:"autor"`
		Link         string `json:"link"`
	}
)

func init() {
	register := []Parser{ // 要注册的解析器
		BaseParser{},
		EhentaiParser{},
		NhentaiParser{},
	}

	sort.SliceStable(register, func(i, j int) bool {
		return register[i].GetMeta().Priority < register[j].GetMeta().Priority
	})

	RegisteredParsers = register
}
