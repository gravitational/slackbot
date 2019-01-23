package main

import (
	"log"

	"context"

	"github.com/shomali11/slacker"
)

func main() {
	var config Config
	err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	bot := slacker.NewClient(config.Slack.Token)
	// defining which function handles the bot Init phase
	bot.Init(
		func(){
			slackbot.Init()
		}
	)

	// error raised by the Bot will be handled by this function
	bot.Err(
		func(err string){
			slackbot.Err(err)
		}
	)

	// function tied to sentences sent to the Bot and starting with "open emergency" followed by some text
	bot.Command("open emergency <msg>",
		"Open an EMERGENCY incident to Gravitational Customer Support",
		func(request slacker.Request, response slacker.ResponseWriter) {
			slackbot.Emergency(request, response)
		}
	)

	// when no other "Command" matches and text is sent to the Bot, this function will be run instead
	bot.DefaultCommand(
		func(request slacker.Request, response slacker.ResponseWriter) {
			slackbot.Default(request, response)
		}
	)

	// function run for all events received by the bot (including time ticks)
	bot.DefaultEvent(
		func(event interface{}) {
			//log.Println(event)
		}
	)

	// set the "help" message handling function
	bot.Help("help", 
		slacker.WithHandler(
			func(request slacker.Request, response slacker.ResponseWriter) {
				slackbot.help(response, &config)
			}
		)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
