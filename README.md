# slackbot - Slack to PagerDuty integration Bot for Customer support

## Intro

`slackbot` is an automated Bot for Slack that connects to PagerDuty to provide some basic self-service Support features to Customers.

## Settings

In order to do so a few basic requirements needs to be setup and configured via environment variables:
| Variable Name | Description | Example |
|---------------|-------------|---------|
| CUSTOMER_NAME | Arbitrary name used in messages and needed to identify the customer requesting support | "Awesome Customer Brand" |
| SLACK_PAGERDUTY_DIRECTORY | String that specifies a Directory of email addresses (used in PagerDuty) and their corresponding Slack user (via their Slack username) in a JSON format | { "oncalluser1@example.email" : "slack_username1", "oncalluser2@example.email" : "slack_username2" } |
| SLACK_TOKEN | Token provided by Slack when creating a new Bot app. | 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef |
| SLACK_BOT_USERNAME | Slack username assigned to the bot, mostly used in the `help` message output | "customersupportbot" |
| PAGERDUTY_LINK | Base PagerDuty link used to hook up to the API | "https://mycompany.pagerduty.com" |
| PAGERDUTY_API_KEY | PagerDuty API Key provided when creating a new API Token | 0123456789abcdef0123456789abcdef |
| PAGERDUTY_FROM_EMAIL | Email Address used as the "From:" in Incidents creation on PagerDuty API requests | "customersupport@example.email" |
| PAGERDUTY_SERVICE_API_KEY | TODO | 0123456789abcdef0123456789abcdef |
| PAGERDUTY_SUPPORT_SCHEDULE | ID associated with the PagerDuty Schedule used to provide oncall support the specified customer | P000001 |
| PAGERDUTY_SUPPORT_SERVICE | TODO | P000001 |

## Running slackbot Locally

You can test your `slackbot` by running the following, however all [environment variables](#Settings)
must be set correctly before it will work.

You can build and run the bot locally by running:

    $ go build
    $ ./slackbot

Setting all environment variables in your current shell is tedious so we advise 
to use a script with all the variables set locally.

## Running slackbot in Docker (and with Docker Compose) 

We suggest using Docker Compose for your local testing and development.<br>
Using `slackbot` with Docker Compose is as easy as setting all your variables in
a `.env` file and then running:

     docker-compose build 
     docker-compose up

You can choose a different name for the `docker-compose.yml` and `.env` files 
and specify them through the appropriate settings.

## Running slackbot in Kubernetes

### Through a Deployment

Once you have your Bot tested and working, we suggest you running it in 
Kubernetes in production.
In order to do so you can use the Deployment we offer as an example:

    kubectl apply -f slackbot_deploy.yml


## TODO

* write a simple `helm` chart for alternative Kubernetes deployment
