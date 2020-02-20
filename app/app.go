package app

import (
	"fmt"
	"time"

	mo "./model"
	"./parser"
	"./window"
	"./worker"
	"github.com/sciter-sdk/go-sciter"
)

// App 应用主struct
type App struct {
	meta    *mo.AppMeta
	conf    *mo.AppConfig
	ts      *mo.TaskStorage
	lg      *Logger
	wp      *worker.Pool
	window  *window.Window
	parsers map[string]parser.Parser
}

// New 返回 App 指针对象
func New(meta *mo.AppMeta) (app *App, err error) {
	conf, err := mo.NewAppConfig(meta.ConfigPath)
	if err != nil {
		return
	}

	ts, err := mo.NewTaskStorage(meta.TaskStoragePath)
	if err != nil {
		return
	}

	lg, err := NewLogger(meta.LogLevel, meta.LogPath)
	if err != nil {
		return
	}

	wp := worker.NewPool(conf.MaxRoutines, 10)

	win, err := window.New(meta.AppName, meta.MainPage)
	if err != nil {
		return
	}

	app = &App{
		meta:    meta,
		conf:    conf,
		ts:      ts,
		lg:      lg,
		wp:      wp,
		window:  win,
		parsers: map[string]parser.Parser{},
	}

	return
}

// Init 初始化 App 配置
func (app *App) Init() {
	app.lg.Info("init app")

	initWorkerPool(app)
	initWindow(app)
	initParsers(app)

	initExportedFunctions(app)
	initEventHandlers(app)

	if app.meta.IsDebug {
		app.meta.LogLevel = "DEBUG"
		app.window.InitDebugOptions()
	}
}

func initWorkerPool(app *App) {
	app.lg.Info("init work pool: size:%d", app.wp.MaxGoroutines)
	go func() {
		for err := range app.wp.Errchan {
			app.lg.Error("worker: %s", err)
		}
	}()
}

func initWindow(app *App) {
	app.lg.Info("init window")
	app.window.Init()
}

func initParsers(app *App) {
	for _, p := range parser.RegisteredParsers {
		meta := p.GetMeta()
		app.parsers[meta.URLRgx] = p
		app.lg.Info("register parser: name:%s, version:%s", meta.Name, meta.Version)
	}
}

// Run app
func (app *App) Run() {
	app.wp.Run()
	app.window.Show()
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("发送事件到tis")
		data := sciter.NewValue()
		data.ConvertFromString(`{"key":"value"}`, sciter.CVT_JSON_LITERAL)
		app.window.PostEvent("data_from_go", data)
	}()
	app.window.Run()
}
