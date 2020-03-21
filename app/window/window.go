package window

import (
	"encoding/json"
	"fmt"

	mo "github.com/enify/multi-downloader/app/model"
	sciter "github.com/sciter-sdk/go-sciter"
	sciterwindow "github.com/sciter-sdk/go-sciter/window"
)

// Window powered by sciter
type Window struct {
	EventHandlers map[string][]func(args ...*sciter.Value)

	*sciterwindow.Window
}

// New create a new window
func New() (window *Window, err error) {
	w, err := sciterwindow.New(sciter.DefaultWindowCreateFlag, sciter.DefaultRect)
	if err != nil {
		return
	}
	window = &Window{map[string][]func(args ...*sciter.Value){}, w}

	return
}

// Init window
func (w *Window) Init(title, mainpage string) (err error) {
	err = w.LoadFile(mainpage)
	if err != nil {
		return
	}
	w.SetTitle(title)

	w.SetOption(sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES, sciter.ALLOW_EVAL|sciter.ALLOW_FILE_IO|sciter.ALLOW_SYSINFO)
	return
}

// InitDebugOptions 调试模式下的设置
func (w *Window) InitDebugOptions() {
	w.SetOption(sciter.SCITER_SET_DEBUG_MODE, 1)
}

// AddEventHandler 在go侧监听tis侧事件①
func (w *Window) AddEventHandler(event string, fn func(args ...*sciter.Value)) {
	if handlers, ok := w.EventHandlers[event]; ok {
		handlers = append(handlers, fn)
		w.EventHandlers[event] = handlers
	} else {
		w.EventHandlers[event] = []func(args ...*sciter.Value){fn}
	}
}

// PostEvent 发送事件到tis侧②
func (w *Window) PostEvent(event string, data *sciter.Value) {
	var evt = sciter.NewValue()

	evt.SetString(event)
	w.Call("postEvent", evt, data)
}

// Toast 发送toast消息到tis侧
func (w *Window) Toast(mtype, msg string, v ...interface{}) {
	var data = sciter.NewValue()

	data.Set("type", mtype)
	data.Set("msg", fmt.Sprintf(msg, v...))
	w.PostEvent("toast", data)
}

// NotifyTask 任务完成时
func (w *Window) NotifyTask(task *mo.Task) {
	var r = sciter.NewValue()

	data, _ := json.Marshal(task)
	r.ConvertFromString(string(data), sciter.CVT_JSON_LITERAL)
	w.PostEvent("notify-task", r)
}

// Close tis端通过handle此事件来关闭窗口
func (w *Window) Close() {
	w.PostEvent("exit-app", sciter.NullValue())
}
