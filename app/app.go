package app

import (
	"errors"

	rice "github.com/GeertJohan/go.rice"
	mo "github.com/enify/multi-downloader/app/model"
	"github.com/enify/multi-downloader/app/parser"
	"github.com/enify/multi-downloader/app/request"
	"github.com/enify/multi-downloader/app/window"
	"github.com/enify/multi-downloader/app/worker"
	ricepack "github.com/sciter-sdk/go-sciter/rice"
)

// App 应用主struct
type App struct {
	meta       *mo.AppMeta
	conf       *mo.AppConfig
	ts         *mo.TaskStorage
	lg         *Logger
	wp         *worker.Pool
	window     *window.Window
	httpclient *request.HTTPClient
	parsers    map[string]parser.Parser
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

	win, err := window.New()
	if err != nil {
		return
	}

	hc := request.NewHTTPClient()

	app = &App{
		meta:       meta,
		conf:       conf,
		ts:         ts,
		lg:         lg,
		wp:         wp,
		window:     win,
		httpclient: hc,
		parsers:    map[string]parser.Parser{},
	}

	return
}

// Init 初始化 App 配置
func (app *App) Init() {
	app.lg.Info("init app")

	initWorkerPool(app)
	initHTTPClient(app)
	initParsers(app)
	initTaskStatus(app.ts)

	initExportedFunctions(app)
	initEventHandlers(app)
	initWindow(app)

	if app.meta.IsDebug {
		app.meta.LogLevel = "DEBUG"
		app.lg.SetLevel(DEBUG)
		app.window.InitDebugOptions()
		rice.Debug = true
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

func initHTTPClient(app *App) {
	app.lg.Info("init httpclient: useproxy:%s", app.conf.UseProxy)
	c := app.httpclient
	c.SetUserAgent(app.conf.UserAgent)
	switch app.conf.UseProxy {
	case "", "off":
		//
	case "system":
		c.SetEnvProxy()
	case "user":
		c.SetProxy(app.conf.Proxy)
	}
}

func initParsers(app *App) {
	for _, pr := range parser.RegisteredParsers {
		meta := pr.GetMeta()
		app.parsers[meta.InternalName] = pr
		app.lg.Info("register parser: name:%s, internal name:%s, version:%s", meta.Name, meta.InternalName, meta.Version)
	}
}

// reset task status when app restart
func initTaskStatus(ts *mo.TaskStorage) {
	for _, task := range ts.Tasks {
		if task.Status == mo.StatusPending {
			task.Status = mo.StatusError
			task.Err = errors.New("app restarted")
		} else if task.Status == mo.StatusRunning {
			task.Status = mo.StatusPause
			for _, subtask := range task.SubTasks {
				if subtask.Status == mo.StatusRunning {
					subtask.Status = mo.StatusPending
				}
			}
		}
		ts.Save()
	}
}

func initWindow(app *App) {
	app.lg.Info("init window")
	ricepack.HandleDataLoad(app.window.Sciter)
	err := app.window.Init(app.meta.AppName, app.meta.MainPage)
	if err != nil {
		app.lg.Error("init window: E:%s", err)
	}
}

// Run app
func (app *App) Run() {
	app.wp.Run()
	app.window.Show()
	app.window.Run()
}

// Exit app
func (app *App) Exit(msg string) (code int) {
	app.lg.Info("exit: msg:%s", msg)

	app.lg.Info("close work pool...")
	app.wp.ShutDown()

	app.lg.Info("storage tasks...")
	app.ts.Save()

	app.lg.Info("close window...")
	app.window.Close()

	app.lg.Info("exited")
	return 0
}
