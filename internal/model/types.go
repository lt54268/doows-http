package model

import "time"

type SetPermissionRequest struct {
	UserID   int  `json:"user_id"`
	IsCreate bool `json:"is_create"`
}

type APIResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type WorkspacePermission struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	IsCreate    bool      `json:"is_create"`
	WorkspaceID string    `json:"workspace_id"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

type CreateWorkspaceRequest struct {
	UserID int `json:"user_id"`
}

type ExternalAPIResponse struct {
	Workspace struct {
		Slug string `json:"slug"`
	} `json:"workspace"`
}

type NewThreadRequest struct {
	Slug   string `json:"slug"`
	Model  string `json:"model"`
	Avatar string `json:"avatar"`
}

type ChatHistoryRequest struct {
	WorkspaceSlug string `json:"workspaceSlug"`
	ThreadSlug    string `json:"threadSlug"`
}

type ChatHistory struct {
	History []ChatMessage `json:"history"`
}

type ChatMessage struct {
	Content string `json:"content"`
}

type HistoryChat struct {
	ID           int64     `json:"id"`
	SessionID    string    `json:"session_id"`
	Model        string    `json:"model"`
	UserID       string    `json:"user_id"`
	LastMessages string    `json:"last_messages"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
	Avatar       string    `json:"avatar"`
}
