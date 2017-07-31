package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Names map[string]int

type ReminderLog map[string]Reminder

type Reminder struct {
	Name      string
	Completed bool
	Time      int64
}

type Config struct {
	SlackAPIKey string
}

func main() {
	var config Config

	configBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Could not read config:", err)
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal("Could not read config:", err)
	}

	api := slack.New(config.SlackAPIKey)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			// Replace #general with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			if strings.Contains(strings.ToLower(ev.Text), "called") {
				parts := strings.Split(ev.Text, " ")
				if len(parts) > 2 {
					rtm.SendMessage(rtm.NewOutgoingMessage("Please mark reminders in the form of `called [name]`", ev.Channel))
					return
				}
				//markCalled(parts[1])
			}

		case *slack.RTMError:
			fmt.Printf("Slack Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		}
	}

	//names, err := readNames()

	//if err != nil {
	//	log.Fatal("Error reading names:", err)
	//}

	//reminderLog, err := readReminders()

	//if err != nil {
	//	log.Fatal("Error reading reminders:", err)
	//}

	//remind(names, &reminderLog)

}

func remind(names Names, reminders ReminderLog) {
	now := time.Now().Unix()
	for name, interval := range names {
		if r, ok := reminders[name]; ok {
			if now-r.Time > int64(interval*86400) {
				if r.Completed {
					r.Completed = false
				} else {
					r.Time = now
					//notify(name, interval)
				}
			}
		}
	}
}

func readNames() (Names, error) {
	var names Names

	// Get list of names
	namesBytes, err := ioutil.ReadFile("./names.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(namesBytes, &names)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func readReminders() (ReminderLog, error) {
	var rlog ReminderLog

	// Read reminder log
	reminderBytes, err := ioutil.ReadFile("./reminder-log.json")
	if err != nil {
		return rlog, err
	}

	if len(reminderBytes) != 0 {
		err = json.Unmarshal(reminderBytes, &rlog)
		if err != nil {
			return rlog, err
		}
	}

	return rlog, nil
}
