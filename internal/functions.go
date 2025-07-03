package internal

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
)

func (t *Tempoo) validateIssueKey(issueKey string) error {
	issueURL := fmt.Sprintf("%s/issue/%s", JiraAPIRootURL, issueKey)
	log.Debugf("Validating issue key: %s", issueURL)

	resp, err := t.client.R().Get(issueURL)
	if err != nil {
		log.Errorf("Request failed: %v", err)
		return &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() != 200 {
		return &InvalidIssueKeyError{IssueKey: issueKey}
	}
	log.Debugf("Validated issue key: %s", issueKey)

	return nil
}

func (t *Tempoo) GetUserAccountID() (string, error) {
	log.Info("Getting current user Atlassian account ID...")

	resp, err := t.client.R().Get(fmt.Sprintf("%s/myself", JiraAPIRootURL))
	if err != nil {
		log.Errorf("Request failed: %v", err)
		return "", &TempooError{Message: "API request failed", Cause: err}
	}
	log.Debugf("Response: %s", resp.StatusCode())

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
	log.Infof("Current user Atlassian account ID: %s", accountID)

	return accountID, nil
}

func (t *Tempoo) GetWorklogs(issueKey, userID string) ([]string, error) {
	log.Infof("Getting worklogs for %s", issueKey)

	if err := t.validateIssueKey(issueKey); err != nil {
		return nil, err
	}

	resp, err := t.client.R().Get(fmt.Sprintf("%s/issue/%s/worklog", JiraAPIRootURL, issueKey))
	if err != nil {
		log.Errorf("Request failed: %v", err)
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

		// check if this worklog belongs to the user
		author, ok := worklog["author"].(map[string]interface{})
		if !ok {
			continue
		}

		authorAccountID, ok := author["accountId"].(string)
		if !ok || authorAccountID != userID {
			continue
		}

		// get the worklog ID
		worklogID, ok := worklog["id"]
		if !ok {
			continue
		}

		// convert to string (it might be a number)
		worklogIDStr := fmt.Sprintf("%v", worklogID)
		worklogsForUser = append(worklogsForUser, worklogIDStr)
	}

	log.Debugf("Found %d worklogs for user %s", len(worklogsForUser), userID)
	return worklogsForUser, nil
}

func (t *Tempoo) AddWorklog(issueKey, worklogTime string, dateStr *string) error {
	log.Infof("Adding worklog to %s", issueKey)

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

	// create ISO timestamp for 08:30 AM on the specified date
	started := time.Date(
		workDate.Year(), workDate.Month(), workDate.Day(),
		8, 30, 0, 751000000, // 08:30:00.751
		time.UTC,
	).Format("2006-01-02T15:04:05.000-0700")

	log.Debugf("Started timestamp: %s", started)

	payload := map[string]string{
		"timeSpent": worklogTime,
		"started":   started,
	}

	log.Debugf("Payload: %+v", payload)

	resp, err := t.client.R().
		SetBody(payload).
		Post(fmt.Sprintf("%s/issue/%s/worklog", JiraAPIRootURL, issueKey))

	if err != nil {
		log.Errorf("Request failed: %v", err)
		return &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() == 201 {
		log.Infof("Worklog added to %s", issueKey)
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

func (t *Tempoo) DeleteWorklog(issueKey, worklogID string) error {
	log.Debugf("Deleting worklog %s for %s", worklogID, issueKey)

	resp, err := t.client.R().Delete(fmt.Sprintf("%s/issue/%s/worklog/%s", JiraAPIRootURL, issueKey, worklogID))
	if err != nil {
		log.Errorf("Request failed: %v", err)
		return &TempooError{Message: "API request failed", Cause: err}
	}

	if resp.StatusCode() != 204 {
		return &TempooError{Message: fmt.Sprintf("Failed to delete worklog: %s", resp.Status())}
	}

	log.Infof("Deleted worklog %s for %s", worklogID, issueKey)
	return nil
}
