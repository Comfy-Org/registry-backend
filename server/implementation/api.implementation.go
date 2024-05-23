package implementation

import (
	"registry-backend/config"
	"registry-backend/ent"
	gateway "registry-backend/gateways/slack"
	"registry-backend/gateways/storage"
	dripservices_comfyci "registry-backend/services/comfy_ci"
	dripservices_registry "registry-backend/services/registry"
)

type DripStrictServerImplementation struct {
	Client          *ent.Client
	ComfyCIService  *dripservices_comfyci.ComfyCIService
	RegistryService *dripservices_registry.RegistryService
}

func NewStrictServerImplementation(client *ent.Client, config *config.Config, storageService storage.StorageService, slackService gateway.SlackService) *DripStrictServerImplementation {

	return &DripStrictServerImplementation{
		Client:          client,
		ComfyCIService:  dripservices_comfyci.NewComfyCIService(config),
		RegistryService: dripservices_registry.NewRegistryService(storageService, slackService),
	}
}
