package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tkanos/gonfig"

	"github.com/slack-go/slack"
)

var (
	config Configuration
	api    *slack.Client
)

type state struct {
	auth string
	ts   time.Time
}

// globalState is an example of how to store a state between calls
var globalState state

// writeError writes an error to the reply
func writeError(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte("Something went wrong, please try again or contact Max@ClockworkCoding.com if the problem persists."))
	log.Output(1, fmt.Sprintf("Err: %s", err))
}

// func responseError(responseURL, message, token string) {
// 	log.Output(1, fmt.Sprintf("Err: %s", message))
// 	simpleResponse(responseURL, "Something went wrong, please try again or contact Max@ClockworkCoding.com if the problem persists.", false, token)
// }

// func simpleResponse(responseURL, message string, replace bool, token string) {
// 	params := slack.mes .NewResponseMessageParameters()
// 	params.ResponseType = "ephemeral"
// 	params.ReplaceOriginal = replace
// 	params.Text = message

// 	api := slack.New(token)
// 	err := api.PostResponse(responseURL, params)
// 	if err != nil {
// 		log.Output(0, fmt.Sprintf("Err: %s", err.Error()))
// 	}

// }

//Configuration config struct
type Configuration struct {
	SlackClientID          string
	SlackClientSecret      string
	SlackVerificationToken string
	SlackSigningSecret     string
	BotToken               string
	Key1                   string
	Key2                   string
	URL                    string
	PORT                   string
	RedirectURL            string
	Patreon                string
	DbUser                 string
	DbPass                 string
	DbPath                 string
}

func main() {
	routing()
}

func routing() {

	mux := http.NewServeMux()

	mux.Handle("/event", http.HandlerFunc(eventHandler))
	mux.Handle("/interactive", http.HandlerFunc(interactiveHandler))
	mux.Handle("/addtrigger", http.HandlerFunc(addTriggerCommand))
	mux.Handle("/edittrigger", http.HandlerFunc(editTriggerCommand))
	mux.Handle("/", http.HandlerFunc(redirect))
	err := http.ListenAndServe(":"+config.PORT, mux)
	log.Println("Running on port: " + config.PORT)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}

func buttonPressed(w http.ResponseWriter, r *http.Request) {
	log.Println("button")

}

func redirect(w http.ResponseWriter, r *http.Request) {
	if url := os.Getenv(strings.Replace(r.URL.Path, "/", "URL_", 1)); url != "" {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, config.RedirectURL+r.URL.Path, http.StatusTemporaryRedirect)
	}
}

func init() {
	settingsFile := "settings.json"
	_, err := os.Stat("usersettings.json")
	if err == nil {
		settingsFile = "usersettings.json"
	}
	err = gonfig.GetConf(settingsFile, &config)
	if err != nil {
		log.Fatalln("Could not load configuration")
	}

	api = slack.New(config.BotToken)
}
