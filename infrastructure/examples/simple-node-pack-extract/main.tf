terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.14.1"
    }
  }
}

variable "prefix" {
  type = string
}

variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "us-central1"
}

provider "google" {
  region  = var.region
  project = var.project_id
}

resource "google_storage_bucket" "bucket" {
  name     = "${var.prefix}-comfy-registry-bucket"
  location = var.region
}

resource "google_service_account" "service_account" {
  account_id = "${var.prefix}-comfy-registry-sa"
}

module "node_pack_extract_trigger" {
  depends_on = [google_service_account.service_account, google_storage_bucket.bucket]
  source     = "../../modules/node-pack-extract-trigger"
  providers = {
    google = google
  }
  project_id                  = var.project_id
  region                      = var.region
  bucket_name                 = google_storage_bucket.bucket.name
  cloud_build_service_account = google_service_account.service_account.email
  topic_name                  = "${var.prefix}-comfy-registry-event"
  trigger_name                = "${var.prefix}-comfy-registry-event"
  registry_backend_url        = "https://stagingapi.comfy.org"
}

output "trigger_id" {
  value = module.node_pack_extract_trigger.trigger_id
}
output "topic_id" {
  value = module.node_pack_extract_trigger.topic_id
}
output "bucket_notification_id" {
  value = module.node_pack_extract_trigger.bucket_notification_id
}
output "backfill_scheduler_id" {
  value = module.node_pack_extract_trigger.backfill_scheduler_id
}
output "bucket_name" {
  value = google_storage_bucket.bucket.name
}
output "service_account" {
  value = google_service_account.service_account.email
}
