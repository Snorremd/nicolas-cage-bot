# nicolas-cage-bot
A Slack bot written in golang for all things Nicolas Cage.

The bot currently provides recommendations for movies starring Nicolas Cage
available on Netflix. It uses the
[NetflixRoulette](https://netflixroulette.net/api/) to achieve this.

The implementation uses the [nlopes/slack](https://github.com/nlopes/slack)
go library to subscribe to the [Slack rtm api](https://api.slack.com/rtm) for
real time messaging.

## Running the bot

Store your slack token in a file to avoid having it end up in your sh history:

```bash
$ vim ~/slacktoken
export SLACK_TOKEN=xoxb-<somevalue>
$ . ~/slacktoken
```

Then run it with Docker:

```bash
docker run -e SLACK_TOKEN=$SLACK_TOKEN snorremd/nicolas-cage-bot:latest
```

Or compile and run it yourself:

```bash
go build main.go
SLACK_TOKEN=$SLACK_TOKEN ./main
```

You can also use docker-compose:

```yaml
version: "3"

services:
  cagebot:
    image: snorremd/nicolas-cage-bot:latest
    restart: unless-stopped
    environment:
      SLACK_TOKEN: "<yourtoken>"
```