package add

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/query"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Add adds a remotelink to an issue.`
	examples = `$ jira issue remotelink add

# Pass required parameters to skip prompt 
$ jira issue remotelink add ISSUE-1 "My useful link" http://example.com`
)

// NewCmdCommentAdd is a comment add command.
func NewCmdRemotelinkAdd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "add ISSUE-KEY TITLE URL",
		Short:   "Add a comment to an issue",
		Long:    helpText,
		Example: examples,
		Annotations: map[string]string{
			"help:args": "ISSUE-KEY\tIssue key of the source issue, eg: ISSUE-1\n" +
				"TITLE\t\tTitle of the link\n" +
				"URL\t\tUrl of the link",
		},
		Run: add,
	}

	cmd.Flags().Bool("web", false, "Open issue in web browser after adding comment")

	return &cmd
}

func add(cmd *cobra.Command, args []string) {
	params := parseArgsAndFlags(args, cmd.Flags())
	client := api.Client(jira.Config{Debug: params.debug})
	ac := addCmd{
		client:    client,
		linkTypes: nil,
		params:    params,
	}

	cmdutil.ExitIfError(ac.setIssueKey())
	cmdutil.ExitIfError(ac.setTitle())
	cmdutil.ExitIfError(ac.setUrl())

	err := func() error {
		s := cmdutil.Info("Adding link to issue")
		defer s.Stop()

		return client.AddIssueRemotelink(ac.params.issueKey, ac.params.title, ac.params.url)
	}()
	cmdutil.ExitIfError(err)

	server := viper.GetString("server")

	cmdutil.Success("Added remotelink \"%s\"", ac.params.title)
	fmt.Printf("%s/browse/%s\n", server, ac.params.issueKey)

	if web, _ := cmd.Flags().GetBool("web"); web {
		err := cmdutil.Navigate(server, ac.params.issueKey)
		cmdutil.ExitIfError(err)
	}

}

type addParams struct {
	issueKey string
	title    string
	url      string

	debug bool
}

func parseArgsAndFlags(args []string, flags query.FlagParser) *addParams {
	var issueKey, title, url string

	nargs := len(args)
	if nargs >= 1 {
		issueKey = cmdutil.GetJiraIssueKey(viper.GetString("project.key"), args[0])
	}
	if nargs >= 2 {
		title = args[1]
	}
	if nargs >= 3 {
		url = args[2]
	}

	debug, err := flags.GetBool("debug")
	cmdutil.ExitIfError(err)

	return &addParams{
		issueKey: issueKey,
		title:    title,
		url:      url,
		debug:    debug,
	}
}

type addCmd struct {
	client    *jira.Client
	linkTypes []*jira.IssueLinkType
	params    *addParams
}

func (ac *addCmd) setIssueKey() error {
	if ac.params.issueKey != "" {
		return nil
	}

	var ans string

	qs := &survey.Question{
		Name:     "issueKey",
		Prompt:   &survey.Input{Message: "Issue key"},
		Validate: survey.Required,
	}
	if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
		return err
	}
	ac.params.issueKey = cmdutil.GetJiraIssueKey(viper.GetString("project.key"), ans)

	return nil
}

func (ac *addCmd) setTitle() error {
	if ac.params.title != "" {
		return nil
	}

	var ans string

	qs := &survey.Question{
		Name:     "title",
		Prompt:   &survey.Input{Message: "Title"},
		Validate: survey.Required,
	}
	if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
		return err
	}
	ac.params.title = ans

	return nil
}

func (ac *addCmd) setUrl() error {
	if ac.params.url != "" {
		return nil
	}

	var ans string

	qs := &survey.Question{
		Name:     "url",
		Prompt:   &survey.Input{Message: "Url"},
		Validate: survey.Required,
	}
	if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
		return err
	}
	ac.params.url = ans

	return nil
}
