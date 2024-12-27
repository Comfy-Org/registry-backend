# Trigger for node-pack-extract staging

Terraform modules to setup trigger for cloud build that will run [node-pack-extract](../../../node-pack-extract/)

## Requirements

- Google Cloud Account

## Configuration

This use the following configuration value:

- bucket_name: "comfy-registry "
- service account: "<cloud-scheduler@dreamboothy.iam.gserviceaccount.com>"
- topic_name: "comfy-registry-event-staging"

## Apply

```bash
terraform apply 
    -var project_id=dreamboothy-dev
    -var region=us-central1
```
