package main

import (
	"fmt"
	"os"
	"tempoo/internal"

	"github.com/alecthomas/kong"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/willabides/kongplete"
)

// version will be set by goreleaser at build time
var version = "0.0.1-dev"

// Global factory instance
var tempooFactory *internal.TempooFactory

// VersionCmd represents the version command
type VersionCmd struct{}

// Run executes the version command
func (cmd *VersionCmd) Run() error {
	fmt.Println(version)
	return nil
}

// AddWorklogCmd represents the add worklog command
type AddWorklogCmd struct {
	IssueKey string  `help:"Jira issue key (e.g., PROJ-123)" short:"i"`
	Time     string  `help:"Time to log (e.g., 1h 30m, 2h, 45m)" short:"t"`
	Date     *string `help:"Date for the worklog in DD.MM.YYYY format (defaults to today)" short:"D"`
}

// getFactory initializes and returns the tempoo factory
func getFactory() (*internal.TempooFactory, error) {
	if tempooFactory == nil {
		var err error
		tempooFactory, err = internal.NewTempooFactory()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Tempoo factory: %w", err)
		}
	}
	return tempooFactory, nil
}

// Run executes the add worklog command
func (cmd *AddWorklogCmd) Run(ctx *kong.Context) error {
	// Check if required parameters are provided
	if cmd.IssueKey == "" || cmd.Time == "" {
		fmt.Fprintf(ctx.Stderr, "Usage: %s\n", ctx.Command())
		ctx.PrintUsage(false)
		return nil
	}

	factory, err := getFactory()
	if err != nil {
		return err
	}
	tempoo := factory.GetClient()
	return tempoo.AddWorklog(cmd.IssueKey, cmd.Time, cmd.Date)
}

// RemoveWorklogCmd represents the remove worklog command
type RemoveWorklogCmd struct {
	IssueKey string `help:"Jira issue key (e.g., PROJ-123)" short:"i"`
}

// Run executes the remove worklog command
func (cmd *RemoveWorklogCmd) Run(ctx *kong.Context) error {
	// Check if required parameters are provided
	if cmd.IssueKey == "" {
		fmt.Fprintf(ctx.Stderr, "Usage: %s\n", ctx.Command())
		ctx.PrintUsage(false)
		return nil
	}

	factory, err := getFactory()
	if err != nil {
		return err
	}
	tempoo := factory.GetClient()

	// get current user's account ID
	userID, err := tempoo.GetUserAccountID()
	if err != nil {
		return fmt.Errorf("failed to get user account ID: %w", err)
	}
	log.Debugf("User ID: %s", userID)

	// get worklogs for the user
	worklogIDs, err := tempoo.GetWorklogs(cmd.IssueKey, userID)
	if err != nil {
		return fmt.Errorf("failed to get worklogs: %w", err)
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
			return fmt.Errorf("failed to delete worklog %s: %w", worklogID, err)
		}
		log.Infof("Worklog ID %s deleted", worklogID)
	}

	return nil
}

// ListWorklogsCmd represents the list worklogs command
type ListWorklogsCmd struct {
	IssueKey string `help:"Jira issue key (e.g., PROJ-123)" short:"i"`
}

// Run executes the list worklogs command
func (cmd *ListWorklogsCmd) Run(ctx *kong.Context) error {
	// Check if required parameters are provided
	if cmd.IssueKey == "" {
		fmt.Fprintf(ctx.Stderr, "Usage: %s\n", ctx.Command())
		ctx.PrintUsage(false)
		return nil
	}

	factory, err := getFactory()
	if err != nil {
		return err
	}
	tempoo := factory.GetClient()
	return tempoo.ListWorklogs(cmd.IssueKey)
}

// Kong CLI struct
var CLI struct {
	AddWorklog    AddWorklogCmd    `cmd:"add-worklog" help:"Add a worklog to a Jira issue"`
	RemoveWorklog RemoveWorklogCmd `cmd:"remove-worklog" help:"Remove all user worklogs from a Jira issue"`
	ListWorklogs  ListWorklogsCmd  `cmd:"list-worklogs" help:"List all worklogs for a Jira issue"`
	Version       VersionCmd       `cmd:"version" help:"Print the version of the CLI"`

	// Add the completion installation command
	InstallCompletions kongplete.InstallCompletions `cmd:"install-completions" help:"Install shell completions"`

	Verbose bool `help:"Enable debug logging"`
}

// main function
func main() {
	// print help by default unless -v flag is set
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "--help")
	}

	// set up kong CLI
	parser := kong.Must(&CLI,
		kong.Name("tempoo"),
		kong.Description("tem💩, because life is too short.\n\nA CLI tool for managing Jira worklogs."),
		kong.UsageOnError(),
	)

	// Add completion support
	kongplete.Complete(parser)

	// Parse the arguments
	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		parser.FatalIfErrorf(err)
	}

	// set up apex/log with CLI handler
	log.SetHandler(cli.New(os.Stderr))
	if CLI.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// execute kong
	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
