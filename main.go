package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	mention "github.com/woogea/slack-memobot/pkg"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var mentions = map[string]*mention.Mention{}
var openais = map[string]*mention.OpenAI{}

func main() {
	err := godotenv.Load()
	if err != nil {
		// for debug environment. Usualy it will be set by environment variable.
		log.Println(".env is not found")
	}
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		//socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						// Storage must be provided in each channel.
						m, ok := mentions[ev.Channel]
						if !ok {
							mentions[ev.Channel] = mention.NewMention(ev.Channel)
							m = mentions[ev.Channel]
						}
						fmt.Printf("channel: %v\n ", ev.Channel)
						fmt.Printf("text: %v\n ", ev.Text)
						// remove mentions
						re := regexp.MustCompile(`<@.*?>`)
						body := strings.TrimSpace(re.ReplaceAllString(ev.Text, ""))
						strs := strings.Split(body, " ")
						//first is command, second and follows are parameters.
						resp, err := m.Exec(strs[0], strings.Join(strs[1:], " "))
						if err != nil {
							resp = err.Error()
						}
						//commands return some message including empty.
						_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText(resp, false))
						if err != nil {
							fmt.Printf("failed posting message: %v", err)
						}
					case *slackevents.MemberJoinedChannelEvent:
						fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
					case *slackevents.MessageEvent:
						fmt.Printf("message is  %q", ev.Text)
						// if includes atmark into message, break.
						if strings.Contains(ev.Text, "<@") {
							break
						}
						// if text is empty, break.
						if ev.Text == "" {
							break
						}
						if openais[ev.Channel] == nil {
							openais[ev.Channel] = mention.EnableOpenapi()
							r := "enable OpenAPI"
							_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText(r, false))
							if err != nil {
								fmt.Printf("failed posting message: %v\n", err)
							}
						}
						if openais[ev.Channel].CanUseOpenai() {
							//ignoe bot's message
							isbot := ev.BotID != ""
							r, err := openais[ev.Channel].CallOpenai(ev.Text, isbot)
							if err != nil {
								fmt.Printf("failed calling openai: %v", err)
							}
							if !isbot {
								_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText(r, false))
								if err != nil {
									fmt.Printf("failed posting message: %v\n", err)
								}
							}
						}
					}

				default:
					client.Debugf("unsupported Events API event received")
				}
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()
}
