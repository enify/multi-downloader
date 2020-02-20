package model

// TaskStatus set
const (
	StatusPending TaskStatus = "pending"
	StatusRunning TaskStatus = "running"
	StatusPause   TaskStatus = "pause"
	StatusError   TaskStatus = "error"
	StatusDone    TaskStatus = "done"
)

// TaskStatus 任务状态
type TaskStatus string
