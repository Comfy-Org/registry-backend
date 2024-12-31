# Trigger for node-pack-extract

Terraform modules to setup trigger for cloud build that will run [node-pack-extract](../../../node-pack-extract/)

## Requirements

- Google Cloud Account
- Existing Google Cloud Storage public bucket where the Registry backend store the comfy node packs.
- Existing Service Account that is whitelisted in [service_account_auth](../../../server/middleware/authentication/service_account_auth.go#65) middleware. This service account will be attached to `Service Account Token Creator` and `Logs Writer` Role inside the module.
- [Connected repositories](https://cloud.google.com/build/docs/repositories) contains the [node-pack-extract](../../../node-pack-extract/) folder
