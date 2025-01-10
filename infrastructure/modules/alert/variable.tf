variable "slack_notification_channel_name" {
  type        = string
  description = "Existing slack notification channel name"
  default     = "Google Cloud Monitoring"
}

variable "environment" {
  type        = string
  description = "Environment name"
  default     = "prod"
  validation {
    condition     = contains(["prod", "staging"], var.environment)
    error_message = "Environment name must be either 'prod' or 'staging'."
  }
}

variable "prefix" {
  type        = string
  description = "Prefix to be added to alerts name"
  default     = "[Terraform]"
}
