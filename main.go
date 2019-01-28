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
	"context"
	"fmt"

	"github.com/shomali11/slacker"
)

func main() {
	var config Config
	err := config.FromEnv()
	if err != nil {
		fmt.Println(err.Error())
	}

	bot := slacker.NewClient(config.Slack.Token)
	// defining which function handles the bot Init phase
	bot.Init(
		func() {
			Init(&config)
		})

	// error raised by the Bot will be handled by this function
	bot.Err(
		func(err string) {
			Err(err)
		})

	// function tied to sentences sent to the Bot and starting with "open emergency" followed by some text
	bot.Command("open emergency <msg>",
		"Open an EMERGENCY incident to Customer Support",
		func(request slacker.Request, response slacker.ResponseWriter) {
			Emergency(request, response, &config)
		})

	// when no other "Command" matches and text is sent to the Bot, this function will be run instead
	bot.DefaultCommand(
		func(request slacker.Request, response slacker.ResponseWriter) {
			Default(request, response, &config)
		})

	// function run for all events received by the bot (including time ticks)
	bot.DefaultEvent(
		func(event interface{}) {
			//log.Println(event)
		})

	// set the "help" message handling function
	bot.Help("help",
		slacker.WithHandler(
			func(request slacker.Request, response slacker.ResponseWriter) {
				help(response, &config)
			}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
}
