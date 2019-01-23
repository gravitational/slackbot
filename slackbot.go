package main

import (
	"log"
	"time"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/shomali11/slacker"
)

// Init TODO
func Init() {
	log.Println("Connected!")
	client := pagerduty.NewClient(config.PagerDuty.APIKey)
	var opts pagerduty.GetScheduleOptions
	if schedule, err := client.GetSchedule(config.PagerDuty.Schedule, opts); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Configured schedule is \"" + schedule.Name + "\" with ID: " + config.PagerDuty.Schedule)
	}
}

// Err TODO
func Err(err string) {
	log.Println(err)
}

// Emergency TODO
func Emergency(request slacker.Request, response slacker.ResponseWriter) {
	client := pagerduty.NewClient(config.PagerDuty.APIKey)
	var scheduleOpts pagerduty.GetScheduleOptions
	if schedule, err := client.GetSchedule(config.PagerDuty.Schedule, scheduleOpts); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Opening incident on schedule \"" + schedule.Name +
			"\"/" + config.PagerDuty.Schedule)
	}

	newIncident := pagerduty.CreateIncidentOptions{
		Type: "incident",
		Title: "Incident opened by " + config.CustomerName +
			" via Slack/@" + config.Slack.BotUsername,
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
		log.Println(errText)
		response.Reply(errText)
	} else {
		log.Println("Incident created by " + config.CustomerName +
			" via @" + config.Slack.BotUsername +
			" > https://gravitational.pagerduty.com/incidents/" + incident.Id)
		response.Reply("Incident created successfully, " +
			"please refer to incident " +
			"https://gravitational.pagerduty.com/incidents/" + incident.Id)
	}
}

// Default TODO
func Default(request slacker.Request, response slacker.ResponseWriter) {
	client := pagerduty.NewClient(config.PagerDuty.APIKey)
	var opts pagerduty.ListOnCallUsersOptions
	opts.Since = time.Now().Format(time.RFC3339)
	opts.Until = time.Now().Add(time.Minute * 1).Format(time.RFC3339)
	if onCallUserList, err := client.ListOnCallUsers(config.PagerDuty.Schedule, opts); err != nil {
		log.Fatal(err)
	} else {
		for _, p := range onCallUserList {
			onCallSlackUsername := config.Directory[p.Email].(string)
			responseText := "@" + onCallSlackUsername + " I think that " +
				config.CustomerName + " may need some help ASAP! " +
				":point_up: :fire: :helmet_with_white_cross:"
			log.Println(config.CustomerName + " requested help via @" +
				config.Slack.BotUsername + " and @" + onCallSlackUsername +
				" was pinged via Slack.")
			response.Reply(responseText)
		}
	}
}

func help(resp slacker.ResponseWriter, c *Config) {
	help_text := `
> *Gravity SlackBot - HELP*
>
> _@` + c.Slack.BotUsername + ` help_ - Prints the help message (if the word help is anywhere in the sentence)
> _@` + c.Slack.BotUsername + ` open emergency ` + "`<msg>`" + `_ - Open an EMERGENCY incident to Gravitational Customer Support
> _@` + c.Slack.BotUsername + ` <anything else>_ - Any other message that will be sent directly to the Bot or mentioning the 
>                                             Bot name in other channels, will result in a ping (mention) to the current 
>                                             person on call.
	`
	resp.Reply(help_text)
}
