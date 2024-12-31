package config

type Config struct {
	ProjectID                     string
	DripEnv                       string
	SlackRegistryChannelWebhook   string
	JWTSecret                     string
	DiscordSecurityChannelWebhook string
	DiscordSecurityPrivateChannelWebhook string
	SecretScannerURL              string
	IDTokenAudience               string
	AlgoliaAppID                  string
	AlgoliaAPIKey                 string
	CloudStorageBucketName        string
}
