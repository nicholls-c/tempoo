package main

import (
	"tempoo/internal"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

// AddWorklogCmd represents the add worklog command
type AddWorklogCmd struct {
	IssueKey string  `help:"Jira issue key (e.g., PROJ-123)" short:"i" required:""`
	Time     string  `help:"Time to log (e.g., 1h 30m, 2h, 45m)" short:"t" required:""`
	Date     *string `help:"Date for the worklog in DD.MM.YYYY format (defaults to today)" short:"d"`
}

// RemoveWorklogCmd represents the remove worklog command
type RemoveWorklogCmd struct {
	IssueKey string `help:"Jira issue key (e.g., PROJ-123)" short:"i" required:""`
}

// Run executes the add worklog command
func (cmd *AddWorklogCmd) Run() error {
	tempoo, err := internal.NewTempoo()
	if err != nil {
		return err
	}

	return tempoo.AddWorklog(cmd.IssueKey, cmd.Time, cmd.Date)
}

// Run executes the remove worklog command
func (cmd *RemoveWorklogCmd) Run() error {
	tempoo, err := internal.NewTempoo()
	if err != nil {
		return err
	}

	// Get current user's account ID
	userID, err := tempoo.GetUserAccountID()
	if err != nil {
		return err
	}
	logrus.Debugf("User ID: %s", userID)

	// Get worklogs for the user
	worklogIDs, err := tempoo.GetWorklogs(cmd.IssueKey, userID)
	if err != nil {
		return err
	}
	logrus.Debugf("Worklog IDs: %+v", worklogIDs)

	// check if there are worklogs to remove
	if len(worklogIDs) == 0 {
		logrus.Infof("No worklogs found for issue %s", cmd.IssueKey)
		return nil
	}

	// Delete all worklogs for the user
	for _, worklogID := range worklogIDs {
		logrus.Debugf("Deleting worklog ID: %s", worklogID)
		if err := tempoo.DeleteWorklog(cmd.IssueKey, worklogID); err != nil {
			return err
		}
		logrus.Infof("Worklog ID %s deleted", worklogID)
	}

	return nil
}

// Kong CLI struct
var CLI struct {
	AddWorklog    AddWorklogCmd    `cmd:"add-worklog" help:"Add a worklog to a Jira issue"`
	RemoveWorklog RemoveWorklogCmd `cmd:"remove-worklog" help:"Remove a worklog from a Jira issue"`

	Debug bool `help:"Enable debug logging" short:"d"`
}

// Main function
func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("tempoo"),
		kong.Description("A CLI for the Tempoo API"),
		kong.UsageOnError(),
	)

	// Set up logrus logging
	//TODO: replace with zap/apex log library
	if CLI.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetReportCaller(false)

	// Execute the command
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
