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

type trigger struct {
	Triggers            []string
	Explanation         string
	DefaultResponseType int
	Creator             string
	Editor              string
	ID                  int
	Enabled             bool
}

func addTriggerCommand(w http.ResponseWriter, r *http.Request) {

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

	showTriggerModal(s.TriggerID)

	params := &slack.Msg{Text: s.Text}
	b, err := json.Marshal(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}

func showTriggerModal(triggerID string) {

	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", "Add a response", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Please enter your name", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	triggerText := slack.NewTextBlockObject("plain_text", "trigger Words (comma sepearted variations)", false, false)
	triggerPlaceholder := slack.NewTextBlockObject("plain_text", "toaster, toasters", false, false)
	triggerElement := slack.NewPlainTextInputBlockElement(triggerPlaceholder, "trigger_list")
	// Notice that blockID is a unique identifier for a block
	trigger := slack.NewInputBlock("trigger_list", triggerText, triggerElement)

	explanationText := slack.NewTextBlockObject("plain_text", "Response", false, false)
	explanationPlaceholder := slack.NewTextBlockObject("plain_text", "The proper term is Cylon, please let them live peacefully among you.", false, false)
	explanationElement := slack.NewPlainTextInputBlockElement(explanationPlaceholder, "explanation")
	explanation := slack.NewInputBlock("explanation", explanationText, explanationElement)

	manageTxt := slack.NewTextBlockObject("plain_text", "pick a default response type", true, false)
	threadTxt := slack.NewTextBlockObject("plain_text", "in a thread", true, false)
	channelTxt := slack.NewTextBlockObject("plain_text", "in the channel", false, false)
	ephemeralTxt := slack.NewTextBlockObject("plain_text", "in the channel, but only you see it", false, false)
	dmTxt := slack.NewTextBlockObject("plain_text", "in a direct message", false, false)

	threadOpt := slack.NewOptionBlockObject(strconv.Itoa(thread), threadTxt)
	channelOpt := slack.NewOptionBlockObject(strconv.Itoa(channel), channelTxt)
	ephemeralOpt := slack.NewOptionBlockObject(strconv.Itoa(ephemeral), ephemeralTxt)
	dmOpt := slack.NewOptionBlockObject(strconv.Itoa(directMessage), dmTxt)

	responseTypeOption := slack.NewOptionsSelectBlockElement("static_select", manageTxt, "response_type", threadOpt, channelOpt, ephemeralOpt, dmOpt)
	responseType := slack.NewInputBlock("response_type", manageTxt, responseTypeOption)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			trigger,
			explanation,
			responseType,
		},
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	modalRequest.CallbackID = "trigger_modal"

	_, err := api.OpenView(triggerID, modalRequest)
	if err != nil {
		fmt.Printf("Error opening view: %s", err)
	}
}

func handleAddTriggerModalAction(w http.ResponseWriter, r *http.Request) {
	var i slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &i)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	triggerString := i.View.State.Values["trigger_list"]["trigger_list"].Value
	explanation := i.View.State.Values["explanation"]["explanation"].Value
	responseType, err := strconv.Atoi(i.View.State.Values["response_type"]["response_type"].SelectedOption.Value)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	triggerList := strings.Split(triggerString, ",")
	for index, triggerWord := range triggerList {
		triggerList[index] = strings.TrimSpace(triggerWord)
	}

	newTrigger := trigger{
		Triggers:            triggerList,
		Explanation:         explanation,
		Creator:             i.User.Name,
		DefaultResponseType: responseType,
		Enabled:             true,
	}

	insertTrigger(i.Team.ID, newTrigger)

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

func createTriggerTable() {
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS triggers (id INTEGER PRIMARY KEY AUTOINCREMENT, teamid TEXT, trigger TEXTs, enabled INTEGER)")
	statement.Exec()
}
