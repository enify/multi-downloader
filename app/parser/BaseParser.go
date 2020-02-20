package parser

import (
	mo "../model"
)

// BaseParser support base http download
type BaseParser struct{}

// GetMeta return meta of this Parser
func (parser BaseParser) GetMeta() Meta {
	return Meta{
		URLRgx:   `^(http|https)://.*$`,
		Priority: 10,

		Name:        "HTTP下载",
		Version:     "0.1",
		Description: "支持基本HTTP下载功能",
		Author:      "",
		Link:        "",
	}
}

// Prepare task meta and subtask with task url
func (parser BaseParser) Prepare(task *mo.Task) (err error) { return }

// AtSubTaskDone will call when subtask complete
func (parser BaseParser) AtSubTaskDone() { return }

// AtTaskDone will call when task complete
func (parser BaseParser) AtTaskDone() { return }
