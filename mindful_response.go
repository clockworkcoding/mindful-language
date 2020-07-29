package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handleEvent(event slackevents.EventsAPIEvent, messageEvent *slackevents.MessageEvent) {
	teamTriggers := getTriggers(event.TeamID, true)
	rand.Seed(time.Now().UnixNano())
	msgText := strings.ToLower(messageEvent.Text)
	for _, t := range teamTriggers {
		for _, tword := range t.Triggers {
			if containsTrigger(msgText, tword) {
				num := rand.Intn(len(t.Explanations))

				headerText := slack.NewTextBlockObject("mrkdwn", t.Explanations[num], false, false)
				headerSection := slack.NewSectionBlock(headerText, nil, nil)

				footerString := "Triggered by \"" + tword + "\". Created by " + t.Creator
				if t.Editor != "" {
					footerString += " and last edited by " + t.Editor
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
					options := []slack.MsgOption{blocks}
					if messageEvent.ThreadTimeStamp != "" {
						options = append(options, slack.MsgOptionTS(messageEvent.TimeStamp))
					}
					_, _, err := api.PostMessage(messageEvent.Channel, options...)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				case thread:
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
