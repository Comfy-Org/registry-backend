resource "google_monitoring_alert_policy" "cloud_deploy_release_render_failure" {
  count = local.is_prod ? 1 : 0

  display_name = "${local.prefix}Cloud Deploy - Release Render Failure (comfy-backend-api-pipeline pipeline) (${var.environment})"

  conditions {
    display_name = "Release in comfy-backend-api-pipeline pipeline has render failure"

    condition_matched_log {
      filter = <<EOT
        logName="projects/dreamboothy/logs/clouddeploy.googleapis.com%2Frelease_render" 
        AND resource.type="clouddeploy.googleapis.com/DeliveryPipeline" 
        AND resource.labels.pipeline_id="comfy-backend-api-pipeline" 
        AND resource.labels.location="us-central1" 
        AND jsonPayload.releaseRenderState="FAILED"
      EOT
    }
  }
  combiner = "OR"
  alert_strategy {
    notification_rate_limit {
      period = "300s"
    }
    auto_close = "604800s"
  }

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}



resource "google_monitoring_alert_policy" "cloud_deploy_rollout_failure" {
  count = local.is_prod ? 1 : 0

  display_name = "${local.prefix}Cloud Deploy - Rollout Failure (comfy-backend-api-pipeline pipeline) (${var.environment})"

  conditions {
    display_name = "Rollout in comfy-backend-api-pipeline pipeline has failure"

    condition_matched_log {
      filter = <<EOT
        logName="projects/dreamboothy/logs/clouddeploy.googleapis.com%2Frollout_update" 
        AND resource.type="clouddeploy.googleapis.com/DeliveryPipeline" 
        AND resource.labels.pipeline_id="comfy-backend-api-pipeline" 
        AND resource.labels.location="us-central1" 
        AND jsonPayload.rolloutUpdateType="FAILED"
      EOT
    }
  }
  combiner = "OR"
  alert_strategy {
    notification_rate_limit {
      period = "300s"
    }
    auto_close = "604800s"
  }

  notification_channels = [
    data.google_monitoring_notification_channel.slack.id
  ]
}
