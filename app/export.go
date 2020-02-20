package app

import (
	sciter "github.com/sciter-sdk/go-sciter"
)

// 初始化导出到sciter的函数
func initExportedFunctions(app *App) {
	w := app.window
	def := map[string]func(args ...*sciter.Value) *sciter.Value{}

	// postEvent 发送事件到go侧（供tis调用）①
	def["postEvent"] = func(args ...*sciter.Value) *sciter.Value {
		event := args[0].String()
		params := args[1:]

		if handlers, ok := w.EventHandlers[event]; ok {
			for _, h := range handlers {
				fn := h
				ps := copyParams(params)
				go func() {
					fn(ps...)
				}()
			}
		}
		return sciter.NullValue()
	}

	for name, f := range def {
		w.DefineFunction(name, f)
	}
}

// deepcopy sciter Value slice
func copyParams(src []*sciter.Value) (dst []*sciter.Value) {
	for _, p := range src {
		dst = append(dst, p.Clone())
	}

	return
}
