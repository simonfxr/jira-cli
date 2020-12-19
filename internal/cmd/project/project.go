package project

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
)

// NewCmdProject is a project command.
func NewCmdProject() *cobra.Command {
	return &cobra.Command{
		Use:     "project",
		Short:   "All accessible jira projects",
		Long:    "Project lists all jira projects that a user has access to.",
		Aliases: []string{"projects"},
		Run:     projects,
	}
}

func projects(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	resp, err := api.Client(debug).Project()
	cmdutil.ExitIfError(err)

	v := view.NewProject(resp)

	cmdutil.ExitIfError(v.Render())
}
