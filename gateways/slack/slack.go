package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"registry-backend/config"
)

type SlackService interface {
	SendRegistryMessageToSlack(msg string) error
}

type DripSlackService struct {
	registrySlackChannelWebhook string
	config                      *config.Config
}

func NewSlackService(config *config.Config) *DripSlackService {
	return &DripSlackService{
		config:                      config,
		registrySlackChannelWebhook: config.SlackRegistryChannelWebhook,
	}

}

type slackRequestBody struct {
	Text string `json:"text"`
}

func (s *DripSlackService) SendRegistryMessageToSlack(msg string) error {
	if s.config.DripEnv == "prod" {
		return sendSlackNotification(msg, s.registrySlackChannelWebhook)
	}
	return nil
}

func sendSlackNotification(msg string, slackWebhookURL string) error {
	if slackWebhookURL == "" {
		println("No Slack webhook URL provided, skipping sending message to Slack")
		return nil
	}

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
