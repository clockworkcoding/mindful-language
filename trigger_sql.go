package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func getTriggers(teamID string, enabled bool) (triggers []trigger) {
	enabledString := ""
	if enabled {
		enabledString = " enabled = 1 AND "
	}
	rows, err := db.Query(fmt.Sprintf("SELECT id, trigger FROM triggers WHERE %s teamid = ?", enabledString), teamID)

	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}

	defer rows.Close()
	var jsonTrigger string
	var id int
	for rows.Next() {
		var rowTrigger trigger
		rows.Scan(&id, &jsonTrigger)
		json.Unmarshal([]byte(jsonTrigger), &rowTrigger)
		rowTrigger.ID = id
		triggers = append(triggers, rowTrigger)
	}
	return
}

func getTrigger(triggerID int) (trigger trigger) {
	rows, err := db.Query("SELECT id, trigger FROM triggers WHERE id = ?", triggerID)

	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}

	defer rows.Close()
	var jsonTrigger string
	var id int
	for rows.Next() {
		rows.Scan(&id, &jsonTrigger)
		json.Unmarshal([]byte(jsonTrigger), &trigger)
		trigger.ID = id
	}
	return
}

func insertTrigger(teamid string, newTrigger trigger) {
	triggerJSON, _ := json.Marshal(newTrigger)

	statement, err := db.Prepare("INSERT INTO triggers (teamid, trigger, enabled) values(?,?,?)")
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}

	_, err = statement.Exec(teamid, triggerJSON, newTrigger.Enabled)
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}
}

func updateTrigger(trigger trigger) {
	triggerJSON, _ := json.Marshal(trigger)

	statement, err := db.Prepare("update triggers set trigger = ?, enabled = ? where id = ?")
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}

	_, err = statement.Exec(triggerJSON, trigger.Enabled, trigger.ID)
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}
}

func deleteTrigger(triggerID int) {
	statement, err := db.Prepare("delete from triggers where id = ?")
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}

	_, err = statement.Exec(triggerID)
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}
	statement, err = db.Prepare("delete from user_settings where triggerid = ?")
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}

	_, err = statement.Exec(triggerID)
	if err != nil {
		log.Output(0, fmt.Sprintf("Triggers storage Err: %s", err.Error()))
		return
	}
}

func createTriggerTable() {
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS triggers (id INTEGER PRIMARY KEY AUTOINCREMENT, teamid TEXT, trigger TEXTs, enabled INTEGER)")
	statement.Exec()
}
