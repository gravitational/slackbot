/*
Copyright 2019 Gravitational, Inc.

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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gravitational/trace"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/shomali11/slacker"
)

// Start will run the main code for the Slack->PD Bot
func Start(config *config) error {
	bot := slacker.NewClient(config.slack.token)
	// defining which function handles the bot Init phase
	bot.Init(
		func() {
			Init(config)
		})

	// error raised by the Bot will be handled by this function
	bot.Err(
		func(err string) {
			Err(err)
		})

	// function tied to sentences sent to the Bot and starting with "open emergency" followed by some text
	emergencyCmdDefinition := &slacker.CommandDefinition{
		Description: "Open an EMERGENCY incident to Customer Support",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			Emergency(request, response, config)
		},
	}
	bot.Command("open emergency <msg>", emergencyCmdDefinition)

	// when no other "Command" matches and text is sent to the Bot, this function will be run instead
	bot.DefaultCommand(
		func(request slacker.Request, response slacker.ResponseWriter) {
			Default(request, response, config)
		})

	// function run for all events received by the bot (including time ticks)
	bot.DefaultEvent(
		func(event interface{}) {
			//log.Println(event)
		})

	// set the "help" message handling function
	helpCmdDefinition := &slacker.CommandDefinition{
		Description: "Help function",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			help(response, config)
		},
	}
	bot.Help(helpCmdDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		Err(err.Error())
	}

	return nil
}

// Init is called upon Bot creation (first startup)
func Init(config *config) {
	fmt.Printf("Connected!\n")
	client := pagerduty.NewClient(config.pagerDuty.aPIKey)
	var opts pagerduty.GetScheduleOptions

	schedule, err := client.GetSchedule(config.pagerDuty.schedule, opts)
	if err != nil {
		trace.Wrap(err)
	}
	fmt.Printf("Configured schedule is \"%s\" with ID: %s\n", schedule.Name,
		config.pagerDuty.schedule)
}

// Err function is used to handle all Error reported by the Bot
func Err(err string) {
	fmt.Fprint(os.Stderr, err)
}

// Emergency is used to open Emergency Incidents on PagerDuty
func Emergency(request slacker.Request, response slacker.ResponseWriter, config *config) {
	client := pagerduty.NewClient(config.pagerDuty.aPIKey)
	var scheduleOpts pagerduty.GetScheduleOptions

	schedule, err := client.GetSchedule(config.pagerDuty.schedule, scheduleOpts)
	if err != nil {
		textErr := fmt.Sprintf("Error encountered while fetching schedules: %s", err.Error())
		response.Reply(textErr)
		trace.Wrap(err)
	}
	fmt.Printf(`Opening incident on schedule "%s"/%s`, schedule.Name,
		config.pagerDuty.schedule)

	newIncident := pagerduty.CreateIncidentOptions{
		Type: "incident",
		Title: fmt.Sprintf("Incident opened by %s, via Slack/@%s",
			config.customerName, config.slack.botUsername),
	}

	newIncident.Body = pagerduty.APIDetails{
		Type:    "incident_body",
		Details: request.Param("msg"),
	}

	newIncident.Service = pagerduty.APIReference{
		Type: "service_reference",
		ID:   config.pagerDuty.service,
	}

	createIncidentOpts := pagerduty.CreateIncident{
		Incident: newIncident,
	}

	incident, err := client.CreateIncident(config.pagerDuty.fromEmail, &createIncidentOpts)
	if err != nil {
		errText := "There was an error while creating a new incident created, please try again and report the following error" + err.Error()
		Err(errText)
		response.Reply(errText)
		return
	}

	incidentURL := config.pagerDuty.link + "/incidents/" + incident.Id
	fmt.Printf("Incident created by %s via @%s > %s\n", config.customerName,
		config.slack.botUsername, incidentURL)
	response.Reply("Incident created successfully, please refer to incident " + incidentURL)
}

// Default function handles all messages that won't match the other Commands
func Default(request slacker.Request, response slacker.ResponseWriter, config *config) {
	client := pagerduty.NewClient(config.pagerDuty.aPIKey)
	var opts pagerduty.ListOnCallUsersOptions

	now := time.Now().UTC()
	opts.Since = now.Format(time.RFC3339)
	opts.Until = now.Add(time.Minute * 1).Format(time.RFC3339)

	onCallUserList, err := client.ListOnCallUsers(config.pagerDuty.schedule, opts)
	if err != nil {
		errText := "There was an error while fetching oncall users, please try again and report the following error" + err.Error()
		response.Reply(errText)
		trace.Wrap(err)
		return
	}

	for _, p := range onCallUserList {
		if config.directory[p.Email] == nil {
			fmt.Printf("Oncall %s user not found. Please report this error\n", p.Email)
			textErr := fmt.Sprintf("Oncall user not found. Please report this error")
			response.Reply(textErr)
			continue
		}

		onCallSlackUsername := config.directory[p.Email].(string)
		responseText := fmt.Sprintf("<@%s> I think that %s may need some help ASAP! :point_up: :fire: :helmet_with_white_cross:",
			onCallSlackUsername, config.customerName)
		fmt.Printf("%s requested help via @%s and @%s was pinged via Slack.\n",
			config.customerName, config.slack.botUsername, onCallSlackUsername)
		response.Reply(responseText)
	}
}

// help is used to print the Help message text
func help(resp slacker.ResponseWriter, c *config) {
	help_text := `
> *SlackBot - HELP*
>
> _@` + c.slack.botUsername + ` help_ - Prints the help message (if the word help is anywhere in the sentence)
> _@` + c.slack.botUsername + ` open emergency ` + "`<msg>`" + `_ - Open an EMERGENCY incident to Customer Support
> _@` + c.slack.botUsername + ` <anything else>_ - Any other message that will be sent directly to the Bot or mentioning the 
>                                             Bot name in other channels, will result in a ping (mention) to the current 
>                                             person on call.
	`
	resp.Reply(help_text)
}
