package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	// JiraFQDN is the FQDN of the Jira instance
	JiraFQDN = "esendex.atlassian.net"
	// JiraAPIRootURL is the root URL of the Jira API
	JiraAPIRootURL = "https://" + JiraFQDN + "/rest/api/3"
)

// Custom error types
type TempooError struct {
	Message string
	Cause   error
}

// Error returns the error message
func (e *TempooError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// InvalidIssueKeyError is an error type for invalid issue keys
type InvalidIssueKeyError struct {
	IssueKey string
}

// Error returns the error message
func (e *InvalidIssueKeyError) Error() string {
	return fmt.Sprintf("Issue key %s is not valid", e.IssueKey)
}

// Type aliases for better readability
type JiraResponse map[string]interface{}
type WorklogData map[string]interface{}

// Tempoo client struct
type Tempoo struct {
	email    string
	apiToken string
	client   *resty.Client
}

// NewTempoo creates a new Tempoo client
func NewTempoo() (*Tempoo, error) {
	email := os.Getenv("JIRA_EMAIL")
	if email == "" {
		return nil, &TempooError{Message: "JIRA_EMAIL environment variable is not set"}
	}
	logrus.Debugf("Read JIRA_EMAIL from env: %s", email)

	apiToken := os.Getenv("JIRA_API_TOKEN")
	if apiToken == "" {
		return nil, &TempooError{Message: "JIRA_API_TOKEN environment variable is not set"}
	}
	logrus.Debug("Read JIRA_API_TOKEN from env")

	client := resty.New()
	client.SetBasicAuth(email, apiToken)
	client.SetHeader("Content-Type", "application/json")
	client.SetTimeout(10 * time.Second)
	logrus.Debugf("Created Resty client: %+v", client)

	t := &Tempoo{
		email:    email,
		apiToken: apiToken,
		client:   client,
	}

	logrus.Debug("Tempoo initialized")
	return t, nil
}

// validateIssueKey validates that an issue key exists
func (t *Tempoo) validateIssueKey(issueKey string) error {
	issueURL := fmt.Sprintf("%s/issue/%s", JiraAPIRootURL, issueKey)

	resp, err := t.client.R().Get(issueURL)
	if err != nil {
		logrus.Errorf("Request failed: %v", err)
		return &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() != 200 {
		return &InvalidIssueKeyError{IssueKey: issueKey}
	}

	return nil
}

// GetUserAccountID gets the current user's Atlassian account ID
func (t *Tempoo) GetUserAccountID() (string, error) {
	logrus.Info("Getting current user Atlassian account ID...")

	resp, err := t.client.R().Get(fmt.Sprintf("%s/myself", JiraAPIRootURL))
	if err != nil {
		logrus.Errorf("Request failed: %v", err)
		return "", &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() != 200 {
		return "", &TempooError{Message: fmt.Sprintf("Failed to get user info: %s", resp.Status())}
	}

	var userData JiraResponse
	if err := json.Unmarshal(resp.Body(), &userData); err != nil {
		return "", &TempooError{Message: "Failed to parse user data", Cause: err}
	}

	accountID, ok := userData["accountId"].(string)
	if !ok || accountID == "" {
		return "", &TempooError{Message: "Account ID not found in user data"}
	}

	logrus.Infof("Current user Atlassian account ID: %s", accountID)
	return accountID, nil
}

// GetWorklogs gets worklog IDs for a given issue key and user
func (t *Tempoo) GetWorklogs(issueKey, userID string) ([]string, error) {
	logrus.Infof("Getting worklogs for %s", issueKey)

	if err := t.validateIssueKey(issueKey); err != nil {
		return nil, err
	}

	resp, err := t.client.R().Get(fmt.Sprintf("%s/issue/%s/worklog", JiraAPIRootURL, issueKey))
	if err != nil {
		logrus.Errorf("Request failed: %v", err)
		return nil, &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() != 200 {
		return nil, &TempooError{Message: fmt.Sprintf("Failed to get worklogs: %s", resp.Status())}
	}

	var responseData JiraResponse
	if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
		return nil, &TempooError{Message: "Failed to parse worklog data", Cause: err}
	}

	worklogsInterface, ok := responseData["worklogs"]
	if !ok {
		return []string{}, nil
	}

	worklogs, ok := worklogsInterface.([]interface{})
	if !ok {
		return []string{}, nil
	}

	var worklogsForUser []string
	for _, worklogInterface := range worklogs {
		worklog, ok := worklogInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if this worklog belongs to the user
		author, ok := worklog["author"].(map[string]interface{})
		if !ok {
			continue
		}

		authorAccountID, ok := author["accountId"].(string)
		if !ok || authorAccountID != userID {
			continue
		}

		// Get the worklog ID
		worklogID, ok := worklog["id"]
		if !ok {
			continue
		}

		// Convert to string (it might be a number)
		worklogIDStr := fmt.Sprintf("%v", worklogID)
		worklogsForUser = append(worklogsForUser, worklogIDStr)
	}

	logrus.Debugf("Found %d worklogs for user %s", len(worklogsForUser), userID)
	return worklogsForUser, nil
}

// DeleteWorklog deletes a worklog for a given issue key
func (t *Tempoo) DeleteWorklog(issueKey, worklogID string) error {
	logrus.Debugf("Deleting worklog %s for %s", worklogID, issueKey)

	resp, err := t.client.R().Delete(fmt.Sprintf("%s/issue/%s/worklog/%s", JiraAPIRootURL, issueKey, worklogID))
	if err != nil {
		logrus.Errorf("Request failed: %v", err)
		return &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() != 204 {
		return &TempooError{Message: fmt.Sprintf("Failed to delete worklog: %s", resp.Status())}
	}

	logrus.Infof("Deleted worklog %s for %s", worklogID, issueKey)
	return nil
}

// AddWorklog adds a worklog to a Jira issue
func (t *Tempoo) AddWorklog(issueKey, worklogTime string, dateStr *string) error {
	logrus.Infof("Adding worklog to %s", issueKey)

	// Use current date if none provided
	var workDate time.Time
	if dateStr == nil || *dateStr == "" {
		workDate = time.Now()
	} else {
		parsedDate, err := parseDateString(*dateStr)
		if err != nil {
			return err
		}
		workDate = parsedDate
	}

	// Create ISO timestamp for 08:30 AM on the specified date
	started := time.Date(
		workDate.Year(), workDate.Month(), workDate.Day(),
		8, 30, 0, 751000000, // 08:30:00.751
		time.UTC,
	).Format("2006-01-02T15:04:05.000-0700")

	logrus.Debugf("Started timestamp: %s", started)

	payload := map[string]string{
		"timeSpent": worklogTime,
		"started":   started,
	}

	logrus.Debugf("Payload: %+v", payload)

	resp, err := t.client.R().
		SetBody(payload).
		Post(fmt.Sprintf("%s/issue/%s/worklog", JiraAPIRootURL, issueKey))

	if err != nil {
		logrus.Errorf("Request failed: %v", err)
		return &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() == 201 {
		logrus.Infof("Worklog added to %s", issueKey)
		return nil
	}

	return &TempooError{Message: fmt.Sprintf("Failed to add worklog: %s", resp.Status())}
}

// parseDateString parses date string in DD.MM.YYYY format
func parseDateString(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, ".")
	if len(parts) != 3 {
		return time.Time{}, &TempooError{Message: fmt.Sprintf("Invalid date format '%s'. Expected DD.MM.YYYY", dateStr)}
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, &TempooError{Message: fmt.Sprintf("Invalid day in date '%s'", dateStr)}
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, &TempooError{Message: fmt.Sprintf("Invalid month in date '%s'", dateStr)}
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, &TempooError{Message: fmt.Sprintf("Invalid year in date '%s'", dateStr)}
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date, nil
}
