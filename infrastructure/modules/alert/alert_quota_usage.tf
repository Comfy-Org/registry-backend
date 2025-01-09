
resource "google_monitoring_alert_policy" "quota_usage_compute_engine_api" {
  display_name = "${local.prefix}Quota usage - Compute Engine API - NVIDIA T4 GPUs - compute.googleapis.com/nvidia_t4_gpus (${var.environment})"

  conditions {
    display_name = "Quota usage reached defined threshold"

    condition_monitoring_query_language {
      duration = "60s"
      query    = <<EOT
        fetch consumer_quota
        | filter resource.project_id == 'dreamboothy'
        | {
            metric serviceruntime.googleapis.com/quota/allocation/usage
            | filter metric.quota_metric == 'compute.googleapis.com/nvidia_t4_gpus'
            && resource.location == 'us-central1'
            | map add [metric.limit_name: 'NVIDIA-T4-GPUS-per-project-region']
            | align next_older(1d)
            | group_by [resource.project_id, resource.service, metric.quota_metric, metric.limit_name, resource.location], .max;

            metric serviceruntime.googleapis.com/quota/limit
            | filter metric.quota_metric == 'compute.googleapis.com/nvidia_t4_gpus'
            && metric.limit_name == 'NVIDIA-T4-GPUS-per-project-region'
            && resource.location == 'us-central1'
            | align next_older(1d)
            | group_by [resource.project_id, resource.service, metric.quota_metric, metric.limit_name, resource.location], .min
        }
        | ratio
        | every 30s
        | condition gt(val(), 0.8 '1')
      EOT

      trigger {
        count = 1
      }
    }
  }
  combiner = "OR"
  alert_strategy {
    auto_close = "604800s"
  }
  severity = "ERROR"

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}

