package internal

import "github.com/go-resty/resty/v2"

// type aliases for better readability
type JiraResponse map[string]interface{}
type WorklogData map[string]interface{}

// tempoo client struct
type Tempoo struct {
	email    string
	apiToken string
	client   *resty.Client // resty client for making HTTP requests to the Jira API
}
