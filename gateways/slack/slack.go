package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"registry-backend/config"
)

// SlackService defines the interface for interacting with Slack notifications.
type SlackService interface {
	SendRegistryMessageToSlack(msg string) error
}

// Ensure slackService struct implements SlackService interface
var _ SlackService = (*slackService)(nil)

// slackService struct holds the configuration and webhook URL.
type slackService struct {
	registrySlackChannelWebhook string
	config                      *config.Config
}

// NewSlackService creates a new Slack service using the provided config or returns a noop implementation if the config is missing.
func NewSlackService(cfg *config.Config) SlackService {
	if cfg == nil || cfg.SlackRegistryChannelWebhook == "" {
		// Return a noop implementation if config is nil or missing keys
		log.Info().Msg("No Slack configuration found, using noop implementation")
		return &slackNoop{}
	}

	return &slackService{
		config:                      cfg,
		registrySlackChannelWebhook: cfg.SlackRegistryChannelWebhook,
	}
}

type slackRequestBody struct {
	Text string `json:"text"`
}

// SendRegistryMessageToSlack sends a message to the registry Slack channel.
func (s *slackService) SendRegistryMessageToSlack(msg string) error {
	// Skip sending messages in non-production environments
	if s.config.DripEnv != "prod" {
		log.Info().Msg("Skipping Slack operations in non-prod environment")
		return nil
	}

	return sendSlackNotification(msg, s.registrySlackChannelWebhook)
}

// sendSlackNotification sends the message to the provided Slack webhook URL.
func sendSlackNotification(msg string, slackWebhookURL string) error {
	if slackWebhookURL == "" {
		return fmt.Errorf("no Slack webhook URL provided, skipping sending message to Slack")
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

// slackNoop is a noop implementation of the SlackService interface.
// It does nothing and is used when no valid config or the environment is non-production.
type slackNoop struct{}

// Implement all SlackService methods for noop behavior.

func (s *slackNoop) SendRegistryMessageToSlack(msg string) error {
	// No-op, just return nil to avoid any side-effects
	return nil
}
