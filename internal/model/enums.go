package model

// project_status 表对应的状态名
type ProjectStatus int

const (
	ProjectStatusPending ProjectStatus = iota + 1
	ProjectStatusDoing
	ProjectStatusPaused
	ProjectStatusCompleted
	ProjectStatusCancelled
)

// process_status 表对应的状态名
type ProcessStatus int

const (
	ProcessStatusWaiting ProcessStatus = iota + 1
	ProcessStatusReady
	ProcessStatusDoing
	ProcessStatusPaused
	ProcessStatusRework
	ProcessStatusCompleted
	ProcessStatusCancelled
)

// todo_status 表对应的状态名
type TodoStatus int

const (
	TodoStatusPending TodoStatus = iota + 1
	TodoStatusDoing
	TodoStatusCompleted
	TodoStatusCancelled
)

// 角色名（roles 表）
type Role int

const (
	RoleAdmin Role = iota + 1
	RoleManager
	RoleEmployee
	RoleViewer
)
