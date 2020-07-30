package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func interactiveHandler(w http.ResponseWriter, r *http.Request) {
	err := verifySigningSecret(r)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var interaction slack.InteractionCallback
	payload := r.FormValue("payload")

	err = json.Unmarshal([]byte(payload), &interaction)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write([]byte(""))

	log.Println(interaction.Type)
	switch {
	case interaction.Type == "view_submission" && interaction.View.CallbackID == "trigger_modal":
		handleAddTriggerModalAction(interaction)
	case interaction.Type == "view_submission" && interaction.View.CallbackID == "trigger_edit_modal":
		handleEditTriggerModalSelect(interaction, w, r)
	case interaction.Type == "view_submission" && interaction.View.CallbackID == "trigger_update_save_modal":
		handleEditTriggerModalSave(interaction)
	case interaction.Type == "block_actions" && interaction.ActionCallback.BlockActions[0].ActionID == "user_settings":
		handleUserSettingAction(interaction)
	case interaction.Type == "block_actions" && interaction.ActionCallback.BlockActions[0].BlockID == "user_setting_selection":
		handleUserSettingSelection(interaction)
	case interaction.Type == "block_actions" && interaction.ActionCallback.BlockActions[0].ActionID == "add_response":
		showTriggerModal(interaction.TriggerID)
	case interaction.Type == "block_actions" && interaction.ActionCallback.BlockActions[0].ActionID == "edit_response":
		showEditTriggerModal(interaction.TriggerID, interaction.Team.ID)
	default:
		log.Println("looking for type: " + interaction.Type)
	}

}
