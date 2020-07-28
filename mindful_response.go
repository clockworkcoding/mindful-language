package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handleEvent(event slackevents.EventsAPIEvent, messageEvent *slackevents.MessageEvent) {
	teamTriggers := getTriggers(event.TeamID, true)
	for _, t := range teamTriggers {
		for _, tword := range t.Triggers {
			if containsTrigger(messageEvent.Text, tword) {

				headerText := slack.NewTextBlockObject("mrkdwn", t.Explanation, false, false)
				headerSection := slack.NewSectionBlock(headerText, nil, nil)

				footerString := "Created by @" + t.Creator
				if t.Editor != "" {
					footerString += " and last edited by @" + t.Editor
				}
				footerText := slack.NewTextBlockObject("mrkdwn", footerString, false, false)
				footer := slack.NewContextBlock("", []slack.MixedElement{footerText}...)

				settingsBtnTxt := slack.NewTextBlockObject("plain_text", "my settings", false, false)
				settingsBtn := slack.NewButtonBlockElement("user_settings", "user_settings", settingsBtnTxt)
				ignoreBtnTxt := slack.NewTextBlockObject("plain_text", "ignore when I write "+tword, false, false)
				ignoreBtn := slack.NewButtonBlockElement("ignore", "ignore", ignoreBtnTxt)
				buttons := slack.NewActionBlock("triggered_buttons", settingsBtn, ignoreBtn)

				blocks := slack.MsgOptionBlocks(headerSection, footer, buttons)

				switch t.DefaultResponseType {
				case ephemeral:
					//channel only visible to you
					_, err := api.PostEphemeral(messageEvent.Channel, messageEvent.User, blocks)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				case channel:
					//channel
					_, _, err := api.PostMessage(messageEvent.Channel, blocks)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				case thread:
					log.Println(messageEvent.TimeStamp, messageEvent.ThreadTimeStamp)
					//thread
					_, _, err := api.PostMessage(messageEvent.Channel, slack.MsgOptionTS(messageEvent.TimeStamp), blocks)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				case directMessage:
					//im
					conversationParmas := slack.OpenConversationParameters{Users: []string{messageEvent.User}, ReturnIM: true}
					imChannel, _, _, err := api.OpenConversation(&conversationParmas)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}

					_, _, err = api.PostMessage(imChannel.Conversation.ID, blocks)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				}
				break
			}
		}
	}
}

func getTriggers(teamID string, enabled bool) (triggers []trigger) {
	enabledString := ""
	if enabled {
		enabledString = " enabled = 1 AND "
	}
	rows, err := db.Query(fmt.Sprintf("SELECT id, trigger FROM triggers WHERE %s teamid = \"%s\"", enabledString, teamID))

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
