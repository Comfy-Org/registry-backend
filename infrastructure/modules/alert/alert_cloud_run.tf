resource "google_monitoring_alert_policy" "cloud_run_uptime_check_comfy_backend" {
  display_name = "${local.prefix}comfy-backend Cloud Run Uptime Check uptime failure (${var.environment})"

  conditions {
    display_name = "Failure of uptime check_id prod-comfy-backend-cloud-run-uptime-check-xDxQ_xAD38Y"

    condition_threshold {
      filter          = <<-EOT
        metric.type="monitoring.googleapis.com/uptime_check/check_passed" 
        AND metric.label.check_id="prod-comfy-backend-cloud-run-uptime-check-xDxQ_xAD38Y" 
        AND resource.type="cloud_run_revision"
      EOT
      comparison      = "COMPARISON_GT"
      duration        = "60s"
      threshold_value = 10

      aggregations {
        alignment_period     = "1200s"
        per_series_aligner   = "ALIGN_NEXT_OLDER"
        cross_series_reducer = "REDUCE_COUNT_FALSE"
        group_by_fields      = ["resource.label.*"]
      }

      trigger {
        count = 1
      }
    }
  }
  combiner = "OR"

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}

resource "google_monitoring_alert_policy" "cloud_run_uptime_check_uptime_check" {
  display_name = "${local.prefix}electron-updater Cloud Run Uptime Check uptime failure (${var.environment})"

  conditions {
    display_name = "Failure of uptime check_id electron-updater-cloud-run-uptime-check-_b2ucd4qPlM"

    condition_threshold {
      filter          = <<-EOT
        metric.type="monitoring.googleapis.com/uptime_check/check_passed" 
        AND metric.label.check_id="electron-updater-cloud-run-uptime-check-_b2ucd4qPlM" 
        AND resource.type="cloud_run_revision"
      EOT
      comparison      = "COMPARISON_GT"
      duration        = "60s"
      threshold_value = 1

      aggregations {
        alignment_period     = "1200s"
        per_series_aligner   = "ALIGN_NEXT_OLDER"
        cross_series_reducer = "REDUCE_COUNT_FALSE"
        group_by_fields      = ["resource.label.*"]
      }

      trigger {
        count = 1
      }
    }
  }
  combiner = "OR"

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}


resource "google_monitoring_alert_policy" "mixpanel_tracking_proxy_uptime_check" {
  display_name = "${local.prefix}mixpanel-tracking-proxy Cloud Run Uptime Check uptime failure (${var.environment})"

  conditions {
    display_name = "Failure of uptime check_id mixpanel-tracking-proxy-cloud-run-uptime-check-jaNJ8pRRdh4"

    condition_threshold {
      filter          = <<-EOT
        metric.type="monitoring.googleapis.com/uptime_check/check_passed" 
        AND metric.label.check_id="mixpanel-tracking-proxy-cloud-run-uptime-check-jaNJ8pRRdh4" 
        AND resource.type="cloud_run_revision"
      EOT
      comparison      = "COMPARISON_GT"
      duration        = "60s"
      threshold_value = 1

      aggregations {
        alignment_period     = "1200s"
        per_series_aligner   = "ALIGN_NEXT_OLDER"
        cross_series_reducer = "REDUCE_COUNT_FALSE"
        group_by_fields      = ["resource.label.*"]
      }

      trigger {
        count = 1
      }
    }
  }
  combiner = "OR"

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}
