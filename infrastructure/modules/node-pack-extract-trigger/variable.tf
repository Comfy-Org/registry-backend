# REQUIRED VARIABLE
variable "project_id" {
  type        = string
  description = "google cloud project id"
}

variable "bucket_name" {
  type        = string
  description = "Existing public bucket that store the comfy node-packs."
}

variable "cloud_build_service_account" {
  type        = string
  description = "Existing service account used to run the cloud build and used to access registry backend, e.g. cloud-build@my-project.iam.gserviceaccount.com. Note that this service account needs to have 'Service Account Token Creator' role."
}

# OPTIONAL VARIABLE
variable "region" {
  type        = string
  description = "Google Cloud region"
  default     = "us-central1"
}

variable "topic_name" {
  type        = string
  description = "Google Cloudpub/sub topic to be created"
  default     = "comfy-registry-event"
}

variable "trigger_name" {
  type        = string
  description = "Cloud build trigger name"
  default     = "comfy-node-pack-extract"
}

variable "git_repo_uri" {
  type        = string
  description = "Connected git repo containing the cloud build pipeline. See https://cloud.google.com/build/docs/repositories"
  default     = "https://github.com/Comfy-Org/registry-backend"
}

variable "git_repo_branch" {
  type        = string
  description = "Git repo branch."
  default     = "main"
}

variable "registry_backend_url" {
  type        = string
  description = "The base url where registry backend can be reached"
  default     = "https://api.comfy.org"
}

