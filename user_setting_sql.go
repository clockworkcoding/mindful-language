package main

import (
	"fmt"
	"log"
)

func getUserSettings(userID string) (settings map[int]int) {
	settings = make(map[int]int)
	rows, err := db.Query("SELECT triggerid, responsetype FROM user_settings WHERE userid = ?", userID)

	if err != nil {
		log.Output(0, fmt.Sprintf("user_settings storage Err: %s", err.Error()))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var triggerID, responseType int
		rows.Scan(&triggerID, &responseType)
		settings[triggerID] = responseType
	}
	return
}

func insertUserSetting(setting userSetting) {

	statement, err := db.Prepare("DELETE from user_settings WHERE userid = ? AND triggerid = ?")
	if err != nil {
		log.Output(0, fmt.Sprintf("user_settings storage Err: %s", err.Error()))
		return
	}

	_, err = statement.Exec(setting.UserID, setting.TriggerID)
	if err != nil {
		log.Output(0, fmt.Sprintf("user_settings storage Err: %s", err.Error()))
		return
	}
	statement, err = db.Prepare("INSERT INTO user_settings (userid, triggerid, responsetype) values(?,?,?)")
	if err != nil {
		log.Output(0, fmt.Sprintf("user_settings storage Err: %s", err.Error()))
		return
	}
	_, err = statement.Exec(setting.UserID, setting.TriggerID, setting.ResponseType)
	if err != nil {
		log.Output(0, fmt.Sprintf("user_settings storage Err: %s", err.Error()))
		return
	}
}

func createUserSettingsTable() {
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS user_settings (userid TEXT, triggerid INT, responsetype INT)")
	statement.Exec()
}
