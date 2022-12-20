package cmd

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getConfig() (*viper.Viper, error) {
	c := viper.New()
	c.SetConfigName("config")
	c.SetConfigType("yml")
	c.AddConfigPath("$HOME/.config/slack-mastodon")
	err := c.ReadInConfig()
	return c, err
}

func updateRunner(cmd *cobra.Command, _ []string) error {
	flags := cmd.Flags()

	noop, err = flags.GetBool("noop")
	if err != nil {
		return err
	}

	verbose, err = flags.GetBool("verbose")
	if err != nil {
		return err
	}

	c, err := getConfig()
	if err != nil {
		return err
	}

	slackClient := slack.New(c.GetString("slack.bot_token"))
	slackChannel := c.GetString("slack.channel")

	mastodonClient := mastodon.NewClient(&mastodon.Config{
		Server:      c.GetString("mastodon.server"),
		AccessToken: c.GetString("mastodon.access_token"),
	})

	lastPostedID := c.GetString("last_status_id")
	if verbose {
		fmt.Printf("last_status_id: %s\n", lastPostedID)
	}

	pg := mastodon.Pagination{
		SinceID: mastodon.ID(lastPostedID),
	}
	timeline, err := mastodonClient.GetTimelineHome(context.Background(), &pg)
	if err != nil {
		return err
	}

	for _, post := range timeline {
		if verbose {
			fmt.Println(post.URL)
		}
		if !noop {
			slackClient.PostMessage(
				slackChannel,
				slack.MsgOptionText(post.URL),
				slack.MsgOptionEnableLinkUnfurl(),
			)
			c.Set("last_status_id", timeline[i].ID)
			c.WriteConfig()
		}
	}

	// update recorded last_status_id
	config.WriteConfig()

	return nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Run an update based on the provided configuration",
	RunE:  updateRunner,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("noop", "n", false, "Don't actually post")
	updateCmd.Flags().BoolP("verbose", "v", false, "Verbose logging")
}
