package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type slackRequestBody struct {
	Text string `json:"text"`
}

type Slack struct {
	Hook string
}

func (s Slack) Notify(message string) error {
	fmt.Println("Sending notification")

	SendSlackNotification(s.Hook, message)
	return nil
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func SendSlackNotification(webhookURL string, msg string) error {

	slackBody, _ := json.Marshal(slackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
