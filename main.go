package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"./app"
	"./app/model"
	rice "github.com/GeertJohan/go.rice"
)

var (
	meta = &model.AppMeta{
		AppName:         "multi-Downloader",
		Version:         "0.0.1",
		Description:     "多功能下载器项目，可通过解析器来适配更多网站支持",
		MainPage:        "rice://res/App.htm",
		ConfigPath:      "./Config.json",
		TaskStoragePath: "./TaskStorage.json",
		LogPath:         "./app.log",
		LogLevel:        "DEBUG",
		IsDebug:         true,
	}
)

func main() {
	rice.MustFindBox("res")

	svr, err := app.New(meta)
	if err != nil {
		panic(err)
	}

	svr.Init()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				msg := fmt.Sprintf("signal received:%s", s.String())
				os.Exit(svr.Exit(msg))
			case syscall.SIGHUP:
			default:
				return
			}
		}
	}()

	svr.Run()
}

func init() {
	runtime.LockOSThread()
}
