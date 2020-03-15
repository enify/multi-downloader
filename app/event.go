package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	mo "./model"
	"./parser"
	"./worker"
	"github.com/sciter-sdk/go-sciter"
)

// 初始化sciter侧事件监听器
func initEventHandlers(app *App) {
	w := app.window

	subTaskDoneEventWrap := func(wraped func(subtask *mo.SubTask) error) func(subtask *mo.SubTask) error {
		return func(subtask *mo.SubTask) error {
			app.ts.Save()
			return wraped(subtask)
		}
	}

	taskDoneEventWrap := func(wraped func(task *mo.Task) error) func(task *mo.Task) error {
		return func(task *mo.Task) error {
			app.lg.Info("taskDone: id:%s, title:%s", task.ID, task.Title)
			if task.Status == mo.StatusDone && app.conf.NotifyAtTaskDone {
				w.NotifyTask(task)
			} else if task.Status == mo.StatusError && app.conf.NotifyAtTaskError {
				w.NotifyTask(task)
			}
			return wraped(task)
		}
	}

	w.AddEventHandler("debug-msg", func(args ...*sciter.Value) {
		app.lg.Debug("msg from tis:%v", args)
	})

	w.AddEventHandler("close-request", func(args ...*sciter.Value) {
		app.Exit("exit from window")
	})

	w.AddEventHandler("task-added", func(args ...*sciter.Value) {
		url := args[0].String()

		var support = false
		var pr parser.Parser
		for _, p := range app.parsers {
			exp := regexp.MustCompile(p.GetMeta().URLRgx)
			if exp.MatchString(url) {
				support = true
				pr = p
				break
			}
		}

		if !support {
			w.Toast("warn", "不支持的链接格式！")
			return
		}

		if app.ts.HasTask(url) {
			w.Toast("warn", "任务:%s已存在", url)
			return
		}

		var task = &mo.Task{
			ID:         strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
			URL:        url,
			Status:     mo.StatusPending,
			Path:       app.conf.SaveDir,
			Meta:       map[string]string{},
			CreateAt:   time.Now(),
			ParserName: pr.GetMeta().InternalName,
			SubTasks:   []*mo.SubTask{},
		}
		app.ts.AddTask(task)
		app.lg.Info("AddTask: id:%s, url:%s", task.ID, url)

		err := pr.Prepare(task, app.httpclient)
		if err != nil {
			task.Status = mo.StatusError
			task.Err = fmt.Errorf("parseTask: %w", err)
			w.Toast("warn", "任务:%s 解析失败!解析器:%s E:%s", task.URL, pr.GetMeta().Name, err)
			app.lg.Warn("parseTask: id:%s, parser:%s E:%s", task.ID, pr.GetMeta().InternalName, err)
		} else {
			task.Status = mo.StatusRunning
			app.lg.Info("parseTask: id:%s, parser:%s", task.ID, pr.GetMeta().InternalName)
			for _, subtask := range task.SubTasks {
				app.wp.Submit(&worker.TaskDownloadWork{
					Task:          task,
					SubTask:       subtask,
					Client:        app.httpclient,
					AtTaskDone:    taskDoneEventWrap(pr.AtTaskDone),
					AtSubTaskDone: subTaskDoneEventWrap(pr.AtSubTaskDone),
				})
				app.lg.Debug("submitTask: id:%s, subtask:%s", task.ID, subtask.FileName)
			}

		}
		app.ts.Save()
	})

	w.AddEventHandler("task-started", func(args ...*sciter.Value) {
		list := args[0]

		var ids = []string{}
		for i := 0; i < list.Length(); i++ {
			ids = append(ids, list.Index(i).String())
		}

		for _, id := range ids {
			task := app.ts.Find(id)
			if task != nil {
				pr := app.parsers[task.ParserName]
				if task.Status == mo.StatusPause {
					app.lg.Info("startTask: type:paused, id:%s, do: re submut subtask", task.ID)
					task.Status = mo.StatusRunning
					for _, subtask := range task.SubTasks {
						app.wp.Submit(&worker.TaskDownloadWork{
							Task:          task,
							SubTask:       subtask,
							Client:        app.httpclient,
							AtTaskDone:    taskDoneEventWrap(pr.AtTaskDone),
							AtSubTaskDone: subTaskDoneEventWrap(pr.AtSubTaskDone),
						})
						app.lg.Debug("submitTask: id:%s, subtask:%s", task.ID, subtask.FileName)
					}
				} else if task.Status == mo.StatusError {
					app.lg.Info("startTask: type:error, id:%s, do: re start task", task.ID)
					task.Status = mo.StatusPending
					task.Path = app.conf.SaveDir // clean task
					task.Meta = map[string]string{}
					task.FinishAt = time.Time{}
					task.SubTasks = []*mo.SubTask{}
					task.Err = nil
					err := pr.Prepare(task, app.httpclient)
					if err != nil {
						task.Status = mo.StatusError
						task.Err = fmt.Errorf("parseTask: %w", err)
						w.Toast("warn", "任务:%s 再解析失败!解析器:%s E:%s", task.URL, pr.GetMeta().Name, err)
						app.lg.Warn("parseTask: re parse id:%s, parser:%s, E:%s", task.ID, pr.GetMeta().InternalName, err)
					} else {
						task.Status = mo.StatusRunning
						app.lg.Info("parseTask: re parse id:%s, parser:%s", task.ID, pr.GetMeta().InternalName)
						for _, subtask := range task.SubTasks {
							app.wp.Submit(&worker.TaskDownloadWork{
								Task:          task,
								SubTask:       subtask,
								Client:        app.httpclient,
								AtTaskDone:    taskDoneEventWrap(pr.AtTaskDone),
								AtSubTaskDone: subTaskDoneEventWrap(pr.AtSubTaskDone),
							})
							app.lg.Debug("submitTask: id:%s, subtask:%s", task.ID, subtask.FileName)
						}
					}
				}
			}
			app.ts.Save()
		}
	})

	w.AddEventHandler("task-paused", func(args ...*sciter.Value) {
		list := args[0]

		var ids = []string{}
		for i := 0; i < list.Length(); i++ {
			ids = append(ids, list.Index(i).String())
		}

		for _, id := range ids {
			task := app.ts.Find(id)
			if task != nil {
				if task.Status == mo.StatusRunning {
					task.Status = mo.StatusPause
					app.lg.Info("pauseTask: id%s", task.ID)
				}
			}
			app.ts.Save()
		}
	})

	w.AddEventHandler("task-deleted", func(args ...*sciter.Value) {
		list := args[0]

		var ids = []string{}
		for i := 0; i < list.Length(); i++ {
			ids = append(ids, list.Index(i).String())
		}

		for _, id := range ids {
			task := app.ts.Find(id)
			if task != nil {
				app.ts.DeleteTask(task)
				app.lg.Info("deleteTask: id:%s", task.ID)
			}
		}
	})

	w.AddEventHandler("task-removed", func(args ...*sciter.Value) {
		list := args[0]

		var ids = []string{}
		for i := 0; i < list.Length(); i++ {
			ids = append(ids, list.Index(i).String())
		}

		for _, id := range ids {
			task := app.ts.Find(id)
			if task != nil {
				app.ts.DeleteTask(task)
				if task.Path == app.conf.SaveDir {
					for _, subtask := range task.SubTasks {
						os.Remove(filepath.Join(app.conf.SaveDir, subtask.FileName))
					}
				} else if filepath.Dir(task.Path) == app.conf.SaveDir {
					for _, subtask := range task.SubTasks {
						os.Remove(filepath.Join(task.Path, subtask.FileName))
					}
					os.Remove(task.Path)
				}
				app.lg.Info("removeTask: id:%s", task.ID)
			}
		}
	})

	w.AddEventHandler("config-changed", func(args ...*sciter.Value) {
		data := args[0]

		err := data.ConvertToString(sciter.CVT_JSON_LITERAL)
		if err != nil {
			w.Toast("warn", "配置格式有误：%s", err)
		} else {
			err := json.Unmarshal([]byte(data.String()), app.conf)
			if err != nil {
				w.Toast("warn", "配置保存失败：%s", err)
			} else {
				app.conf.Save()
				w.Toast("info", "配置保存成功，重启后生效")
				app.lg.Info("changeConfig: done")
			}
		}
	})

}
