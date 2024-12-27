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
  source = "../../node-pack-extract/trigger"
  providers = {
    google = google
  }
  region                      = var.region
  bucket_name                 = "comfy-registry"
  cloud_build_service_account = "cloud-scheduler@dreamboothy.iam.gserviceaccount.com"
  topic_name                  = "comfy-registry-event-stage"
}
