resource "google_monitoring_alert_policy" "cloud_function_security_scan_error" {
  display_name = "${local.prefix}Cloud Function - Security Scan Error (${var.environment})"

  conditions {
    display_name = "Cloud Function - Security Scan Error"

    condition_threshold {
      filter = <<EOT
        resource.type = "cloud_function" 
        AND resource.labels.function_name = "security_scan" 
        AND metric.type = "cloudfunctions.googleapis.com/function/execution_count" 
        AND metric.labels.status != "ok"
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 0.001

      aggregations {
        alignment_period     = "3600s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "resource.label.function_name"
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

resource "google_monitoring_alert_policy" "cloud_function_security_scan_high_latency" {
  display_name = "${local.prefix}Cloud Function - Security Scan High Latency (${var.environment})"

  conditions {
    display_name = "Cloud Function - Execution times"

    condition_threshold {
      filter = <<EOT
        resource.type = "cloud_function" 
        AND resource.labels.function_name = "security_scan" 
        AND metric.type = "cloudfunctions.googleapis.com/function/execution_times" 
        AND metric.labels.status != "ok"
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 90000000000

      aggregations {
        alignment_period     = "300s"
        per_series_aligner   = "ALIGN_PERCENTILE_95"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "resource.label.function_name"
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
  severity = "WARNING"

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}
