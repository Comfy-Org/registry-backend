package config

type Config struct {
	ProjectID                     string
	DripEnv                       string
	SlackRegistryChannelWebhook   string
	JWTSecret                     string
	ReindexNodesCrontab           string
	DiscordSecurityChannelWebhook string
	SecretScannerURL              string
}
