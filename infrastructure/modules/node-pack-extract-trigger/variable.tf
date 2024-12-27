# REQUIRED VARIABLE
variable "bucket_name" {
  type        = string
  description = "existing bucket name"
}

variable "cloud_build_service_account" {
  type        = string
  description = "Existing service account used to run the cloud build and used to access registry backend, e.g. cloud-build@my-project.iam.gserviceaccount.com. Note that this service account needs to have 'Service Account Token Creator' role."
}

# OPTIONAL VARIABLE
variable "region" {
  type        = string
  default     = "us-central1"
  description = "google cloud region"
}

variable "topic_name" {
  type        = string
  description = "pub/sub topic to be created"
  default     = "comfy-registry-event"
}

variable "git_repo_uri" {
  type        = string
  description = "git repo containing the cloud build pipeline"
  default     = "https://github.com/Comfy-Org/registry-backend"
}

variable "git_repo_branch" {
  type        = string
  description = "git repo branch"
  default     = "master"
}

variable "registry_backend_url" {
  type        = string
  description = "the url where registry backend can be reached"
  default     = "https://api.comfy.org"
}

