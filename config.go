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
	"encoding/json"
	"os"

	log "github.com/gravitational/logrus"
	"github.com/gravitational/trace"
)

// Config is the configuration structure used everywhere in the code to pass settings and configuration
type Config struct {
	CustomerName string
	Directory    map[string]interface{}
	Slack        slackConfig
	PagerDuty    pagerDutyConfig
}

const varCustomerName = "CUSTOMER_NAME"
const varJSONDirectory = "SLACK_PAGERDUTY_DIRECTORY"
const varToken = "SLACK_TOKEN"
const varBotUsername = "SLACK_BOT_USERNAME"
const varLink = "PAGERDUTY_LINK"
const varAPIKey = "PAGERDUTY_API_KEY"
const varSchedule = "PAGERDUTY_SUPPORT_SCHEDULE"
const varService = "PAGERDUTY_SUPPORT_SERVICE"
const varFromEmail = "PAGERDUTY_FROM_EMAIL"

// FromEnv gathers configuration from the Environment variables and merge them into the Config structure
func (c *Config) FromEnv() error {
	c.CustomerName = getVarFromEnv(varCustomerName)

	JSONDirectory := getVarFromEnv(varJSONDirectory)
	err := json.Unmarshal([]byte(JSONDirectory), &c.Directory)
	if err != nil {
		trace.Wrap(err)
	}

	err = c.Slack.fromEnv()
	if err != nil {
		trace.Wrap(err)
	}

	err = c.PagerDuty.fromEnv()
	if err != nil {
		trace.Wrap(err)
	}

	return nil
}

type slackConfig struct {
	Token       string
	BotUsername string
}

// fromEnv handles the Slack part of the configuration, fetching values from Env variables
func (c *slackConfig) fromEnv() error {
	c.Token = getVarFromEnv("SLACK_TOKEN")

	c.BotUsername = getVarFromEnv("SLACK_BOT_USERNAME")

	return nil
}

// pagerDutyConfig struct is a PagerDuty config
type pagerDutyConfig struct {
	Link      string
	APIKey    string
	Schedule  string
	Service   string
	FromEmail string
}

// fromEnv handles the PagerDuty part of the configuration, fetching values from Env variables
func (c *pagerDutyConfig) fromEnv() error {
	c.Link = getVarFromEnv("PAGERDUTY_LINK")

	c.APIKey = getVarFromEnv("PAGERDUTY_API_KEY")

	c.Schedule = getVarFromEnv("PAGERDUTY_SUPPORT_SCHEDULE")

	c.Service = getVarFromEnv("PAGERDUTY_SUPPORT_SERVICE")

	c.FromEmail = getVarFromEnv("PAGERDUTY_FROM_EMAIL")

	return nil
}

// getVarFromEnv is a wrapper function that just gets variable from the Env and return an error if no value is passed
func getVarFromEnv(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatal(varName + " ENV variable must be set")
	}
	return value
}
