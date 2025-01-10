terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.14.1"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

module "node_pack_extract_trigger" {
  source = "../modules/node-pack-extract-trigger"
  providers = {
    google = google
  }
  project_id                  = var.project_id
  region                      = var.region
  bucket_name                 = "staging-comfy-registry"
  cloud_build_service_account = "cloud-scheduler@dreamboothy.iam.gserviceaccount.com"
  trigger_name                = "comfy-node-pack-extract-staging"
  topic_name                  = "comfy-registry-event-staging"
  backfill_job_name           = "comfy-node-pack-backfill-staging"
  registry_backend_url        = "https://stagingapi.comfy.org"
}

module "alert" {
  source = "../modules/alert"
  providers = {
    google = google
  }
  environment = "staging"
}
