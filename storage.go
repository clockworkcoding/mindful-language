package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/slack-go/slack"
)

var (
	db *sql.DB
)

func init() {
	_, err := os.Stat(config.DbPath)
	if err != nil {
		log.Println("Database not found, attempting to create database")
	}
	db, err = sql.Open("sqlite3", config.DbPath)
	if err != nil {
		log.Fatalln("could not access database: ", err.Error())
	}
	createTriggerTable()
	createUserSettingsTable()
	incrementStartCounter()
}

func incrementStartCounter() {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS startcounter(counter INTEGER, version INTEGER)")
	if err != nil {
		log.Fatalln("create table error: " + err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln("create table error: " + err.Error())
	}

	row, err := db.Query("SELECT counter from startcounter")
	if err != nil {
		log.Fatalln("Counter read error: " + err.Error())
	}

	defer row.Close()
	count := 0
	if row.Next() {
		row.Scan(&count)
		row.Close()
	} else {
		db.Exec(fmt.Sprintf("insert into startcounter (counter, version) values (1, 1)"))
		return
	}
	count++

	db.Exec(fmt.Sprintf("update startcounter set counter = %d", count))

	log.Println(fmt.Sprintf("Start count: %d", count))
}

func saveSlackAuth(oAuth *slack.OAuthV2Response) (err error) {
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS slack_auth (
		id serial,
		team varchar(200),
		teamid varchar(20),
		token varchar(200),
		url varchar(200),
		configUrl varchar(200),
		channel varchar(200),
		channelid varchar(200),
		createdtime	timestamp
		)`); err != nil {
		fmt.Println("Error creating database table: " + err.Error())
		return
	}
	if _, err = db.Exec(`INSERT INTO slack_auth (
		team ,
		teamid,
		token ,
		url ,
		configUrl ,
		channel ,
		channelid,
		createdtime	)
		VALUES ($1,$2,$3,$4,$5,$6,$7, now())`, oAuth.Team.Name, oAuth.Team.ID,
		"", "", "", "", ""); err != nil {
		fmt.Println("Error saving slack auth: " + err.Error())
		return
	}

	return
}

func getSlackAuth(teamID string) (id int, token, channelid string, err error) {
	rows, err := db.Query("SELECT id, token, channelid FROM slack_auth WHERE teamid = $1 ORDER BY createdtime DESC FETCH FIRST 1 ROWS ONLY", teamID)
	if err != nil {
		log.Output(0, fmt.Sprintf("Storage Err: %s", err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &token, &channelid); err != nil {
			fmt.Println("Error scanning auth:" + err.Error())
			return
		}
		return
	}

	return id, token, channelid, errors.New("Team not found")
}
