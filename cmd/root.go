package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "slack-mastodon",
	Short:         "Mirror Mastodon posts to Slack channel",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute function is the entrypoint for the CLI
func Execute() error {
	return rootCmd.Execute()
}
