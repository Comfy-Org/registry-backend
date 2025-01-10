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
  bucket_name                 = "comfy-registry"
  cloud_build_service_account = "cloud-scheduler@dreamboothy.iam.gserviceaccount.com"
  topic_name                  = "comfy-registry-event"
  registry_backend_url        = "https://api.comfy.org"
  backfill_job_name           = "comfy-node-pack-backfill"
}

module "alert" {
  source = "../modules/alert"
  providers = {
    google = google
  }
}
