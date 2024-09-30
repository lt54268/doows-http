package api

import (
	"doows/internal/model"
	"doows/internal/repository"
	"doows/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// handleIndex 处理主页路由
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

// handleSync 处理同步路由
func handleSync(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	go service.SyncUsers() // 启动同步任务
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Sync started",
	}

	JsonResponse(w, response, http.StatusOK)
}

// handleSetPermission 处理设置权限的路由
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

// handleCheckWorkspaceID 处理检查 workspace_id 的路由
func handleCheckWorkspaceID(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	if r.Method == "OPTIONS" {
		return
	}
	count, err := service.CheckWorkspacePermissions(repository.DB)
	if err != nil {
		JsonResponse(w, map[string]string{"error": "Failed to fetch data"}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]int{"count": count}, http.StatusOK)
}

// handleCreateWorkspace 处理创建工作区的请求
func handleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var req model.CreateWorkspaceRequest
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

	// 将 slug 更新到数据库
	if err := service.UpdateWorkspaceID(req.UserID, slug); err != nil {
		JsonResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, map[string]string{"slug": slug}, http.StatusOK)
}

func handleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 从数据库获取 slug
	slug, err := service.GetWorkspaceSlug(req.UserID)
	if err != nil {
		http.Error(w, "Failed to retrieve workspace slug", http.StatusInternalServerError)
		return
	}

	// 调用外部 API 删除工作区
	if err := service.DeleteExternalWorkspace(slug); err != nil {
		http.Error(w, "Failed to delete workspace", http.StatusInternalServerError)
		return
	}

	// 将数据库中的 workspace_id 设置为 NULL
	if err := service.ResetWorkspaceID(req.UserID); err != nil {
		http.Error(w, "Failed to reset workspace ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Workspace deleted successfully"))
}
