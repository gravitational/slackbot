package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"context"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/shomali11/slacker"
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		panic("SLACK_TOKEN ENV variable must be set")
	}
	PDApiKey := os.Getenv("PAGERDUTY_API_KEY")
	if PDApiKey == "" {
		panic("PAGERDUTY_API_KEY ENV variable must be set")
	}
	PDSupportSchedule := os.Getenv("PAGERDUTY_SUPPORT_SCHEDULE")
	if PDSupportSchedule == "" {
		panic("PAGERDUTY_SUPPORT_SCHEDULE ENV variable must be set")
	}
	PDSupportService := os.Getenv("PAGERDUTY_SUPPORT_SERVICE")
	if PDSupportSchedule == "" {
		panic("PAGERDUTY_SUPPORT_SERVICE ENV variable must be set")
	}
	PDCustomerName := os.Getenv("PAGERDUTY_CUSTOMER_NAME")
	if PDCustomerName == "" {
		panic("PAGERDUTY_CUSTOMER_NAME ENV variable must be set")
	}

	slackPDDirectoryJSON := os.Getenv("SLACK_PAGERDUTY_DIRECTORY")
	if slackPDDirectoryJSON == "" {
		panic("SLACK_PAGERDUTY_DIRECTORY ENV variable must be set")
	}
	var slackPDDirectory map[string]interface{}
	err := json.Unmarshal([]byte(slackPDDirectoryJSON), &slackPDDirectory)
	if err != nil {
		panic(err)
	}
	PDFromEmail := os.Getenv("PAGERDUTY_FROM_EMAIL")
	if slackPDDirectoryJSON == "" {
		panic("PAGERDUTY_FROM_EMAIL ENV variable must be set")
	}
	//PDServiceApiKey	:= os.Getenv("PAGERDUTY_SERVICE_API_KEY")

	bot := slacker.NewClient(slackToken)

	bot.Init(func() {
		log.Println("Connected!")
		client := pagerduty.NewClient(PDApiKey)
		var opts pagerduty.GetScheduleOptions
		if schedule, err := client.GetSchedule(PDSupportSchedule, opts); err != nil {
			panic(err)
		} else {
			log.Println("Enable schedule is \"" + schedule.Name +
				"\" with ID: " + PDSupportSchedule)
		}
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.Command("emergency <msg>", "Open an EMERGENCY incident to Gravitational Customer Support",
		func(request slacker.Request, response slacker.ResponseWriter) {
			client := pagerduty.NewClient(PDApiKey)
			var scheduleOpts pagerduty.GetScheduleOptions
			if schedule, err := client.GetSchedule(PDSupportSchedule, scheduleOpts); err != nil {
				panic(err)
			} else {
				log.Println("Opening incident on schedule \"" + schedule.Name +
					"\"/" + PDSupportSchedule)
			}

			// TODO: NOT WORKING
			var newIncident pagerduty.CreateIncidentOptions
			newIncident["Type"] = "incident"
			newIncident["Title"] = "Incident opened via Slack by " + PDCustomerName

			var newIncidentBody pagerduty.APIDetails
			newIncidentBody["Type"] = "incident_body"
			newIncidentBody["Details"] = request.Param("msg")

			newIncidentService := pagerduty.APIReference{
				"Type": "service_reference",
				"ID":   PDSupportService,
			}

			var createIncidentOpts *pagerduty.CreateIncident
			createIncidentOpts["Incident"] = newIncident

			if incident, err := client.CreateIncident(PDFromEmail, createIncidentOpts); err != nil {
				panic(err)
			} else {
				log.Println("Incident created successfully" + incident.IncidentKey)
				response.Reply("Incident created successfully" + incident.IncidentKey)
			}
		})

	bot.DefaultCommand(func(request slacker.Request, response slacker.ResponseWriter) {
		client := pagerduty.NewClient(PDApiKey)
		var opts pagerduty.ListOnCallUsersOptions
		opts.Since = time.Now().Format(time.RFC3339)
		opts.Until = time.Now().Add(time.Minute * 1).Format(time.RFC3339)
		if onCallUserList, err := client.ListOnCallUsers(PDSupportSchedule, opts); err != nil {
			panic(err)
		} else {
			for _, p := range onCallUserList {
				response_text := (" :point_up: :fire: :helmet_with_white_cross: " +
					"@" + slackPDDirectory[p.Email].(string) +
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
