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
// type RemoveWorklogCmd struct {
// 	IssueKey  string `arg:"" help:"Jira issue key (e.g., PROJ-123)"`
// 	WorklogID string `arg:"" help:"ID of the worklog to remove"`
// }

var CLI struct {
	AddWorklog AddWorklogCmd `cmd:"add-worklog" help:"Add a worklog to a Jira issue"`
	//RemoveWorklog RemoveWorklogCmd `cmd:"" help:"Remove a worklog from a Jira issue"`

	Debug bool `help:"Enable debug logging" short:"d"`
}

// Run executes the add worklog command
func (cmd *AddWorklogCmd) Run() error {
	tempoo, err := internal.NewTempoo()
	if err != nil {
		return err
	}

	return tempoo.AddWorklog(cmd.IssueKey, cmd.Time, cmd.Date)
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("tempoo"),
		kong.Description("A CLI for the Tempoo API"),
		kong.UsageOnError(),
	)

	// Set up logging
	if CLI.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Execute the command
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
