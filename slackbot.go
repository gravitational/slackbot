/*
Copyright 2018 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gravitational/trace"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/shomali11/slacker"
)

// Init is called upon Bot creation (first startup)
func Init(config *Config) {
	fmt.Printf("Connected!\n")
	client := pagerduty.NewClient(config.PagerDuty.APIKey)
	var opts pagerduty.GetScheduleOptions
	if schedule, err := client.GetSchedule(config.PagerDuty.Schedule, opts); err != nil {
		textErr := fmt.Sprintf("Error encountered while fetching schedules: %s", err.Error())
		response.Reply(textErr)
		trace.Wrap(err)
	} else {
		fmt.Printf("Configured schedule is \"%s\" with ID: %s\n", schedule.Name,
			config.PagerDuty.Schedule)
	}
}

// Err function is used to handle all Error reported by the Bot
func Err(err string) {
	fmt.Fprint(os.Stderr, err)
}

// Emergency is used to open Emergency Incidents on PagerDuty
func Emergency(request slacker.Request, response slacker.ResponseWriter, config *Config) {
	client := pagerduty.NewClient(config.PagerDuty.APIKey)
	var scheduleOpts pagerduty.GetScheduleOptions
	if schedule, err := client.GetSchedule(config.PagerDuty.Schedule, scheduleOpts); err != nil {
		textErr := fmt.Sprintf("Error encountered while fetching schedules: %s", err.Error())
		response.Reply(textErr)
		trace.Wrap(err)
	} else {
		fmt.Printf(`Opening incident on schedule "%s"/%s`, schedule.Name,
			config.PagerDuty.Schedule)
	}

	newIncident := pagerduty.CreateIncidentOptions{
		Type: "incident",
		Title: fmt.Sprintf("Incident opened by %s, via Slack/@%s",
			config.CustomerName, config.Slack.BotUsername),
	}

	newIncidentBody := pagerduty.APIDetails{
		Type:    "incident_body",
		Details: request.Param("msg"),
	}
	newIncident.Body = newIncidentBody

	newIncidentService := pagerduty.APIReference{
		Type: "service_reference",
		ID:   config.PagerDuty.Service,
	}
	newIncident.Service = newIncidentService

	createIncidentOpts := pagerduty.CreateIncident{
		Incident: newIncident,
	}

	if incident, err := client.CreateIncident(config.PagerDuty.FromEmail, &createIncidentOpts); err != nil {
		errText := "There was an error while creating a new incident created, please try again and report the following error" + err.Error()
		Err(errText)
		response.Reply(errText)
	} else {
		incidentURL := config.PagerDuty.Link + "/incidents/" + incident.Id
		fmt.Printf("Incident created by %s via @%s > %s", config.CustomerName,
			config.Slack.BotUsername, incidentURL)
		response.Reply("Incident created successfully, " +
			"please refer to incident " + incidentURL)

	}
}

// Default function handles all messages that won't match the other Commands
func Default(request slacker.Request, response slacker.ResponseWriter, config *Config) {
	client := pagerduty.NewClient(config.PagerDuty.APIKey)
	var opts pagerduty.ListOnCallUsersOptions
	opts.Since = time.Now().UTC().Format(time.RFC3339)
	opts.Until = time.Now().UTC().Add(time.Minute * 1).Format(time.RFC3339)
	if onCallUserList, err := client.ListOnCallUsers(config.PagerDuty.Schedule, opts); err != nil {
		trace.Wrap(err)
	} else {
		for _, p := range onCallUserList {
			onCallSlackUsername := config.Directory[p.Email].(string)
			responseText := fmt.Sprintf("<@%s> I think that %s may need some help ASAP! :point_up: :fire: :helmet_with_white_cross:",
				onCallSlackUsername, config.CustomerName)
			fmt.Printf("%s requested help via @%s and @%s was pinged via Slack.",
				config.CustomerName, config.Slack.BotUsername, onCallSlackUsername)
			response.Reply(responseText)
		}
	}
}

// help is used to print the Help message text
func help(resp slacker.ResponseWriter, c *Config) {
	help_text := `
> *SlackBot - HELP*
>
> _@` + c.Slack.BotUsername + ` help_ - Prints the help message (if the word help is anywhere in the sentence)
> _@` + c.Slack.BotUsername + ` open emergency ` + "`<msg>`" + `_ - Open an EMERGENCY incident to Customer Support
> _@` + c.Slack.BotUsername + ` <anything else>_ - Any other message that will be sent directly to the Bot or mentioning the 
>                                             Bot name in other channels, will result in a ping (mention) to the current 
>                                             person on call.
	`
	resp.Reply(help_text)
}
