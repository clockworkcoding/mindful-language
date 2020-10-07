package main

import (
	"log"
	"strconv"

	"github.com/slack-go/slack"
)

type userSetting struct {
	UserID       string
	TriggerID    int
	ResponseType int
}

func handleUserSettingAction(payload slack.InteractionCallback) {
	action := payload.ActionCallback.BlockActions[0]
	headerText := slack.NewTextBlockObject("mrkdwn", "How would you like responses to your use of this word to be posted?", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	channelBtnTxt := slack.NewTextBlockObject("plain_text", "in channel", false, false)
	channelBtn := slack.NewButtonBlockElement(strconv.Itoa(channelResponse), action.Value, channelBtnTxt)
	ephemeralBtnTxt := slack.NewTextBlockObject("plain_text", "only visible to me in channel", false, false)
	ephemeralBtn := slack.NewButtonBlockElement(strconv.Itoa(ephemeralResponse), action.Value, ephemeralBtnTxt)
	threadBtnTxt := slack.NewTextBlockObject("plain_text", "thread the response", false, false)
	threadBtn := slack.NewButtonBlockElement(strconv.Itoa(threadResponse), action.Value, threadBtnTxt)
	directMessageBtnTxt := slack.NewTextBlockObject("plain_text", "direct message", false, false)
	directMessageBtn := slack.NewButtonBlockElement(strconv.Itoa(directMessageResponse), action.Value, directMessageBtnTxt)
	noneBtnTxt := slack.NewTextBlockObject("plain_text", "don't show", false, false)
	noneBtn := slack.NewButtonBlockElement(strconv.Itoa(noResponse), action.Value, noneBtnTxt)
	buttons := slack.NewActionBlock("user_setting_selection", channelBtn, ephemeralBtn, threadBtn, directMessageBtn, noneBtn)
  
  if payload.User.ID == payload.User.ID { // "UG5DH19EX"{
	  deleteBtnTxt := slack.NewTextBlockObject("plain_text", "delete this instance", false, false)
    deleteBtn := slack.NewButtonBlockElement(strconv.Itoa(deleteInstance), payload.ResponseURL, deleteBtnTxt)
	  buttons = slack.NewActionBlock("user_setting_selection", channelBtn, ephemeralBtn, threadBtn, directMessageBtn, noneBtn, deleteBtn)
  }
	blocks := slack.MsgOptionBlocks(headerSection, buttons)
	options := []slack.MsgOption{blocks}
	if payload.Message.ThreadTimestamp != "" {
		options = append(options, slack.MsgOptionTS(payload.Message.ThreadTimestamp))
	}

	_, err := api.PostEphemeral(payload.Container.ChannelID, payload.User.ID, options...)
	if err != nil {
		log.Println("Error posting: ", err)
		return
	}
}

func handleUserSettingSelection(payload slack.InteractionCallback) {
	action := payload.ActionCallback.BlockActions[0]
	responseType, _ := strconv.Atoi(action.ActionID)
  if responseType == deleteInstance{

    _, err := api.PostEphemeral(payload.Container.ChannelID, payload.User.ID, slack.MsgOptionText(":heavy_check_mark:", false), slack.MsgOptionDeleteOriginal(payload.ResponseURL))
    if err != nil {
      log.Println("Error posting: ", err)
      return
	  }
	  _, err = api.PostEphemeral(payload.Container.ChannelID, payload.User.ID, slack.MsgOptionText(":heavy_check_mark:", false), slack.MsgOptionDeleteOriginal(action.Value))
	  if err != nil {
		  log.Println("Error posting: ", err)
		  return
	  }
    return
  }
	triggerID, _ := strconv.Atoi(action.Value)
	setting := userSetting{
		ResponseType: responseType,
		TriggerID:    triggerID,
		UserID:       payload.User.ID,
	}
  
	insertUserSetting(setting)

	_, err := api.PostEphemeral(payload.Container.ChannelID, payload.User.ID, slack.MsgOptionText(":heavy_check_mark:", false), slack.MsgOptionDeleteOriginal(payload.ResponseURL))
	if err != nil {
		log.Println("Error posting: ", err)
		return
	}
}
