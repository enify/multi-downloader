package app

import (
	"encoding/json"

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

	def["getAppMeta"] = func(args ...*sciter.Value) *sciter.Value {
		var r = sciter.NewValue()

		data, _ := json.Marshal(app.meta)
		r.ConvertFromString(string(data), sciter.CVT_JSON_LITERAL)

		return r
	}

	def["getAppConfig"] = func(args ...*sciter.Value) *sciter.Value {
		var r = sciter.NewValue()

		data, _ := json.Marshal(app.conf)
		r.ConvertFromString(string(data), sciter.CVT_JSON_LITERAL)

		return r
	}

	def["getTasks"] = func(args ...*sciter.Value) *sciter.Value {
		var r = sciter.NewValue()

		data, _ := json.Marshal(app.ts.Tasks)
		r.ConvertFromString(string(data), sciter.CVT_JSON_LITERAL)

		return r
	}

	def["getUrlFilters"] = func(args ...*sciter.Value) *sciter.Value {
		var r = sciter.NewValue()

		filters := []string{}
		for _, pr := range app.parsers {
			filters = append(filters, pr.GetMeta().URLRgx)
		}
		data, _ := json.Marshal(filters)
		r.ConvertFromString(string(data), sciter.CVT_JSON_LITERAL)

		return r
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
