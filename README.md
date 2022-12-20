slack-mastodon
=========

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/akerl/slack-mastodon/build.yml?branch=main)](https://github.com/akerl/slack-mastodon/actions)
[![GitHub release](https://img.shields.io/github/release/akerl/slack-mastodon.svg)](https://github.com/akerl/slack-mastodon/releases)
[![License](https://img.shields.io/github/license/akerl/slack-mastodon)](https://github.com/akerl/slack-mastodon/blob/master/LICENSE)

slack-mastodon mirrors your Mastodon feed to a Slack channel

slack-mastodon is based on [pansapiens/masto2slack](https://github.com/pansapiens/masto2slack), which performs a similar function but for your own posts. I started from that repo, adjusted for my own build process, and tweaked the code to use the user's timeline instead of their own posts.

## Usage

### Setup

You'll need to:

* On your Mastodon server, create a 'new application' and get an access token (`access_token`) (you can do this under `https://<your_instance>/settings/applications/new` in the web interface). You'll need `read` access.
* Create a [Slack App](https://api.slack.com/apps?new_app=1), grant it the chat:write and link:write bot scopes, and install to your workspace.
* Create a configuration file, `~/.config/slack-mastodon/config.yml`, as below.

Example `~/.config/slack-mastodon/config.yml`:

```yaml
slack:
  bot_token: xoxb-BLAH
  channel: C04GDJLJASD
mastodon:
  server: https://mastodon.social
  access_token: TOKEN
```

## Running

Either create a cronjob (`crontab -e`) to call `slack-mastodon` every 5 mins (or longer):

```
*/5 * * * * /usr/local/bin/slack-mastodon update >/dev/null 2>&1
```
(this assumes you've copied `slack-mastodon` to `/usr/local/bin/slack-mastodon`)

Alternatively, run:

`watch -n 300 ./slack-mastodon update`

to check for new posts every 5 mins (300 seconds).

## Installation

```
go install github.com/akerl/slack-mastodon@latest
```

## License

slack-mastodon is released under the MIT License. See the bundled LICENSE file for details.

