package main

import (
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handleMentionEvent(event slackevents.EventsAPIEvent, messageEvent *slackevents.AppMentionEvent) {

	headerText := slack.NewTextBlockObject("mrkdwn", "What would you like to do?", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	addResponseBtnTxt := slack.NewTextBlockObject("plain_text", "add a response", false, false)
	addResponseBtn := slack.NewButtonBlockElement("add_response", "add_response", addResponseBtnTxt)
	editResponseBtnTxt := slack.NewTextBlockObject("plain_text", "edit a response", false, false)
	editResponseBtn := slack.NewButtonBlockElement("edit_response", "edit_response", editResponseBtnTxt)
	buttons := slack.NewActionBlock("mention_buttons", addResponseBtn, editResponseBtn)

	blocks := slack.MsgOptionBlocks(headerSection, buttons)
	options := []slack.MsgOption{blocks}

	if messageEvent.ThreadTimeStamp != "" {
		options = append(options, slack.MsgOptionTS(messageEvent.TimeStamp))
	}
	_, _, err := api.PostMessage(messageEvent.Channel, options...)
	if err != nil {
		log.Println("Error posting: ", err)
		return

	}
}
