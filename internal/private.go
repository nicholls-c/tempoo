package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
)

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

func (t *Tempoo) printWorklogs(worklogs []interface{}) {

	log.Infof("Total worklogs found: %d", len(worklogs))

	for i, worklogInterface := range worklogs {
		worklog, ok := worklogInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Use timeSpentSeconds for accurate time calculation
		timeSpentSeconds, ok := worklog["timeSpentSeconds"].(float64)
		var timeDisplay string
		if ok && timeSpentSeconds > 0 {
			hours := timeSpentSeconds / 3600
			if hours == float64(int(hours)) {
				// Whole hours
				timeDisplay = fmt.Sprintf("%.0fh", hours)
			} else {
				// Hours with minutes
				wholeHours := int(hours)
				minutes := int((hours - float64(wholeHours)) * 60)
				timeDisplay = fmt.Sprintf("%dh %dm", wholeHours, minutes)
			}
		} else {
			// Fallback to timeSpent string
			timeSpent, ok := worklog["timeSpent"].(string)
			if ok {
				timeDisplay = timeSpent
			} else {
				timeDisplay = "Unknown"
			}
		}

		// Extract other fields...
		worklogID := "Unknown"
		if id, ok := worklog["id"].(string); ok {
			worklogID = id
		} else if id, ok := worklog["id"].(float64); ok {
			worklogID = fmt.Sprintf("%.0f", id)
		}

		startedStr, ok := worklog["started"].(string)
		var dateStr string
		if ok {
			if startedTime, err := time.Parse("2006-01-02T15:04:05.000-0700", startedStr); err == nil {
				dateStr = startedTime.Format("02.01.2006")
			} else {
				dateStr = "Unknown"
			}
		} else {
			dateStr = "Unknown"
		}

		authorName := "Unknown"
		if author, ok := worklog["author"].(map[string]interface{}); ok {
			if displayName, ok := author["displayName"].(string); ok {
				authorName = displayName
			}
		}

		log.Infof("  %d. %s - %s (by %s) [ID: %s]", i+1, timeDisplay, dateStr, authorName, worklogID)
	}
}
