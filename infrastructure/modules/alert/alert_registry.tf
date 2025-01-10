
resource "google_monitoring_alert_policy" "registry_security_scan_high_latency" {
  display_name = "${local.prefix}Security Scan High Latency (${var.environment})"

  conditions {
    display_name = "Global - custom/comfy_api_frontend/request_duration"

    condition_threshold {
      filter = <<EOT
        resource.type = "global" 
        AND metric.type = "custom.googleapis.com/comfy_api_frontend/request_duration" 
        AND (
          metric.labels.endpoint = "/security-scan" 
          AND metric.labels.env = "${var.environment}"
          AND metric.labels.method = "GET"
        )
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 300

      aggregations {
        alignment_period     = "3600s"
        per_series_aligner   = "ALIGN_MEAN"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.endpoint",
          "metric.label.statusCode",
          "metric.label.method"
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

resource "google_monitoring_alert_policy" "registry_security_scan_error" {
  display_name = "${local.prefix}Security Scan Error (${var.environment})"

  conditions {
    display_name = "Security Scan Error (5xx)"

    condition_threshold {
      filter = <<EOT
        resource.type = "global" 
        AND metric.type = "custom.googleapis.com/comfy_api_frontend/request_errors" 
        AND (
          metric.labels.endpoint = "/security-scan" 
          AND metric.labels.env = "${var.environment}"
          AND metric.labels.method = "GET" 
          AND metric.labels.statusCode = monitoring.regex.full_match("5.*")
        )
      EOT

      comparison = "COMPARISON_GT"
      duration   = "0s"

      aggregations {
        alignment_period     = "3600s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.endpoint",
          "metric.label.statusCode",
          "metric.label.method"
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

resource "google_monitoring_alert_policy" "registry_security_scan_abnormal_request_count" {
  display_name = "${local.prefix}Security Scan Abnormal Request Count (${var.environment})"

  conditions {
    display_name = "Global - custom/comfy_api_frontend/request_count"

    condition_threshold {
      filter = <<EOT
        resource.type = "global" 
        AND metric.type = "custom.googleapis.com/comfy_api_frontend/request_count" 
        AND (
          metric.labels.endpoint = "/security-scan" 
          AND metric.labels.env = "${var.environment}" 
          AND metric.labels.method = "GET"
        )
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 3

      aggregations {
        alignment_period     = "3600s"
        per_series_aligner   = "ALIGN_DELTA"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.endpoint",
          "metric.label.method",
          "metric.label.statusCode"
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

resource "google_monitoring_alert_policy" "registry_metrics_5xx" {
  display_name = "${local.prefix}Registry Metrics 5xx Server Side Error (${var.environment})"

  conditions {
    display_name = "Request Server Errors (5xx)"

    condition_threshold {
      filter = <<EOT
        resource.type = "global" 
        AND metric.type = "custom.googleapis.com/comfy_api_frontend/request_errors" 
        AND (
          metric.labels.endpoint = monitoring.regex.full_match("/nodes.*|/publishers.*") 
          AND metric.labels.env = "${var.environment}" 
          AND metric.labels.statusCode = monitoring.regex.full_match("5.*")
        )
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 0

      aggregations {
        alignment_period     = "60s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.endpoint",
          "metric.label.statusCode",
          "metric.label.method"
        ]
      }

      trigger {
        count = 1
      }
    }
  }
  combiner = "OR"
  alert_strategy {
    auto_close           = "3600s"
    notification_prompts = ["OPENED"]
  }
  severity = "ERROR"

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}

resource "google_monitoring_alert_policy" "registry_node_reindex_error" {
  display_name = "${local.prefix}Node Reindex Error (${var.environment})"

  conditions {
    display_name = "Nodes Reindex Error (5xx)"

    condition_threshold {
      filter = <<EOT
        resource.type = "global" 
        AND metric.type = "custom.googleapis.com/comfy_api_frontend/request_errors" 
        AND (
          metric.labels.endpoint = "/nodes/reindex" 
          AND metric.labels.method = "POST" 
          AND metric.labels.env = "${var.environment}" 
          AND metric.labels.statusCode = monitoring.regex.full_match("5.*")
        )
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 0

      aggregations {
        alignment_period     = "86400s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.endpoint",
          "metric.label.statusCode",
          "metric.label.method"
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

resource "google_monitoring_alert_policy" "registry_node_reindex_high_latency" {
  display_name = "${local.prefix}Node Reindex High Latency (${var.environment})"

  conditions {
    display_name = "Global - custom/comfy_api_frontend/request_duration"

    condition_threshold {
      filter = <<EOT
        resource.type = "global" 
        AND metric.type = "custom.googleapis.com/comfy_api_frontend/request_duration" 
        AND (
          metric.labels.endpoint = "/nodes/reindex" 
          AND metric.labels.method = "POST"
          AND metric.labels.env = "${var.environment}" 
        )
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 60

      aggregations {
        alignment_period     = "86400s"
        per_series_aligner   = "ALIGN_MEAN"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.endpoint",
          "metric.label.statusCode",
          "metric.label.method"
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

resource "google_monitoring_alert_policy" "registry_node_reindex_abnormal_request_count" {
  display_name = "${local.prefix}Node Reindex Abnormal Request Count (${var.environment})"

  conditions {
    display_name = "Global - custom/comfy_api_frontend/request_count"

    condition_threshold {
      filter = <<EOT
        resource.type = "global" 
        AND metric.type = "custom.googleapis.com/comfy_api_frontend/request_count" 
        AND (
            metric.labels.endpoint = "/nodes/reindex" 
            AND metric.labels.method = "POST"
            AND metric.labels.env = "${var.environment}" 
        )
      EOT

      comparison      = "COMPARISON_GT"
      duration        = "0s"
      threshold_value = 3

      aggregations {
        alignment_period     = "86400s"
        per_series_aligner   = "ALIGN_DELTA"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields = [
          "metric.label.endpoint",
          "metric.label.statusCode",
          "metric.label.method"
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
