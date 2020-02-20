package window

import (
	"path/filepath"

	sciter "github.com/sciter-sdk/go-sciter"
	sciterwindow "github.com/sciter-sdk/go-sciter/window"
)

// Window powered by sciter
type Window struct {
	EventHandlers map[string][]func(args ...*sciter.Value)

	*sciterwindow.Window
}

// New create a new window
func New(title, mainpage string) (window *Window, err error) {
	w, err := sciterwindow.New(sciter.DefaultWindowCreateFlag, sciter.DefaultRect)
	if err != nil {
		return
	}

	fp, _ := filepath.Abs(mainpage)
	err = w.LoadFile(fp)
	if err != nil {
		return
	}

	w.SetTitle(title)
	window = &Window{map[string][]func(args ...*sciter.Value){}, w}

	return
}

// Init window
func (w *Window) Init() {
	w.SetOption(sciter.SCITER_SET_SCRIPT_RUNTIME_FEATURES, sciter.ALLOW_EVAL|sciter.ALLOW_FILE_IO)
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
