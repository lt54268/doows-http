package api

import (
	"net/http"
)

// 设置所需要的 CORS headers
func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // 注意，这里设置为 "*" 仅适用于开发环境
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// 如果是 OPTIONS 请求，直接返回并结束处理
	if req.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return
	}
}

func SetupRoutes() {
	http.HandleFunc("/", handleIndex)
	// http.HandleFunc("/sync", handleSync)
	http.HandleFunc("/set", handleSetPermission)
	http.HandleFunc("/check", handleCheckWorkspaceID)
	http.HandleFunc("/create", handleCreateWorkspace)
	http.HandleFunc("/delete-ws", handleDeleteWorkspace)
	http.HandleFunc("/new", handleNewThread)
	http.HandleFunc("/get-user", handleGetWorkspaceUsers)
	http.HandleFunc("/update-last", handleUpdateLastChat)
	http.HandleFunc("/get-list", handleGetChatList)
	http.HandleFunc("/get-sessionid", handleGetSession)
	http.HandleFunc("/delete-session", handleDeleteSession)
	http.HandleFunc("/is-admin", handleIsAdmin)
}
