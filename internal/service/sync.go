package service

import (
	"doows/internal/repository"
	"log"
	"time"
)

func SyncUsers() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Syncing users...")
		syncUsers()
	}
}

func syncUsers() {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	tx, err := repository.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return
	}

	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM workspace_permission WHERE user_id NOT IN (SELECT userid FROM pre_users)")
	if err != nil {
		log.Println("Failed to delete users:", err)
		return
	}

	_, err = tx.Exec(`
        INSERT INTO workspace_permission (user_id, create_time, update_time)
        SELECT userid, ?, ? FROM pre_users
        WHERE userid NOT IN (SELECT user_id FROM workspace_permission)
    `, currentTime, currentTime)
	if err != nil {
		log.Println("Failed to insert new users:", err)
		return
	}

	_, err = tx.Exec(`
        UPDATE workspace_permission 
        SET update_time = ?
        WHERE user_id IN (SELECT userid FROM pre_users)
    `, currentTime)
	if err != nil {
		log.Println("Failed to update users:", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit transaction:", err)
		return
	}

	log.Println("Sync completed successfully")
}
