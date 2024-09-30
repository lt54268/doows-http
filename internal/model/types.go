package model

import "time"

// SetPermissionRequest 定义用于`/set`路由的请求体结构
type SetPermissionRequest struct {
	UserID   int  `json:"user_id"`
	IsCreate bool `json:"is_create"`
}

// APIResponse 定义通用的API响应结构
type APIResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// WorkspacePermission 对应数据库中的 workspace_permission 表
type WorkspacePermission struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	IsCreate    bool      `json:"is_create"`
	WorkspaceID string    `json:"workspace_id"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

// CreateWorkspaceRequest 定义创建工作区的请求结构
type CreateWorkspaceRequest struct {
	UserID int `json:"user_id"`
}

// ExternalAPIResponse 定义外部 API 的响应结构
type ExternalAPIResponse struct {
	Workspace struct {
		Slug string `json:"slug"`
	} `json:"workspace"`
}
