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
	Default

	Hook string
}

func (s Slack) Notify(message string) error {
	fmt.Println("Sending notification")

	slackBody, _ := json.Marshal(slackRequestBody{Text: message})
	req, err := http.NewRequest(http.MethodPost, s.Hook, bytes.NewBuffer(slackBody))
	if err != nil {
		s.FallbackNotify(message)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		s.FallbackNotify(message)
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		s.FallbackNotify(message)
		return errors.New("Unable to send notification to slack")
	}

	return nil
}