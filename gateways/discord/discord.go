package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"registry-backend/config"
)

type DiscordService interface {
	SendSecurityCouncilMessage(msg string) error
}

type DripDiscordService struct {
	securityDiscordChannelWebhook string
	config                        *config.Config
}

func NewDiscordService(config *config.Config) *DripDiscordService {
	return &DripDiscordService{
		config:                        config,
		securityDiscordChannelWebhook: config.DiscordSecurityChannelWebhook,
	}
}

type discordRequestBody struct {
	Content string `json:"content"`
}

func (s *DripDiscordService) SendSecurityCouncilMessage(msg string) error {
	if s.config.DripEnv == "prod" {
		return sendDiscordNotification(msg, s.securityDiscordChannelWebhook)
	} else {
		println("Skipping sending message to Discord in non-prod environment. " + msg)
	}
	return nil
}

func sendDiscordNotification(msg string, discordWebhookURL string) error {
	if discordWebhookURL == "" {
		return fmt.Errorf("no Discord webhook URL provided, skipping sending message to Discord")
	}

	body, err := json.Marshal(discordRequestBody{Content: msg})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, discordWebhookURL, bytes.NewBuffer(body))
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
		return fmt.Errorf("request to Discord returned error status: %d", resp.StatusCode)
	}

	return nil
}
