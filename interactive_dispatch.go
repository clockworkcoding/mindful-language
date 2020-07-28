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

	switch interaction.View.CallbackID {
	case "trigger_modal":
		handleAddTriggerModalAction(w, r)
	default:
		log.Println("looking for callback: " + interaction.CallbackID)
	}

}
