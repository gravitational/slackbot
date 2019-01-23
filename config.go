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
	Slack        _SlackConfig
	PagerDuty    _PagerDutyConfig
}

// FromEnv TODO
func (c *Config) FromEnv() error {
	c.CustomerName = _GetVarFromEnv("CUSTOMER_NAME")

	JSONDirectory := _GetVarFromEnv("SLACK_PAGERDUTY_DIRECTORY")
	err := json.Unmarshal([]byte(JSONDirectory), &c.Directory)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Slack._FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = c.PagerDuty._FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

type _SlackConfig struct {
	Token       string
	BotUsername string
}

func (c *_SlackConfig) _FromEnv() error {
	c.Token = _GetVarFromEnv("SLACK_TOKEN")

	c.BotUsername = _GetVarFromEnv("SLACK_BOT_USERNAME")

	return nil
}

type _PagerDutyConfig struct {
	Link      string
	APIKey    string
	Schedule  string
	Service   string
	FromEmail string
}

func (c *_PagerDutyConfig) _FromEnv() error {
	c.Link = _GetVarFromEnv("PAGERDUTY_LINK")

	c.APIKey = _GetVarFromEnv("PAGERDUTY_API_KEY")

	c.Schedule = _GetVarFromEnv("PAGERDUTY_SUPPORT_SCHEDULE")

	c.Service = _GetVarFromEnv("PAGERDUTY_SUPPORT_SERVICE")

	c.FromEmail = _GetVarFromEnv("PAGERDUTY_FROM_EMAIL")

	return nil
}

func _GetVarFromEnv(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		log.Fatal(varName + " ENV variable must be set")
	}
	return value
}
