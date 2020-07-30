package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handleEvent(event slackevents.EventsAPIEvent, messageEvent *slackevents.MessageEvent) {
	teamTriggers := getTriggers(event.TeamID, true)
	userSettings := getUserSettings(messageEvent.User)
	rand.Seed(time.Now().UnixNano())
	msgText := strings.ToLower(messageEvent.Text)
	for _, t := range teamTriggers {
		for _, tword := range t.Triggers {
			responseType := userSettings[t.ID]
			if responseType >= 0 && containsTrigger(msgText, tword) {
				num := rand.Intn(len(t.Explanations))

				headerText := slack.NewTextBlockObject("mrkdwn", t.Explanations[num], false, false)
				headerSection := slack.NewSectionBlock(headerText, nil, nil)

				footerString := "Triggered because of \"" + tword + "\". Created by " + t.Creator
				if t.Editor != "" {
					footerString += " and last edited by " + t.Editor
				}
				footerText := slack.NewTextBlockObject("mrkdwn", footerString, false, false)
				footer := slack.NewContextBlock("", []slack.MixedElement{footerText}...)

				settingsBtnTxt := slack.NewTextBlockObject("plain_text", "my settings", false, false)
				settingsBtn := slack.NewButtonBlockElement("user_settings", strconv.Itoa(t.ID), settingsBtnTxt)
				buttons := slack.NewActionBlock("triggered_buttons", settingsBtn)

				blocks := slack.MsgOptionBlocks(headerSection, footer, buttons)
				unfurlLinks := slack.MsgOptionEnableLinkUnfurl()
				unfurlMedia := slack.MsgOptionDisableMediaUnfurl()
				options := []slack.MsgOption{blocks, unfurlLinks, unfurlMedia}

				if responseType == 0 {
					responseType = t.DefaultResponseType
				}

				switch responseType {
				case ephemeralResponse:
					//channel only visible to you
					_, err := api.PostEphemeral(messageEvent.Channel, messageEvent.User, options...)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				case channelResponse:
					if messageEvent.ThreadTimeStamp != "" {
						options = append(options, slack.MsgOptionTS(messageEvent.TimeStamp))
					}
					_, _, err := api.PostMessage(messageEvent.Channel, options...)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				case threadResponse:
					//thread
					options = append(options, slack.MsgOptionTS(messageEvent.TimeStamp))
					_, _, err := api.PostMessage(messageEvent.Channel, options...)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}
				case directMessageResponse:
					//im
					conversationParmas := slack.OpenConversationParameters{Users: []string{messageEvent.User}, ReturnIM: true}
					imChannel, _, _, err := api.OpenConversation(&conversationParmas)
					if err != nil {
						log.Println("Error posting: ", err)
						return
					}

					_, _, err = api.PostMessage(imChannel.Conversation.ID, options...)
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
