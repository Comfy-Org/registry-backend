package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const registrySlackChannelWebhook = "https://hooks.slack.com/services/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"

type SlackService interface {
	SendRegistryMessageToSlack(msg string) error
}

type DripSlackService struct {
}

func NewSlackService() *DripSlackService {
	return &DripSlackService{}

}

type slackRequestBody struct {
	Text string `json:"text"`
}

func (s *DripSlackService) SendRegistryMessageToSlack(msg string) error {
	return sendSlackNotification(msg, registrySlackChannelWebhook)
}

func sendSlackNotification(msg string, slackWebhookURL string) error {
	body, err := json.Marshal(slackRequestBody{Text: msg})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, slackWebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// You can handle or log the HTTP error status code here
		return fmt.Errorf("request to Slack returned error status: %d", resp.StatusCode)
	}

	return nil
}
