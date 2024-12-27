# Trigger for node-pack-extract

Terraform modules to setup trigger for cloud build that will run [node-pack-extract](../../../node-pack-extract/)

## Requirements

- Google Cloud Account
- Existing Google Cloud Storage bucket where the Registry backend store the comfy node packs.
- Existing Service Account with `Service Account Token Creator` Role plus the one that is whitelisted in [service_account_auth](../../../server/middleware/authentication/service_account_auth.go#65) middleware.
