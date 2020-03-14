package model

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// AppConfig 应用配置详情
type AppConfig struct {
	path string
	mux  sync.Mutex

	SaveDir           string `json:"save_dir"`
	MaxRoutines       int    `json:"max_routines"`
	UseProxy          string `json:"use_proxy"` // "off":关闭，"system":环境代理，"user":自定代理
	Proxy             string `json:"proxy"`
	UserAgent         string `json:"user_agent"`
	NotifyAtTaskDone  bool   `json:"notify_at_task_done"`
	NotifyAtTaskError bool   `json:"notify_at_task_error"`
}

// NewAppConfig 返回 AppConfig 指针对象
func NewAppConfig(path string) (conf *AppConfig, err error) {
	conf = &AppConfig{path: path}

	if !conf.FileExist() {
		conf.initDefaultAppConfig()
		err = conf.Save()
	} else {
		data, err := ioutil.ReadFile(conf.path)
		if err != nil {
			return conf, err
		}

		err = json.Unmarshal(data, conf)
	}

	return
}

// initDefaulAppConfig 初始化默认配置
func (c *AppConfig) initDefaultAppConfig() {
	c.SaveDir, _ = filepath.Abs("./Downloads")
	c.MaxRoutines = 10
	c.UseProxy = "off"
	c.NotifyAtTaskDone = true
	c.NotifyAtTaskError = true
}

// FileExist 检查配置文件是否存在
func (c *AppConfig) FileExist() bool {
	filePath, _ := filepath.Abs(c.path)
	_, err := os.Stat(filePath)

	return !os.IsNotExist(err)
}

// Save 保存配置信息到配置文件
func (c *AppConfig) Save() (err error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(c.path, data, 0666)
	if err != nil {
		return
	}

	return
}
