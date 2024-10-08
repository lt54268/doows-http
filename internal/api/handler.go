package api

import (
	"doows/internal/model"
	"doows/internal/repository"
	"doows/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// 处理主页路由
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

// 处理同步路由
func handleSync(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "GET" {
		JsonResponse(w, map[string]string{"error": "Only GET method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	go service.SyncUsers()
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Sync started",
	}

	JsonResponse(w, response, http.StatusOK)
}

// 处理设置权限的路由
func handleSetPermission(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	// 检查 workspace_id 不为空的记录的总数
	count, err := service.CheckWorkspacePermissions(repository.DB)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to check workspace permissions"}, http.StatusInternalServerError)
		return
	}

	// 定义一个阈值
	const maxAllowed = 3
	if count >= maxAllowed {
		JsonResponse(w, map[string]string{"error": "The limit of non-empty workspace_ids has been reached"}, http.StatusForbidden)
		return
	}

	var req model.SetPermissionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, err := service.SetPermission(req)
	if err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, resp, http.StatusOK)
}

// 处理检查 workspace_id 的路由
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

// 处理创建工作区的请求
func handleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var req model.CreateWorkspaceRequest
	if r.Method != "POST" {
		JsonResponse(w, map[string]string{"error": "Only POST method is allowed"}, http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
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

func handleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
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

func handleNewThread(w http.ResponseWriter, r *http.Request) {
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

func handleGetWorkspaceUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		JsonResponse(w, map[string]string{"error": "Only GET method is allowed"}, http.StatusMethodNotAllowed)
		return
	}

	users, err := service.GetUsersWithCreatePermission(repository.DB)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to retrieve users"}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string][]int{"users": users}, http.StatusOK)
}
