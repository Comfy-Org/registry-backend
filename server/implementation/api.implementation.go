package implementation

import (
	"registry-backend/config"
	"registry-backend/ent"
	"registry-backend/gateways/algolia"
	"registry-backend/gateways/discord"
	"registry-backend/gateways/pubsub"
	gateway "registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	dripservices_comfyci "registry-backend/services/comfy_ci"
	dripservices "registry-backend/services/registry"

	"github.com/mixpanel/mixpanel-go"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type DripStrictServerImplementation struct {
	Client          *ent.Client
	ComfyCIService  *dripservices_comfyci.ComfyCIService
	RegistryService *dripservices.RegistryService
	MixpanelService *mixpanel.ApiClient
	NewRelicApp     *newrelic.Application
}

func NewStrictServerImplementation(
	client *ent.Client,
	config *config.Config,
	storageService storage.StorageService,
	pubsubService pubsub.PubSubService,
	slackService gateway.SlackService,
	discordService discord.DiscordService,
	algolia algolia.AlgoliaService,
	newRelicApp *newrelic.Application) *DripStrictServerImplementation {
	return &DripStrictServerImplementation{
		Client:          client,
		ComfyCIService:  dripservices_comfyci.NewComfyCIService(config),
		RegistryService: dripservices.NewRegistryService(storageService, pubsubService, slackService, discordService, algolia, config, newRelicApp),
		MixpanelService: mixpanel.NewApiClient("f919d1b9da9a57482453c72ef7b16d88"),
		NewRelicApp:     newRelicApp,
	}
}
