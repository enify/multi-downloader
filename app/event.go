package app

import (
	"github.com/sciter-sdk/go-sciter"
)

// 初始化sciter侧事件监听器
func initEventHandlers(app *App) {
	w := app.window

	w.AddEventHandler("data_from_tis", func(args ...*sciter.Value) {

		app.lg.Debug("process event: name:test, data:%v", args)
	})
}
