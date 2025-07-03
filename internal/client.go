package internal

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/go-resty/resty/v2"
)

// NewTempoo creates a new client for the Jira API
func NewTempoo() (*Tempoo, error) {
	email := os.Getenv("JIRA_EMAIL")
	if email == "" {
		return nil, &TempooError{Message: "JIRA_EMAIL environment variable is not set"}
	}
	log.Debugf("Read JIRA_EMAIL from env: %s", email)

	apiToken := os.Getenv("JIRA_API_TOKEN")
	if apiToken == "" {
		return nil, &TempooError{Message: "JIRA_API_TOKEN environment variable is not set"}
	}
	log.Debug("Read JIRA_API_TOKEN from env")

	// create a new resty client
	client := resty.New()
	// set auth
	client.SetBasicAuth(email, apiToken)
	// build header
	client.SetHeader("Content-Type", "application/json")
	// set default timeout
	client.SetTimeout(10 * time.Second)

	log.Debugf("Created Resty client: %+v", client)

	t := &Tempoo{
		email:    email,
		apiToken: apiToken,
		client:   client,
	}

	log.Debug("Tempoo initialized")
	return t, nil
}
