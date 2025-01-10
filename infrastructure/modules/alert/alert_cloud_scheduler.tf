resource "google_monitoring_alert_policy" "cloud_scheduler_job_error" {
  count = local.is_prod ? 1 : 0

  display_name = "${local.prefix}Cloud Scheduler Job Error (${var.environment})"

  conditions {
    display_name = "Cloud Scheduler Job Error"

    condition_threshold {
      filter = <<EOT
        resource.type = "cloud_scheduler_job" AND metric.type = "logging.googleapis.com/user/cloud-scheduler-job-error"
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 0

      aggregations {
        alignment_period     = "300s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.jobName",
          "metric.label.status"
        ]
      }

      trigger {
        count = 1
      }
    }
  }
  combiner = "OR"
  alert_strategy {
    auto_close           = "1800s"
    notification_prompts = ["OPENED", "CLOSED"]
  }
  severity = "ERROR"

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}
