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

	slackPDDirectoryJSON := os.Getenv("SLACK_PAGERDUTY_DIRECTORY")
	if slackPDDirectoryJSON == "" {
		panic("SLACK_PAGERDUTY_DIRECTORY ENV variable must be set")
	}
	var slackPDDirectory map[string]interface{}
	err := json.Unmarshal([]byte(slackPDDirectoryJSON), &slackPDDirectory)
	if err != nil {
		panic(err)
	}
	//PDServiceApiKey	:= os.Getenv("PAGERDUTY_SERVICE_API_KEY")
	//PDFromEmail		:= os.Getenv("PAGERDUTY_FROM_EMAIL")

	bot := slacker.NewClient(slackToken)
	var opts pagerduty.ListEscalationPoliciesOptions
	bot.Init(func() {
		log.Println("Connected!")
		client := pagerduty.NewClient(PDApiKey)
		if eps, err := client.ListEscalationPolicies(opts); err != nil {
			panic(err)
		} else {
			for _, p := range eps.EscalationPolicies {
				log.Println(p.Name)
			}
		}
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.Command("emergency <msg>", "Open an EMERGENCY incident to Gravitational Customer Support",
		func(request slacker.Request, response slacker.ResponseWriter) {
			msg := request.Param("msg")
			response.Reply(msg)
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
		log.Println(event)
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
