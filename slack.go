package checkup

import (
	"log"

	slack "github.com/ashwanthkumar/slack-go-webhook"
)

// Slack consist of all the sub components required to use Slack API
type Slack struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Webhook  string `json:"webhook"`
}

// Notify implements notifier interface
func (s Slack) Notify(results []Result) error {
	attach := []slack.Attachment{}
	for _, result := range results {
		if !result.Healthy {

			attach = append(attach, FormatAttachments(result))
		}
	}
	if len(attach) > 0 {
		s.Send(attach)
	}
	return nil
}

func FormatAttachments(result Result) (slack.Attachment)  {
	color := "danger"
	attach := slack.Attachment{}
	attach.AddField(slack.Field{Title: "Endpoint", Value: result.Title})
	attach.AddField(slack.Field{Title: "Status", Value: result.Times[0].Status})
	attach.Color = &color

	return attach
}

func (s Slack) Send(attachments []slack.Attachment) error {
	payload := slack.Payload{
		Text:        "'Oh bother' sighed Pooh, 'not again.'",
		Username:    s.Username,
		Channel:     s.Channel,
		Attachments: attachments,
	}

	err := slack.Send(s.Webhook, "", payload)
	if len(err) > 0 {
		log.Printf("ERROR: %s", err)
	}
	return nil
}

// Send request via Slack API to create incident
// func (s Slack) Send(result Result) error {
// 	color := "danger"
// 	attach := slack.Attachment{}
// 	attach.AddField(slack.Field{Title: "Endpoint", Value: result.Title})
// 	attach.AddField(slack.Field{Title: "Status", Value: strings.ToUpper(fmt.Sprint(result.Status()))})
// 	attach.Color = &color
// 	payload := slack.Payload{
// 		Text:        "'Oh bother' sighed Pooh, 'not again.'",
// 		Username:    s.Username,
// 		Channel:     s.Channel,
// 		Attachments: []slack.Attachment{attach},
// 	}
//
// 	err := slack.Send(s.Webhook, "", payload)
// 	if len(err) > 0 {
// 		log.Printf("ERROR: %s", err)
// 	}
// 	log.Printf("Create request for %s", result.Endpoint)
// 	return nil
// }
