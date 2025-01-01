package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"registry-backend/config"
)

// DiscordService defines the interface for interacting with Discord notifications.
type DiscordService interface {
	SendSecurityCouncilMessage(msg string, private bool) error
}

// Ensure discordService struct implements DiscordService interface
var _ DiscordService = (*discordService)(nil)

// discordService struct holds the configuration and webhook URLs.
type discordService struct {
	securityDiscordChannelWebhook        string
	securityDiscordPrivateChannelWebhook string
	config                               *config.Config
}

// NewDiscordService creates a new Discord service using the provided config or returns a no-op implementation if the config is missing.
func NewDiscordService(cfg *config.Config) DiscordService {
	if cfg == nil || cfg.DiscordSecurityChannelWebhook == "" || cfg.DiscordSecurityPrivateChannelWebhook == "" {
		log.Info().Msg("No Discord configuration found, using no-op implementation")
		return &discordNoop{}
	}

	return &discordService{
		config:                               cfg,
		securityDiscordChannelWebhook:        cfg.DiscordSecurityChannelWebhook,
		securityDiscordPrivateChannelWebhook: cfg.DiscordSecurityPrivateChannelWebhook,
	}
}

type discordRequestBody struct {
	Content string `json:"content"`
}

// SendSecurityCouncilMessage sends a message to the appropriate Discord channel.
func (s *discordService) SendSecurityCouncilMessage(msg string, private bool) error {
	webhookURL := s.securityDiscordChannelWebhook
	if private {
		webhookURL = s.securityDiscordPrivateChannelWebhook
	}

	return sendDiscordNotification(msg, webhookURL)
}

// sendDiscordNotification sends the message to the provided Discord webhook URL.
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
		return fmt.Errorf("request to Discord returned error status: %d", resp.StatusCode)
	}

	return nil
}

// discordNoop is a no-op implementation of the DiscordService interface.
// It does nothing and is used when no valid config is provided.
type discordNoop struct{}

// SendSecurityCouncilMessage is a no-op implementation that simply returns nil.
func (s *discordNoop) SendSecurityCouncilMessage(msg string, private bool) error {
	log.Info().Msgf("No-op: Skipping Discord message: %s (private: %v)", msg, private)
	return nil
}
