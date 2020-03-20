package parser

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	mo "../model"
	"../request"
	"../util"
)

// BaseParser support base http download
type BaseParser struct{}

// GetMeta return meta of this Parser
func (parser BaseParser) GetMeta() Meta {
	return Meta{
		URLRgx:   `^(http|https)://.*$`,
		Priority: 10,

		Name:         "HTTP下载",
		InternalName: "base-parser",
		Version:      "0.1",
		Description:  "支持基本HTTP下载功能",
		Author:       "",
		Link:         "",
	}
}

// Prepare task Path, Title, FileSize, Preview, Meta, Subtasks,
func (parser BaseParser) Prepare(task *mo.Task, client *request.HTTPClient) (err error) {
	resp, err := client.Req("HEAD", task.URL, nil, "", nil)
	if err != nil {
		return fmt.Errorf("request task url: E:%w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request task url: E:status:%s", resp.Status)
	}

	filename, err := request.FilenameFromResponse(resp)
	if err != nil {
		err = nil
		filename = "index.html"
	}
	if util.IsFileExist(filepath.Join(task.Path, filename)) {
		filename = util.GenUnusedFilename(task.Path, filename)
	}

	task.Title = filename
	task.FileSize = resp.ContentLength
	if exp := regexp.MustCompile(`.+(.jpg|.jpeg|.png|.bmp|.gif)$`); exp.MatchString(filename) {
		task.Preview = filepath.Join(task.Path, filename)
	}

	subtask := &mo.SubTask{
		FileName: task.Title,
		URL:      task.URL,
		Status:   mo.StatusPending,
	}

	task.AddSubTask(subtask)

	return
}

// AtSubTaskDone will call when subtask complete
func (parser BaseParser) AtSubTaskDone(subtask *mo.SubTask) error {
	return nil
}

// AtTaskDone will call when task complete
func (parser BaseParser) AtTaskDone(task *mo.Task) error {
	return nil
}
