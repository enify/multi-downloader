package worker

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	mo "github.com/enify/multi-downloader/app/model"
	"github.com/enify/multi-downloader/app/request"
)

// TaskDownloadWork subtask download Work type
type TaskDownloadWork struct {
	Task    *mo.Task
	SubTask *mo.SubTask
	Client  *request.HTTPClient

	AtTaskDone    func(task *mo.Task) error
	AtSubTaskDone func(subtask *mo.SubTask) error
}

// Do execute single subtask download job
func (w *TaskDownloadWork) Do() (err error) {
	if w.Task.Status != mo.StatusRunning {
		return
	}

	if w.SubTask.Status != mo.StatusPending {
		return
	}

	w.SubTask.Status = mo.StatusRunning
	err = func() error {
		err := os.MkdirAll(w.Task.Path, 0755)
		if err != nil {
			return fmt.Errorf("createTaskPath: path:%s, E:%s", w.Task.Path, err)
		}

		w.Client.SetTimeout(10 * time.Minute)

		resp, err := w.Client.Req("GET", w.SubTask.URL, nil, "", nil)
		if err != nil {
			return fmt.Errorf("requestTaskUrl: url:%s, E:%s", w.SubTask.URL, err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("requestTaskUrl: url:%s, E:status:%s", w.SubTask.URL, resp.Status)
		}

		filePath := filepath.Join(w.Task.Path, w.SubTask.FileName)
		f, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("createTaskFile: path:%s, E:%s", filePath, err)
		}

		defer f.Close()
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return fmt.Errorf("writeTaskFile: path:%s, E:%s", filePath, err)
		}

		return nil
	}()

	if err != nil {
		w.SubTask.Status = mo.StatusError
	} else {
		w.SubTask.Status = mo.StatusDone
	}
	w.SubTask.Err = err

	var count = map[mo.TaskStatus]int{}
	for _, t := range w.Task.SubTasks {
		count[t.Status]++
	}
	if count[mo.StatusDone] == len(w.Task.SubTasks) {
		w.Task.FinishAt = time.Now()
		w.Task.Status = mo.StatusDone
	} else if count[mo.StatusError] > 0 && count[mo.StatusError]+count[mo.StatusDone] == len(w.Task.SubTasks) {
		w.Task.FinishAt = time.Now()
		w.Task.Status = mo.StatusError
		w.Task.Err = errors.New("some subtask faild")
	}

	if w.AtSubTaskDone != nil {
		w.AtSubTaskDone(w.SubTask)
	}

	if w.AtTaskDone != nil && (w.Task.Status == mo.StatusDone || w.Task.Status == mo.StatusError) {
		w.AtTaskDone(w.Task)
	}

	return nil
}
