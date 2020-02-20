package main

import (
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
	svr.Run()
}
