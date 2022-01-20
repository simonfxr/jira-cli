package remotelink

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/remotelink/add"
)

const helpText = `remotelink command helps you add weblinks to issues. See available commands below.`

// NewCmdRemotelink is a comment command.
func NewCmdRemotelink() *cobra.Command {
	cmd := cobra.Command{
		Use:     "remotelink",
		Short:   "Manage issue comments",
		Long:    helpText,
		Aliases: []string{"remotelinks", "weblink", "weblinks"},
		RunE:    remotelink,
	}

	cmd.AddCommand(add.NewCmdRemotelinkAdd())

	return &cmd
}

func remotelink(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
