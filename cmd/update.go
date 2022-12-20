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

func getStamp() (*viper.Viper, error) {
	s := viper.New()
	s.SetConfigName("stamp")
	s.SetConfigType("yml")
	s.AddConfigPath("$HOME/.config/slack-mastodon")
	s.SafeWriteConfig()
	err := s.ReadInConfig()
	return s, err
}

func updateRunner(cmd *cobra.Command, _ []string) error { //revive:disable-line cyclomatic
	flags := cmd.Flags()

	noop, err := flags.GetBool("noop")
	if err != nil {
		return err
	}

	verbose, err := flags.GetBool("verbose")
	if err != nil {
		return err
	}

	c, err := getConfig()
	if err != nil {
		return err
	}

	s, err := getStamp()
	if err != nil {
		return err
	}

	slackClient := slack.New(c.GetString("slack.bot_token"))
	slackChannel := c.GetString("slack.channel")

	mastodonClient := mastodon.NewClient(&mastodon.Config{
		Server:      c.GetString("mastodon.server"),
		AccessToken: c.GetString("mastodon.access_token"),
	})

	lastPostedID := s.GetString("last_status_id")
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

	for i := len(timeline) - 1; i >= 0; i-- {
		post := timeline[i]
		text := post.URL
		if text == "" {
			if post.Reblog == nil {
				return fmt.Errorf("empty content found: %s", post.ID)
			}
			text = fmt.Sprintf("%s (boosted by %s)", post.Reblog.URL, post.Account.Username)
		}
		if verbose {
			fmt.Println(text)
		}
		if !noop {
			_, _, err := slackClient.PostMessage(
				slackChannel,
				slack.MsgOptionText(text, false),
				slack.MsgOptionEnableLinkUnfurl(),
			)
			if err != nil {
				return err
			}
			s.Set("last_status_id", post.ID)

			err = s.WriteConfig()
			if err != nil {
				return err
			}
		}
	}
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
