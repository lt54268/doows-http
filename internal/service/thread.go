package service

import (
	"database/sql"
	"doows/internal/model"
	"doows/internal/repository"
	"doows/pkg/settime"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// 获取某个对话窗口的所有历史记录
func FetchChatHistory(workspaceSlug, threadSlug string) ([]model.ChatMessage, error) {
	url := fmt.Sprintf("http://103.63.139.165:3001/api/v1/workspace/%s/thread/%s/chats", workspaceSlug, threadSlug)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer CM34YVB-3HJM2RS-PRGK1D2-ECZD4R6")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var history model.ChatHistory
	if err := json.Unmarshal(body, &history); err != nil {
		return nil, err
	}

	return history.History, nil
}

// 更新最后一条消息并返回所有字段
func UpdateLastMessage(sessionID, lastMessage string) (model.HistoryChat, error) {
	currentTime := settime.GetCurrentFormattedTime()
	query := `UPDATE pre_history_aichats SET last_messages = ?, update_time = ? WHERE session_id = ?`
	_, err := repository.DB.Exec(query, lastMessage, currentTime, sessionID)
	if err != nil {
		log.Printf("Error updating last message for session_id %s: %v", sessionID, err)
		return model.HistoryChat{}, fmt.Errorf("error updating last message: %v", err)
	}

	return GetHistoryChat(sessionID)
}

// 根据 sessionID 获取记录
func GetHistoryChat(sessionID string) (model.HistoryChat, error) {
	var chat model.HistoryChat
	query := `SELECT id, session_id, model, user_id, last_messages, create_time, update_time, avatar FROM pre_history_aichats WHERE session_id = ?`
	err := repository.DB.QueryRow(query, sessionID).Scan(&chat.ID, &chat.SessionID, &chat.Model, &chat.UserID, &chat.LastMessages, &chat.CreateTime, &chat.UpdateTime, &chat.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.HistoryChat{}, fmt.Errorf("no record found: %v", err)
		}
		return model.HistoryChat{}, fmt.Errorf("error fetching record: %v", err)
	}
	return chat, nil
}

// 根据用户 ID 获取所有对话窗口的最后一条记录, 并按照更新时间降序排列
func FetchChatsByUserID(userID int) ([]model.HistoryChat, error) {
	var chats []model.HistoryChat
	query := `
			SELECT id, session_id, model, user_id, last_messages, create_time, update_time, avatar
			FROM pre_history_aichats WHERE user_id = ?
			ORDER BY update_time DESC
			`
	rows, err := repository.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying chats: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chat model.HistoryChat
		if err := rows.Scan(&chat.ID, &chat.SessionID, &chat.Model, &chat.UserID, &chat.LastMessages, &chat.CreateTime, &chat.UpdateTime, &chat.Avatar); err != nil {
			return nil, fmt.Errorf("error scanning chat: %v", err)
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// 根据用户 ID 查询最新 update_time 的 session_id 并返回所有字段
func FetchLatestSessionIDByUserID(userID int) (*model.HistoryChat, error) {
	var session model.HistoryChat
	query := `SELECT * FROM pre_history_aichats WHERE user_id = ? ORDER BY update_time DESC LIMIT 1`
	err := repository.DB.QueryRow(query, userID).Scan(&session.ID, &session.UserID, &session.SessionID, &session.Model, &session.LastMessages, &session.CreateTime, &session.UpdateTime, &session.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying latest session ID: %v", err)
	}
	return &session, nil
}

// 删除某个会话
func DeleteExternalSession(workspaceSlug, threadSlug string) error {
	url := fmt.Sprintf("http://103.63.139.165:3001/api/v1/workspace/%s/thread/%s", workspaceSlug, threadSlug)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer CM34YVB-3HJM2RS-PRGK1D2-ECZD4R6")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete session: %s", body)
	}

	return nil
}

// 从数据库中删除该会话
func DeleteSessionFromDatabase(sessionID string) error {
	query := "DELETE FROM pre_history_aichats WHERE session_id = ?"
	_, err := repository.DB.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("error deleting session from database: %v", err)
	}
	return nil
}
