package config

type Config struct {
	ProjectID                     string
	DripEnv                       string
	SlackRegistryChannelWebhook   string
	JWTSecret                     string
	DiscordSecurityChannelWebhook string
	SecretScannerURL              string
}
