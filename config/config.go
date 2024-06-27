package config

type Config struct {
	ProjectID                     string
	DripEnv                       string
	SlackRegistryChannelWebhook   string
	DiscordSecurityChannelWebhook string
	JWTSecret                     string
	SecretScannerURL              string
}
