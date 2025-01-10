data "google_monitoring_notification_channel" "slack" {
  type         = "slack"
  display_name = var.slack_notification_channel_name
}

locals {
  prefix  = var.prefix != "" ? join("", [var.prefix, " "]) : ""
  is_prod = var.environment == "prod"
}
