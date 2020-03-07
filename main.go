package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./app"
	"./app/model"
)

var (
	meta = &model.AppMeta{
		AppName:         "multi-Downloader",
		Version:         "0.0.1",
		MainPage:        "./res/App.htm",
		ConfigPath:      "./Config.json",
		TaskStoragePath: "./TaskStorage.json",
		LogPath:         "./app.log",
		LogLevel:        "DEBUG",
		IsDebug:         true,
	}
)

func main() {
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
