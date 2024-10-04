package service

import (
	"bytes"
	"database/sql"
	"doows/internal/model"
	"doows/internal/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 发送请求到外部 API 创建工作区
func CreateExternalWorkspace(name string) (string, error) {
	url := "http://103.63.139.165:3001/api/v1/workspace/new"
	payload := map[string]string{"name": name}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer CM34YVB-3HJM2RS-PRGK1D2-ECZD4R6")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var apiResp model.ExternalAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	return apiResp.Workspace.Slug, nil
}

// 更新 workspace_id
func UpdateWorkspaceID(userID int, slug string) error {
	query := `UPDATE workspace_permission SET workspace_id = ? WHERE user_id = ?`
	_, err := repository.DB.Exec(query, slug, userID)
	if err != nil {
		return fmt.Errorf("failed to update workspace_id: %v", err)
	}
	return nil
}

func DeleteExternalWorkspace(slug string) error {
	url := fmt.Sprintf("http://103.63.139.165:3001/api/v1/workspace/%s", slug)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer CM34YVB-3HJM2RS-PRGK1D2-ECZD4R6")
	req.Header.Set("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete workspace with status: %d", resp.StatusCode)
	}

	return nil
}

func ResetWorkspaceID(userID int) error {
	query := `UPDATE workspace_permission SET workspace_id = NULL WHERE user_id = ?`
	_, err := repository.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset workspace_id: %v", err)
	}
	return nil
}

// 根据 userID 获取 workspace_id (slug)
func GetWorkspaceSlug(userID int) (string, error) {
	var slug string
	query := `SELECT workspace_id FROM workspace_permission WHERE user_id = ?`
	err := repository.DB.QueryRow(query, userID).Scan(&slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no workspace found for user ID %d", userID)
		}
		return "", fmt.Errorf("error querying workspace slug: %v", err)
	}
	return slug, nil
}
