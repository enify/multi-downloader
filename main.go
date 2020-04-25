package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	rice "github.com/GeertJohan/go.rice"
	"github.com/enify/multi-downloader/app"
	"github.com/enify/multi-downloader/app/model"
)

// will be inject at building
var (
	AppName = "Debug Mode"
	AppVer  = "0.0.0"
	AppDesc = "你正处于debug模式"
)

func main() {
	rice.MustFindBox("res")

	meta := &model.AppMeta{
		AppName:         AppName,
		Version:         AppVer,
		Description:     AppDesc,
		MainPage:        "rice://res/App.htm",
		ConfigPath:      "./Config.json",
		TaskStoragePath: "./TaskStorage.json",
		LogPath:         "./app.log",
		LogLevel:        "INFO",
		IsDebug:         AppVer == "0.0.0",
	}

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
