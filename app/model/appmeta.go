package model

// AppMeta 应用元信息
type AppMeta struct {
	AppName         string `json:"appname"`
	Version         string `json:"version"`
	Description     string `json:"description"`
	MainPage        string `json:"main_page"`
	ConfigPath      string `json:"config_path"`
	TaskStoragePath string `json:"task_storage_path"`
	LogPath         string `json:"log_path"`
	LogLevel        string `json:"log_level"`
	IsDebug         bool   `json:"is_debug"`
}
