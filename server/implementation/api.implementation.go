package implementation

import (
	"registry-backend/config"
	"registry-backend/ent"
	"registry-backend/gateways/algolia"
	"registry-backend/gateways/discord"
	gateway "registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	dripservices_comfyci "registry-backend/services/comfy_ci"
	dripservices "registry-backend/services/registry"
)

type DripStrictServerImplementation struct {
	Client          *ent.Client
	ComfyCIService  *dripservices_comfyci.ComfyCIService
	RegistryService *dripservices.RegistryService
}

func NewStrictServerImplementation(client *ent.Client, config *config.Config, storageService storage.StorageService, slackService gateway.SlackService, discordService discord.DiscordService, algolia algolia.AlgoliaService) *DripStrictServerImplementation {
	return &DripStrictServerImplementation{
		Client:          client,
		ComfyCIService:  dripservices_comfyci.NewComfyCIService(config),
		RegistryService: dripservices.NewRegistryService(storageService, slackService, discordService, algolia, config),
	}
}
