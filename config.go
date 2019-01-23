package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config TODO
type Config struct {
	CustomerName string
	Directory    map[string]interface{}
	Slack        slackConfig
	PagerDuty    pagerDutyConfig
}

// FromEnv TODO
func (c *Config) FromEnv() error {
	c.CustomerName = getVarFromEnv("CUSTOMER_NAME")

	JSONDirectory := getVarFromEnv("SLACK_PAGERDUTY_DIRECTORY")
	err := json.Unmarshal([]byte(JSONDirectory), &c.Directory)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Slack.fromEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = c.PagerDuty.fromEnv()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

type slackConfig struct {
	Token       string
	BotUsername string
}

func (c *slackConfig) fromEnv() error {
	c.Token = getVarFromEnv("SLACK_TOKEN")

	c.BotUsername = getVarFromEnv("SLACK_BOT_USERNAME")

	return nil
}

type pagerDutyConfig struct {
	Link      string
	APIKey    string
	Schedule  string
	Service   string
	FromEmail string
}

func (c *pagerDutyConfig) fromEnv() error {
	c.Link = getVarFromEnv("PAGERDUTY_LINK")

	c.APIKey = getVarFromEnv("PAGERDUTY_API_KEY")

	c.Schedule = getVarFromEnv("PAGERDUTY_SUPPORT_SCHEDULE")

	c.Service = getVarFromEnv("PAGERDUTY_SUPPORT_SERVICE")

	c.FromEmail = getVarFromEnv("PAGERDUTY_FROM_EMAIL")

	return nil
}

func getVarFromEnv(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatal(varName + " ENV variable must be set")
	}
	return value
}
