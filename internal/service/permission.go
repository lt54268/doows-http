package service

import (
	"database/sql"
	"doows/internal/model"
	"doows/internal/repository"
	"fmt"
)

// 更新数据库中的权限设置
func SetPermission(req model.SetPermissionRequest) (model.APIResponse, error) {
	isCreateVal := 0
	if req.IsCreate {
		isCreateVal = 1
	}

	query := "UPDATE pre_workspace_permissions SET is_create = ? WHERE user_id = ?"
	result, err := repository.DB.Exec(query, isCreateVal, req.UserID)
	if err != nil {
		return model.APIResponse{}, fmt.Errorf("error updating permission: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.APIResponse{}, fmt.Errorf("error getting affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return model.APIResponse{Message: "No rows affected", Success: false}, nil
	}

	return model.APIResponse{Message: "Permission updated successfully", Success: true}, nil
}

// 检查 workspace_id 不为空的记录的数量
func CheckWorkspacePermissions(db *sql.DB) (int, error) {
	query := "SELECT COUNT(*) FROM pre_workspace_permissions WHERE workspace_id IS NOT NULL"
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error querying count: %v", err)
	}

	return count, nil
}

func GetUsersWithCreatePermission(userID int) (bool, error) {
	query := "SELECT is_create FROM pre_workspace_permissions WHERE user_id = ?"
	var isCreate int
	err := repository.DB.QueryRow(query, userID).Scan(&isCreate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error querying user create permission: %v", err)
	}

	return isCreate == 1, nil
}
