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
	c.CustomerName = os.Getenv("PAGERDUTY_CUSTOMER_NAME")
	if c.CustomerName == "" {
		log.Fatal("PAGERDUTY_CUSTOMER_NAME ENV variable must be set")
	}

	JSONDirectory := os.Getenv("SLACK_PAGERDUTY_DIRECTORY")
	if JSONDirectory == "" {
		log.Fatal("SLACK_PAGERDUTY_DIRECTORY ENV variable must be set")
	}
	err := json.Unmarshal([]byte(JSONDirectory), &c.Directory)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Slack._fromEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = c.PagerDuty._fromEnv()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

type _SlackConfig struct {
	Token string
}

func (c *_SlackConfig) _fromEnv() error {
	c.Token = os.Getenv("SLACK_TOKEN")
	if c.Token == "" {
		log.Fatal("SLACK_TOKEN ENV variable must be set")
	}

	return nil
}

type _PagerDutyConfig struct {
	APIKey    string
	Schedule  string
	Service   string
	FromEmail string
}

func (c *_PagerDutyConfig) _fromEnv() error {
	c.APIKey = os.Getenv("PAGERDUTY_API_KEY")
	if c.APIKey == "" {
		log.Fatal("PAGERDUTY_API_KEY ENV variable must be set")
	}
	c.Schedule = os.Getenv("PAGERDUTY_SUPPORT_SCHEDULE")
	if c.Schedule == "" {
		log.Fatal("PAGERDUTY_SUPPORT_SCHEDULE ENV variable must be set")
	}
	c.Service = os.Getenv("PAGERDUTY_SUPPORT_SERVICE")
	if c.Service == "" {
		log.Fatal("PAGERDUTY_SUPPORT_SERVICE ENV variable must be set")
	}
	c.FromEmail = os.Getenv("PAGERDUTY_FROM_EMAIL")
	if c.FromEmail == "" {
		log.Fatal("PAGERDUTY_FROM_EMAIL ENV variable must be set")
	}

	return nil
}
