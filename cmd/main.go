package main

import (
	"fmt"
	"os"
	"tempoo/internal"

	"github.com/alecthomas/kong"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// version will be set by goreleaser at build time
var version = "0.0.1-dev"

// VersionCmd represents the version command
type VersionCmd struct {
}

// Run executes the version command
func (cmd *VersionCmd) Run() error {
	fmt.Println(version)
	return nil
}

// AddWorklogCmd represents the add worklog command
type AddWorklogCmd struct {
	IssueKey string  `help:"Jira issue key (e.g., PROJ-123)" short:"i" required:""`
	Time     string  `help:"Time to log (e.g., 1h 30m, 2h, 45m)" short:"t" required:""`
	Date     *string `help:"Date for the worklog in DD.MM.YYYY format (defaults to today)" short:"D"`
}

// RemoveWorklogCmd represents the remove worklog command
type RemoveWorklogCmd struct {
	IssueKey string `help:"Jira issue key (e.g., PROJ-123)" short:"i" required:""`
}

// run executes the add worklog command
func (cmd *AddWorklogCmd) Run() error {
	tempoo, err := internal.NewTempoo()
	if err != nil {
		return err
	}

	return tempoo.AddWorklog(cmd.IssueKey, cmd.Time, cmd.Date)
}

// run executes the remove worklog command
func (cmd *RemoveWorklogCmd) Run() error {
	tempoo, err := internal.NewTempoo()
	if err != nil {
		return err
	}

	// get current user's account ID
	userID, err := tempoo.GetUserAccountID()
	if err != nil {
		return err
	}
	log.Debugf("User ID: %s", userID)

	// get worklogs for the user
	worklogIDs, err := tempoo.GetWorklogs(cmd.IssueKey, userID)
	if err != nil {
		return err
	}
	log.Debugf("Worklog IDs: %+v", worklogIDs)

	// check if there are worklogs to remove
	if len(worklogIDs) == 0 {
		log.Infof("No worklogs found for issue %s", cmd.IssueKey)
		return nil
	}

	// delete all worklogs for the user
	for _, worklogID := range worklogIDs {
		log.Debugf("Deleting worklog ID: %s", worklogID)
		if err := tempoo.DeleteWorklog(cmd.IssueKey, worklogID); err != nil {
			return err
		}
		log.Infof("Worklog ID %s deleted", worklogID)
	}

	return nil
}

// Kong CLI struct
var CLI struct {
	AddWorklog    AddWorklogCmd    `cmd:"add-worklog" help:"Add a worklog to a Jira issue"`
	RemoveWorklog RemoveWorklogCmd `cmd:"remove-worklog" help:"Remove all user worklogs from a Jira issue"`

	Verbose bool `help:"Enable debug logging" short:"v"`
	Version bool `help:"Show version" short:"V"`
}

// main function
func main() {
	// check for version flag before parsing and print
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Println(version)
			os.Exit(0)
		}
	}

	// print help by default unless -v flag is set
	if len(os.Args) == 1 && !CLI.Version {
		os.Args = append(os.Args, "--help")
	}

	// set up kong CLI
	ctx := kong.Parse(&CLI,
		kong.Name("tempoo"),
		kong.Description("tem💩, because life is too short.\n\nA CLI tool for managing Jira worklogs."),
		kong.UsageOnError(),
	)

	// set up apex/log with CLI handler
	log.SetHandler(cli.New(os.Stderr))
	if CLI.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// execute kong
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
