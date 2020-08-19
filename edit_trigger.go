package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
)

func editTriggerCommand(w http.ResponseWriter, r *http.Request) {
	err := verifySigningSecret(r)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	showEditTriggerModal(s.TriggerID, s.TeamID)

	params := &slack.Msg{Text: "Okay, let's set one up"}
	b, err := json.Marshal(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func showEditTriggerModal(triggerID string, teamID string) {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "Edit a response", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Edit", false, false)

	manageTxt := slack.NewTextBlockObject("plain_text", "Pick a response to edit", true, false)
	var options []*slack.OptionBlockObject
	triggers := getTriggers(teamID, false)
	for _, t := range triggers {
		text := t.Triggers[0] + " | " + t.Explanations[0]
		length := 74
		if len(text) < 74 {
			length = len(text)
		}
		triggerTxt := slack.NewTextBlockObject("plain_text", text[:length], true, false)
		triggerOpt := slack.NewOptionBlockObject(strconv.Itoa(t.ID), triggerTxt)
		options = append(options, triggerOpt)
	}
	triggerSelectOption := slack.NewOptionsSelectBlockElement("static_select", manageTxt, "trigger", options...)
	triggerSelect := slack.NewInputBlock("trigger", manageTxt, triggerSelectOption)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			triggerSelect,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	modalRequest.CallbackID = "trigger_edit_modal"

	_, err := api.OpenView(triggerID, modalRequest)
	if err != nil {
		fmt.Printf("Error opening view: %s", err)
	}
}

func handleEditTriggerModalSelect(i slack.InteractionCallback, w http.ResponseWriter, r *http.Request) {
	triggerID, err := strconv.Atoi(i.View.State.Values["trigger"]["trigger"].SelectedOption.Value)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	trigger := getTrigger(triggerID)
	triggers := strings.Join(trigger.Triggers, ", ")
	explanations := strings.Join(trigger.Explanations, "\n")

	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "Update response", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	triggerText := slack.NewTextBlockObject("plain_text", "Trigger Words (comma sepearted variations)", false, false)
	triggerPlaceholder := slack.NewTextBlockObject("plain_text", "toaster, toasters", false, false)
	triggerElement := slack.NewPlainTextInputBlockElement(triggerPlaceholder, "trigger_list")
	triggerElement.InitialValue = triggers

	triggerBlock := slack.NewInputBlock("trigger_list", triggerText, triggerElement)

	explanationText := slack.NewTextBlockObject("plain_text", "Responses (seperated by new lines, randomly chosen)", false, false)
	explanationPlaceholder := slack.NewTextBlockObject("plain_text", "The proper term is Cylon, please let them live peacefully among you.", false, false)
	explanationElement := &slack.PlainTextInputBlockElement{
		Type:         slack.METPlainTextInput,
		ActionID:     "explanation",
		Placeholder:  explanationPlaceholder,
		InitialValue: explanations,
		Multiline:    true,
	}
	explanation := slack.NewInputBlock("explanation", explanationText, explanationElement)

	manageTxt := slack.NewTextBlockObject("plain_text", "pick a default response type", true, false)
	threadTxt := slack.NewTextBlockObject("plain_text", "in a thread", true, false)
	channelTxt := slack.NewTextBlockObject("plain_text", "in the channel", false, false)
	ephemeralTxt := slack.NewTextBlockObject("plain_text", "in the channel, but only you see it", false, false)
	dmTxt := slack.NewTextBlockObject("plain_text", "in a direct message", false, false)

	threadOpt := slack.NewOptionBlockObject(strconv.Itoa(threadResponse), threadTxt)
	channelOpt := slack.NewOptionBlockObject(strconv.Itoa(channelResponse), channelTxt)
	ephemeralOpt := slack.NewOptionBlockObject(strconv.Itoa(ephemeralResponse), ephemeralTxt)
	dmOpt := slack.NewOptionBlockObject(strconv.Itoa(directMessageResponse), dmTxt)

	responseTypeOption := slack.NewOptionsSelectBlockElement("static_select", manageTxt, "response_type", threadOpt, channelOpt, ephemeralOpt, dmOpt)
	var selectedOption *slack.OptionBlockObject
	switch trigger.DefaultResponseType {
	case threadResponse:
		selectedOption = threadOpt
	case channelResponse:
		selectedOption = channelOpt
	case ephemeralResponse:
		selectedOption = ephemeralOpt
	case directMessageResponse:
		selectedOption = dmOpt
	}
	responseTypeOption.InitialOption = selectedOption
	responseType := slack.NewInputBlock("response_type", manageTxt, responseTypeOption)

	enabledTxt := slack.NewTextBlockObject("plain_text", "Enabled", false, false)
	enabledOpt := slack.NewOptionBlockObject("true", enabledTxt)
	enabledGroup := slack.NewCheckboxGroupsBlockElement("enabled", enabledOpt)
	if trigger.Enabled {
		enabledGroup.InitialOptions = enabledGroup.Options
	}
	enabledBlock := slack.NewInputBlock("enabled", enabledTxt, enabledGroup)
	enabledBlock.Optional = true

	deleteTxt := slack.NewTextBlockObject("plain_text", "Delete", false, false)
	deleteOpt := slack.NewOptionBlockObject("true", deleteTxt)
	deleteGroup := slack.NewCheckboxGroupsBlockElement("delete", deleteOpt)

	confirmTitleTxt := slack.NewTextBlockObject("plain_text", "Delete", false, false)
	confirmTxt := slack.NewTextBlockObject("plain_text", "This cannot be undone after saving. Are you sure?", false, false)
	confirmCancelTxt := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	confirmApproveTxt := slack.NewTextBlockObject("plain_text", "Delete", false, false)
	deleteGroup.Confirm = slack.NewConfirmationBlockObject(confirmTitleTxt, confirmTxt, confirmApproveTxt, confirmCancelTxt)
	deleteBlock := slack.NewInputBlock("delete", deleteTxt, deleteGroup)
	deleteBlock.Optional = true

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			triggerBlock,
			explanation,
			responseType,
			enabledBlock,
			deleteBlock,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	modalRequest.CallbackID = "trigger_update_save_modal"
	modalRequest.PrivateMetadata = strconv.Itoa(trigger.ID)

	response := slack.NewUpdateViewSubmissionResponse(&modalRequest)
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
	_, err = api.UpdateView(modalRequest, "", i.View.Hash, i.View.ID)
	//_, err := api.OpenView(i.TriggerID, modalRequest)
	if err != nil {
		fmt.Printf("Error opening view: %s", err)
	}
}

func handleEditTriggerModalSave(i slack.InteractionCallback) {
	triggerID, _ := strconv.Atoi(i.View.PrivateMetadata)
	triggerString := i.View.State.Values["trigger_list"]["trigger_list"].Value
	explanation := i.View.State.Values["explanation"]["explanation"].Value
	responseType, err := strconv.Atoi(i.View.State.Values["response_type"]["response_type"].SelectedOption.Value)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	var enabled bool

	if len(i.View.State.Values["enabled"]["enabled"].SelectedOptions) > 0 {
		enabled = true
	}
	if len(i.View.State.Values["delete"]["delete"].SelectedOptions) > 0 {
		deleteTrigger(triggerID)
		return
	}
	oldTrigger := getTrigger(triggerID)

	triggerList := strings.Split(strings.ToLower(triggerString), ",")
	for index, triggerWord := range triggerList {
		triggerList[index] = strings.TrimSpace(triggerWord)
	}

	newTrigger := trigger{
		ID:                  triggerID,
		Triggers:            triggerList,
		Explanations:        strings.Split(explanation, "\n"),
		Creator:             oldTrigger.Creator,
		Editor:              i.User.Name,
		DefaultResponseType: responseType,
		Enabled:             enabled,
	}

	updateTrigger(newTrigger)

	msg := fmt.Sprintf("I saved the response for %s and variants", newTrigger.Triggers[0])

	conversationParmas := slack.OpenConversationParameters{Users: []string{i.User.ID}, ReturnIM: true}
	imChannel, _, _, err := api.OpenConversation(&conversationParmas)
	if err != nil {
		log.Println("Error posting: ", err)
		return
	}

	_, err = api.PostEphemeral(imChannel.Conversation.ID, msg)
	if err != nil {
		log.Println("Error posting: ", err)
		return
	}

}
