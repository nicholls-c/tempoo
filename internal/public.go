package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/apex/log"
)

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

// ListWorklogs lists all worklogs for a given issue key for the current user
func (t *Tempoo) ListWorklogs(issueKey string) error {
	log.Infof("Listing worklogs for %s", issueKey)

	// validate the issue key
	if err := t.validateIssueKey(issueKey); err != nil {
		return err
	}

	// get current user's account ID
	userID, err := t.GetUserAccountID()
	if err != nil {
		return fmt.Errorf("failed to get user account ID: %w", err)
	}
	log.Debugf("User ID: %s", userID)

	// get the worklogs for the issue
	resp, err := t.client.R().Get(fmt.Sprintf("%s/issue/%s/worklog", JiraAPIRootURL, issueKey))
	if err != nil {
		log.Errorf("Request failed: %v", err)
		return &TempooError{Message: "API request failed", Cause: err}
	}

	// check if the request was successful
	if resp.StatusCode() != 200 {
		return &TempooError{Message: fmt.Sprintf("Failed to list worklogs: %s", resp.Status())}
	}

	// responseData is a map[string]interface{} type
	var responseData JiraResponse
	// unmarshal the response body into responseData
	if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
		return &TempooError{Message: "Failed to parse worklog data", Cause: err}
	}

	// parse responseData as a list of worklogs
	worklogsInterface, ok := responseData["worklogs"]
	if !ok {
		return &TempooError{Message: "Worklogs not found in response data"}
	}

	worklogs, ok := worklogsInterface.([]interface{})
	if !ok {
		return &TempooError{Message: "Worklogs not found in response data"}
	}

	// filter worklogs for the current user
	var userWorklogs []interface{}
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

		// this worklog belongs to the current user
		userWorklogs = append(userWorklogs, worklogInterface)
	}

	if len(userWorklogs) == 0 {
		log.Infof("No worklogs found for issue %s for current user", issueKey)
		return nil
	} else {
		log.Infof("Found %d worklog(s) for issue %s for current user:", len(userWorklogs), issueKey)
		t.printWorklogs(userWorklogs)
	}

	return nil
}
