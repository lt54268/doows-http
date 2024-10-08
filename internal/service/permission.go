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

func GetUsersWithCreatePermission(db *sql.DB) ([]int, error) {
	query := "SELECT user_id FROM pre_workspace_permissions WHERE is_create <> 'false'"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
