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

		ID         string            `json:"id"`
		URL        string            `json:"url"`
		Status     TaskStatus        `json:"status"`
		Path       string            `json:"path"`
		Title      string            `json:"title"`
		FileSize   int64             `json:"file_size"`
		Preview    string            `json:"preview"` // 预览图的路径
		Meta       map[string]string `json:"meta"`
		ParserName string            `json:"parser_name"` // Parser内部名
		Err        error             `json:"error"`

		SubTasks []*SubTask `json:"subtasks"`
	}

	// SubTask 任务的子任务
	SubTask struct {
		FileName string     `json:"filename"`
		URL      string     `json:"url"`
		Status   TaskStatus `json:"status"`
		Err      error      `json:"error"`
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

	return !os.IsNotExist(err)
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

// Find find task by id
func (ts *TaskStorage) Find(id string) (t *Task) {
	for _, task := range ts.Tasks {
		if id == task.ID {
			return task
		}
	}
	return nil
}

// AddTask add task to storage
func (ts *TaskStorage) AddTask(t *Task) {
	ts.mux.Lock()
	ts.Tasks = append(ts.Tasks, t)
	ts.mux.Unlock()

	ts.Save()
}

// DeleteTask remove task from storage
func (ts *TaskStorage) DeleteTask(t *Task) (deleted bool) {
	for idx, task := range ts.Tasks {
		if t.ID == task.ID {
			ts.mux.Lock()
			ts.Tasks = append(ts.Tasks[:idx], ts.Tasks[idx+1:]...)
			ts.mux.Unlock()
			deleted = true
			break
		}
	}
	ts.Save()
	return
}

// HasTask check task exist or not
func (ts *TaskStorage) HasTask(url string) bool {
	for _, task := range ts.Tasks {
		if url == task.URL {
			return true
		}
	}

	return false
}

// AddSubTask add subtask to task
func (t *Task) AddSubTask(st *SubTask) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.SubTasks = append(t.SubTasks, st)
}
