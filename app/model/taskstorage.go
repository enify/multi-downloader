package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type (
	// TaskStorage 任务存储
	TaskStorage struct {
		path string
		mux  sync.Mutex

		Tasks []*Task `json:"tasks"`
	}

	// Task 任务
	Task struct {
		mux sync.Mutex

		URL    string            `json:"url"`
		Status TaskStatus        `json:"status"`
		Path   string            `json:"path"`
		Meta   map[string]string `json:"meta"`

		SubTasks []*SubTask `json:"subtasks"`
	}

	// SubTask 任务的子任务
	SubTask struct {
		URL      string     `json:"url"`
		Status   TaskStatus `json:"status"`
		FileName string     `json:"filename"`
		ErrMsg   string     `json:"errmsg"`
	}
)

// NewTaskStorage 返回 TaskStorage 指针对象
func NewTaskStorage(path string) (ts *TaskStorage, err error) {
	ts = &TaskStorage{path: path, Tasks: []*Task{}}

	if !ts.FileExist() {
		err = ts.Save()
	} else {
		data, err := ioutil.ReadFile(ts.path)
		if err != nil {
			return ts, err
		}

		err = json.Unmarshal(data, ts)
	}

	return
}

// FileExist 检查任务存储文件是否存在
func (ts *TaskStorage) FileExist() bool {
	filePath, _ := filepath.Abs(ts.path)
	_, err := os.Stat(filePath)

	return os.IsExist(err)
}

// Save 保存任务到任务存储文件
func (ts *TaskStorage) Save() (err error) {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	data, err := json.MarshalIndent(ts, "", "    ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(ts.path, data, 0666)
	if err != nil {
		return
	}

	return
}
