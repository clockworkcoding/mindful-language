package main

import (
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handleMentionEvent(event slackevents.EventsAPIEvent, messageEvent *slackevents.AppMentionEvent) {

	headerText := slack.NewTextBlockObject("mrkdwn", "Hi there!  \nAnybody can add or edit responses with the buttons below. You can also disable or delete responses in the edit menu. You might see some connection error messages because @max is still working on me, but it should all work ok! \nIf you want to test out responses, #response-bot is a great place for it.\nTo see this again, just tag me in any channel!", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	addResponseBtnTxt := slack.NewTextBlockObject("plain_text", "add a response", false, false)
	addResponseBtn := slack.NewButtonBlockElement("add_response", "add_response", addResponseBtnTxt)
	editResponseBtnTxt := slack.NewTextBlockObject("plain_text", "edit a response", false, false)
	editResponseBtn := slack.NewButtonBlockElement("edit_response", "edit_response", editResponseBtnTxt)
	buttons := slack.NewActionBlock("mention_buttons", addResponseBtn, editResponseBtn)

	blocks := slack.MsgOptionBlocks(headerSection, buttons)
	options := []slack.MsgOption{blocks}

	if messageEvent.ThreadTimeStamp != "" {
		options = append(options, slack.MsgOptionTS(messageEvent.ThreadTimeStamp))
	}
	_, _, err := api.PostMessage(messageEvent.Channel, options...)
	if err != nil {
		log.Println("Error posting: ", err)
		return

	}
}
