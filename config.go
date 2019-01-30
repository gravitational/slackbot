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
	"encoding/json"
	"os"

	"github.com/gravitational/trace"
)

// config is the configuration structure used everywhere in the code to pass settings and configuration
type config struct {
	customerName string
	directory    map[string]interface{}
	slack        slackConfig
	pagerDuty    pagerDutyConfig
}

const (
	botName        = "slackbot"
	botDescription = "Slack to PagerDuty Self Support Bot"

	customerName = "CUSTOMER_NAME"

	jSONDirectory = "SLACK_PAGERDUTY_DIRECTORY"

	token = "SLACK_TOKEN"

	botUsername = "SLACK_BOT_USERNAME"

	link = "PAGERDUTY_LINK"

	aPIKey = "PAGERDUTY_API_KEY"

	schedule = "PAGERDUTY_SUPPORT_SCHEDULE"

	service = "PAGERDUTY_SUPPORT_SERVICE"

	fromEmail = "PAGERDUTY_FROM_EMAIL"
)

// FromEnv gathers configuration from the Environment variables and merge them into the Config structure
func (c *config) FromEnv() error {
	c.customerName = getVarFromEnv(customerName)

	JSONDirectory := getVarFromEnv(jSONDirectory)
	err := json.Unmarshal([]byte(JSONDirectory), &c.directory)
	if err != nil {
		return trace.Wrap(err)
	}

	err = c.slack.fromEnv()
	if err != nil {
		return trace.Wrap(err)
	}

	err = c.pagerDuty.fromEnv()
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

type slackConfig struct {
	token       string
	botUsername string
}

// fromEnv handles the Slack part of the configuration, fetching values from Env variables
func (c *slackConfig) fromEnv() error {
	c.token = getVarFromEnv(token)

	c.botUsername = getVarFromEnv(botUsername)

	return nil
}

// pagerDutyconfig struct is a pagerDuty config
type pagerDutyConfig struct {
	link      string
	aPIKey    string
	schedule  string
	service   string
	fromEmail string
}

// fromEnv handles the pagerDuty part of the configuration, fetching values from Env variables
func (c *pagerDutyConfig) fromEnv() error {
	c.link = getVarFromEnv(link)

	c.aPIKey = getVarFromEnv(aPIKey)

	c.schedule = getVarFromEnv(schedule)

	c.service = getVarFromEnv(service)

	c.fromEmail = getVarFromEnv(fromEmail)

	return nil
}

// getVarFromEnv is a wrapper function that just gets variable from the Env and return an error if no value is passed
func getVarFromEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		Err(name + " ENV variable must be set\n")
	}
	return value
}
