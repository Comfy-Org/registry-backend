# get the existing GCS bucket
data "google_storage_bucket" "bucket" {
  name = var.bucket_name
}

# create a Pub/Sub topic
resource "google_pubsub_topic" "topic" {
  name = var.topic_name
}

# get the default GCS service account
data "google_storage_project_service_account" "gcs_account" {
}

# Grant the GCS service account permission to publish to the Pub/Sub topic
resource "google_pubsub_topic_iam_binding" "gcs_publisher" {
  topic   = google_pubsub_topic.topic.name
  role    = "roles/pubsub.publisher"
  members = ["serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"]
}

# enable GCS Bucket Notification to Pub/Sub
resource "google_storage_notification" "notification" {
  bucket         = data.google_storage_bucket.bucket.name
  topic          = google_pubsub_topic.topic.id
  payload_format = "JSON_API_V1"
  depends_on     = [google_pubsub_topic_iam_binding.gcs_publisher]
  event_types = [
    "OBJECT_FINALIZE", # Triggered when an object is successfully created or overwritten
  ]
}


# Get the existing cloudbuild service account
data "google_service_account" "cloudbuild_service_account" {
  account_id = var.cloud_build_service_account
}

resource "google_project_iam_member" "logs_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${data.google_service_account.cloudbuild_service_account.email}"
}

resource "google_project_iam_member" "token_creator" {
  project = var.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:${data.google_service_account.cloudbuild_service_account.email}"
}

# Create the cloud build trigger
resource "google_cloudbuild_trigger" "trigger" {
  name            = var.trigger_name
  location        = var.region
  service_account = data.google_service_account.cloudbuild_service_account.id

  pubsub_config {
    topic = google_pubsub_topic.topic.id
  }

  source_to_build {
    uri       = var.git_repo_uri
    ref       = "refs/heads/${var.git_repo_branch}"
    repo_type = "GITHUB"
  }

  git_file_source {
    uri       = var.git_repo_uri
    revision  = "refs/heads/${var.git_repo_branch}"
    repo_type = "GITHUB"
    path      = "node-pack-extract/cloudbuild.yaml"
  }

  substitutions = {
    _CUSTOM_NODE_NAME     = "custom-node"
    _CUSTOM_NODE_URL      = "https://storage.googleapis.com/$(body.message.data.bucket)/$(body.message.data.name)"
    _REGISTRY_BACKEND_URL = var.registry_backend_url
  }
}
