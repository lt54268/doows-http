package api

import (
	"doows/internal/model"
	"doows/internal/repository"
	"doows/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// 主页路由
func handleIndex(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Hello, dootask!")
}

// 同步路由
// func handleSync(w http.ResponseWriter, r *http.Request) {
// 	setupCORS(&w, r)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}
// 	if r.Method != "GET" {
// 		JsonResponse(w, map[string]string{"error": "Only GET method is allowed"}, http.StatusMethodNotAllowed)
// 		return
// 	}

// 	go service.SyncUsers()
// 	response := struct {
// 		Message string `json:"message"`
// 	}{
// 		Message: "Sync started",
// 	}

// 	JsonResponse(w, response, http.StatusOK)
// }

// 设置权限路由
func handleSetPermission(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.SetPermissionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	service.SyncUsers()

	if req.IsCreate {
		count, err := service.CheckWorkspacePermissions(repository.DB)
		if err != nil {
			JsonResponse(w, map[string]string{"error": "Failed to check workspace permissions"}, http.StatusInternalServerError)
			return
		}

		// 允许的最大工作区数量
		const maxAllowed = 3
		if count >= maxAllowed {
			JsonResponse(w, map[string]string{"error": "The limit of non-empty workspace_ids has been reached"}, http.StatusForbidden)
			return
		}
	}

	resp, err := service.SetPermission(req)
	if err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, resp, http.StatusOK)
}

// 统计已创建工作区路由
func handleCheckWorkspaceID(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "GET" {
		JsonResponse(w, map[string]string{"error": "Only GET method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	count, err := service.CheckWorkspacePermissions(repository.DB)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to fetch data"}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]int{"count": count}, http.StatusOK)
}

// 创建工作区路由
func handleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	var req model.CreateWorkspaceRequest
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	hasPermission, err := service.CheckUserCreatePermission(req.UserID)
	if err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	if !hasPermission {
		JsonResponse(w, map[string]string{"error": "User does not have permission to create workspace"}, http.StatusForbidden)
		return
	}

	workspaceName := fmt.Sprintf("Workspace for User %d", req.UserID)
	slug, err := service.CreateExternalWorkspace(workspaceName)
	if err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if err := service.UpdateWorkspaceID(req.UserID, slug); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{"slug": slug}, http.StatusOK)
}

// 删除工作区路由
func handleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "DELETE" {
		JsonResponse(w, map[string]string{"error": "Only DELETE method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	slug, err := service.GetWorkspaceSlug(req.UserID)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to retrieve workspace slug"}, http.StatusInternalServerError)
		return
	}

	if err := service.DeleteExternalWorkspace(slug); err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to delete workspace"}, http.StatusInternalServerError)
		return
	}

	if err := service.ResetWorkspaceID(req.UserID); err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to reset workspace ID"}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{"message": "Workspace deleted successfully"}, http.StatusOK)
}

// 新建会话路由
func handleNewThread(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.NewThreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	threadSlug, err := service.CreateNewThread(req.Slug)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to create new thread"}, http.StatusInternalServerError)
		return
	}

	userID, err := service.ExtractUserID(req.Slug)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to extract user ID"}, http.StatusInternalServerError)
		return
	}

	if err := service.StoreChatData(req.Model, req.Avatar, threadSlug, userID); err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to store chat data"}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{"thread_slug": threadSlug}, http.StatusOK)
}

// 获取用户是否有权限路由
func handleGetWorkspaceUsers(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateWorkspaceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	isCreate, err := service.GetUsersWithCreatePermission(req.UserID)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to retrieve user create permission"}, http.StatusInternalServerError)
		return
	}

	response := map[string]bool{"is_create": isCreate}
	JsonResponse(w, response, http.StatusOK)
}

// 更新最后一条对话路由
func handleUpdateLastChat(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.ChatHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	chatHistory, err := service.FetchChatHistory(req.WorkspaceSlug, req.ThreadSlug)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to fetch chat history"}, http.StatusInternalServerError)
		return
	}

	lastMessage := chatHistory[len(chatHistory)-1]
	chat, err := service.UpdateLastMessage(req.ThreadSlug, lastMessage.Content)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to store last message"}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, chat, http.StatusOK)
}

// 获取对话列表路由
func handleGetChatList(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	chats, err := service.FetchChatsByUserID(req.UserID)
	if err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, chats, http.StatusOK)
}

// 获取当前会话路由
func handleGetSession(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	session, err := service.FetchLatestSessionIDByUserID(req.UserID)
	if err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, session, http.StatusOK)
}

// 删除会话路由
func handleDeleteSession(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "DELETE" {
		JsonResponse(w, map[string]string{"error": "Only DELETE method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.ChatHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	if err := service.DeleteExternalSession(req.WorkspaceSlug, req.ThreadSlug); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if err := service.DeleteSessionFromDatabase(req.ThreadSlug); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{"message": "Session deleted successfully"}, http.StatusOK)
}

// 判断是否为管理员路由
func handleIsAdmin(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": "Invalid request body"}, http.StatusBadRequest)
		return
	}

	isAdmin, err := service.CheckIfUserIsAdmin(req.UserID)
	if err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]bool{"is_admin": isAdmin}, http.StatusOK)
}
