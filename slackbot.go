package main

import (
	"log"
	"time"

	"context"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/shomali11/slacker"
)

func main() {
	var config Config
	err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	bot := slacker.NewClient(config.Slack.Token)

	bot.Init(func() {
		log.Println("Connected!")
		client := pagerduty.NewClient(config.PagerDuty.APIKey)
		var opts pagerduty.GetScheduleOptions
		if schedule, err := client.GetSchedule(config.PagerDuty.Schedule, opts); err != nil {
			panic(err)
		} else {
			log.Println("Enable schedule is \"" + schedule.Name + "\" with ID: " + config.PagerDuty.Schedule)
		}
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.Command("emergency <msg>", "Open an EMERGENCY incident to Gravitational Customer Support",
		func(request slacker.Request, response slacker.ResponseWriter) {
			client := pagerduty.NewClient(config.PagerDuty.APIKey)
			var scheduleOpts pagerduty.GetScheduleOptions
			if schedule, err := client.GetSchedule(config.PagerDuty.Schedule, scheduleOpts); err != nil {
				log.Fatal(err)
			} else {
				log.Println("Opening incident on schedule \"" + schedule.Name +
					"\"/" + config.PagerDuty.Schedule)
			}

			// TODO: NOT WORKING
			newIncident := pagerduty.CreateIncidentOptions{
				Type:  "incident",
				Title: "Incident opened via Slack by " + config.CustomerName,
			}

			newIncidentBody := pagerduty.APIDetails{
				Type:    "incident_body",
				Details: request.Param("msg"),
			}
			log.Println(newIncidentBody.Type)

			newIncidentService := pagerduty.APIReference{
				Type: "service_reference",
				ID:   config.PagerDuty.Service,
			}
			log.Println(newIncidentService.Type)

			createIncidentOpts := pagerduty.CreateIncident{
				Incident: newIncident,
			}

			if incident, err := client.CreateIncident(config.PagerDuty.FromEmail, &createIncidentOpts); err != nil {
				log.Fatal(err)
			} else {
				log.Println("Incident created successfully" + incident.IncidentKey)
				response.Reply("Incident created successfully" + incident.IncidentKey)
			}
		})

	bot.DefaultCommand(func(request slacker.Request, response slacker.ResponseWriter) {
		client := pagerduty.NewClient(config.PagerDuty.APIKey)
		var opts pagerduty.ListOnCallUsersOptions
		opts.Since = time.Now().Format(time.RFC3339)
		opts.Until = time.Now().Add(time.Minute * 1).Format(time.RFC3339)
		if onCallUserList, err := client.ListOnCallUsers(config.PagerDuty.Schedule, opts); err != nil {
			panic(err)
		} else {
			for _, p := range onCallUserList {
				response_text := (" :point_up: :fire: :helmet_with_white_cross: " +
					"@" + config.Directory[p.Email].(string) +
					" has been summoned and will soon be here to help")
				log.Println(response_text)
				response.Reply(response_text)
			}
		}
	})

	bot.DefaultEvent(func(event interface{}) {
		//log.Println(event)
	})

	bot.Help("help", slacker.WithHandler(func(request slacker.Request, response slacker.ResponseWriter) {
		help(response)
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func help(resp slacker.ResponseWriter) {
	help_text := `
	"* Gravity SlackBot - HELP *"

	emergency <msg> - Open an EMERGENCY incident to Gravitational Customer Support

	Any other message that will be sent directly to the Bot or mentioning the 
	Bot name in other channels, will result in a ping (mention) to the current 
	person on call.
	`
	resp.Reply(help_text)
}
